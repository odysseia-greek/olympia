package hippokrates

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/stretchr/testify/assert"
)

const (
	AnalyzeContext    string = "analyzeContext"
	TextContext       string = "textContext"
	CheckedContext    string = "checkedContext"
	DictionaryContext string = "dictionaryContext"
	MediaContext      string = "mediaContext"
)

type dictionarySmokeResponse struct {
	Results []struct {
		Headword string `json:"headword"`
	} `json:"results"`
	FoundInText *struct {
		Texts []struct {
			Reference string `json:"reference"`
		} `json:"texts"`
	} `json:"foundInText"`
}

type mediaSmokeContext struct {
	Theme    string
	Segment  string
	Set      string
	QuizWord string
	Answer   string
}

type graphQLError struct {
	Message string `json:"message"`
}

func (l *OdysseiaFixture) theGatewayIsUp() error {
	// Send the GraphQL query as an HTTP GET request
	response, err := http.Get(l.homeros.health)
	if err != nil {
		return err
	}

	var health models.Health
	err = json.NewDecoder(response.Body).Decode(&health)

	if !health.Healthy {
		return fmt.Errorf("service was %v were a healthy status was expected", health.Healthy)
	}

	return nil
}

func (l *OdysseiaFixture) theGrammarIsCheckedForWordThroughTheGateway(word string) error {
	// Define your GraphQL query
	query := fmt.Sprintf(`query {
	grammar(word: "%s") {
		results {
			translations
			word
			rule
			rootWord
		}
	}
}`, word)
	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var dionysiosResponse struct {
		Data struct {
			Grammar struct {
				Results []models.Result `json:"results"`
			} `json:"grammar"`
		} `json:"data"`
	}

	err = json.NewDecoder(response.Body).Decode(&dionysiosResponse)

	l.ctx = context.WithValue(l.ctx, ResponseBody, dionysiosResponse.Data.Grammar.Results)

	return nil
}

func (l *OdysseiaFixture) aFoundInTextResponseShouldIncludeResults() error {
	results := l.ctx.Value(ResponseBody).([]models.AnalyzeResult)
	err := assertTrue(
		assert.True, len(results) >= 1,
		"results %v when more than 1 expected", len(results),
	)
	if err != nil {
		return err
	}

	return nil

}

