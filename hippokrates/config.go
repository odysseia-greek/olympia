package hippokrates

import (
	"context"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
)

type OdysseiaFixture struct {
	ctx        context.Context
	client     service.OdysseiaClient
	homeros    Homeros
	randomizer randomizer.Random
}

type Homeros struct {
	graphql string
	health  string
}

func New() (*OdysseiaFixture, error) {
	svc, err := config.CreateOdysseiaClient()
	if err != nil {
		return nil, err
	}

	var gqlEndpoint string
	var healthEndpoint string

	homerosBaseUrl := config.StringFromEnv("HOMEROS_SERVICE", "")

	if homerosBaseUrl == "" {
		gqlEndpoint = "http://k3d-odysseia.greek:8080/graphql"
		healthEndpoint = "http://k3d-odysseia.api.greek:8080/homeros/v1/health"
	} else {
		gqlEndpoint = homerosBaseUrl + "/graphql"
		healthEndpoint = homerosBaseUrl + "/homeros/v1/health"
	}

	randomizerClient, err := randomizer.NewRandomizerClient()
	if err != nil {
		return nil, err
	}

	return &OdysseiaFixture{
		client: svc,
		homeros: Homeros{
			graphql: gqlEndpoint,
			health:  healthEndpoint,
		},
		ctx:        context.Background(),
		randomizer: randomizerClient,
	}, nil
}
