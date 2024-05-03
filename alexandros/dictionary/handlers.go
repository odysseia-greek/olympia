package dictionary

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	elasticmodels "github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/agora/plato/transform"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AlexandrosHandler struct {
	Elastic  aristoteles.Client
	Index    string
	Client   service.OdysseiaClient
	Streamer pb.TraceService_ChorusClient
	Cancel   context.CancelFunc
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
	//requestId := req.Context().Value(plato.HeaderKey).(string)
	var requestId string
	fromContext := req.Context().Value(config.DefaultTracingName)
	if fromContext == nil {
		requestId = req.Header.Get(config.HeaderKey)
	} else {
		requestId = fromContext.(string)
	}
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

	queryWord := req.URL.Query().Get("word")
	mode := req.URL.Query().Get("mode")
	language := req.URL.Query().Get("lang")
	text := req.URL.Query().Get("searchInText")
	searchInText, _ := strconv.ParseBool(text)

	if a.validateAndRespondError(w, "word", queryWord, nil, traceID) {
		return
	}

	if language == "" {
		language = defaultLang
	}

	if a.validateAndRespondError(w, "lang", language, allowedLanguages, traceID) {
		return
	}

	if mode == "" {
		mode = fuzzy
	}

	if a.validateAndRespondError(w, "mode", mode, allowedMatchModes, traceID) {
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
			if traceCall {
				herodotosSpan := &pb.ParabasisRequest{
					TraceId:      traceID,
					ParentSpanId: spanID,
					SpanId:       spanID,
					RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
						Action: "analyseText",
						Status: fmt.Sprintf("querying Herodotos for word: %s", baseWord),
					}},
				}

				err := a.Streamer.Send(herodotosSpan)
				if err != nil {
					logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
				}
			}

			startTime := time.Now()
			foundInText, err := a.Client.Herodotos().AnalyseText(baseWord, requestId)
			endTime := time.Since(startTime)

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
					herodotosSpan := &pb.ParabasisRequest{
						TraceId:      traceID,
						ParentSpanId: spanID,
						SpanId:       spanID,
						RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
							Action: "analyseText",
							Took:   fmt.Sprintf("%v", endTime),
							Status: fmt.Sprintf("querying Herodotos returned: %d", foundInText.StatusCode),
						}},
					}

					err := a.Streamer.Send(herodotosSpan)
					if err != nil {
						logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
					}
				}

				results.Hits[i].FoundInText = &source
			}

		}
	}

	middleware.ResponseWithCustomCode(w, http.StatusOK, results)
}

func (a *AlexandrosHandler) databaseSpan(response *elasticmodels.Response, query map[string]interface{}, traceID, spanID string) {
	parsedQuery, _ := json.Marshal(query)
	dataBaseSpan := &pb.ParabasisRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		SpanId:       spanID,
		RequestType: &pb.ParabasisRequest_DatabaseSpan{DatabaseSpan: &pb.DatabaseSpanRequest{
			Action:   "search",
			Query:    string(parsedQuery),
			Hits:     response.Hits.Total.Value,
			TimeTook: response.Took,
		}},
	}

	err := a.Streamer.Send(dataBaseSpan)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}
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

func (a *AlexandrosHandler) validateAndRespondError(w http.ResponseWriter, field, value string, allowedValues []string, traceID string) bool {
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

func (a *AlexandrosHandler) Close() {
	a.Cancel()
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
