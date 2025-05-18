package gateway

import (
	"encoding/json"
	"github.com/odysseia-greek/olympia/homeros/graph/model"
)

func (h *HomerosHandler) CreateText(body []byte, requestID string) (*model.Text, error) {
	response, err := h.HttpClients.Herodotos().Create(body, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}

	defer response.Body.Close()

	var sentence model.Text
	err = json.NewDecoder(response.Body).Decode(&sentence)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, sentence)

	return &sentence, nil
}

func (h *HomerosHandler) CheckText(body []byte, requestID string) (*model.CheckTextResponse, error) {
	response, err := h.HttpClients.Herodotos().Check(body, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var sentence model.CheckTextResponse
	err = json.NewDecoder(response.Body).Decode(&sentence)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, sentence)

	return &sentence, nil
}

func (h *HomerosHandler) HerodotosOptions(requestID string) (*model.AggregationResult, error) {
	response, err := h.HttpClients.Herodotos().Options(requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var aggregate model.AggregationResult
	err = json.NewDecoder(response.Body).Decode(&aggregate)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, aggregate)

	return &aggregate, nil
}

func (h *HomerosHandler) Analyze(body []byte, requestID string) (*model.AnalyzeTextResponse, error) {
	response, err := h.HttpClients.Herodotos().Analyze(body, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var analyzeResult model.AnalyzeTextResponse
	err = json.NewDecoder(response.Body).Decode(&analyzeResult)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, analyzeResult)

	return &analyzeResult, nil
}
