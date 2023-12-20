package scholar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	pb "github.com/odysseia-greek/olympia/aristarchos/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func (a *AggregatorServiceImpl) CreateNewEntry(ctx context.Context, request *pb.AggregatorCreationRequest) (*pb.AggregatorCreationResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var requestId string
	if ok {
		headerValue := md.Get(service.HeaderKey)
		if len(headerValue) > 0 {
			requestId = headerValue[0]
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

		if traceCall {
			traceReceived := &pbar.TraceRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				Method:       "CreateNewEntry",
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			go a.Tracer.Trace(ctx, traceReceived)
		}

		logging.Trace(fmt.Sprintf("found traceId: %s", traceID))

		responseMd := metadata.New(map[string]string{service.HeaderKey: traceID})
		grpc.SendHeader(ctx, responseMd)
	}

	createNewWord := false
	query := a.Elastic.Builder().MatchQuery(ROOTWORD, request.RootWord)
	response, err := a.Elastic.Query().Match(a.Index, query)
	if err != nil {
		return &pb.AggregatorCreationResponse{Created: false, Updated: false}, err
	} else if len(response.Hits.Hits) == 0 {
		createNewWord = true
	}

	entry, err := a.mapAndHandleGrammaticalCategories(request)
	if err != nil {
		return nil, err
	}

	if createNewWord {
		entryAsJson, _ := json.Marshal(entry)
		createDocument, err := a.Elastic.Index().CreateDocument(a.Index, entryAsJson)
		if err != nil {
			return nil, err
		}

		logging.Debug(fmt.Sprintf("created document with id: %s and rootWord: %s", createDocument.ID, request.RootWord))
		return &pb.AggregatorCreationResponse{
			Created: true,
			Updated: false,
		}, nil
	}

	jsonHit, _ := json.Marshal(response.Hits.Hits[0].Source)
	rootWord, _ := UnmarshalRootWordEntry(jsonHit)

	updateNeeded := false
	conjucationFound := false
	for i, conjugation := range rootWord.Categories {
		if conjugation.Mood == entry.Categories[0].Mood && conjugation.Tense == entry.Categories[0].Tense && conjugation.Aspect == entry.Categories[0].Aspect {
			conjucationFound = true
			formFound := false
			for _, cform := range conjugation.Forms {
				if cform.Person == entry.Categories[0].Forms[0].Person && cform.Number == entry.Categories[0].Forms[0].Number && cform.Word == entry.Categories[0].Forms[0].Word {
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
	}

	if !conjucationFound {
		rootWord.Categories = append(rootWord.Categories, entry.Categories[0])
		updateNeeded = true
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
			return nil, err
		}

		logging.Debug(fmt.Sprintf("updated document with id: %s and rootWord: %s", createDocument.ID, request.RootWord))
		return &pb.AggregatorCreationResponse{
			Created: false,
			Updated: true,
		}, nil
	}

	logging.Debug("no action needed since document and grammar exists")

	return &pb.AggregatorCreationResponse{
		Created: false,
		Updated: false,
	}, nil
}

func (a *AggregatorServiceImpl) RetrieveEntry(ctx context.Context, request *pb.AggregatorRequest) (*pb.RootWordResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var requestId string
	if ok {
		headerValue := md.Get(service.HeaderKey)
		if len(headerValue) > 0 {
			requestId = headerValue[0]
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

		if traceCall {
			traceReceived := &pbar.TraceRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				Method:       "RetrieveEntry",
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			go a.Tracer.Trace(ctx, traceReceived)
		}

		logging.Trace(fmt.Sprintf("found traceId: %s", traceID))

		responseMd := metadata.New(map[string]string{service.HeaderKey: traceID})
		grpc.SendHeader(ctx, responseMd)
	}

	query := a.Elastic.Builder().MatchQuery(ROOTWORD, request.RootWord)
	response, err := a.Elastic.Query().Match(a.Index, query)
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
		conjPB := &pb.GrammaticalCategory{
			Tense:  conj.Tense,
			Mood:   conj.Mood,
			Aspect: conj.Aspect,
		}

		for _, form := range conj.Forms {
			formPB := &pb.GrammaticalForm{
				Person: form.Person,
				Number: form.Number,
				Gender: form.Gender,
				Case:   form.Case,
				Word:   form.Word,
			}
			conjPB.Forms = append(conjPB.Forms, formPB)
		}

		responsePB.Categories = append(responsePB.Categories, conjPB)
	}

	return &responsePB, nil
}

func (a *AggregatorServiceImpl) RetrieveSearchWords(ctx context.Context, request *pb.AggregatorRequest) (*pb.SearchWordResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var requestId string
	if ok {
		headerValue := md.Get(service.HeaderKey)
		if len(headerValue) > 0 {
			requestId = headerValue[0]
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

		if traceCall {
			traceReceived := &pbar.TraceRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				Method:       "RetrieveSearchWords",
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			go a.Tracer.Trace(ctx, traceReceived)
		}

		logging.Trace(fmt.Sprintf("found traceId: %s", traceID))
		responseMd := metadata.New(map[string]string{service.HeaderKey: traceID})
		grpc.SendHeader(ctx, responseMd)
	}

	query := a.Elastic.Builder().MatchQuery(ROOTWORD, request.RootWord)
	response, err := a.Elastic.Query().Match(a.Index, query)
	if err != nil {
		return nil, err
	} else if len(response.Hits.Hits) == 0 {
		return nil, fmt.Errorf("no entry can be found")
	}

	var responsePB pb.SearchWordResponse
	jsonHit, _ := json.Marshal(response.Hits.Hits[0].Source)
	rootWord, _ := UnmarshalRootWordEntry(jsonHit)

	for _, conj := range rootWord.Categories {
		for _, form := range conj.Forms {
			responsePB.Word = append(responsePB.Word, form.Word)
		}
	}

	return &responsePB, nil
}

func (a *AggregatorServiceImpl) mapAndHandleGrammaticalCategories(request *pb.AggregatorCreationRequest) (*RootWordEntry, error) {
	// example: 3th sing - impf - ind - act
	// example: 1st plur - aor - ind - act
	// example: noun - plural - masc - nom
	// example: pres act part - sing - masc - nom
	var entry RootWordEntry
	ruleSet := strings.Split(request.Rule, "-")

	entry.RootWord = request.RootWord
	entry.PartOfSpeech = mapEnumToCategory(request.PartOfSpeech)
	entry.Translations = []string{request.Translation}

	var conjForm GrammaticalForm
	var conj GrammaticalCategory

	switch request.PartOfSpeech {
	case pb.PartOfSpeech_VERB:
		if len(ruleSet) >= 4 {
			form := strings.Split(strings.TrimSpace(ruleSet[0]), " ")
			conjForm = GrammaticalForm{
				Person: strings.TrimSpace(form[0]),
				Number: strings.TrimSpace(form[1]),
				Word:   request.Word,
			}
			conj = GrammaticalCategory{
				Tense:  strings.TrimSpace(ruleSet[3]),
				Mood:   strings.TrimSpace(ruleSet[2]),
				Aspect: strings.TrimSpace(ruleSet[1]),
				Forms:  []GrammaticalForm{conjForm},
			}
		}
	case pb.PartOfSpeech_NOUN:
		if len(ruleSet) >= 4 {
			conjForm = GrammaticalForm{
				Number: strings.TrimSpace(ruleSet[1]),
				Gender: strings.TrimSpace(ruleSet[2]),
				Case:   strings.TrimSpace(ruleSet[3]),
				Word:   request.Word,
			}
			conj = GrammaticalCategory{
				Forms: []GrammaticalForm{conjForm},
			}
		}
	case pb.PartOfSpeech_PARTICIPLE:
		if len(ruleSet) >= 4 {
			form := strings.Split(strings.TrimSpace(ruleSet[0]), " ")
			conjForm = GrammaticalForm{
				Number: strings.TrimSpace(ruleSet[1]),
				Gender: strings.TrimSpace(ruleSet[2]),
				Case:   strings.TrimSpace(ruleSet[3]),
				Word:   request.Word,
			}
			conj = GrammaticalCategory{
				Tense:  strings.TrimSpace(form[1]),
				Aspect: strings.TrimSpace(form[0]),
				Forms:  []GrammaticalForm{conjForm},
			}
		}
	default:
		return nil, errors.New("unsupported grammatical category")
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
	default:
		return "UNKNOWN"
	}
}
