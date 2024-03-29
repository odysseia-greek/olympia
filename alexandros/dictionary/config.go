package dictionary

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
	"log"
	"os"
	"strings"
	"time"
)

const (
	defaultIndex string = "dictionary"
)

// CreateNewConfig creates a new configuration based on the provided environment.
//
// The function performs the following steps:
//  1. Determines whether health check should be enabled based on the environment.
//  2. Retrieves configuration values from the Vault if health check is enabled.
//  3. Initializes the Elasticsearch service based on the TLS setting.
//  4. Creates the configuration object with the Elasticsearch service, username, password, and certificate.
//     - If health check is disabled, the configuration is created using the environment, testOverWrite, and TLS settings.
//     - If health check is enabled, the configuration is created using the retrieved Vault configuration values.
//  5. Creates a new Elasticsearch client using the configuration.
//  6. Performs a health check on the Elasticsearch client if health check is enabled.
//  7. Retrieves the index name from the environment variables or uses the default index name.
//  8. Returns the created configuration with the Elasticsearch client and index name.
//
// Parameters:
//   - env: The environment name (e.g., "LOCAL", "TEST").
//
// Returns:
//   - *Config: The created configuration containing the Elasticsearch client and index name.
//   - error: An error if any occurred during the configuration creation process.
func CreateNewConfig(env string) (*AlexandrosHandler, error) {
	healthCheck := true
	if env == "DEVELOPEMENT" {
		healthCheck = false
	}

	testOverWrite := config.BoolFromEnv(config.EnvTestOverWrite)
	tls := config.BoolFromEnv(config.EnvTlSKey)

	tracer := aristophanes.NewClientTracer()
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
				logging.Info("ambassador service not ready - restarting seems the only option")
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

	client, err := config.CreateOdysseiaClient()
	if err != nil {
		return nil, err
	}

	index := config.StringFromEnv(config.EnvIndex, defaultIndex)
	return &AlexandrosHandler{
		Index:   index,
		Elastic: elastic,
		Tracer:  tracer,
		Client:  client,
	}, nil
}
