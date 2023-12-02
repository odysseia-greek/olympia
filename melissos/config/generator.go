package config

import (
	"context"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/agora/thales"
	ptolemaios "github.com/odysseia-greek/delphi/ptolemaios/app"
	pbp "github.com/odysseia-greek/delphi/ptolemaios/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os"
	"time"
)

const (
	defaultIndex    string = "dictionary"
	EnvWaitTime     string = "WAIT_TIME"
	DefaultWaitTime string = "120"
)

type EupalinosClient interface {
	DequeueMessage(ctx context.Context, in *pb.ChannelInfo, opts ...grpc.CallOption) (*pb.Epistello, error)
	GetQueueLength(ctx context.Context, in *pb.ChannelInfo, opts ...grpc.CallOption) (*pb.QueueLength, error)
}

func CreateNewConfig(env string) (*Config, *grpc.ClientConn, error) {
	healthCheck := true
	if env == "DEVELOPMENT" {
		healthCheck = false
	}
	testOverWrite := config.BoolFromEnv(config.EnvTestOverWrite)
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config

	kube, err := thales.CreateKubeClient(healthCheck)
	if err != nil {
		return nil, nil, err
	}
	ns := config.StringFromEnv(config.EnvNamespace, config.DefaultNamespace)
	job := config.StringFromEnv(config.EnvJobName, config.DefaultJobName)

	ambassador := ptolemaios.NewClientAmbassador()
	if healthCheck {
		if healthCheck {
			healthy := ambassador.WaitForHealthyState()
			if !healthy {
				logging.Info("tracing service not ready - restarting seems the only option")
				os.Exit(1)
			}
		}

		traceId := uuid.New().String()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		md := metadata.New(map[string]string{service.HeaderKey: traceId})
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		vaultConfig, err := ambassador.GetSecret(ctx, &pbp.VaultRequest{})
		if err != nil {
			logging.Error(err.Error())
			return nil, nil, err
		}

		elasticService := aristoteles.ElasticService(tls)

		cfg = models.Config{
			Service:     elasticService,
			Username:    vaultConfig.ElasticUsername,
			Password:    vaultConfig.ElasticPassword,
			ElasticCERT: vaultConfig.ElasticCERT,
		}
	} else {
		cfg = aristoteles.ElasticConfig(env, testOverWrite, tls)
	}

	elastic, err := aristoteles.NewClient(cfg)
	if err != nil {
		return nil, nil, err
	}

	channel := config.StringFromEnv(config.EnvChannel, config.DefaultParmenidesChannel)
	eupalinosAddress := config.StringFromEnv(config.EnvEupalinosService, config.DefaultEupalinosService)

	index := config.StringFromEnv(config.EnvIndex, defaultIndex)
	wait := os.Getenv(EnvWaitTime)

	client, conn, err := createEupalinosClient(eupalinosAddress)
	if err != nil {
		return nil, nil, err
	}
	logging.Debug("client config created")

	var waitDuration time.Duration

	if wait == "" {
		waitDuration, _ = time.ParseDuration(DefaultWaitTime + "s")
	} else {
		waitDuration, _ = time.ParseDuration(wait + "ms")
	}

	return &Config{
		Index:        index,
		Created:      0,
		Updated:      0,
		Processed:    0,
		Elastic:      elastic,
		Eupalinos:    client,
		Channel:      channel,
		DutchChannel: config.DefaultDutchChannel,
		WaitTime:     waitDuration,
		Kube:         kube,
		Job:          job,
		Namespace:    ns,
		Ambassador:   ambassador,
	}, conn, nil
}

func createEupalinosClient(serverAddress string) (pb.EupalinosClient, *grpc.ClientConn, error) {
	logging.Debug("creating client config")
	logging.Debug(serverAddress)
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewEupalinosClient(conn)
	return client, conn, nil
}
