package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/odysseia-greek/agora/plato/config"
	"io"
	"net/http"
)

func (h *HomerosHandler) ForwardToSokrates(ctx context.Context) (json.RawMessage, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionID, _ := ctx.Value(config.SessionIdKey).(string)
	requestContext := graphql.GetOperationContext(ctx)
	if requestContext == nil {
		return nil, fmt.Errorf("failed to retrieve GraphQL operation context")
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     requestContext.RawQuery,
		"variables": requestContext.Variables,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", h.SokratesGraphqlUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(config.HeaderKey, requestID)
	req.Header.Set(config.SessionIdKey, sessionID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Sokrates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Sokrates response: %w", err)
	}

	h.CloseTrace(resp, body)

	return body, nil
}
