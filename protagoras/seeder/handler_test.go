package seeder

import (
	"context"
	"testing"

	"github.com/odysseia-greek/agora/plato/models"
	dionysiosv1 "github.com/odysseia-greek/alexandreia/dionysios/gen/go/v1"
	"google.golang.org/grpc"
)

type fakeDionysiosClient struct {
	grammarRequests []*dionysiosv1.CheckGrammarRequest
}

func (f *fakeDionysiosClient) Health(context.Context, *dionysiosv1.HealthRequest, ...grpc.CallOption) (*dionysiosv1.HealthResponse, error) {
	return &dionysiosv1.HealthResponse{Healthy: true}, nil
}

func (f *fakeDionysiosClient) CheckGrammar(_ context.Context, request *dionysiosv1.CheckGrammarRequest, _ ...grpc.CallOption) (*dionysiosv1.CheckGrammarResponse, error) {
	f.grammarRequests = append(f.grammarRequests, request)
	return &dionysiosv1.CheckGrammarResponse{Results: []*dionysiosv1.DeclensionResult{{Word: request.Word}}}, nil
}

func TestLoopOverAndDeclineWordsUsesDionysiosGRPC(t *testing.T) {
	client := &fakeDionysiosClient{}
	handler := &ProtagorasHandler{Dionysios: client}
	text := models.Text{Rhemai: []models.Rhema{{Greek: "χαῖρε, κόσμε!"}}}

	if err := handler.loopOverAndDeclineWords(text); err != nil {
		t.Fatalf("loopOverAndDeclineWords returned an error: %v", err)
	}
	if len(client.grammarRequests) != 2 || client.grammarRequests[0].GetWord() != "χαῖρε" || client.grammarRequests[1].GetWord() != "κόσμε" {
		t.Fatalf("unexpected grammar requests: %+v", client.grammarRequests)
	}
	if len(handler.wordsDone) != 2 || len(handler.wordsNotFound) != 0 {
		t.Fatalf("unexpected seeding result: done=%v notFound=%v", handler.wordsDone, handler.wordsNotFound)
	}
}
