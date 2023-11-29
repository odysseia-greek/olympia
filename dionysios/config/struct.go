package config

import (
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/app"
)

type Config struct {
	Elastic          aristoteles.Client
	Cache            archytas.Client
	Index            string
	Client           service.OdysseiaClient
	DeclensionConfig models.DeclensionConfig
	Tracer           *aristophanes.ClientTracer
}
