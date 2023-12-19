package scholar

import (
	"context"
	"github.com/odysseia-greek/agora/aristoteles"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/olympia/aristarchos/proto"
	"google.golang.org/grpc"
	"time"
)

type AggregatorService interface {
	WaitForHealthyState() bool
	CreateNewEntry(ctx context.Context, in *pb.AggregatorCreationRequest) (*pb.AggregatorCreationResponse, error)
	RetrieveEntry(ctx context.Context, request *pb.AggregatorRequest) (*pb.RootWordResponse, error)
	RetrieveSearchWords(ctx context.Context, in *pb.AggregatorRequest) (*pb.SearchWordResponse, error)
}

const (
	DEFAULTADDRESS string = "localhost:50053"
)

type AggregatorServiceImpl struct {
	Elastic aristoteles.Client
	Index   string
	Tracer  *aristophanes.ClientTracer
	pb.UnimplementedAristarchosServer
}

type AggregatorServiceClient struct {
	Impl AggregatorService
}

type ClientAggregator struct {
	scholar pb.AristarchosClient
}

func NewClientAggregator() *ClientAggregator {
	conn, _ := grpc.Dial(DEFAULTADDRESS, grpc.WithInsecure())
	client := pb.NewAristarchosClient(conn)
	return &ClientAggregator{scholar: client}
}

func (c *ClientAggregator) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := c.Health(context.Background(), &pb.HealthRequest{})
		if err == nil && response.Health {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (c *ClientAggregator) Health(ctx context.Context, request *pb.HealthRequest) (*pb.HealthResponse, error) {
	return c.scholar.Health(ctx, request)
}

func (c *ClientAggregator) CreateNewEntry(ctx context.Context, request *pb.AggregatorCreationRequest) (*pb.AggregatorCreationResponse, error) {
	return c.scholar.CreateNewEntry(ctx, request)
}

func (c *ClientAggregator) RetrieveEntry(ctx context.Context, request *pb.AggregatorRequest) (*pb.RootWordResponse, error) {
	return c.scholar.RetrieveEntry(ctx, request)
}

func (c *ClientAggregator) RetrieveSearchWords(ctx context.Context, request *pb.AggregatorRequest) (*pb.SearchWordResponse, error) {
	return c.scholar.RetrieveSearchWords(ctx, request)
}
