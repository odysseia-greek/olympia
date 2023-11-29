package config

import (
	elastic "github.com/odysseia-greek/agora/aristoteles"
)

type Config struct {
	Index      string
	MaxAge     string
	PolicyName string
	Elastic    elastic.Client
}
