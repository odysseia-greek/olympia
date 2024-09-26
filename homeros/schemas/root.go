package schemas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	plato "github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/olympia/homeros/gateway"
	"log"
	"sync"
)

var (
	handler            *gateway.HomerosHandler
	homerosHandlerOnce sync.Once
)

func HomerosHandler() *gateway.HomerosHandler {
	homerosHandlerOnce.Do(func() {
		ctx := context.Background()
		homerosHandler, err := gateway.CreateNewConfig(ctx)
		if err != nil {
			log.Print(err)
		}
		handler = homerosHandler
	})
	return handler
}

var HomerosSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		// shared
		"status": &graphql.Field{
			Type:        status,
			Description: "See if the backendApis are healthy",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context

				// Get the traceID
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				return handler.Health(traceID)
			},
		},

		// Herodotos
		"textOptions": &graphql.Field{
			Type:        aggregationResultType,
			Description: "Fetch options from Herodotos",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context

				// Get the traceID
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				return handler.HerodotosOptions(traceID)
			},
		},

		"analyze": &graphql.Field{
			Type:        analyzeTextResponseType,
			Description: "analyze text based on a rootword",
			Args: graphql.FieldConfigArgument{
				"rootword": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context

				// Get the traceID
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				rootword, isOK := p.Args["rootword"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument rootword")
				}

				r := models.AnalyzeTextRequest{Rootword: rootword}
				jsonBody, err := json.Marshal(r)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal input to JSON: %v", err)
				}

				return handler.Analyze(jsonBody, traceID)
			},
		},

		"create": &graphql.Field{
			Type:        textType,
			Description: "Create a new Text in Herodotos",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{
					Type: createTextInputType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context

				// Get the traceID
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				input, isOK := p.Args["input"].(map[string]interface{})
				if !isOK {
					return nil, fmt.Errorf("expected argument input")
				}

				jsonBody, err := json.Marshal(input)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal input to JSON: %v", err)
				}

				return handler.CreateText(jsonBody, traceID)
			},
		},

		"check": &graphql.Field{
			Type:        checkTextResponseType,
			Description: "Check the text given",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{
					Type: checkTextRequestInputType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context

				// Get the traceID
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				input, isOK := p.Args["input"].(map[string]interface{})
				if !isOK {
					return nil, fmt.Errorf("expected argument input")
				}

				jsonBody, err := json.Marshal(input)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal input to JSON: %v", err)
				}

				return handler.CheckText(jsonBody, traceID)
			},
		},

		// Sokrates
		"quiz": &graphql.Field{
			Type: graphql.NewUnion(graphql.UnionConfig{
				Name:  "QuizResponseUnion",
				Types: []*graphql.Object{quizResponseType, dialogueQuizType, authorBasedQuizType},
				ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
					if _, ok := p.Value.(*models.QuizResponse); ok {
						return quizResponseType
					}
					if _, ok := p.Value.(*models.DialogueQuiz); ok {
						return dialogueQuizType
					}
					if _, ok := p.Value.(*models.AuthorbasedQuizResponse); ok {
						return authorBasedQuizType
					}
					return nil
				},
			}),
			Args: graphql.FieldConfigArgument{
				"theme": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"set": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"segment": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"quizType": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"order": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"excludeWords": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				theme, _ := p.Args["theme"].(string)
				segment, _ := p.Args["segment"].(string)
				order, _ := p.Args["order"].(string)
				excludeWords, _ := p.Args["excludeWords"].([]interface{})
				excludeWordsStr := make([]string, len(excludeWords))

				if excludeWords != nil {
					for i, word := range excludeWords {
						excludeWordsStr[i], _ = word.(string)
					}
				}

				set, isOK := p.Args["set"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument set")
				}
				quizType, isOK := p.Args["quizType"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument quizType")
				}

				if quizType == models.DIALOGUE {
					return handler.CreateDialogueQuiz(theme, set, segment, quizType, traceID)
				} else if quizType == models.AUTHORBASED {
					return handler.CreateAuthorBasedQuiz(theme, set, segment, quizType, traceID, excludeWordsStr)
				} else {
					return handler.CreateQuiz(theme, set, segment, quizType, order, traceID, excludeWordsStr)
				}

			},
		},

		"answer": &graphql.Field{
			Type: graphql.NewUnion(graphql.UnionConfig{
				Name:  "AnswerUnion",
				Types: []*graphql.Object{comprehensiveAnswer, dialogueAnswerType, authorBasedAnswer},
				ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
					if _, ok := p.Value.(*models.ComprehensiveResponse); ok {
						return comprehensiveAnswer
					}
					if _, ok := p.Value.(*models.DialogueAnswer); ok {
						return dialogueAnswerType
					}
					if _, ok := p.Value.(*models.AuthorBasedResponse); ok {
						return authorBasedAnswer
					}
					return nil
				},
			}),
			Args: graphql.FieldConfigArgument{
				"theme": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"set": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"segment": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"quizType": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"quizWord": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"answer": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"comprehensive": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
				"dialogue": &graphql.ArgumentConfig{
					Type: graphql.NewList(dialogueInputType),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context

				// Get the traceID
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				set, isOK := p.Args["set"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument set")
				}
				quizType, isOK := p.Args["quizType"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument quizType")
				}

				theme, _ := p.Args["theme"].(string)
				segment, _ := p.Args["segment"].(string)
				quizWord, _ := p.Args["quizWord"].(string)
				answer, _ := p.Args["answer"].(string)
				comprehensive, _ := p.Args["comprehensive"].(bool)
				dialogueList, _ := p.Args["dialogue"].([]interface{})

				var dialogue []models.DialogueContent
				for _, item := range dialogueList {
					itemMap, ok := item.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("each dialogue item must be a map")
					}

					var dialogueItem models.DialogueContent
					if translation, ok := itemMap["translation"].(string); ok {
						dialogueItem.Translation = translation
					}
					if greek, ok := itemMap["greek"].(string); ok {
						dialogueItem.Greek = greek
					}
					if place, ok := itemMap["place"].(int); ok {
						dialogueItem.Place = place
					}
					if speaker, ok := itemMap["speaker"].(string); ok {
						dialogueItem.Speaker = speaker
					}

					dialogue = append(dialogue, dialogueItem)
				}

				answerRequest := models.AnswerRequest{
					Theme:         theme,
					Set:           set,
					Segment:       segment,
					QuizType:      quizType,
					Comprehensive: comprehensive,
					Answer:        answer,
					Dialogue:      dialogue,
					QuizWord:      quizWord,
				}

				if quizType == models.DIALOGUE {
					return handler.CheckDialogue(answerRequest, traceID)
				} else if quizType == models.AUTHORBASED {
					return handler.CheckAuthorBased(answerRequest, traceID)
				} else {
					return handler.Check(answerRequest, traceID)
				}
			},
		},

		"options": &graphql.Field{
			Type:        aggregateResultType,
			Description: "returns the options for a specific quiztype",
			Args: graphql.FieldConfigArgument{
				"quizType": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context
				// Get the traceID
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				quizType, isOK := p.Args["quizType"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument quizType")
				}
				return handler.Options(quizType, traceID)
			},
		},

		// Alexandros
		"dictionary": &graphql.Field{
			Type:        extendedDictionary,
			Description: "Search Alexandros dictionary for a word",
			Args: graphql.FieldConfigArgument{
				"word": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"language": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"mode": &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "",
				},
				"searchInText": &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					DefaultValue: false,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context

				// Get the traceID
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				word, isOK := p.Args["word"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument word")
				}
				language, isOK := p.Args["language"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument language")
				}
				mode, isOK := p.Args["mode"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument word")
				}
				searchInText, isOK := p.Args["searchInText"].(bool)
				if !isOK {
					return nil, fmt.Errorf("expected argument searchInText")
				}
				return handler.Dictionary(word, language, mode, traceID, searchInText)
			},
		},

		// Dionysios
		"grammar": &graphql.Field{
			Type:        graphql.NewList(grammar),
			Description: "Search Dionysios for grammar results",
			Args: graphql.FieldConfigArgument{
				"word": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context

				// Get the traceID
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				word, isOK := p.Args["word"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument word")
				}
				return handler.Grammar(word, traceID)
			},
		},
		"convert": &graphql.Field{
			Type:        convertWordResponseType, // The response type
			Description: "Convert root word to Greek and other details",
			Args: graphql.FieldConfigArgument{
				"rootword": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Get rootword from args
				rootword, isOK := p.Args["rootword"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument rootword")
				}

				// Create the EdgecaseRequest model from rootword
				r := models.EdgecaseRequest{Rootword: rootword}

				// Convert request to JSON
				jsonData, err := json.Marshal(r)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal request: %v", err)
				}

				// Get the traceID from context
				ctx := p.Context
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get traceID from context")
				}

				// Call the Convert method on the handler and pass JSON data
				response, err := handler.Convert(jsonData, traceID)
				if err != nil {
					return nil, fmt.Errorf("failed to convert word: %v", err)
				}

				// Return the handler response (should be of type EdgecaseResponse)
				return response, nil
			},
		},
	},
})
