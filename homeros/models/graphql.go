package models

import (
	"encoding/json"
	"github.com/odysseia-greek/agora/plato/models"
)

type Health struct {
	Overall    bool          `json:"overallHealth"`
	Herodotos  models.Health `json:"herodotos"`
	Sokrates   models.Health `json:"sokrates"`
	Dionysios  models.Health `json:"dionysios"`
	Alexandros models.Health `json:"alexandros"`
}

type AuthorGraph struct {
	AuthorTree []AuthorTree `json:"authors"`
}

func (r *AuthorGraph) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalAuthorGraph(data []byte) (*AuthorGraph, error) {
	var r AuthorGraph
	err := json.Unmarshal(data, &r)
	return &r, err
}

type AuthorTree struct {
	Name  string        `json:"name"`
	Books []models.Book `json:"books"`
}

type SentenceGraph struct {
	Author string `json:"author"`
	Book   string `json:"book"`
	Greek  string `json:"greek"`
	Id     string `json:"id"`
}

func UnmarshalMethodGraph(data []byte) (*MethodGraph, error) {
	var r MethodGraph
	err := json.Unmarshal(data, &r)
	return &r, err
}

func (r *MethodGraph) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MethodGraph struct {
	MethodTree []MethodTree `json:"methods"`
}

type MethodTree struct {
	Name       string     `json:"name"`
	Categories []Category `json:"categories"`
}

type Category struct {
	Name           string `json:"name"`
	HighestChapter int64  `json:"highestChapter"`
}

type Answer struct {
	LevenshteinPercentage string            `json:"levenshtein"`
	Input                 string            `json:"input"`
	Quiz                  string            `json:"quiz"`
	SplitQuizSentence     []SplitQuiz       `json:"splitQuiz"`
	SplitAnswerSentence   []SplitAnswer     `json:"splitAnswer"`
	MatchingWords         []MatchingWord    `json:"matches,omitempty"`
	NonMatchingWords      []NonMatchingWord `json:"mistakes,omitempty"`
}

type SplitQuiz struct {
	Word string `json:"word"`
}

type SplitAnswer struct {
	Word string `json:"word"`
}

type MatchingWord struct {
	Word        string `json:"word"`
	SourceIndex int    `json:"index"`
}

type NonMatchingWord struct {
	Word        string  `json:"word"`
	SourceIndex int     `json:"index"`
	Matches     []Match `json:"nonMatches"`
}

type Match struct {
	Match       string `json:"nonMatch"`
	Levenshtein int    `json:"levenshtein"`
	AnswerIndex int    `json:"index"`
	Percentage  string `json:"percentage"`
}

func ParseSentenceResponseToAnswer(sentence models.CheckSentenceResponse) Answer {
	a := Answer{
		LevenshteinPercentage: sentence.LevenshteinPercentage,
		Input:                 sentence.AnswerSentence,
		Quiz:                  sentence.QuizSentence,
	}

	for _, word := range sentence.SplitAnswerSentence {
		m := SplitAnswer{Word: word}
		a.SplitAnswerSentence = append(a.SplitAnswerSentence, m)
	}

	for _, word := range sentence.SplitQuizSentence {
		m := SplitQuiz{Word: word}
		a.SplitQuizSentence = append(a.SplitQuizSentence, m)
	}

	for _, match := range sentence.MatchingWords {
		m := MatchingWord{
			Word:        match.Word,
			SourceIndex: match.SourceIndex,
		}
		a.MatchingWords = append(a.MatchingWords, m)
	}

	for _, mistake := range sentence.NonMatchingWords {
		m := NonMatchingWord{
			Word:        mistake.Word,
			SourceIndex: mistake.SourceIndex,
			Matches:     nil,
		}
		for _, match := range mistake.Matches {
			innerMatch := Match{
				Match:       match.Match,
				Levenshtein: match.Levenshtein,
				AnswerIndex: match.AnswerIndex,
				Percentage:  match.Percentage,
			}
			m.Matches = append(m.Matches, innerMatch)
		}
		a.NonMatchingWords = append(a.NonMatchingWords, m)
	}

	return a
}
func UnmarshalAnswer(data []byte) (*Answer, error) {
	var r Answer
	err := json.Unmarshal(data, &r)
	return &r, err
}

func (r *Answer) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
