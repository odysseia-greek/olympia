package app

import (
	elastic "github.com/odysseia-greek/agora/aristoteles"
	configs "github.com/odysseia-greek/olympia/anaximenes/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlerDeleteIndex(t *testing.T) {
	index := "test"

	t.Run("Deleted", func(t *testing.T) {
		file := "deleteIndex"
		status := 201
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic:    mockElasticClient,
			MaxAge:     "30",
			PolicyName: "TestName",
			Index:      index,
		}

		testHandler := AnaximenesConfig{Config: &testConfig}
		err = testHandler.DeleteIndexAtStartUp()
		assert.Nil(t, err)
	})

	t.Run("IndexDoesNotExist", func(t *testing.T) {
		file := "deleteIndex404"
		status := 404
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic:    mockElasticClient,
			MaxAge:     "30",
			PolicyName: "TestName",
			Index:      index,
		}

		testHandler := AnaximenesConfig{Config: &testConfig}
		err = testHandler.DeleteIndexAtStartUp()
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "error"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic:    mockElasticClient,
			MaxAge:     "30",
			PolicyName: "TestName",
			Index:      index,
		}

		testHandler := AnaximenesConfig{Config: &testConfig}
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
			MaxAge:     "30",
			PolicyName: "TestName",
			Index:      index,
		}

		testHandler := AnaximenesConfig{Config: &testConfig}
		err = testHandler.CreateIndexAtStartup()
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "error"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic:    mockElasticClient,
			MaxAge:     "30",
			PolicyName: "TestName",
			Index:      index,
		}
		testHandler := AnaximenesConfig{Config: &testConfig}
		err = testHandler.CreateIndexAtStartup()
		assert.NotNil(t, err)
	})
}
