package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	"os"
)

const (
	defaultSokratesAddress = "http://sokrates:8080/sokrates/graphql"
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

	tracer, err := aristophanes.NewClientTracer(aristophanes.DefaultAddress)
	if err != nil {
		logging.Error(err.Error())
	}

	sokratesGraphqlAddress := config.StringFromEnv("SOKRATES_GRAPHQL_ADDRESS", defaultSokratesAddress)

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
		Cache:              cache,
		HttpClients:        service,
		Streamer:           streamer,
		Randomizer:         randomizer,
		Cancel:             cancel,
		SokratesGraphqlUrl: sokratesGraphqlAddress,
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
	defaultTraceConfig := &TraceConfig{
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
				Operation: "authorBasedAnswer",
				Score:     100,
			},
			{
				Operation: "authorBasedQuiz",
				Score:     100,
			},
			{
				Operation: "dialogueQuiz",
				Score:     100,
			}, {
				Operation: "dialogueAnswer",
				Score:     100,
			},
			{
				Operation: "multipleChoiceAnswer",
				Score:     100,
			},
			{
				Operation: "multipleChoiceQuiz",
				Score:     100,
			},
			{
				Operation: "mediaQuiz",
				Score:     100,
			},
			{
				Operation: "mediaAnswer",
				Score:     100,
			},
			//shared
			{
				Operation: "status",
				Score:     50,
			},
		},
	}
	var traceConfig *TraceConfig

	traceConfigPath := os.Getenv("TRACE_CONFIG_PATH")
	if traceConfigPath == "" {
		traceConfig = defaultTraceConfig
	} else {
		traceConfigData, err := os.ReadFile(traceConfigPath)
		if err != nil {
			logging.Warn("could not load trace config. Returning Default")
			traceConfig = defaultTraceConfig
		}

		if err := json.Unmarshal(traceConfigData, &traceConfig); err != nil {
			logging.Error("error unmarshalling data")
			traceConfig = defaultTraceConfig
		}
	}

	jsonPayload, _ := json.MarshalIndent(traceConfig, "", "  ")

	logging.Debug(fmt.Sprintf("found the following trace config: %s", string(jsonPayload)))

	return traceConfig
}
