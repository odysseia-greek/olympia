package api

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

// This is your new mock data for ElasticSearch's response
const (
	edgecaseModel = `{
	  "took":4,
	  "timed_out":false,
	  "_shards":{
		"total":1,
		"successful":1,
		"skipped":0,
		"failed":0
	  },
	  "hits":{
		"total":{
		  "value":176,
		  "relation":"eq"
		},
		"max_score":21.65853,
		"hits":[
		  {
			"_index":"dictionary",
			"_id":"ky0c1pEBX4KIXZ5B04dm",
			"_score":21.65853,
			"_source":{
			  "english":"carry out",
			  "greek":"ἐκφέρω"
			}
		  },
		  {
			"_index":"dictionary",
			"_id":"hy0c1pEBX4KIXZ5B04pq",
			"_score":19.556658,
			"_source":{
			  "english":"carry on, make a difference",
			  "greek":"διαφέρω"
			}
		  }
		]
	  }
	}`
)

func TestCreateWordsEndpoint(t *testing.T) {
	englishToGreekDict, err := createEnglishToGreekDict()
	assert.Nil(t, err)

	t.Run("Words Convert", func(t *testing.T) {
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient([][]byte{[]byte(edgecaseModel)}, mockCode)
		assert.Nil(t, err)

		// This is the request payload
		request := models.EdgecaseRequest{
			Rootword: "ferw",
		}

		jsonBody, err := json.Marshal(request)
		bodyInBytes := bytes.NewReader(jsonBody)
		assert.Nil(t, err)

		// Create the DiogenesHandler and initialize it with necessary mock dependencies
		testConfig := DiogenesHandler{
			Elastic:            mockElasticClient,
			EnglishToGreekDict: englishToGreekDict, // This is set for all tests
		}

		// Initialize routes with the DiogenesHandler
		router := InitRoutes(&testConfig)

		// Perform the POST request to the /diogenes/v1/words/_convert endpoint
		response := performPostRequest(router, "/diogenes/v1/words/_convert", bodyInBytes)

		// Assert the expected response structure and values
		var edgecaseResponse models.EdgecaseResponse
		err = json.NewDecoder(response.Body).Decode(&edgecaseResponse)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)

		// Check expected values in the response
		expectedWord := "φερω"
		expectedStrongPasswordLength := 16
		expectedSimilarWordsLength := 2

		// Assertions on the response
		assert.Equal(t, expectedWord, edgecaseResponse.GreekWord)
		assert.Equal(t, expectedStrongPasswordLength, len(edgecaseResponse.StrongPassword))
		assert.Equal(t, expectedSimilarWordsLength, len(edgecaseResponse.SimilarWords))
	})

}

// Helper function to simulate POST requests
func performPostRequest(r http.Handler, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
