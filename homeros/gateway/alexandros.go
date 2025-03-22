package gateway

import (
	"encoding/json"
	"github.com/odysseia-greek/agora/plato/transform"
	"github.com/odysseia-greek/olympia/homeros/graph/model"

	"strconv"
)

func (h *HomerosHandler) Dictionary(word, language, mode, traceID string, searchInText bool) (*model.ExtendedDictionary, error) {
	searchInTextAsString := strconv.FormatBool(searchInText)
	formattedWord := transform.RemoveAccents(word)
	response, err := h.HttpClients.Alexandros().Search(formattedWord, language, mode, searchInTextAsString, traceID)
	if err != nil {
		h.CloseTraceWithError(err, traceID)
		return nil, err
	}

	var extendedResponse model.ExtendedDictionary
	err = json.NewDecoder(response.Body).Decode(&extendedResponse)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, extendedResponse)

	return &extendedResponse, nil
}
