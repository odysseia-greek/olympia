package scholar

import "encoding/json"

type GrammaticalForm struct {
	Person string `json:"person,omitempty"`
	Number string `json:"number,omitempty"`
	Gender string `json:"gender,omitempty"`
	Case   string `json:"case,omitempty"`
	Word   string `json:"word"`
}

type GrammaticalCategory struct {
	Tense  string            `json:"tense,omitempty"`
	Mood   string            `json:"mood,omitempty"`
	Aspect string            `json:"aspect,omitempty"`
	Forms  []GrammaticalForm `json:"forms"`
}

type RootWordEntry struct {
	RootWord     string                `json:"rootWord"`
	PartOfSpeech string                `json:"partOfSpeech"` //  E.g., "verb", "noun", "participle"
	Translations []string              `json:"translations"`
	Categories   []GrammaticalCategory `json:"categories"`
}

func UnmarshalRootWordEntry(data []byte) (RootWordEntry, error) {
	var r RootWordEntry
	err := json.Unmarshal(data, &r)
	return r, err
}
