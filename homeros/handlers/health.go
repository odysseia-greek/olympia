package handlers

import (
	"context"
	"encoding/json"
	"github.com/odysseia-greek/agora/plato/middleware"
	plato "github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/odysseia-greek/olympia/homeros/models"
	"log"
	"net/http"
	"sync"
	"time"
)

type healthChannel struct {
	name      string
	apiHealth *plato.Health
}

func HealthProbe(w http.ResponseWriter, req *http.Request) {
	health := plato.Health{
		Healthy: true,
		Time:    time.Now().String(),
	}
	middleware.ResponseWithJson(w, health)
}

func (h *HomerosHandler) Health(requestId string) (*models.Health, error) {
	var waitGroup sync.WaitGroup
	c := make(chan *healthChannel)

	waitGroup.Add(4)

	go func() {
		waitGroup.Wait()
		close(c)
	}()

	go func() {
		defer waitGroup.Done()
		response, err := h.HttpClients.Herodotos().Health(requestId)
		if err != nil {
			msg := healthChannel{
				name: "dionysios",
				apiHealth: &plato.Health{
					Healthy: false,
				},
			}
			c <- &msg
			return
		}
		id := response.Header.Get(service.HeaderKey)
		log.Printf("route: %s | %s: %s |", response.Request.URL.RequestURI(), service.HeaderKey, id)
		defer response.Body.Close()

		var health plato.Health
		err = json.NewDecoder(response.Body).Decode(&health)
		if err != nil {
			c <- nil
		}

		msg := healthChannel{
			name:      "herodotos",
			apiHealth: &health,
		}
		c <- &msg
	}()

	go func() {
		defer waitGroup.Done()
		response, err := h.HttpClients.Sokrates().Health(requestId)
		if err != nil {
			msg := healthChannel{
				name: "dionysios",
				apiHealth: &plato.Health{
					Healthy: false,
				},
			}
			c <- &msg
			return
		}
		id := response.Header.Get(service.HeaderKey)
		log.Printf("route: %s | %s: %s |", response.Request.URL.RequestURI(), service.HeaderKey, id)
		defer response.Body.Close()

		var health plato.Health
		err = json.NewDecoder(response.Body).Decode(&health)
		if err != nil {
			c <- nil
		}
		msg := healthChannel{
			name:      "sokrates",
			apiHealth: &health,
		}
		c <- &msg
	}()

	go func() {
		defer waitGroup.Done()
		response, err := h.HttpClients.Alexandros().Health(requestId)
		if err != nil {
			msg := healthChannel{
				name: "dionysios",
				apiHealth: &plato.Health{
					Healthy: false,
				},
			}
			c <- &msg
			return
		}
		id := response.Header.Get(service.HeaderKey)
		log.Printf("route: %s | %s: %s |", response.Request.URL.RequestURI(), service.HeaderKey, id)
		defer response.Body.Close()

		var health plato.Health
		err = json.NewDecoder(response.Body).Decode(&health)
		if err != nil {
			c <- nil
		}

		msg := healthChannel{
			name:      "alexandros",
			apiHealth: &health,
		}
		c <- &msg
	}()

	go func() {
		defer waitGroup.Done()
		response, err := h.HttpClients.Dionysios().Health(requestId)
		if err != nil {
			msg := healthChannel{
				name: "dionysios",
				apiHealth: &plato.Health{
					Healthy: false,
				},
			}
			c <- &msg
			return
		}
		id := response.Header.Get(service.HeaderKey)
		log.Printf("route: %s | %s: %s |", response.Request.URL.RequestURI(), service.HeaderKey, id)
		defer response.Body.Close()

		var health plato.Health
		err = json.NewDecoder(response.Body).Decode(&health)
		if err != nil {
			c <- nil
		}

		msg := healthChannel{
			name:      "dionysios",
			apiHealth: &health,
		}
		c <- &msg
	}()

	healthy := models.Health{
		Overall: true,
	}

	for apiHealth := range c {
		if !apiHealth.apiHealth.Healthy {
			healthy.Overall = false
		}
		switch apiHealth.name {
		case "dionysios":
			healthy.Dionysios = *apiHealth.apiHealth
		case "herodotos":
			healthy.Herodotos = *apiHealth.apiHealth
		case "alexandros":
			healthy.Alexandros = *apiHealth.apiHealth
		case "sokrates":
			healthy.Sokrates = *apiHealth.apiHealth
		}
	}

	traceId, _, traceCall := ParseHeaderID(requestId)
	if traceCall {

		traceCloser := &aristophanes.CloseTraceRequest{
			TraceId:      traceId,
			ResponseCode: 200,
		}

		go h.Tracer.CloseTrace(context.Background(), traceCloser)
	}

	return &healthy, nil
}
