package gateway

import (
	"context"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	v1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
)

type HomerosHandler struct {
	HttpClients          service.OdysseiaClient
	Cache                archytas.Client
	Streamer             v1.TraceService_ChorusClient
	Cancel               context.CancelFunc
	Randomizer           randomizer.Random
	SokratesGraphqlUrl   string
	AlexandrosGraphqlUrl string
}
