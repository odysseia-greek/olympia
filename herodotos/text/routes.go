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
	serveMux.HandleFunc("/herodotos/v1/createQuestion", middleware.Adapt(herodotosHandler.createQuestion, middleware.ValidateRestMethod("GET"), middleware.Adapter(comedy.TraceWithLogAndSpan(herodotosHandler.Streamer))))
	serveMux.HandleFunc("/herodotos/v1/authors", middleware.Adapt(herodotosHandler.queryAuthors, middleware.ValidateRestMethod("GET"), middleware.Adapter(comedy.TraceWithLogAndSpan(herodotosHandler.Streamer))))
	serveMux.HandleFunc("/herodotos/v1/authors/{author}/books", middleware.Adapt(herodotosHandler.queryBooks, middleware.ValidateRestMethod("GET"), middleware.Adapter(comedy.TraceWithLogAndSpan(herodotosHandler.Streamer))))
	serveMux.HandleFunc("/herodotos/v1/checkSentence", middleware.Adapt(herodotosHandler.checkSentence, middleware.ValidateRestMethod("POST"), middleware.Adapter(comedy.TraceWithLogAndSpan(herodotosHandler.Streamer))))
	serveMux.HandleFunc("/herodotos/v1/texts", middleware.Adapt(herodotosHandler.analyseText, middleware.ValidateRestMethod("GET"), middleware.Adapter(comedy.TraceWithLogAndSpan(herodotosHandler.Streamer))))

	return serveMux
}
