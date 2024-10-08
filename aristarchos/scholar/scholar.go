package scholar

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/transform"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	pb "github.com/odysseia-greek/olympia/aristarchos/proto"
	"io"
	"strings"
	"time"
)

const (
	ROOTWORD = "rootWord"
)

func (a *AggregatorServiceImpl) Health(context.Context, *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Health: true,
	}, nil
}

func (a *AggregatorServiceImpl) CreateNewEntry(stream pb.Aristarchos_CreateNewEntryServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.AggregatorStreamResponse{
				Ack: "acknowledged",
			})
		}
		if err != nil {
			return err
		}

		go a.createOrUpdate(in)
	}
}

func (a *AggregatorServiceImpl) createOrUpdate(request *pb.AggregatorCreationRequest) {
	startTime := time.Now()
	splitID := strings.Split(request.TraceId, "+")

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

	parsedWord := transform.RemoveAccents(request.RootWord)
	request.RootWord = parsedWord

	createNewWord := false
	query := a.Elastic.Builder().MatchQuery(ROOTWORD, request.RootWord)
	response, err := a.Elastic.Query().Match(a.Index, query)

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			createNewWord = true
		} else {
			logging.Error(err.Error())
			return
		}
	} else if len(response.Hits.Hits) == 0 {
		createNewWord = true
	}

	if traceCall {
		go func() {
			parsedQuery, _ := json.Marshal(query)
			hits := int64(0)
			took := int64(0)
			if response != nil {
				hits = response.Hits.Total.Value
				took = response.Took
			}
			dataBaseSpan := &pbar.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				SpanId:       spanID,
				RequestType: &pbar.ParabasisRequest_DatabaseSpan{DatabaseSpan: &pbar.DatabaseSpanRequest{
					Action:   "search",
					Query:    string(parsedQuery),
					Hits:     hits,
					TimeTook: took,
				}},
			}

			err := streamer.Send(dataBaseSpan)
			if err != nil {
				logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
			}
		}()
	}

	entry, err := a.mapAndHandleGrammaticalCategories(request)
	if err != nil {
		logging.Error(err.Error())
		return
	}

	if entry.Categories == nil {
		logging.Error(fmt.Sprintf("could not map the word %s to a workable form", parsedWord))
		return
	}

	if createNewWord {
		entryAsJson, _ := json.Marshal(entry)
		createDocument, err := a.Elastic.Index().CreateDocument(a.Index, entryAsJson)
		if err != nil {
			logging.Error(err.Error())
			return
		}

		logging.Debug(fmt.Sprintf("created document with id: %s and rootWord: %s", createDocument.ID, request.RootWord))
		return
	}

	jsonHit, _ := json.Marshal(response.Hits.Hits[0].Source)
	rootWord, _ := UnmarshalRootWordEntry(jsonHit)

	updateNeeded := false
	for i, conjugation := range rootWord.Categories {
		formFound := false
		for _, cform := range conjugation.Forms {
			if cform.Word == entry.Categories[0].Forms[0].Word {
				formFound = true
				break
			}
		}
		if !formFound {
			rootWord.Categories[i].Forms = append(rootWord.Categories[i].Forms, entry.Categories[0].Forms[0])
			updateNeeded = true
		}
		break
	}

	translationFound := false

	for _, trans := range rootWord.Translations {
		if trans == request.Translation {
			translationFound = true
			break
		}
	}

	if !translationFound || rootWord.Translations == nil {
		rootWord.Translations = append(rootWord.Translations, request.Translation)
		updateNeeded = true
	}

	if updateNeeded {
		entryAsJson, _ := json.Marshal(rootWord)
		createDocument, err := a.Elastic.Document().Update(a.Index, response.Hits.Hits[0].ID, entryAsJson)
		if err != nil {
			logging.Error(err.Error())
			return
		}

		if traceCall {
			go func() {
				parabasis := &pbar.ParabasisRequest{
					TraceId:      traceID,
					ParentSpanId: spanID,
					SpanId:       comedy.GenerateSpanID(),
					RequestType: &pbar.ParabasisRequest_Span{
						Span: &pbar.SpanRequest{
							Action: "CloseSpan",
							Took:   fmt.Sprintf("%v", time.Since(startTime)),
							Status: "updated document",
						},
					},
				}
				if err := streamer.Send(parabasis); err != nil {
					logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
				}
			}()
		}

		logging.Debug(fmt.Sprintf("updated document with id: %s and rootWord: %s", createDocument.ID, request.RootWord))
		return
	}

	if traceCall {
		go func() {
			parabasis := &pbar.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				SpanId:       comedy.GenerateSpanID(),
				RequestType: &pbar.ParabasisRequest_Span{
					Span: &pbar.SpanRequest{
						Action: "CloseSpan",
						Took:   fmt.Sprintf("%v", time.Since(startTime)),
						Status: "no update done",
					},
				},
			}
			if err := streamer.Send(parabasis); err != nil {
				logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
			}
		}()
	}

	logging.Debug(fmt.Sprintf("no action needed since document and grammar exists: %s", rootWord))

	return
}

