package config

import (
	"bytes"
	elastic "github.com/odysseia-greek/agora/aristoteles"
)

type Config struct {
	Index      string
	SearchWord string
	Created    int
	Elastic    elastic.Client
	MinNGram   int
	MaxNGram   int
	PolicyName string
	Buf        bytes.Buffer
}
