package text

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	aristarchos "github.com/odysseia-greek/alexandreia/aristarchos/scholar"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"github.com/odysseia-greek/delphi/aristides/diplomat"
	pb "github.com/odysseia-greek/delphi/aristides/proto"
	"google.golang.org/grpc/metadata"
)

const (
	defaultIndex string = "text"
)

func CreateNewConfig(ctx context.Context) (*HerodotosHandler, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	tracer, err := aristophanes.NewClientTracer(aristophanes.DefaultAddress)
	if err != nil {
		logging.Error(err.Error())
	}

	healthy := tracer.WaitForHealthyState()
	if !healthy {
		logging.Error("tracing service not ready - setting tracer to nil and starting backup process")
		tracer = nil
	}

	streamer, err := tracer.Chorus(ctx)
	if err != nil {
		logging.Error(err.Error())
	}

	ambassador, err := diplomat.NewClientAmbassador(diplomat.DEFAULTADDRESS)
	ambassadorHealthy := ambassador.WaitForHealthyState()
	if !ambassadorHealthy {
		logging.Info("ambassador service not ready - restarting seems the only option")
		os.Exit(1)
	}

	traceID := uuid.New().String()
	spanID := aristophanes.GenerateSpanID()
	combinedID := fmt.Sprintf("%s+%s+%d", traceID, spanID, 1)

	ambassadorCtx, ctxCancel := context.WithTimeout(ctx, 30*time.Second)
	defer ctxCancel()

	payload := &arv1.ObserveTraceStart{
		Method:        "GetSecret",
		Url:           diplomat.DEFAULTADDRESS,
		Host:          "",
		RemoteAddress: "",
		Operation:     "/delphi_ptolemaios.Ptolemaios/GetSecret",
	}

	go func() {
		parabasis := &arv1.ObserveRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			SpanId:       spanID,
			Kind: &arv1.ObserveRequest_TraceStart{
				TraceStart: payload,
			},
		}
		if err := streamer.Send(parabasis); err != nil {
			logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
		}

		logging.Trace(fmt.Sprintf("trace with requestID: %s and span: %s", traceID, spanID))
	}()

	md := metadata.New(map[string]string{service.HeaderKey: combinedID})
	ambassadorCtx = metadata.NewOutgoingContext(ambassadorCtx, md)
	vaultConfig, err := ambassador.GetSecret(ambassadorCtx, &pb.VaultRequest{})
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	go func() {
		parabasis := &arv1.ObserveRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			SpanId:       spanID,
			Kind: &arv1.ObserveRequest_TraceStop{
				TraceStop: &arv1.ObserveTraceStop{
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

	cache, err := archytas.CreateBadgerClient()
	if err != nil {
		return nil, err
	}
	index := config.StringFromEnv(config.EnvIndex, defaultIndex)

	logging.Debug("creating new aggregator client")

	aggregatorAddress := config.StringFromEnv(config.EnvAggregatorAddress, config.DefaultAggregatorAddress)
	aggregator, err := aristarchos.NewClientAggregator(aggregatorAddress)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}
	aggregatorHealthy := aggregator.WaitForHealthyState()
	if !aggregatorHealthy {
		logging.Debug("aggregator service not ready - restarting seems the only option")
		os.Exit(1)
	}

	logging.Debug("aggregator client created and healthy")

	ctx, cancel := context.WithCancel(ctx)

	return &HerodotosHandler{
		Index:      index,
		Elastic:    elastic,
		Cache:      cache,
		Aggregator: aggregator,
		Streamer:   streamer,
		Cancel:     cancel,
	}, nil
}
