package api

import (
	"strings"
	"testing"
)

func TestConvertToGreek(t *testing.T) {
	eToGDict, err := createEnglishToGreekDict()
	if err != nil {
		t.Fatalf("Failed to create English to Greek dictionary: %v", err)
	}
	d := DiogenesHandler{
		EnglishToGreekDict: eToGDict,
	}

	var tests = []struct {
		name  string
		input string
		want  string
	}{
		{"Empty string", "", ""},
		{"All uppercase word", "UPPERCASE", "ΥΠΠΕΡΚΑΣΕ"},
		{"All lowercase word", "lowercase", "λοωερκασε"},
		{"Word ending in sigma", "sometimes", "σομετιμες"},
		{"Word with a sigma", "some", "σομε"},
		{"Complex word with accents", "A)qh=naios", "Ἁθῆναιος"},
		{"Simple word with accents", "lo/gos", "λόγος"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := d.ConvertToGreek(tc.input)
			if err != nil {
				t.Fatalf("An error occurred: %v", err)
			}

			if output != tc.want {
				t.Errorf("got %v; want %v", output, tc.want)
			}
		})
	}
}

func TestGenerateStrongPassword(t *testing.T) {
	d := DiogenesHandler{}
	tests := []struct {
		name       string
		word       string
		length     int
		hasSpecial bool
	}{
		{"Short", "pass", 6, true},
		{"Long", "password", 16, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			password, err := d.GenerateStrongPassword(tc.word, tc.length)

			if err != nil {
				t.Fatalf("An error occurred: %v", err)
			}

			if len(password) != tc.length {
				t.Errorf("got length %v; want length %v", len(password), tc.length)
			}

			hasSpecial := strings.ContainsAny(password, "@81307$")
			if hasSpecial != tc.hasSpecial {
				t.Errorf("got special characters %v; want special characters %v", hasSpecial, tc.hasSpecial)
			}
		})
	}
}
