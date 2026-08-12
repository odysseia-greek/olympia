package gateway

import (
	"net/http"
	"os"
	"time"

	"github.com/odysseia-greek/agora/plato/middleware"
	plato "github.com/odysseia-greek/agora/plato/models"
)

func HealthProbe(w http.ResponseWriter, _ *http.Request) {
	version := os.Getenv("VERSION")
	if version == "" {
		version = "unknown"
	}
	middleware.ResponseWithJson(w, plato.Health{Healthy: true, Version: version, Time: time.Now().String()})
}
