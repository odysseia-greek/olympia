package schemas

import "github.com/graphql-go/graphql"

var authors = graphql.NewObject(graphql.ObjectConfig{
	Name: "Authors",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"books": &graphql.Field{
			Type: graphql.NewList(book),
		},
	},
})

var book = graphql.NewObject(graphql.ObjectConfig{
	Name: "Book",
	Fields: graphql.Fields{
		"book": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var sentence = graphql.NewObject(graphql.ObjectConfig{
	Name: "Sentence",
	Fields: graphql.Fields{
		"author": &graphql.Field{
			Type: graphql.String,
		},
		"book": &graphql.Field{
			Type: graphql.String,
		},
		"greek": &graphql.Field{
			Type: graphql.String,
		},
		"id": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var text = graphql.NewObject(graphql.ObjectConfig{
	Name: "Text",
	Fields: graphql.Fields{
		"levenshtein": &graphql.Field{
			Type: graphql.String,
		},
		"quiz": &graphql.Field{
			Type: graphql.String,
		},
		"input": &graphql.Field{
			Type: graphql.String,
		},
		"splitQuiz": &graphql.Field{
			Type: graphql.NewList(split),
		},
		"splitAnswer": &graphql.Field{
			Type: graphql.NewList(split),
		},
		"matches": &graphql.Field{
			Type: graphql.NewList(matchingWords),
		},
		"mistakes": &graphql.Field{
			Type: graphql.NewList(nonMatchingWords),
		},
	},
})

var split = graphql.NewObject(graphql.ObjectConfig{
	Name: "word",
	Fields: graphql.Fields{
		"word": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var matchingWords = graphql.NewObject(graphql.ObjectConfig{
	Name: "matches",
	Fields: graphql.Fields{
		"word": &graphql.Field{
			Type: graphql.String,
		},
		"index": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var nonMatchingWords = graphql.NewObject(graphql.ObjectConfig{
	Name: "mistakes",
	Fields: graphql.Fields{
		"word": &graphql.Field{
			Type: graphql.String,
		},
		"index": &graphql.Field{
			Type: graphql.Int,
		},
		"nonMatches": &graphql.Field{
			Type: graphql.NewList(nonMatch),
		},
	},
})

var nonMatch = graphql.NewObject(graphql.ObjectConfig{
	Name: "nonMatches",
	Fields: graphql.Fields{
		"match": &graphql.Field{
			Type: graphql.String,
		},
		"levenshtein": &graphql.Field{
			Type: graphql.Int,
		},
		"index": &graphql.Field{
			Type: graphql.Int,
		},
		"percentage": &graphql.Field{
			Type: graphql.String,
		},
	},
})
