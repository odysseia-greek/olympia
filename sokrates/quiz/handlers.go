package quiz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	elasticmodels "github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/agora/plato/transform"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SokratesHandler struct {
	Elastic            elastic.Client
	Randomizer         randomizer.Random
	Client             service.OdysseiaClient
	SearchWord         string
	Index              string
	QuizAttempts       chan models.QuizAttempt
	AggregatedAttempts map[string]models.QuizAttempt
	Ticker             *time.Ticker
	Streamer           pb.TraceService_ChorusClient
	Cancel             context.CancelFunc
}

const (
	THEME    string = "theme"
	SET      string = "set"
	QUIZTYPE string = "quizType"
)

func (s *SokratesHandler) Health(w http.ResponseWriter, req *http.Request) {
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

func (s *SokratesHandler) Create(w http.ResponseWriter, req *http.Request) {
	var requestId string
	fromContext := req.Context().Value(config.DefaultTracingName)
	if fromContext == nil {
		requestId = req.Header.Get(config.HeaderKey)
	} else {
		requestId = fromContext.(string)
	}

	var createQuizRequest models.CreationRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&createQuizRequest)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: requestId},
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

	switch createQuizRequest.QuizType {
	case models.MEDIA:
		quiz, err := s.mediaQuiz(createQuizRequest.Theme, createQuizRequest.Set, requestId)
		if err != nil {
			e := models.ValidationError{
				ErrorModel: models.ErrorModel{UniqueCode: requestId},
				Messages: []models.ValidationMessages{
					{
						Field:   "creating quiz error",
						Message: err.Error(),
					},
				},
			}
			middleware.ResponseWithJson(w, e)
			return
		}

		middleware.ResponseWithCustomCode(w, http.StatusOK, quiz)
		return
	case models.AUTHORBASED:
		quiz, err := s.authorBasedQuiz(createQuizRequest.Theme, createQuizRequest.Set, requestId)
		if err != nil {
			e := models.ValidationError{
				ErrorModel: models.ErrorModel{UniqueCode: requestId},
				Messages: []models.ValidationMessages{
					{
						Field:   "creating quiz error",
						Message: err.Error(),
					},
				},
			}
			middleware.ResponseWithJson(w, e)
			return
		}

		middleware.ResponseWithCustomCode(w, http.StatusOK, quiz)
		return
	case models.DIALOGUE:
		quiz, err := s.dialogueQuiz(createQuizRequest.Theme, createQuizRequest.Set, requestId)
		if err != nil {
			e := models.ValidationError{
				ErrorModel: models.ErrorModel{UniqueCode: requestId},
				Messages: []models.ValidationMessages{
					{
						Field:   "creating quiz error",
						Message: err.Error(),
					},
				},
			}
			middleware.ResponseWithJson(w, e)
			return
		}

		middleware.ResponseWithCustomCode(w, http.StatusOK, quiz)
		return
	}
}

