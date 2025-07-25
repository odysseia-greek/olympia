package pneuma

import (
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/delphi/aristides/diplomat"
	"strings"
)

type AnaximenesHandler struct {
	Indices    []string
	MaxAge     string
	PolicyName string
	Elastic    elastic.Client
	Ambassador *diplomat.ClientAmbassador
}

func (a *AnaximenesHandler) CreateAttikeIndices() {
	for _, index := range a.Indices {
		deleted, err := a.Elastic.Index().Delete(index)

		logging.Info(fmt.Sprintf("delete response for %s: success=%v, err=%v",
			index, deleted, err))

		if err != nil {
			if strings.Contains(err.Error(), "index_not_found_exception") {
				logging.Info(fmt.Sprintf("index %s not found, creating it", index))
				err = a.createIndexAtStartup(index)
				if err != nil {
					logging.Error(fmt.Sprintf("failed to create index: %v", err))
				}
				continue
			}

			// Log any other error in detail
			logging.Debug(fmt.Sprintf("error deleting index %s: %v", index, err))
		}

		if deleted {
			logging.Info(fmt.Sprintf("recreating index %s after deletion", index))
			err = a.createIndexAtStartup(index)
			if err != nil {
				logging.Error(fmt.Sprintf("failed to recreate index: %v", err))
			}
		}
	}
}

func (a *AnaximenesHandler) createIndexAtStartup(index string) error {
	request := a.createMapping(index, a.PolicyName)
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
