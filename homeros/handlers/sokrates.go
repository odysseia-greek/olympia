package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	plato "github.com/odysseia-greek/agora/plato/models"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/odysseia-greek/olympia/homeros/models"
	"time"
)

func (h *HomerosHandler) Methods(requestID string) ([]models.MethodTree, error) {
	cacheItem, _ := h.Cache.Read(cacheNameSokratesTree)
	if cacheItem != nil {
		cachedGraph, err := models.UnmarshalMethodGraph(cacheItem)
		if err != nil {
			return nil, err
		}

		traceId, parentspanId, traceCall := ParseHeaderID(requestID)
		if traceCall {
			traceResponse, err := json.Marshal(cachedGraph)
			if err != nil {
				return nil, err
			}

			span := &aristophanes.SpanRequest{
				TraceId:      traceId,
				ParentSpanId: parentspanId,
				Action:       "Cached",
				ResponseBody: string(traceResponse),
			}

			h.Tracer.Span(context.Background(), span)

			traceCloser := &aristophanes.CloseTraceRequest{
				TraceId:      traceId,
				ParentSpanId: parentspanId,
				ResponseCode: 200,
			}

			h.Tracer.CloseTrace(context.Background(), traceCloser)
		}

		logging.Info(fmt.Sprintf("taking from cache | traceID: %s | responseCode: %d", traceId, 200))
		return cachedGraph.MethodTree, nil
	}

	response, err := h.HttpClients.Sokrates().GetMethods(requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var methods plato.Methods
	err = json.NewDecoder(response.Body).Decode(&methods)
	if err != nil {
		return nil, err
	}

	var tree models.MethodGraph

	for _, method := range methods.Method {
		innerMethod := models.MethodTree{
			Name: method.Method,
		}

		resp, err := h.HttpClients.Sokrates().GetCategories(method.Method, requestID)
		if err != nil {
			h.CloseTrace(resp, nil)
			return nil, err
		}
		defer resp.Body.Close()

		var categories plato.Categories
		err = json.NewDecoder(resp.Body).Decode(&categories)
		if err != nil {
			return nil, err
		}

		for _, category := range categories.Category {

			categoryResp, err := h.HttpClients.Sokrates().GetChapters(method.Method, category.Category, requestID)
			if err != nil {
				h.CloseTrace(categoryResp, nil)
				return nil, err
			}
			defer resp.Body.Close()

			var highestChapter plato.LastChapterResponse
			err = json.NewDecoder(categoryResp.Body).Decode(&highestChapter)
			if err != nil {
				return nil, err
			}

			cat := models.Category{
				Name:           category.Category,
				HighestChapter: highestChapter.LastChapter,
			}
			innerMethod.Categories = append(innerMethod.Categories, cat)
		}

		tree.MethodTree = append(tree.MethodTree, innerMethod)
	}

	stringifiedMethodTree, _ := tree.Marshal()
	ttl := time.Hour
	err = h.Cache.SetWithTTL(cacheNameSokratesTree, string(stringifiedMethodTree), ttl)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, tree)

	return tree.MethodTree, nil
}

func (h *HomerosHandler) CreateQuestion(method, category, chapter, requestID string) (*plato.QuizResponse, error) {
	response, err := h.HttpClients.Sokrates().CreateQuestion(method, category, chapter, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var quiz plato.QuizResponse
	err = json.NewDecoder(response.Body).Decode(&quiz)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, quiz)

	return &quiz, nil
}

func (h *HomerosHandler) CheckQuestion(quiz, providedAnswer, requestID string) (*plato.CheckAnswerResponse, error) {
	check := plato.CheckAnswerRequest{
		QuizWord:       quiz,
		AnswerProvided: providedAnswer,
	}

	jsonCheck, err := check.Marshal()
	if err != nil {
		return nil, err
	}

	response, err := h.HttpClients.Sokrates().CheckAnswer(jsonCheck, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var answer plato.CheckAnswerResponse
	err = json.NewDecoder(response.Body).Decode(&answer)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, answer)

	return &answer, nil
}
