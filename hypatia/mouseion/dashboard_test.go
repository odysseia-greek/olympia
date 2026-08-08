package mouseion

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pb "github.com/odysseia-greek/olympia/hypatia/proto/v1"
)

func newDashboardTestService(t *testing.T) *HypatiaServiceImpl {
	t.Helper()
	store := NewInMemoryStore()
	if err := store.Add([]*pb.RequestEvent{
		{Timestamp: "2026-08-03T08:00:00Z", Path: "/search", Method: "GET", Status: 200, Ip: "203.0.113.7", SessionId: "session-1", TraceId: "trace-1"},
		{Timestamp: "2026-08-03T08:01:00Z", Path: "/search", Method: "GET", Status: 404, Ip: "203.0.113.7", SessionId: "session-2"},
		{Timestamp: "2026-08-03T08:02:00Z", Path: "/health", Method: "GET", Status: 200, Ip: "198.51.100.9", SessionId: "session-3"},
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	return &HypatiaServiceImpl{store: store}
}

func TestDashboardSummary(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/summary", nil)
	response := httptest.NewRecorder()
	newDashboardTestService(t).DashboardHandler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	var summary dashboardSummary
	if err := json.NewDecoder(response.Body).Decode(&summary); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if summary.Events != 3 || summary.Visitors != 2 || summary.Sessions != 3 || summary.Paths != 2 || summary.TracedEvents != 1 {
		t.Fatalf("summary = %+v", summary)
	}
	if summary.StatusCounts["2xx"] != 2 || summary.StatusCounts["4xx"] != 1 {
		t.Fatalf("status counts = %v", summary.StatusCounts)
	}
}

func TestDashboardEventsAreNewestFirstAndLimited(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/events?limit=2", nil)
	response := httptest.NewRecorder()
	newDashboardTestService(t).DashboardHandler().ServeHTTP(response, request)

	var events []dashboardEvent
	if err := json.NewDecoder(response.Body).Decode(&events); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(events) != 2 || events[0].Path != "/health" || events[1].Status != 404 {
		t.Fatalf("events = %+v", events)
	}
}

func TestDashboardAggregatesAndServesEmbeddedUI(t *testing.T) {
	handler := newDashboardTestService(t).DashboardHandler()
	visitorsResponse := httptest.NewRecorder()
	handler.ServeHTTP(visitorsResponse, httptest.NewRequest(http.MethodGet, "/api/visitors", nil))
	var visitors []dashboardCount
	if err := json.NewDecoder(visitorsResponse.Body).Decode(&visitors); err != nil {
		t.Fatalf("decode visitors: %v", err)
	}
	if len(visitors) != 2 || visitors[0].Value != "203.0.113.7" || visitors[0].Count != 2 {
		t.Fatalf("visitors = %+v", visitors)
	}

	pathsResponse := httptest.NewRecorder()
	handler.ServeHTTP(pathsResponse, httptest.NewRequest(http.MethodGet, "/api/paths", nil))
	var paths []dashboardCount
	if err := json.NewDecoder(pathsResponse.Body).Decode(&paths); err != nil {
		t.Fatalf("decode paths: %v", err)
	}
	if len(paths) != 2 || paths[0].Value != "/search" || paths[0].Count != 2 {
		t.Fatalf("paths = %+v", paths)
	}

	indexResponse := httptest.NewRecorder()
	handler.ServeHTTP(indexResponse, httptest.NewRequest(http.MethodGet, "/", nil))
	if indexResponse.Code != http.StatusOK {
		t.Fatalf("dashboard status = %d, want %d", indexResponse.Code, http.StatusOK)
	}
	if contentType := indexResponse.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Fatalf("content type = %q", contentType)
	}
}
