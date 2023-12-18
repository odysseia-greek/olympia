package seeder

import (
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	"strings"
)

type AnaximenesHandler struct {
	Index      string
	MaxAge     string
	PolicyName string
	Elastic    elastic.Client
	Ambassador *diplomat.ClientAmbassador
}

func (a *AnaximenesHandler) DeleteIndexAtStartUp() error {
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

func (a *AnaximenesHandler) createPolicyAtStartup() error {
	policyCreated, err := a.Elastic.Policy().CreatePolicyWithRollOver(a.PolicyName, a.MaxAge, "hot")
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", a.PolicyName, policyCreated.Acknowledged))

	return nil
}

func (a *AnaximenesHandler) CreateIndexAtStartup() error {
	logging.Info(fmt.Sprintf("creating policy: %s", a.PolicyName))
	err := a.createPolicyAtStartup()
	if err != nil {
		return err
	}
	request := a.Elastic.Builder().CreateTraceIndexMapping(a.PolicyName)
	created, err := a.Elastic.Index().CreateWithAlias(a.Index, request)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", a.Index, created.Acknowledged))

	return nil
}
