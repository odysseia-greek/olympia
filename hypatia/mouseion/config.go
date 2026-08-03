package mouseion

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/odysseia-greek/agora/plato/config"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
)

func CreateNewConfig(ctx context.Context) (*HypatiaServiceImpl, error) {
	tracer, err := aristophanes.NewClientTracer(aristophanes.DefaultAddress)
	if err != nil {
		return nil, fmt.Errorf("create tracing client: %w", err)
	}

	healthy := tracer.WaitForHealthyState()
	if !healthy {
		return nil, fmt.Errorf("tracing service is not ready")
	}

	streamer, err := tracer.Chorus(ctx)
	if err != nil {
		return nil, fmt.Errorf("open tracing stream: %w", err)
	}

	version := os.Getenv(config.EnvVersion)
	maxEvents := DefaultMaxEvents
	if configured := os.Getenv("MAX_EVENTS"); configured != "" {
		maxEvents, err = strconv.Atoi(configured)
		if err != nil || maxEvents <= 0 {
			return nil, fmt.Errorf("MAX_EVENTS must be a positive integer")
		}
	}

	return &HypatiaServiceImpl{
		Version:  version,
		store:    NewInMemoryStoreWithLimit(maxEvents),
		Streamer: streamer,
	}, nil
}
