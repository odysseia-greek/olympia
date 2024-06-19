package text

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	pab "github.com/odysseia-greek/olympia/aristarchos/proto"
	aristarchos "github.com/odysseia-greek/olympia/aristarchos/scholar"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
	"time"
)

type HerodotosHandler struct {
	Aggregator *aristarchos.ClientAggregator
	Elastic    aristoteles.Client
	Index      string
	Streamer   pb.TraceService_ChorusClient
	Cancel     context.CancelFunc
	Cache      archytas.Client
}

const (
	AuthorReq    string = "author"
	BookReq      string = "book"
	ReferenceReq string = "reference"
	SectionReq   string = "section"
	RootWord     string = "rootword"
	Options      string = "texts/options"
	Id           string = "_id"
)

// PingPong pongs the ping
func (h *HerodotosHandler) pingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

// returns the health of the api
func (h *HerodotosHandler) health(w http.ResponseWriter, req *http.Request) {
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
func (h *HerodotosHandler) create(w http.ResponseWriter, req *http.Request) {
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

	var createTextRequest models.CreateTextRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&createTextRequest)
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

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							AuthorReq: map[string]interface{}{
								"query":         createTextRequest.Author,
								"operator":      "and",
								"fuzziness":     "AUTO",
								"prefix_length": 0,
							},
						},
					},
					{
						"match": map[string]interface{}{
							BookReq: createTextRequest.Book,
						},
					},
					{
						"match": map[string]interface{}{
							ReferenceReq: createTextRequest.Reference,
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
		hits := int64(0)
		took := int64(0)
		if response != nil {
			hits = response.Hits.Total.Value
			took = response.Took
		}
		go h.databaseSpan(hits, took, query, traceID, spanID)
	}

	if len(response.Hits.Hits) == 0 {
		e := models.NotFoundError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Message: models.NotFoundMessage{
				Type:   fmt.Sprintf("author: %s and book: %s", createTextRequest.Author, createTextRequest.Book),
				Reason: "no hits for combination",
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	} else if len(response.Hits.Hits) > 1 {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "response.Hits.Hits",
					Message: fmt.Sprintf("more than one entry found for author: %s book: %s and reference: %s", createTextRequest.Author, createTextRequest.Book, createTextRequest.Reference),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	elasticJson, _ := json.Marshal(response.Hits.Hits[0].Source)
	var text models.Text
	err = json.Unmarshal(elasticJson, &text)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "unmarshal json",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	if createTextRequest.Section != "" {
		for _, rhema := range text.Rhemai {
			if rhema.Section == createTextRequest.Section {
				var sectionedText = models.Text{
					Author:          text.Author,
					Book:            text.Book,
					Type:            text.Type,
					Reference:       text.Reference,
					PerseusTextLink: text.PerseusTextLink,
					Rhemai: []models.Rhema{
						rhema,
					},
				}

				middleware.ResponseWithCustomCode(w, http.StatusOK, sectionedText)
				return
			}
		}
	}

	middleware.ResponseWithCustomCode(w, http.StatusOK, text)
}

// checks the validity of an answer
func (h *HerodotosHandler) check(w http.ResponseWriter, req *http.Request) {
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

	var checkTextRequest models.CheckTextRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&checkTextRequest)
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

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							AuthorReq: map[string]interface{}{
								"query":         checkTextRequest.Author,
								"operator":      "and",
								"fuzziness":     "AUTO",
								"prefix_length": 0,
							},
						},
					},
					{
						"match": map[string]interface{}{
							BookReq: checkTextRequest.Book,
						},
					},
					{
						"match": map[string]interface{}{
							ReferenceReq: checkTextRequest.Reference,
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
		hits := int64(0)
		took := int64(0)
		if response != nil {
			hits = response.Hits.Total.Value
			took = response.Took
		}
		go h.databaseSpan(hits, took, query, traceID, spanID)
	}

	if len(response.Hits.Hits) == 0 {
		e := models.NotFoundError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Message: models.NotFoundMessage{
				Type:   fmt.Sprintf("author: %s and book: %s", checkTextRequest.Author, checkTextRequest.Book),
				Reason: "no hits for combination",
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	} else if len(response.Hits.Hits) > 1 {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "response.Hits.Hits",
					Message: fmt.Sprintf("more than one entry found for author: %s book: %s and reference: %s", checkTextRequest.Author, checkTextRequest.Book, checkTextRequest.Reference),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	elasticJson, _ := json.Marshal(response.Hits.Hits[0].Source)
	var text models.Text
	err = json.Unmarshal(elasticJson, &text)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "unmarshal json",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	var answers []models.AnswerSection
	var possibleTypos []models.Typo
	var totalLevenshtein float64

	for _, section := range checkTextRequest.Translations {
		for _, sect := range text.Rhemai {
			if section.Section == sect.Section {
				for _, translation := range sect.Translations {
					levenshteinDist := levenshteinDistance(translation, section.Translation)
					lenOfLongestSentence := longestStringOfTwo(translation, section.Translation)
					levenshteinPerc := levenshteinDistanceInPercentage(levenshteinDist, lenOfLongestSentence)
					roundedPercentage := fmt.Sprintf("%.2f", levenshteinPerc)
					answerSection := models.AnswerSection{
						Section:               sect.Section,
						LevenshteinPercentage: roundedPercentage,
						QuizSentence:          translation,
						AnswerSentence:        section.Translation,
					}

					if levenshteinPerc < 100.0 {
						typos := findTypos(section.Translation, translation)
						possibleTypos = append(possibleTypos, typos...)
					}

					totalLevenshtein += levenshteinPerc

					answers = append(answers, answerSection)
				}
			}
		}
	}

	averageLevenshtein := fmt.Sprintf("%.2f", totalLevenshtein/float64(len(answers)))

	result := models.CheckTextResponse{
		AverageLevenshteinPercentage: averageLevenshtein,
		Sections:                     answers,
		PossibleTypos:                possibleTypos,
	}

	middleware.ResponseWithCustomCode(w, http.StatusOK, result)
}

func (h *HerodotosHandler) options(w http.ResponseWriter, req *http.Request) {
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

	cacheItem, _ := h.Cache.Read(Options)
	if cacheItem != nil {
		var agg models.AggregationResult
		err := json.Unmarshal(cacheItem, &agg)
		if err != nil {
			e := models.ValidationError{
				ErrorModel: models.ErrorModel{UniqueCode: requestId},
				Messages: []models.ValidationMessages{
					{
						Field:   "unmarshal json cache",
						Message: err.Error(),
					},
				},
			}

			middleware.ResponseWithJson(w, e)
			return
		}

		if traceCall {
			parabasis := &pb.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				SpanId:       comedy.GenerateSpanID(),
				RequestType: &pb.ParabasisRequest_Span{
					Span: &pb.SpanRequest{
						Action: "TakenFromCache",
						Status: fmt.Sprintf("status code: %d", http.StatusOK),
					},
				},
			}
			if err := h.Streamer.Send(parabasis); err != nil {
				logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
			}
		}

		middleware.ResponseWithCustomCode(w, http.StatusOK, agg)
		return
	}

	query := textAggregationQuery()

	elasticResult, err := h.Elastic.Query().MatchRaw(h.Index, query)
	if err != nil {
		e := models.ElasticSearchError{
			ErrorModel: models.ErrorModel{UniqueCode: requestId},
			Message: models.ElasticErrorMessage{
				ElasticError: err.Error(),
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	var agg map[string]interface{}
	err = json.Unmarshal(elasticResult, &agg)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: requestId},
			Messages: []models.ValidationMessages{
				{
					Field:   "unmarshall action failed internally",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	result, err := parseAggregationResults(agg)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: requestId},
			Messages: []models.ValidationMessages{
				{
					Field:   "querying aggregator failed",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	ttl := time.Hour
	setCache, err := json.Marshal(result)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: requestId},
			Messages: []models.ValidationMessages{
				{
					Field:   "unmarshall action failed internally",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	err = h.Cache.SetWithTTL(Options, string(setCache), ttl)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: requestId},
			Messages: []models.ValidationMessages{
				{
					Field:   "cache",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	middleware.ResponseWithCustomCode(w, http.StatusOK, result)
}

// analyseText fetches words and queries them in all the texts
func (h *HerodotosHandler) analyze(w http.ResponseWriter, req *http.Request) {
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

	var analyzeTextRequest models.AnalyzeTextRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&analyzeTextRequest)
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

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	md := metadata.New(map[string]string{service.HeaderKey: requestId})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	aggregatorRequest := pab.AggregatorRequest{RootWord: analyzeTextRequest.Rootword}
	entry, err := h.Aggregator.RetrieveEntry(ctx, &aggregatorRequest)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "querying aggregator failed",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	var conjugations []models.Conjugations
	var words []string
	var rootWordfound bool

	for _, category := range entry.Categories {
		for _, form := range category.Forms {
			if form.Word == entry.RootWord {
				rootWordfound = true
			}
			words = append(words, form.Word)
			conjugation := models.Conjugations{
				Word: form.Word,
				Rule: form.Rule,
			}
			conjugations = append(conjugations, conjugation)
		}
	}

	if !rootWordfound {
		words = append(words, entry.RootWord)
	}

	query := createGreekTextQuery(words)

	logging.Debug(fmt.Sprintf("%v", query))
	response, err := h.Elastic.Query().MatchWithScroll(h.Index, query)
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

	if len(response.Hits.Hits) == 0 {
		e := models.NotFoundError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Message: models.NotFoundMessage{
				Type:   fmt.Sprintf("words: %v", words),
				Reason: "no hits",
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	if traceCall {
		hits := int64(0)
		took := int64(0)
		if response != nil {
			hits = response.Hits.Total.Value
			took = response.Took
		}
		go h.databaseSpan(hits, took, query, traceID, spanID)
	}

	result := models.AnalyzeTextResponse{
		Rootword:     analyzeTextRequest.Rootword,
		Conjugations: conjugations,
	}

	for _, hit := range response.Hits.Hits {
		elasticJson, _ := json.Marshal(hit.Source)
		var text models.Text
		err = json.Unmarshal(elasticJson, &text)
		if err != nil {
			e := models.ValidationError{
				ErrorModel: models.ErrorModel{UniqueCode: traceID},
				Messages: []models.ValidationMessages{
					{
						Field:   "unmarshal json",
						Message: err.Error(),
					},
				},
			}
			middleware.ResponseWithJson(w, e)
			return
		}

		for _, section := range text.Rhemai {
			var processedWords []string
			wordFound := false // Variable to track if the word is found
			for _, word := range words {
				processedInLoop := false
				for _, processed := range processedWords {
					if processed == word {
						processedInLoop = true
					}
				}
				if !processedInLoop {
					if strings.Contains(section.Greek, word) {
						wordFound = true // Mark that the word is found in the section
						highlight := fmt.Sprintf("&&&%s&&&", word)
						section.Greek = strings.ReplaceAll(section.Greek, word, highlight)
					}
					processedWords = append(processedWords, word)
				}
			}

			// Only add the result if the word was found
			if wordFound {
				res := models.AnalyzeResult{
					ReferenceLink: fmt.Sprintf("/texts?author=%s&book=%s&reference=%s", text.Author, text.Book, text.Reference),
					Author:        text.Author,
					Book:          text.Book,
					Reference:     text.Reference,
					Text:          section,
				}
				result.Results = append(result.Results, res)
			}
		}
	}

	middleware.ResponseWithCustomCode(w, 200, result)
}

func (h *HerodotosHandler) databaseSpan(hits, took int64, query map[string]interface{}, traceID, spanID string) {
	parsedQuery, _ := json.Marshal(query)

	dataBaseSpan := &pb.ParabasisRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		SpanId:       spanID,
		RequestType: &pb.ParabasisRequest_DatabaseSpan{DatabaseSpan: &pb.DatabaseSpanRequest{
			Action:   "search",
			Query:    string(parsedQuery),
			Hits:     hits,
			TimeTook: took,
		}},
	}

	err := h.Streamer.Send(dataBaseSpan)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}
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

// Function to find typos between two sentences
func findTypos(provided, source string) []models.Typo {
	s := removeCharacters(source, ",`~<>/?!.;:'\"")
	p := removeCharacters(provided, ",`~<>/?!.;:'\"")

	sourceSentence := strings.Split(s, " ")
	providedSentence := strings.Split(p, " ")

	var possibleTypos []models.Typo

	for _, wordProvided := range providedSentence {
		wordMatched := false

		for _, wordInSource := range sourceSentence {
			sourceWord := strings.ToLower(wordInSource)
			providedWord := strings.ToLower(wordProvided)

			levenshtein := levenshteinDistance(providedWord, sourceWord)

			// If an exact match is found, break out of the loop
			if levenshtein == 0 {
				wordMatched = true
				break
			}

			// Check for a typo with a Levenshtein distance of 1
			if levenshtein == 1 {
				possibleTypo := models.Typo{
					Source:   wordInSource,
					Provided: wordProvided,
				}
				possibleTypos = append(possibleTypos, possibleTypo)
			}

			// Check for a possible typo with longer words and a small Levenshtein distance
			if levenshtein > 1 && len(sourceWord) > 10 && levenshtein <= 3 {
				possibleTypo := models.Typo{
					Source:   wordInSource,
					Provided: wordProvided,
				}
				possibleTypos = append(possibleTypos, possibleTypo)
			}
		}

		// If an exact match was found for the word, remove any typos added for this word
		if wordMatched {
			for i := len(possibleTypos) - 1; i >= 0; i-- {
				if strings.ToLower(possibleTypos[i].Provided) == strings.ToLower(wordProvided) {
					possibleTypos = append(possibleTypos[:i], possibleTypos[i+1:]...)
				}
			}
		}
	}

	return possibleTypos
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
