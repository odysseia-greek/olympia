package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
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
	traceID, parentspanID, traceCall := ParseHeaderID(requestId)

	var jsonBody []byte
	var err error

	// Handle body based on its type
	switch v := body.(type) {
	case []byte:
		jsonBody = v
	case json.RawMessage:
		jsonBody = []byte(v)
	default:
		jsonBody, err = json.Marshal(v)
		if err != nil {
			logging.Error(fmt.Sprintf("failed to marshal response body: %v", err))
			jsonBody = []byte(`{"error": "failed to serialize response body"}`)
		}
	}

	logging.Info(fmt.Sprintf("RESPONSE | traceID: %s | responseCode: %d", traceID, response.StatusCode))

	if traceCall {
		parabasis := &pb.ParabasisRequest{
			TraceId:      traceID,
			ParentSpanId: parentspanID,
			SpanId:       parentspanID,
			RequestType: &pb.ParabasisRequest_CloseTrace{
				CloseTrace: &pb.CloseTraceRequest{
					ResponseCode: int32(response.StatusCode),
					ResponseBody: string(jsonBody),
				},
			},
		}

		err := h.Streamer.Send(parabasis)
		if err != nil {
			logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
			return
		}

		logging.Trace(fmt.Sprintf("trace closed with id: %s", traceID))
	}
}

func (h *HomerosHandler) CloseTraceWithError(err error, requestId string) {
	traceID, parentspanID, traceCall := ParseHeaderID(requestId)
	statusCode := parseActualStatusCodeFromErrorMessage(err)

	logging.Error(fmt.Sprintf("RESPONSE | traceID: %s | responseCode: %d | error: %s", traceID, statusCode, err.Error()))

	if traceCall {
		parabasis := &pb.ParabasisRequest{
			TraceId:      traceID,
			ParentSpanId: parentspanID,
			SpanId:       parentspanID,
			RequestType: &pb.ParabasisRequest_CloseTrace{
				CloseTrace: &pb.CloseTraceRequest{
					ResponseCode: int32(statusCode),
					ResponseBody: err.Error(),
				},
			},
		}

		err := h.Streamer.Send(parabasis)
		if err != nil {
			logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
			return
		}

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
