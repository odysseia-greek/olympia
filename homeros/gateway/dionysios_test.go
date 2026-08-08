package gateway

import (
	"context"
	"testing"

	"github.com/odysseia-greek/agora/plato/service"
	dionysiosv1 "github.com/odysseia-greek/alexandreia/dionysios/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type fakeDionysiosResearchClient struct {
	request         *dionysiosv1.ResearchRequest
	requestID       string
	response        *dionysiosv1.ResearchResponse
	grammarRequest  *dionysiosv1.CheckGrammarRequest
	grammarResponse *dionysiosv1.CheckGrammarResponse
	textRequest     *dionysiosv1.TextModeRequest
	textResponse    *dionysiosv1.TextModeResponse
	healthResponse  *dionysiosv1.HealthResponse
}

func (f *fakeDionysiosResearchClient) Health(_ context.Context, _ *dionysiosv1.HealthRequest, _ ...grpc.CallOption) (*dionysiosv1.HealthResponse, error) {
	return f.healthResponse, nil
}

func (f *fakeDionysiosResearchClient) CheckGrammar(_ context.Context, request *dionysiosv1.CheckGrammarRequest, _ ...grpc.CallOption) (*dionysiosv1.CheckGrammarResponse, error) {
	f.grammarRequest = request
	return f.grammarResponse, nil
}

func (f *fakeDionysiosResearchClient) TextMode(_ context.Context, request *dionysiosv1.TextModeRequest, _ ...grpc.CallOption) (*dionysiosv1.TextModeResponse, error) {
	f.textRequest = request
	return f.textResponse, nil
}

func (f *fakeDionysiosResearchClient) Research(ctx context.Context, request *dionysiosv1.ResearchRequest, _ ...grpc.CallOption) (*dionysiosv1.ResearchResponse, error) {
	f.request = request
	if outgoing, ok := metadata.FromOutgoingContext(ctx); ok {
		values := outgoing.Get(service.HeaderKey)
		if len(values) > 0 {
			f.requestID = values[0]
		}
	}
	return f.response, nil
}

func TestAnalyzeUsesDionysiosResearch(t *testing.T) {
	client := &fakeDionysiosResearchClient{response: &dionysiosv1.ResearchResponse{
		Rootword:     "λόγος",
		Conjugations: []*dionysiosv1.Conjugation{{Word: "λόγου", Rule: "genitive"}},
		Results: []*dionysiosv1.AnalyzeResult{{
			Author: "Plato", Reference: "1.1",
			Text: &dionysiosv1.Rhema{Greek: "λόγος", Translations: []string{"word"}},
		}},
	}}
	handler := &HomerosHandler{Dionysios: client}

	response, err := handler.Analyze(context.Background(), "λόγος", "trace-id")
	if err != nil {
		t.Fatalf("Analyze returned an error: %v", err)
	}
	if client.request.GetRootword() != "λόγος" || client.request.GetLimit() != 5 {
		t.Fatalf("unexpected research request: %+v", client.request)
	}
	if client.requestID != "trace-id" {
		t.Fatalf("expected trace ID to be propagated, got %q", client.requestID)
	}
	if response.Rootword == nil || *response.Rootword != "λόγος" {
		t.Fatalf("unexpected rootword: %v", response.Rootword)
	}
	if len(response.Conjugations) != 1 || response.Conjugations[0].Word == nil || *response.Conjugations[0].Word != "λόγου" {
		t.Fatalf("unexpected conjugations: %+v", response.Conjugations)
	}
	if len(response.Texts) != 1 || response.Texts[0].Text == nil || response.Texts[0].Text.Greek == nil || *response.Texts[0].Text.Greek != "λόγος" {
		t.Fatalf("unexpected text results: %+v", response.Texts)
	}
}

func TestGrammarUsesDionysiosCheckGrammar(t *testing.T) {
	client := &fakeDionysiosResearchClient{grammarResponse: &dionysiosv1.CheckGrammarResponse{
		Results: []*dionysiosv1.DeclensionResult{{
			Word: "λόγου", Rule: "genitive", RootWord: "λόγος", Translation: []string{"word"},
		}},
		Audit: &dionysiosv1.GrammarAudit{Outcome: "success", Events: []*dionysiosv1.GrammarAuditEvent{{Step: "grammar.lookup", ResultCount: 1}}},
	}}
	handler := &HomerosHandler{Dionysios: client}

	response, err := handler.Grammar(context.Background(), "λόγου", "trace-id", true)
	if err != nil {
		t.Fatalf("Grammar returned an error: %v", err)
	}
	if client.grammarRequest.GetWord() != "λόγου" {
		t.Fatalf("unexpected grammar request: %+v", client.grammarRequest)
	}
	if !client.grammarRequest.GetIncludeAudit() {
		t.Fatal("expected grammar audit to be requested")
	}
	if len(response.Results) != 1 || response.Results[0].RootWord == nil || *response.Results[0].RootWord != "λόγος" {
		t.Fatalf("unexpected grammar response: %+v", response.Results)
	}
	if response.Audit == nil || response.Audit.Outcome != "success" || len(response.Audit.Events) != 1 || response.Audit.Events[0].ResultCount != 1 {
		t.Fatalf("unexpected grammar audit: %+v", response.Audit)
	}
}

func TestSentenceUsesDionysiosTextMode(t *testing.T) {
	client := &fakeDionysiosResearchClient{textResponse: &dionysiosv1.TextModeResponse{
		OriginalText: "χαῖρε κόσμε",
		TextSearch: &dionysiosv1.TextSearchStatus{Searched: true, Found: true, MatchCount: 1, Matches: []*dionysiosv1.AnalyzeResult{{
			Author: "Herodotus", Reference: "1.1", Text: &dionysiosv1.Rhema{Greek: "χαῖρε κόσμε"},
		}}},
	}}
	handler := &HomerosHandler{Dionysios: client}

	response, err := handler.Sentence(context.Background(), "χαῖρε κόσμε", "session-id", "trace-id", true)
	if err != nil {
		t.Fatalf("Sentence returned an error: %v", err)
	}
	if client.textRequest.GetText() != "χαῖρε κόσμε" || client.textRequest.GetSessionId() != "session-id" || !client.textRequest.GetIncludeAudit() {
		t.Fatalf("unexpected text-mode request: %+v", client.textRequest)
	}
	if response.OriginalText != "χαῖρε κόσμε" {
		t.Fatalf("unexpected text-mode response: %+v", response)
	}
	if response.TextSearch == nil || !response.TextSearch.Found || len(response.TextSearch.Matches) != 1 || response.TextSearch.Matches[0].Author == nil || *response.TextSearch.Matches[0].Author != "Herodotus" {
		t.Fatalf("unexpected text search response: %+v", response.TextSearch)
	}
}
