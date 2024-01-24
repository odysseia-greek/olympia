package quiz

import (
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/agora/plato/middleware"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(sokratesHandler *SokratesHandler) *mux.Router {
	serveMux := mux.NewRouter()

	serveMux.HandleFunc("/sokrates/v1/health", middleware.Adapt(sokratesHandler.Health, middleware.ValidateRestMethod("GET")))
	serveMux.HandleFunc("/sokrates/v1/quiz/create", middleware.Adapt(sokratesHandler.Create, middleware.ValidateRestMethod("POST"), middleware.LogRequestDetails()))
	serveMux.HandleFunc("/sokrates/v1/quiz/answer", middleware.Adapt(sokratesHandler.Check, middleware.ValidateRestMethod("POST"), middleware.LogRequestDetails()))
	serveMux.HandleFunc("/sokrates/v1/quiz/options", middleware.Adapt(sokratesHandler.Options, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails()))

	// start handling updates
	go sokratesHandler.updateElasticsearch()

	return serveMux
}
