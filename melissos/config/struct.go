package config

import (
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/thales"
	ptolemaios "github.com/odysseia-greek/delphi/ptolemaios/app"
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
	Kube         thales.KubeClient
	Namespace    string
	Job          string
	Ambassador   *ptolemaios.ClientAmbassador
}
