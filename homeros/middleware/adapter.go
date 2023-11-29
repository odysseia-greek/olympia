package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/randomizer"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/app"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/odysseia-greek/olympia/homeros/handlers"
	"io"
	"io/ioutil"
	"log"
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
			//allow all CORS w.Header().Set("Access-Control-Allow-Origin", "*")
			allowedOrigin := "localhost"

			origin := r.Header.Get("Origin")
			if strings.Contains(origin, allowedOrigin) {
				logging.Debug(fmt.Sprintf("setting CORS header for origin: %s", origin))
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
				if r.Method == "OPTIONS" {
					return
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
func LogRequestDetails(tracer *aristophanes.ClientTracer, traceConfig *handlers.TraceConfig, randomizer randomizer.Random) Adapter {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes)) // Set the original request body

			var bodyClone map[string]interface{}
			decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
			if err := decoder.Decode(&bodyClone); err != nil {
				http.Error(w, "Failed to parse request body", http.StatusInternalServerError)
				return
			}

			operationName, _ := bodyClone["operationName"].(string)
			query, _ := bodyClone["query"].(string)

			var traceID string
			payload := &pb.StartTraceRequest{
				Method:        r.Method,
				Url:           r.URL.RequestURI(),
				Host:          r.Host,
				RemoteAddress: r.RemoteAddr,
				RootQuery:     query,
				Operation:     operationName,
			}

			for _, service := range traceConfig.OperationScores {
				if service.Operation == operationName && shouldTrace(service.Score, randomizer) {
					trace, err := tracer.StartTrace(context.Background(), payload)
					if err != nil {
						log.Print(err)
						break
					}

					traceID = trace.CombinedId
					go Log(trace.CombinedId) // Integrated log here
					break
				}
			}

			if traceID == "" {
				traceID = uuid.New().String()
			}

			if operationName != "IntrospectionQuery" {
				jsonPayload, err := json.MarshalIndent(payload, "", "  ")
				if err != nil {
					logging.Error(err.Error())
				}

				logLine := fmt.Sprintf("REQUEST | traceId: %s and params:\n%s", traceID, string(jsonPayload))
				logging.Info(logLine)
			}

			ctx := context.WithValue(r.Context(), config.HeaderKey, traceID)
			f.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ... (shouldTrace and Log functions remain the same)

func Log(combinedId string) {
	splitID := strings.Split(combinedId, "+")

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

	logging.Trace(fmt.Sprintf("traceID: %s | parentSpanID: %s | callTraced: %v", traceID, spanID, traceCall))
}

func shouldTrace(score int, random randomizer.Random) bool {
	return random.RandomNumberBaseOne(100) < score
}
