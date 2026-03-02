package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	v1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"github.com/odysseia-greek/olympia/homeros/gateway"

	"github.com/odysseia-greek/agora/plato/logging"
)

type bodyRecorder struct {
	http.ResponseWriter
	status int
	buf    bytes.Buffer
	limit  int // max bytes to store (0 = no limit)
}

func (r *bodyRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *bodyRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}

	// capture body (bounded)
	if r.limit == 0 {
		r.buf.Write(b)
	} else if r.buf.Len() < r.limit {
		remain := r.limit - r.buf.Len()
		if len(b) > remain {
			r.buf.Write(b[:remain])
		} else {
			r.buf.Write(b)
		}
	}

	return r.ResponseWriter.Write(b)
}

// LogRequestDetails is a middleware function that captures and logs details of incoming requests,
// and initiates traces based on the configured trace probabilities for specific GraphQL operations.
// It reads the incoming request body to extract the operation name and query from GraphQL requests.
// The middleware then checks the trace configuration to determine whether to initiate a trace for
// the given operation. If the trace probability condition is met, a trace is started using the
// provided tracer's StartTrace method. The trace ID is logged, and the middleware creates a new
// context with the trace ID to pass it along to downstream handlers.
//
// Parameters:
// - tracer: The tracer instance used to initiate traces.
// - traceConfig: The configuration specifying the trace probabilities for specific operations.
//
// Returns:
// An Adapter that wraps an http.Handler and performs the described middleware actions.
func LogRequestDetails(tracer v1.TraceService_ChorusClient, traceConfig *gateway.TraceConfig, randomizer randomizer.Random) middleware.GraphqlAdapter {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionId := r.Header.Get(config.SessionIdKey)
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // Set the original request body

			var bodyClone map[string]interface{}
			decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
			if err := decoder.Decode(&bodyClone); err != nil {
				http.Error(w, "Failed to parse request body", http.StatusInternalServerError)
				return
			}

			operationName, _ := bodyClone["operationName"].(string)
			query, _ := bodyClone["query"].(string)

			if operationName == "" {
				splitQuery := strings.Split(query, "{")
				if len(splitQuery) != 0 {
					if strings.Contains(splitQuery[1], "(") {
						splitsStringPart := strings.Split(splitQuery[1], "(")[0]
						operationName = strings.TrimSpace(splitsStringPart)
						logging.Debug(fmt.Sprintf("extracted operationName from query: %s", operationName))
					}
				}
			}

			traceID := uuid.New().String()
			spanID := comedy.GenerateSpanID()
			traceRequest := 0

			payload := &v1.ObserveTraceStart{
				Method:        r.Method,
				Url:           r.URL.RequestURI(),
				Host:          r.Host,
				RemoteAddress: getRealIP(r),
				RootQuery:     query,
				Operation:     operationName,
			}

			for _, service := range traceConfig.OperationScores {
				if service.Operation == operationName && shouldTrace(service.Score, randomizer) {
					traceRequest = 1

					go func() {
						parabasis := &v1.ObserveRequest{
							TraceId:      traceID,
							ParentSpanId: spanID,
							SpanId:       spanID,
							Kind:         &v1.ObserveRequest_TraceStart{TraceStart: payload},
						}

						if err := tracer.Send(parabasis); err != nil {
							logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
						}

						logging.Trace(fmt.Sprintf("trace with requestID: %s and span: %s", traceID, spanID))
					}()
					break
				}
			}
			requestID := fmt.Sprintf("%s+%s+%d", traceID, spanID, traceRequest)

			// set headers + ctx as you do
			w.Header().Set(config.HeaderKey, requestID)
			w.Header().Set(config.SessionIdKey, sessionId)
			ctx := context.WithValue(r.Context(), config.HeaderKey, requestID)
			ctx = context.WithValue(ctx, config.SessionIdKey, sessionId)

			rec := &bodyRecorder{
				ResponseWriter: w,
				limit:          64 * 1024, // 64KB cap so you don't blow memory on giant responses
			}

			start := time.Now()
			f.ServeHTTP(rec, r.WithContext(ctx))
			dur := time.Since(start)

			status := rec.status
			if status == 0 {
				status = http.StatusOK
			}

			if traceRequest == 1 {
				stop := &v1.ObserveRequest{
					TraceId:      traceID,
					SpanId:       spanID, // root span
					ParentSpanId: spanID, // root parent == root
					Kind: &v1.ObserveRequest_TraceStop{
						TraceStop: &v1.ObserveTraceStop{
							ResponseBody: rec.buf.String(),
							ResponseCode: int32(status),
						},
					},
				}

				if err := tracer.Send(stop); err != nil {
					logging.Error(fmt.Sprintf("failed to send trace stop: %v", err))
				}

				logging.Trace(fmt.Sprintf(
					"trace closed | traceId=%s span=%s status=%d tookMs=%d",
					traceID, spanID, status, dur.Milliseconds(),
				))
			}
		})
	}
}

func shouldTrace(score int, random randomizer.Random) bool {
	return random.RandomNumberBaseOne(100) < score
}

func getRealIP(r *http.Request) string {
	// Check if the X-Real-IP header is set by Traefik or another proxy
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// If X-Real-IP is not present, check the X-Forwarded-For header
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// The X-Forwarded-For header can contain a comma-separated list of IP addresses.
		// The left-most IP address is the original client IP.
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// If neither header is present, fall back to the standard remote address
	return r.RemoteAddr
}
