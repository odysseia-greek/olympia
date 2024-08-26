package monos

import (
	"context"
	"encoding/json"
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/transform"
	"github.com/odysseia-greek/agora/thales"
	ptolemaios "github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"sync"
	"time"
)

type MelissosHandler struct {
	Duration     time.Duration
	TimeFinished int64
	Index        string
	Created      int
	Updated      int
	Processed    int
	Elastic      elastic.Client
	Eupalinos    EupalinosClient
	Channel      string
	DutchChannel string
	WaitTime     time.Duration
	Kube         *thales.KubeClient
	Namespace    string
	Job          string
	Ambassador   *ptolemaios.ClientAmbassador
}

func (m *MelissosHandler) HandleParmenides() bool {
	finished := false

	//handle Parmenides channel
	for {
		payload := &pb.ChannelInfo{Name: m.Channel}
		msg, err := m.Eupalinos.DequeueMessage(context.Background(), payload)
		if err != nil {
			logging.Info("Queue is empty. Waiting...")
			time.Sleep(m.WaitTime)
			queueLength, _ := m.Eupalinos.GetQueueLength(context.Background(), payload)
			if queueLength.Length == 0 {
				finished = true
				break
			}

			continue
		}

		var word models.Meros
		err = json.Unmarshal([]byte(msg.Data), &word)
		if err != nil {
			logging.Error(err.Error())
		}

		m.Processed++

		found, err := m.queryWord(word)
		if err != nil {
			logging.Error(err.Error())
		}

		if !found {
			m.addWord(word)
		}
	}

	return finished
}

func (m *MelissosHandler) HandleDutch() bool {
	finished := false
	//handle Dutch channel
	for {
		payload := &pb.ChannelInfo{Name: m.DutchChannel}
		msg, err := m.Eupalinos.DequeueMessage(context.Background(), payload)
		if err != nil {
			logging.Info("Queue is empty. Waiting...")
			time.Sleep(m.WaitTime)
			queueLength, _ := m.Eupalinos.GetQueueLength(context.Background(), payload)
			if queueLength.Length == 0 {
				finished = true
				break
			}

			continue
		}

		var word models.Meros
		err = json.Unmarshal([]byte(msg.Data), &word)
		if err != nil {
			logging.Error(err.Error())
		}

		err = m.addDutchWord(word)
		if err != nil {
			logging.Error(err.Error())
		}

		m.Processed++
	}

	return finished
}

func (m *MelissosHandler) addDutchWord(word models.Meros) error {
	if word.Dutch == "de, het" {
		return nil
	}
	s := m.stripMouseionWords(word.Greek)
	strippedWord := transform.RemoveAccents(s)

	term := "greek"
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"prefix": map[string]interface{}{
							fmt.Sprintf("%s.keyword", term): fmt.Sprintf("%s,", strippedWord),
						},
					},
					map[string]interface{}{
						"prefix": map[string]interface{}{
							fmt.Sprintf("%s.keyword", term): fmt.Sprintf("%s ", strippedWord),
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							fmt.Sprintf("%s.keyword", term): strippedWord,
						},
					},
				},
			},
		},
		"size": 100,
	}
	response, err := m.Elastic.Query().MatchWithScroll(m.Index, query)

	if err != nil {
		return err
	}

	for _, hit := range response.Hits.Hits {
		jsonHit, _ := json.Marshal(hit.Source)
		meros, _ := models.UnmarshalMeros(jsonHit)
		baseWord := m.extractBaseWord(meros.Greek)
		if len(response.Hits.Hits) > 1 {
			if baseWord != strippedWord && baseWord != s {
				continue
			}
		}
		meros.Dutch = word.Dutch
		jsonifiedLogos, _ := meros.Marshal()
		_, err := m.Elastic.Document().Update(m.Index, hit.ID, jsonifiedLogos)
		if err != nil {
			return err
		}

		m.Updated++
	}

	return nil
}

