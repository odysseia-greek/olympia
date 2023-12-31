package dictionary

import (
	"encoding/json"
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	elasticmodels "github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPongRoute(t *testing.T) {
	testConfig := &AlexandrosHandler{}
	router := InitRoutes(testConfig)
	expected := "{\"result\":\"pong\"}"

	w := performGetRequest(router, "/alexandros/v1/ping")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestHealthEndpointHealthy(t *testing.T) {
	fixtureFile := "info"
	mockCode := 200
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := &AlexandrosHandler{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/alexandros/v1/health")

	var healthModel models.Health
	err = json.NewDecoder(response.Body).Decode(&healthModel)
	assert.Nil(t, err)
	//models.Health
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, healthModel.Healthy)
}

func TestHealthEndpointElasticDown(t *testing.T) {
	fixtureFile := "infoServiceDown"
	mockCode := 502
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := &AlexandrosHandler{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/alexandros/v1/health")

	var healthModel models.Health
	err = json.NewDecoder(response.Body).Decode(&healthModel)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.False(t, healthModel.Healthy)
}

func TestSearchShardFailure(t *testing.T) {
	fixtureFile := "shardFailure"
	mockCode := 500
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := &AlexandrosHandler{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	searchWord := "αγο"
	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/alexandros/v1/search?word=%s", searchWord))

	var searchResults models.ElasticSearchError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)
	expectedText := "500 Internal Server Error"

	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Contains(t, searchResults.Message.ElasticError, expectedText)
}

func TestSearchEndPointHappyPath(t *testing.T) {
	fixtureFile := "searchWord"
	mockCode := 200
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := &AlexandrosHandler{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	searchWord := "αγο"
	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/alexandros/v1/search?word=%s", searchWord))

	var searchResults []models.Meros
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, 2, len(searchResults))

	expectedGreek := [2]string{"ἀγορεύω", "ἀγορά, -ᾶς, ἡ"}

	for _, word := range searchResults {
		assert.Contains(t, expectedGreek, word.Greek)
	}
}

func TestSearchEndPointElasticDown(t *testing.T) {
	config := elasticmodels.Config{
		Service:     "hhttttt://sjdsj.com",
		Username:    "",
		Password:    "",
		ElasticCERT: "",
	}
	testClient, err := elastic.NewClient(config)
	assert.Nil(t, err)

	testConfig := &AlexandrosHandler{
		Elastic: testClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/alexandros/v1/search?word=αγο")

	var searchResults models.ElasticSearchError
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadGateway, response.Code)
}

func TestSearchEndPointNoResults(t *testing.T) {
	fixtureFile := "searchWordNoResults"
	mockCode := 200
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := &AlexandrosHandler{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	searchWord := "αγο"
	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/alexandros/v1/search?word=%s", searchWord))

	var searchResults []models.Meros
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, len(searchResults) == 0)
}

func TestSearchEndpointEmptyWord(t *testing.T) {
	fixtureFile := "searchWordNoResults"
	mockCode := 200
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := &AlexandrosHandler{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	searchWord := ""
	router := InitRoutes(testConfig)
	response := performGetRequest(router, fmt.Sprintf("/alexandros/v1/search?word=%s", searchWord))

	var searchResults models.ValidationError
	err = json.NewDecoder(response.Body).Decode(&searchResults)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, searchResults.Messages[0].Message, "cannot be empty")
	assert.Equal(t, "word", searchResults.Messages[0].Field)
}

func performGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
