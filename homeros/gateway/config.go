package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	"log"
	"os"
)

func CreateNewConfig(env string) (*HomerosHandler, error) {
	healthCheck := true
	if env == "DEVELOPMENT" {
		healthCheck = false
	}

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

	tracer := aristophanes.NewClientTracer()

	if healthCheck {
		healthy := tracer.WaitForHealthyState()
		if !healthy {
			logging.Error("tracing service not ready - restarting seems the only option")
			os.Exit(1)
		}
	}

	return &HomerosHandler{
		Cache:       cache,
		HttpClients: service,
		Tracer:      tracer,
		Randomizer:  randomizer,
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
					Score:     20, // 20% chance of tracing
				},
				// dionysios
				{
					Operation: "grammar",
					Score:     50, // 50% chance of tracing
				},
				//herodotos
				{
					Operation: "authors",
					Score:     100, // 100% chance of tracing
				},
				{
					Operation: "sentence",
					Score:     100, // 50% chance of tracing
				},
				{
					Operation: "text",
					Score:     100, // 50% chance of tracing
				},
				//sokrates
				{
					Operation: "methods",
					Score:     100, // 100% chance of tracing
				},
				{
					Operation: "answer",
					Score:     100, // 50% chance of tracing
				},
				{
					Operation: "quiz",
					Score:     100, // 50% chance of tracing
				},
				//shared
				{
					Operation: "status",
					Score:     50, // 50% chance of tracing
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
