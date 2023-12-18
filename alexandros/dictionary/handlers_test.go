package dictionary

import (
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleanSearchResult(t *testing.T) {
	tests := []struct {
		name     string
		results  []models.Meros
		expected []models.Meros
	}{
		{
			name: "No duplicate entries",
			results: []models.Meros{
				{English: "apple", Greek: "μήλο"},
				{English: "banana", Greek: "μπανάνα"},
				{English: "carrot", Greek: "καρότο"},
				{English: "dog", Greek: "σκύλος"},
			},
			expected: []models.Meros{
				{English: "apple", Greek: "μήλο"},
				{English: "banana", Greek: "μπανάνα"},
				{English: "carrot", Greek: "καρότο"},
				{English: "dog", Greek: "σκύλος"},
			},
		},
		{
			name: "Duplicate entries but English is different",
			results: []models.Meros{
				{English: "carrot", Greek: "καρότο"},
				{English: "dog", Greek: "σκύλος"},
				{English: "hound", Greek: "σκύλος"},
			},
			expected: []models.Meros{
				{English: "carrot", Greek: "καρότο"},
				{English: "dog", Greek: "σκύλος"},
				{English: "hound", Greek: "σκύλος"},
			},
		},
		{
			name: "No Dutch set but including doubles",
			results: []models.Meros{
				{English: "holy", Greek: "ἱερός –ᾶ -ον"},
				{English: "revered, august, holy, awful", Greek: "σεμνός"},
				{English: "revered, august, holy, awful", Greek: "σεμνός"},
				{English: "revered, august, holy, awful", Greek: "σεμνός"},
			},
			expected: []models.Meros{
				{English: "holy", Greek: "ἱερός –ᾶ -ον"},
				{English: "revered, august, holy, awful", Greek: "σεμνός"},
			},
		},
		{
			name: "With duplicate entries, prefer non-empty Dutch",
			results: []models.Meros{
				{English: "holy", Greek: "ἱερός –ᾶ -ον", Dutch: ""},
				{English: "holy", Greek: "ἱερός –ᾶ -ον", Dutch: "heilig"},
				{English: "revered, august, holy, awful", Greek: "σεμνός", Dutch: ""},
				{English: "revered, august, holy, awful", Greek: "σεμνός", Dutch: "eerbiedwaardig, indrukwekkend"},
			},
			expected: []models.Meros{
				{English: "holy", Greek: "ἱερός –ᾶ -ον", Dutch: "heilig"},
				{English: "revered, august, holy, awful", Greek: "σεμνός", Dutch: "eerbiedwaardig, indrukwekkend"},
			},
		},
	}

	// Create an instance of AlexandrosHandler
	handler := &AlexandrosHandler{}

	// Iterate over the test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call the cleanSearchResult function
			filteredResults := handler.cleanSearchResult(tc.results)

			// Use assert.ElementsMatch to compare the slices
			assert.ElementsMatch(t, tc.expected, filteredResults)
		})
	}
}

func TestValidateQueryParam(t *testing.T) {
	tests := []struct {
		queryParam    string
		field         string
		allowedValues []string
		expectedErr   string
	}{
		{
			queryParam:    "english",
			field:         "lang",
			allowedValues: []string{"english", "greek", "dutch"},
			expectedErr:   "",
		},
		{
			queryParam:    "spanish",
			field:         "lang",
			allowedValues: []string{"english", "greek", "dutch"},
			expectedErr:   "invalid lang value. Please choose one of the following: english, greek, dutch",
		},
		{
			queryParam:    "",
			field:         "lang",
			allowedValues: []string{"english", "greek", "dutch"},
			expectedErr:   "lang cannot be empty",
		},
		{
			queryParam:    "any_value",
			field:         "lang",
			allowedValues: nil,
			expectedErr:   "",
		},
	}

	// Create an instance of AlexandrosHandler
	handler := &AlexandrosHandler{}

	for _, test := range tests {
		err := handler.validateQueryParam(test.queryParam, test.field, test.allowedValues)
		if test.expectedErr == "" {
			assert.NoError(t, err, "Expected no error")
		} else {
			assert.EqualError(t, err, test.expectedErr, "Expected error message")
		}
	}
}

func TestProcessQuery(t *testing.T) {
	tests := []struct {
		option        string
		language      string
		queryWord     string
		expectedQuery map[string]interface{}
	}{
		{
			option:    fuzzy,
			language:  "greek",
			queryWord: "test",
			expectedQuery: map[string]interface{}{
				"query": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":    "test",
						"type":     "most_fields",
						"analyzer": "greek_analyzer",
						"fields": []string{
							"greek",
							"original",
						},
					},
				},
				"size": 50,
			},
		},
		{
			option:    phrase,
			language:  "greek",
			queryWord: "test",
			expectedQuery: map[string]interface{}{
				"query": map[string]interface{}{
					"match_phrase": map[string]string{
						"greek": "test",
					},
				},
			},
		},
		{
			option:    exact,
			language:  "greek",
			queryWord: "test",
			expectedQuery: map[string]interface{}{
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"should": []interface{}{
							map[string]interface{}{
								"prefix": map[string]interface{}{
									"greek.keyword": "test",
								},
							},
							map[string]interface{}{
								"term": map[string]interface{}{
									"greek.keyword": "test",
								},
							},
						},
					},
				},
			},
		},
		{
			option:        "unknown",
			language:      "greek",
			queryWord:     "test",
			expectedQuery: nil,
		},
	}

	// Create an instance of AlexandrosHandler
	handler := &AlexandrosHandler{}

	for _, test := range tests {
		query := handler.processQuery(test.option, test.language, test.queryWord)
		assert.Equal(t, test.expectedQuery, query, "Expected query")
	}
}
