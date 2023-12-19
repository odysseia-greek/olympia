package scholar

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	pb "github.com/odysseia-greek/olympia/aristarchos/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
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
	var traceID string
	if ok {
		headerValue := md.Get(service.HeaderKey)
		if len(headerValue) > 0 {
			traceID = headerValue[0]
		}

		logging.Trace(fmt.Sprintf("found traceId: %s", traceID))
	}

	createNewWord := false
	query := a.Elastic.Builder().MatchQuery(ROOTWORD, request.RootWord)
	response, err := a.Elastic.Query().Match(a.Index, query)
	if err != nil {
		return &pb.AggregatorCreationResponse{Created: false, Updated: false}, err
	} else if len(response.Hits.Hits) == 0 {
		createNewWord = true
	}

	// example: 3th sing - impf - ind - act
	ruleSet := strings.Split(request.Rule, "-")
	form := strings.Split(ruleSet[0], " ")
	var entry RootWordEntry
	conjForm := ConjugationForm{
		Person: strings.TrimSpace(form[0]),
		Number: strings.TrimSpace(form[1]),
		Word:   request.Word,
	}
	conj := Conjugation{
		Tense:  strings.TrimSpace(ruleSet[len(ruleSet)-1]),
		Mood:   strings.TrimSpace(ruleSet[len(ruleSet)-2]),
		Aspect: strings.TrimSpace(ruleSet[len(ruleSet)-3]),
		Forms: []ConjugationForm{
			conjForm,
		},
	}
	entry = RootWordEntry{
		RootWord:     request.RootWord,
		Translations: []string{request.Translation},
		Conjugations: []Conjugation{
			conj,
		},
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
	for i, conjugation := range rootWord.Conjugations {
		if conjugation.Mood == conj.Mood && conjugation.Tense == conj.Tense && conjugation.Aspect == conj.Aspect {
			conjucationFound = true
			formFound := false
			for _, cform := range conjugation.Forms {
				if cform.Person == conjForm.Person && cform.Number == conjForm.Number && cform.Word == conjForm.Word {
					formFound = true
					break
				}
			}
			if !formFound {
				rootWord.Conjugations[i].Forms = append(rootWord.Conjugations[i].Forms, conjForm)
				updateNeeded = true
			}
			break
		}
	}

	if !conjucationFound {
		rootWord.Conjugations = append(rootWord.Conjugations, conj)
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

	responseMd := metadata.New(map[string]string{service.HeaderKey: traceID})
	grpc.SendHeader(ctx, responseMd)
	return &pb.AggregatorCreationResponse{
		Created: false,
		Updated: false,
	}, nil
}

func (a *AggregatorServiceImpl) RetrieveEntry(ctx context.Context, request *pb.AggregatorRequest) (*pb.RootWordResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var traceID string
	if ok {
		headerValue := md.Get(service.HeaderKey)
		if len(headerValue) > 0 {
			traceID = headerValue[0]
		}

		logging.Trace(fmt.Sprintf("found traceId: %s", traceID))
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
	for _, conj := range rootWord.Conjugations {
		conjPB := &pb.Conjugation{
			Tense:  conj.Tense,
			Mood:   conj.Mood,
			Aspect: conj.Aspect,
		}

		for _, form := range conj.Forms {
			formPB := &pb.ConjugationForm{
				Number: form.Number,
				Person: form.Person,
				Word:   form.Word,
			}
			conjPB.Forms = append(conjPB.Forms, formPB)
		}

		responsePB.Conjugations = append(responsePB.Conjugations, conjPB)
	}

	responseMd := metadata.New(map[string]string{service.HeaderKey: traceID})
	grpc.SendHeader(ctx, responseMd)
	return &responsePB, nil
}

func (a *AggregatorServiceImpl) RetrieveSearchWords(ctx context.Context, request *pb.AggregatorRequest) (*pb.SearchWordResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var traceID string
	if ok {
		headerValue := md.Get(service.HeaderKey)
		if len(headerValue) > 0 {
			traceID = headerValue[0]
		}

		logging.Trace(fmt.Sprintf("found traceId: %s", traceID))
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

	for _, conj := range rootWord.Conjugations {
		for _, form := range conj.Forms {
			responsePB.Word = append(responsePB.Word, form.Word)
		}
	}

	responseMd := metadata.New(map[string]string{service.HeaderKey: traceID})
	grpc.SendHeader(ctx, responseMd)
	return &responsePB, nil
}
