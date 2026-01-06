package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/middleware"
	plato "github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"github.com/odysseia-greek/olympia/homeros/graph/model"
)

type healthChannel struct {
	name      string
	apiHealth *model.Health
}

func HealthProbe(w http.ResponseWriter, req *http.Request) {
	version := os.Getenv("VERSION")
	if version == "" {
		version = "unknown"
	}
	health := plato.Health{
		Healthy: true,
		Version: version,
		Time:    time.Now().String(),
	}
	middleware.ResponseWithJson(w, health)
}

func (h *HomerosHandler) Health(requestId string) (*model.Status, error) {
	var waitGroup sync.WaitGroup
	c := make(chan *healthChannel)

	waitGroup.Add(3)

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
				apiHealth: &model.Health{
					Healthy: BoolPtr(false),
				},
			}
			c <- &msg
			return
		}
		id := response.Header.Get(service.HeaderKey)
		logging.Info(fmt.Sprintf("route: %s | %s: %s |", response.Request.URL.RequestURI(), service.HeaderKey, id))
		defer response.Body.Close()

		var health model.Health
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
		response, err := h.HttpClients.Alexandros().Health(requestId)
		if err != nil {
			msg := healthChannel{
				name: "dionysios",
				apiHealth: &model.Health{
					Healthy: BoolPtr(false),
				},
			}
			c <- &msg
			return
		}
		id := response.Header.Get(service.HeaderKey)
		logging.Info(fmt.Sprintf("route: %s | %s: %s |", response.Request.URL.RequestURI(), service.HeaderKey, id))
		defer response.Body.Close()

		var health model.Health
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
				apiHealth: &model.Health{
					Healthy: BoolPtr(false),
				},
			}
			c <- &msg
			return
		}
		id := response.Header.Get(service.HeaderKey)
		logging.Info(fmt.Sprintf("route: %s | %s: %s |", response.Request.URL.RequestURI(), service.HeaderKey, id))
		defer response.Body.Close()

		var health model.Health
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

	healthy := model.Status{
		OverallHealth: BoolPtr(true),
	}

	for apiHealth := range c {
		if !*apiHealth.apiHealth.Healthy {
			healthy.OverallHealth = BoolPtr(false)
		}
		switch apiHealth.name {
		case "dionysios":
			healthy.Dionysios = apiHealth.apiHealth
		case "herodotos":
			healthy.Herodotos = apiHealth.apiHealth
		}
	}

	traceID, parentSpanID, traceCall := ParseHeaderID(requestId)
	if traceCall {
		health, err := json.Marshal(healthy)
		parabasis := &arv1.ObserveRequest{
			TraceId:      traceID,
			ParentSpanId: parentSpanID,
			SpanId:       parentSpanID,
			Kind: &arv1.ObserveRequest_TraceStop{
				TraceStop: &arv1.ObserveTraceStop{
					ResponseCode: 200,
					ResponseBody: string(health),
				},
			},
		}

		err = h.Streamer.Send(parabasis)
		if err != nil {
			logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
		}
	}

	return &healthy, nil
}

func BoolPtr(b bool) *bool {
	return &b
}
