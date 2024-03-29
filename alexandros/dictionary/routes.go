package dictionary

import (
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/agora/plato/middleware"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(alexandrosHandler *AlexandrosHandler) *mux.Router {
	serveMux := mux.NewRouter()

	serveMux.HandleFunc("/alexandros/v1/ping", middleware.Adapt(alexandrosHandler.pingPong, middleware.ValidateRestMethod("GET")))
	serveMux.HandleFunc("/alexandros/v1/health", middleware.Adapt(alexandrosHandler.health, middleware.ValidateRestMethod("GET")))
	serveMux.HandleFunc("/alexandros/v1/search", middleware.Adapt(alexandrosHandler.searchWord, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails()))

	return serveMux
}
