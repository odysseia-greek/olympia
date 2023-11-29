package hippokrates

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
	"strconv"
)

const (
	ContextCategories string = "contextCategories"
	ContextChapter    string = "contextChapter"
	ContextMethod     string = "contextMethod"
	ContextAnswer     string = "contextAnswer"
	TraceId           string = "hippokrates-traceid"
)

func (l *OdysseiaFixture) aQueryIsMadeForAllMethods() error {
	response, err := l.client.Sokrates().GetMethods(TraceId)
	if err != nil {
		return err
	}

	var methods models.Methods
	err = json.NewDecoder(response.Body).Decode(&methods)

	l.ctx = context.WithValue(l.ctx, ResponseBody, methods)

	return nil
}

func (l *OdysseiaFixture) aRandomMethodIsQueriedForCategories() error {
	methods := l.ctx.Value(ResponseBody).(models.Methods)

	randNumber := GenerateRandomNumber(len(methods.Method))
	method := methods.Method[randNumber].Method

	categoriesResponse, err := l.client.Sokrates().GetCategories(method, TraceId)
	if err != nil {
		return err
	}

	var categories models.Categories
	err = json.NewDecoder(categoriesResponse.Body).Decode(&categories)

	l.ctx = context.WithValue(l.ctx, ContextCategories, categories)
	l.ctx = context.WithValue(l.ctx, ContextMethod, method)

	return nil
}

func (l *OdysseiaFixture) aRandomCategoryIsQueriedForTheLastChapter() error {
	categories := l.ctx.Value(ContextCategories).(models.Categories)
	method := l.ctx.Value(ContextMethod).(string)

	randNumber := GenerateRandomNumber(len(categories.Category))
	category := categories.Category[randNumber].Category

	lastChapterResponse, err := l.client.Sokrates().GetChapters(method, category, TraceId)
	if err != nil {
		return err
	}

	var lastChapter models.LastChapterResponse
	err = json.NewDecoder(lastChapterResponse.Body).Decode(&lastChapter)

	l.ctx = context.WithValue(l.ctx, ContextChapter, lastChapter)

	return nil
}

func (l *OdysseiaFixture) aNewQuizQuestionIsRequested() error {
	resp, err := l.client.Sokrates().GetMethods(TraceId)
	if err != nil {
		return err
	}

	var methods models.Methods
	err = json.NewDecoder(resp.Body).Decode(&methods)

	randomMethod := GenerateRandomNumber(len(methods.Method))
	method := methods.Method[randomMethod].Method

	categoryResponse, err := l.client.Sokrates().GetCategories(method, TraceId)
	if err != nil {
		return err
	}

	var categories models.Categories
	err = json.NewDecoder(categoryResponse.Body).Decode(&categories)

	randomCategory := GenerateRandomNumber(len(categories.Category))
	category := categories.Category[randomCategory].Category

	lastChapterResponse, err := l.client.Sokrates().GetChapters(method, category, TraceId)
	if err != nil {
		return err
	}

	var lastChapter models.LastChapterResponse
	err = json.NewDecoder(lastChapterResponse.Body).Decode(&lastChapter)

	randomChapter := GenerateRandomNumber(int(lastChapter.LastChapter)) + 1

	chapter := strconv.Itoa(randomChapter)
	quizResponse, err := l.client.Sokrates().CreateQuestion(method, category, chapter, TraceId)
	if err != nil {
		return err
	}

	var quizQuestion models.QuizResponse
	err = json.NewDecoder(quizResponse.Body).Decode(&quizQuestion)

	l.ctx = context.WithValue(l.ctx, ContextCategories, category)
	l.ctx = context.WithValue(l.ctx, ResponseBody, quizQuestion)

	return nil
}

func (l *OdysseiaFixture) thatQuestionIsAnsweredWithAAnswer(correctAnswer string) error {
	quiz := l.ctx.Value(ResponseBody).(models.QuizResponse)

	checkAnswerRequest := models.CheckAnswerRequest{
		QuizWord:       quiz.Question,
		AnswerProvided: "",
	}

	parsedAnswer, err := strconv.ParseBool(correctAnswer)
	if err != nil {
		return err
	}

	if parsedAnswer {
		checkAnswerRequest.AnswerProvided = quiz.Answer
	} else {
		for _, answer := range quiz.QuizQuestions {
			if answer != quiz.Answer {
				checkAnswerRequest.AnswerProvided = answer
				break
			}
		}
	}

	request, err := checkAnswerRequest.Marshal()
	if err != nil {
		return err
	}

	answerResponse, err := l.client.Sokrates().CheckAnswer(request, TraceId)
	if err != nil {
		return err
	}

	var answer models.CheckAnswerResponse
	err = json.NewDecoder(answerResponse.Body).Decode(&answer)

	l.ctx = context.WithValue(l.ctx, ContextAnswer, answer.Correct)

	return nil
}

func (l *OdysseiaFixture) theMethodShouldBeIncluded(method string) error {
	methods := l.ctx.Value(ResponseBody).(models.Methods)

	found := false

	for _, result := range methods.Method {
		if result.Method == method {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("could not find book %v in slice", method)
	}

	return nil
}

func (l *OdysseiaFixture) theNumberOfMethodsShouldExceed(results int) error {
	methods := l.ctx.Value(ResponseBody).(models.Methods)
	numberOfMethods := len(methods.Method)
	if numberOfMethods <= results {
		return fmt.Errorf("expected results to be equal to or more than %v but was %v", results, numberOfMethods)
	}

	return nil
}

func (l *OdysseiaFixture) aCategoryShouldBeReturned() error {
	categories := l.ctx.Value(ContextCategories).(models.Categories)

	if len(categories.Category) == 0 {
		return fmt.Errorf("expected categories to be returned but non were found")
	}

	return nil
}

func (l *OdysseiaFixture) thatChapterShouldBeANumberAbove(number int) error {
	lastChapter := l.ctx.Value(ContextChapter).(models.LastChapterResponse)
	if lastChapter.LastChapter < int64(number) {
		return fmt.Errorf("expected lastchapter to be higher than %v but was %v", number, lastChapter.LastChapter)
	}

	return nil
}

func (l *OdysseiaFixture) theResultShouldBe(correct string) error {
	answer := l.ctx.Value(ContextAnswer).(bool)

	parsedCorrectness, err := strconv.ParseBool(correct)
	if err != nil {
		return err
	}

	if answer != parsedCorrectness {
		return fmt.Errorf("expected answer %v to be equal to correctness %v", answer, correct)
	}
	return nil
}
