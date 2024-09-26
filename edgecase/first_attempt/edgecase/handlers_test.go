package edgecase

import (
	"testing"
)

func TestDiogenesHandler_TranslateToGreek(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"All uppercase word", "UPPERCASE", "ΥΠΠΕΡΚΑΣΕ"},
		{"All lowercase word", "lowercase", "λοωερκασε"},
		{"Word ending in sigma", "sometimes", "σομετιμες"},
		{"Word with a sigma", "some", "σομε"},
		{"Complex word with accents", "A)qh=naios", "Ἁθῆναιος"},
		{"Simple word with accents", "lo/gos", "λόγος"},
	}

	eToGDict, err := createEnglishToGreekDict()
	if err != nil {
		t.Fatalf("Failed to create English to Greek dictionary: %v", err)
	}
	d := DiogenesHandler{
		EnglishToGreekDict: eToGDict,
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := d.translateToGreek(tc.input)
			if got != tc.expected {
				t.Errorf("TranslateToGreek() = %v; want %v", got, tc.expected)
			}
		})
	}
}
func TestDiogenesHandler_GenerateStrongPassword(t *testing.T) {
	tests := []struct {
		name           string
		word           string
		passwordLength int
		expectedLength int
	}{
		{"Empty word", "", 10, 10},
		{"Word with length less than passwordLength", "abc", 10, 10},
		{"Word with length equals to passwordLength", "abcdefghij", 10, 10},
		{"Word with length more than passwordLength", "abcdefghijklmno", 10, 10},
		{"Zero passwordLength", "abc", 0, 0},
		{"PasswordLength less than word length", "abcdefghij", 5, 5},
	}

	d := DiogenesHandler{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := d.generateStrongPassword(tc.word, tc.passwordLength)
			if len(got) != tc.expectedLength {
				t.Errorf("GenerateStrongPassword() length = %v; want %v", len(got), tc.expectedLength)
			}
		})
	}
}
