package seeder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	dionysiosv1 "github.com/odysseia-greek/alexandreia/dionysios/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultDionysiosAddress = "dionysios.alexandreia.svc:50060"

func CreateNewConfig() (*ProtagorasHandler, error) {
	client, err := config.CreateOdysseiaClient()
	saveToDisk := os.Getenv("SAVE_TO_DISK") == "true"
	if err != nil {
		return nil, err
	}

	dionysiosAddress := config.StringFromEnv("DIONYSIOS_GRPC_ADDRESS", defaultDionysiosAddress)
	dionysiosConnection, err := grpc.NewClient(dionysiosAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("create Dionysios gRPC client: %w", err)
	}
	dionysios := dionysiosv1.NewDionysiosServiceClient(dionysiosConnection)

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
			dionysiosResponse, err := dionysios.Health(ctx, &dionysiosv1.HealthRequest{})
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

			logging.Debug(fmt.Sprintf("elapsed Time: %vs, Dionysios Healthy: %v, Herodotos Healthy: %v", elapsed, dionysiosResponse.GetHealthy(), healthHerodotos.Healthy))

			if dionysiosResponse.GetHealthy() && healthHerodotos.Healthy {
				return &ProtagorasHandler{
					Client:    client,
					Dionysios: dionysios,
					Save:      saveToDisk,
				}, nil
			}
		}
	}
}
