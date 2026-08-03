package mouseion

import (
	"context"
	"fmt"
	"time"

	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	pb "github.com/odysseia-greek/olympia/hypatia/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	DefaultAddress string = "localhost:50061"
)

// HypatiaService is the interface for the server-side implementation.
type HypatiaService interface {
	TrackEvents(ctx context.Context, in *pb.RequestEventBatch) (*pb.TrackResponse, error)
	GetEventsBySession(ctx context.Context, in *pb.SessionRequest) (*pb.EventsResponse, error)
	GetEventsByPath(ctx context.Context, in *pb.PathRequest) (*pb.EventsResponse, error)
	GetRecentEvents(ctx context.Context, in *pb.RecentRequest) (*pb.EventsResponse, error)
	Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthResponse, error)
}

// HypatiaServiceImpl holds the runtime dependencies.
type HypatiaServiceImpl struct {
	Version  string
	store    EventStore
	Streamer arv1.TraceService_ChorusClient
	pb.UnimplementedHypatiaServer
}

// HypatiaClient wraps the gRPC client for use by other services.
type HypatiaClient struct {
	client pb.HypatiaClient
	conn   *grpc.ClientConn
}

func NewHypatiaClient(address string) (*HypatiaClient, error) {
	if address == "" {
		address = DefaultAddress
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to hypatia: %w", err)
	}
	return &HypatiaClient{client: pb.NewHypatiaClient(conn), conn: conn}, nil
}

// Close releases the underlying gRPC connection.
func (c *HypatiaClient) Close() error {
	return c.conn.Close()
}

func (c *HypatiaClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := c.client.Health(context.Background(), &pb.HealthRequest{})
		if err == nil && resp.Healthy {
			return true
		}
		time.Sleep(checkInterval)
	}
	return false
}

func (c *HypatiaClient) TrackEvents(ctx context.Context, in *pb.RequestEventBatch) (*pb.TrackResponse, error) {
	return c.client.TrackEvents(ctx, in)
}

func (c *HypatiaClient) GetEventsBySession(ctx context.Context, in *pb.SessionRequest) (*pb.EventsResponse, error) {
	return c.client.GetEventsBySession(ctx, in)
}

func (c *HypatiaClient) GetEventsByPath(ctx context.Context, in *pb.PathRequest) (*pb.EventsResponse, error) {
	return c.client.GetEventsByPath(ctx, in)
}

func (c *HypatiaClient) GetRecentEvents(ctx context.Context, in *pb.RecentRequest) (*pb.EventsResponse, error) {
	return c.client.GetRecentEvents(ctx, in)
}

func (c *HypatiaClient) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthResponse, error) {
	return c.client.Health(ctx, in)
}
