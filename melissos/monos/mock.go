package monos

import (
	"context"
	"fmt"
	uuid2 "github.com/google/uuid"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"google.golang.org/grpc"
	"os"
	"strconv"
)

const (
	TestLength = "TEST_LENGTH"
	TestData   = "TEST_DATA"
)

var NumberOfDequeue int

// MockEupalinosClient is a mock implementation of EupalinosClient
type MockEupalinosClient struct {
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
