package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	plato "github.com/odysseia-greek/agora/plato/models"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/odysseia-greek/olympia/homeros/models"
	"time"
)

const (
	cacheNameHerodotosTree string = "herodotosTree"
)

func (h *HomerosHandler) Books(requestID string) ([]models.AuthorTree, error) {
	cacheItem, _ := h.Cache.Read(cacheNameHerodotosTree)
	if cacheItem != nil {
		cachedGraph, err := models.UnmarshalAuthorGraph(cacheItem)
		if err != nil {
			return nil, err
		}

		traceID, parentspanID, traceCall := ParseHeaderID(requestID)
		if traceCall {
			parabasis := &pb.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: parentspanID,
				SpanId:       parentspanID,
				RequestType: &pb.ParabasisRequest_CloseTrace{
					CloseTrace: &pb.CloseTraceRequest{
						ResponseCode: 200,
						ResponseBody: "taken from cache",
					},
				},
			}

			err := h.Streamer.Send(parabasis)
			if err != nil {
				logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
			}

			logging.Info(fmt.Sprintf("taking from cache | traceID: %s | responseCode: %d", traceID, 200))
		}

		return cachedGraph.AuthorTree, nil
	}

	var graph models.AuthorGraph

	response, err := h.HttpClients.Herodotos().GetAuthors(requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var authors plato.Authors
	err = json.NewDecoder(response.Body).Decode(&authors)
	if err != nil {
		return nil, err
	}

	for _, author := range authors.Authors {

		resp, err := h.HttpClients.Herodotos().GetBooks(author.Author, requestID)
		if err != nil {
			h.CloseTrace(resp, nil)
			return nil, err
		}

		var books plato.Books
		err = json.NewDecoder(resp.Body).Decode(&books)
		if err != nil {
			return nil, err
		}

		authorTree := models.AuthorTree{
			Name:  author.Author,
			Books: books.Books,
		}

		graph.AuthorTree = append(graph.AuthorTree, authorTree)
	}

	h.CloseTrace(response, graph)

	stringifiedAuthorTree, _ := graph.Marshal()
	ttl := time.Hour
	err = h.Cache.SetWithTTL(cacheNameHerodotosTree, string(stringifiedAuthorTree), ttl)
	if err != nil {
		return nil, err
	}

	return graph.AuthorTree, nil
}

func (h *HomerosHandler) Sentence(author, book, requestID string) (*models.SentenceGraph, error) {
	response, err := h.HttpClients.Herodotos().CreateQuestion(author, book, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}

	defer response.Body.Close()

	var sentence plato.CreateSentenceResponse
	err = json.NewDecoder(response.Body).Decode(&sentence)
	if err != nil {
		return nil, err
	}

	graph := models.SentenceGraph{
		Author: author,
		Book:   book,
		Greek:  sentence.Sentence,
		Id:     sentence.SentenceId,
	}

	h.CloseTrace(response, graph)

	return &graph, nil
}

func (h *HomerosHandler) Answer(id, author, answer, requestID string) (*models.Answer, error) {
	answerModel := plato.CheckSentenceRequest{
		SentenceId:       id,
		ProvidedSentence: answer,
		Author:           author,
	}
	response, err := h.HttpClients.Herodotos().CheckSentence(answerModel, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var sentence plato.CheckSentenceResponse
	err = json.NewDecoder(response.Body).Decode(&sentence)
	if err != nil {
		return nil, err
	}

	a := models.ParseSentenceResponseToAnswer(sentence)

	h.CloseTrace(response, sentence)

	return &a, nil
}
