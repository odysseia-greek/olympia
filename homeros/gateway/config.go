package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	"log"
	"os"
)

func CreateNewConfig(ctx context.Context) (*HomerosHandler, error) {
	cache, err := archytas.CreateBadgerClient()
	if err != nil {
		return nil, err
	}

	service, err := config.CreateOdysseiaClient()
	if err != nil {
		return nil, err
	}

	randomizer, err := config.CreateNewRandomizer()
	if err != nil {
		return nil, err
	}

	tracer, err := aristophanes.NewClientTracer()
	if err != nil {
		logging.Error(err.Error())
	}

	streamer, err := tracer.Chorus(ctx)
	if err != nil {
		logging.Error(err.Error())
	}

	healthy := tracer.WaitForHealthyState()
	if !healthy {
		logging.Error("tracing service not ready - restarting seems the only option")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(ctx)

	return &HomerosHandler{
		Cache:       cache,
		HttpClients: service,
		Streamer:    streamer,
		Randomizer:  randomizer,
		Cancel:      cancel,
	}, nil
}

type OperationScore struct {
	Operation string `json:"operation"`
	Score     int    `json:"score"`
}

type TraceConfig struct {
	OperationScores []OperationScore `json:"operationScores"`
}

func InitTracingConfig() *TraceConfig {
	var traceConfig TraceConfig

	traceConfigPath := os.Getenv("TRACE_CONFIG_PATH")
	if traceConfigPath == "" {
		traceConfig = TraceConfig{
			OperationScores: []OperationScore{
				// alexandros
				{
					Operation: "dictionary",
					Score:     100,
				},
				// dionysios
				{
					Operation: "grammar",
					Score:     100,
				},
				//herodotos
				{
					Operation: "authors",
					Score:     100,
				},
				{
					Operation: "sentence",
					Score:     100,
				},
				{
					Operation: "text",
					Score:     100,
				},
				//sokrates
				{
					Operation: "methods",
					Score:     100,
				},
				{
					Operation: "answer",
					Score:     100,
				},
				{
					Operation: "quiz",
					Score:     100,
				},
				//shared
				{
					Operation: "status",
					Score:     50,
				},
			},
		}
	} else {
		traceConfigData, err := os.ReadFile(traceConfigPath)
		if err != nil {
			log.Print("could not load trace config...")
		}

		if err := json.Unmarshal(traceConfigData, &traceConfig); err != nil {
			log.Print("error unmarshalling data...")
		}
	}

	jsonPayload, _ := json.MarshalIndent(traceConfig, "", "  ")

	logging.Debug(fmt.Sprintf("found the following trace config: %s", string(jsonPayload)))

	return &traceConfig
}
