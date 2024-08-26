package apeiron

import (
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	ptolemaios "github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	"strings"
	"time"
)

type AnaximanderHandler struct {
	Index      string
	Created    int
	PolicyName string
	Elastic    elastic.Client
	Ambassador *ptolemaios.ClientAmbassador
}

func (a *AnaximanderHandler) DeleteIndexAtStartUp() error {
	deleted, err := a.Elastic.Index().Delete(a.Index)
	logging.Info(fmt.Sprintf("deleted index: %s success: %v", a.Index, deleted))
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

func (a *AnaximanderHandler) CreateIndexAtStartup() error {
	logging.Info(fmt.Sprintf("creating policy: %s", a.PolicyName))
	err := a.createPolicyAtStartup()
	if err != nil {
		return err
	}

	indexMapping := a.Elastic.Builder().GrammarIndex(a.PolicyName)
	created, err := a.Elastic.Index().Create(a.Index, indexMapping)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", a.Index, created.Acknowledged))

	return nil
}

func (a *AnaximanderHandler) createPolicyAtStartup() error {
	policyCreated, err := a.Elastic.Policy().CreateWarmPolicy(a.PolicyName)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", a.PolicyName, policyCreated.Acknowledged))

	return nil
}

func (a *AnaximanderHandler) AddToElastic(declension models.Declension) error {
	upload, _ := declension.Marshal()

	_, err := a.Elastic.Index().CreateDocument(a.Index, upload)
	a.Created++
	if err != nil {
		return err
	}

	return nil
}

func (a *AnaximanderHandler) PrintProgress(total int) {
	for {
		percentage := float64(a.Created) / float64(total) * 100
		logging.Info(fmt.Sprintf("Progress: %d/%d documents created (%.2f%%)", a.Created, total, percentage))
		time.Sleep(1000 * time.Second)
	}
}
