package text

import (
	"context"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	"github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	pb "github.com/odysseia-greek/delphi/ptolemaios/proto"
	aristarchos "github.com/odysseia-greek/olympia/aristarchos/scholar"
	"os"
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

	vaultConfig, err := ambassador.GetSecret(ctx, &pb.VaultRequest{})
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

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