func (a *AggregatorServiceImpl) RetrieveEntry(ctx context.Context, request *pb.AggregatorRequest) (*pb.RootWordResponse, error) {
	startTime := time.Now()
	requestID, ok := ctx.Value(config.DefaultTracingName).(string)
	if !ok {
		logging.Error("could not extract combinedId")
		requestID = "donot+trace+0"
	}

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

	parsedWord := transform.RemoveAccents(request.RootWord)
	request.RootWord = parsedWord

	query := a.Elastic.Builder().MatchQuery(ROOTWORD, request.RootWord)
	response, err := a.Elastic.Query().Match(a.Index, query)

	if traceCall {
		go func() {
			hits := int64(0)
			took := int64(0)
			if response != nil {
				hits = response.Hits.Total.Value
				took = response.Took
			}
			parsedQuery, _ := json.Marshal(query)
			dataBaseSpan := &pbar.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				SpanId:       spanID,
				RequestType: &pbar.ParabasisRequest_DatabaseSpan{DatabaseSpan: &pbar.DatabaseSpanRequest{
					Action:   "search",
					Query:    string(parsedQuery),
					Hits:     hits,
					TimeTook: took,
				}},
			}

			err := streamer.Send(dataBaseSpan)
			if err != nil {
				logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
			}
		}()
	}

	if err != nil {
		return nil, err
	} else if len(response.Hits.Hits) == 0 {
		return nil, fmt.Errorf("no entry can be found")
	}

	var responsePB pb.RootWordResponse
	jsonHit, _ := json.Marshal(response.Hits.Hits[0].Source)
	rootWord, _ := UnmarshalRootWordEntry(jsonHit)

	responsePB.RootWord = rootWord.RootWord
	responsePB.Translations = rootWord.Translations
	responsePB.PartOfSpeech = mapCategoryToEnum(rootWord.PartOfSpeech)
	for _, conj := range rootWord.Categories {
		conjPB := &pb.GrammaticalCategory{}

		for _, form := range conj.Forms {
			formPB := &pb.GrammaticalForm{
				Word: form.Word,
				Rule: form.Rule,
			}
			conjPB.Forms = append(conjPB.Forms, formPB)
		}

		responsePB.Categories = append(responsePB.Categories, conjPB)
	}

	if traceCall {
		go func() {
			parabasis := &pbar.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				SpanId:       comedy.GenerateSpanID(),
				RequestType: &pbar.ParabasisRequest_Span{
					Span: &pbar.SpanRequest{
						Action: "CloseSpan",
						Took:   fmt.Sprintf("%v", time.Since(startTime)),
					},
				},
			}
			if err := streamer.Send(parabasis); err != nil {
				logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
			}
		}()
	}

	return &responsePB, nil
}

func (a *AggregatorServiceImpl) RetrieveSearchWords(ctx context.Context, request *pb.AggregatorRequest) (*pb.SearchWordResponse, error) {
	startTime := time.Now()
	requestID, ok := ctx.Value(config.DefaultTracingName).(string)
	if !ok {
		logging.Error("could not extract combinedId")
		requestID = "donot+trace+0"
	}

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

	parsedWord := transform.RemoveAccents(request.RootWord)
	request.RootWord = parsedWord

	query := a.Elastic.Builder().MatchQuery(ROOTWORD, request.RootWord)
	response, err := a.Elastic.Query().Match(a.Index, query)

	if err != nil {
		return nil, err
	} else if len(response.Hits.Hits) == 0 {
		return nil, fmt.Errorf("no entry can be found")
	}

	if traceCall {
		go func() {
			parsedQuery, _ := json.Marshal(query)
			hits := int64(0)
			took := int64(0)
			if response != nil {
				hits = response.Hits.Total.Value
				took = response.Took
			}

			dataBaseSpan := &pbar.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				SpanId:       spanID,
				RequestType: &pbar.ParabasisRequest_DatabaseSpan{DatabaseSpan: &pbar.DatabaseSpanRequest{
					Action:   "search",
					Query:    string(parsedQuery),
					Hits:     hits,
					TimeTook: took,
				}},
			}

			err := streamer.Send(dataBaseSpan)
			if err != nil {
				logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
			}
		}()
	}

	var responsePB pb.SearchWordResponse
	jsonHit, _ := json.Marshal(response.Hits.Hits[0].Source)
	rootWord, _ := UnmarshalRootWordEntry(jsonHit)

	for _, conj := range rootWord.Categories {
		for _, form := range conj.Forms {
			responsePB.Word = append(responsePB.Word, form.Word)
		}
	}

	if traceCall {
		go func() {
			parabasis := &pbar.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				SpanId:       comedy.GenerateSpanID(),
				RequestType: &pbar.ParabasisRequest_Span{
					Span: &pbar.SpanRequest{
						Action: "CloseSpan",
						Took:   fmt.Sprintf("%v", time.Since(startTime)),
					},
				},
			}
			if err := streamer.Send(parabasis); err != nil {
				logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
			}
		}()
	}

	return &responsePB, nil
}

