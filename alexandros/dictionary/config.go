package dictionary

import (
	"context"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	"github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	pb "github.com/odysseia-greek/delphi/ptolemaios/proto"
	"os"
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
func CreateNewConfig(ctx context.Context) (*AlexandrosHandler, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	tracer, err := aristophanes.NewClientTracer()
	if err != nil {
		logging.Error(err.Error())
	}

	healthy := tracer.WaitForHealthyState()
	if !healthy {
		logging.Error("tracing service not ready - restarting seems the only option")
		os.Exit(1)
	}

	streamer, err := tracer.Chorus(ctx)
	if err != nil {
		logging.Error(err.Error())
	}

	ambassador := diplomat.NewClientAmbassador()
	ambassadorHealthy := ambassador.WaitForHealthyState()
	if !ambassadorHealthy {
		logging.Info("ambassador service not ready - restarting seems the only option")
		os.Exit(1)
	}

	vaultConfig, err := ambassador.GetSecret(ctx, &pb.VaultRequest{})
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	elasticService := aristoteles.ElasticService(tls)

	cfg := models.Config{
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

	client, err := config.CreateOdysseiaClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	return &AlexandrosHandler{
		Index:    config.StringFromEnv(config.EnvIndex, defaultIndex),
		Elastic:  elastic,
		Streamer: streamer,
		Client:   client,
		Cancel:   cancel,
	}, nil
}
