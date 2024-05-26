package text

import (
	"fmt"
)

func textAggregationQuery() map[string]interface{} {
	return map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"authors": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "author",
					"size":  100,
				},
				"aggs": map[string]interface{}{
					"books": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "book",
							"size":  100,
						},
						"aggs": map[string]interface{}{
							"references": map[string]interface{}{
								"nested": map[string]interface{}{
									"path": "biblos",
								},
								"aggs": map[string]interface{}{
									"reference_ids": map[string]interface{}{
										"terms": map[string]interface{}{
											"field": "biblos.reference",
											"size":  100,
										},
										"aggs": map[string]interface{}{
											"sections": map[string]interface{}{
												"nested": map[string]interface{}{
													"path": "biblos.rhemai",
												},
												"aggs": map[string]interface{}{
													"section_ids": map[string]interface{}{
														"terms": map[string]interface{}{
															"field": "biblos.rhemai.section",
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
			},
		},
	}
}

func parseAggregationResults(agg map[string]interface{}) (AggregationResult, error) {
	var result AggregationResult

	authorsAgg, ok := agg["aggregations"].(map[string]interface{})["authors"].(map[string]interface{})["buckets"].([]interface{})
	if !ok {
		return result, fmt.Errorf("error parsing authors aggregation")
	}

	for _, authorBucket := range authorsAgg {
		authorMap := authorBucket.(map[string]interface{})
		author := ESAuthor{
			Key: authorMap["key"].(string),
		}

		booksAgg, ok := authorMap["books"].(map[string]interface{})["buckets"].([]interface{})
		if !ok {
			return result, fmt.Errorf("error parsing books aggregation")
		}

		for _, bookBucket := range booksAgg {
			bookMap := bookBucket.(map[string]interface{})
			book := ESBook{
				Key: bookMap["key"].(string),
			}

			referencesAgg, ok := bookMap["references"].(map[string]interface{})["reference_ids"].(map[string]interface{})["buckets"].([]interface{})
			if !ok {
				return result, fmt.Errorf("error parsing references aggregation")
			}

			for _, referenceBucket := range referencesAgg {
				referenceMap := referenceBucket.(map[string]interface{})
				reference := Reference{
					Key: referenceMap["key"].(string),
				}

				sectionsAgg, ok := referenceMap["sections"].(map[string]interface{})["section_ids"].(map[string]interface{})["buckets"].([]interface{})
				if !ok {
					return result, fmt.Errorf("error parsing sections aggregation")
				}

				for _, sectionBucket := range sectionsAgg {
					sectionMap := sectionBucket.(map[string]interface{})
					section := Section{
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

type Section struct {
	Key string `json:"key"`
}

type Reference struct {
	Key      string    `json:"key"`
	Sections []Section `json:"sections"`
}

type ESBook struct {
	Key        string      `json:"key"`
	References []Reference `json:"references"`
}

type ESAuthor struct {
	Key   string   `json:"key"`
	Books []ESBook `json:"books"`
}

type AggregationResult struct {
	Authors []ESAuthor `json:"authors"`
}
