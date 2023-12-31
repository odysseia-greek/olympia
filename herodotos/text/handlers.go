package text

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/agora/aristoteles"
	elasticmodels "github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/helpers"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/models"
	plato "github.com/odysseia-greek/agora/plato/service"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"log"
	"net/http"
	"strings"
	"time"
)

type HerodotosHandler struct {
	Tracer  *aristophanes.ClientTracer
	Elastic aristoteles.Client
	Index   string
}

const (
	Author  string = "author"
	Authors string = "authors"
	Book    string = "book"
	Books   string = "books"
	Id      string = "_id"
)

// PingPong pongs the ping
func (h *HerodotosHandler) pingPong(w http.ResponseWriter, req *http.Request) {
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
func (h *HerodotosHandler) health(w http.ResponseWriter, req *http.Request) {
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

		go h.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)

	elasticHealth := h.Elastic.Health().Info()
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

// creates a new sentence for translation
func (h *HerodotosHandler) createQuestion(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /createQuestion sentence createSentence
	//
	// Creates a new sentence for translation
	//
	//	Consumes:
	//	- application/json
	//
	//	Produces:
	//	- application/json
	//
	//   Parameters:
	//     + name: author
	//       in: query
	//       description: author to be used for creating a sentence
	//		 example: herodotos
	//       required: true
	//       type: string
	//       format: author
	//		 title: author
	//     + name: book
	//       in: query
	//       description: book to be used for creating a sentence
	//		 example: 2
	//       required: true
	//       type: string
	//       format: book
	//		 title: book
	//
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: CreateSentenceResponse
	//    400: ValidationError
	//	  404: NotFoundError
	//	  405: MethodError
	//    502: ElasticSearchError
	author := req.URL.Query().Get(Author)
	book := req.URL.Query().Get(Book)

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

		go h.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)

	if author == "" || book == "" {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "author and book",
					Message: "cannot be empty",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							Author: map[string]interface{}{
								"query":         author,
								"operator":      "and",
								"fuzziness":     "AUTO",
								"prefix_length": 0,
							},
						},
					},
					{
						"match": map[string]interface{}{
							Book: book,
						},
					},
				},
			},
		},
	}

	response, err := h.Elastic.Query().Match(h.Index, query)

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

	if traceCall {
		go h.databaseSpan(response, query, traceID, spanID)
	}

	if len(response.Hits.Hits) == 0 {
		e := models.NotFoundError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Message: models.NotFoundMessage{
				Type:   fmt.Sprintf("author: %s and book: %s", author, book),
				Reason: "no hits for combination",
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	randNumber := helpers.GenerateRandomNumber(len(response.Hits.Hits))
	questionItem := response.Hits.Hits[randNumber]
	id := questionItem.ID

	elasticJson, _ := json.Marshal(questionItem.Source)
	rhemaSource, err := models.UnmarshalRhema(elasticJson)
	if err != nil || rhemaSource.Translations == nil {
		errorMessage := fmt.Errorf("an error occurred while parsing %s", elasticJson)
		log.Print(errorMessage.Error())
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "createQuestion",
					Message: errorMessage.Error(),
				},
				{
					Field:   "translation",
					Message: "cannot be nil",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	question := models.CreateSentenceResponse{Sentence: rhemaSource.Greek,
		SentenceId: id}
	middleware.ResponseWithJson(w, question)
}

// swagger:parameters checkSentence
type checkSentenceParameters struct {
	// in:body
	Application models.CheckSentenceRequest
}

// checks the validity of an answer
func (h *HerodotosHandler) checkSentence(w http.ResponseWriter, req *http.Request) {
	// swagger:route POST /checkSentence sentence checkSentence
	//
	// Checks the sentence against a set of provided accepted answers
	//
	//	Consumes:
	//	- application/json
	//
	//	Produces:
	//	- application/json
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: CheckSentenceResponse
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
		go h.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)

	log.Printf("received %s code with value: %s", plato.HeaderKey, traceID)

	var checkSentenceRequest models.CheckSentenceRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&checkSentenceRequest)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "decoding",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	if traceCall {
		jsonCheckSentence, _ := checkSentenceRequest.Marshal()
		span := &pb.SpanRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			Action:       "RequestBodyFromPost",
			RequestBody:  string(jsonCheckSentence),
		}

		go h.Tracer.Span(context.Background(), span)
	}

	query := h.Elastic.Builder().MatchQuery(Id, checkSentenceRequest.SentenceId)
	elasticResult, err := h.Elastic.Query().Match(h.Index, query)
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

	if len(elasticResult.Hits.Hits) == 0 {
		e := models.NotFoundError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Message: models.NotFoundMessage{
				Type:   "no reults",
				Reason: "elastic results 0 results",
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	if traceCall {
		go h.databaseSpan(elasticResult, query, traceID, spanID)
	}

	elasticJson, _ := json.Marshal(elasticResult.Hits.Hits[0].Source)
	original, err := models.UnmarshalRhema(elasticJson)
	if err != nil || original.Translations == nil {
		errorMessage := fmt.Errorf("an error occurred while parsing %s", elasticJson)
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "createQuestion",
					Message: errorMessage.Error(),
				},
				{
					Field:   "translation",
					Message: "cannot be nil",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	var sentence string
	var percentage float64
	for _, solution := range original.Translations {
		levenshteinDist := levenshteinDistance(solution, checkSentenceRequest.ProvidedSentence)
		lenOfLongestSentence := longestStringOfTwo(solution, checkSentenceRequest.ProvidedSentence)
		levenshteinPerc := levenshteinDistanceInPercentage(levenshteinDist, lenOfLongestSentence)
		if levenshteinPerc > percentage {
			sentence = solution
			percentage = levenshteinPerc
		}
	}

	roundedPercentage := fmt.Sprintf("%.2f", percentage)

	model := findMatchingWordsWithSpellingAllowance(sentence, checkSentenceRequest.ProvidedSentence)

	response := models.CheckSentenceResponse{
		LevenshteinPercentage: roundedPercentage,
		QuizSentence:          sentence,
		AnswerSentence:        checkSentenceRequest.ProvidedSentence,
		MatchingWords:         model.MatchingWords,
		NonMatchingWords:      model.NonMatchingWords,
		SplitQuizSentence:     model.SplitQuizSentence,
		SplitAnswerSentence:   model.SplitAnswerSentence,
	}

	middleware.ResponseWithJson(w, response)
}

func (h *HerodotosHandler) queryAuthors(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /authors authors authors
	//
	// Finds all authors
	//
	//	Consumes:
	//	- application/json
	//
	//	Produces:
	//	- application/json
	//
	//   Parameters:
	//     + name: author
	//       in: path
	//       description: author to be used for creating a sentence
	//		 example: herodotos
	//       required: true
	//       type: string
	//       format: author
	//		 title: author
	//
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: Authors
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

		go h.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)

	log.Printf("received %s code with value: %s", plato.HeaderKey, traceID)

	query := h.Elastic.Builder().Aggregate(Authors, Author)

	elasticResult, err := h.Elastic.Query().MatchAggregate(h.Index, query)

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

	var authors models.Authors
	for _, bucket := range elasticResult.Aggregations.AuthorAggregation.Buckets {
		author := models.Author{Author: strings.ToLower(fmt.Sprintf("%v", bucket.Key))}
		authors.Authors = append(authors.Authors, author)
	}

	middleware.ResponseWithJson(w, authors)
}

func (h *HerodotosHandler) queryBooks(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /authors/{author}/books authors books
	//
	// Finds all books
	//
	//	Consumes:
	//	- application/json
	//
	//	Produces:
	//	- application/json
	//
	//   Parameters:
	//     + name: author
	//       in: path
	//       description: author to be used for creating a sentence
	//		 example: herodotos
	//       required: true
	//       type: string
	//       format: author
	//		 title: author
	//
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: Books
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

		go h.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)

	log.Printf("received %s code with value: %s", plato.HeaderKey, traceID)

	pathParams := mux.Vars(req)
	author := pathParams[Author]

	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			Books: map[string]interface{}{
				"terms": map[string]interface{}{
					"field": Book,
					"size":  500,
				},
			},
		},
		"query": map[string]interface{}{
			"wildcard": map[string]interface{}{
				Author: map[string]interface{}{
					"value":            fmt.Sprintf("*%s*", author),
					"case_insensitive": true,
				},
			},
		},
	}

	elasticResult, err := h.Elastic.Query().MatchAggregate(h.Index, query)

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

	var books models.Books
	for _, bucket := range elasticResult.Aggregations.BookAggregation.Buckets {

		book := models.Book{Book: int64(bucket.Key.(float64))}
		books.Books = append(books.Books, book)

	}

	middleware.ResponseWithJson(w, books)
}

