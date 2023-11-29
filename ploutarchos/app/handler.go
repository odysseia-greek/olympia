package app

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type PloutarchosHandler struct {
}

// PingPong pongs the ping
func (p *PloutarchosHandler) pingPong(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /ping status ping
	//
	// Checks if api is reachable
	//
	//	Consumes:
	//	- application/json
	//
	//	Produces:
	//	- application/json
	//
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: ResultModel
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

// returns the health of the api
func (p *PloutarchosHandler) health(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /health status health
	//
	// Checks if api is healthy
	//
	//	Consumes:
	//	- application/json
	//
	//	Produces:
	//	- application/json
	//
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: Health
	healthy := models.Health{
		Healthy: true,
		Time:    time.Now().String(),
	}

	middleware.ResponseWithJson(w, healthy)
}

func (p *PloutarchosHandler) landingPage(w http.ResponseWriter, req *http.Request) {
	name := "docs/index.html"
	absolutePath, _ := filepath.Abs(name)

	if _, err := os.Stat(absolutePath); errors.Is(err, os.ErrNotExist) {
		absolutePath = filepath.Join("/app", name)
	}

	http.ServeFile(w, req, absolutePath)
}

func (p *PloutarchosHandler) grpcPages(w http.ResponseWriter, req *http.Request) {
	pathParams := mux.Vars(req)
	api := pathParams["api"]

	name := fmt.Sprintf("docs/grpc/%s.html", api)
	absolutePath, _ := filepath.Abs(name)

	if _, err := os.Stat(absolutePath); errors.Is(err, os.ErrNotExist) {
		absolutePath = filepath.Join("/app", name)
	}

	http.ServeFile(w, req, absolutePath)
}

func (p *PloutarchosHandler) docs(w http.ResponseWriter, req *http.Request) {
	pathParams := mux.Vars(req)
	api := pathParams["api"]

	name := "docs/"
	absolutePath, _ := filepath.Abs(name)

	if _, err := os.Stat(absolutePath); errors.Is(err, os.ErrNotExist) {
		absolutePath = filepath.Join("/app", name)
	}

	tmplFile := filepath.Join(absolutePath, "templates", "redoc.tmpl")
	tmpl, err := template.New("redoc").ParseFiles(tmplFile)
	if err != nil {
		log.Print(err)
	}

	err = tmpl.ExecuteTemplate(w, "redoc", api)
	if err != nil {
		log.Print(err)
	}
}

func (p *PloutarchosHandler) serveYamlFiles(w http.ResponseWriter, req *http.Request) {
	pathParams := mux.Vars(req)
	api := pathParams["api"]

	name := fmt.Sprintf("docs/templates/%s.yaml", api)
	absolutePath, _ := filepath.Abs(name)

	if _, err := os.Stat(absolutePath); errors.Is(err, os.ErrNotExist) {
		absolutePath = filepath.Join("/app", name)
	}

	http.ServeFile(w, req, absolutePath)
}
