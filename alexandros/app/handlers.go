package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	elasticmodels "github.com/odysseia-greek/agora/aristoteles/models"
	plato "github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/models"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/app"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"net/http"
	"strings"
	"time"
)

type AlexandrosHandler struct {
	Elastic aristoteles.Client
	Index   string
	Tracer  *aristophanes.ClientTracer
}

const (
	fuzzy       string = "fuzzy"
	phrase      string = "phrase"
	exact       string = "exact"
	defaultLang string = "greek"
)

var allowedLanguages = []string{"greek", "english", "dutch"}
var allowedMatchModes = []string{fuzzy, exact, phrase}

// pingPong is used to check the reachability of the API. It returns a "pong" response when called, indicating that the API is reachable.
func (a *AlexandrosHandler) pingPong(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /ping status ping
	//
	// Checks the reachability of the API.
	//
	// This endpoint returns a "pong" response to indicate that the API is reachable.
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

// health returns the health of the API
// The `health` function is used to check the health status of the API. It checks if the underlying infrastructure, such as the database and Elasticsearch, is healthy.
func (a *AlexandrosHandler) health(w http.ResponseWriter, req *http.Request) {
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

		go a.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)

	elasticHealth := a.Elastic.Health().Info()
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

// searchWord searches the dictionary for a word in Greek, English or Dutch
// It handles the HTTP GET request to the "/search" endpoint.
//
// The function performs the following steps:
//  1. Validates the query parameters and handles any validation errors.
//  2. Constructs the appropriate Elasticsearch query based on the provided mode and language.
//  3. Executes the Elasticsearch query and handles any errors.
//  4. Processes the query results and prepares the response.
//  5. Sends the response back to the client.
//
// Responses:
//   - 200: Returns an array of models.Meros representing the search results.
//   - 400: Returns a models.ValidationError if there are validation errors in the query parameters.
//   - 404: Returns a models.NotFoundError if no results are found for the given query.
//   - 405: Returns a models.MethodError if the HTTP method is not allowed for the endpoint.
//   - 502: Returns a models.ElasticSearchError if there is an error in the Elasticsearch query.
func (a *AlexandrosHandler) searchWord(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /search search search
	//
	// Searches the dictionary for a word in Greek (English wip)
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
	//     + name: mode
	//       in: query
	//       description: Determines a number of query modes; fuzzy, exact or phrase
	//		 example: false
	//       required: false
	//       type: string
	//       format: mode
	//		 title: mode
	//     + name: lang
	//       in: query
	//       description: language to use (greek, english, dutch)
	//		 example: greek
	//       required: false
	//       type: string
	//       format: lang
	//		 title: lang
	//
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: []Meros
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

		go a.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)

	queryWord := req.URL.Query().Get("word")
	mode := req.URL.Query().Get("mode")
	language := req.URL.Query().Get("lang")

	if a.validateAndRespondError(w, "word", queryWord, nil, traceID, spanID, traceCall) {
		return
	}

	if language == "" {
		language = defaultLang
	}

	if a.validateAndRespondError(w, "lang", language, allowedLanguages, traceID, spanID, traceCall) {
		return
	}

	if mode == "" {
		mode = fuzzy
	}

	if a.validateAndRespondError(w, "mode", mode, allowedMatchModes, traceID, spanID, traceCall) {
		return
	}

	query := a.processQuery(mode, language, queryWord)

	response, err := a.Elastic.Query().Match(a.Index, query)
	if err != nil {
		e := models.ElasticSearchError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Message: models.ElasticErrorMessage{
				ElasticError: err.Error(),
			},
		}
		if traceCall {
			parsedError, _ := json.Marshal(e)
			span := &pb.SpanRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				Action:       "MatchElasticQuery",
				ResponseBody: string(parsedError),
			}
			go a.Tracer.Span(context.Background(), span)
		}

		middleware.ResponseWithJson(w, e)
		return
	}

	var meroi []models.Meros

	for _, hit := range response.Hits.Hits {
		jsonHit, _ := json.Marshal(hit.Source)
		meros, _ := models.UnmarshalMeros(jsonHit)
		if meros.Original != "" {
			meros.Greek = meros.Original
			meros.Original = ""
		}
		meroi = append(meroi, meros)
	}

	if traceCall {
		go a.databaseSpan(response, query, traceID, spanID)
	}

	results := a.cleanSearchResult(meroi)
	middleware.ResponseWithJson(w, results)
}

