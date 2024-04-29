package dictionary

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	elasticmodels "github.com/odysseia-greek/agora/aristoteles/models"
	plato "github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/agora/plato/transform"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AlexandrosHandler struct {
	Elastic aristoteles.Client
	Index   string
	Client  service.OdysseiaClient
	Tracer  *aristophanes.ClientTracer
}

const (
	partial  string = "partial"
	extended string = "extended"
	exact    string = "exact"
	fuzzy    string = "fuzzy"

	defaultLang string = "greek"
)

var allowedLanguages = []string{"greek", "english", "dutch"}
var allowedMatchModes = []string{partial, exact, extended, fuzzy}

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
func (a *AlexandrosHandler) searchWord(w http.ResponseWriter, req *http.Request) {
	requestId := req.Header.Get(plato.HeaderKey)
	splitID := strings.Split(requestId, "+")

	traceCall := false
	var traceID, parentSpanID, spanID string

	if len(splitID) >= 3 {
		traceCall = splitID[2] == "1"
	}

	if len(splitID) >= 1 {
		traceID = splitID[0]
	}
	if len(splitID) >= 2 {
		parentSpanID = splitID[1]
	}

	w.Header().Set(plato.HeaderKey, requestId)

	if traceCall {
		traceReceived := &pb.TraceRequest{
			TraceId:      traceID,
			ParentSpanId: parentSpanID,
			Method:       req.Method,
			Url:          req.URL.RequestURI(),
			Host:         req.Host,
		}

		go a.Tracer.Trace(context.Background(), traceReceived)

		spanStart := &pb.StartSpanRequest{
			TraceId:      traceID,
			ParentSpanId: parentSpanID,
			Action:       "searchWord",
			RequestBody:  req.URL.RequestURI(),
		}

		resp, _ := a.Tracer.StartSpan(context.Background(), spanStart)
		requestId = resp.CombinedId
		split := strings.Split(resp.CombinedId, "+")
		if len(split) >= 2 {
			spanID = split[1]
		}
	}

	queryWord := req.URL.Query().Get("word")
	mode := req.URL.Query().Get("mode")
	language := req.URL.Query().Get("lang")
	text := req.URL.Query().Get("searchInText")
	searchInText, _ := strconv.ParseBool(text)

	if a.validateAndRespondError(w, "word", queryWord, nil, traceID, parentSpanID, spanID, traceCall) {
		return
	}

	if language == "" {
		language = defaultLang
	}

	if a.validateAndRespondError(w, "lang", language, allowedLanguages, traceID, parentSpanID, spanID, traceCall) {
		return
	}

	if mode == "" {
		mode = fuzzy
	}

	if a.validateAndRespondError(w, "mode", mode, allowedMatchModes, traceID, parentSpanID, spanID, traceCall) {
		return
	}

	var wordToQuery string
	if mode == fuzzy {
		wordToQuery = extractBaseWord(queryWord)
	} else {
		wordToQuery = strings.ToLower(queryWord)
	}

	query := a.processQuery(mode, language, wordToQuery)

	response, err := a.Elastic.Query().Match(a.Index, query)
	if err != nil {
		e := models.ElasticSearchError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Message: models.ElasticErrorMessage{
				ElasticError: err.Error(),
			},
		}
		if traceCall {
			span := &pb.CloseSpanRequest{
				TraceId:      traceID,
				ParentSpanId: parentSpanID,
				SpanId:       spanID,
				ResponseCode: http.StatusOK,
			}
			go a.Tracer.CloseSpan(context.Background(), span)
		}

		middleware.ResponseWithJson(w, e)
		return
	}

	var meroi []models.Meros

	if len(response.Hits.Hits) == 0 {
		noResponseQuery := a.processQuery(mode, language, queryWord)

		noResponseResponse, err := a.Elastic.Query().Match(a.Index, noResponseQuery)
		if err != nil {
			e := models.ElasticSearchError{
				ErrorModel: models.ErrorModel{UniqueCode: traceID},
				Message: models.ElasticErrorMessage{
					ElasticError: err.Error(),
				},
			}
			if traceCall {
				span := &pb.CloseSpanRequest{
					TraceId:      traceID,
					ParentSpanId: parentSpanID,
					SpanId:       spanID,
					ResponseCode: http.StatusOK,
				}
				go a.Tracer.CloseSpan(context.Background(), span)
			}

			middleware.ResponseWithJson(w, e)
			return
		}

		for _, hit := range noResponseResponse.Hits.Hits {
			jsonHit, _ := json.Marshal(hit.Source)
			meros, _ := models.UnmarshalMeros(jsonHit)
			if meros.Original != "" {
				meros.Greek = meros.Original
				meros.Original = ""
			}
			meroi = append(meroi, meros)
		}

	} else {
		for _, hit := range response.Hits.Hits {
			jsonHit, _ := json.Marshal(hit.Source)
			meros, _ := models.UnmarshalMeros(jsonHit)
			if meros.Original != "" {
				meros.Greek = meros.Original
				meros.Original = ""
			}
			meroi = append(meroi, meros)
		}
	}

	if traceCall {
		go a.databaseSpan(response, query, traceID, spanID)
	}

	results := a.cleanSearchResult(meroi)

	if searchInText {
		for i, hit := range results.Hits {
			baseWord := a.extractBaseWord(hit.Hit.Greek)
			herodotosRequestId := requestId
			var herodotosSpanID string
			if traceCall {
				herodotosSpan := &pb.StartSpanRequest{
					TraceId:      traceID,
					ParentSpanId: spanID,
					Action:       "analyseText",
				}

				resp, _ := a.Tracer.StartSpan(context.Background(), herodotosSpan)
				herodotosRequestId = resp.CombinedId
				split := strings.Split(resp.CombinedId, "+")
				if len(split) >= 2 {
					herodotosSpanID = split[1]
				}
			}

			foundInText, err := a.Client.Herodotos().AnalyseText(baseWord, herodotosRequestId)

			if err != nil {
				logging.Error(fmt.Sprintf("could not query any texts for word: %s error: %s", hit.Hit.Greek, err.Error()))
			} else {
				var source models.Rhema
				defer foundInText.Body.Close()
				err = json.NewDecoder(foundInText.Body).Decode(&source)
				if err != nil {
					logging.Error(fmt.Sprintf("error while decoding: %s", err.Error()))
				}

				if traceCall {
					span := &pb.CloseSpanRequest{
						TraceId:      traceID,
						ParentSpanId: parentSpanID,
						SpanId:       herodotosSpanID,
						ResponseCode: int32(foundInText.StatusCode),
					}
					go a.Tracer.CloseSpan(context.Background(), span)
				}

				results.Hits[i].FoundInText = &source
			}

		}
	}

	if traceCall {
		span := &pb.CloseSpanRequest{
			TraceId:      traceID,
			ParentSpanId: parentSpanID,
			SpanId:       spanID,
			ResponseCode: http.StatusOK,
		}

		go a.Tracer.CloseSpan(context.Background(), span)
	}

	middleware.ResponseWithCustomCode(w, http.StatusOK, results)
}

