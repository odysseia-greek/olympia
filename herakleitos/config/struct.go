package config

import (
	elastic "github.com/odysseia-greek/agora/aristoteles"
)

type Config struct {
	Index      string
	Created    int
	Elastic    elastic.Client
	PolicyName string
}
