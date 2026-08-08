package monos

import (
	"context"
	"fmt"

	uuid2 "github.com/google/uuid"
	pb "github.com/odysseia-greek/agora/eupalinos/v1"
	"google.golang.org/grpc"
)

// MockEupalinosClient is a mock implementation of EupalinosClient
type MockEupalinosClient struct {
	LastDequeue *pb.ChannelInfo
	LastAck     *pb.AcknowledgeRequest
	LastNack    *pb.NackRequest
	Data        string
	Length      int32
	dequeues    int
}

func (m *MockEupalinosClient) Health(ctx context.Context, in *pb.HealthRequest, opts ...grpc.CallOption) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Healthy: true}, nil
}

func (m *MockEupalinosClient) StreamQueueUpdates(ctx context.Context, opts ...grpc.CallOption) (pb.Eupalinos_StreamQueueUpdatesClient, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockEupalinosClient) EnqueueMessage(ctx context.Context, in *pb.Epistello) (*pb.EnqueueResponse, error) {
	return &pb.EnqueueResponse{}, nil
}

func (m *MockEupalinosClient) EnqueueMessageBytes(ctx context.Context, in *pb.EpistelloBytes, opts ...grpc.CallOption) (*pb.EnqueueResponse, error) {
	return &pb.EnqueueResponse{}, nil
}

// DequeueMessage is the mock implementation for the EnqueueMessage method
func (m *MockEupalinosClient) DequeueMessage(ctx context.Context, in *pb.ChannelInfo) (*pb.Epistello, error) {
	m.LastDequeue = in
	uuid := uuid2.New()
	data := m.Data
	if data == "" {
		data = "{\"method\":\"\",\"category\":\"\",\"greek\":\"Ἄβδηρα\",\"translation\":\"town of Abdera, known for stupidity of inhabitants\",\"chapter\":57}"
	}

	if m.dequeues == 1 {
		return nil, fmt.Errorf("some error")
	}

	m.dequeues++

	return &pb.Epistello{
		Id:      uuid.String(),
		Data:    data,
		Channel: in.Name,
	}, nil
}

func (m *MockEupalinosClient) DequeueMessageBytes(ctx context.Context, in *pb.ChannelInfo, opts ...grpc.CallOption) (*pb.EpistelloBytes, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockEupalinosClient) AcknowledgeMessage(ctx context.Context, in *pb.AcknowledgeRequest) (*pb.AcknowledgeResponse, error) {
	m.LastAck = in
	return &pb.AcknowledgeResponse{Acknowledged: true}, nil
}

func (m *MockEupalinosClient) NackMessage(ctx context.Context, in *pb.NackRequest) (*pb.NackResponse, error) {
	m.LastNack = in
	return &pb.NackResponse{Requeued: true}, nil
}

// EnqueueMessage is the mock implementation for the EnqueueMessage method
func (m *MockEupalinosClient) GetQueueLength(ctx context.Context, in *pb.ChannelInfo) (*pb.QueueLength, error) {
	return &pb.QueueLength{
		Length: m.Length,
	}, nil
}