func (h *HerodotosHandler) databaseSpan(response *elasticmodels.Response, query map[string]interface{}, traceID, spanID string) {
	parsedResult, _ := json.Marshal(response)
	parsedQuery, _ := json.Marshal(query)
	dataBaseSpan := &pb.DatabaseSpanRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		Action:       "search",
		Query:        string(parsedQuery),
		ResultJson:   string(parsedResult),
	}

	h.Tracer.DatabaseSpan(context.Background(), dataBaseSpan)
}

// calculates the amount of changes needed to have two sentences match
// example: Distance from Python to Peithen is 3
func levenshteinDistance(question, answer string) int {
	questionLen := len(question)
	answerLen := len(answer)
	column := make([]int, len(question)+1)

	for y := 1; y <= questionLen; y++ {
		column[y] = y
	}
	for x := 1; x <= answerLen; x++ {
		column[0] = x
		lastKey := x - 1
		for y := 1; y <= questionLen; y++ {
			oldKey := column[y]
			var incr int
			if question[y-1] != answer[x-1] {
				incr = 1
			}

			column[y] = minimum(column[y]+1, column[y-1]+1, lastKey+incr)
			lastKey = oldKey
		}
	}
	return column[questionLen]
}

// creates a percentage based on the levenshtein and the longest string
func levenshteinDistanceInPercentage(levenshteinDistance int, longestString int) float64 {
	return (1.00 - float64(levenshteinDistance)/float64(longestString)) * 100.00
}

