package hippokrates

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/stretchr/testify/assert"
	"math/rand/v2"
	"net/http"
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

func (l *OdysseiaFixture) aQuizIsPlayedInComprehensiveModeForTheWordAndTheCorrectAnswerWithTypeSetAndTheme(fullWord, grammarResponse, quizType, set, theme string) error {
	answerRequest := models.AnswerRequest{
		Theme:         theme,
		Set:           set,
		QuizType:      quizType,
		Comprehensive: true,
		Answer:        grammarResponse,
		Dialogue:      nil,
		QuizWord:      fullWord,
	}

	body, err := json.Marshal(answerRequest)
	if err != nil {
		return err
	}
	answer, err := l.client.Sokrates().Check(body, TraceId)
	if err != nil {
		return err
	}

	if answer.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 but got %v", answer.StatusCode)
	}

	var answerResponse models.ComprehensiveResponse
	err = json.NewDecoder(answer.Body).Decode(&answerResponse)
	if err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, QuizContext, answerResponse)

	return nil
}

func (l *OdysseiaFixture) theQuizresponseIsExpandedWithTextAndSimilarWords() error {
	answer, ok := l.ctx.Value(QuizContext).(models.ComprehensiveResponse)
	if !ok {
		return fmt.Errorf("failed to assert quizResponse type")
	}

	err := assertTrue(
		assert.True, len(answer.SimilarWords) >= 1,
		"expected length of similar word the be one or greater but was: %v", len(answer.SimilarWords),
	)
	if err != nil {
		return err
	}

	err = assertTrue(
		assert.True, len(answer.FoundInText.Results) >= 1,
		"expected length of rhemai to be one or greater but was: %v", len(answer.FoundInText.Results),
	)

	return err
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
