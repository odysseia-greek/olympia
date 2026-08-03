package mouseion

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"sort"
	"strconv"
	"time"

	pb "github.com/odysseia-greek/olympia/hypatia/proto/v1"
)

const (
	defaultDashboardLimit = 100
	maxDashboardLimit     = 1000
)

//go:embed dashboard
var dashboardFiles embed.FS

type dashboardEvent struct {
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
	Method    string `json:"method"`
	Status    int32  `json:"status"`
	IP        string `json:"ip"`
	UserAgent string `json:"userAgent"`
	Referrer  string `json:"referrer"`
	SessionID string `json:"sessionId,omitempty"`
	TraceID   string `json:"traceId,omitempty"`
}

type dashboardSummary struct {
	Events       int            `json:"events"`
	Sessions     int            `json:"sessions"`
	Paths        int            `json:"paths"`
	TracedEvents int            `json:"tracedEvents"`
	StatusCounts map[string]int `json:"statusCounts"`
	GeneratedAt  string         `json:"generatedAt"`
}

type dashboardCount struct {
	Value     string `json:"value"`
	Count     int    `json:"count"`
	LastSeen  string `json:"lastSeen,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
}

// DashboardHandler serves the embedded dashboard and its read-only JSON API.
func (h *HypatiaServiceImpl) DashboardHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/summary", h.dashboardSummary)
	mux.HandleFunc("GET /api/events", h.dashboardEvents)
	mux.HandleFunc("GET /api/sessions", h.dashboardSessions)
	mux.HandleFunc("GET /api/paths", h.dashboardPaths)

	assets, err := fs.Sub(dashboardFiles, "dashboard")
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(assets)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		mux.ServeHTTP(w, r)
	})
}

func (h *HypatiaServiceImpl) dashboardSummary(w http.ResponseWriter, _ *http.Request) {
	events, err := h.store.GetRecent(0)
	if err != nil {
		writeDashboardError(w, err)
		return
	}

	sessions := make(map[string]struct{})
	paths := make(map[string]struct{})
	statusCounts := make(map[string]int)
	traced := 0
	for _, event := range events {
		if event.SessionId != "" {
			sessions[event.SessionId] = struct{}{}
		}
		if event.Path != "" {
			paths[event.Path] = struct{}{}
		}
		statusCounts[statusGroup(event.Status)]++
		if event.TraceId != "" {
			traced++
		}
	}

	writeDashboardJSON(w, dashboardSummary{
		Events:       len(events),
		Sessions:     len(sessions),
		Paths:        len(paths),
		TracedEvents: traced,
		StatusCounts: statusCounts,
		GeneratedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	})
}

func (h *HypatiaServiceImpl) dashboardEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.store.GetRecent(dashboardLimit(r))
	if err != nil {
		writeDashboardError(w, err)
		return
	}
	out := make([]dashboardEvent, 0, len(events))
	for _, event := range events {
		out = append(out, eventForDashboard(event))
	}
	writeDashboardJSON(w, out)
}

func (h *HypatiaServiceImpl) dashboardSessions(w http.ResponseWriter, r *http.Request) {
	events, err := h.store.GetRecent(0)
	if err != nil {
		writeDashboardError(w, err)
		return
	}
	counts := aggregateDashboard(events, func(event *pb.RequestEvent) string { return event.SessionId })
	writeDashboardJSON(w, limitCounts(counts, dashboardLimit(r)))
}

func (h *HypatiaServiceImpl) dashboardPaths(w http.ResponseWriter, r *http.Request) {
	events, err := h.store.GetRecent(0)
	if err != nil {
		writeDashboardError(w, err)
		return
	}
	counts := aggregateDashboard(events, func(event *pb.RequestEvent) string { return event.Path })
	writeDashboardJSON(w, limitCounts(counts, dashboardLimit(r)))
}

func aggregateDashboard(events []*pb.RequestEvent, key func(*pb.RequestEvent) string) []dashboardCount {
	byValue := make(map[string]*dashboardCount)
	for _, event := range events {
		value := key(event)
		if value == "" {
			continue
		}
		entry := byValue[value]
		if entry == nil {
			entry = &dashboardCount{Value: value, LastSeen: event.Timestamp, UserAgent: event.UserAgent}
			byValue[value] = entry
		}
		entry.Count++
	}
	counts := make([]dashboardCount, 0, len(byValue))
	for _, count := range byValue {
		counts = append(counts, *count)
	}
	sort.Slice(counts, func(i, j int) bool {
		if counts[i].Count == counts[j].Count {
			return counts[i].Value < counts[j].Value
		}
		return counts[i].Count > counts[j].Count
	})
	return counts
}

func limitCounts(counts []dashboardCount, limit int) []dashboardCount {
	if len(counts) > limit {
		return counts[:limit]
	}
	return counts
}

func dashboardLimit(r *http.Request) int {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		return defaultDashboardLimit
	}
	if limit > maxDashboardLimit {
		return maxDashboardLimit
	}
	return limit
}

func eventForDashboard(event *pb.RequestEvent) dashboardEvent {
	return dashboardEvent{
		Timestamp: event.Timestamp,
		Path:      event.Path,
		Method:    event.Method,
		Status:    event.Status,
		IP:        event.Ip,
		UserAgent: event.UserAgent,
		Referrer:  event.Referrer,
		SessionID: event.SessionId,
		TraceID:   event.TraceId,
	}
}

func statusGroup(status int32) string {
	if status < 100 {
		return "unknown"
	}
	return strconv.Itoa(int(status/100)) + "xx"
}

func writeDashboardJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		return
	}
}

func writeDashboardError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
