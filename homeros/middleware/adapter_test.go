package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/olympia/homeros/gateway"
	pb "github.com/odysseia-greek/olympia/hypatia/proto/v1"
)

type eventTrackerStub struct {
	events chan *pb.RequestEvent
}

func (s *eventTrackerStub) TrackEvents(_ context.Context, batch *pb.RequestEventBatch) (*pb.TrackResponse, error) {
	s.events <- batch.Events[0]
	return &pb.TrackResponse{Ack: "received"}, nil
}

type randomStub struct{}

func (randomStub) RandomNumberBaseZero(int) int { return 0 }
func (randomStub) RandomNumberBaseOne(int) int  { return 100 }

func TestLogRequestDetailsTracksCompletedRequest(t *testing.T) {
	tracker := &eventTrackerStub{events: make(chan *pb.RequestEvent, 1)}
	traceConfig := &gateway.TraceConfig{}
	handler := LogRequestDetails(nil, tracker, traceConfig, randomStub{})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	request := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"operationName":"dictionary","query":"query dictionary { dictionary(word: \"λόγος\") }"}`))
	request.Header.Set(config.SessionIdKey, "visitor-1")
	request.Header.Set("User-Agent", "Hypatia test agent")
	request.Header.Set("Referer", "https://example.test/reader")
	request.Header.Set("X-Forwarded-For", "203.0.113.7, 10.0.0.1")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	select {
	case event := <-tracker.events:
		if event.Path != "dictionary" || event.Method != http.MethodPost || event.Status != http.StatusCreated {
			t.Fatalf("request fields = %+v", event)
		}
		if event.SessionId != "visitor-1" || event.Ip != "203.0.113.7" {
			t.Fatalf("visitor fields = %+v", event)
		}
		if event.UserAgent != "Hypatia test agent" || event.Referrer != "https://example.test/reader" {
			t.Fatalf("client fields = %+v", event)
		}
		if event.TraceId != "" {
			t.Fatalf("trace_id = %q, want empty for an unsampled request", event.TraceId)
		}
		if _, err := time.Parse(time.RFC3339Nano, event.Timestamp); err != nil {
			t.Fatalf("timestamp = %q: %v", event.Timestamp, err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Hypatia event")
	}
}

func TestLogRequestDetailsIgnoresInternalUnidentifiedRequest(t *testing.T) {
	tracker := &eventTrackerStub{events: make(chan *pb.RequestEvent, 1)}
	handler := LogRequestDetails(nil, tracker, &gateway.TraceConfig{}, randomStub{})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	request := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"operationName":"status","query":"query status { status { healthy } }"}`))

	handler.ServeHTTP(httptest.NewRecorder(), request)

	select {
	case event := <-tracker.events:
		t.Fatalf("unexpected event for internal request: %+v", event)
	case <-time.After(50 * time.Millisecond):
	}
}
