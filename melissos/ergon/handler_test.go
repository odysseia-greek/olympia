package ergon

import (
	"context"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/thales"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

		testHandler := MelissosHandler{
			Elastic:      mockElasticClient,
			Index:        index,
			Created:      0,
			Eupalinos:    mockClient,
			Channel:      channel,
			DutchChannel: dutchChannel,
		}

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

		testHandler := MelissosHandler{
			Elastic:      mockElasticClient,
			Index:        index,
			Created:      0,
			Eupalinos:    mockClient,
			Channel:      channel,
			DutchChannel: dutchChannel,
		}

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

		testHandler := MelissosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		body.English = "a market place"

		found, err := testHandler.queryWord(body)
		assert.True(t, found)
		assert.Nil(t, err)
	})

	t.Run("WordWithAPronoun", func(t *testing.T) {
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := MelissosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}
		found, err := testHandler.queryWord(body)
		assert.True(t, found)
		assert.Nil(t, err)
	})

	t.Run("WordFoundButDifferentMeaning", func(t *testing.T) {
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := MelissosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		body.English = "notthesame"
		found, err := testHandler.queryWord(body)
		assert.False(t, found)
		assert.Nil(t, err)
	})

	t.Run("WordFoundDifferentMeaningWithoutAPronoun", func(t *testing.T) {
		file := "thalesSingleHit"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := MelissosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		body.English = "notthesame but multiple words"

		found, err := testHandler.queryWord(body)
		assert.False(t, found)
		assert.Nil(t, err)
	})

	t.Run("DoesNotExist", func(t *testing.T) {
		file := "searchWordNoResults"
		status := 200
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := MelissosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		found, err := testHandler.queryWord(body)
		assert.False(t, found)
		assert.Nil(t, err)
	})

	t.Run("DoesNotExist", func(t *testing.T) {
		file := "shardFailure"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := MelissosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

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

		testHandler := MelissosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		testHandler.addWord(body)
		assert.Equal(t, testHandler.Created, 0)
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

		testHandler := MelissosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		wait.Add(1)

		testHandler.transformWord(body, &wait)
		assert.Equal(t, testHandler.Created, 1)
	})

	t.Run("DocumentNotCreated", func(t *testing.T) {
		file := "shardFailure"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		var wait sync.WaitGroup

		testHandler := MelissosHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		wait.Add(1)

		testHandler.transformWord(body, &wait)
		assert.Equal(t, testHandler.Created, 0)
	})

}

func TestHandlerParser(t *testing.T) {
	splitWord := "ἀκούω + gen."
	pronounSplit := "μῦθος, ὁ"
	pronounSplitTwo := "ὁ δοῦλος"
	testHandler := MelissosHandler{
		Elastic: nil,
		Index:   "",
		Created: 0,
	}

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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	t.Run("JobFinished", func(t *testing.T) {
		testClient := thales.NewFakeKubeClient()

		jobSpec := thales.TestJobObject(expectedName, ns, true)
		job, err := testClient.BatchV1().Jobs(ns).Create(ctx, jobSpec, metav1.CreateOptions{})
		assert.Nil(t, err)
		assert.Equal(t, job.Name, expectedName)

		testHandler := MelissosHandler{
			Kube:         testClient,
			Job:          expectedName,
			Namespace:    ns,
			Duration:     duration,
			TimeFinished: timeFinished,
		}

		jobExit := make(chan bool, 1)
		go testHandler.WaitForJobsToFinish(jobExit)

		select {

		case <-jobExit:
			exitStatus := <-jobExit
			assert.True(t, exitStatus)
		}
	})

	t.Run("JobNotFinished", func(t *testing.T) {
		testClient := thales.NewFakeKubeClient()

		jobSpec := thales.TestJobObject(expectedName, ns, false)
		job, err := testClient.BatchV1().Jobs(ns).Create(ctx, jobSpec, metav1.CreateOptions{})
		assert.Nil(t, err)
		assert.Equal(t, job.Name, expectedName)

		testHandler := MelissosHandler{
			Duration:  duration,
			Kube:      testClient,
			Job:       expectedName,
			Namespace: ns,
		}

		timeFinished = duration.Milliseconds() * 2

		jobExit := make(chan bool, 1)
		go testHandler.WaitForJobsToFinish(jobExit)

		select {

		case <-jobExit:
			exitStatus := <-jobExit
			assert.False(t, exitStatus)
		}
	})
}
