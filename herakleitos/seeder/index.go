package seeder

func textIndex(policyName string) map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 1,
				"lifecycle.name":     policyName,
			},
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"greek_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"greek_stop",
							"greek_stemmer",
						},
					},
				},
				"filter": map[string]interface{}{
					"greek_stop": map[string]interface{}{
						"type":      "stop",
						"stopwords": "_greek_",
					},
					"greek_stemmer": map[string]interface{}{
						"type":     "stemmer",
						"language": "greek",
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"author": map[string]interface{}{
					"type": "keyword",
				},
				"book": map[string]interface{}{
					"type": "keyword", // Assuming book is a string now
				},
				"type": map[string]interface{}{
					"type": "keyword",
				},
				"reference": map[string]interface{}{
					"type": "keyword",
				},
				"perseusTextLink": map[string]interface{}{
					"type": "keyword",
				},

				"rhemai": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"greek": map[string]interface{}{
							"type":     "text",
							"analyzer": "greek_analyzer",
						},
						"translations": map[string]interface{}{
							"type": "text",
						},
						"section": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
			},
		},
	}
}
