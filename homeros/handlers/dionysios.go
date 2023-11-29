package handlers

import (
	"encoding/json"
	plato "github.com/odysseia-greek/agora/plato/models"
)

func (h *HomerosHandler) Grammar(word, traceID string) ([]plato.Result, error) {
	response, err := h.HttpClients.Dionysios().Grammar(word, traceID)
	if err != nil {
		h.CloseTraceWithError(err, traceID)
		return nil, err
	}

	defer response.Body.Close()

	var results plato.DeclensionTranslationResults
	err = json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, results)

	return results.Results, nil
}
