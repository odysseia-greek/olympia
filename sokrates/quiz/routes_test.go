package quiz

import (
	"bytes"
	"encoding/json"
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPongRoute(t *testing.T) {
	testConfig := &SokratesHandler{}
	router := InitRoutes(testConfig)
	expected := "{\"result\":\"pong\"}"

	w := performGetRequest(router, "/sokrates/v1/ping")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestHealth(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		fixtureFile := "info"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			SearchWord: "test",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/sokrates/v1/health")

		var healthModel models.Health
		err = json.NewDecoder(response.Body).Decode(&healthModel)
		assert.Nil(t, err)
		//models.Health
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, healthModel.Healthy)
	})

	t.Run("ElasticDown", func(t *testing.T) {
		fixtureFile := "shardFailure"
		mockCode := 500
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			SearchWord: "test",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/sokrates/v1/health")

		var healthModel models.Health
		err = json.NewDecoder(response.Body).Decode(&healthModel)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadGateway, response.Code)
		assert.False(t, healthModel.Healthy)
	})
}

func TestLastChapter(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		fixtureFile := "lastChapterSokrates"
		mockCode := 200
		expectedChapter := int64(15)
		category := "nomina"
		method := "mousieon"
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/methods/%s/categories/%s/chapters", method, category))

		var searchResults models.LastChapterResponse
		err = json.NewDecoder(response.Body).Decode(&searchResults)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expectedChapter, searchResults.LastChapter)
	})
	t.Run("ElasticDown", func(t *testing.T) {
		fixtureFile := "infoServiceDown"
		mockCode := 502
		category := "nomina"
		method := "method"
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/methods/%s/categories/%s/chapters", method, category))

		var searchResults models.ElasticSearchError
		err = json.NewDecoder(response.Body).Decode(&searchResults)
		assert.Nil(t, err)

		expectedText := "502 Bad Gateway"

		assert.Equal(t, http.StatusBadGateway, response.Code)
		assert.Contains(t, searchResults.Message.ElasticError, expectedText)
	})
	t.Run("BadRequest", func(t *testing.T) {
		fixtureFile := "lastChapterSokrates"
		mockCode := 200
		category := "f"
		method := "s"
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/methods/%s/categories/%s/chapters", method, category))

		var searchResults models.ValidationError
		err = json.NewDecoder(response.Body).Decode(&searchResults)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, "must be longer than 1", searchResults.Messages[0].Message)
		assert.Equal(t, "category", searchResults.Messages[0].Field)
	})
}