func (a *AggregatorServiceImpl) RetrieveRootFromGrammarForm(ctx context.Context, request *pb.AggregatorRequest) (*pb.FormsResponse, error) {
	startTime := time.Now()
	requestID, ok := ctx.Value(config.DefaultTracingName).(string)
	if !ok {
		logging.Error("could not extract combinedId")
		requestID = "donot+trace+0"
	}

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
	var responsePB pb.FormsResponse
	responsePB.Word = request.RootWord
	parsedWord := transform.RemoveAccents(request.RootWord)
	request.RootWord = parsedWord
	responsePB.UnaccentedWord = parsedWord

	query := a.Elastic.Builder().MatchQuery("categories.forms.word", request.RootWord)
	response, err := a.Elastic.Query().Match(a.Index, query)

	if err != nil {
		return nil, err
	} else if len(response.Hits.Hits) == 0 {
		return nil, fmt.Errorf("no entry can be found")
	}

	if traceCall {
		go func() {
			parsedQuery, _ := json.Marshal(query)
			hits := int64(0)
			took := int64(0)
			if response != nil {
				hits = response.Hits.Total.Value
				took = response.Took
			}

			dataBaseSpan := &pbar.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				SpanId:       spanID,
				RequestType: &pbar.ParabasisRequest_DatabaseSpan{DatabaseSpan: &pbar.DatabaseSpanRequest{
					Action:   "search",
					Query:    string(parsedQuery),
					Hits:     hits,
					TimeTook: took,
				}},
			}

			err := streamer.Send(dataBaseSpan)
			if err != nil {
				logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
			}
		}()
	}

	jsonHit, _ := json.Marshal(response.Hits.Hits[0].Source)

	logging.Debug(fmt.Sprintf("only taken the first hit: %s", string(jsonHit)))
	rootEntry, _ := UnmarshalRootWordEntry(jsonHit)

	responsePB.RootWord = rootEntry.RootWord
	responsePB.Translation = rootEntry.Translations
	responsePB.PartOfSpeech = rootEntry.PartOfSpeech

	for _, conj := range rootEntry.Categories {
		for _, form := range conj.Forms {
			wordFormForm := transform.RemoveAccents(form.Word)
			if wordFormForm == parsedWord {
				responsePB.Rule = form.Rule
				responsePB.Word = form.Word
			}
		}
	}

	if traceCall {
		go func() {
			parabasis := &pbar.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				SpanId:       comedy.GenerateSpanID(),
				RequestType: &pbar.ParabasisRequest_Span{
					Span: &pbar.SpanRequest{
						Action: "CloseSpan",
						Took:   fmt.Sprintf("%v", time.Since(startTime)),
					},
				},
			}
			if err := streamer.Send(parabasis); err != nil {
				logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
			}
		}()
	}

	return &responsePB, nil
}

func (a *AggregatorServiceImpl) mapAndHandleGrammaticalCategories(request *pb.AggregatorCreationRequest) (*RootWordEntry, error) {
	// example: 3th sing - impf - ind - act
	// example: 1st plur - aor - ind - act
	// example: noun - plural - masc - nom
	// example: pres act part - sing - masc - nom
	// example: inf - pres - act
	var entry RootWordEntry

	entry.RootWord = request.RootWord
	entry.PartOfSpeech = mapEnumToCategory(request.PartOfSpeech)
	entry.Translations = []string{request.Translation}

	conjForm := GrammaticalForm{
		Word: request.Word,
		Rule: request.Rule,
	}
	conj := GrammaticalCategory{
		Forms: []GrammaticalForm{conjForm},
	}

	entry.Categories = append(entry.Categories, conj)

	return &entry, nil
}

func mapCategoryToEnum(category string) pb.PartOfSpeech {
	switch category {
	case "verb":
		return pb.PartOfSpeech_VERB
	case "noun":
		return pb.PartOfSpeech_NOUN
	case "participle":
		return pb.PartOfSpeech_PARTICIPLE
	case "preposition":
		return pb.PartOfSpeech_PREPOSITION
	case "adverb":
		return pb.PartOfSpeech_ADVERB
	case "article":
		return pb.PartOfSpeech_ARTICLE
	case "conjunction":
		return pb.PartOfSpeech_CONJUNCTION
	case "pronoun":
		return pb.PartOfSpeech_PRONOUN
	case "particle":
		return pb.PartOfSpeech_PARTICLE
	default:
		return pb.PartOfSpeech_UNKNOWN_CATEGORY
	}
}

func mapEnumToCategory(category pb.PartOfSpeech) string {
	switch category {
	case pb.PartOfSpeech_VERB:
		return "verb"
	case pb.PartOfSpeech_NOUN:
		return "noun"
	case pb.PartOfSpeech_PARTICIPLE:
		return "participle"
	case pb.PartOfSpeech_PREPOSITION:
		return "preposition"
	case pb.PartOfSpeech_ADVERB:
		return "adverb"
	case pb.PartOfSpeech_ARTICLE:
		return "article"
	case pb.PartOfSpeech_CONJUNCTION:
		return "conjunction"
	case pb.PartOfSpeech_PRONOUN:
		return "pronoun"
	case pb.PartOfSpeech_PARTICLE:
		return "particle"
	default:
		return "UNKNOWN"
	}
}