func (a *AlexandrosHandler) databaseSpan(response *elasticmodels.Response, query map[string]interface{}, traceID, spanID string) {
	parsedResult, _ := json.Marshal(response)
	parsedQuery, _ := json.Marshal(query)
	dataBaseSpan := &pb.DatabaseSpanRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		Action:       "search",
		Query:        string(parsedQuery),
		ResultJson:   string(parsedResult),
	}

	a.Tracer.DatabaseSpan(context.Background(), dataBaseSpan)
}

// CleanSearchResult removes duplicate entries and cleans up the search results.
func (a *AlexandrosHandler) cleanSearchResult(results []models.Meros) []models.Meros {
	filteredSlice := make([]models.Meros, 0)

	// Create a map to track entries by Greek and English
	entryMap := make(map[string]models.Meros)

	for _, meros := range results {
		// Check if the entry already exists in the map
		if existingEntry, exists := entryMap[meros.Greek+meros.English]; exists {
			// If the existing entry has a non-empty Dutch field and the current entry has an empty Dutch field, skip the current entry
			if existingEntry.Dutch != "" && meros.Dutch == "" {
				continue
			}
		}

		entryMap[meros.Greek+meros.English] = meros
	}

	// Collect the filtered entries
	for _, entry := range entryMap {
		filteredSlice = append(filteredSlice, entry)
	}

	return filteredSlice
}

// ProcessQuery constructs the appropriate Elasticsearch query based on the provided mode and language.
func (a *AlexandrosHandler) processQuery(option string, language string, queryWord string) map[string]interface{} {
	switch option {
	case fuzzy:
		// Process fuzzy matching query
		// Example implementation:
		return map[string]interface{}{
			"query": map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":    queryWord,
					"type":     "most_fields",
					"analyzer": "greek_analyzer",
					"fields": []string{
						language,
						"original",
					},
				},
			},
			"size": 50,
		}

	case phrase:
		// Process phrase matching query
		// Example implementation:
		return map[string]interface{}{
			"query": map[string]interface{}{
				"match_phrase": map[string]string{
					language: queryWord,
				},
			},
		}

	case exact:
		// Process exact matching query
		// Example implementation:
		return map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []interface{}{
						map[string]interface{}{
							"prefix": map[string]interface{}{
								fmt.Sprintf("%s.keyword", language): queryWord,
							},
						},
						map[string]interface{}{
							"term": map[string]interface{}{
								fmt.Sprintf("%s.keyword", language): queryWord,
							},
						},
					},
				},
			},
		}

	default:
		// Handle unknown option
		return nil
	}
}

func (a *AlexandrosHandler) validateAndRespondError(w http.ResponseWriter, field, value string, allowedValues []string, traceID, spanID string, traceCall bool) bool {
	if err := a.validateQueryParam(value, field, allowedValues); err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   field,
					Message: err.Error(),
				},
			},
		}
		if traceCall {
			parsedResult, _ := json.Marshal(e)
			span := &pb.SpanRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				Action:       "validateQueryParam",
				ResponseBody: string(parsedResult),
			}
			go a.Tracer.Span(context.Background(), span)
		}
		middleware.ResponseWithJson(w, e)
		return true
	}
	return false
}

// ValidateQueryParam validates the query parameters based on the provided allowed values.
func (a *AlexandrosHandler) validateQueryParam(queryParam, field string, allowedValues []string) error {
	if queryParam == "" {
		return fmt.Errorf("%s cannot be empty", field)
	}

	if allowedValues == nil {
		return nil
	}

	for _, value := range allowedValues {
		if value == queryParam {
			return nil
		}
	}

	return fmt.Errorf("invalid %s value. Please choose one of the following: %s", field, strings.Join(allowedValues, ", "))
}
