package app

import (
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/models"
	configs "github.com/odysseia-greek/olympia/demokritos/config"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestHandlerCreateDocuments(t *testing.T) {
	index := "test"
	body := models.Biblos{Biblos: []models.Meros{
		{
			Greek:   "ἀγγέλλω",
			English: "to bear a message",
		},
	},
	}

	t.Run("Created", func(t *testing.T) {
		file := "createDocument"
		status := 200
		var wg sync.WaitGroup
		wg.Add(1)
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := DemokritosHandler{Config: &testConfig}
		testHandler.AddDirectoryToElastic(body, &wg)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "createIndex"
		status := 502
		var wg sync.WaitGroup
		wg.Add(1)
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := DemokritosHandler{Config: &testConfig}
		testHandler.AddDirectoryToElastic(body, &wg)
		assert.Equal(t, 0, testConfig.Created)
	})
}

func TestTransformWord(t *testing.T) {
	index := "test"
	body := models.Meros{
		Greek:   "ἀγγέλλω",
		English: "to bear a message",
	}

	t.Run("Created", func(t *testing.T) {
		file := "createDocument"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := DemokritosHandler{Config: &testConfig}
		testHandler.transformWord(body)
		assert.Equal(t, 0, testConfig.Created)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "createIndex"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := DemokritosHandler{Config: &testConfig}
		testHandler.transformWord(body)
		assert.Equal(t, 0, testConfig.Created)
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

		testHandler := DemokritosHandler{Config: &testConfig}
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

		testHandler := DemokritosHandler{Config: &testConfig}
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

		testHandler := DemokritosHandler{Config: &testConfig}
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
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := DemokritosHandler{Config: &testConfig}
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

		testHandler := DemokritosHandler{Config: &testConfig}
		err = testHandler.CreateIndexAtStartup()
		assert.NotNil(t, err)
	})
}
