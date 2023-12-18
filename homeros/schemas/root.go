package schemas

import (
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	plato "github.com/odysseia-greek/agora/plato/config"
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
		"methods": &graphql.Field{
			Type:        graphql.NewList(methods),
			Description: "Ask Sokrates for the methods",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context

				// Get the traceID
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				return handler.Methods(traceID)
			},
		},

		"quiz": &graphql.Field{
			Type:        quiz,
			Description: "Create a new question from Sokrates",
			Args: graphql.FieldConfigArgument{
				"method": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"category": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"chapter": &graphql.ArgumentConfig{
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

				method, isOK := p.Args["method"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument method")
				}
				category, isOK := p.Args["category"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument category")
				}
				chapter, isOK := p.Args["chapter"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument chapter")
				}
				return handler.CreateQuestion(method, category, chapter, traceID)
			},
		},

		"answer": &graphql.Field{
			Type:        answer,
			Description: "Check the answer given",
			Args: graphql.FieldConfigArgument{
				"quizWord": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"answerProvided": &graphql.ArgumentConfig{
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

				quizWord, isOK := p.Args["quizWord"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument method")
				}
				answerProvided, isOK := p.Args["answerProvided"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument chapter")
				}
				return handler.CheckQuestion(quizWord, answerProvided, traceID)
			},
		},

		// Alexandros
		"dictionary": &graphql.Field{
			Type:        graphql.NewList(dictionary),
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
				return handler.Dictionary(word, language, mode, traceID)
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
