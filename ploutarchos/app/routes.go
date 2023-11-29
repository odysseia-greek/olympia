package app

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/agora/plato/middleware"
	"net/http"
	"os"
	"path/filepath"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes() *mux.Router {
	serveMux := mux.NewRouter()

	ploutarchosHandler := PloutarchosHandler{}

	serveMux.HandleFunc("/ploutarchos/v1/ping", middleware.Adapt(ploutarchosHandler.pingPong, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/ploutarchos/v1/health", middleware.Adapt(ploutarchosHandler.health, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/ploutarchos/", middleware.Adapt(ploutarchosHandler.landingPage, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/ploutarchos", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ploutarchos/", http.StatusFound)
	})
	serveMux.HandleFunc("/ploutarchos/redoc/{api}", middleware.Adapt(ploutarchosHandler.docs, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))
	serveMux.HandleFunc("/ploutarchos/redoc/{api}/yaml", middleware.Adapt(ploutarchosHandler.serveYamlFiles, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))

	serveMux.HandleFunc("/ploutarchos/grpc/{api}", middleware.Adapt(ploutarchosHandler.grpcPages, middleware.ValidateRestMethod("GET"), middleware.LogRequestDetails(), middleware.SetCorsHeaders()))

	name := "docs/"
	absolutePath, _ := filepath.Abs(name)
	if _, err := os.Stat(absolutePath); errors.Is(err, os.ErrNotExist) {
		absolutePath = filepath.Join("/app", name)
	}

	public := filepath.Join(absolutePath, "public")

	fs := http.FileServer(http.Dir(public))
	serveMux.PathPrefix("/ploutarchos/homeros/").Handler(http.StripPrefix("/ploutarchos/homeros/", fs))

	return serveMux
}
