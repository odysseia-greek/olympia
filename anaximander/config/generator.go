package config

import (
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/olympia/eratosthenes"
)

const (
	defaultIndex string = "grammar"
)

func CreateNewConfig(env string) (*Config, error) {
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
			logging.Error(err.Error())
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

	index := config.StringFromEnv(config.EnvIndex, defaultIndex)

	policyName := fmt.Sprintf("%s_policy", index)

	return &Config{
		Index:      index,
		Created:    0,
		Elastic:    elastic,
		PolicyName: policyName,
	}, nil
}
