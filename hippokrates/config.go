package hippokrates

import (
	"context"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/service"
)

type OdysseiaFixture struct {
	ctx     context.Context
	client  service.OdysseiaClient
	homeros Homeros
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
		gqlEndpoint = "http://k3s-odysseia.greek/graphql"
		healthEndpoint = "http://k3s-odysseia.api.greek/homeros/v1/health"
	} else {
		gqlEndpoint = homerosBaseUrl + "/graphql"
		healthEndpoint = homerosBaseUrl + "/homeros/v1/health"
	}

	return &OdysseiaFixture{
		client: svc,
		homeros: Homeros{
			graphql: gqlEndpoint,
			health:  healthEndpoint,
		},
		ctx: context.Background(),
	}, nil
}
