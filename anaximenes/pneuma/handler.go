package pneuma

import (
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	"strings"
)

type AnaximenesHandler struct {
	Indices    []string
	MaxAge     string
	Elastic    elastic.Client
	Ambassador *diplomat.ClientAmbassador
}

func (a *AnaximenesHandler) CreateAttikeIndices() {
	for _, index := range a.Indices {
		deleted, err := a.Elastic.Index().Delete(index)
		logging.Info(fmt.Sprintf("deleted index: %s success: %v", index, deleted))
		if err != nil {
			logging.Error(err.Error())
			if strings.Contains(err.Error(), "index_not_found_exception") {
				err = a.createIndexAtStartup(index)
				if err != nil {
					logging.Error(err.Error())
				}
				continue
			}
			logging.Debug(fmt.Sprintf("cannot delete index: %s which means an aliased version exist and should not be deleted", index))
		}

		if deleted {
			err = a.createIndexAtStartup(index)
			if err != nil {
				logging.Error(err.Error())
			}
			continue
		}
	}
}

func (a *AnaximenesHandler) createPolicyAtStartup(policyName string) error {
	policyCreated, err := a.Elastic.Policy().CreatePolicyWithRollOver(policyName, a.MaxAge, "hot")
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", policyName, policyCreated.Acknowledged))

	return nil
}

func (a *AnaximenesHandler) createIndexAtStartup(index string) error {
	policyName := fmt.Sprintf("%s_policy", index)
	err := a.createPolicyAtStartup(policyName)
	if err != nil {
		return err
	}

	request := a.createMapping(index, policyName)
	created, err := a.Elastic.Index().CreateWithAlias(index, request)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", index, created.Acknowledged))

	return nil
}

func (a *AnaximenesHandler) createMapping(indexName, policyName string) map[string]interface{} {
	switch indexName {
	case config.TracingElasticIndex:
		return a.createTraceIndexMapping(policyName)
	case config.MetricsElasticIndex:
		return a.createMetricsIndexMapping(policyName)
	}

	return nil
}

func (a *AnaximenesHandler) createTraceIndexMapping(policyName string) map[string]interface{} {
	return map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"items": map[string]interface{}{
					"type": "nested", // Use "nested" type for arrays of complex objects
				},
				"isActive": map[string]interface{}{
					"type": "boolean",
				},
				"timeStarted": map[string]interface{}{
					"type":   "date",
					"format": "yyyy-MM-dd'T'HH:mm:ss.SSS",
				},
				"timeEnded": map[string]interface{}{
					"type":       "date",
					"format":     "yyyy-MM-dd'T'HH:mm:ss.SSS",
					"null_value": "1970-01-01T00:00:00.000",
				},
				"totalTime": map[string]interface{}{
					"type": "long",
				},
				"responseCode": map[string]interface{}{
					"type": "short",
				},
				"metrics": map[string]interface{}{
					"properties": map[string]interface{}{
						"cpu_units":             map[string]interface{}{"type": "keyword"},
						"memory_units":          map[string]interface{}{"type": "keyword"},
						"name":                  map[string]interface{}{"type": "keyword"},
						"cpu_raw":               map[string]interface{}{"type": "integer"},
						"memory_raw":            map[string]interface{}{"type": "integer"},
						"cpu_human_readable":    map[string]interface{}{"type": "keyword"},
						"memory_human_readable": map[string]interface{}{"type": "keyword"},
					},
				},
				// Add additional fields here if needed
			},
		},
		"settings": map[string]interface{}{
			"index.lifecycle.name":                   policyName,
			"index.lifecycle.rollover_alias":         "trace",
			"index.lifecycle.parse_origination_date": true,
		},
	}
}

func (a *AnaximenesHandler) createMetricsIndexMapping(policyName string) map[string]interface{} {
	return map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"metrics": map[string]interface{}{
					"properties": map[string]interface{}{
						"cpu_units":    map[string]interface{}{"type": "keyword"},
						"memory_units": map[string]interface{}{"type": "keyword"},
						"pods":         map[string]interface{}{"type": "nested", "properties": podProperties()},
						"nodes":        map[string]interface{}{"type": "nested", "properties": nodeProperties()},
						"grouped":      map[string]interface{}{"type": "nested", "properties": groupedProperties()},
					},
				},
				"timeStamp": map[string]interface{}{
					"type":   "date",
					"format": "yyyy-MM-dd'T'HH:mm:ss.SSS",
				},
			},
		},
		"settings": map[string]interface{}{
			"index.lifecycle.name":                   policyName,
			"index.lifecycle.rollover_alias":         "metrics",
			"index.lifecycle.parse_origination_date": true,
		},
	}
}

func podProperties() map[string]interface{} {
	return map[string]interface{}{
		"name":                  map[string]interface{}{"type": "keyword"},
		"cpu_raw":               map[string]interface{}{"type": "integer"},
		"memory_raw":            map[string]interface{}{"type": "integer"},
		"cpu_human_readable":    map[string]interface{}{"type": "keyword"},
		"memory_human_readable": map[string]interface{}{"type": "keyword"},
	}
}

func nodeProperties() map[string]interface{} {
	return map[string]interface{}{
		"node_name":                        map[string]interface{}{"type": "keyword"},
		"cpu_raw":                          map[string]interface{}{"type": "integer"},
		"memory_raw":                       map[string]interface{}{"type": "integer"},
		"cpu_percentage":                   map[string]interface{}{"type": "float"},
		"memory_percentage":                map[string]interface{}{"type": "float"},
		"cpu_human_readable":               map[string]interface{}{"type": "keyword"},
		"memory_human_readable":            map[string]interface{}{"type": "keyword"},
		"cpu_percentage_human_readable":    map[string]interface{}{"type": "keyword"},
		"memory_percentage_human_readable": map[string]interface{}{"type": "keyword"},
	}
}

func groupedProperties() map[string]interface{} {
	return map[string]interface{}{
		"name":                  map[string]interface{}{"type": "keyword"},
		"cpu_raw":               map[string]interface{}{"type": "integer"},
		"memory_raw":            map[string]interface{}{"type": "integer"},
		"cpu_human_readable":    map[string]interface{}{"type": "keyword"},
		"memory_human_readable": map[string]interface{}{"type": "keyword"},
	}
}
