package app

import (
	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	plato "github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/randomizer"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/app"
	"github.com/odysseia-greek/olympia/homeros/handlers"
	"github.com/odysseia-greek/olympia/homeros/middleware"
	"github.com/odysseia-greek/olympia/homeros/schemas"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(tracer *aristophanes.ClientTracer, config *handlers.TraceConfig, random randomizer.Random) *mux.Router {
	serveMux := mux.NewRouter()

	srv := handler.New(&handler.Config{
		Schema:   &schemas.HomerosSchema,
		Pretty:   true,
		GraphiQL: false,
	})

	serveMux.HandleFunc("/homeros/v1/health", plato.Adapt(handlers.HealthProbe, plato.ValidateRestMethod("GET"), plato.SetCorsHeaders()))
	serveMux.Handle("/graphql", middleware.Adapt(srv, middleware.LogRequestDetails(tracer, config, random), middleware.SetCorsHeaders()))

	return serveMux
}
