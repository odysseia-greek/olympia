package schemas

import (
	"github.com/graphql-go/graphql"
)

// Define Section Type
var sectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Section",
	Fields: graphql.Fields{
		"key": &graphql.Field{Type: graphql.String},
	},
})

// Define Reference Type
var referenceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Reference",
	Fields: graphql.Fields{
		"key": &graphql.Field{Type: graphql.String},
		"sections": &graphql.Field{
			Type: graphql.NewList(sectionType),
		},
	},
})

// Define ESBook Type
var esBookType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ESBook",
	Fields: graphql.Fields{
		"key": &graphql.Field{Type: graphql.String},
		"references": &graphql.Field{
			Type: graphql.NewList(referenceType),
		},
	},
})

// Define ESAuthor Type
var esAuthorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ESAuthor",
	Fields: graphql.Fields{
		"key": &graphql.Field{Type: graphql.String},
		"books": &graphql.Field{
			Type: graphql.NewList(esBookType),
		},
	},
})

// Define AggregationResult Type
var aggregationResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AggregationResult",
	Fields: graphql.Fields{
		"authors": &graphql.Field{
			Type: graphql.NewList(esAuthorType),
		},
	},
})

// Define Rhema Type
var rhemaType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Rhema",
	Fields: graphql.Fields{
		"greek": &graphql.Field{Type: graphql.String},
		"translations": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"section": &graphql.Field{Type: graphql.String},
	},
})

// Define Text Type
var textType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Text",
	Fields: graphql.Fields{
		"author":          &graphql.Field{Type: graphql.String},
		"book":            &graphql.Field{Type: graphql.String},
		"type":            &graphql.Field{Type: graphql.String},
		"reference":       &graphql.Field{Type: graphql.String},
		"perseusTextLink": &graphql.Field{Type: graphql.String},
		"rhemai": &graphql.Field{
			Type: graphql.NewList(rhemaType),
		},
	},
})

// Define AnalyzeResult Type
var analyzeResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AnalyzeResult",
	Fields: graphql.Fields{
		"referenceLink": &graphql.Field{Type: graphql.String},
		"text":          &graphql.Field{Type: rhemaType},
		"author":        &graphql.Field{Type: graphql.String},
		"book":          &graphql.Field{Type: graphql.String},
		"reference":     &graphql.Field{Type: graphql.String},
	},
})

// Define AnalyzeTextResponse Type
var analyzeTextResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AnalyzeTextResponse",
	Fields: graphql.Fields{
		"rootword": &graphql.Field{Type: graphql.String},
		"conjugations": &graphql.Field{
			Type: graphql.NewList(conjugationResponseType),
		},
		"results": &graphql.Field{
			Type: graphql.NewList(analyzeResultType),
		},
	},
})

var conjugationResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ConjugationResponse",
	Fields: graphql.Fields{
		"word": &graphql.Field{Type: graphql.String},
		"rule": &graphql.Field{Type: graphql.String},
	},
})

// Define CreateTextRequest Type
var createTextInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateTextInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"author":    &graphql.InputObjectFieldConfig{Type: graphql.String},
		"book":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		"reference": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"section":   &graphql.InputObjectFieldConfig{Type: graphql.String}, // Optional field
	},
})

// Define Translations Type
var translationsInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "TranslationsInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"section":     &graphql.InputObjectFieldConfig{Type: graphql.String},
		"translation": &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// Define CheckTextRequest Type
var checkTextRequestInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CheckTextRequestInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"translations": &graphql.InputObjectFieldConfig{
			Type: graphql.NewList(translationsInputType),
		},
		"author":    &graphql.InputObjectFieldConfig{Type: graphql.String},
		"book":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		"reference": &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// Define Typo Type
var typoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Typo",
	Fields: graphql.Fields{
		"source":   &graphql.Field{Type: graphql.String},
		"provided": &graphql.Field{Type: graphql.String},
	},
})

// Define AnswerSection Type
var answerSectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AnswerSection",
	Fields: graphql.Fields{
		"section":               &graphql.Field{Type: graphql.String},
		"levenshteinPercentage": &graphql.Field{Type: graphql.String},
		"quizSentence":          &graphql.Field{Type: graphql.String},
		"answerSentence":        &graphql.Field{Type: graphql.String},
	},
})

// Define CheckTextResponse Type
var checkTextResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CheckTextResponse",
	Fields: graphql.Fields{
		"averageLevenshteinPercentage": &graphql.Field{Type: graphql.String},
		"sections": &graphql.Field{
			Type: graphql.NewList(answerSectionType),
		},
		"possibleTypos": &graphql.Field{
			Type: graphql.NewList(typoType),
		},
	},
})
