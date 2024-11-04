package scholar

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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
	"time"
)

const (
	defaultIndex string = "aggregator"
)

var streamer pbar.TraceService_ChorusClient

func CreateNewConfig(ctx context.Context) (*AggregatorServiceImpl, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	tracer, err := aristophanes.NewClientTracer()
	healthy := tracer.WaitForHealthyState()
	if !healthy {
		logging.Error("tracing service not ready - restarting seems the only option")
		os.Exit(1)
	}

	streamer, err = tracer.Chorus(ctx)
	if err != nil {
		logging.Error(err.Error())
	}

	var cfg models.Config
	ambassador := diplomat.NewClientAmbassador()

	healthy = ambassador.WaitForHealthyState()
	if !healthy {
		logging.Info("ambassador service not ready - restarting seems the only option")
		os.Exit(1)
	}

	traceID := uuid.New().String()
	spanID := aristophanes.GenerateSpanID()
	combinedID := fmt.Sprintf("%s+%s+%d", traceID, spanID, 1)

	ambassadorCtx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	payload := &pbar.StartTraceRequest{
		Method:        "GetSecret",
		Url:           diplomat.DEFAULTADDRESS,
		Host:          "",
		RemoteAddress: "",
		Operation:     "/delphi_ptolemaios.Ptolemaios/GetSecret",
	}

	go func() {
		parabasis := &pbar.ParabasisRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			SpanId:       spanID,
			RequestType: &pbar.ParabasisRequest_StartTrace{
				StartTrace: payload,
			},
		}
		if err := streamer.Send(parabasis); err != nil {
			logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
		}

		logging.Trace(fmt.Sprintf("trace with requestID: %s and span: %s", traceID, spanID))
	}()

	md := metadata.New(map[string]string{service.HeaderKey: combinedID})
	ambassadorCtx = metadata.NewOutgoingContext(context.Background(), md)
	vaultConfig, err := ambassador.GetSecret(ambassadorCtx, &pb.VaultRequest{})
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	go func() {
		parabasis := &pbar.ParabasisRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			SpanId:       spanID,
			RequestType: &pbar.ParabasisRequest_CloseTrace{
				CloseTrace: &pbar.CloseTraceRequest{
					ResponseBody: fmt.Sprintf("user retrieved from vault: %s", vaultConfig.ElasticUsername),
				},
			},
		}

		err := streamer.Send(parabasis)
		if err != nil {
			logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
		}

		logging.Trace(fmt.Sprintf("trace closed with id: %s", traceID))
	}()

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

	index := config.StringFromEnv(config.EnvIndex, defaultIndex)

	err = createIndexAtStartup(elastic, index)
	if err != nil {
		logging.Error(fmt.Sprintf("index creation returned an error, this most likely means the index already exists and this function should be moved to a job: %s", err.Error()))
	}

	return &AggregatorServiceImpl{
		Index:   index,
		Elastic: elastic,
	}, nil
}

// perhaps it would be best to move this to a different job so that an index is created and aristarchos can be switched from hybrid to a regular api
func createIndexAtStartup(elastic aristoteles.Client, indexName string) error {
	policyName := fmt.Sprintf("%s_policy", indexName)
	logging.Info(fmt.Sprintf("creating policy: %s", policyName))
	err := createPolicyAtStartup(elastic, policyName)
	if err != nil {
		return err
	}

	indexMapping := createScholarIndexMapping()
	created, err := elastic.Index().Create(indexName, indexMapping)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", indexName, created.Acknowledged))

	return nil
}

func createPolicyAtStartup(elastic aristoteles.Client, policyName string) error {
	policyCreated, err := elastic.Policy().CreateHotPolicy(policyName)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", policyName, policyCreated.Acknowledged))

	return nil
}

func createScholarIndexMapping() map[string]interface{} {
	return map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"rootWord": map[string]interface{}{
					"type": "text",
				},
				"unaccented": map[string]interface{}{
					"type": "text",
				},
				"variants": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"searchTerm": map[string]interface{}{
							"type": "text",
						},
						"score": map[string]interface{}{
							"type": "integer",
						},
					},
				},
				"partOfSpeech": map[string]interface{}{
					"type": "keyword",
				},
				"translations": map[string]interface{}{
					"type": "text",
				},
				"categories": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"forms": map[string]interface{}{
							"type": "nested",
							"properties": map[string]interface{}{
								"rule": map[string]interface{}{
									"type": "text",
								},
								"word": map[string]interface{}{
									"type": "text",
								},
							},
						},
					},
				},
			},
		},
	}
}
