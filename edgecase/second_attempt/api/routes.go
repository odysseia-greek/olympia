package api

import (
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(apiHandler *DiogenesHandler) *mux.Router {
	serveMux := mux.NewRouter()

	serveMux.HandleFunc("/diogenes/v1/ping", middleware.Adapt(apiHandler.pingPong, middleware.ValidateRestMethod("GET")))
	serveMux.HandleFunc("/diogenes/v1/health", middleware.Adapt(apiHandler.health, middleware.ValidateRestMethod("GET")))

	serveMux.HandleFunc("/diogenes/v1/words/_convert", middleware.Adapt(apiHandler.convert, middleware.ValidateRestMethod("POST"), middleware.Adapter(comedy.TraceWithLogAndSpan(apiHandler.Streamer))))

	return serveMux
}
