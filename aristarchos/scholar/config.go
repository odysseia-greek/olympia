package scholar

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	pb "github.com/odysseia-greek/delphi/ptolemaios/proto"
	"google.golang.org/grpc/metadata"
	"os"
	"strings"
	"time"
)

const (
	defaultIndex string = "aggregator"
)

var Tracer *aristophanes.ClientTracer

func CreateNewConfig() (*AggregatorServiceImpl, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	Tracer = aristophanes.NewClientTracer()
	healthy := Tracer.WaitForHealthyState()
	if !healthy {
		logging.Error("tracing service not ready - restarting seems the only option")
		os.Exit(1)
	}

	var cfg models.Config
	ambassador := diplomat.NewClientAmbassador()

	healthy = ambassador.WaitForHealthyState()
	if !healthy {
		logging.Info("ambassador service not ready - restarting seems the only option")
		os.Exit(1)
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

	trace, err := Tracer.StartTrace(ctx, payload)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
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

	_, err = Tracer.CloseTrace(context.Background(), traceCloser)
	logging.Trace(fmt.Sprintf("trace closed with id: %s", traceID))
	if err != nil {
		return nil, err
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
	return &AggregatorServiceImpl{
		Index:   index,
		Elastic: elastic,
	}, nil
}
