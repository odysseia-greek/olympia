package monos

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	eupalinos "github.com/odysseia-greek/agora/eupalinos/stomion"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/delphi/aristides/diplomat"
	pbp "github.com/odysseia-greek/delphi/aristides/proto"
	"google.golang.org/grpc/metadata"
)

const (
	defaultIndex    string = "dictionary"
	EnvWaitTime     string = "WAIT_TIME"
	DefaultWaitTime string = "120"
)

func CreateNewConfig(duration time.Duration, finished int64) (*MelissosHandler, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config
	ambassador, err := diplomat.NewClientAmbassador(diplomat.DEFAULTADDRESS)
	healthy := ambassador.WaitForHealthyState()
	if !healthy {
		logging.Info("tracing service not ready - restarting seems the only option")
		os.Exit(1)
	}

	traceId := uuid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	md := metadata.New(map[string]string{service.HeaderKey: traceId})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	vaultConfig, err := ambassador.GetSecret(ctx, &pbp.VaultRequest{})
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	elasticService := aristoteles.ElasticService(tls)

	cfg = models.Config{
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

	channel := config.StringFromEnv(config.EnvChannel, config.DefaultParmenidesChannel)
	jobName := config.StringFromEnv(config.EnvJobName, config.DefaultJobName)
	index := config.StringFromEnv(config.EnvIndex, defaultIndex)
	wait := os.Getenv(EnvWaitTime)

	eupalinosAddress := config.StringFromEnv(config.EnvEupalinosService, config.DefaultEupalinosService)
	logging.Debug(fmt.Sprintf("creating new eupalinos client: %s", eupalinosAddress))
	queue, err := eupalinos.NewEupalinosClient(eupalinosAddress)
	if err != nil {
		logging.Error(err.Error())
	}

	logging.Debug("waiting for queue to be ready")
	queueHealthy := queue.WaitForHealthyState()
	if !queueHealthy {
		logging.Debug("no queue that is healthy")
	}

	var waitDuration time.Duration

	if wait == "" {
		waitDuration, _ = time.ParseDuration(DefaultWaitTime + "s")
	} else {
		waitDuration, _ = time.ParseDuration(wait + "s")
	}

	return &MelissosHandler{
		Duration:             duration,
		TimeFinished:         finished,
		Index:                index,
		Created:              0,
		Updated:              0,
		Processed:            0,
		Elastic:              elastic,
		Eupalinos:            queue,
		Channel:              channel,
		JobCompletionChannel: jobName,
		DutchChannel:         config.DefaultDutchChannel,
		WaitTime:             waitDuration,
		Ambassador:           ambassador,
	}, nil
}
