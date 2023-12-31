package seeder

import (
	"bytes"
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	ptolemaios "github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	"strings"
	"sync"
)

type HerakleitosHandler struct {
	Index      string
	Created    int
	Elastic    elastic.Client
	PolicyName string
	Ambassador *ptolemaios.ClientAmbassador
}

func (h *HerakleitosHandler) DeleteIndexAtStartUp() error {
	deleted, err := h.Elastic.Index().Delete(h.Index)
	logging.Info(fmt.Sprintf("deleted index: %s success: %v", h.Index, deleted))
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
	policyCreated, err := h.Elastic.Policy().CreateWarmPolicy(h.PolicyName)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", h.PolicyName, policyCreated.Acknowledged))

	return nil
}

func (h *HerakleitosHandler) CreateIndexAtStartup() error {
	logging.Info(fmt.Sprintf("creating policy: %s", h.PolicyName))
	err := h.createPolicyAtStartup()
	if err != nil {
		return err
	}

	indexMapping := h.Elastic.Builder().TextIndex(h.PolicyName)
	created, err := h.Elastic.Index().Create(h.Index, indexMapping)
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
			res, err := h.Elastic.Document().Bulk(buf, h.Index)
			if err != nil {
				logging.Error(err.Error())
				return err
			}

			h.Created = h.Created + len(res.Items)
		}
	}

	return nil
}
