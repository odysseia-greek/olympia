package config

import (
	elastic "github.com/odysseia-greek/agora/aristoteles"
	ptolemaios "github.com/odysseia-greek/delphi/ptolemaios/app"
)

type Config struct {
	Index      string
	Created    int
	Elastic    elastic.Client
	PolicyName string
	Ambassador *ptolemaios.ClientAmbassador
}
