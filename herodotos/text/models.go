package text

import (
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
)

func textAggregationQuery() map[string]interface{} {
	return map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"authors": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "author.keyword",
					"size":  100,
				},
				"aggs": map[string]interface{}{
					"books": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "book.keyword",
							"size":  100,
						},
						"aggs": map[string]interface{}{
							"references": map[string]interface{}{
								"terms": map[string]interface{}{
									"field": "reference",
									"size":  100,
								},
								"aggs": map[string]interface{}{
									"sections": map[string]interface{}{
										"nested": map[string]interface{}{
											"path": "rhemai",
										},
										"aggs": map[string]interface{}{
											"section_ids": map[string]interface{}{
												"terms": map[string]interface{}{
													"field": "rhemai.section",
													"size":  100,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func parseAggregationResults(agg map[string]interface{}) (models.AggregationResult, error) {
	var result models.AggregationResult

	authorsAgg, ok := agg["aggregations"].(map[string]interface{})["authors"].(map[string]interface{})["buckets"].([]interface{})
	if !ok {
		return result, fmt.Errorf("error parsing authors aggregation")
	}

	for _, authorBucket := range authorsAgg {
		authorMap := authorBucket.(map[string]interface{})
		author := models.ESAuthor{
			Key: authorMap["key"].(string),
		}

		booksAgg, ok := authorMap["books"].(map[string]interface{})["buckets"].([]interface{})
		if !ok {
			return result, fmt.Errorf("error parsing books aggregation")
		}

		for _, bookBucket := range booksAgg {
			bookMap := bookBucket.(map[string]interface{})
			book := models.ESBook{
				Key: bookMap["key"].(string),
			}

			referencesAgg, ok := bookMap["references"].(map[string]interface{})["buckets"].([]interface{})
			if !ok {
				return result, fmt.Errorf("error parsing references aggregation")
			}

			for _, referenceBucket := range referencesAgg {
				referenceMap := referenceBucket.(map[string]interface{})
				reference := models.Reference{
					Key: referenceMap["key"].(string),
				}

				sectionsAgg, ok := referenceMap["sections"].(map[string]interface{})["section_ids"].(map[string]interface{})["buckets"].([]interface{})
				if !ok {
					return result, fmt.Errorf("error parsing sections aggregation")
				}

				for _, sectionBucket := range sectionsAgg {
					sectionMap := sectionBucket.(map[string]interface{})
					section := models.Section{
						Key: sectionMap["key"].(string),
					}
					reference.Sections = append(reference.Sections, section)
				}

				book.References = append(book.References, reference)
			}

			author.Books = append(author.Books, book)
		}

		result.Authors = append(result.Authors, author)
	}

	return result, nil
}

func createGreekTextQuery(words []string) map[string]interface{} {
	shouldClauses := make([]map[string]interface{}, len(words))

	for i, word := range words {
		shouldClauses[i] = map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "rhemai",
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"should": []map[string]interface{}{
							{
								"match": map[string]interface{}{
									"rhemai.greek": word,
								},
							},
						},
					},
				},
			},
		}
	}

	boolQuery := map[string]interface{}{
		"bool": map[string]interface{}{
			"should": shouldClauses,
		},
	}

	return map[string]interface{}{
		"query": boolQuery,
	}
}