func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}

func longestStringOfTwo(a, b string) int {
	if len(a) >= len(b) {
		return len(a)
	}
	return len(b)
}

// take two sentences and creates a list of matching words and words with a typo (1 levenshtein)
func findMatchingWordsWithSpellingAllowance(source, target string) (response models.CheckSentenceResponse) {
	s := removeCharacters(source, ",`~<>/?!.;:'\"")
	t := removeCharacters(target, ",`~<>/?!.;:'\"")

	sourceSentence := strings.Split(s, " ")
	targetSentence := strings.Split(t, " ")

	response.SplitQuizSentence = sourceSentence
	response.SplitAnswerSentence = targetSentence

	//todo probably we want to check the len of the word and determine what levenshtein distance is allowed
	//for example: the vs | tha | then | thee | dhe which one is misspelled and which should we ignore?
	for i, wordInSource := range sourceSentence {
		for _, wordInTarget := range targetSentence {
			sourceWord := strings.ToLower(wordInSource)
			targetWord := strings.ToLower(wordInTarget)

			levenshtein := levenshteinDistance(sourceWord, targetWord)

			// might be changed to one levenshtein as being typo's
			if levenshtein == 0 {
				response.MatchingWords = append(response.MatchingWords, models.MatchingWord{
					Word:        wordInSource,
					SourceIndex: i,
				})
				break
			}
		}
	}

	var slice []string
	for _, word := range response.MatchingWords {
		slice = append(slice, word.Word)
	}

	for i, wordInSource := range sourceSentence {
		wordInSlice := checkSliceForItem(slice, wordInSource)
		if !wordInSlice {
			levenshteinModel := models.NonMatchingWord{
				Word:        wordInSource,
				SourceIndex: i,
				Matches:     nil,
			}
			for j, wordInTarget := range targetSentence {
				sourceWord := strings.ToLower(wordInSource)
				targetWord := strings.ToLower(wordInTarget)

				levenshtein := levenshteinDistance(sourceWord, targetWord)
				percentage := levenshteinDistanceInPercentage(levenshtein, longestStringOfTwo(sourceWord, targetWord))
				roundedPercentage := fmt.Sprintf("%.2f", percentage)

				matchModel := models.Match{
					Match:       wordInTarget,
					Levenshtein: levenshtein,
					AnswerIndex: j,
					Percentage:  roundedPercentage,
				}
				levenshteinModel.Matches = append(levenshteinModel.Matches, matchModel)
			}
			response.NonMatchingWords = append(response.NonMatchingWords, levenshteinModel)
		}
	}

	return
}

// takes a slice and returns a bool if value is part of the slice
func checkSliceForItem(slice []string, sourceWord string) bool {
	for _, item := range slice {
		if item == sourceWord {
			return true
		}
	}

	return false
}

func removeCharacters(input string, characters string) string {
	filter := func(r rune) rune {
		if !strings.ContainsRune(characters, r) {
			return r
		}
		return -1
	}

	return strings.Map(filter, input)

}

func streamlineSentenceBeforeCompare(matchingWords []string, sentence string) string {
	newSentence := ""

	sourceSentence := strings.Split(sentence, " ")

	for index, wordInSource := range sourceSentence {
		wordInSlice := checkSliceForItem(matchingWords, wordInSource)
		if !wordInSlice {
			newSentence += wordInSource
			if len(sourceSentence) != index+1 {
				newSentence += " "
			}
		}
	}

	return newSentence
}
