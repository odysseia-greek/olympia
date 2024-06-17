package seeder

import (
	"github.com/odysseia-greek/agora/plato/config"
	"os"
)

func CreateNewConfig() (*ProtagorasHandler, error) {
	client, err := config.CreateOdysseiaClient()
	saveToDisk := os.Getenv("SAVE_TO_DISK") == "true"
	if err != nil {
		return nil, err
	}

	return &ProtagorasHandler{
		Client: client,
		Save:   saveToDisk,
	}, nil
}