func (s *SokratesHandler) Check(w http.ResponseWriter, req *http.Request) {
	var requestId string
	fromContext := req.Context().Value(config.DefaultTracingName)
	if fromContext == nil {
		requestId = req.Header.Get(config.HeaderKey)
	} else {
		requestId = fromContext.(string)
	}

	var answerRequest models.AnswerRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&answerRequest)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: requestId},
			Messages: []models.ValidationMessages{
				{
					Field:   "decoding",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithCustomCode(w, 400, e)
		return
	}

	switch answerRequest.QuizType {
	case models.MEDIA:
		quiz, err := s.mediaQuizAnswer(answerRequest, requestId)
		if err != nil {
			e := models.ValidationError{
				ErrorModel: models.ErrorModel{UniqueCode: requestId},
				Messages: []models.ValidationMessages{
					{
						Field:   "creating quiz error",
						Message: err.Error(),
					},
				},
			}
			middleware.ResponseWithCustomCode(w, 400, e)
			return
		}

		middleware.ResponseWithCustomCode(w, 200, quiz)
		return
	case models.AUTHORBASED:
		quiz, err := s.authorBasedQuizAnswer(answerRequest, requestId)
		if err != nil {
			e := models.ValidationError{
				ErrorModel: models.ErrorModel{UniqueCode: requestId},
				Messages: []models.ValidationMessages{
					{
						Field:   "creating quiz error",
						Message: err.Error(),
					},
				},
			}
			middleware.ResponseWithCustomCode(w, 400, e)
			return
		}

		middleware.ResponseWithCustomCode(w, http.StatusOK, quiz)
		return
	case models.DIALOGUE:
		quiz, err := s.dialogueAnswer(answerRequest, requestId)
		if err != nil {
			e := models.ValidationError{
				ErrorModel: models.ErrorModel{UniqueCode: requestId},
				Messages: []models.ValidationMessages{
					{
						Field:   "creating quiz error",
						Message: err.Error(),
					},
				},
			}
			middleware.ResponseWithCustomCode(w, 400, e)
			return
		}

		middleware.ResponseWithCustomCode(w, http.StatusOK, quiz)
		return
	}
}

func (s *SokratesHandler) Options(w http.ResponseWriter, req *http.Request) {
	var requestId string
	fromContext := req.Context().Value(config.DefaultTracingName)
	if fromContext == nil {
		requestId = req.Header.Get(config.HeaderKey)
	} else {
		requestId = fromContext.(string)
	}

	quizType := req.URL.Query().Get("quizType")

	options, err := s.options(quizType)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: requestId},
			Messages: []models.ValidationMessages{
				{
					Field:   "creating quiz error",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	middleware.ResponseWithCustomCode(w, http.StatusOK, options)
	return

}

func (s *SokratesHandler) mediaQuizAnswer(req models.AnswerRequest, requestID string) (*models.ComprehensiveResponse, error) {
	splitID := strings.Split(requestID, "+")

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

	mustQuery := []map[string]string{
		{
			QUIZTYPE: models.MEDIA,
		},
		{
			THEME: req.Theme,
		},
		{
			SET: req.Set,
		},
	}

	query := s.Elastic.Builder().MultipleMatch(mustQuery)
	elasticResponse, err := s.Elastic.Query().Match(s.Index, query)
	if err != nil {
		return nil, err
	}
	if len(elasticResponse.Hits.Hits) == 0 {
		return nil, fmt.Errorf("no hits found in Elastic")
	}

	if traceCall {
		go s.databaseSpan(elasticResponse, query, traceID, spanID)
	}

	var option models.MediaQuiz
	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &option)
	if err != nil {
		return nil, err
	}

	answer := models.ComprehensiveResponse{Correct: false, QuizWord: req.QuizWord}

	if req.Comprehensive {
		s.gatherComprehensiveData(&answer, requestID)
	}

	for _, content := range option.Content {
		if content.Greek == req.QuizWord {
			if content.Translation == req.Answer {
				answer.Correct = true
			}
		}
	}

	return &answer, nil
}

func (s *SokratesHandler) authorBasedQuizAnswer(req models.AnswerRequest, requestID string) (*models.ComprehensiveResponse, error) {
	splitID := strings.Split(requestID, "+")

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

	mustQuery := []map[string]string{
		{
			QUIZTYPE: models.AUTHORBASED,
		},
		{
			THEME: req.Theme,
		},
		{
			SET: req.Set,
		},
	}

	query := s.Elastic.Builder().MultipleMatch(mustQuery)
	elasticResponse, err := s.Elastic.Query().Match(s.Index, query)
	if err != nil {
		return nil, err
	}
	if len(elasticResponse.Hits.Hits) == 0 {
		return nil, fmt.Errorf("no hits found in Elastic")
	}

	if traceCall {
		go s.databaseSpan(elasticResponse, query, traceID, spanID)
	}

	var option models.AuthorBasedQuiz
	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &option)
	if err != nil {
		return nil, err
	}

	answer := models.ComprehensiveResponse{Correct: false, QuizWord: req.QuizWord}

	if req.Comprehensive {
		s.gatherComprehensiveData(&answer, requestID)
	}

	for _, content := range option.Content {
		if content.Greek == req.QuizWord {
			if content.Translation == req.Answer {
				answer.Correct = true
			}
		}
	}

	s.QuizAttempts <- models.QuizAttempt{Correct: answer.Correct, Set: req.Set, Theme: req.Theme}
	answer.Progress.AverageAccuracy = option.Progress.AverageAccuracy
	answer.Progress.TimesCorrect = option.Progress.TimesCorrect
	answer.Progress.TimesIncorrect = option.Progress.TimesIncorrect

	return &answer, nil
}

