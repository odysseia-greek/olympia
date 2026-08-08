package gateway

import (
	"context"
	"fmt"

	"github.com/odysseia-greek/agora/plato/service"
	dionysiosv1 "github.com/odysseia-greek/alexandreia/dionysios/gen/go/v1"
	"github.com/odysseia-greek/olympia/homeros/graph/model"
	"google.golang.org/grpc/metadata"
)

func (h *HomerosHandler) Grammar(ctx context.Context, word, requestID string, includeAudit bool) (*model.DeclensionTranslationResult, error) {
	if h.Dionysios == nil {
		return nil, fmt.Errorf("dionysios client is not configured")
	}

	response, err := h.Dionysios.CheckGrammar(dionysiosContext(ctx, requestID), &dionysiosv1.CheckGrammarRequest{Word: word, IncludeAudit: includeAudit})
	if err != nil {
		return nil, fmt.Errorf("dionysios grammar check failed: %w", err)
	}

	return mapGrammarResponse(response), nil
}

// Analyze uses Dionysios Research as the single source for grammatical forms
// and occurrences in the corpus. Herodotos remains responsible for text
// creation and translation checking only.
func (h *HomerosHandler) Analyze(ctx context.Context, rootword, requestID string) (*model.AnalyzeTextResponse, error) {
	if h.Dionysios == nil {
		return nil, fmt.Errorf("dionysios research client is not configured")
	}

	response, err := h.Dionysios.Research(dionysiosContext(ctx, requestID), &dionysiosv1.ResearchRequest{
		Rootword: rootword,
		Limit:    5,
	})
	if err != nil {
		return nil, fmt.Errorf("dionysios research failed: %w", err)
	}

	return mapResearchResponse(response), nil
}

func (h *HomerosHandler) Sentence(ctx context.Context, text, sessionID, requestID string, includeAudit bool) (*model.DionysiosSentenceResponse, error) {
	if h.Dionysios == nil {
		return nil, fmt.Errorf("dionysios client is not configured")
	}

	response, err := h.Dionysios.TextMode(dionysiosContext(ctx, requestID), &dionysiosv1.TextModeRequest{
		Text: text, SessionId: sessionID, IncludeAudit: includeAudit,
	})
	if err != nil {
		return nil, fmt.Errorf("dionysios sentence analysis failed: %w", err)
	}

	return mapSentenceResponse(response), nil
}

func dionysiosContext(ctx context.Context, requestID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, service.HeaderKey, requestID)
}

func mapGrammarResponse(response *dionysiosv1.CheckGrammarResponse) *model.DeclensionTranslationResult {
	result := &model.DeclensionTranslationResult{}
	if response == nil {
		return result
	}
	result.Results = make([]*model.Result, 0, len(response.Results))
	for _, grammarResult := range response.Results {
		if grammarResult == nil {
			continue
		}
		result.Results = append(result.Results, &model.Result{
			Word: stringPointer(grammarResult.Word), Rule: stringPointer(grammarResult.Rule),
			RootWord: stringPointer(grammarResult.RootWord), Translations: stringPointers(grammarResult.Translation),
		})
	}
	result.Audit = mapGrammarAudit(response.Audit)
	return result
}

func mapGrammarAudit(audit *dionysiosv1.GrammarAudit) *model.GrammarAudit {
	if audit == nil {
		return nil
	}

	result := &model.GrammarAudit{
		RequestID: audit.RequestId, Word: audit.Word, Outcome: audit.Outcome,
		Source: audit.Source, Reason: audit.Reason, Events: []*model.GrammarAuditEvent{},
	}
	for _, event := range audit.Events {
		if event == nil {
			continue
		}
		result.Events = append(result.Events, &model.GrammarAuditEvent{
			Step: event.Step, Status: event.Status, Reason: event.Reason, Source: event.Source,
			Rule: event.Rule, RootWord: event.RootWord, SearchTerm: event.SearchTerm,
			ResultCount: int32(event.ResultCount), CandidateCount: int32(event.CandidateCount), Details: event.Details,
		})
	}
	return result
}

func mapSentenceResponse(response *dionysiosv1.TextModeResponse) *model.DionysiosSentenceResponse {
	result := &model.DionysiosSentenceResponse{Tokens: []*model.DionysiosTextToken{}}
	if response == nil {
		return result
	}

	result.SessionID = response.SessionId
	result.OriginalText = response.OriginalText
	result.LiteralTranslation = response.LiteralTranslation
	if response.RateLimit != nil {
		result.RateLimit = &model.DionysiosRateLimit{
			UpstreamIP: response.RateLimit.UpstreamIp, WindowSeconds: int32(response.RateLimit.WindowSeconds), NextAllowedTime: response.RateLimit.NextAllowedTime,
		}
	}
	for _, token := range response.Tokens {
		if token == nil {
			continue
		}
		result.Tokens = append(result.Tokens, &model.DionysiosTextToken{
			Token: token.Token, Position: int32(token.Position), Gloss: token.Gloss, Resolved: token.Resolved, Message: token.Message,
			Results: mapGrammarResponse(&dionysiosv1.CheckGrammarResponse{Results: token.Results}).Results,
		})
	}
	if response.TextSearch != nil {
		result.TextSearch = &model.DionysiosTextSearch{
			Searched: response.TextSearch.Searched, Found: response.TextSearch.Found,
			Status: response.TextSearch.Status, Message: response.TextSearch.Message,
			MatchCount: int32(response.TextSearch.MatchCount), Query: response.TextSearch.Query,
			Matches: mapAnalyzeResults(response.TextSearch.Matches),
		}
	}
	return result
}

func mapResearchResponse(response *dionysiosv1.ResearchResponse) *model.AnalyzeTextResponse {
	if response == nil {
		return &model.AnalyzeTextResponse{}
	}

	result := &model.AnalyzeTextResponse{Rootword: stringPointer(response.Rootword)}
	result.Conjugations = make([]*model.ConjugationResponse, 0, len(response.Conjugations))
	for _, conjugation := range response.Conjugations {
		if conjugation == nil {
			continue
		}
		result.Conjugations = append(result.Conjugations, &model.ConjugationResponse{
			Word: stringPointer(conjugation.Word),
			Rule: stringPointer(conjugation.Rule),
		})
	}

	result.Texts = mapAnalyzeResults(response.Results)

	return result
}

func mapAnalyzeResults(results []*dionysiosv1.AnalyzeResult) []*model.AnalyzeResult {
	mappedResults := make([]*model.AnalyzeResult, 0, len(results))
	for _, researchResult := range results {
		if researchResult == nil {
			continue
		}
		mapped := &model.AnalyzeResult{
			ReferenceLink: stringPointer(researchResult.ReferenceLink),
			Author:        stringPointer(researchResult.Author),
			Book:          stringPointer(researchResult.Book),
			Reference:     stringPointer(researchResult.Reference),
		}
		if researchResult.Text != nil {
			mapped.Text = &model.Rhema{
				Greek:        stringPointer(researchResult.Text.Greek),
				Section:      stringPointer(researchResult.Text.Section),
				Translations: stringPointers(researchResult.Text.Translations),
			}
		}
		mappedResults = append(mappedResults, mapped)
	}
	return mappedResults
}

func stringPointer(value string) *string {
	return &value
}

func stringPointers(values []string) []*string {
	result := make([]*string, 0, len(values))
	for _, value := range values {
		result = append(result, stringPointer(value))
	}
	return result
}
