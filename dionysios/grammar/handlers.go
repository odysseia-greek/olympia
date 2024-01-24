package grammar

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	plato "github.com/odysseia-greek/agora/plato/service"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	pba "github.com/odysseia-greek/olympia/aristarchos/proto"
	aristarchos "github.com/odysseia-greek/olympia/aristarchos/scholar"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
	"time"
)

type DionysosHandler struct {
	Elastic          aristoteles.Client
	Cache            archytas.Client
	Index            string
	Client           service.OdysseiaClient
	DeclensionConfig models.DeclensionConfig
	Tracer           *aristophanes.ClientTracer
	Aggregator       *aristarchos.ClientAggregator
}

// PingPong pongs the ping
func (d *DionysosHandler) pingPong(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /ping status ping
	//
	// Checks if api is reachable
	//
	//	Consumes:
	//	- application/json
	//
	//	Produces:
	//	- application/json
	//
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: ResultModel
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

// returns the health of the api
func (d *DionysosHandler) health(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /health status health
	//
	// Checks if api is healthy
	//
	//	Consumes:
	//	- application/json
	//
	//	Produces:
	//	- application/json
	//
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: Health
	//	  502: Health
	requestId := req.Header.Get(plato.HeaderKey)
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

	if traceCall {
		traceReceived := &pb.TraceRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			Method:       req.Method,
			Url:          req.URL.RequestURI(),
			Host:         req.Host,
		}

		go d.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)

	elasticHealth := d.Elastic.Health().Info()
	dbHealth := models.DatabaseHealth{
		Healthy:       elasticHealth.Healthy,
		ClusterName:   elasticHealth.ClusterName,
		ServerName:    elasticHealth.ServerName,
		ServerVersion: elasticHealth.ServerVersion,
	}
	healthy := models.Health{
		Healthy:  dbHealth.Healthy,
		Time:     time.Now().String(),
		Database: dbHealth,
	}
	if !healthy.Healthy {
		middleware.ResponseWithCustomCode(w, http.StatusBadGateway, healthy)
		return
	}
	middleware.ResponseWithJson(w, healthy)
}

