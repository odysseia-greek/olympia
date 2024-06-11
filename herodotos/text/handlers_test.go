package text

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestLevenshtein(t *testing.T) {
	t.Run("TestLevenshteinAsInt", func(t *testing.T) {
		sourceString := "This is the display of the inquiry of Herodotos of Halikarnassos"
		targetString := "This is the historical inquiry of Herodotos of Halikarnassos"
		expected := 11

		levenshteinDist := levenshteinDistance(sourceString, targetString)

		assert.Equal(t, expected, levenshteinDist)
	})

	t.Run("TestLevenshteinToPercentageZero", func(t *testing.T) {
		source := 2
		longestWord := 2
		percentage := levenshteinDistanceInPercentage(source, longestWord)

		expected := float64(0)

		assert.Equal(t, expected, percentage)
	})

	t.Run("TestLevenshteinToPercentageHundred", func(t *testing.T) {
		source := 0
		longestWord := 20
		percentage := levenshteinDistanceInPercentage(source, longestWord)

		expected := float64(100)

		assert.Equal(t, expected, percentage)
	})

	t.Run("TestLevenshteinToPercentageMixed", func(t *testing.T) {
		source := 11
		longestWord := 64
		percentage := levenshteinDistanceInPercentage(source, longestWord)

		expected := float64(82.8125)

		assert.Equal(t, expected, percentage)
	})

	t.Run("TestLevenshteinToPercentageMixed", func(t *testing.T) {
		source := 11
		longestWord := 64
		percentage := levenshteinDistanceInPercentage(source, longestWord)

		expected := float64(82.8125)

		assert.Equal(t, expected, percentage)
	})
}

func TestLongestSentence(t *testing.T) {
	sourceString := "this is a lot longer"
	targetString := "short"

	longestSentence := longestStringOfTwo(sourceString, targetString)
	expected := 20

	assert.Equal(t, expected, longestSentence)
}

func TestFindMatchingWords(t *testing.T) {
	sourceString := "This inquiry the Herodotos of Halikarnassos"
	targetString := "This inquiry tHe Herodotus of Halikarnassus"

	response := findTypos(sourceString, targetString)

	expectedTypos := 2
	assert.Equal(t, expectedTypos, len(response))
}

func TestStreamlineString(t *testing.T) {
	sourceString := "This inquiry the Herodotos of Halikarnassos"
	matchingWords := []string{"the", "of"}

	newSentence := streamlineSentenceBeforeCompare(matchingWords, sourceString)

	expected := "This inquiry Herodotos Halikarnassos"
	assert.Equal(t, expected, newSentence)
}
