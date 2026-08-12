package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/odysseia-greek/agora/plato/config"
)

type graphqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// QueryHerodotos calls the Ionia Herodotos GraphQL API and decodes one root
// field into target. Keeping the upstream query here lets Homeros expose names
// such as corpusText without requiring the upstream service to rename text.
func (h *HomerosHandler) QueryHerodotos(ctx context.Context, query string, variables map[string]any, field string, target any) error {
	payload, err := json.Marshal(map[string]any{"query": query, "variables": variables})
	if err != nil {
		return fmt.Errorf("marshal Herodotos request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.HerodotosGraphqlUrl, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create Herodotos request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if requestID, ok := ctx.Value(config.HeaderKey).(string); ok {
		req.Header.Set(config.HeaderKey, requestID)
	}
	if sessionID, ok := ctx.Value(config.SessionIdKey).(string); ok {
		req.Header.Set(config.SessionIdKey, sessionID)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("call Herodotos: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("Herodotos returned %s", response.Status)
	}

	var envelope graphqlResponse
	if err := json.NewDecoder(response.Body).Decode(&envelope); err != nil {
		return fmt.Errorf("decode Herodotos response: %w", err)
	}
	if len(envelope.Errors) != 0 {
		return fmt.Errorf("Herodotos GraphQL error: %s", envelope.Errors[0].Message)
	}
	var data map[string]json.RawMessage
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		return fmt.Errorf("decode Herodotos data: %w", err)
	}
	value, ok := data[field]
	if !ok {
		return fmt.Errorf("Herodotos response is missing %q", field)
	}
	if err := json.Unmarshal(value, target); err != nil {
		return fmt.Errorf("decode Herodotos %s: %w", field, err)
	}
	return nil
}
