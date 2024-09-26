package edgecase

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

const (
	edgecaseModel = `{
	  "took": 4,
	  "timed_out": false,
	  "_shards": {
	    "total": 1,
	    "successful": 1,
	    "skipped": 0,
	    "failed": 0
	  },
	  "hits": {
	    "total": {
	      "value": 176,
	      "relation": "eq"
	    },
	    "max_score": 21.65853,
	    "hits": [
	      {
	        "_index": "dictionary",
	        "_id": "ky0c1pEBX4KIXZ5B04dm",
	        "_score": 21.65853,
	        "_source": {
	          "english": "carry out",
	          "greek": "ἐκφέρω"
	        }
	      },
	      {
	        "_index": "dictionary",
	        "_id": "hy0c1pEBX4KIXZ5B04pq",
	        "_score": 19.556658,
	        "_source": {
	          "english": "carry on, make a difference",
	          "greek": "διαφέρω"
	        }
	      }
	    ]
	  }
	}`
)

func TestConvertWordEndpoint(t *testing.T) {
	// Create the English to Greek dictionary once for all tests
	englishToGreekDict, dictErr := createEnglishToGreekDict()
	assert.Nil(t, dictErr)

	t.Run("Word Conversion", func(t *testing.T) {
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient([][]byte{[]byte(edgecaseModel)}, mockCode)
		assert.Nil(t, err)

		request := models.EdgecaseRequest{
			Rootword: "ferw",
		}

		jsonBody, err := json.Marshal(request)
		bodyInBytes := bytes.NewReader(jsonBody)
		assert.Nil(t, err)

		testConfig := DiogenesHandler{
			Elastic:            mockElasticClient,
			EnglishToGreekDict: englishToGreekDict,
		}

		router := InitRoutes(&testConfig)
		response := performPostRequest(router, "/diogenes/v1/words/_convert", bodyInBytes)

		var sut models.EdgecaseResponse
		err = json.NewDecoder(response.Body).Decode(&sut)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)

		// Test cases
		expectedWord := "φερω"
		expectedStrongPasswordLength := 16
		expectedSimilarWordsLength := 2

		// Check if the Greek word matches the expected word
		assert.Equal(t, expectedWord, sut.GreekWord, "Expected Greek word 'φέρω' does not match the response")

		// Check that the strong password length matches the expected length
		assert.Equal(t, expectedStrongPasswordLength, len(sut.StrongPassword), "The length of the StrongPassword does not match 16")

		// Check that the length of similar words matches the expected length
		assert.Equal(t, expectedSimilarWordsLength, len(sut.SimilarWords), "The number of SimilarWords does not match 2")
	})
}

func performPostRequest(r http.Handler, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
