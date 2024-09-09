package gateway

import (
	"encoding/json"
	plato "github.com/odysseia-greek/agora/plato/models"
)

func (h *HomerosHandler) Convert(body []byte, requestID string) (*plato.EdgecaseResponse, error) {
	response, err := h.HttpClients.Diogenes().Convert(body, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var edgecaseResponse plato.EdgecaseResponse
	err = json.NewDecoder(response.Body).Decode(&edgecaseResponse)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, edgecaseResponse)

	return &edgecaseResponse, nil
}
