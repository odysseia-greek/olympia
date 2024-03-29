package schemas

import "github.com/graphql-go/graphql"

var extendedDictionary = graphql.NewObject(graphql.ObjectConfig{
	Name: "ExtendedDictionary",
	Fields: graphql.Fields{
		"hits": &graphql.Field{
			Type: graphql.NewList(extendedDictionaryHit),
		},
	},
})

var extendedDictionaryHit = graphql.NewObject(graphql.ObjectConfig{
	Name: "ExtendedDictionaryEntry",
	Fields: graphql.Fields{
		"foundInText": &graphql.Field{
			Type: foundInTextType,
		},
		"hit": &graphql.Field{
			Type: dictionary,
		},
	},
})

var dictionary = graphql.NewObject(graphql.ObjectConfig{
	Name: "Hit",
	Fields: graphql.Fields{
		"greek": &graphql.Field{
			Type: graphql.String,
		},
		"english": &graphql.Field{
			Type: graphql.String,
		},
		"dutch": &graphql.Field{
			Type: graphql.String,
		},
		"linkedWord": &graphql.Field{
			Type: graphql.String,
		},
		"original": &graphql.Field{
			Type: graphql.String,
		},
	},
})
