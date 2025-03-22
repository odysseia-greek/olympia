package flux

import (
	"bytes"
	"encoding/json"
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/delphi/aristides/diplomat"
	"strings"
	"sync"
)

type HerakleitosHandler struct {
	Index      string
	Created    int
	Elastic    elastic.Client
	PolicyName string
	Ambassador *diplomat.ClientAmbassador
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

	indexMapping := textIndex(h.PolicyName)
	created, err := h.Elastic.Index().Create(h.Index, indexMapping)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", created.Index, created.Acknowledged))

	return nil
}

func (h *HerakleitosHandler) Add(rhemai []Text, wg *sync.WaitGroup) error {
	defer wg.Done()
	var buf bytes.Buffer

	var currBatch int

	for _, rhema := range rhemai {
		currBatch++

		meta := []byte(fmt.Sprintf(`{ "index": {} }%s`, "\n"))
		jsonifiedText, _ := json.Marshal(rhema)
		jsonifiedText = append(jsonifiedText, "\n"...)
		buf.Grow(len(meta) + len(jsonifiedText))
		buf.Write(meta)
		buf.Write(jsonifiedText)

		if currBatch == len(rhemai) {
			res, err := h.Elastic.Document().Bulk(buf, h.Index)
			if err != nil {
				return err
			}

			h.Created = h.Created + len(res.Items)
		}
	}
	return nil
}
