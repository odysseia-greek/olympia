package gateway

import (
	"context"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	dionysiosv1 "github.com/odysseia-greek/alexandreia/dionysios/gen/go/v1"
	v1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"github.com/odysseia-greek/olympia/hypatia/mouseion"
	"google.golang.org/grpc"
)

type DionysiosClient interface {
	Health(context.Context, *dionysiosv1.HealthRequest, ...grpc.CallOption) (*dionysiosv1.HealthResponse, error)
	CheckGrammar(context.Context, *dionysiosv1.CheckGrammarRequest, ...grpc.CallOption) (*dionysiosv1.CheckGrammarResponse, error)
	Research(context.Context, *dionysiosv1.ResearchRequest, ...grpc.CallOption) (*dionysiosv1.ResearchResponse, error)
	TextMode(context.Context, *dionysiosv1.TextModeRequest, ...grpc.CallOption) (*dionysiosv1.TextModeResponse, error)
}

type HomerosHandler struct {
	HttpClients          service.OdysseiaClient
	Dionysios            DionysiosClient
	Cache                archytas.Client
	Streamer             v1.TraceService_ChorusClient
	Hypatia              *mouseion.HypatiaClient
	Cancel               context.CancelFunc
	Randomizer           randomizer.Random
	SokratesGraphqlUrl   string
	AlexandrosGraphqlUrl string
	HerodotosGraphqlUrl  string
	Version              string
	Environment          string
}
