package config

import (
	elastic "github.com/odysseia-greek/agora/aristoteles"
)

type Config struct {
	Index        string
	Created      int
	Elastic      elastic.Client
	Eupalinos    EupalinosClient
	Channel      string
	DutchChannel string
	ExitCode     string
	PolicyName   string
}
