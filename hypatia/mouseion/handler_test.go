package mouseion

import (
	"context"
	"testing"

	"github.com/odysseia-greek/agora/plato/config"
	pb "github.com/odysseia-greek/olympia/hypatia/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestTrackEventsAddsTraceIDFromContext(t *testing.T) {
	service := &HypatiaServiceImpl{store: NewInMemoryStore()}
	ctx := context.WithValue(context.Background(), config.DefaultTracingName, "trace-123+span-456+1")
	event := &pb.RequestEvent{Path: "/search", SessionId: "session-1"}

	if _, err := service.TrackEvents(ctx, &pb.RequestEventBatch{Events: []*pb.RequestEvent{event}}); err != nil {
		t.Fatalf("TrackEvents() error = %v", err)
	}
	if event.TraceId != "" {
		t.Fatalf("TrackEvents() mutated caller event trace_id to %q", event.TraceId)
	}

	response, err := service.GetEventsBySession(context.Background(), &pb.SessionRequest{SessionId: "session-1"})
	if err != nil {
		t.Fatalf("GetEventsBySession() error = %v", err)
	}
	if got := response.Events[0].TraceId; got != "trace-123" {
		t.Fatalf("trace_id = %q, want %q", got, "trace-123")
	}
}

func TestTrackEventsPreservesExplicitTraceID(t *testing.T) {
	service := &HypatiaServiceImpl{store: NewInMemoryStore()}
	ctx := context.WithValue(context.Background(), config.DefaultTracingName, "context-trace+span-456+1")
	event := &pb.RequestEvent{Path: "/search", TraceId: "event-trace"}

	if _, err := service.TrackEvents(ctx, &pb.RequestEventBatch{Events: []*pb.RequestEvent{event}}); err != nil {
		t.Fatalf("TrackEvents() error = %v", err)
	}

	response, err := service.GetEventsByPath(context.Background(), &pb.PathRequest{Path: "/search"})
	if err != nil {
		t.Fatalf("GetEventsByPath() error = %v", err)
	}
	if got := response.Events[0].TraceId; got != "event-trace" {
		t.Fatalf("trace_id = %q, want %q", got, "event-trace")
	}
}

func TestHandlersRejectInvalidRequests(t *testing.T) {
	service := &HypatiaServiceImpl{store: NewInMemoryStore()}
	tests := []struct {
		name string
		call func() error
	}{
		{name: "empty batch", call: func() error { _, err := service.TrackEvents(context.Background(), &pb.RequestEventBatch{}); return err }},
		{name: "nil event", call: func() error {
			_, err := service.TrackEvents(context.Background(), &pb.RequestEventBatch{Events: []*pb.RequestEvent{nil}})
			return err
		}},
		{name: "empty session", call: func() error {
			_, err := service.GetEventsBySession(context.Background(), &pb.SessionRequest{})
			return err
		}},
		{name: "empty path", call: func() error { _, err := service.GetEventsByPath(context.Background(), &pb.PathRequest{}); return err }},
		{name: "negative limit", call: func() error {
			_, err := service.GetRecentEvents(context.Background(), &pb.RecentRequest{Limit: -1})
			return err
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := status.Code(test.call()); got != codes.InvalidArgument {
				t.Fatalf("status code = %s, want %s", got, codes.InvalidArgument)
			}
		})
	}
}