func (s *SokratesHandler) dialogueAnswer(req models.AnswerRequest, requestID string) (*models.DialogueAnswer, error) {
	splitID := strings.Split(requestID, "+")

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

	mustQuery := []map[string]string{
		{
			QUIZTYPE: models.DIALOGUE,
		},
		{
			THEME: req.Theme,
		},
		{
			SET: req.Set,
		},
	}

	query := s.Elastic.Builder().MultipleMatch(mustQuery)
	elasticResponse, err := s.Elastic.Query().Match(s.Index, query)
	if err != nil {
		return nil, err
	}

	if len(elasticResponse.Hits.Hits) == 0 {
		return nil, fmt.Errorf("no hits found in Elastic")
	}

	if traceCall {
		go s.databaseSpan(elasticResponse, query, traceID, spanID)
	}

	var option models.DialogueQuiz
	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &option)
	if err != nil {
		return nil, err
	}

	answer := models.DialogueAnswer{
		Percentage:   0,
		Input:        req.Dialogue,
		Answer:       option.Content,
		InWrongPlace: []models.DialogueCorrection{},
	}
	var correctPlace int
	var wrongPlace int

	for _, dialogue := range req.Dialogue {
		verifiedContent := option.Content[dialogue.Place-1]
		if verifiedContent.Greek == dialogue.Greek && verifiedContent.Place == dialogue.Place {
			correctPlace++
		} else {
			correctedPlacing := models.DialogueCorrection{
				Translation:  dialogue.Translation,
				Greek:        dialogue.Greek,
				Place:        dialogue.Place,
				Speaker:      dialogue.Speaker,
				CorrectPlace: 0,
			}

			for _, corrected := range option.Content {
				if corrected.Greek == dialogue.Greek && corrected.Speaker == dialogue.Speaker {
					correctedPlacing.CorrectPlace = corrected.Place
				}
			}

			answer.InWrongPlace = append(answer.InWrongPlace, correctedPlacing)
			wrongPlace++
		}
	}

	total := correctPlace + wrongPlace
	totalProgress := 0.0
	if total > 0 {
		totalProgress = float64(correctPlace) / float64(total) * 100
	}

	answer.Percentage = totalProgress

	return &answer, nil
}

