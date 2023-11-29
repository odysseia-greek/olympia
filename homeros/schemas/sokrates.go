package schemas

import "github.com/graphql-go/graphql"

var methods = graphql.NewObject(graphql.ObjectConfig{
	Name: "Methods",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"categories": &graphql.Field{
			Type: graphql.NewList(category),
		},
	},
})

var category = graphql.NewObject(graphql.ObjectConfig{
	Name: "Category",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"highestChapter": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var quiz = graphql.NewObject(graphql.ObjectConfig{
	Name: "Quiz",
	Fields: graphql.Fields{
		"question": &graphql.Field{
			Type: graphql.String,
		},
		"answer": &graphql.Field{
			Type: graphql.String,
		},
		"quiz": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
	},
})

var answer = graphql.NewObject(graphql.ObjectConfig{
	Name: "Answer",
	Fields: graphql.Fields{
		"correct": &graphql.Field{
			Type: graphql.Boolean,
		},
		"quizWord": &graphql.Field{
			Type: graphql.String,
		},
		"possibilities": &graphql.Field{
			Type: graphql.NewList(possibilities),
		},
	},
})

var possibilities = graphql.NewObject(graphql.ObjectConfig{
	Name: "Possibilities",
	Fields: graphql.Fields{
		"category": &graphql.Field{
			Type: graphql.String,
		},
		"greek": &graphql.Field{
			Type: graphql.String,
		},
		"translation": &graphql.Field{
			Type: graphql.String,
		},
		"chapter": &graphql.Field{
			Type: graphql.String,
		},
		"method": &graphql.Field{
			Type: graphql.String,
		},
	},
})