func (d *DionysosHandler) checkGrammar(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /check grammar check
	//
	// Tries to determine what declensions a word might be or what form it takes when it is a verb
	//
	//	Consumes:
	//	- application/json
	//
	//	Produces:
	//	- application/json
	//
	//   Parameters:
	//     + name: word
	//       in: query
	//       description: word or part of word being queried
	//		 example: test
	//       required: true
	//       type: string
	//       format: word
	//		 title: word
	//
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: DeclensionTranslationResults
	//    400: ValidationError
	//	  404: NotFoundError
	//	  405: MethodError
	//    502: ElasticSearchError
	requestId := req.Header.Get(plato.HeaderKey)
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

	if traceCall {
		traceReceived := &pb.TraceRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			Method:       req.Method,
			Url:          req.URL.RequestURI(),
			Host:         req.Host,
		}

		go d.Tracer.Trace(context.Background(), traceReceived)
		logging.Trace(fmt.Sprintf("found traceId: %s", traceID))
	}

	w.Header().Set(plato.HeaderKey, requestId)

	queryWord := req.URL.Query().Get("word")

	if queryWord == "" {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "word",
					Message: "cannot be empty",
				},
			},
		}
		if traceCall {
			parsedResult, _ := json.Marshal(e)
			span := &pb.SpanRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				Action:       "Validating",
				ResponseBody: string(parsedResult),
			}
			go d.Tracer.Span(context.Background(), span)
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	cacheItem, _ := d.Cache.Read(queryWord)
	if cacheItem != nil {
		cache, err := models.UnmarshalDeclensionTranslationResults(cacheItem)
		if err != nil {
			e := models.ValidationError{
				ErrorModel: models.ErrorModel{UniqueCode: traceID},
				Messages: []models.ValidationMessages{
					{
						Field:   "cache",
						Message: err.Error(),
					},
				},
			}

			for _, declension := range cache.Results {
				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()
				md := metadata.New(map[string]string{service.HeaderKey: requestId})
				ctx = metadata.NewOutgoingContext(context.Background(), md)

				speech := pba.PartOfSpeech_VERB

				if strings.Contains(declension.Rule, "noun") {
					speech = pba.PartOfSpeech_NOUN
				}

				if strings.Contains(declension.Rule, "part") {
					speech = pba.PartOfSpeech_PARTICIPLE
				}

				request := pba.AggregatorCreationRequest{
					Word:         declension.Word,
					Rule:         declension.Rule,
					RootWord:     declension.RootWord,
					Translation:  declension.Translation,
					PartOfSpeech: speech,
				}

				test, _ := d.Aggregator.RetrieveEntry(ctx, &pba.AggregatorRequest{
					RootWord: declension.RootWord,
				})

				l, _ := json.Marshal(test)
				logging.Debug(string(l))

				entry, err := d.Aggregator.CreateNewEntry(ctx, &request)
				if err != nil {
					logging.Error(fmt.Sprintf("failed to created entry in aggregator: %s", err.Error()))
					continue
				}

				logging.Debug(fmt.Sprintf("new entry in aggregator created: %v updated: %v", entry.Created, entry.Created))
			}

			if traceCall {
				parsedResult, _ := json.Marshal(e)
				span := &pb.SpanRequest{
					TraceId:      traceID,
					ParentSpanId: spanID,
					Action:       "Cache",
					ResponseBody: string(parsedResult),
				}
				go d.Tracer.Span(context.Background(), span)
			}
			middleware.ResponseWithJson(w, e)
			return
		}
		middleware.ResponseWithJson(w, cache)
		return
	}

	declensions, _ := d.StartFindingRules(queryWord, requestId)
	if len(declensions.Results) == 0 || declensions.Results == nil {
		e := models.NotFoundError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Message: models.NotFoundMessage{
				Type:   queryWord,
				Reason: "no options found",
			},
		}
		if traceCall {
			parsedResult, _ := json.Marshal(e)
			span := &pb.SpanRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				Action:       "Rules Found",
				ResponseBody: string(parsedResult),
			}
			go d.Tracer.Span(context.Background(), span)
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	for _, declension := range declensions.Results {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		md := metadata.New(map[string]string{service.HeaderKey: requestId})
		ctx = metadata.NewOutgoingContext(context.Background(), md)

		speech := pba.PartOfSpeech_VERB

		if strings.Contains(declension.Rule, "noun") {
			speech = pba.PartOfSpeech_NOUN
		}

		if strings.Contains(declension.Rule, "part") {
			speech = pba.PartOfSpeech_PARTICIPLE
		}

		request := pba.AggregatorCreationRequest{
			Word:         declension.Word,
			Rule:         declension.Rule,
			RootWord:     declension.RootWord,
			Translation:  declension.Translation,
			PartOfSpeech: speech,
		}

		test, _ := d.Aggregator.RetrieveEntry(ctx, &pba.AggregatorRequest{
			RootWord: declension.RootWord,
		})

		l, _ := json.Marshal(test)
		logging.Debug(string(l))

		entry, err := d.Aggregator.CreateNewEntry(ctx, &request)
		if err != nil {
			logging.Error(fmt.Sprintf("failed to created entry in aggregator: %s", err.Error()))
			continue
		}

		logging.Debug(fmt.Sprintf("new entry in aggregator created: %v updated: %v", entry.Created, entry.Created))
	}

	stringifiedDeclension, _ := declensions.Marshal()
	ttl := time.Hour
	err := d.Cache.SetWithTTL(queryWord, string(stringifiedDeclension), ttl)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "cache",
					Message: err.Error(),
				},
			},
		}
		if traceCall {
			parsedResult, _ := json.Marshal(e)
			span := &pb.SpanRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				Action:       "Cache",
				ResponseBody: string(parsedResult),
			}
			go d.Tracer.Span(context.Background(), span)
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	middleware.ResponseWithJson(w, *declensions)
}
