package routing

import (
	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	plato "github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/odysseia-greek/olympia/homeros/gateway"
	"github.com/odysseia-greek/olympia/homeros/middleware"
	"github.com/odysseia-greek/olympia/homeros/schemas"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(tracer proto.TraceService_ChorusClient, config *gateway.TraceConfig, random randomizer.Random) *mux.Router {
	serveMux := mux.NewRouter()

	srv := handler.New(&handler.Config{
		Schema:   &schemas.HomerosSchema,
		Pretty:   true,
		GraphiQL: false,
	})

	serveMux.HandleFunc("/homeros/v1/health", plato.Adapt(gateway.HealthProbe, plato.ValidateRestMethod("GET"), plato.SetCorsHeaders()))
	serveMux.Handle("/graphql", middleware.Adapt(srv, middleware.LogRequestDetails(tracer, config, random), middleware.SetCorsHeaders()))

	return serveMux
}
