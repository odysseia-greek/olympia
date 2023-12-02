package config

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	ptolemaios "github.com/odysseia-greek/delphi/ptolemaios/app"
	pbp "github.com/odysseia-greek/delphi/ptolemaios/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os"
	"time"
)

const (
	defaultIndex string = "quiz"
)

type EupalinosClient interface {
	EnqueueMessage(ctx context.Context, in *pb.Epistello, opts ...grpc.CallOption) (*pb.EnqueueResponse, error)
}

func CreateNewConfig(env string) (*Config, *grpc.ClientConn, error) {
	healthCheck := true
	if env == "DEVELOPMENT" {
		healthCheck = false
	}

	testOverWrite := config.BoolFromEnv(config.EnvTestOverWrite)
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config
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
		Ambassador:   ambassador,
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
