package hippokrates

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/stretchr/testify/assert"
	"net/http"
)

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
	query := fmt.Sprintf(`query grammar {
	grammar(word: "%s") {
		translation
		word
		rule
		rootWord
	}
}`, word)
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

	l.ctx = context.WithValue(l.ctx, ResponseBody, dionysiosResponse.Data.Response)

	return nil
}

func (l *OdysseiaFixture) theWordIsQueriedUsingAndThroughTheGateway(word, mode, language string) error {
	// Define your GraphQL query
	query := fmt.Sprintf(`query dictionary {
	dictionary(word: "%s", language: "%s", mode: "%s", searchInText: false) {
			hits{
				hit{
					english
					greek
					dutch
					linkedWord
					original
				}
		}
	}
}
`, word, language, mode)
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

	l.ctx = context.WithValue(l.ctx, ResponseBody, alexandrosResponse.Data.Response)

	return nil
}

func (l *OdysseiaFixture) theWordIsQueriedUsingAndAndSearchInTextThroughTheGateway(word, mode, language string) error {
	// Define your GraphQL query
	query := fmt.Sprintf(`query dictionary {
	dictionary(word: "%s", language: "%s", mode: "%s", searchInText: true) {
			hits{
				hit{
					english
					greek
					dutch
					linkedWord
					original
				}
		foundInText{
					rootword
					conjugations {
						word
						rule
					}
					results{
						author
						book
						reference
						referenceLink
						text{
							translations
							greek
						}
					}
		}
		}
	}
}
`, word, language, mode)
	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var alexandrosResponse struct {
		Data struct {
			Response struct {
				Hits []struct {
					FoundInText struct {
						Conjugations []models.Conjugations  `json:"conjugations"`
						Results      []models.AnalyzeResult `json:"results"`
						Rootword     string                 `json:"rootword"`
					}
					Hit models.Meros `json:"hit"`
				} `json:"hits"`
			} `json:"dictionary"`
		} `json:"data"`
	}

	err = json.NewDecoder(response.Body).Decode(&alexandrosResponse)

	l.ctx = context.WithValue(l.ctx, ResponseBody, alexandrosResponse.Data.Response.Hits[0].FoundInText.Results)

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
	_, ok := l.ctx.Value(QuizContext).(models.ComprehensiveResponse)
	if !ok {
		return fmt.Errorf("the answer was not a correct format")
	}

	return nil
}

func (l *OdysseiaFixture) aQueryIsMadeForAllTextOptions() error {
	query := `query textOptions {
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

	query := fmt.Sprintf(`query create {
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

	query := fmt.Sprintf(`query check {
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
	query := fmt.Sprintf(`query analyze {
	analyze(rootword: "%s") {
		rootword
		conjugations{
			word
			rule
		}
		results{
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
				Results      []models.AnalyzeResult `json:"results"`
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

func (l *OdysseiaFixture) graphqlHelper(query string) (*http.Response, error) {
	// Define any required variables (if applicable)
	variables := map[string]interface{}{
		// If your query requires variables, define them here
	}

	// Create the JSON payload for the GraphQL query
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	// Convert the payload to JSON
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil
	}

	// Send the GraphQL query as an HTTP POST request
	return http.Post(l.homeros.graphql, "application/json", bytes.NewBuffer(requestBodyBytes))
}
