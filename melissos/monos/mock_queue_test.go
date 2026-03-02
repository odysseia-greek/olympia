package monos

import (
	"unsafe"

	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/eupalinos/stomion"
)

type queueClientAlias struct {
	queue pb.EupalinosClient
}

func newMockQueueClient(mock pb.EupalinosClient) *stomion.QueueClient {
	alias := &queueClientAlias{queue: mock}
	return (*stomion.QueueClient)(unsafe.Pointer(alias))
}
