package mouseion

import (
	"context"
	"time"

	pb "github.com/odysseia-greek/olympia/hypatia/proto/v1"
)

func (h *HypatiaServiceImpl) Health(ctx context.Context, _ *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Healthy: true,
		Time:    time.Now().String(),
		Version: h.Version,
	}, nil
}

func (h *HypatiaServiceImpl) TrackEvents(ctx context.Context, in *pb.RequestEventBatch) (*pb.TrackResponse, error) {
	if err := h.store.Add(in.Events); err != nil {
		return nil, err
	}
	return &pb.TrackResponse{Ack: "received"}, nil
}

func (h *HypatiaServiceImpl) GetEventsBySession(ctx context.Context, in *pb.SessionRequest) (*pb.EventsResponse, error) {
	events, err := h.store.GetBySession(in.SessionId)
	if err != nil {
		return nil, err
	}
	return &pb.EventsResponse{Events: events}, nil
}

func (h *HypatiaServiceImpl) GetEventsByPath(ctx context.Context, in *pb.PathRequest) (*pb.EventsResponse, error) {
	events, err := h.store.GetByPath(in.Path)
	if err != nil {
		return nil, err
	}
	return &pb.EventsResponse{Events: events}, nil
}

func (h *HypatiaServiceImpl) GetRecentEvents(ctx context.Context, in *pb.RecentRequest) (*pb.EventsResponse, error) {
	events, err := h.store.GetRecent(int(in.Limit))
	if err != nil {
		return nil, err
	}
	return &pb.EventsResponse{Events: events}, nil
}
