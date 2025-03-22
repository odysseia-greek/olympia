package routing

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/mux"
	plato "github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/olympia/homeros/gateway"
	"github.com/odysseia-greek/olympia/homeros/graph"
	"github.com/odysseia-greek/olympia/homeros/middleware"
	"github.com/vektah/gqlparser/v2/ast"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(handlerConfig *gateway.HomerosHandler, config *gateway.TraceConfig, random randomizer.Random) *mux.Router {
	serveMux := mux.NewRouter()

	srv := handler.New(graph.NewExecutableSchema(
		graph.Config{Resolvers: &graph.Resolver{Handler: handlerConfig}},
	))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](100))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	graphqlHandler := middleware.Adapt(
		srv,
		middleware.LogRequestDetails(handlerConfig.Streamer, config, random),
		middleware.SetCorsHeaders(),
	)

	serveMux.HandleFunc("/homeros/v1/health", plato.Adapt(gateway.HealthProbe, plato.ValidateRestMethod("GET"), plato.SetCorsHeaders()))
	serveMux.Handle("/graphql", graphqlHandler)

	return serveMux
}
