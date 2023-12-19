package scholar

import "encoding/json"

type ConjugationForm struct {
	Person string `json:"person"`
	Number string `json:"number"`
	Word   string `json:"word"`
}

type Conjugation struct {
	Tense  string            `json:"tense"`
	Mood   string            `json:"mood"`
	Aspect string            `json:"aspect"`
	Forms  []ConjugationForm `json:"forms"`
}

type RootWordEntry struct {
	RootWord     string        `json:"rootWord"`
	Translations []string      `json:"translations"`
	Conjugations []Conjugation `json:"conjugations"`
}

func UnmarshalRootWordEntry(data []byte) (RootWordEntry, error) {
	var r RootWordEntry
	err := json.Unmarshal(data, &r)
	return r, err
}
