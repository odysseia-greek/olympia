package seeder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"os"
	"time"
)

func CreateNewConfig() (*ProtagorasHandler, error) {
	client, err := config.CreateOdysseiaClient()
	saveToDisk := os.Getenv("SAVE_TO_DISK") == "true"
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout: one or both APIs are unhealthy")
		case <-ticker.C:
			elapsed := time.Since(startTime).Seconds()
			logging.Debug(fmt.Sprintf("%vs", elapsed))
			dionysiosResponse, err := client.Dionysios().Health("")
			if err != nil {
				continue
			}

			defer dionysiosResponse.Body.Close()

			var healthDionysios models.Health
			err = json.NewDecoder(dionysiosResponse.Body).Decode(&healthDionysios)
			if err != nil {
				continue
			}

			herodotosResponse, err := client.Herodotos().Health("")
			if err != nil {
				continue
			}

			defer herodotosResponse.Body.Close()
			var healthHerodotos models.Health
			err = json.NewDecoder(herodotosResponse.Body).Decode(&healthHerodotos)
			if err != nil {
				continue
			}

			logging.Debug(fmt.Sprintf("elapsed Time: %vs, Dionysios Healthy: %v, Herodotos Healthy: %v", elapsed, healthDionysios.Healthy, healthHerodotos.Healthy))

			if healthDionysios.Healthy && healthHerodotos.Healthy {
				return &ProtagorasHandler{
					Client: client,
					Save:   saveToDisk,
				}, nil
			}
		}
	}
}
