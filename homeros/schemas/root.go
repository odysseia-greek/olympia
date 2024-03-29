package schemas

import (
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	plato "github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/olympia/homeros/gateway"
	"log"
	"os"
	"sync"
)

var (
	handler            *gateway.HomerosHandler
	homerosHandlerOnce sync.Once
)

func HomerosHandler() *gateway.HomerosHandler {
	homerosHandlerOnce.Do(func() {
		env := os.Getenv("ENV")
		homerosHandler, err := gateway.CreateNewConfig(env)
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
		"authors": &graphql.Field{
			Type:        graphql.NewList(authors),
			Description: "Get the author and books tree from Herodotos",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context

				// Get the traceID
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				return handler.Books(traceID)
			},
		},

		"sentence": &graphql.Field{
			Type:        sentence,
			Description: "Create a new Question in Herodotos",
			Args: graphql.FieldConfigArgument{
				"book": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"author": &graphql.ArgumentConfig{
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

				book, isOK := p.Args["book"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument book")
				}
				author, isOK := p.Args["author"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument book")
				}
				return handler.Sentence(author, book, traceID)
			},
		},

		"text": &graphql.Field{
			Type:        text,
			Description: "Check the text given",
			Args: graphql.FieldConfigArgument{
				"sentenceId": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"author": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"answer": &graphql.ArgumentConfig{
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

				sentenceId, isOK := p.Args["sentenceId"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument id")
				}
				author, isOK := p.Args["author"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument book")
				}
				answer, isOK := p.Args["answer"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument answer")
				}
				return handler.Answer(sentenceId, author, answer, traceID)
			},
		},

		// Sokrates
		"quiz": &graphql.Field{
			Type: graphql.NewUnion(graphql.UnionConfig{
				Name:  "QuizResponseUnion",
				Types: []*graphql.Object{quizResponseType, dialogueQuizType},
				ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
					if _, ok := p.Value.(*models.QuizResponse); ok {
						return quizResponseType
					}
					if _, ok := p.Value.(*models.DialogueQuiz); ok {
						return dialogueQuizType
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
				"quizType": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				theme, _ := p.Args["theme"].(string)
				set, isOK := p.Args["set"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument set")
				}
				quizType, isOK := p.Args["quizType"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument quizType")
				}

				if quizType == models.DIALOGUE {
					return handler.CreateDialogueQuiz(theme, set, quizType, traceID)
				} else {
					return handler.CreateQuiz(theme, set, quizType, traceID)
				}

			},
		},

		"answer": &graphql.Field{
			Type: graphql.NewUnion(graphql.UnionConfig{
				Name:  "AnswerUnion",
				Types: []*graphql.Object{comprehensiveAnswer, dialogueAnswerType},
				ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
					if _, ok := p.Value.(*models.ComprehensiveResponse); ok {
						return comprehensiveAnswer
					}
					if _, ok := p.Value.(*models.DialogueAnswer); ok {
						return dialogueAnswerType
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
					QuizType:      quizType,
					Comprehensive: comprehensive,
					Answer:        answer,
					Dialogue:      dialogue,
					QuizWord:      quizWord,
				}

				if quizType == models.DIALOGUE {
					return handler.CheckDialogue(answerRequest, traceID)
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
	},
})
