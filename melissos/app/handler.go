package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/transform"
	pb "github.com/odysseia-greek/eupalinos/proto"
	configs "github.com/odysseia-greek/olympia/melissos/config"
	"strings"
	"sync"
	"time"
)

type MelissosHandler struct {
	Config *configs.Config
}

func (m *MelissosHandler) HandleParmenides() bool {
	finished := false

	//handle Parmenides channel
	for {
		payload := &pb.ChannelInfo{Name: m.Config.Channel}
		msg, err := m.Config.Eupalinos.DequeueMessage(context.Background(), payload)
		if err != nil {
			logging.Info("Queue is empty. Waiting...")
			time.Sleep(m.Config.WaitTime)
			queueLength, _ := m.Config.Eupalinos.GetQueueLength(context.Background(), payload)
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

		m.Config.Processed++

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
		payload := &pb.ChannelInfo{Name: m.Config.DutchChannel}
		msg, err := m.Config.Eupalinos.DequeueMessage(context.Background(), payload)
		if err != nil {
			logging.Info("Queue is empty. Waiting...")
			time.Sleep(m.Config.WaitTime)
			queueLength, _ := m.Config.Eupalinos.GetQueueLength(context.Background(), payload)
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

		m.Config.Processed++
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
	query := m.Config.Elastic.Builder().MatchQuery(term, strippedWord)
	response, err := m.Config.Elastic.Query().Match(m.Config.Index, query)

	if err != nil {
		return err
	}

	for _, hit := range response.Hits.Hits {
		jsonHit, _ := json.Marshal(hit.Source)
		meros, _ := models.UnmarshalMeros(jsonHit)
		if len(response.Hits.Hits) > 1 {
			if meros.Greek != strippedWord && meros.Greek != s {
				continue
			}
		}
		meros.Dutch = word.Dutch
		jsonifiedLogos, _ := meros.Marshal()
		_, err := m.Config.Elastic.Document().Update(m.Config.Index, hit.ID, jsonifiedLogos)
		if err != nil {
			return err
		}

		m.Config.Updated++
	}

	return nil
}

func (m *MelissosHandler) queryWord(word models.Meros) (bool, error) {
	found := false

	strippedWord := transform.RemoveAccents(word.Greek)

	term := "greek"
	query := m.Config.Elastic.Builder().MatchQuery(term, strippedWord)
	response, err := m.Config.Elastic.Query().Match(m.Config.Index, query)

	if err != nil {
		logging.Error(err.Error())
		return found, err
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
	_, err := m.Config.Elastic.Index().CreateDocument(m.Config.Index, jsonifiedLogos)

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
	_, err := m.Config.Elastic.Index().CreateDocument(m.Config.Index, jsonifiedLogos)

	if err != nil {
		logging.Error(err.Error())
		return
	}

	m.Config.Created++

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
		logging.Info(fmt.Sprintf("documents processed: %d | documents created: %d | documents updated: %d", m.Config.Processed, m.Config.Created, m.Config.Updated))
		time.Sleep(20 * time.Second)
	}
}
