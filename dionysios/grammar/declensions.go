package grammar

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// queryWordInAlexandros tries to find results for given words in the dictionary.
// It queries the Alexandros dictionary for the stripped word and returns the search results.
func (d *DionysosHandler) queryWordInAlexandros(word, traceID string) ([]models.Hit, error) {
	// Remove accents from the word
	strippedWord := d.removeAccents(word)

	// Set the search term and mode
	term := "greek"
	mode := "exact"

	// Send a search request to the Alexandros dictionary
	response, err := d.Client.Alexandros().Search(strippedWord, "greek", mode, "false", traceID)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		// Create a not found error with a unique code and reason
		e := models.NotFoundError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Message: models.NotFoundMessage{
				Type:   term,
				Reason: "not found",
			},
		}
		return nil, &e
	}

	// Decode the response body into search results
	var extendedResponse models.ExtendedResponse
	err = json.NewDecoder(response.Body).Decode(&extendedResponse)
	if err != nil {
		return nil, err
	}

	// Return the search results
	return extendedResponse.Hits, nil
}

// removeAccents removes accents from a given string and returns the transformed string.
// It uses golang.org/x/text/transform package to normalize and remove combining diacritical marks.
func (d *DionysosHandler) removeAccents(s string) string {
	// Create a transformation chain to normalize and remove combining diacritical marks
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	// Apply the transformation to the input string
	output, _, e := transform.String(t, s)
	if e != nil {
		log.Print(e.Error())
	}

	// Return the transformed string
	return output
}

// parseDictResults parses the dictionary hit and extracts the translation and article (if available).
// It splits the Greek term to extract the article and returns the translation and article strings.
func (d *DionysosHandler) parseDictResults(dictionaryHits models.Meros) (translation, article string) {
	// Set the translation to the English field in the dictionary hit
	translation = dictionaryHits.English

	// Split the Greek term by comma
	greek := strings.Split(dictionaryHits.Greek, ",")

	// Check if there is an article present in the Greek term
	if len(greek) > 1 {
		// Remove spaces from the article and assign it to the article variable
		article = strings.Replace(greek[1], " ", "", -1)
	}

	// Return the translation and article
	return
}

