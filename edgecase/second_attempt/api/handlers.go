package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	elasticmodels "github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/models"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"net/http"
	"strings"
	"time"
)

type DiogenesHandler struct {
	Elastic            aristoteles.Client
	Index              string
	Streamer           pb.TraceService_ChorusClient
	Cancel             context.CancelFunc
	EnglishToGreekDict map[string]string
}

// PingPong pongs the ping
func (d *DiogenesHandler) pingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

// returns the health of the api
func (d *DiogenesHandler) health(w http.ResponseWriter, req *http.Request) {
	elasticHealth := d.Elastic.Health().Info()
	dbHealth := models.DatabaseHealth{
		Healthy:       elasticHealth.Healthy,
		ClusterName:   elasticHealth.ClusterName,
		ServerName:    elasticHealth.ServerName,
		ServerVersion: elasticHealth.ServerVersion,
	}

	healthy := models.Health{
		Healthy:  dbHealth.Healthy,
		Time:     time.Now().String(),
		Database: dbHealth,
	}
	if !healthy.Healthy {
		middleware.ResponseWithCustomCode(w, http.StatusBadGateway, healthy)
		return
	}

	middleware.ResponseWithJson(w, healthy)
}

// Example
func (d *DiogenesHandler) convert(w http.ResponseWriter, req *http.Request) {
	var requestId string
	fromContext := req.Context().Value(config.DefaultTracingName)
	if fromContext == nil {
		requestId = req.Header.Get(config.HeaderKey)
	} else {
		requestId = fromContext.(string)
	}
	splitID := strings.Split(requestId, "+")

	traceCall := false
	var traceID, spanID string

	if len(splitID) >= 3 {
		traceCall = splitID[2] == "1"
	}

	if len(splitID) >= 1 {
		traceID = splitID[0]
	}
	if len(splitID) >= 2 {
		spanID = splitID[1]
	}

	var edgecaseRequest models.EdgecaseRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&edgecaseRequest)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "decoding",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	greekWord, err := d.ConvertToGreek(edgecaseRequest.Rootword)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "parsing greek word",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	query := buildFuzzyQuery(greekWord)

	elasticResult, err := d.Elastic.Query().Match(d.Index, query)
	jsonResult, _ := json.Marshal(elasticResult)
	logging.Info(fmt.Sprintf("elastic result: %s", jsonResult))
	if err != nil {
		e := models.ElasticSearchError{
			ErrorModel: models.ErrorModel{UniqueCode: requestId},
			Message: models.ElasticErrorMessage{
				ElasticError: err.Error(),
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}

	var meroi []models.Meros

	for _, hit := range elasticResult.Hits.Hits {
		jsonHit, _ := json.Marshal(hit.Source)
		meros, _ := models.UnmarshalMeros(jsonHit)
		meroi = append(meroi, meros)
	}
	if traceCall {
		go d.databaseSpan(elasticResult, query, traceID, spanID)
	}

	strongPassword, _ := d.GenerateStrongPassword(edgecaseRequest.Rootword, 16)

	response := models.EdgecaseResponse{
		OriginalWord:   edgecaseRequest.Rootword,
		GreekWord:      greekWord,
		StrongPassword: strongPassword,
		SimilarWords:   meroi,
	}

	middleware.ResponseWithCustomCode(w, http.StatusOK, response)
}

func (d *DiogenesHandler) databaseSpan(response *elasticmodels.Response, query map[string]interface{}, traceID, spanID string) {
	parsedQuery, _ := json.Marshal(query)
	hits := int64(0)
	took := int64(0)
	if response != nil {
		hits = response.Hits.Total.Value
		took = response.Took
	}
	dataBaseSpan := &pb.ParabasisRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		SpanId:       spanID,
		RequestType: &pb.ParabasisRequest_DatabaseSpan{DatabaseSpan: &pb.DatabaseSpanRequest{
			Action:   "search",
			Query:    string(parsedQuery),
			Hits:     hits,
			TimeTook: took,
		}},
	}

	err := d.Streamer.Send(dataBaseSpan)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}
}

// ConvertToGreek This method converts an English word to Greek using the given mapping.
func (d *DiogenesHandler) ConvertToGreek(word string) (string, error) {
	var greekWord strings.Builder
	wordLen := len(word)

	for i := 0; i < wordLen; i++ {
		currentChar := string(word[i])

		// If 's' is the last character of the word, replace it with 's_end'
		if i == wordLen-1 && currentChar == "s" {
			if greekLetter, exists := d.EnglishToGreekDict["s_end"]; exists {
				greekWord.WriteString(greekLetter)
				continue
			}
		}

		// Look for diacritical marks or multi-character sequences
		if i < wordLen-1 {
			possibleMultiChar := word[i : i+2]
			if greekLetter, exists := d.EnglishToGreekDict[possibleMultiChar]; exists {
				greekWord.WriteString(greekLetter)
				i++ // Skip next character
				continue
			}
		}

		// Otherwise, map the current character
		if greekLetter, exists := d.EnglishToGreekDict[currentChar]; exists {
			greekWord.WriteString(greekLetter)
		} else {
			// If the character is not found in the dictionary, just add it as is
			greekWord.WriteString(currentChar)
		}
	}

	return greekWord.String(), nil
}

// GenerateStrongPassword combines the given word with the current timestamp
// and generates a strong password of the desired length
func (d *DiogenesHandler) GenerateStrongPassword(word string, length int) (string, error) {
	// Get the current timestamp in Unix format
	timestamp := time.Now().Unix()

	// Combine the word with the timestamp
	combined := fmt.Sprintf("%s%d", word, timestamp)

	// Hash the combined word + timestamp using SHA-256 for strong randomness
	hash := sha256.New()
	hash.Write([]byte(combined))
	hashedBytes := hash.Sum(nil)

	// Convert the hash to a hex string
	hexString := hex.EncodeToString(hashedBytes)

	// Ensure the password length doesn't exceed the available hex string length
	if length > len(hexString) {
		length = len(hexString)
	}

	// Take the desired portion of the hashed string for the password
	password := hexString[:length]

	// Optionally, mix in uppercase, numbers, and special characters for a stronger password
	password = d.addSpecialCharacters(password)

	return password, nil
}

// addSpecialCharacters adds some variety to the password by replacing some characters
func (d *DiogenesHandler) addSpecialCharacters(password string) string {
	replacements := map[string]string{
		"a": "@", "b": "8", "e": "3", "l": "1", "o": "0", "s": "$", "t": "7",
	}

	var finalPassword strings.Builder
	for _, char := range password {
		if replacement, ok := replacements[string(char)]; ok {
			finalPassword.WriteString(replacement)
		} else {
			finalPassword.WriteRune(char)
		}
	}

	return finalPassword.String()
}

// buildFuzzyQuery builds a fuzzy query for Elasticsearch that matches
// the value of rootWord with a fuzziness of 2 and contains the term "greek".
// It limits the results to 5 documents.
func buildFuzzyQuery(rootWord string) map[string]interface{} {
	query := map[string]interface{}{
		"size": 5,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"fuzzy": map[string]interface{}{
							"greek": map[string]interface{}{
								"value":     rootWord,
								"fuzziness": 2,
							},
						},
					},
				},
			},
		},
	}
	return query
}
