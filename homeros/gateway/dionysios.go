package gateway

import (
	"encoding/json"
	"github.com/odysseia-greek/olympia/homeros/graph/model"
)

func (h *HomerosHandler) Grammar(word, traceID string) (*model.DeclensionTranslationResult, error) {
	response, err := h.HttpClients.Dionysios().Grammar(word, traceID)
	if err != nil {
		h.CloseTraceWithError(err, traceID)
		return nil, err
	}

	defer response.Body.Close()

	var results model.DeclensionTranslationResult
	err = json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, results)

	return &results, nil
}
