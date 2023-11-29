package handlers

import (
	"encoding/json"
	plato "github.com/odysseia-greek/agora/plato/models"
	"sync"
)

func (h *HomerosHandler) Dictionary(word, language, mode, traceID string) ([]plato.Meros, error) {
	response, err := h.HttpClients.Alexandros().Search(word, language, mode, traceID)
	if err != nil {
		h.CloseTraceWithError(err, traceID)
		return nil, err
	}

	var meroi []plato.Meros
	err = json.NewDecoder(response.Body).Decode(&meroi)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		h.CloseTrace(response, meroi)
	}()

	wg.Wait()

	return meroi, nil
}