func (s *SokratesHandler) gatherComprehensiveData(answer *models.ComprehensiveResponse, requestID string) {
	splitID := strings.Split(requestID, "+")

	traceCall := false
	var traceID, parentSpanID string

	if len(splitID) >= 3 {
		traceCall = splitID[2] == "1"
	}

	if len(splitID) >= 1 {
		traceID = splitID[0]
	}
	if len(splitID) >= 2 {
		parentSpanID = splitID[1]
	}

	wordToBeSend := extractBaseWord(answer.QuizWord)

	// Use a WaitGroup to wait for both goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Buffered channels to capture 1 response
	foundInTextChan := make(chan *http.Response, 1)
	similarWordsChan := make(chan *http.Response, 1)
	errChan := make(chan error, 2) // Buffered to hold potential errors from both calls

	go func() {
		defer wg.Done()
		if traceCall {
			herodotosSpan := &pb.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: parentSpanID,
				SpanId:       comedy.GenerateSpanID(),
				RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
					Action: "analyseText",
					Status: fmt.Sprintf("querying Herodotos for word: %s", wordToBeSend),
				}},
			}

			err := s.Streamer.Send(herodotosSpan)
			if err != nil {
				logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
			}
		}
		foundInText, err := s.Client.Herodotos().AnalyseText(wordToBeSend, requestID)
		if err != nil {
			logging.Error(fmt.Sprintf("could not query any texts for word: %s error: %s", answer.QuizWord, err.Error()))
			errChan <- err
			return
		}
		foundInTextChan <- foundInText
	}()

	go func() {
		defer wg.Done()
		if traceCall {
			alexandrosSpan := &pb.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: parentSpanID,
				SpanId:       comedy.GenerateSpanID(),
				RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
					Action: "analyseText",
					Status: fmt.Sprintf("querying Alexandros for word: %s", wordToBeSend),
				}},
			}

			err := s.Streamer.Send(alexandrosSpan)
			if err != nil {
				logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
			}
		}
		similarWords, err := s.Client.Alexandros().Search(wordToBeSend, "greek", "fuzzy", "false", requestID)
		if err != nil {
			logging.Error(fmt.Sprintf("could not query any similar words for word: %s error: %s", answer.QuizWord, err.Error()))
			errChan <- err
			return
		}
		similarWordsChan <- similarWords
	}()

	// Wait for both goroutines to complete
	wg.Wait()

	// Process responses
	close(errChan)
	close(foundInTextChan)
	close(similarWordsChan)

	for err := range errChan {
		logging.Error(err.Error())
	}

	for foundInText := range foundInTextChan {
		defer foundInText.Body.Close()
		err := json.NewDecoder(foundInText.Body).Decode(&answer.FoundInText)
		if err != nil {
			logging.Error(fmt.Sprintf("error while decoding: %s", err.Error()))
		}
	}

	for similarWords := range similarWordsChan {
		defer similarWords.Body.Close()
		var extended models.ExtendedResponse
		err := json.NewDecoder(similarWords.Body).Decode(&extended)
		if err != nil {
			logging.Error(fmt.Sprintf("error while decoding: %s", err.Error()))
		}

		for _, meros := range extended.Hits {
			answer.SimilarWords = append(answer.SimilarWords, meros.Hit)
		}
	}
}

func (s *SokratesHandler) mediaQuiz(theme, set, requestID string) (*models.QuizResponse, error) {
	splitID := strings.Split(requestID, "+")

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

	mustQuery := []map[string]string{
		{
			THEME: theme,
		},
		{
			SET: set,
		},
		{
			QUIZTYPE: models.MEDIA,
		},
	}

	query := s.Elastic.Builder().MultipleMatch(mustQuery)
	elasticResponse, err := s.Elastic.Query().Match(s.Index, query)
	if err != nil {
		return nil, err
	}

	if elasticResponse.Hits.Hits == nil || len(elasticResponse.Hits.Hits) == 0 {
		return nil, errors.New("no hits found in query")
	}

	var option models.MediaQuiz

	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &option)
	if err != nil {
		return nil, err
	}

	if traceCall {
		go s.databaseSpan(elasticResponse, query, traceID, spanID)
	}

	var quiz models.QuizResponse
	randNumber := s.Randomizer.RandomNumberBaseZero(len(option.Content))

	question := option.Content[randNumber]
	quiz.QuizItem = question.Greek
	quiz.Options = append(quiz.Options, models.Options{
		Option:   question.Translation,
		ImageUrl: question.ImageURL,
	})

	numberOfNeededAnswers := 4

	if len(option.Content) < numberOfNeededAnswers {
		numberOfNeededAnswers = len(option.Content)
	}

	for len(quiz.Options) != numberOfNeededAnswers {
		randNumber = s.Randomizer.RandomNumberBaseZero(len(option.Content))
		randEntry := option.Content[randNumber]

		exists := findQuizWord(quiz.Options, randEntry.Translation)
		if !exists {
			option := models.Options{
				Option:   randEntry.Translation,
				ImageUrl: randEntry.ImageURL,
			}
			quiz.Options = append(quiz.Options, option)
		}
	}

	rand.Shuffle(len(quiz.Options), func(i, j int) {
		quiz.Options[i], quiz.Options[j] = quiz.Options[j], quiz.Options[i]
	})

	return &quiz, nil
}

