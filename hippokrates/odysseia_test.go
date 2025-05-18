package hippokrates

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
	"math/rand/v2"
	"strings"
)

var randomWords = []string{"τρέφοντος", "Ἡλίου", "δημοκρατίαν", "λοιπὰ", "ὀλιγαρχίας", "μεταβάλλει", "λόγιοι", "Βούλει", "ἔβαλλε", "λύουσι", "Λακεδαιμονίους", "Ἀττικῆς", "στρατιὰ", "Πελοποννησίους", "ἐστί", "ἀγορεύεσθαι"}

func (l *OdysseiaFixture) aGrammarEntryIsMadeForTheWord(word string) error {
	response, err := l.client.Dionysios().Grammar(word, TraceId)
	if err != nil {
		return err
	}

	var declensions models.DeclensionTranslationResults
	err = json.NewDecoder(response.Body).Decode(&declensions)

	l.ctx = context.WithValue(l.ctx, ResponseBody, declensions)

	return nil
}

func (l *OdysseiaFixture) theOptionsReturnedFromTheGrammarApiShouldInclude(grammarResponse string) error {
	declensions := l.ctx.Value(ResponseBody).(models.DeclensionTranslationResults)

	found := false
	for _, decResult := range declensions.Results {
		for _, translation := range decResult.Translation {
			if strings.Contains(translation, grammarResponse) || translation == grammarResponse {
				found = true
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("could not find declension %v in slice", grammarResponse)
	}
	return nil
}

func (l *OdysseiaFixture) theGrammarIsCheckedForARandomWordInTheList() error {
	n := rand.N(len(randomWords))
	randomWord := randomWords[n]

	query := fmt.Sprintf(`query grammar {
	grammar(word: "%s") {
		translation
		word
		rule
		rootWord
	}
}`, randomWord)
	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var dionysiosResponse struct {
		Data struct {
			Response []models.Result `json:"grammar"`
		} `json:"data"`
	}

	err = json.NewDecoder(response.Body).Decode(&dionysiosResponse)
	if err != nil {
		return err
	}
	l.ctx = context.WithValue(l.ctx, ResponseBody, dionysiosResponse.Data.Response)

	return nil
}

func (l *OdysseiaFixture) aResponseWithARootwordIsReturned() error {
	declensions := l.ctx.Value(ResponseBody).([]models.Result)

	if len(declensions) == 0 {
		return fmt.Errorf("the number of results from dionysios was 0 were 1 or more was expected")
	}
	var rootWord string
	for _, declension := range declensions {
		if declension.RootWord != "" {
			rootWord = declension.RootWord
			break
		}
	}

	l.ctx = context.WithValue(l.ctx, Rootword, rootWord)
	return nil
}

func (l *OdysseiaFixture) thatRootwordIsQueriedInAlexandrosWith(expand string) error {
	rootWord := l.ctx.Value(Rootword).(string)

	query := fmt.Sprintf(`query dictionary{
	dictionary(word: "%s", language: "greek", mode: "exact", searchInText: %v) {
			hits{
				hit{
					english
					greek
				}
				foundInText{
					rhemai{
						author
						greek
						translations
					}
				}
			}
		}
	}

`, rootWord, expand)
	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var alexandrosResponse struct {
		Data struct {
			Response models.ExtendedResponse `json:"dictionary"`
		} `json:"data"`
	}

	err = json.NewDecoder(response.Body).Decode(&alexandrosResponse)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, ResponseBody, alexandrosResponse.Data.Response)

	return nil
}

func (l *OdysseiaFixture) theQueryResultHasTextsIncluded() error {
	response := l.ctx.Value(ResponseBody).(models.ExtendedResponse)

	for _, hit := range response.Hits {
		if hit.FoundInText == nil {
			return fmt.Errorf("expected each response to have text but found none")
		}
	}

	return nil
}

func (l *OdysseiaFixture) aResultShouldBeReturned() error {
	declensions := l.ctx.Value(ResponseBody).(models.DeclensionTranslationResults)

	if len(declensions.Results) == 0 || declensions.Results == nil {
		return fmt.Errorf("expected declensions to have results but it did not")
	}

	return nil
}

func (l *OdysseiaFixture) thatWordIsSearchedForInTheGrammarComponent() error {
	grammarWord := l.ctx.Value(GrammarContext).(string)
	response, err := l.client.Dionysios().Grammar(grammarWord, TraceId)
	if err != nil {
		return err
	}

	var declensions models.DeclensionTranslationResults
	err = json.NewDecoder(response.Body).Decode(&declensions)

	l.ctx = context.WithValue(l.ctx, ResponseBody, declensions)

	return nil
}
