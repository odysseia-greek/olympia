package gateway

import (
	"context"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
)

type HomerosHandler struct {
	HttpClients        service.OdysseiaClient
	Cache              archytas.Client
	Streamer           pb.TraceService_ChorusClient
	Cancel             context.CancelFunc
	Randomizer         randomizer.Random
	SokratesGraphqlUrl string
}