func (m *MelissosHandler) queryWord(word models.Meros) (bool, error) {
	found := false

	strippedWord := transform.RemoveAccents(word.Greek)

	term := "greek"
	query := m.Elastic.Builder().MatchQuery(term, strippedWord)
	response, err := m.Elastic.Query().Match(m.Index, query)

	if err != nil {
		logging.Error(err.Error())
		return false, err
	}

	if len(response.Hits.Hits) >= 1 {
		found = true
	}

	var parsedEnglishWord string
	pronouns := []string{"a", "an", "the"}
	splitEnglish := strings.Split(word.English, " ")
	numberOfWords := len(splitEnglish)
	if numberOfWords > 1 {
		for _, pronoun := range pronouns {
			if splitEnglish[0] == pronoun {
				toJoin := splitEnglish[1:numberOfWords]
				parsedEnglishWord = strings.Join(toJoin, " ")
				break
			} else {
				parsedEnglishWord = word.English
			}
		}
	} else {
		parsedEnglishWord = word.English
	}

	for _, hit := range response.Hits.Hits {
		jsonHit, _ := json.Marshal(hit.Source)
		meros, _ := models.UnmarshalMeros(jsonHit)
		if meros.English == parsedEnglishWord || meros.English == word.English {
			return true, nil
		} else {
			found = false
		}
	}

	return found, nil
}

func (m *MelissosHandler) addWord(word models.Meros) {
	var innerWaitGroup sync.WaitGroup
	jsonifiedLogos, _ := word.Marshal()
	_, err := m.Elastic.Index().CreateDocument(m.Index, jsonifiedLogos)

	if err != nil {
		logging.Error(err.Error())
		return
	} else {
		innerWaitGroup.Add(1)
		go m.transformWord(word, &innerWaitGroup)
	}
}

func (m *MelissosHandler) transformWord(word models.Meros, wg *sync.WaitGroup) {
	defer wg.Done()
	strippedWord := transform.RemoveAccents(word.Greek)
	meros := models.Meros{
		Greek:      strippedWord,
		English:    word.English,
		Dutch:      word.Dutch,
		LinkedWord: word.LinkedWord,
		Original:   word.Greek,
	}

	jsonifiedLogos, _ := meros.Marshal()
	_, err := m.Elastic.Index().CreateDocument(m.Index, jsonifiedLogos)

	if err != nil {
		logging.Error(err.Error())
		return
	}

	m.Created++

	return
}

func (m *MelissosHandler) stripMouseionWords(word string) string {
	if !strings.Contains(word, " ") {
		return word
	}

	splitKeys := []string{"+", "("}
	greekPronous := []string{"ἡ", "ὁ", "τὸ", "τό"}
	var w string

	for _, key := range splitKeys {
		if strings.Contains(word, key) {
			split := strings.Split(word, key)
			w = strings.TrimSpace(split[0])
		}
	}

	for _, pronoun := range greekPronous {
		if strings.Contains(word, pronoun) {
			if strings.Contains(word, ",") {
				split := strings.Split(word, ",")
				w = strings.TrimSpace(split[0])
			} else {
				split := strings.Split(word, pronoun)
				w = strings.TrimSpace(split[1])
			}
		}
	}

	if w == "" {
		return word
	}

	return w
}

func (m *MelissosHandler) PrintProgress() {
	for {
		logging.Info(fmt.Sprintf("documents processed: %d | documents created: %d | documents updated: %d", m.Processed, m.Created, m.Updated))
		time.Sleep(20 * time.Second)
	}
}

func (m *MelissosHandler) WaitForJobsToFinish(c chan bool) {
	start := time.Now()
	ticker := time.NewTicker(m.Duration)
	defer ticker.Stop()

	for ts := range ticker.C {
		if ts.Sub(start).Milliseconds() >= m.TimeFinished {
			c <- false
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		job, err := m.Kube.BatchV1().Jobs(m.Namespace).Get(ctx, m.Job, metav1.GetOptions{})
		if err != nil {
			logging.Error(err.Error())
		}

		conditionFound := false
		if job.Status.Active == 0 {
			for _, condition := range job.Status.Conditions {
				if condition.Type == "Complete" {
					conditionFound = true
				}
			}
		}

		if conditionFound {
			c <- true
		}
	}
}

func (m *MelissosHandler) extractBaseWord(queryWord string) string {
	// Normalize and split the input
	strippedWord := transform.RemoveAccents(strings.ToLower(queryWord))
	splitWord := strings.Split(strippedWord, " ")

	// Known Greek pronouns
	greekPronouns := map[string]bool{"η": true, "ο": true, "το": true}

	// Function to clean punctuation from a word
	cleanWord := func(word string) string {
		return strings.Trim(word, ",.!?-") // Add any other punctuation as needed
	}

	// Iterate through the words
	for _, word := range splitWord {
		cleanedWord := cleanWord(word)

		if strings.HasPrefix(cleanedWord, "-") {
			// Skip words starting with "-"
			continue
		}

		if _, isPronoun := greekPronouns[cleanedWord]; !isPronoun {
			// If the word is not a pronoun, it's likely the correct word
			return cleanedWord
		}
	}

	return queryWord
}