func (l *OdysseiaFixture) theDeclensionShouldBeIncludedInTheResponseAsAGatewayStruct(declension string) error {
	declensions := l.ctx.Value(ResponseBody).([]models.Result)

	found := false
	for _, decResult := range declensions {
		if decResult.Rule == declension {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("could not find declension %v in slice", declension)
	}
	return nil
}

func (l *OdysseiaFixture) theGatewayShouldRespondWithACorrectness() error {
	_, ok := l.ctx.Value(QuizContext).(bool)
	if !ok {
		return fmt.Errorf("the media answer did not include a correctness value")
	}

	return nil
}

func (l *OdysseiaFixture) theDictionaryEndpointIsSearchedForIn(endpoint, word, language string) error {
	return l.searchDictionary(endpoint, word, language, false)
}

func (l *OdysseiaFixture) theDictionaryEndpointIsSearchedForInWithTextExpansion(endpoint, word, language string) error {
	return l.searchDictionary(endpoint, word, language, true)
}

func (l *OdysseiaFixture) searchDictionary(endpoint, word, language string, expand bool) error {
	languageValues := map[string]string{
		"Greek":   "LANG_GREEK",
		"English": "LANG_ENGLISH",
		"Dutch":   "LANG_DUTCH",
	}
	languageValue, ok := languageValues[language]
	if !ok {
		return fmt.Errorf("unsupported dictionary language %q", language)
	}

	var inputType, selection string
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"word": word, "language": languageValue, "page": 1, "size": 10,
		},
	}
	switch endpoint {
	case "exact", "text":
		inputType = "ExpandableSearchQueryInput!"
		variables["input"].(map[string]interface{})["expand"] = expand
		selection = `results { headword } foundInText { texts { reference } }`
	case "fuzzy", "partial", "phrase":
		if expand {
			return fmt.Errorf("dictionary endpoint %q does not support text expansion", endpoint)
		}
		inputType = "SearchQueryInput!"
		selection = `results { headword }`
	default:
		return fmt.Errorf("unsupported dictionary endpoint %q", endpoint)
	}

	operation := strings.ToUpper(endpoint[:1]) + endpoint[1:]
	query := fmt.Sprintf(`query %s($input: %s) { %s(input: $input) { %s } }`, operation, inputType, endpoint, selection)
	response, err := l.graphqlRequest(operation, query, variables)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var result struct {
		Data   map[string]dictionarySmokeResponse `json:"data"`
		Errors []graphQLError                     `json:"errors"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return err
	}
	if err := graphqlErrors(result.Errors); err != nil {
		return err
	}

	l.ctx = context.WithValue(l.ctx, DictionaryContext, result.Data[endpoint])
	return nil
}

func (l *OdysseiaFixture) theDictionaryResponseShouldContainResults() error {
	response, ok := l.ctx.Value(DictionaryContext).(dictionarySmokeResponse)
	if !ok {
		return fmt.Errorf("the dictionary response was not available")
	}
	if len(response.Results) == 0 {
		return fmt.Errorf("the dictionary response contained no results")
	}
	return nil
}

func (l *OdysseiaFixture) theDictionaryResponseShouldContainTextOccurrences() error {
	response, ok := l.ctx.Value(DictionaryContext).(dictionarySmokeResponse)
	if !ok || response.FoundInText == nil {
		return fmt.Errorf("the dictionary response did not include foundInText")
	}
	if len(response.FoundInText.Texts) == 0 {
		return fmt.Errorf("the dictionary response contained no text occurrences")
	}
	return nil
}

func (l *OdysseiaFixture) theMediaQuizOptionsAreRequested() error {
	query := `query MediaOptions { mediaOptions { themes { name segments { name maxSet } } } }`
	response, err := l.graphqlRequest("MediaOptions", query, nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var result struct {
		Data struct {
			MediaOptions struct {
				Themes []struct {
					Name     string `json:"name"`
					Segments []struct {
						Name   string  `json:"name"`
						MaxSet float64 `json:"maxSet"`
					} `json:"segments"`
				} `json:"themes"`
			} `json:"mediaOptions"`
		} `json:"data"`
		Errors []graphQLError `json:"errors"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return err
	}
	if err := graphqlErrors(result.Errors); err != nil {
		return err
	}
	if len(result.Data.MediaOptions.Themes) == 0 || len(result.Data.MediaOptions.Themes[0].Segments) == 0 {
		return fmt.Errorf("mediaOptions contained no usable theme and segment")
	}

	theme := result.Data.MediaOptions.Themes[0]
	segment := theme.Segments[0]
	l.ctx = context.WithValue(l.ctx, MediaContext, mediaSmokeContext{
		Theme: theme.Name, Segment: segment.Name, Set: "1",
	})
	return nil
}

