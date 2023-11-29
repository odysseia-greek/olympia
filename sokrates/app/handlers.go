package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	elasticmodels "github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/randomizer"
	plato "github.com/odysseia-greek/agora/plato/service"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/app"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type SokratesHandler struct {
	Tracer     *aristophanes.ClientTracer
	Elastic    elastic.Client
	Randomizer randomizer.Random
	SearchWord string
	Index      string
}

const (
	Method     string = "method"
	Authors    string = "authors"
	Category   string = "category"
	Categories string = "categories"
	Chapter    string = "chapter"
)

// PingPong pongs the ping
func (s *SokratesHandler) pingPong(w http.ResponseWriter, req *http.Request) {
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
func (s *SokratesHandler) health(w http.ResponseWriter, req *http.Request) {
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

		go s.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)

	elasticHealth := s.Elastic.Health().Info()
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

func (s *SokratesHandler) findHighestChapter(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /methods/{method}/categories/{category}/chapters methods highestChapter
	//
	// Finds the highest chapter for a combination of method and categories
	//
	//	Consumes:
	//	- application/json
	//
	//	Produces:
	//	- application/json
	//
	//   Parameters:
	//     + name: method
	//       in: query
	//       description: method to be used
	//		 example: aristophanes
	//       required: true
	//       type: string
	//       format: method
	//		 title: method
	//     + name: category
	//       in: query
	//       description: category to be used
	//		 example: frogs
	//       required: true
	//       type: string
	//       format: category
	//		 title: category
	//
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: LastChapterResponse
	//    400: ValidationError
	//	  405: MethodError
	//	  502: ElasticSearchError
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

		go s.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)
	log.Printf("received %s code with value: %s", plato.HeaderKey, traceID)

	pathParams := mux.Vars(req)
	category := pathParams[Category]
	method := pathParams[Method]

	if len(category) < 2 {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   Category,
					Message: "must be longer than 1",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	mustQuery := []map[string]string{
		{
			Method: method,
		},
		{
			Category: category,
		},
	}

	query := s.Elastic.Builder().MultipleMatch(mustQuery)
	mode := "desc"

	elasticResult, err := s.Elastic.Query().MatchWithSort(s.Index, mode, Chapter, 1, query)
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
				Type:   "",
				Reason: fmt.Sprintf("no chapters found for category: %s and method: %s", category, method),
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	elasticJson, _ := json.Marshal(elasticResult.Hits.Hits[0].Source)
	chapter, _ := models.UnmarshalWord(elasticJson)
	response := models.LastChapterResponse{LastChapter: chapter.Chapter}

	middleware.ResponseWithCustomCode(w, 200, response)
}

// swagger:parameters checkQuestion
type checkQuestionParameters struct {
	// in:body
	Application models.CheckAnswerRequest
}

func (s *SokratesHandler) checkAnswer(w http.ResponseWriter, req *http.Request) {
	// swagger:route POST /answer questions checkQuestion
	//
	// Checks whether the provided answer is right or wrong
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
	//	  502: ElasticSearchError
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

		go s.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)
	log.Printf("received %s code with value: %s", plato.HeaderKey, traceID)

	var checkAnswerRequest models.CheckAnswerRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&checkAnswerRequest)
	if err != nil || checkAnswerRequest.AnswerProvided == "" {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "body",
					Message: "error parsing",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	if traceCall {
		jsonCheckSentence, _ := checkAnswerRequest.Marshal()
		span := &pb.SpanRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			Action:       "RequestBodyFromPost",
			RequestBody:  string(jsonCheckSentence),
		}

		go s.Tracer.Span(context.Background(), span)
	}

	query := s.Elastic.Builder().MatchQuery(s.SearchWord, checkAnswerRequest.QuizWord)
	elasticResult, err := s.Elastic.Query().Match(s.Index, query)
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
		go s.databaseSpan(elasticResult, query, traceID, spanID)
	}

	var logoi models.Logos
	answer := models.CheckAnswerResponse{Correct: false, QuizWord: checkAnswerRequest.QuizWord}
	for _, hit := range elasticResult.Hits.Hits {
		elasticJson, _ := json.Marshal(hit.Source)
		logos, _ := models.UnmarshalWord(elasticJson)
		logoi.Logos = append(logoi.Logos, logos)
	}

	for _, logos := range logoi.Logos {
		if logos.Translation == checkAnswerRequest.AnswerProvided {
			answer.Correct = true
		}

		answer.Possibilities = append(answer.Possibilities, logos)
	}

	middleware.ResponseWithCustomCode(w, 200, answer)
}

