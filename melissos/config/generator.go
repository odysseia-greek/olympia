package config

import (
	"context"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/thales"
	"github.com/odysseia-greek/olympia/eratosthenes"
	"google.golang.org/grpc"
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
	if env == "LOCAL" || env == "TEST" {
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

	if healthCheck {
		vaultConfig, err := eratosthenes.ConfigFromVault()
		if err != nil {
			return nil, nil, err
		}

		service := aristoteles.ElasticService(tls)

		cfg = models.Config{
			Service:     service,
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
