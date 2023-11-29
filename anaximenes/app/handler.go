package app

import (
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	configs "github.com/odysseia-greek/olympia/anaximenes/config"
	"strings"
)

type AnaximenesConfig struct {
	Config *configs.Config
}

func (a *AnaximenesConfig) DeleteIndexAtStartUp() error {
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

func (a *AnaximenesConfig) createPolicyAtStartup() error {
	policyCreated, err := a.Config.Elastic.Policy().CreatePolicyWithRollOver(a.Config.PolicyName, a.Config.MaxAge, "hot")
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", a.Config.PolicyName, policyCreated.Acknowledged))

	return nil
}

func (a *AnaximenesConfig) CreateIndexAtStartup() error {
	logging.Info(fmt.Sprintf("creating policy: %s", a.Config.PolicyName))
	err := a.createPolicyAtStartup()
	if err != nil {
		return err
	}
	request := a.Config.Elastic.Builder().CreateTraceIndexMapping(a.Config.PolicyName)
	created, err := a.Config.Elastic.Index().CreateWithAlias(a.Config.Index, request)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", a.Config.Index, created.Acknowledged))

	return nil
}
