package gateway

import (
	"encoding/json"
	plato "github.com/odysseia-greek/agora/plato/models"
)

func (h *HomerosHandler) CreateText(body []byte, requestID string) (*plato.Text, error) {
	response, err := h.HttpClients.Herodotos().Create(body, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}

	defer response.Body.Close()

	var sentence plato.Text
	err = json.NewDecoder(response.Body).Decode(&sentence)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, sentence)

	return &sentence, nil
}

func (h *HomerosHandler) CheckText(body []byte, requestID string) (*plato.CheckTextResponse, error) {
	response, err := h.HttpClients.Herodotos().Check(body, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var sentence plato.CheckTextResponse
	err = json.NewDecoder(response.Body).Decode(&sentence)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, sentence)

	return &sentence, nil
}

func (h *HomerosHandler) HerodotosOptions(requestID string) (*plato.AggregationResult, error) {
	response, err := h.HttpClients.Herodotos().Options(requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var aggregate plato.AggregationResult
	err = json.NewDecoder(response.Body).Decode(&aggregate)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, aggregate)

	return &aggregate, nil
}

func (h *HomerosHandler) Analyze(body []byte, requestID string) (*plato.AnalyzeTextResponse, error) {
	response, err := h.HttpClients.Herodotos().Analyze(body, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var analyzeResult plato.AnalyzeTextResponse
	err = json.NewDecoder(response.Body).Decode(&analyzeResult)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, analyzeResult)

	return &analyzeResult, nil
}