func (l *OdysseiaFixture) aMediaQuizIsCreatedFromThoseOptions() error {
	media, ok := l.ctx.Value(MediaContext).(mediaSmokeContext)
	if !ok {
		return fmt.Errorf("media quiz options were not available")
	}
	query := `query MediaQuiz($input: MediaQuizInput!) {
		mediaQuiz(input: $input) { quizItem options { option } }
	}`
	variables := map[string]interface{}{"input": map[string]interface{}{
		"theme": media.Theme, "segment": media.Segment, "set": media.Set,
		"doneAfter": 1, "resetProgress": false, "archiveProgress": false,
	}}
	response, err := l.graphqlRequest("MediaQuiz", query, variables)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	var result struct {
		Data struct {
			MediaQuiz struct {
				QuizItem string `json:"quizItem"`
				Options  []struct {
					Option string `json:"option"`
				} `json:"options"`
			} `json:"mediaQuiz"`
		} `json:"data"`
		Errors []graphQLError `json:"errors"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return err
	}
	if err := graphqlErrors(result.Errors); err != nil {
		return err
	}
	if result.Data.MediaQuiz.QuizItem == "" || len(result.Data.MediaQuiz.Options) == 0 {
		return fmt.Errorf("mediaQuiz did not return a quiz item and answer options")
	}
	media.QuizWord = result.Data.MediaQuiz.QuizItem
	media.Answer = result.Data.MediaQuiz.Options[0].Option
	l.ctx = context.WithValue(l.ctx, MediaContext, media)
	return nil
}

func (l *OdysseiaFixture) theFirstMediaQuizOptionIsAnswered() error {
	media, ok := l.ctx.Value(MediaContext).(mediaSmokeContext)
	if !ok || media.QuizWord == "" || media.Answer == "" {
		return fmt.Errorf("a media quiz question was not available")
	}
	query := `query MediaAnswer($input: MediaAnswerInput!) {
		mediaAnswer(input: $input) { correct }
	}`
	variables := map[string]interface{}{"input": map[string]interface{}{
		"theme": media.Theme, "segment": media.Segment, "set": media.Set,
		"quizWord": media.QuizWord, "answer": media.Answer,
		"comprehensive": false, "doneAfter": 1,
	}}
	response, err := l.graphqlRequest("MediaAnswer", query, variables)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	var result struct {
		Data struct {
			MediaAnswer *struct {
				Correct *bool `json:"correct"`
			} `json:"mediaAnswer"`
		} `json:"data"`
		Errors []graphQLError `json:"errors"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return err
	}
	if err := graphqlErrors(result.Errors); err != nil {
		return err
	}
	if result.Data.MediaAnswer == nil || result.Data.MediaAnswer.Correct == nil {
		return fmt.Errorf("mediaAnswer did not return correct")
	}
	l.ctx = context.WithValue(l.ctx, QuizContext, *result.Data.MediaAnswer.Correct)
	return nil
}