// StartFindingRules initiates the process of finding declension rules and translations for a given word.
// It returns the declension translation results.
func (d *DionysosHandler) StartFindingRules(word, requestID string) (*models.DeclensionTranslationResults, error) {
	// Initialize the results variable
	splitID := strings.Split(requestID, "+")

	traceCall := false
	var traceID, spanID string

	if len(splitID) >= 3 {
		traceCall = splitID[2] == "1"
	}

	if len(splitID) >= 1 {
		traceID = splitID[0]
	}
	if len(splitID) >= 2 {
		msg := &pb.StartSpanRequest{
			TraceId: traceID,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		combinedID, err := d.Tracer.StartSpan(ctx, msg)
		if err != nil {
			logging.Error(err.Error())
			spanID = splitID[1]
		} else {
			split := strings.Split(combinedID.CombinedId, "+")
			if len(split) > 1 {
				spanID = split[1]
			}
		}
	}

	var results models.DeclensionTranslationResults

	// Search for declensions for the given word
	declensions, err := d.searchForDeclensions(word)
	if err != nil {
		return nil, err
	}

	if traceCall {
		span := &pb.SpanRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			Action:       "searchForDeclensions",
		}
		go d.Tracer.Span(context.Background(), span)
	}

	// If declensions are found, process each declension
	if len(declensions.Rules) > 0 {
		for _, declension := range declensions.Rules {
			// Process each search term in the declension
			if len(declension.SearchTerms) > 0 {
				for _, term := range declension.SearchTerms {
					// Query the word in the Alexandros dictionary

					dictionaryHits, err := d.queryWordInAlexandros(term, requestID)
					if err != nil {
					}

					// Parse the dictionary results and create result objects
					for _, hit := range dictionaryHits {
						translation, article := d.parseDictResults(hit.Hit)

						result := models.Result{
							Word:        word,
							Rule:        declension.Rule,
							RootWord:    term,
							Translation: translation,
						}

						// Skip adding the result if it already exists with the same translation
						if len(results.Results) > 0 {
							if translation == "" {
								continue
							}
						}

						// Handle articles and filter the results accordingly
						if article != "" {
							switch article {
							case "ὁ":
								if !strings.Contains(declension.Rule, "masc") {
									continue
								}
							case "ἡ":
								if !strings.Contains(declension.Rule, "fem") {
									continue
								}
							case "τό":
								if !strings.Contains(declension.Rule, "neut") {
									continue
								}
							}
						}

						// Append the result to the results slice
						results.Results = append(results.Results, result)
					}
				}
			}
		}
	} else {
		// If no declensions are found, query the word in the Alexandros dictionary directly
		singleSearchResult, err := d.queryWordInAlexandros(word, requestID)
		if err != nil {
			logging.Debug(fmt.Sprintf("single search result gave an error: %s", err.Error()))
		}

		if len(singleSearchResult) > 0 {
			// Parse the dictionary results and create result objects
			for _, searchResult := range singleSearchResult {
				translation, _ := d.parseDictResults(searchResult.Hit)

				// Skip adding the result if it already exists with the same translation
				doNotAdd := false
				for _, res := range results.Results {
					if res.Translation == translation {
						doNotAdd = true
						break
					}
				}

				if doNotAdd {
					continue
				}

				result := models.Result{
					Word:        word,
					Rule:        "preposition",
					RootWord:    searchResult.Hit.Greek,
					Translation: translation,
				}
				results.Results = append(results.Results, result)
			}
		}
	}

	// Handle additional filtering and removal of redundant results
	if results.Results == nil {
		singleSearchResult, err := d.queryWordInAlexandros(word, requestID)
		if err != nil {
			logging.Debug(fmt.Sprintf("single search result gave an error: %s", err.Error()))
		}

		if len(singleSearchResult) > 0 {
			// Parse the dictionary results and create result objects
			for _, searchResult := range singleSearchResult {
				translation, _ := d.parseDictResults(searchResult.Hit)

				// Skip adding the result if it already exists with the same translation
				doNotAdd := false
				for _, res := range results.Results {
					if res.Translation == translation {
						doNotAdd = true
						break
					}
				}

				if doNotAdd {
					continue
				}

				result := models.Result{
					Word:        word,
					Rule:        "preposition",
					RootWord:    searchResult.Hit.Greek,
					Translation: translation,
				}
				results.Results = append(results.Results, result)
			}
		}
	}

	if traceCall {
		span := &pb.SpanRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			Action:       "unfilteredRules",
		}
		go d.Tracer.Span(context.Background(), span)
	}

	// Perform additional filtering and removal of redundant results
	if len(results.Results) > 1 {
		filteredResults := make([]models.Result, 0)
		seen := make(map[string]int) // Map to track unique combinations and their indices in filteredResults

		for _, result := range results.Results {
			// Create a key that represents the unique combination of Rule, Translation, and RootWord (with accents removed)
			key := result.Rule + "|" + result.Translation + "|" + d.removeAccents(result.RootWord)

			if result.Rule == "preposition" || result.Translation == "" ||
				(result.RootWord == "η" || result.RootWord == "ο" || result.RootWord == "το") && result.Word != result.RootWord {
				continue // Skip adding this result
			}

			if index, found := seen[key]; found {
				// Check if the current result has accents but the one in filteredResults does not
				if result.RootWord != d.removeAccents(result.RootWord) && filteredResults[index].RootWord == d.removeAccents(filteredResults[index].RootWord) {
					// Replace the no-accent version with the current result (which has accents)
					filteredResults[index] = result
				}
			} else {
				// Add the current result to filteredResults and record its index in seen
				filteredResults = append(filteredResults, result)
				seen[key] = len(filteredResults) - 1 // Store the index of this result in filteredResults
			}
		}

		// Replace the original results with the filtered results
		results.Results = filteredResults
	}

	// Return the declension translation results
	return &results, nil
}

// searchForDeclensions searches for declensions of a given word.
// It iterates over each declension and declension form, processes them, and returns the found declension rules.
func (d *DionysosHandler) searchForDeclensions(word string) (*models.FoundRules, error) {
	// Initialize the foundRules variable
	var foundRules models.FoundRules

	// Iterate over each declension
	for _, declension := range d.DeclensionConfig.Declensions {
		var contract bool

		// Determine if the declension type is "past" (where you might see contraction)
		switch declension.Type {
		case "past":
			contract = true
		case "irregular":
			// Process irregular verbs separately
			rules := d.loopOverIrregularVerbs(word, declension.Declensions)
			for _, rule := range rules.Rules {
				foundRules.Rules = append(foundRules.Rules, rule)
			}
			continue
		default:
			contract = false
		}

		// Iterate over each declension form
		for _, declensionForm := range declension.Declensions {
			// Process the declension form and get the result
			result := d.loopOverDeclensions(word, declensionForm, contract)

			// Check if any rules were found
			if len(result.Rules) >= 1 {
				for _, rule := range result.Rules {
					d.addRuleIfDifferent(&foundRules, rule)
				}
			}
		}
	}

	// Return the found declension rules
	return &foundRules, nil
}

