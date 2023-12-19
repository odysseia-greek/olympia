package grammar

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	plato "github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	"github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	pb "github.com/odysseia-greek/delphi/ptolemaios/proto"
	aristarchos "github.com/odysseia-greek/olympia/aristarchos/scholar"
	"google.golang.org/grpc/metadata"
	"os"
	"time"
)

const (
	defaultIndex             string = "grammar"
	EnvAggregatorAddress            = "ARISTARCHOS_ADDRESS"
	DEFAULTAGGREGATORADDRESS        = "aristarchos:50053"
)

func CreateNewConfig(env string) (*DionysosHandler, error) {
	healthCheck := true
	if env == "DEVELOPMENT" {
		healthCheck = false
	}
	testOverWrite := config.BoolFromEnv(config.EnvTestOverWrite)
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config
	ambassador := diplomat.NewClientAmbassador()
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
		vaultConfig, err := ambassador.GetSecret(ctx, &pb.VaultRequest{})
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
	} else {
		cfg = aristoteles.ElasticConfig(env, testOverWrite, tls)
	}

	elastic, err := aristoteles.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	if healthCheck {
		err := aristoteles.HealthCheck(elastic)
		if err != nil {
			return nil, err
		}
	}

	index := config.StringFromEnv(config.EnvIndex, defaultIndex)
	cache, err := archytas.CreateBadgerClient()
	if err != nil {
		return nil, err
	}

	client, err := config.CreateOdysseiaClient()
	if err != nil {
		return nil, err
	}

	tracer := aristophanes.NewClientTracer()
	if healthCheck {
		healthy := tracer.WaitForHealthyState()
		if !healthy {
			logging.Debug("tracing service not ready - restarting seems the only option")
			os.Exit(1)
		}
	}

	aggregatorAddress := config.StringFromEnv(EnvAggregatorAddress, DEFAULTAGGREGATORADDRESS)
	aggregator := aristarchos.NewClientAggregator(aggregatorAddress)
	if healthCheck {
		healthy := aggregator.WaitForHealthyState()
		if !healthy {
			logging.Debug("aggregator service not ready - restarting seems the only option")
			os.Exit(1)
		}
	}

	return &DionysosHandler{
		Elastic:          elastic,
		Cache:            cache,
		Index:            index,
		Client:           client,
		DeclensionConfig: plato.DeclensionConfig{},
		Tracer:           tracer,
		Aggregator:       aggregator,
	}, nil
}

func QueryRuleSet(es elastic.Client, index string) (*plato.DeclensionConfig, error) {
	query := es.Builder().MatchAll()
	response, err := es.Query().MatchWithScroll(index, query)

	if err != nil {
		return nil, err
	}
	var declensionConfig plato.DeclensionConfig
	for _, jsonHit := range response.Hits.Hits {
		byteJson, err := json.Marshal(jsonHit.Source)
		if err != nil {
			return nil, err
		}
		declension, err := plato.UnmarshalDeclension(byteJson)
		if err != nil {
			return nil, err
		}

		declensionConfig.Declensions = append(declensionConfig.Declensions, declension)
	}
	return &declensionConfig, nil

}
