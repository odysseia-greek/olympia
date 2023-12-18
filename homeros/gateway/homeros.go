package gateway

import (
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
)

type HomerosHandler struct {
	HttpClients service.OdysseiaClient
	Cache       archytas.Client
	Tracer      *aristophanes.ClientTracer
	Randomizer  randomizer.Random
}
