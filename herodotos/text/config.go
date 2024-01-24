package text

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	pb "github.com/odysseia-greek/delphi/ptolemaios/proto"
	aristarchos "github.com/odysseia-greek/olympia/aristarchos/scholar"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"strings"
	"time"
)

const (
	defaultIndex string = "text"
)

func CreateNewConfig(env string) (*HerodotosHandler, error) {
	healthCheck := true
	if env == "DEVELOPMENT" {
		healthCheck = false
	}
	testOverWrite := config.BoolFromEnv(config.EnvTestOverWrite)
	tls := config.BoolFromEnv(config.EnvTlSKey)

	tracer := comedy.NewClientTracer()
	if healthCheck {
		healthy := tracer.WaitForHealthyState()
		if !healthy {
			log.Print("tracing service not ready - restarting seems the only option")
			os.Exit(1)
		}
	}

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

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		payload := &pbar.StartTraceRequest{
			Method:        "GetSecret",
			Url:           diplomat.DEFAULTADDRESS,
			Host:          "",
			RemoteAddress: "",
			Operation:     "/delphi_ptolemaios.Ptolemaios/GetSecret",
		}

		trace, err := tracer.StartTrace(ctx, payload)
		if err != nil {
			logging.Error(err.Error())
		}

		logging.Trace(fmt.Sprintf("traceID: %s |", trace.CombinedId))

		md := metadata.New(map[string]string{service.HeaderKey: trace.CombinedId})
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

		splitID := strings.Split(trace.CombinedId, "+")

		var traceID, spanID string

		if len(splitID) >= 1 {
			traceID = splitID[0]
		}
		if len(splitID) >= 2 {
			spanID = splitID[1]
		}
		traceCloser := &pbar.CloseTraceRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			ResponseCode: 200,
			ResponseBody: "redacted",
		}

		_, err = tracer.CloseTrace(context.Background(), traceCloser)
		logging.Trace(fmt.Sprintf("trace closed with id: %s", traceID))
		if err != nil {
			return nil, err
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

	aggregatorAddress := config.StringFromEnv(config.EnvAggregatorAddress, config.DefaultAggregatorAddress)
	aggregator := aristarchos.NewClientAggregator(aggregatorAddress)
	if healthCheck {
		healthy := aggregator.WaitForHealthyState()
		if !healthy {
			logging.Debug("aggregator service not ready - restarting seems the only option")
			os.Exit(1)
		}
	}

	return &HerodotosHandler{
		Index:      index,
		Elastic:    elastic,
		Tracer:     tracer,
		Aggregator: aggregator,
	}, nil
}
