package app

import (
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/app"
	"github.com/odysseia-greek/olympia/eratosthenes"
	"log"
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
func CreateNewConfig(env string) (*AlexandrosHandler, error) {
	healthCheck := true
	if env == "LOCAL" || env == "TEST" {
		healthCheck = false
	}
	testOverWrite := config.BoolFromEnv(config.EnvTestOverWrite)
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config

	if healthCheck {
		vaultConfig, err := eratosthenes.ConfigFromVault()
		if err != nil {
			log.Print(err)
			return nil, err
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
		return nil, err
	}

	if healthCheck {
		err := aristoteles.HealthCheck(elastic)
		if err != nil {
			return nil, err
		}
	}

	tracer := aristophanes.NewClientTracer()
	if healthCheck {
		healthy := tracer.WaitForHealthyState()
		if !healthy {
			log.Print("tracing service not ready - restarting seems the only option")
			os.Exit(1)
		}
	}

	index := config.StringFromEnv(config.EnvIndex, defaultIndex)
	return &AlexandrosHandler{
		Index:   index,
		Elastic: elastic,
		Tracer:  tracer,
	}, nil
}
