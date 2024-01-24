package gateway

import (
	"encoding/json"
	plato "github.com/odysseia-greek/agora/plato/models"
)

func (h *HomerosHandler) CreateDialogueQuiz(theme, set, quizType, requestID string) (*plato.DialogueQuiz, error) {
	request := plato.CreationRequest{
		Theme:    theme,
		Set:      set,
		QuizType: quizType,
	}

	body, err := json.Marshal(request)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}

	response, err := h.HttpClients.Sokrates().Create(body, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var quiz plato.DialogueQuiz
	err = json.NewDecoder(response.Body).Decode(&quiz)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, quiz)

	return &quiz, nil
}

func (h *HomerosHandler) CreateQuiz(theme, set, quizType, requestID string) (*plato.QuizResponse, error) {
	request := plato.CreationRequest{
		Theme:    theme,
		Set:      set,
		QuizType: quizType,
	}

	body, err := json.Marshal(request)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}

	response, err := h.HttpClients.Sokrates().Create(body, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var quiz plato.QuizResponse
	err = json.NewDecoder(response.Body).Decode(&quiz)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, quiz)

	return &quiz, nil
}

func (h *HomerosHandler) Check(answerRequest plato.AnswerRequest, requestID string) (*plato.ComprehensiveResponse, error) {
	jsonCheck, err := json.Marshal(answerRequest)
	if err != nil {
		return nil, err
	}

	response, err := h.HttpClients.Sokrates().Check(jsonCheck, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var answer plato.ComprehensiveResponse
	err = json.NewDecoder(response.Body).Decode(&answer)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, answer)

	return &answer, nil
}

func (h *HomerosHandler) CheckDialogue(answerRequest plato.AnswerRequest, requestID string) (*plato.DialogueAnswer, error) {
	jsonCheck, err := json.Marshal(answerRequest)
	if err != nil {
		return nil, err
	}

	response, err := h.HttpClients.Sokrates().Check(jsonCheck, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var answer plato.DialogueAnswer
	err = json.NewDecoder(response.Body).Decode(&answer)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, answer)

	return &answer, nil
}

func (h *HomerosHandler) Options(quizType string, requestID string) (*plato.AggregateResult, error) {
	response, err := h.HttpClients.Sokrates().Options(quizType, requestID)
	if err != nil {
		h.CloseTraceWithError(err, requestID)
		return nil, err
	}
	defer response.Body.Close()

	var aggregate plato.AggregateResult
	err = json.NewDecoder(response.Body).Decode(&aggregate)
	if err != nil {
		return nil, err
	}

	h.CloseTrace(response, aggregate)

	return &aggregate, nil
}
