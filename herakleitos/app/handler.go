package app

import (
	"bytes"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	configs "github.com/odysseia-greek/olympia/herakleitos/config"
	"strings"
	"sync"
)

type HerakleitosHandler struct {
	Config *configs.Config
}

func (h *HerakleitosHandler) DeleteIndexAtStartUp() error {
	deleted, err := h.Config.Elastic.Index().Delete(h.Config.Index)
	logging.Info(fmt.Sprintf("deleted index: %s success: %v", h.Config.Index, deleted))
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

func (h *HerakleitosHandler) createPolicyAtStartup() error {
	policyCreated, err := h.Config.Elastic.Policy().CreateWarmPolicy(h.Config.PolicyName)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", h.Config.PolicyName, policyCreated.Acknowledged))

	return nil
}

func (h *HerakleitosHandler) CreateIndexAtStartup() error {
	logging.Info(fmt.Sprintf("creating policy: %s", h.Config.PolicyName))
	err := h.createPolicyAtStartup()
	if err != nil {
		return err
	}

	indexMapping := h.Config.Elastic.Builder().TextIndex(h.Config.PolicyName)
	created, err := h.Config.Elastic.Index().Create(h.Config.Index, indexMapping)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", created.Index, created.Acknowledged))

	return nil
}

func (h *HerakleitosHandler) Add(rhema models.Rhema, wg *sync.WaitGroup) error {
	defer wg.Done()
	var buf bytes.Buffer

	var currBatch int

	for _, text := range rhema.Rhemai {
		currBatch++

		meta := []byte(fmt.Sprintf(`{ "index": {} }%s`, "\n"))
		jsonifiedText, _ := text.Marshal()
		jsonifiedText = append(jsonifiedText, "\n"...)
		buf.Grow(len(meta) + len(jsonifiedText))
		buf.Write(meta)
		buf.Write(jsonifiedText)

		if currBatch == len(rhema.Rhemai) {
			res, err := h.Config.Elastic.Document().Bulk(buf, h.Config.Index)
			if err != nil {
				logging.Error(err.Error())
				return err
			}

			h.Config.Created = h.Config.Created + len(res.Items)
		}
	}

	return nil
}
