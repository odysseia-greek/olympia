package config

import (
	elastic "github.com/odysseia-greek/agora/aristoteles"
	ptolemaios "github.com/odysseia-greek/delphi/ptolemaios/app"
)

type Config struct {
	Index      string
	MaxAge     string
	PolicyName string
	Elastic    elastic.Client
	Ambassador *ptolemaios.ClientAmbassador
}
