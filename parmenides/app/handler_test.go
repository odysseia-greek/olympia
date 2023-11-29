package app

import (
	"context"
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/models"
	pb "github.com/odysseia-greek/eupalinos/proto"
	configs "github.com/odysseia-greek/olympia/parmenides/config"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestParmenidesHandlerAdd(t *testing.T) {
	index := "test"
	body := models.Logos{Logos: []models.Word{{
		Method:      "",
		Category:    "",
		Greek:       "ἀγγέλλω",
		Translation: "to bear a message ",
		Chapter:     0,
	},
	},
	}

	method := "testmethod"
	category := "testcategory"
	channel := "testchannel"
	mockClient := &MockEupalinosClient{}

	t.Run("CreatedWithQueue", func(t *testing.T) {
		file := "createDocument"
		status := 200
		var wg sync.WaitGroup
		wg.Add(1)
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic:   mockElasticClient,
			Index:     index,
			Created:   0,
			Channel:   channel,
			Eupalinos: mockClient,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.Add(body, &wg, method, category)
		assert.Nil(t, err)
	})

	t.Run("EnqueueTask", func(t *testing.T) {
		testConfig := configs.Config{
			Index:     index,
			Created:   0,
			Channel:   channel,
			Eupalinos: mockClient,
		}

		// Use the mock client implementation

		testHandler := ParmenidesHandler{
			Config: &testConfig,
		}

		message, err := body.Marshal()
		assert.Nil(t, err)
		msg := &pb.Epistello{
			Data:    string(message),
			Channel: channel,
		}

		err = testHandler.EnqueueTask(context.Background(), msg)
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "createIndex"
		status := 502
		var wg sync.WaitGroup
		wg.Add(1)
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic:   mockElasticClient,
			Index:     index,
			Created:   0,
			Eupalinos: mockClient,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.Add(body, &wg, method, category)
		assert.NotNil(t, err)
	})
}

func TestHandlerDeleteIndex(t *testing.T) {
	index := "test"

	t.Run("Deleted", func(t *testing.T) {
		file := "deleteIndex"
		status := 201
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.DeleteIndexAtStartUp()
		assert.Nil(t, err)
	})

	t.Run("IndexDoesNotExist", func(t *testing.T) {
		file := "deleteIndex404"
		status := 404
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.DeleteIndexAtStartUp()
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "error"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.DeleteIndexAtStartUp()
		assert.NotNil(t, err)
	})
}

func TestHandlerCreateIndex(t *testing.T) {
	index := "test"

	t.Run("Created", func(t *testing.T) {
		files := []string{"createIndex", "createIndex"}
		status := 201
		mockElasticClient, err := elastic.NewMockClient(files, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic:    mockElasticClient,
			Index:      index,
			Created:    0,
			PolicyName: fmt.Sprintf("%s-policy", index),
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.CreateIndexAtStartup()
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "error"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := ParmenidesHandler{Config: &testConfig}
		err = testHandler.CreateIndexAtStartup()
		assert.NotNil(t, err)
	})
}