func (a *AlexandrosHandler) databaseSpan(response *elasticmodels.Response, query map[string]interface{}, traceID, spanID string) {
	parsedQuery, _ := json.Marshal(query)
	dataBaseSpan := &pb.DatabaseSpanRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		Action:       "search",
		Query:        string(parsedQuery),
		Hits:         response.Hits.Total.Value,
		TimeTook:     response.Took,
	}

	a.Tracer.DatabaseSpan(context.Background(), dataBaseSpan)
}

func (a *AlexandrosHandler) extractBaseWord(queryWord string) string {
	// Normalize and split the input
	strippedWord := transform.RemoveAccents(strings.ToLower(queryWord))
	splitWord := strings.Split(strippedWord, " ")

	greekPronouns := map[string]bool{"η": true, "ο": true, "το": true}
	cleanWord := func(word string) string {
		return strings.Trim(word, ",.!?-") // Add any other punctuation as needed
	}

	for _, word := range splitWord {
		cleanedWord := cleanWord(word)

		if strings.HasPrefix(cleanedWord, "-") {
			continue
		}

		if _, isPronoun := greekPronouns[cleanedWord]; !isPronoun {
			// If the word is not a pronoun, it's likely the correct word
			return cleanedWord
		}
	}

	return queryWord
}

// CleanSearchResult removes duplicate entries and cleans up the search results.
func (a *AlexandrosHandler) cleanSearchResult(results []models.Meros) *models.ExtendedResponse {
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

	var extendedResponse models.ExtendedResponse

	for _, filtered := range filteredSlice {
		hit := models.Hit{
			Hit:         filtered,
			FoundInText: nil,
		}

		extendedResponse.Hits = append(extendedResponse.Hits, hit)
	}
	return &extendedResponse
}

// ProcessQuery constructs the appropriate Elasticsearch query based on the provided mode and language.
func (a *AlexandrosHandler) processQuery(option string, language string, queryWord string) map[string]interface{} {
	switch option {
	case partial:
		// Process partial matching query
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
			"size": 20,
		}

	case extended:
		// Process phrase matching query
		// Example implementation:
		return map[string]interface{}{
			"query": map[string]interface{}{
				"match_phrase": map[string]string{
					language: queryWord,
				},
			},
			"size": 20,
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
								fmt.Sprintf("%s.keyword", language): fmt.Sprintf("%s,", queryWord),
							},
						},
						map[string]interface{}{
							"prefix": map[string]interface{}{
								fmt.Sprintf("%s.keyword", language): fmt.Sprintf("%s ", queryWord),
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
			"size": 10,
		}

	case fuzzy:
		// Process fuzzy matching query (with 2 levenshtein)
		// Example implementation:
		return map[string]interface{}{
			"query": map[string]interface{}{
				"fuzzy": map[string]interface{}{
					language: map[string]interface{}{
						"value":     queryWord,
						"fuzziness": 2,
					},
				},
			},
			"size": 5,
		}

	default:
		// Handle unknown option
		return nil
	}
}

func (a *AlexandrosHandler) validateAndRespondError(w http.ResponseWriter, field, value string, allowedValues []string, traceID, parentSpanID, spanID string, traceCall bool) bool {
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
			span := &pb.CloseSpanRequest{
				TraceId:      traceID,
				ParentSpanId: parentSpanID,
				SpanId:       spanID,
				ResponseCode: http.StatusOK,
			}
			go a.Tracer.CloseSpan(context.Background(), span)
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

func extractBaseWord(queryWord string) string {
	// Normalize and split the input
	strippedWord := transform.RemoveAccents(strings.ToLower(queryWord))
	splitWord := strings.Split(strippedWord, " ")

	// Known Greek pronouns
	greekPronouns := map[string]bool{"η": true, "ο": true, "το": true}

	// Function to clean punctuation from a word
	cleanWord := func(word string) string {
		return strings.Trim(word, ",.!?-") // Add any other punctuation as needed
	}

	// Iterate through the words
	for _, word := range splitWord {
		cleanedWord := cleanWord(word)

		if strings.HasPrefix(cleanedWord, "-") {
			// Skip words starting with "-"
			continue
		}

		if _, isPronoun := greekPronouns[cleanedWord]; !isPronoun {
			// If the word is not a pronoun, it's likely the correct word
			return cleanedWord
		}
	}

	return queryWord
}
