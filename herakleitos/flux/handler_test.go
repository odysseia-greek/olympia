package flux

import (
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestHandlerCreateDocuments(t *testing.T) {
	index := "test"
	body := models.Rhema{Rhemai: []models.RhemaSource{{
		Author:          "test",
		Greek:           "test",
		Translations:    []string{"test"},
		Book:            0,
		Chapter:         0,
		Section:         0,
		PerseusTextLink: "test",
	},
	}}

	t.Run("Created", func(t *testing.T) {
		file := "createDocument"
		status := 200
		var wg sync.WaitGroup
		wg.Add(1)
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := HerakleitosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		err = testHandler.Add(body, &wg)
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "createIndex"
		status := 502
		var wg sync.WaitGroup
		wg.Add(1)
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := HerakleitosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}
		err = testHandler.Add(body, &wg)
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

		testHandler := HerakleitosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}
		err = testHandler.DeleteIndexAtStartUp()
		assert.Nil(t, err)
	})

	t.Run("IndexDoesNotExist", func(t *testing.T) {
		file := "deleteIndex404"
		status := 404
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := HerakleitosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}
		err = testHandler.DeleteIndexAtStartUp()
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "error"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := HerakleitosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}
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

		testHandler := HerakleitosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}
		err = testHandler.CreateIndexAtStartup()
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "error"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := HerakleitosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}
		err = testHandler.CreateIndexAtStartup()
		assert.NotNil(t, err)
	})
}
