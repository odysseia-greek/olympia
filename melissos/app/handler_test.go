package app

import (
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/thales"
	configs "github.com/odysseia-greek/olympia/melissos/config"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
)

func TestHandlerHandle(t *testing.T) {
	index := "test"

	body := models.Meros{
		Greek:      "ἀγορά",
		English:    "market place",
		LinkedWord: "",
		Original:   "",
	}

	bodyDutch := models.Meros{
		Greek:      "ἀγορά",
		Dutch:      "marktplaats",
		LinkedWord: "",
		Original:   "",
	}

	channel := "testchannel"
	dutchChannel := "testkanaal"

	t.Run("OneRunWithoutAction", func(t *testing.T) {
		mockClient := &MockEupalinosClient{}
		bodyString, err := body.Marshal()
		assert.Nil(t, err)
		os.Setenv("WAIT_TIME", "1")
		os.Setenv(TestData, string(bodyString))
		os.Setenv(TestLength, "0")
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic:      mockElasticClient,
			Index:        index,
			Created:      0,
			Eupalinos:    mockClient,
			Channel:      channel,
			DutchChannel: dutchChannel,
		}

		testHandler := MelissosHandler{Config: &testConfig}
		exitCode := testHandler.HandleParmenides()
		assert.True(t, exitCode)
		assert.Nil(t, err)

		os.Clearenv()
	})

	t.Run("DutchUpdate", func(t *testing.T) {
		mockClient := &MockEupalinosClient{}
		bodyString, err := bodyDutch.Marshal()
		assert.Nil(t, err)
		os.Setenv("WAIT_TIME", "1")
		os.Setenv(TestData, string(bodyString))
		os.Setenv(TestLength, "0")
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic:      mockElasticClient,
			Index:        index,
			Created:      0,
			Eupalinos:    mockClient,
			Channel:      channel,
			DutchChannel: dutchChannel,
		}

		testHandler := MelissosHandler{Config: &testConfig}
		exitCode := testHandler.HandleDutch()
		assert.True(t, exitCode)
		assert.Nil(t, err)

		os.Clearenv()
	})

}

func TestHandlerCreateDocuments(t *testing.T) {
	index := "test"

	body := models.Meros{
		Greek:      "ἀγορά",
		English:    "market place",
		LinkedWord: "",
		Original:   "",
	}

	t.Run("WordIsTheSame", func(t *testing.T) {
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		body.English = "a market place"

		testHandler := MelissosHandler{Config: &testConfig}
		found, err := testHandler.queryWord(body)
		assert.True(t, found)
		assert.Nil(t, err)
	})

	t.Run("WordWithAPronoun", func(t *testing.T) {
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := MelissosHandler{Config: &testConfig}
		found, err := testHandler.queryWord(body)
		assert.True(t, found)
		assert.Nil(t, err)
	})

	t.Run("WordFoundButDifferentMeaning", func(t *testing.T) {
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		body.English = "notthesame"

		testHandler := MelissosHandler{Config: &testConfig}
		found, err := testHandler.queryWord(body)
		assert.False(t, found)
		assert.Nil(t, err)
	})

	t.Run("WordFoundDifferentMeaningWithoutAPronoun", func(t *testing.T) {
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		body.English = "notthesame but multiple words"

		testHandler := MelissosHandler{Config: &testConfig}
		found, err := testHandler.queryWord(body)
		assert.False(t, found)
		assert.Nil(t, err)
	})

	t.Run("DoesNotExist", func(t *testing.T) {
		file := "searchWordNoResults"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := MelissosHandler{Config: &testConfig}
		found, err := testHandler.queryWord(body)
		assert.False(t, found)
		assert.Nil(t, err)
	})

	t.Run("DoesNotExist", func(t *testing.T) {
		file := "shardFailure"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := MelissosHandler{Config: &testConfig}
		found, err := testHandler.queryWord(body)
		assert.False(t, found)
		assert.NotNil(t, err)
	})
}

