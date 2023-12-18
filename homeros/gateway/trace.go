package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/proto"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func ParseHeaderID(requestId string) (string, string, bool) {
	splitID := strings.Split(requestId, "+")

	traceCall := false
	var traceID, spanID string

	if len(splitID) >= 3 {
		traceCall = splitID[2] == "1"
	}

	if len(splitID) >= 1 {
		traceID = splitID[0]
	}
	if len(splitID) >= 2 {
		spanID = splitID[1]
	}

	return traceID, spanID, traceCall
}

func (h *HomerosHandler) CloseTrace(response *http.Response, body interface{}) {
	requestId := response.Header.Get(config.HeaderKey)

	traceId, parentspanId, traceCall := ParseHeaderID(requestId)

	jsonBody, _ := json.Marshal(body)

	logging.Info(fmt.Sprintf("RESPONSE | traceID: %s | responseCode: %d", traceId, response.StatusCode))

	if traceCall {
		traceCloser := &aristophanes.CloseTraceRequest{
			TraceId:      traceId,
			ParentSpanId: parentspanId,
			ResponseCode: int32(response.StatusCode),
			ResponseBody: string(jsonBody),
		}

		h.Tracer.CloseTrace(context.Background(), traceCloser)
		logging.Trace(fmt.Sprintf("trace closed with id: %s", traceId))
	}
}

// todo refactor to actual have a meaningful code this has to be changed in the plato client
func (h *HomerosHandler) CloseTraceWithError(err error, requestId string) {
	traceID, parentspanId, traceCall := ParseHeaderID(requestId)
	statusCode := parseActualStatusCodeFromErrorMessage(err)

	logging.Error(fmt.Sprintf("RESPONSE | traceID: %s | responseCode: %d | error: %s", traceID, statusCode, err.Error()))

	if traceCall {
		traceCloser := &aristophanes.CloseTraceRequest{
			TraceId:      traceID,
			ParentSpanId: parentspanId,
			ResponseCode: int32(statusCode),
			ResponseBody: err.Error(),
		}

		h.Tracer.CloseTrace(context.Background(), traceCloser)
		logging.Trace(fmt.Sprintf("trace closed with id: %s", traceID))
	}
}

func parseActualStatusCodeFromErrorMessage(err error) int {
	// Define a regular expression to extract the actual status code from the error message
	re := regexp.MustCompile(`got (\d+)`)

	// Use the regular expression to find matches in the error message
	matches := re.FindStringSubmatch(err.Error())

	// If matches are found, parse the actual status code
	if len(matches) > 1 {
		actualStatusCode, parseErr := strconv.Atoi(matches[1])
		if parseErr != nil {
			return 0
		}

		return actualStatusCode
	}

	return 0
}
