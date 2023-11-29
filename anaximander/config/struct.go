package config

import (
	elastic "github.com/odysseia-greek/agora/aristoteles"
)

type Config struct {
	Index      string
	Created    int
	PolicyName string
	Elastic    elastic.Client
}
