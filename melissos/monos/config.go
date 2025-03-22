package monos

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
	"github.com/odysseia-greek/delphi/aristides/diplomat"
	pbp "github.com/odysseia-greek/delphi/aristides/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func CreateNewConfig(duration time.Duration, finished int64) (*MelissosHandler, *grpc.ClientConn, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	kube, err := thales.CreateKubeClient(false)
	if err != nil {
		return nil, nil, err
	}
	ns := config.StringFromEnv(config.EnvNamespace, config.DefaultNamespace)
	job := config.StringFromEnv(config.EnvJobName, config.DefaultJobName)

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
		return nil, nil, err
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
		return nil, nil, err
	}

	err = aristoteles.HealthCheck(elastic)
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
		waitDuration, _ = time.ParseDuration(wait + "s")
	}

	return &MelissosHandler{
		Duration:     duration,
		TimeFinished: finished,
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
		Namespace:    ns,
		Job:          job,
		Ambassador:   ambassador,
	}, conn, nil
}

func createEupalinosClient(serverAddress string) (pb.EupalinosClient, *grpc.ClientConn, error) {
	logging.Debug("creating client config for Eupalinos")
	logging.Debug(serverAddress)
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewEupalinosClient(conn)
	return client, conn, nil
}
