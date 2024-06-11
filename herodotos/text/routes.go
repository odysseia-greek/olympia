package text

import (
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(herodotosHandler *HerodotosHandler) *mux.Router {
	serveMux := mux.NewRouter()

	serveMux.HandleFunc("/herodotos/v1/ping", middleware.Adapt(herodotosHandler.pingPong, middleware.ValidateRestMethod("GET")))
	serveMux.HandleFunc("/herodotos/v1/health", middleware.Adapt(herodotosHandler.health, middleware.ValidateRestMethod("GET")))

	serveMux.HandleFunc("/herodotos/v1/texts/_create", middleware.Adapt(herodotosHandler.create, middleware.ValidateRestMethod("POST"), middleware.Adapter(comedy.TraceWithLogAndSpan(herodotosHandler.Streamer))))
	serveMux.HandleFunc("/herodotos/v1/texts/_analyze", middleware.Adapt(herodotosHandler.analyze, middleware.ValidateRestMethod("POST"), middleware.Adapter(comedy.TraceWithLogAndSpan(herodotosHandler.Streamer))))
	serveMux.HandleFunc("/herodotos/v1/texts/_check", middleware.Adapt(herodotosHandler.check, middleware.ValidateRestMethod("POST"), middleware.Adapter(comedy.TraceWithLogAndSpan(herodotosHandler.Streamer))))
	serveMux.HandleFunc("/herodotos/v1/texts/options", middleware.Adapt(herodotosHandler.options, middleware.ValidateRestMethod("GET"), middleware.Adapter(comedy.TraceWithLogAndSpan(herodotosHandler.Streamer))))

	return serveMux
}