func TestCheckQuestion(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		fixtureFile := "checkQuestionSokrates"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		body := map[string]string{"answerProvided": "godin", "quizWord": "θεός", "category": "nomina"}

		jsonBody, _ := json.Marshal(body)
		bodyInBytes := bytes.NewReader(jsonBody)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performPostRequest(router, "/sokrates/v1/answer", bodyInBytes)

		var searchResults models.CheckAnswerResponse
		err = json.NewDecoder(response.Body).Decode(&searchResults)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, searchResults.Correct)
	})
	t.Run("ElasticDown", func(t *testing.T) {
		fixtureFile := "shardFailure"
		mockCode := 500
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		body := map[string]string{"answerProvided": "godin", "quizWord": "θεός", "category": "nomina"}

		jsonBody, _ := json.Marshal(body)
		bodyInBytes := bytes.NewReader(jsonBody)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performPostRequest(router, "/sokrates/v1/answer", bodyInBytes)

		var searchResults models.ElasticSearchError
		err = json.NewDecoder(response.Body).Decode(&searchResults)
		assert.Nil(t, err)

		expectedText := "500 Internal Server Error"

		assert.Equal(t, http.StatusBadGateway, response.Code)
		assert.Contains(t, searchResults.Message.ElasticError, expectedText)
	})
	t.Run("BadJson", func(t *testing.T) {
		fixtureFile := "shardFailure"
		mockCode := 500
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		body := map[int]string{1: "34"}

		jsonBody, _ := json.Marshal(body)
		bodyInBytes := bytes.NewReader(jsonBody)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performPostRequest(router, "/sokrates/v1/answer", bodyInBytes)

		var searchResults models.ValidationError
		err = json.NewDecoder(response.Body).Decode(&searchResults)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, "error parsing", searchResults.Messages[0].Message)
		assert.Contains(t, searchResults.Messages[0].Field, "body")
	})
}
func TestCreateQuestions(t *testing.T) {
	category := "verba"
	method := "mousieon"
	chapter := "1"
	mockRandomizer, err := randomizer.NewRandomizerClient()
	assert.Nil(t, err)

	t.Run("HappyPath", func(t *testing.T) {
		fixtureFile := "createQuestionSokrates"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			Randomizer: mockRandomizer,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/createQuestion?method=%s&category=%s&chapter=%s", method, category, chapter))

		var searchResults models.QuizResponse
		err = json.NewDecoder(response.Body).Decode(&searchResults)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)

		expectedGreek := [5]string{"ἀναβαίνω", "λέγω", "προβαίνω", "πονέω", "φέπω"}

		assert.Contains(t, expectedGreek, searchResults.Question)
	})

	t.Run("ShortQuiz", func(t *testing.T) {
		fixtureFile := "createQuestionSokratesShort"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			Randomizer: mockRandomizer,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/createQuestion?method=%s&category=%s&chapter=%s", method, category, chapter))

		var searchResults models.QuizResponse
		err = json.NewDecoder(response.Body).Decode(&searchResults)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)

		expectedGreek := [5]string{"ἀναβαίνω", "λέγω", "προβαίνω", "πονέω", "φέπω"}

		assert.Contains(t, expectedGreek, searchResults.Question)
	})

	t.Run("EmptyResponse", func(t *testing.T) {
		fixtureFile := "matchEmptyScroll"
		mockCode := 404
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			Randomizer: mockRandomizer,
			SearchWord: "greek",
		}

		expectedText := "a non 200 error message was returned"

		router := InitRoutes(testConfig)
		response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/createQuestion?method=%s&category=%s&chapter=%s", method, category, chapter))

		var searchResults models.ElasticSearchError
		err = json.NewDecoder(response.Body).Decode(&searchResults)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadGateway, response.Code)
		assert.Contains(t, searchResults.Message.ElasticError, expectedText)
	})

	t.Run("EmptyQuery", func(t *testing.T) {
		fixtureFile := "createQuestionSokratesShort"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			Randomizer: mockRandomizer,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/createQuestion?category=%s", category))

		var searchResults models.ValidationError
		err = json.NewDecoder(response.Body).Decode(&searchResults)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, "cannot be empty", searchResults.Messages[0].Message)
		assert.Contains(t, searchResults.Messages[0].Field, "chapter")
	})

	t.Run("ElasticDown", func(t *testing.T) {
		fixtureFile := "shardFailure"
		mockCode := 500
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			Randomizer: mockRandomizer,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, fmt.Sprintf("/sokrates/v1/createQuestion?method=%s&category=%s&chapter=%s", method, category, chapter))

		var searchResults models.ElasticSearchError
		err = json.NewDecoder(response.Body).Decode(&searchResults)
		assert.Nil(t, err)

		expectedText := "a non 200 error message was returned"

		assert.Equal(t, http.StatusBadGateway, response.Code)
		assert.Contains(t, searchResults.Message.ElasticError, expectedText)
	})
}

func TestMethods(t *testing.T) {
	expected := [3]string{"aristophanes", "logos", "herodotos"}

	t.Run("Pass", func(t *testing.T) {
		fixtureFile := "sokratesMethods"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/sokrates/v1/methods")

		var searchResults models.Methods
		err = json.NewDecoder(response.Body).Decode(&searchResults)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		for _, method := range searchResults.Method {
			assert.Contains(t, expected, method.Method)
		}
	})

	t.Run("ElasticDown", func(t *testing.T) {
		fixtureFile := "shardFailure"
		mockCode := 500
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/sokrates/v1/methods")

		var searchResults models.ElasticSearchError
		err = json.NewDecoder(response.Body).Decode(&searchResults)
		assert.Nil(t, err)

		expectedText := "a non 200 error message was returned"

		assert.Equal(t, http.StatusBadGateway, response.Code)
		assert.Contains(t, searchResults.Message.ElasticError, expectedText)
	})
}

func TestCategories(t *testing.T) {
	expected := "frogs"

	t.Run("Pass", func(t *testing.T) {
		fixtureFile := "sokratesCategories"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/sokrates/v1/methods/aristophanes/categories")

		var searchResults models.Categories
		err = json.NewDecoder(response.Body).Decode(&searchResults)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, searchResults.Category[0].Category, expected)
	})

	t.Run("ElasticDown", func(t *testing.T) {
		fixtureFile := "shardFailure"
		mockCode := 500
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := &SokratesHandler{
			Elastic:    mockElasticClient,
			SearchWord: "greek",
		}

		router := InitRoutes(testConfig)
		response := performGetRequest(router, "/sokrates/v1/methods/aristophanes/categories")

		var searchResults models.ElasticSearchError
		err = json.NewDecoder(response.Body).Decode(&searchResults)
		assert.Nil(t, err)

		expectedText := "a non 200 error message was returned"

		assert.Equal(t, http.StatusBadGateway, response.Code)
		assert.Contains(t, searchResults.Message.ElasticError, expectedText)
	})
}

func performGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performPostRequest(r http.Handler, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
