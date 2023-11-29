package config

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	pb "github.com/odysseia-greek/eupalinos/proto"
	"google.golang.org/grpc"
)

const (
	defaultIndex string = "quiz"
)

type EupalinosClient interface {
	EnqueueMessage(ctx context.Context, in *pb.Epistello, opts ...grpc.CallOption) (*pb.EnqueueResponse, error)
}

func CreateNewConfig(env string) (*Config, *grpc.ClientConn, error) {
	healthCheck := true
	if env == "LOCAL" || env == "TEST" {
		healthCheck = false
	}
	testOverWrite := config.BoolFromEnv(config.EnvTestOverWrite)
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config

	if healthCheck {
		vaultConfig, err := config.ConfigFromVault()
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

	if healthCheck {
		err := aristoteles.HealthCheck(elastic)
		if err != nil {
			return nil, nil, err
		}
	}

	index := config.StringFromEnv(config.EnvIndex, defaultIndex)

	client, conn, err := createEupalinosClient(eupalinosAddress)
	if err != nil {
		return nil, nil, err
	}

	policyName := fmt.Sprintf("%s_policy", index)

	return &Config{
		Index:        index,
		Created:      0,
		Elastic:      elastic,
		Eupalinos:    client,
		Channel:      channel,
		DutchChannel: config.DefaultDutchChannel,
		PolicyName:   policyName,
	}, conn, nil
}

func createEupalinosClient(serverAddress string) (pb.EupalinosClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewEupalinosClient(conn)
	return client, conn, nil
}
