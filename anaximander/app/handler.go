package app

import (
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	configs "github.com/odysseia-greek/olympia/anaximander/config"
	"strings"
	"time"
)

type AnaximanderHandler struct {
	Config *configs.Config
}

func (a *AnaximanderHandler) DeleteIndexAtStartUp() error {
	deleted, err := a.Config.Elastic.Index().Delete(a.Config.Index)
	logging.Info(fmt.Sprintf("deleted index: %s success: %v", a.Config.Index, deleted))
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
	logging.Info(fmt.Sprintf("creating policy: %s", a.Config.PolicyName))
	err := a.createPolicyAtStartup()
	if err != nil {
		return err
	}

	indexMapping := a.Config.Elastic.Builder().GrammarIndex(a.Config.PolicyName)
	created, err := a.Config.Elastic.Index().Create(a.Config.Index, indexMapping)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", a.Config.Index, created.Acknowledged))

	return nil
}

func (a *AnaximanderHandler) createPolicyAtStartup() error {
	policyCreated, err := a.Config.Elastic.Policy().CreateWarmPolicy(a.Config.PolicyName)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", a.Config.PolicyName, policyCreated.Acknowledged))

	return nil
}

func (a *AnaximanderHandler) AddToElastic(declension models.Declension) error {
	upload, _ := declension.Marshal()

	_, err := a.Config.Elastic.Index().CreateDocument(a.Config.Index, upload)
	a.Config.Created++
	if err != nil {
		return err
	}

	return nil
}

func (a *AnaximanderHandler) PrintProgress(total int) {
	for {
		percentage := float64(a.Config.Created) / float64(total) * 100
		logging.Info(fmt.Sprintf("Progress: %d/%d documents created (%.2f%%)", a.Config.Created, total, percentage))
		time.Sleep(1000 * time.Second)
	}
}