// loopOverDeclensions processes declensions for a given word and declension form.
// It checks if the declension form matches the word and returns the found declension rules.
func (d *DionysosHandler) loopOverDeclensions(word string, form models.DeclensionElement, prefix bool) models.FoundRules {
	// Initialize the declensions variable
	var declensions models.FoundRules

	// Determine the root cutoff based on the prefix flag
	rootCutOff := 0
	if prefix {
		rootCutOff = 1
	}

	// Remove accents from the declension and trim hyphens
	trimmedLetters := d.removeAccents(strings.Replace(form.Declension, "-", "", -1))

	// Determine the length of the declension
	lengthOfDeclension := utf8.RuneCountInString(trimmedLetters)

	// Remove accnets and convert the word to a rune slice
	// The accents can become a rune and so interfere with the rest of the function
	wordInRune := []rune(d.removeAccents(word))

	// If the length of the declension is greater than the word length, return empty declensions
	if lengthOfDeclension > len(wordInRune) {
		return declensions
	}

	// Extract the letters from the end of the word matching the length of the declension
	lettersOfWord := d.removeAccents(string(wordInRune[len(wordInRune)-lengthOfDeclension:]))

	// If the extracted letters match the trimmed declension, proceed with further checks
	if lettersOfWord == trimmedLetters {
		// Extract the root of the word based on the root cutoff
		rootOfWord := string(wordInRune[rootCutOff : len(wordInRune)-lengthOfDeclension])

		// Get the first letter of the word
		firstLetter := d.removeAccents(string(wordInRune[0]))

		// Initialize the words slice to store search terms
		var words []string

		// Iterate over each search term in the form's search terms
		for _, term := range form.SearchTerm {
			if prefix {
				// Handle prefix declensions
				legitimateStartLetters := []string{"η", "ε"}
				legitimate := false

				// Check if the first letter is one of the legitimate start letters
				for _, startLetter := range legitimateStartLetters {
					if startLetter == firstLetter {
						legitimate = true
					}
				}

				// Skip this search term if the first letter is not legitimate
				if !legitimate {
					continue
				}

				// Handle specific cases for "η" as the first letter
				if firstLetter == "η" {
					vowels := []string{"α", "ε"}

					// Append search terms for each vowel-root combination
					for _, vowel := range vowels {
						searchTerm := fmt.Sprintf("%s%s%s", vowel, rootOfWord, term)
						words = append(words, searchTerm)
					}

					continue
				}
			}

			// Append the search term with the root to the words slice
			searchTerm := fmt.Sprintf("%s%s", rootOfWord, term)
			words = append(words, searchTerm)
		}

		// Create a declension rule based on the form and append it to the declensions slice
		if len(words) >= 1 {
			declension := models.Rule{
				Rule:        form.RuleName,
				SearchTerms: words,
			}
			declensions.Rules = append(declensions.Rules, declension)
		}
	}

	// Return the found declension rules
	return declensions
}

// loopOverIrregularVerbs processes irregular verbs for a given word.
// It checks if the stripped word matches the stripped outcome word and returns the found declension rules.
func (d *DionysosHandler) loopOverIrregularVerbs(word string, declensions []models.DeclensionElement) models.FoundRules {
	// Initialize the rules variable
	var rules models.FoundRules

	// Remove accents from the word
	strippedWord := d.removeAccents(word)

	// Iterate over each declension outcome
	for _, outcome := range declensions {
		// Remove accents from the outcome word
		strippedOutcomeWord := d.removeAccents(outcome.Declension)

		// Check if the stripped word matches the stripped outcome word
		if strippedWord == strippedOutcomeWord {
			// Create a declension rule based on the outcome
			declension := models.Rule{
				Rule:        outcome.RuleName,
				SearchTerms: outcome.SearchTerm,
			}

			// Append the declension rule to the rules slice
			rules.Rules = append(rules.Rules, declension)
		}
	}

	// Return the found declension rules
	return rules
}

// containsRule checks if a given rule is already present in the rules slice.
// It compares both the rule string and the search terms to determine if they are the same.
func (d *DionysosHandler) containsRule(rules []models.Rule, rule models.Rule) bool {
	for _, r := range rules {
		if d.slicesEqual(r.SearchTerms, rule.SearchTerms) && r.Rule == rule.Rule {
			return true
		}
	}
	return false
}

// slicesEqual compares two string slices and returns true if they are equal.
func (d *DionysosHandler) slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// addRuleIfDifferent adds a rule to the foundRules.Rules slice only if it is not already present
// or if the search terms are different.
func (d *DionysosHandler) addRuleIfDifferent(foundRules *models.FoundRules, rule models.Rule) {
	if !d.containsRule(foundRules.Rules, rule) {
		foundRules.Rules = append(foundRules.Rules, rule)
	}
}
