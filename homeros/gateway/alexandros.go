package gateway

import (
	"encoding/json"
	plato "github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/transform"
	"strconv"
)

func (h *HomerosHandler) Dictionary(word, language, mode, traceID string, searchInText bool) (*plato.ExtendedResponse, error) {
	searchInTextAsString := strconv.FormatBool(searchInText)
	formattedWord := transform.RemoveAccents(word)
	response, err := h.HttpClients.Alexandros().Search(formattedWord, language, mode, searchInTextAsString, traceID)
	if err != nil {
		h.CloseTraceWithError(err, traceID)
		return nil, err
	}

	var extendedResponse plato.ExtendedResponse
	err = json.NewDecoder(response.Body).Decode(&extendedResponse)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, extendedResponse)

	return &extendedResponse, nil
}
