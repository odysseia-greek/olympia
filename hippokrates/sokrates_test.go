package hippokrates

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"strings"
)

const (
	AnswerContext    string = "answerContext"
	QuizContext      string = "quizContext"
	BodyContext      string = "bodyContext"
	AggregateContext string = "aggregateContext"
	TraceId          string = "hippokrates-traceid"
)

func (l *OdysseiaFixture) aListOfThemesWithTheHighestSetShouldBeReturned() error {
	aggregates := l.ctx.Value(AggregateContext).(models.AggregatedOptions)
	err := assertTrue(
		assert.True, len(aggregates.Themes) >= 1,
		"aggregates %v when more than 1 expected", len(aggregates.Themes),
	)
	if err != nil {
		return err
	}

	for _, aggregate := range aggregates.Themes {
		err := assertTrue(
			assert.True, len(aggregate.Segments) >= 1,
			"segments %v when more than 1 expected", len(aggregates.Themes),
		)
		if err != nil {
			return err
		}
		for _, segment := range aggregate.Segments {
			highestSet := int(segment.MaxSet)

			if highestSet < 1 {
				return fmt.Errorf("number of sets lower than 1")
			}
		}

	}

	return nil
}

func (l *OdysseiaFixture) aQueryIsMadeForTheOptionsForTheQuizType(quizType string) error {
	response, err := l.client.Sokrates().Options(quizType, TraceId)
	if err != nil {
		return err
	}

	var aggregates models.AggregatedOptions
	err = json.NewDecoder(response.Body).Decode(&aggregates)

	l.ctx = context.WithValue(l.ctx, AggregateContext, aggregates)

	return err
}

func (l *OdysseiaFixture) aNewQuizQuestionIsMadeWithTheQuizType(quizType string) error {
	aggregates := l.ctx.Value(AggregateContext).(models.AggregatedOptions)
	randomTheme := aggregates.Themes[l.randomizer.RandomNumberBaseZero(len(aggregates.Themes))]
	randomSegment := randomTheme.Segments[l.randomizer.RandomNumberBaseZero(len(randomTheme.Segments))]
	randomSet := strconv.Itoa(int(randomSegment.MaxSet))

	bodyModel := models.CreationRequest{
		Theme:    randomTheme.Name,
		Set:      randomSet,
		Segment:  randomSegment.Name,
		QuizType: quizType,
	}

	body, err := json.Marshal(bodyModel)
	if err != nil {
		return err
	}

	quiz, err := l.client.Sokrates().Create(body, TraceId)
	if err != nil {
		return err
	}

	if quiz.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 but got %v", quiz.StatusCode)
	}

	var quizResponse interface{}
	err = json.NewDecoder(quiz.Body).Decode(&quizResponse)

	l.ctx = context.WithValue(l.ctx, QuizContext, quizResponse)
	l.ctx = context.WithValue(l.ctx, BodyContext, bodyModel)
	return nil
}

func (l *OdysseiaFixture) theQuestionCanBeAnsweredFromTheResponse() error {
	quizResponse, ok := l.ctx.Value(QuizContext).(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to assert quizResponse type")
	}

	bodyMeta, ok := l.ctx.Value(BodyContext).(models.CreationRequest)
	if !ok {
		return fmt.Errorf("failed to assert bodyMeta type")
	}

	answerRequest := models.AnswerRequest{
		Theme:         bodyMeta.Theme,
		Set:           bodyMeta.Set,
		QuizType:      bodyMeta.QuizType,
		Segment:       bodyMeta.Segment,
		Comprehensive: false,
		Answer:        "",
		Dialogue:      nil,
		QuizWord:      "",
	}

	switch bodyMeta.QuizType {
	case models.MEDIA:
		var quiz models.QuizResponse
		quizBytes, err := json.Marshal(quizResponse)
		if err != nil {
			return fmt.Errorf("failed to marshal quizResponse: %v", err)
		}
		if err := json.Unmarshal(quizBytes, &quiz); err != nil {
			return fmt.Errorf("failed to unmarshal quizResponse into QuizResponse: %v", err)
		}

		for _, item := range quiz.Options {
			err := assertTrue(
				assert.True, strings.Contains(item.ImageUrl, "webp"),
				"expected webp to be included in: %s", item.ImageUrl,
			)
			if err != nil {
				return err
			}
		}

		randomQuizItem := quiz.Options[l.randomizer.RandomNumberBaseZero(len(quiz.Options))]
		answerRequest.Answer = randomQuizItem.Option
		answerRequest.QuizWord = quiz.QuizItem

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

	case models.MULTICHOICE:
		var quiz models.QuizResponse
		quizBytes, err := json.Marshal(quizResponse)
		if err != nil {
			return fmt.Errorf("failed to marshal quizResponse: %v", err)
		}
		if err := json.Unmarshal(quizBytes, &quiz); err != nil {
			return fmt.Errorf("failed to unmarshal quizResponse into QuizResponse: %v", err)
		}

		randomQuizItem := quiz.Options[l.randomizer.RandomNumberBaseZero(len(quiz.Options))]
		answerRequest.Answer = randomQuizItem.Option
		answerRequest.QuizWord = quiz.QuizItem

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

	case models.AUTHORBASED:
		var quiz models.AuthorbasedQuizResponse
		quizBytes, err := json.Marshal(quizResponse)
		if err != nil {
			return fmt.Errorf("failed to marshal quizResponse: %v", err)
		}
		if err := json.Unmarshal(quizBytes, &quiz); err != nil {
			return fmt.Errorf("failed to unmarshal quizResponse into QuizResponse: %v", err)
		}

		randomQuizItem := quiz.Quiz.Options[l.randomizer.RandomNumberBaseZero(len(quiz.Quiz.Options))]
		answerRequest.Answer = randomQuizItem.Option
		answerRequest.QuizWord = quiz.Quiz.QuizItem

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

		var answerResponse models.AuthorbasedQuizResponse
		err = json.NewDecoder(answer.Body).Decode(&answerResponse)
		if err != nil {
			return err
		}

	case models.DIALOGUE:
		var quiz models.DialogueQuiz
		quizBytes, err := json.Marshal(quizResponse)
		if err != nil {
			return fmt.Errorf("failed to marshal quizResponse: %v", err)
		}
		if err := json.Unmarshal(quizBytes, &quiz); err != nil {
			return fmt.Errorf("failed to unmarshal quizResponse into QuizResponse: %v", err)
		}

		answerRequest.Dialogue = quiz.Content

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

		var answerResponse models.DialogueAnswer
		err = json.NewDecoder(answer.Body).Decode(&answerResponse)
		if err != nil {
			return err
		}
		err = assertTrue(
			assert.True, answerResponse.Percentage == 100.00,
			"expected correctness to be 100% but was: %v", answerResponse.Percentage,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
