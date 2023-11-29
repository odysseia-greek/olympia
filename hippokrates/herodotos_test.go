package hippokrates

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
	"golang.org/x/exp/utf8string"
	"strconv"
	"strings"
)

func (l *OdysseiaFixture) aQueryIsMadeForAllAuthors() error {
	response, err := l.client.Herodotos().GetAuthors(TraceId)
	if err != nil {
		return err
	}

	var authors models.Authors
	err = json.NewDecoder(response.Body).Decode(&authors)

	l.ctx = context.WithValue(l.ctx, ResponseBody, authors)

	return nil
}

func (l *OdysseiaFixture) aQueryIsMadeForAllBooksByAuthor(author string) error {
	response, err := l.client.Herodotos().GetBooks(author, TraceId)
	if err != nil {
		return err
	}

	var books models.Books
	err = json.NewDecoder(response.Body).Decode(&books)

	l.ctx = context.WithValue(l.ctx, ResponseBody, books)

	return nil
}

func (l *OdysseiaFixture) anAuthorAndBookCombinationIsQueried() error {
	response, err := l.client.Herodotos().GetAuthors(TraceId)
	if err != nil {
		return err
	}

	var authors models.Authors
	err = json.NewDecoder(response.Body).Decode(&authors)

	randNumber := GenerateRandomNumber(len(authors.Authors))
	author := authors.Authors[randNumber].Author

	resp, err := l.client.Herodotos().GetBooks(author, TraceId)
	if err != nil {
		return err
	}

	var books models.Books
	err = json.NewDecoder(resp.Body).Decode(&books)

	randomBookNumber := GenerateRandomNumber(len(books.Books))
	book := strconv.Itoa(int(books.Books[randomBookNumber].Book))

	question, err := l.client.Herodotos().CreateQuestion(author, book, TraceId)
	if err != nil {
		return err
	}

	var query models.CreateSentenceResponse
	err = json.NewDecoder(question.Body).Decode(&query)

	l.ctx = context.WithValue(l.ctx, ContextAuthor, author)
	l.ctx = context.WithValue(l.ctx, ResponseBody, query)

	return nil
}

func (l *OdysseiaFixture) aTranslationIsReturned() error {
	author := l.ctx.Value(ContextAuthor).(string)
	sentence := l.ctx.Value(ResponseBody).(models.CreateSentenceResponse)

	translation := "this is a random translation that should be long enough to have some matches"

	answerModel := models.CheckSentenceRequest{
		SentenceId:       sentence.SentenceId,
		ProvidedSentence: translation,
		Author:           author,
	}

	answerResponse, err := l.client.Herodotos().CheckSentence(answerModel, TraceId)
	if err != nil {
		return err
	}

	var answer models.CheckSentenceResponse
	err = json.NewDecoder(answerResponse.Body).Decode(&answer)

	l.ctx = context.WithValue(l.ctx, AnswerBody, answer)

	return nil
}

func (l *OdysseiaFixture) theAuthorShouldBeIncluded(author string) error {
	authors := l.ctx.Value(ResponseBody).(models.Authors)

	found := false

	for _, result := range authors.Authors {
		if strings.Contains(result.Author, author) {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("could not find author %v in slice", author)
	}
	return nil
}

func (l *OdysseiaFixture) theNumberOfAuthorsShouldExceed(results int) error {
	authors := l.ctx.Value(ResponseBody).(models.Authors)
	numberOfAuthors := len(authors.Authors)
	if results > numberOfAuthors {
		return fmt.Errorf("expected results to be equal to or more than %v but was %v", results, numberOfAuthors)
	}

	return nil
}

func (l *OdysseiaFixture) theBookShouldBeIncluded(book int) error {
	books := l.ctx.Value(ResponseBody).(models.Books)

	found := false

	for _, result := range books.Books {
		if result.Book == int64(book) {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("could not find book %v in slice", book)
	}
	return nil
}

func (l *OdysseiaFixture) theSentenceIdShouldBeLongerThan(lengthOfId int) error {
	sentence := l.ctx.Value(ResponseBody).(models.CreateSentenceResponse)
	idLen := len(sentence.SentenceId)

	if idLen < lengthOfId {
		return fmt.Errorf("expected id to be longer than %v but was %v - id was %v", lengthOfId, idLen, sentence.SentenceId)
	}

	return nil
}

func (l *OdysseiaFixture) theSentenceShouldIncludeNonASCIIGreekCharacters() error {
	sentence := l.ctx.Value(ResponseBody).(models.CreateSentenceResponse)

	ascii := utf8string.NewString(sentence.Sentence).IsASCII()

	if ascii {
		return fmt.Errorf("expected sentence to not include ASCII - %v", sentence.Sentence)
	}
	return nil
}

func (l *OdysseiaFixture) aCorrectnessPercentage() error {
	answer := l.ctx.Value(AnswerBody).(models.CheckSentenceResponse)

	levenshtein, err := strconv.ParseFloat(answer.LevenshteinPercentage, 32)
	if err != nil {
		return err
	}

	if levenshtein < 0.1 {
		return fmt.Errorf("expected levenshtein to be greater than zero but was %v", answer.LevenshteinPercentage)
	}

	return nil
}

func (l *OdysseiaFixture) aSentenceWithATranslation() error {
	answer := l.ctx.Value(AnswerBody).(models.CheckSentenceResponse)
	if answer.AnswerSentence == "" {
		return fmt.Errorf("expected an answer to have been provided but found none %v", answer.AnswerSentence)
	}

	return nil
}