func (s *SokratesHandler) authorBasedQuiz(theme, set, requestID string) (*models.QuizResponse, error) {
	splitID := strings.Split(requestID, "+")

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

	mustQuery := []map[string]string{
		{
			THEME: theme,
		},
		{
			SET: set,
		},
		{
			QUIZTYPE: models.AUTHORBASED,
		},
	}

	query := s.Elastic.Builder().MultipleMatch(mustQuery)
	elasticResponse, err := s.Elastic.Query().Match(s.Index, query)
	if err != nil {
		return nil, err
	}

	if elasticResponse.Hits.Hits == nil || len(elasticResponse.Hits.Hits) == 0 {
		return nil, errors.New("no hits found in query")
	}

	var option models.AuthorBasedQuiz
	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &option)
	if err != nil {
		return nil, err
	}

	if traceCall {
		go s.databaseSpan(elasticResponse, query, traceID, spanID)
	}

	var quiz models.QuizResponse
	randNumber := s.Randomizer.RandomNumberBaseZero(len(option.Content))

	question := option.Content[randNumber]
	quiz.QuizItem = question.Greek
	quiz.Options = append(quiz.Options, models.Options{
		Option: question.Translation,
	})

	numberOfNeededAnswers := 4

	if len(option.Content) < numberOfNeededAnswers {
		numberOfNeededAnswers = len(option.Content)
	}

	for len(quiz.Options) != numberOfNeededAnswers {
		randNumber = s.Randomizer.RandomNumberBaseZero(len(option.Content))
		randEntry := option.Content[randNumber]

		exists := findQuizWord(quiz.Options, randEntry.Translation)
		if !exists {
			option := models.Options{
				Option: randEntry.Translation,
			}
			quiz.Options = append(quiz.Options, option)
		}
	}

	rand.Shuffle(len(quiz.Options), func(i, j int) {
		quiz.Options[i], quiz.Options[j] = quiz.Options[j], quiz.Options[i]
	})

	return &quiz, nil
}

func (s *SokratesHandler) dialogueQuiz(theme, set, requestID string) (*models.DialogueQuiz, error) {
	splitID := strings.Split(requestID, "+")

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

	mustQuery := []map[string]string{
		{
			THEME: theme,
		},
		{
			SET: set,
		},
		{
			QUIZTYPE: models.DIALOGUE,
		},
	}

	query := s.Elastic.Builder().MultipleMatch(mustQuery)
	elasticResponse, err := s.Elastic.Query().Match(s.Index, query)
	if err != nil {
		return nil, err
	}

	if elasticResponse.Hits.Hits == nil || len(elasticResponse.Hits.Hits) == 0 {
		return nil, errors.New("no hits found in query")
	}

	if traceCall {
		go s.databaseSpan(elasticResponse, query, traceID, spanID)
	}

	var quiz models.DialogueQuiz

	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &quiz)
	if err != nil {
		return nil, err
	}

	return &quiz, nil
}

