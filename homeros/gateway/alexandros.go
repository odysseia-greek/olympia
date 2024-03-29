package gateway

import (
	"encoding/json"
	plato "github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/transform"
	"strconv"
	"sync"
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

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		h.CloseTrace(response, extendedResponse)
	}()

	wg.Wait()

	return &extendedResponse, nil
}
