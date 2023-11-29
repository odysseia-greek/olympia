package config

import (
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"time"
)

type Config struct {
	Index        string
	Created      int
	Updated      int
	Processed    int
	Elastic      elastic.Client
	Eupalinos    EupalinosClient
	Channel      string
	DutchChannel string
	WaitTime     time.Duration
}
