package text

import (
	"bytes"
	"encoding/json"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPongRoute(t *testing.T) {
	testConfig := &HerodotosHandler{}
	router := InitRoutes(testConfig)
	expected := "{\"result\":\"pong\"}"

	w := performGetRequest(router, "/herodotos/v1/ping")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expected, w.Body.String())
}

func TestHealthEndpointHealthy(t *testing.T) {
	fixtureFile := "info"
	mockCode := 200
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := &HerodotosHandler{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/herodotos/v1/health")

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

	testConfig := &HerodotosHandler{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	router := InitRoutes(testConfig)
	response := performGetRequest(router, "/herodotos/v1/health")

	var healthModel models.Health
	err = json.NewDecoder(response.Body).Decode(&healthModel)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.False(t, healthModel.Healthy)
}

func TestCreateQuestionHappyPath(t *testing.T) {
	fixtureFile := "createQuestionHerodotos"
	mockCode := 200
	mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	testConfig := &HerodotosHandler{
		Elastic: mockElasticClient,
		Index:   "test",
	}

	body := models.CreateTextRequest{
		Author:    "thucydides",
		Book:      "",
		Reference: "",
		Section:   "",
	}

	jsonBody, _ := json.Marshal(body)
	bodyInBytes := bytes.NewReader(jsonBody)

	router := InitRoutes(testConfig)
	response := performPostRequest(router, "/herodotos/v1/texts/_create", bodyInBytes)

	var searchResults models.CheckTextResponse
	err = json.NewDecoder(response.Body).Decode(&searchResults)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.Code)
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
