package app

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/service"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/app"
	pb "github.com/odysseia-greek/delphi/ptolemaios/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
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
	if env == "LOCAL" || env == "TEST" {
		healthCheck = false
	}
	testOverWrite := config.BoolFromEnv(config.EnvTestOverWrite)
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config

	if healthCheck {
		vaultConfig, err := configFromVault()
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

func dialGrpcService(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// ConfigFromVault establishes a gRPC connection to the Ptolemaios service,
// checks the health status, and retrieves a secret from the service.
// It retries with an increasing sleep time for establishing the gRPC
// connection and for the health check. The maximum sleep duration for the
// health check is capped at 8 seconds for the last 5 attempts.
// If successful, it returns the retrieved secret and nil error.
// If the maximum number of attempts is reached without success, it returns
// an error indicating the failure.
func configFromVault(optionalName ...string) (*pb.ElasticConfigVault, error) {
	var name string

	if len(optionalName) > 0 && optionalName[0] != "" {
		name = optionalName[0]
	}

	sidecarService := os.Getenv(config.EnvPtolemaiosService)
	if sidecarService == "" {
		log.Printf("defaulting to %s for sidecar", config.DefaultSidecarService)
		sidecarService = config.DefaultSidecarService
	}

	var grpcConnection *grpc.ClientConn
	var err error

	// Retry with increasing sleep time for establishing the gRPC connection
	maxAttempts := 10
	sleepDuration := 500 * time.Millisecond // Starting sleep duration
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		grpcConnection, err = dialGrpcService(sidecarService)
		if err != nil {
			log.Printf("error returned setting up connection to grpc: %s", err.Error())
			if attempt == maxAttempts {
				break // Skip increasing sleep duration on the last attempt
			}
			time.Sleep(sleepDuration)
			// Increase sleep duration for the next attempt
			sleepDuration *= 2
			continue
		}
		break
	}

	if err != nil {
		return nil, fmt.Errorf("failed to establish gRPC connection: %s", err)
	}

	defer func() {
		if e := grpcConnection.Close(); e != nil {
			log.Printf("failed to close connection: %s", e)
		}
	}()

	client := pb.NewPtolemaiosClient(grpcConnection)

	// Health check loop with timeout and increasing sleep time
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	maxAttempts = 10
	sleepDuration = 500 * time.Millisecond // Starting sleep duration
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		r, err := client.Health(ctx, &pb.HealthRequest{})
		if err != nil {
			if attempt == maxAttempts {
				break // Skip increasing sleep duration on the last attempt
			}
			time.Sleep(sleepDuration)
			// Increase sleep duration for the next attempt
			sleepDuration *= 2
			continue
		}

		if r.Health {
			break
		}

		time.Sleep(sleepDuration)
		// Increase sleep duration for the next attempt
		sleepDuration *= 2
		if attempt == maxAttempts-5 {
			// Cap the maximum sleep duration to 8 seconds for the last 5 attempts
			sleepDuration = 8 * time.Second
		}
	}

	// Secret retrieval loop with timeout and increasing sleep time
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	maxAttempts = 10
	sleepDuration = 500 * time.Millisecond // Starting sleep duration

	// Create a context with the custom metadata
	traceId := uuid.New().String()
	md := metadata.New(map[string]string{service.HeaderKey: traceId})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		var r *pb.ElasticConfigVault
		var err error

		if name != "" {
			r, err = client.GetNamedSecret(ctx, &pb.VaultRequestNamed{PodName: name})
		} else {
			r, err = client.GetSecret(ctx, &pb.VaultRequest{})
		}
		if err != nil {
			log.Printf("error getting response from ptolemaios (secret): %s", err)
			if attempt == maxAttempts {
				break // Skip increasing sleep duration on the last attempt
			}
			time.Sleep(sleepDuration)
			// Increase sleep duration for the next attempt
			sleepDuration *= 2
			continue
		}

		return r, nil
	}

	return nil, fmt.Errorf("failed to get a secret from ptolemaios")
}
