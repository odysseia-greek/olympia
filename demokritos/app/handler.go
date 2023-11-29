package app

import (
	"bytes"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/transform"
	configs "github.com/odysseia-greek/olympia/demokritos/config"
	"strings"
	"sync"
)

type DemokritosHandler struct {
	Config *configs.Config
}

func (d *DemokritosHandler) DeleteIndexAtStartUp() error {
	deleted, err := d.Config.Elastic.Index().Delete(d.Config.Index)
	logging.Info(fmt.Sprintf("deleted index: %s success: %v", d.Config.Index, deleted))
	if err != nil {
		if deleted {
			return nil
		}
		if strings.Contains(err.Error(), "index_not_found_exception") {
			logging.Error(err.Error())
			return nil
		}

		return err
	}

	return nil
}

func (d *DemokritosHandler) CreateIndexAtStartup() error {
	logging.Info(fmt.Sprintf("creating policy: %s", d.Config.PolicyName))
	err := d.createPolicyAtStartup()
	if err != nil {
		return err
	}
	logging.Info(fmt.Sprintf("creating index: %s with min: %v and max: %v ngram", d.Config.Index, d.Config.MinNGram, d.Config.MaxNGram))
	query := d.Config.Elastic.Builder().DictionaryIndex(d.Config.MinNGram, d.Config.MaxNGram, d.Config.PolicyName)
	res, err := d.Config.Elastic.Index().Create(d.Config.Index, query)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s", res.Index))
	return nil
}

func (d *DemokritosHandler) createPolicyAtStartup() error {
	policyCreated, err := d.Config.Elastic.Policy().CreateHotPolicy(d.Config.PolicyName)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", d.Config.PolicyName, policyCreated.Acknowledged))

	return nil
}

func (d *DemokritosHandler) AddDirectoryToElastic(biblos models.Biblos, wg *sync.WaitGroup) {
	defer wg.Done()
	var buf bytes.Buffer

	var currBatch int

	for _, word := range biblos.Biblos {
		currBatch++

		meta := []byte(fmt.Sprintf(`{ "index": {} }%s`, "\n"))
		jsonifiedWord, _ := word.Marshal()
		jsonifiedWord = append(jsonifiedWord, "\n"...)
		buf.Grow(len(meta) + len(jsonifiedWord))
		buf.Write(meta)
		buf.Write(jsonifiedWord)

		stripped := d.transformWord(word)
		stripped = append(stripped, "\n"...)
		buf.Grow(len(meta) + len(stripped))
		buf.Write(meta)
		buf.Write(stripped)

		if currBatch == len(biblos.Biblos) {
			res, err := d.Config.Elastic.Document().Bulk(buf, d.Config.Index)
			if err != nil {
				logging.Error(err.Error())
				return
			}

			d.Config.Created = d.Config.Created + len(res.Items)
		}
	}
}

func (d *DemokritosHandler) transformWord(m models.Meros) []byte {
	strippedWord := transform.RemoveAccents(m.Greek)
	word := models.Meros{
		Greek:      strippedWord,
		English:    m.English,
		LinkedWord: m.LinkedWord,
		Original:   m.Greek,
	}

	jsonifiedWord, _ := word.Marshal()

	return jsonifiedWord
}
