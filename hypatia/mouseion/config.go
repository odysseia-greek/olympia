package mouseion

import (
	"context"
	"os"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
)

func CreateNewConfig(ctx context.Context) (*HypatiaServiceImpl, error) {
	tracer, err := aristophanes.NewClientTracer(aristophanes.DefaultAddress)
	if err != nil {
		logging.Error(err.Error())
	}

	healthy := tracer.WaitForHealthyState()
	if !healthy {
		logging.Error("tracing service not ready - restarting seems the only option")
		os.Exit(1)
	}

	streamer, err := tracer.Chorus(ctx)
	if err != nil {
		logging.Error(err.Error())
	}

	version := os.Getenv(config.EnvVersion)

	return &HypatiaServiceImpl{
		Version:  version,
		store:    NewInMemoryStore(),
		Streamer: streamer,
	}, nil
}
