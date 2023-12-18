package seeder

import (
	"bytes"
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/transform"
	ptolemaios "github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	"strings"
	"sync"
)

type DemokritosHandler struct {
	Index      string
	SearchWord string
	Created    int
	Elastic    elastic.Client
	MinNGram   int
	MaxNGram   int
	PolicyName string
	Buf        bytes.Buffer
	Ambassador *ptolemaios.ClientAmbassador
}

func (d *DemokritosHandler) DeleteIndexAtStartUp() error {
	deleted, err := d.Elastic.Index().Delete(d.Index)
	logging.Info(fmt.Sprintf("deleted index: %s success: %v", d.Index, deleted))
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
	logging.Info(fmt.Sprintf("creating policy: %s", d.PolicyName))
	err := d.createPolicyAtStartup()
	if err != nil {
		return err
	}
	logging.Info(fmt.Sprintf("creating index: %s with min: %v and max: %v ngram", d.Index, d.MinNGram, d.MaxNGram))
	query := d.Elastic.Builder().DictionaryIndex(d.MinNGram, d.MaxNGram, d.PolicyName)
	res, err := d.Elastic.Index().Create(d.Index, query)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s", res.Index))
	return nil
}

func (d *DemokritosHandler) createPolicyAtStartup() error {
	policyCreated, err := d.Elastic.Policy().CreateHotPolicy(d.PolicyName)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", d.PolicyName, policyCreated.Acknowledged))

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
			res, err := d.Elastic.Document().Bulk(buf, d.Index)
			if err != nil {
				logging.Error(err.Error())
				return
			}

			d.Created = d.Created + len(res.Items)
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
