package mouseion

import (
	"context"
	"time"

	"github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/olympia/hypatia/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func (h *HypatiaServiceImpl) Health(ctx context.Context, _ *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Healthy: true,
		Time:    time.Now().UTC().Format(time.RFC3339Nano),
		Version: h.Version,
	}, nil
}

func (h *HypatiaServiceImpl) TrackEvents(ctx context.Context, in *pb.RequestEventBatch) (*pb.TrackResponse, error) {
	if in == nil || len(in.GetEvents()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one event is required")
	}

	traceID, _, sampled := comedy.ExtractRequestIds(ctx)
	events := make([]*pb.RequestEvent, 0, len(in.Events))
	for _, event := range in.Events {
		if event == nil {
			continue
		}
		event = proto.Clone(event).(*pb.RequestEvent)
		if sampled && traceID != "" && event.TraceId == "" {
			event.TraceId = traceID
		}
		events = append(events, event)
	}
	if len(events) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one non-nil event is required")
	}

	if err := h.store.Add(events); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.TrackResponse{Ack: "received"}, nil
}

func (h *HypatiaServiceImpl) GetEventsBySession(ctx context.Context, in *pb.SessionRequest) (*pb.EventsResponse, error) {
	if in == nil || in.GetSessionId() == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	events, err := h.store.GetBySession(in.SessionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EventsResponse{Events: events}, nil
}

func (h *HypatiaServiceImpl) GetEventsByPath(ctx context.Context, in *pb.PathRequest) (*pb.EventsResponse, error) {
	if in == nil || in.GetPath() == "" {
		return nil, status.Error(codes.InvalidArgument, "path is required")
	}
	events, err := h.store.GetByPath(in.Path)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EventsResponse{Events: events}, nil
}

func (h *HypatiaServiceImpl) GetRecentEvents(ctx context.Context, in *pb.RecentRequest) (*pb.EventsResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if in.GetLimit() < 0 {
		return nil, status.Error(codes.InvalidArgument, "limit cannot be negative")
	}
	events, err := h.store.GetRecent(int(in.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EventsResponse{Events: events}, nil
}
