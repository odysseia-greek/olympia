package schemas

import "github.com/graphql-go/graphql"

var optionsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Options",
	Fields: graphql.Fields{
		"option":   &graphql.Field{Type: graphql.String},
		"audioUrl": &graphql.Field{Type: graphql.String},
		"imageUrl": &graphql.Field{Type: graphql.String},
	},
})

var quizResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "QuizResponse",
	Fields: graphql.Fields{
		"quizItem": &graphql.Field{Type: graphql.String},
		"options":  &graphql.Field{Type: graphql.NewList(optionsType)},
	},
})

var speakerType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Speaker",
	Fields: graphql.Fields{
		"name":        &graphql.Field{Type: graphql.String},
		"shorthand":   &graphql.Field{Type: graphql.String},
		"translation": &graphql.Field{Type: graphql.String},
	},
})

var dialogueType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Dialogue",
	Fields: graphql.Fields{
		"introduction":  &graphql.Field{Type: graphql.String},
		"speakers":      &graphql.Field{Type: graphql.NewList(speakerType)},
		"section":       &graphql.Field{Type: graphql.String},
		"linkToPerseus": &graphql.Field{Type: graphql.String},
	},
})

var dialogueInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "DialogueInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"translation": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"greek": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"place": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"speaker": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

var quizMetadataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "QuizMetadata",
	Fields: graphql.Fields{
		"language": &graphql.Field{Type: graphql.String},
	},
})

var dialogueContentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DialogueContent",
	Fields: graphql.Fields{
		"translation": &graphql.Field{Type: graphql.String},
		"greek":       &graphql.Field{Type: graphql.String},
		"place":       &graphql.Field{Type: graphql.Int},
		"speaker":     &graphql.Field{Type: graphql.String},
	},
})

var dialogueQuizType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DialogueQuiz",
	Fields: graphql.Fields{
		"quizMetadata": &graphql.Field{Type: quizMetadataType},
		"theme":        &graphql.Field{Type: graphql.String},
		"quizType":     &graphql.Field{Type: graphql.String},
		"set":          &graphql.Field{Type: graphql.Int},
		"dialogue":     &graphql.Field{Type: dialogueType},
		"content":      &graphql.Field{Type: graphql.NewList(dialogueContentType)},
	},
})

var dialogueCorrectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DialogueCorrection",
	Fields: graphql.Fields{
		"translation":  &graphql.Field{Type: graphql.String},
		"greek":        &graphql.Field{Type: graphql.String},
		"place":        &graphql.Field{Type: graphql.Int},
		"speaker":      &graphql.Field{Type: graphql.String},
		"correctPlace": &graphql.Field{Type: graphql.Int},
	},
})

var dialogueAnswerType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DialogueAnswer",
	Fields: graphql.Fields{
		"percentage": &graphql.Field{
			Type: graphql.Float,
		},
		"input": &graphql.Field{
			Type: graphql.NewList(dialogueContentType),
		},
		"answer": &graphql.Field{
			Type: graphql.NewList(dialogueContentType),
		},
		"wronglyPlaced": &graphql.Field{
			Type: graphql.NewList(dialogueCorrectionType),
		},
	},
})

var progressType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Progress",
	Fields: graphql.Fields{
		"timesCorrect":    &graphql.Field{Type: graphql.Int},
		"timesIncorrect":  &graphql.Field{Type: graphql.Int},
		"averageAccuracy": &graphql.Field{Type: graphql.Float},
	},
})

var comprehensiveAnswer = graphql.NewObject(graphql.ObjectConfig{
	Name: "ComprehensiveResponse",
	Fields: graphql.Fields{
		"correct":  &graphql.Field{Type: graphql.Boolean},
		"quizWord": &graphql.Field{Type: graphql.String},
		"foundInText": &graphql.Field{
			Type: analyzeTextResponseType,
		},
		"similarWords": &graphql.Field{
			Type: graphql.NewList(dictionary),
		},
		"progress": &graphql.Field{Type: progressType},
	},
})

var aggregateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Aggregate",
	Fields: graphql.Fields{
		"highestSet": &graphql.Field{Type: graphql.String},
		"name":       &graphql.Field{Type: graphql.String},
	},
})

var aggregateResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AggregateResult",
	Fields: graphql.Fields{
		"aggregates": &graphql.Field{
			Type: graphql.NewList(aggregateType),
		},
	},
})
