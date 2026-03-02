package monos

import (
	"context"
	"fmt"
	"os"
	"strconv"

	uuid2 "github.com/google/uuid"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"google.golang.org/grpc"
)

const (
	TestLength = "TEST_LENGTH"
	TestData   = "TEST_DATA"
)

var NumberOfDequeue int

// MockEupalinosClient is a mock implementation of EupalinosClient
type MockEupalinosClient struct {
}

func (m *MockEupalinosClient) Health(ctx context.Context, in *pb.HealthRequest, opts ...grpc.CallOption) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Healthy: true}, nil
}

func (m *MockEupalinosClient) StreamQueueUpdates(ctx context.Context, opts ...grpc.CallOption) (pb.Eupalinos_StreamQueueUpdatesClient, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockEupalinosClient) EnqueueMessage(ctx context.Context, in *pb.Epistello, opts ...grpc.CallOption) (*pb.EnqueueResponse, error) {
	return &pb.EnqueueResponse{}, nil
}

func (m *MockEupalinosClient) EnqueueMessageBytes(ctx context.Context, in *pb.EpistelloBytes, opts ...grpc.CallOption) (*pb.EnqueueResponse, error) {
	return &pb.EnqueueResponse{}, nil
}

// DequeueMessage is the mock implementation for the EnqueueMessage method
func (m *MockEupalinosClient) DequeueMessage(ctx context.Context, in *pb.ChannelInfo, opts ...grpc.CallOption) (*pb.Epistello, error) {
	uuid := uuid2.New()
	data := os.Getenv(TestData)
	if data == "" {
		data = "{\"method\":\"\",\"category\":\"\",\"greek\":\"Ἄβδηρα\",\"translation\":\"town of Abdera, known for stupidity of inhabitants\",\"chapter\":57}"
	}

	if NumberOfDequeue == 1 {
		return nil, fmt.Errorf("some error")
	}

	NumberOfDequeue++

	return &pb.Epistello{
		Id:      uuid.String(),
		Data:    data,
		Channel: in.Name,
	}, nil
}

func (m *MockEupalinosClient) DequeueMessageBytes(ctx context.Context, in *pb.ChannelInfo, opts ...grpc.CallOption) (*pb.EpistelloBytes, error) {
	return nil, fmt.Errorf("not implemented")
}

// EnqueueMessage is the mock implementation for the EnqueueMessage method
func (m *MockEupalinosClient) GetQueueLength(ctx context.Context, in *pb.ChannelInfo, opts ...grpc.CallOption) (*pb.QueueLength, error) {
	var length int32
	lengthAsString := os.Getenv(TestLength)
	if lengthAsString == "" {
		length = 1
	} else {
		l, _ := strconv.Atoi(lengthAsString)
		length = int32(l)
	}
	return &pb.QueueLength{
		Length: length,
	}, nil
}
