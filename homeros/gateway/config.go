package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	v1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
)

const (
	defaultSokratesAddress   = "http://sokrates.apologia.svc:8080/sokrates/graphql"
	defaultAlexandrosAddress = "http://alexandros.makedonia.svc:8080/alexandros/graphql"
)

func CreateNewConfig(ctx context.Context) (*HomerosHandler, error) {
	start := time.Now()
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

	var tracer *aristophanes.ClientTracer
	var streamer v1.TraceService_ChorusClient

	maxRetries := 10
	retryDelay := 3 * time.Second

	for i := 1; i <= maxRetries; i++ {
		tracer, err = aristophanes.NewClientTracer(aristophanes.DefaultAddress)
		if err == nil {
			break
		}

		logging.Error(fmt.Sprintf("failed to create tracer (attempt %d/%d): %s", i, maxRetries, err.Error()))

		if i < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	for i := 1; i <= maxRetries; i++ {
		streamer, err = tracer.Chorus(ctx)
		if err == nil {
			break
		}

		logging.Error(fmt.Sprintf("failed to create chorus streamer (attempt %d/%d): %s", i, maxRetries, err.Error()))
		if i < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	healthyTracer := false
	if tracer != nil {
		healthyTracer = tracer.WaitForHealthyState()
	}

	sokratesGraphqlAddress := config.StringFromEnv("SOKRATES_GRAPHQL_ADDRESS", defaultSokratesAddress)
	alexandrosGraphqlAddress := config.StringFromEnv("ALEXANDROS_GRAPHQL_ADDRESS", defaultAlexandrosAddress)

	ctx, cancel := context.WithCancel(ctx)
	elapsed := time.Since(start)

	version := os.Getenv(config.EnvVersion)
	env := os.Getenv(config.EnvKey)

	logging.System(fmt.Sprintf(`Homeros Configuration Overview:
- Initialization Time: %s
- Tracer Service:      %v (Address: %s)
- Sokrates GraphQL:    %s
- Alexandros GraphQL:  %s
- Homeros Version:     %s
- Environment:         %s
`,
		elapsed,
		healthyTracer, aristophanes.DefaultAddress,
		sokratesGraphqlAddress,
		alexandrosGraphqlAddress,
		version,
		env,
	))

	return &HomerosHandler{
		Cache:                cache,
		HttpClients:          service,
		Streamer:             streamer,
		Randomizer:           randomizer,
		Cancel:               cancel,
		SokratesGraphqlUrl:   sokratesGraphqlAddress,
		AlexandrosGraphqlUrl: alexandrosGraphqlAddress,
		Version:              version,
		Environment:          env,
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
				Score:     10,
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