func (l *OdysseiaFixture) aQueryIsMadeForAllTextOptions() error {
	query := `query {
	textOptions {
		authors {
			key
			books {
				key
				references{
					key
					sections{
						key
					}
				}
			}
		}
}
}`

	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var optionsResponse struct {
		Data struct {
			Response models.AggregationResult `json:"textOptions"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&optionsResponse)

	l.ctx = context.WithValue(l.ctx, AggregateContext, optionsResponse.Data.Response)

	return nil
}

func (l *OdysseiaFixture) thatResponseIsUsedToCreateANewText() error {
	aggregation, ok := l.ctx.Value(AggregateContext).(models.AggregationResult)
	if !ok {
		return fmt.Errorf("failed to assert AggregationResult type")
	}

	randomAuthor := aggregation.Authors[l.randomizer.RandomNumberBaseZero(len(aggregation.Authors))]
	randomBook := randomAuthor.Books[l.randomizer.RandomNumberBaseZero(len(randomAuthor.Books))]
	randomRef := randomBook.References[l.randomizer.RandomNumberBaseZero(len(randomBook.References))]
	randomSection := randomRef.Sections[l.randomizer.RandomNumberBaseZero(len(randomRef.Sections))]

	query := fmt.Sprintf(`query {
  create(input: {author: "%s", book: "%s", reference: "%s", section: "%s"}) {
    author
    book
    reference
    rhemai {
      greek
      section
      translations
    }
  }
}
`, randomAuthor.Key, randomBook.Key, randomRef.Key, randomSection.Key)

	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var textResponse struct {
		Data struct {
			Response models.Text `json:"create"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&textResponse)

	l.ctx = context.WithValue(l.ctx, TextContext, textResponse.Data.Response)

	return nil
}

func (l *OdysseiaFixture) theTextIsCheckedAgainstTheOfficialTranslation() error {
	text, ok := l.ctx.Value(TextContext).(models.Text)
	if !ok {
		return fmt.Errorf("failed to assert Text type")
	}

	query := fmt.Sprintf(`query {
  check(input: {author: "%s", 
		book: "%s",
		reference: "%s",
	translations: 
	[
		{
	section: "%s",
	translation: "%s"
		}
	]
}) {
    averageLevenshteinPercentage
	sections{
		section
		answerSentence
		levenshteinPercentage
	}
	possibleTypos{
		source
		provided
		
	}
  }
}
`, text.Author, text.Book, text.Reference, text.Rhemai[0].Section, text.Rhemai[0].Translations[0])
	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var checkResponse struct {
		Data struct {
			Response models.CheckTextResponse `json:"check"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&checkResponse)

	l.ctx = context.WithValue(l.ctx, CheckedContext, checkResponse.Data.Response)

	return nil
}

func (l *OdysseiaFixture) theWordIsAnalyzedThroughTheGateway(word string) error {
	query := fmt.Sprintf(`query {
	analyze(rootword: "%s") {
		rootword
		conjugations{
			word
			rule
		}
		texts{
			text{
				greek
				section
			}
			referenceLink
			author
			book
			reference
		}
}
}
`, word)

	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var analyzeResponse struct {
		Data struct {
			Response struct {
				Conjugations []models.Conjugations  `json:"conjugations"`
				Results      []models.AnalyzeResult `json:"texts"`
				Rootword     string                 `json:"rootword"`
			} `json:"analyze"`
		} `json:"data"`
	}

	err = json.NewDecoder(response.Body).Decode(&analyzeResponse)

	analyzeContextResponse := models.AnalyzeTextResponse{
		Rootword:     analyzeResponse.Data.Response.Rootword,
		PartOfSpeech: "",
		Conjugations: analyzeResponse.Data.Response.Conjugations,
		Results:      analyzeResponse.Data.Response.Results,
	}

	l.ctx = context.WithValue(l.ctx, AnalyzeContext, analyzeContextResponse)

	return nil
}

func (l *OdysseiaFixture) theResponseHasACompleteAnalyzesIncluded() error {
	analyse := l.ctx.Value(AnalyzeContext).(models.AnalyzeTextResponse)
	err := assertTrue(
		assert.True, len(analyse.Results) >= 1,
		"analyse.Results %v when more than 1 expected", len(analyse.Results),
	)
	if err != nil {
		return err
	}

	err = assertTrue(
		assert.True, len(analyse.Conjugations) >= 1,
		"analyse.Conjugations %v when more than 1 expected", len(analyse.Conjugations),
	)
	if err != nil {
		return err
	}

	return nil
}

func (l *OdysseiaFixture) theAverageLevenshteinShouldBePerfect() error {
	checked := l.ctx.Value(CheckedContext).(models.CheckTextResponse)
	levhenstein, err := strconv.ParseFloat(checked.AverageLevenshteinPercentage, 32)
	if err != nil {
		return err
	}

	if levhenstein != 100.0 {
		return fmt.Errorf("expected %v to be 100", levhenstein)
	}

	return nil
}

func (l *OdysseiaFixture) anErrorContainingIsReturned(message string) error {
	errorText := l.ctx.Value(ErrorBody).(string)
	if !strings.Contains(errorText, message) {
		return fmt.Errorf("expected %v to contain %v", errorText, message)
	}

	return nil
}

func (l *OdysseiaFixture) graphqlHelper(query string) (*http.Response, error) {
	return l.graphqlRequest("", query, nil)
}

func (l *OdysseiaFixture) graphqlRequest(operationName, query string, variables map[string]interface{}) (*http.Response, error) {
	if variables == nil {
		variables = map[string]interface{}{}
	}
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	if operationName != "" {
		requestBody["operationName"] = operationName
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, l.homeros.graphql, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("boule", uuid.NewString())
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("Homeros returned HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}
	return response, nil
}

func graphqlErrors(errors []graphQLError) error {
	if len(errors) == 0 {
		return nil
	}
	messages := make([]string, 0, len(errors))
	for _, graphErr := range errors {
		messages = append(messages, graphErr.Message)
	}
	return fmt.Errorf("Homeros returned GraphQL errors: %s", strings.Join(messages, "; "))
}
