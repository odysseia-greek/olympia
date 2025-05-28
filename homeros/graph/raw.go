package graph

import (
	"encoding/json"
	"github.com/odysseia-greek/olympia/homeros/graph/model"
)

type RawQuizSection struct {
	Match            *model.MatchQuiz            `json:"match,omitempty"`
	Trivia           *model.TriviaQuiz           `json:"trivia,omitempty"`
	Structure        *model.StructureQuiz        `json:"structure,omitempty"`
	Media            *model.MediaQuiz            `json:"media,omitempty"`
	FinalTranslation *model.FinalTranslationQuiz `json:"final_translation,omitempty"`
}

type RawJourneySegmentQuiz struct {
	Theme       string            `json:"theme"`
	Segment     string            `json:"segment"`
	Number      int32             `json:"number"`
	Sentence    string            `json:"sentence"`
	Translation string            `json:"translation"`
	ContextNote *string           `json:"contextNote"`
	Intro       *model.QuizIntro  `json:"intro"`
	Quiz        []json.RawMessage `json:"quiz"`
}