func (s *SokratesHandler) createQuestion(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /createQuestion questions create
	//
	// Creates a new question from a method - category - chapter combination
	//
	//	Consumes:
	//	- application/json
	//
	//	Produces:
	//	- application/json
	//
	//   Parameters:
	//     + name: method
	//       in: query
	//       description: method to be used
	//		 example: aristophanes
	//       required: true
	//       type: string
	//       format: method
	//		 title: method
	//     + name: category
	//       in: query
	//       description: category to be used
	//		 example: frogs
	//       required: true
	//       type: string
	//       format: category
	//		 title: category
	//     + name: chapter
	//       in: query
	//       description: chapter to be used
	//		 example: 2
	//       required: true
	//       type: string
	//       format: chapter
	//		 title: chapter
	//
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: QuizResponse
	//    400: ValidationError
	//	  405: MethodError
	//	  502: ElasticSearchError
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

		go s.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)
	log.Printf("received %s code with value: %s", plato.HeaderKey, traceID)

	chapter := req.URL.Query().Get("chapter")
	category := req.URL.Query().Get("category")
	method := req.URL.Query().Get("method")

	if category == "" || chapter == "" || method == "" {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "category, chapter, method",
					Message: "cannot be empty",
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	var quiz models.QuizResponse

	mustQuery := []map[string]string{
		{
			Method: method,
		},
		{
			Category: category,
		},
		{
			Chapter: chapter,
		},
	}

	query := s.Elastic.Builder().MultipleMatch(mustQuery)

	elasticResponse, err := s.Elastic.Query().MatchWithScroll(s.Index, query)

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

	var logoi models.Logos
	for _, hit := range elasticResponse.Hits.Hits {
		source, _ := json.Marshal(hit.Source)
		logos, _ := models.UnmarshalWord(source)
		logoi.Logos = append(logoi.Logos, logos)
	}

	if traceCall {
		go s.databaseSpan(elasticResponse, query, traceID, spanID)
	}

	randNumber := s.Randomizer.RandomNumberBaseZero(len(logoi.Logos))

	question := logoi.Logos[randNumber]
	quiz.Question = question.Greek
	quiz.Answer = question.Translation
	quiz.QuizQuestions = append(quiz.QuizQuestions, question.Translation)

	numberOfNeededAnswers := 4

	if len(logoi.Logos) < numberOfNeededAnswers {
		numberOfNeededAnswers = len(logoi.Logos)
	}

	for len(quiz.QuizQuestions) != numberOfNeededAnswers {
		randNumber = s.Randomizer.RandomNumberBaseZero(len(logoi.Logos))
		randEntry := logoi.Logos[randNumber]

		exists := findQuizWord(quiz.QuizQuestions, randEntry.Translation)
		if !exists {
			quiz.QuizQuestions = append(quiz.QuizQuestions, randEntry.Translation)
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(quiz.QuizQuestions), func(i, j int) {
		quiz.QuizQuestions[i], quiz.QuizQuestions[j] = quiz.QuizQuestions[j], quiz.QuizQuestions[i]
	})

	middleware.ResponseWithCustomCode(w, 200, quiz)
}

func (s *SokratesHandler) queryMethods(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /methods methods methods
	//
	// Finds all methods available
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
	//	  200: Methods
	//	  405: MethodError
	//	  502: ElasticSearchError
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

		go s.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)
	log.Printf("received %s code with value: %s", plato.HeaderKey, traceID)

	query := s.Elastic.Builder().Aggregate(Authors, Method)
	elasticResult, err := s.Elastic.Query().MatchAggregate(s.Index, query)

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

	var methods models.Methods
	for _, bucket := range elasticResult.Aggregations.AuthorAggregation.Buckets {
		author := models.Method{Method: strings.ToLower(fmt.Sprintf("%v", bucket.Key))}
		methods.Method = append(methods.Method, author)
	}

	middleware.ResponseWithCustomCode(w, http.StatusOK, methods)
}

func (s *SokratesHandler) queryCategories(w http.ResponseWriter, req *http.Request) {
	// swagger:route GET /methods/{method}/categories methods categories
	//
	// Finds all categories for a method
	//
	//	Consumes:
	//	- application/json
	//
	//	Produces:
	//	- application/json
	//
	//   Parameters:
	//     + name: method
	//       in: query
	//       description: method to be used
	//		 example: aristophanes
	//       required: true
	//       type: string
	//       format: method
	//		 title: method
	//
	//	Schemes: http, https
	//
	//	Responses:
	//	  200: Categories
	//    400: ValidationError
	//	  405: MethodError
	//	  502: ElasticSearchError
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

		go s.Tracer.Trace(context.Background(), traceReceived)
	}

	w.Header().Set(plato.HeaderKey, requestId)
	log.Printf("received %s code with value: %s", plato.HeaderKey, traceID)

	pathParams := mux.Vars(req)
	method := pathParams[Method]

	query := s.Elastic.Builder().FilteredAggregate(Method, method, Categories, Category)
	elasticResult, err := s.Elastic.Query().MatchAggregate(s.Index, query)

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

	var categories models.Categories
	for _, bucket := range elasticResult.Aggregations.CategoryAggregation.Buckets {
		category := models.Category{Category: fmt.Sprintf("%s", bucket.Key)}
		categories.Category = append(categories.Category, category)

	}

	middleware.ResponseWithJson(w, categories)
}

// findQuizWord takes a slice and looks for an element in it
func findQuizWord(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func (s *SokratesHandler) databaseSpan(response *elasticmodels.Response, query map[string]interface{}, traceID, spanID string) {
	parsedResult, _ := json.Marshal(response)
	parsedQuery, _ := json.Marshal(query)
	dataBaseSpan := &pb.DatabaseSpanRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		Action:       "search",
		Query:        string(parsedQuery),
		ResultJson:   string(parsedResult),
	}

	s.Tracer.DatabaseSpan(context.Background(), dataBaseSpan)
}
