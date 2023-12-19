package scholar

import (
	"context"
	"github.com/odysseia-greek/agora/aristoteles"
	pb "github.com/odysseia-greek/olympia/aristarchos/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

const foundModel = `{
  "took": 8,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": {
      "value": 8,
      "relation": "eq"
    },
    "max_score": 1.0,
    "hits": [
      {
        "_index": "dictionary",
        "_type": "_doc",
        "_id": "Wkzd1ocBSKBq_nTnS81A",
        "_score": 1.0,
        "_source": {
          "conjugations":[
            {
              "aspect":"impf",
              "forms":[
                {
                  "number":"sing",
                  "person":"3th",
                  "word":"ἔβαλλε"
                },
                {
                  "number":"sing",
                  "person":"2nd",
                  "word":"ἔβᾰλλες"
                }
              ],
              "mood":"ind",
              "tense":"act"
            },
            {
              "aspect":"pres",
              "forms":[
                {
                  "number":"plur",
                  "person":"2nd",
                  "word":"βᾰ́λλετε"
                }
              ],
              "mood":"ind",
              "tense":"act"
            }
          ],
          "rootWord":"βαλλω",
          "translation":"throw",
          "translations":[
            "throw"
          ]
        }
      }
    ]
  }
}`

func TestHealth(t *testing.T) {
	t.Run("Health", func(t *testing.T) {
		request := pb.HealthRequest{}
		handler := &AggregatorServiceImpl{
			Elastic:                        nil,
			Index:                          "test",
			Tracer:                         nil,
			UnimplementedAristarchosServer: pb.UnimplementedAristarchosServer{},
		}

		response, err := handler.Health(context.Background(), &request)

		assert.Nil(t, err)
		assert.True(t, response.Health)
	})
}

func TestCreateEntry(t *testing.T) {
	request := pb.AggregatorCreationRequest{
		Word:        "βᾰ́λλετε",
		Rule:        "2nd plur - pres - ind - act",
		RootWord:    "βαλλω",
		Translation: "throw",
	}

	updated := `{
  "_index": "tracing-2023.08.15",
  "_id": "841a4f73-ba5b-4c38-9237-e1ad91459028",
  "_version": 2,
  "result": "updated",
  "_shards": {
    "total": 2,
    "successful": 1,
    "failed": 0
  },
  "_seq_no": 119,
  "_primary_term": 4
}`

	t.Run("CreateNewEntry", func(t *testing.T) {
		fixtureFiles := []string{"matchEmptyHits", "createDocument"}
		mockCode := 200
		mockElasticClient, err := aristoteles.NewMockClient(fixtureFiles, mockCode)
		assert.Nil(t, err)
		handler := &AggregatorServiceImpl{
			Elastic:                        mockElasticClient,
			Index:                          "test",
			Tracer:                         nil,
			UnimplementedAristarchosServer: pb.UnimplementedAristarchosServer{},
		}

		response, err := handler.CreateNewEntry(context.Background(), &request)

		assert.Nil(t, err)
		assert.True(t, response.Created)
		assert.False(t, response.Updated)
	})

	t.Run("UpdateEntryBasedOnConjugations", func(t *testing.T) {
		mockCode := 200
		req := pb.AggregatorCreationRequest{
			Word:        "βᾰ́λλετε",
			Rule:        "1st plur - aor - ind - act",
			RootWord:    "βαλλω",
			Translation: "throw",
		}
		mockElasticClient, err := aristoteles.NewMockClient([][]byte{[]byte(foundModel), []byte(updated)}, mockCode)
		assert.Nil(t, err)
		handler := &AggregatorServiceImpl{
			Elastic:                        mockElasticClient,
			Index:                          "test",
			Tracer:                         nil,
			UnimplementedAristarchosServer: pb.UnimplementedAristarchosServer{},
		}

		response, err := handler.CreateNewEntry(context.Background(), &req)

		assert.Nil(t, err)
		assert.False(t, response.Created)
		assert.True(t, response.Updated)
	})

	t.Run("UpdateEntryBasedOnForms", func(t *testing.T) {
		mockCode := 200
		req := pb.AggregatorCreationRequest{
			Word:        "ἔβᾰλλον",
			Rule:        "1st sing - impf - ind - act",
			RootWord:    "βαλλω",
			Translation: "throw",
		}
		mockElasticClient, err := aristoteles.NewMockClient([][]byte{[]byte(foundModel), []byte(updated)}, mockCode)
		assert.Nil(t, err)
		handler := &AggregatorServiceImpl{
			Elastic:                        mockElasticClient,
			Index:                          "test",
			Tracer:                         nil,
			UnimplementedAristarchosServer: pb.UnimplementedAristarchosServer{},
		}

		response, err := handler.CreateNewEntry(context.Background(), &req)

		assert.Nil(t, err)
		assert.False(t, response.Created)
		assert.True(t, response.Updated)
	})

	t.Run("UpdateEntryNotNeeded", func(t *testing.T) {
		jsonData := []byte(foundModel)
		mockCode := 200
		mockElasticClient, err := aristoteles.NewMockClient([][]byte{jsonData}, mockCode)
		assert.Nil(t, err)
		handler := &AggregatorServiceImpl{
			Elastic:                        mockElasticClient,
			Index:                          "test",
			Tracer:                         nil,
			UnimplementedAristarchosServer: pb.UnimplementedAristarchosServer{},
		}

		response, err := handler.CreateNewEntry(context.Background(), &request)

		assert.Nil(t, err)
		assert.False(t, response.Created)
		assert.False(t, response.Updated)
	})

	t.Run("ElasticDown", func(t *testing.T) {
		fixtureFile := "serviceDown"
		mockCode := 502
		mockElasticClient, err := aristoteles.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		handler := &AggregatorServiceImpl{
			Elastic:                        mockElasticClient,
			Index:                          "test",
			Tracer:                         nil,
			UnimplementedAristarchosServer: pb.UnimplementedAristarchosServer{},
		}

		response, err := handler.CreateNewEntry(context.Background(), &request)

		assert.NotNil(t, err)
		assert.False(t, response.Created)
		assert.False(t, response.Updated)
	})
}

func TestRetrieveEntries(t *testing.T) {
	request := pb.AggregatorRequest{
		RootWord: "βαλλω",
	}

	t.Run("Retrieve", func(t *testing.T) {
		mockCode := 200
		mockElasticClient, err := aristoteles.NewMockClient([][]byte{[]byte(foundModel)}, mockCode)
		assert.Nil(t, err)
		handler := &AggregatorServiceImpl{
			Elastic:                        mockElasticClient,
			Index:                          "test",
			Tracer:                         nil,
			UnimplementedAristarchosServer: pb.UnimplementedAristarchosServer{},
		}

		response, err := handler.RetrieveEntry(context.Background(), &request)

		assert.Nil(t, err)
		assert.Equal(t, request.RootWord, response.RootWord)
	})

	t.Run("SearchWords", func(t *testing.T) {
		mockCode := 200
		mockElasticClient, err := aristoteles.NewMockClient([][]byte{[]byte(foundModel)}, mockCode)
		assert.Nil(t, err)
		handler := &AggregatorServiceImpl{
			Elastic:                        mockElasticClient,
			Index:                          "test",
			Tracer:                         nil,
			UnimplementedAristarchosServer: pb.UnimplementedAristarchosServer{},
		}

		response, err := handler.RetrieveSearchWords(context.Background(), &request)

		sut := []string{"ἔβᾰλλες", "βᾰ́λλετε"}
		assert.Nil(t, err)

		for _, word := range sut {
			assert.Contains(t, response.Word, word)
		}
	})

	t.Run("ElasticDown", func(t *testing.T) {
		fixtureFile := "serviceDown"
		mockCode := 502
		mockElasticClient, err := aristoteles.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		handler := &AggregatorServiceImpl{
			Elastic:                        mockElasticClient,
			Index:                          "test",
			Tracer:                         nil,
			UnimplementedAristarchosServer: pb.UnimplementedAristarchosServer{},
		}

		response, err := handler.RetrieveEntry(context.Background(), &request)

		assert.NotNil(t, err)
		assert.Nil(t, response)
	})
}
