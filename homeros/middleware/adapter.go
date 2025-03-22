package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/odysseia-greek/olympia/homeros/gateway"
	"io"
	"net/http"
	"strings"

	"github.com/odysseia-greek/agora/plato/logging"
)

type Adapter func(http.Handler) http.Handler

// Adapt Iterate over adapters and run them one by one
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func SetCorsHeaders() Adapter {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			// Allow any origin in the allowedOrigins slice
			allowedOrigins := []string{"localhost"}

			for _, o := range allowedOrigins {
				if strings.Contains(origin, o) {
					logging.Debug(fmt.Sprintf("setting CORS header for origin: %s", origin))
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					if r.Method == "OPTIONS" {
						w.WriteHeader(http.StatusOK)
						return
					}
					break
				}
			}
			f.ServeHTTP(w, r)
		})
	}
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
func LogRequestDetails(tracer pb.TraceService_ChorusClient, traceConfig *gateway.TraceConfig, randomizer randomizer.Random) Adapter {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			payload := &pb.StartTraceRequest{
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
						parabasis := &pb.ParabasisRequest{
							TraceId:      traceID,
							ParentSpanId: spanID,
							SpanId:       spanID,
							RequestType: &pb.ParabasisRequest_StartTrace{
								StartTrace: payload,
							},
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

			if operationName != "IntrospectionQuery" {
				jsonPayload, err := json.MarshalIndent(payload, "", "  ")
				if err != nil {
					logging.Error(err.Error())
				}

				logLine := fmt.Sprintf("REQUEST | traceId: %s and params:\n%s", traceID, string(jsonPayload))
				logging.Info(logLine)
			}

			ctx := context.WithValue(r.Context(), config.HeaderKey, requestID)
			f.ServeHTTP(w, r.WithContext(ctx))
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