func TestHandlerAddWord(t *testing.T) {
	index := "test"
	body := models.Meros{
		Greek:      "ἀγορά",
		English:    "a market place",
		LinkedWord: "",
		Original:   "",
	}

	t.Run("DocumentNotCreated", func(t *testing.T) {
		file := "shardFailure"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler := MelissosHandler{Config: &testConfig}
		testHandler.addWord(body)
		assert.Equal(t, testConfig.Created, 0)
	})
}

func TestHandlerTransform(t *testing.T) {
	index := "test"
	body := models.Meros{
		Greek:      "ἀγορά",
		English:    "a market place",
		LinkedWord: "",
		Original:   "",
	}

	t.Run("DocumentCreated", func(t *testing.T) {
		file := "createDocument"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		var wait sync.WaitGroup

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		wait.Add(1)

		testHandler := MelissosHandler{Config: &testConfig}
		testHandler.transformWord(body, &wait)
		assert.Equal(t, testConfig.Created, 1)
	})

	t.Run("DocumentNotCreated", func(t *testing.T) {
		file := "shardFailure"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		var wait sync.WaitGroup

		testConfig := configs.Config{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		wait.Add(1)

		testHandler := MelissosHandler{Config: &testConfig}
		testHandler.transformWord(body, &wait)
		assert.Equal(t, testConfig.Created, 0)
	})

}

func TestHandlerParser(t *testing.T) {
	splitWord := "ἀκούω + gen."
	pronounSplit := "μῦθος, ὁ"
	pronounSplitTwo := "ὁ δοῦλος"
	testConfig := configs.Config{
		Elastic: nil,
		Index:   "",
		Created: 0,
	}

	testHandler := MelissosHandler{Config: &testConfig}

	t.Run("SplitWordsWithPlus", func(t *testing.T) {
		sut := "ἀκούω"
		parsedWord := testHandler.stripMouseionWords(splitWord)
		assert.Equal(t, sut, parsedWord)
	})

	t.Run("SplitWordsWithPronoun", func(t *testing.T) {
		sut := "μῦθος"
		parsedWord := testHandler.stripMouseionWords(pronounSplit)
		assert.Equal(t, sut, parsedWord)
	})

	t.Run("SplitWordsWithPronounWithoutComma", func(t *testing.T) {
		sut := "δοῦλος"
		parsedWord := testHandler.stripMouseionWords(pronounSplitTwo)
		assert.Equal(t, sut, parsedWord)
	})
}

func TestJobExit(t *testing.T) {
	ns := "odysseia"
	expectedName := "testpod"
	duration := 10 * time.Millisecond
	timeFinished := int64(1000)

	t.Run("JobFinished", func(t *testing.T) {
		testClient, err := thales.FakeKubeClient(ns)
		assert.Nil(t, err)

		jobSpec := thales.CreateJobObject(expectedName, ns, true)
		job, err := testClient.Workload().CreateJob(ns, jobSpec)
		assert.Nil(t, err)
		assert.Equal(t, job.Name, expectedName)

		testConfig := configs.Config{
			Kube:      testClient,
			Job:       expectedName,
			Namespace: ns,
		}

		handler := MelissosHandler{Config: &testConfig, Duration: duration, TimeFinished: timeFinished}
		jobExit := make(chan bool, 1)
		go handler.WaitForJobsToFinish(jobExit)

		select {

		case <-jobExit:
			exitStatus := <-jobExit
			assert.True(t, exitStatus)
		}
	})

	t.Run("JobNotFinished", func(t *testing.T) {
		testClient, err := thales.FakeKubeClient(ns)
		assert.Nil(t, err)

		jobSpec := thales.CreateJobObject(expectedName, ns, false)
		job, err := testClient.Workload().CreateJob(ns, jobSpec)
		assert.Nil(t, err)
		assert.Equal(t, job.Name, expectedName)

		testConfig := configs.Config{
			Kube:      testClient,
			Job:       expectedName,
			Namespace: ns,
		}

		timeFinished = duration.Milliseconds() * 2

		handler := MelissosHandler{Config: &testConfig, Duration: duration, TimeFinished: timeFinished}
		jobExit := make(chan bool, 1)
		go handler.WaitForJobsToFinish(jobExit)

		select {

		case <-jobExit:
			exitStatus := <-jobExit
			assert.False(t, exitStatus)
		}
	})
}
