package scholar

import (
	"context"
	pb "github.com/odysseia-greek/olympia/aristarchos/proto"
	"github.com/stretchr/testify/mock"
	"time"
)

// MockTraceService is a mock implementation of the TraceService interface
type MockTraceService struct {
	mock.Mock
}

func (m *MockTraceService) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := m.Health(context.Background(), &pb.HealthRequest{})
		if err == nil && response.Health {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (m *MockTraceService) Health(ctx context.Context, request *pb.HealthRequest) (*pb.HealthResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*pb.HealthResponse), args.Error(1)
}
