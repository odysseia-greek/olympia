package schemas

import "github.com/graphql-go/graphql"

var grammar = graphql.NewObject(graphql.ObjectConfig{
	Name: "Result",
	Fields: graphql.Fields{
		"word": &graphql.Field{
			Type: graphql.String,
		},
		"rule": &graphql.Field{
			Type: graphql.String,
		},
		"rootWord": &graphql.Field{
			Type: graphql.String,
		},
		"translation": &graphql.Field{
			Type: graphql.String,
		},
	},
})
