package text

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	pb "github.com/odysseia-greek/delphi/ptolemaios/proto"
	aristarchos "github.com/odysseia-greek/olympia/aristarchos/scholar"
	"google.golang.org/grpc/metadata"
	"os"
	"time"
)

const (
	defaultIndex string = "text"
)

func CreateNewConfig(ctx context.Context) (*HerodotosHandler, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	tracer, err := aristophanes.NewClientTracer()
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

	ambassador := diplomat.NewClientAmbassador()
	ambassadorHealthy := ambassador.WaitForHealthyState()
	if !ambassadorHealthy {
		logging.Info("tracing service not ready - restarting seems the only option")
		os.Exit(1)
	}

	traceID := uuid.New().String()
	spanID := aristophanes.GenerateSpanID()
	combinedID := fmt.Sprintf("%s+%s+%d", traceID, spanID, 1)

	ambassadorCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	payload := &pbar.StartTraceRequest{
		Method:        "GetSecret",
		Url:           diplomat.DEFAULTADDRESS,
		Host:          "",
		RemoteAddress: "",
		Operation:     "/delphi_ptolemaios.Ptolemaios/GetSecret",
	}

	go func() {
		parabasis := &pbar.ParabasisRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			SpanId:       spanID,
			RequestType: &pbar.ParabasisRequest_StartTrace{
				StartTrace: payload,
			},
		}
		if err := streamer.Send(parabasis); err != nil {
			logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
		}

		logging.Trace(fmt.Sprintf("trace with requestID: %s and span: %s", traceID, spanID))
	}()

	md := metadata.New(map[string]string{service.HeaderKey: combinedID})
	ambassadorCtx = metadata.NewOutgoingContext(context.Background(), md)
	vaultConfig, err := ambassador.GetSecret(ambassadorCtx, &pb.VaultRequest{})
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	go func() {
		parabasis := &pbar.ParabasisRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			SpanId:       spanID,
			RequestType: &pbar.ParabasisRequest_CloseTrace{
				CloseTrace: &pbar.CloseTraceRequest{
					ResponseBody: fmt.Sprintf("user retrieved from vault: %s", vaultConfig.ElasticUsername),
				},
			},
		}

		err := streamer.Send(parabasis)
		if err != nil {
			logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
		}

		logging.Trace(fmt.Sprintf("trace closed with id: %s", traceID))
	}()

	elasticService := aristoteles.ElasticService(tls)

	cfg := models.Config{
		Service:     elasticService,
		Username:    vaultConfig.ElasticUsername,
		Password:    vaultConfig.ElasticPassword,
		ElasticCERT: vaultConfig.ElasticCERT,
	}

	elastic, err := aristoteles.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	err = aristoteles.HealthCheck(elastic)
	if err != nil {
		return nil, err
	}

	index := config.StringFromEnv(config.EnvIndex, defaultIndex)

	aggregatorAddress := config.StringFromEnv(config.EnvAggregatorAddress, config.DefaultAggregatorAddress)
	aggregator := aristarchos.NewClientAggregator(aggregatorAddress)
	aggregatorHealthy := aggregator.WaitForHealthyState()
	if !aggregatorHealthy {
		logging.Debug("aggregator service not ready - restarting seems the only option")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(ctx)

	return &HerodotosHandler{
		Index:      index,
		Elastic:    elastic,
		Aggregator: aggregator,
		Streamer:   streamer,
		Cancel:     cancel,
	}, nil
}
