package schemas

import "github.com/graphql-go/graphql"

var convertWordResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "convertWordResponse",
		Fields: graphql.Fields{
			"originalWord": &graphql.Field{
				Type: graphql.String,
			},
			"greekWord": &graphql.Field{
				Type: graphql.String,
			},
			"strongPassword": &graphql.Field{
				Type: graphql.String,
			},
			"similarWords": &graphql.Field{
				Type: graphql.NewList(dictionary), // Use existing type for Meros
			},
		},
	},
)
