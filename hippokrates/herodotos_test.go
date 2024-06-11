package hippokrates

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/stretchr/testify/assert"
	"strconv"
)

const (
	OptionsContext string = "optionsContext"
	AnalyzeContext string = "analyzeContext"
	TextContext    string = "textContext"
	CheckedContext string = "checkedContext"
)

func (l *OdysseiaFixture) aListOfBooksAuthorsAndReferencesShouldBeReturned() error {
	aggregates := l.ctx.Value(OptionsContext).(models.AggregationResult)
	err := assertTrue(
		assert.True, len(aggregates.Authors) >= 1,
		"aggregates %v when more than 1 expected", len(aggregates.Authors),
	)
	if err != nil {
		return err
	}

	for _, author := range aggregates.Authors {
		err := assertTrue(
			assert.True, len(author.Books) >= 1,
			"authors %v when more than 1 expected", len(author.Books),
		)
		if err != nil {
			return err
		}
	}

	return nil

}

func (l *OdysseiaFixture) aQueryIsMadeForOptions() error {
	response, err := l.client.Herodotos().Options(TraceId)
	if err != nil {
		return err
	}

	var options models.AggregationResult
	err = json.NewDecoder(response.Body).Decode(&options)

	l.ctx = context.WithValue(l.ctx, OptionsContext, options)

	return err
}

func (l *OdysseiaFixture) aTheWordIsAnalyzed(word string) error {
	r := models.AnalyzeTextRequest{Rootword: word}
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return err
	}

	response, err := l.client.Herodotos().Analyze(jsonBody, TraceId)
	if err != nil {
		return err
	}

	var analyse models.AnalyzeTextResponse
	err = json.NewDecoder(response.Body).Decode(&analyse)

	l.ctx = context.WithValue(l.ctx, AnalyzeContext, analyse)

	return err
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

func (l *OdysseiaFixture) theResponseIsUsedToCreateANewText() error {
	analyse := l.ctx.Value(AnalyzeContext).(models.AnalyzeTextResponse)
	randomAnalyse := analyse.Results[l.randomizer.RandomNumberBaseZero(len(analyse.Results))]
	r := models.CreateTextRequest{
		Author:    randomAnalyse.Author,
		Book:      randomAnalyse.Book,
		Reference: randomAnalyse.Reference,
		Section:   randomAnalyse.Text.Section,
	}

	jsonBody, err := json.Marshal(r)
	if err != nil {
		return err
	}

	response, err := l.client.Herodotos().Create(jsonBody, TraceId)
	if err != nil {
		return err
	}

	var createdText models.Text
	err = json.NewDecoder(response.Body).Decode(&createdText)

	l.ctx = context.WithValue(l.ctx, TextContext, createdText)

	return nil
}

func (l *OdysseiaFixture) theSentenceIsCheckedAgainstTheOfficialTranslation() error {
	created := l.ctx.Value(TextContext).(models.Text)
	err := assertTrue(
		assert.True, len(created.Rhemai) >= 1,
		"created.Rhemai %v when more than 1 expected", len(created.Rhemai),
	)
	if err != nil {
		return err
	}

	rhema := created.Rhemai[0]
	translation := []models.Translations{{
		Section:     rhema.Section,
		Translation: rhema.Translations[0],
	},
	}

	r := models.CheckTextRequest{
		Translations: translation,
		Author:       created.Author,
		Book:         created.Book,
		Reference:    created.Reference,
	}

	jsonBody, err := json.Marshal(r)
	if err != nil {
		return err
	}

	response, err := l.client.Herodotos().Check(jsonBody, TraceId)
	if err != nil {
		return err
	}

	var checked models.CheckTextResponse
	err = json.NewDecoder(response.Body).Decode(&checked)

	l.ctx = context.WithValue(l.ctx, CheckedContext, checked)

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

func (l *OdysseiaFixture) theAverageLevenshteinShouldBeLessThanPerfect() error {
	checked := l.ctx.Value(CheckedContext).(models.CheckTextResponse)
	levhenstein, err := strconv.ParseFloat(checked.AverageLevenshteinPercentage, 32)
	if err != nil {
		return err
	}

	if levhenstein == 100.0 {
		return fmt.Errorf("expected %v to be less than 100", levhenstein)
	}

	return nil
}

func (l *OdysseiaFixture) theResponseShouldIncludePossibleTypos() error {
	checked := l.ctx.Value(CheckedContext).(models.CheckTextResponse)
	err := assertTrue(
		assert.True, len(checked.PossibleTypos) >= 1,
		"checked.PossibleTypos %v when more than 1 expected", len(checked.PossibleTypos),
	)
	if err != nil {
		return err
	}

	return nil
}

func (l *OdysseiaFixture) theTextWithAuthorAndBookAndReferenceAndSectionIsCheckedWithTypos(author, book, reference, section string) error {
	translation := []models.Translations{{
		Section:     section,
		Translation: "this is the display of the inquiry of Herodotws of Halikarnassws",
	},
	}
	r := models.CheckTextRequest{
		Translations: translation,
		Author:       author,
		Book:         book,
		Reference:    reference,
	}

	jsonBody, err := json.Marshal(r)
	if err != nil {
		return err
	}

	response, err := l.client.Herodotos().Check(jsonBody, TraceId)
	if err != nil {
		return err
	}

	var checked models.CheckTextResponse
	err = json.NewDecoder(response.Body).Decode(&checked)

	l.ctx = context.WithValue(l.ctx, CheckedContext, checked)

	return nil
}
