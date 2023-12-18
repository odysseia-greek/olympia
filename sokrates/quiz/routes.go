package quiz

import (
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/agora/plato/middleware"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(sokratesHandler *SokratesHandler) *mux.Router {
	serveMux := mux.NewRouter()

	serveMux.HandleFunc("/sokrates/v1/ping", middleware.Adapt(sokratesHandler.pingPong, middleware.ValidateRestMethod("GET"), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/health", middleware.Adapt(sokratesHandler.health, middleware.ValidateRestMethod("GET"), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/methods", middleware.Adapt(sokratesHandler.queryMethods, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/methods/{method}/categories", middleware.Adapt(sokratesHandler.queryCategories, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/methods/{method}/categories/{category}/chapters", middleware.Adapt(sokratesHandler.findHighestChapter, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/createQuestion", middleware.Adapt(sokratesHandler.createQuestion, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/sokrates/v1/answer", middleware.Adapt(sokratesHandler.checkAnswer, middleware.ValidateRestMethod("POST"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))

	return serveMux
}
