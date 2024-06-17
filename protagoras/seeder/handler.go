package seeder

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	"os"
	"regexp"
	"strings"
	"time"
)

type ProtagorasHandler struct {
	Save          bool
	wordsDone     []string
	wordsNotFound []string
	Client        service.OdysseiaClient
}

func (p *ProtagorasHandler) Start() error {
	return p.gatherOptions()
}

func (p *ProtagorasHandler) gatherOptions() error {
	response, err := p.Client.Herodotos().Options("")
	if err != nil {
		return err
	}

	var sectionsDone int
	var options models.AggregationResult
	err = json.NewDecoder(response.Body).Decode(&options)

	for _, author := range options.Authors {
		for _, book := range author.Books {
			for _, reference := range book.References {
				for _, section := range reference.Sections {
					err := p.queryText(author.Key, book.Key, reference.Key, section.Key)
					sectionsDone++
					time.Sleep(1 * time.Second)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	logging.Info(fmt.Sprintf("Total Number of sections done: %d", sectionsDone))
	logging.Info(fmt.Sprintf("Done a total of: %d words", len(p.wordsDone)+len(p.wordsNotFound)))
	logging.Info(fmt.Sprintf("Declinded a total of: %d words", len(p.wordsDone)))
	logging.Info(fmt.Sprintf("Could not decline a total of: %d words", len(p.wordsNotFound)))

	if p.Save {
		err = p.saveWords("/tmp/done.json", p.wordsDone)
		if err != nil {
			return err
		}

		err = p.saveWords("/tmp/notdone.json", p.wordsNotFound)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ProtagorasHandler) queryText(author, book, reference, section string) error {
	r := models.CreateTextRequest{
		Author:    author,
		Book:      book,
		Reference: reference,
		Section:   section,
	}

	jsonBody, err := json.Marshal(r)
	if err != nil {
		return err
	}

	response, err := p.Client.Herodotos().Create(jsonBody, "")
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to create a new text for %s", string(jsonBody)))
		return nil
	}

	var text models.Text
	err = json.NewDecoder(response.Body).Decode(&text)
	if err != nil {
		return err
	}

	return p.loopOverAndDeclineWords(text)
}

func (p *ProtagorasHandler) loopOverAndDeclineWords(text models.Text) error {
	if len(text.Rhemai) != 1 {
		logging.Error(fmt.Sprintf("expected 1 rhema, got %d", len(text.Rhemai)))
	}

	greekText := text.Rhemai[0].Greek

	greekWords := p.splitAndCleanGreekWords(greekText)

	for _, word := range greekWords {
		response, err := p.Client.Dionysios().Grammar(word, "")
		if err != nil {
			p.wordsNotFound = append(p.wordsNotFound, word)
			logging.Error(fmt.Sprintf("no result found for word: %s", word))
			continue
		}

		var declensions models.DeclensionTranslationResults
		err = json.NewDecoder(response.Body).Decode(&declensions)
		if err != nil {
			return err
		}

		if len(declensions.Results) > 0 {
			p.wordsDone = append(p.wordsDone, word)
			logging.Info(fmt.Sprintf("found declension result | word: %s | rootword: %s | translation: %s | rule: %s |", declensions.Results[0].Word, declensions.Results[0].RootWord, declensions.Results[0].Translation, declensions.Results[0].Rule))
		} else {
			p.wordsNotFound = append(p.wordsNotFound, word)
			logging.Error(fmt.Sprintf("no result found for word: %s", word))
		}
	}

	return nil
}

func (p *ProtagorasHandler) splitAndCleanGreekWords(text string) []string {
	// Remove punctuation using regex
	re := regexp.MustCompile(`[.,;:-·!?()“”‘’\[\]{}'"` + "`" + `]`)
	cleanText := re.ReplaceAllString(text, "")

	words := strings.Fields(cleanText)

	return words
}

func (p *ProtagorasHandler) saveWords(filename string, words []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(words)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