func (s *SokratesHandler) options(quizType string) (*models.AggregateResult, error) {
	query := s.Elastic.Builder().FilteredAggregate(QUIZTYPE, quizType, THEME, THEME)
	//if quizType == models.MEDIA {
	//	query = map[string]interface{}{
	//		"query": map[string]interface{}{
	//			"match_phrase": map[string]interface{}{
	//				QUIZTYPE: quizType,
	//			},
	//		},
	//		"size": 0,
	//		"aggs": map[string]interface{}{
	//			SET: map[string]interface{}{
	//				"max": map[string]interface{}{
	//					"field": SET,
	//				},
	//			},
	//		},
	//	}
	//}

	elasticResult, err := s.Elastic.Query().MatchAggregate(s.Index, query)
	if err != nil {
		return nil, err
	}

	var result models.AggregateResult

	for _, bucket := range elasticResult.Aggregations.ThemeAggregation.Buckets {
		aggregate := models.Aggregate{
			HighestSet: strconv.Itoa(int(bucket.DocCount)),
			Name:       bucket.Key.(string),
		}

		result.Aggregates = append(result.Aggregates, aggregate)
	}

	if len(elasticResult.Aggregations.ThemeAggregation.Buckets) == 0 {
		aggregate := models.Aggregate{
			HighestSet: fmt.Sprintf("%v", elasticResult.Aggregations.SetAggregation.Value),
		}

		result.Aggregates = append(result.Aggregates, aggregate)
	}

	return &result, nil
}

func (s *SokratesHandler) updateElasticsearch() {
	for {
		select {
		case <-s.Ticker.C:
			for key, attempt := range s.AggregatedAttempts {
				err := s.performUpdate(attempt.Set, attempt.Theme, attempt.Correct)
				if err != nil {
					logging.Error(err.Error())
				}
				delete(s.AggregatedAttempts, key)
			}
		case attempt := <-s.QuizAttempts:
			key := fmt.Sprintf("%s-%s", attempt.Set, attempt.Theme)
			s.AggregatedAttempts[key] = attempt
		}
	}
}

func (s *SokratesHandler) performUpdate(set, theme string, correct bool) error {
	mustQuery := []map[string]string{
		{
			THEME: theme,
		},
		{
			SET: set,
		},
		{
			QUIZTYPE: models.AUTHORBASED,
		},
	}

	query := s.Elastic.Builder().MultipleMatch(mustQuery)
	elasticResponse, err := s.Elastic.Query().Match(s.Index, query)
	if err != nil {
		return err
	}

	var option models.AuthorBasedQuiz
	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &option)
	if err != nil {
		return err
	}

	if correct {
		option.Progress.TimesCorrect += 1
	} else {
		option.Progress.TimesIncorrect += 1
	}

	total := option.Progress.TimesCorrect + option.Progress.TimesIncorrect
	totalProgress := 0.0
	if total > 0 {
		totalProgress = float64(option.Progress.TimesCorrect) / float64(total) * 100
	}

	option.Progress.AverageAccuracy = roundToTwoDecimals(totalProgress)

	entryAsJson, err := json.Marshal(option)
	if err != nil {
		return err
	}
	_, err = s.Elastic.Document().Update(s.Index, elasticResponse.Hits.Hits[0].ID, entryAsJson)
	if err != nil {
		return err
	}

	return nil
}

func roundToTwoDecimals(f float64) float64 {
	return math.Round(f*100) / 100
}

// findQuizWord takes a slice and looks for an element in it
func findQuizWord(slice []models.Options, val string) bool {
	for _, item := range slice {
		if item.Option == val {
			return true
		}
	}
	return false
}

func extractBaseWord(queryWord string) string {
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

func (s *SokratesHandler) databaseSpan(response *elasticmodels.Response, query map[string]interface{}, traceID, spanID string) {
	parsedQuery, _ := json.Marshal(query)
	hits := int64(0)
	if response != nil {
		hits = response.Hits.Total.Value
	}
	dataBaseSpan := &pb.ParabasisRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		SpanId:       spanID,
		RequestType: &pb.ParabasisRequest_DatabaseSpan{DatabaseSpan: &pb.DatabaseSpanRequest{
			Action:   "search",
			Query:    string(parsedQuery),
			Hits:     hits,
			TimeTook: response.Took,
		}},
	}

	err := s.Streamer.Send(dataBaseSpan)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}
}
