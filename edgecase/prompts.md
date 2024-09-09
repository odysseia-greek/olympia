## Prompts


### Elastic Query

```text
I want to build a query in elastic search. I need the raw json that works with kibana and also a golang map[string]interface{}

it should be a fuzzy query that matches the value rootWord with a fuzzines of 2. rootWord is the argument given to the function so it should not be in "" 
the fuzzy query itself should contain the term "greek".

A max of 5 documents should be returned

Please make it a golang function that takes as argument rootWord string and returns the map[string]interface{}
```

### Create Method for password generation

```text
can you make me a method that is part of: func (d *DiogenesHandler) in golang

What it should do is take an argument word string and an int to determine the length.
Combine the word with the current timestamp

It then creates from that greek word a strong password based on the length argument and returns it as a string.
```



### graphql model

```text
Please convert the following model to a graphql model using the golang graphql library:
import "github.com/graphql-go/graphql"

this is the struct

type EdgecaseResponse struct {
	OriginalWord   string  `json:"originalWord"`
	GreekWord      string  `json:"greekWord"`
	StrongPassword string  `json:"strongPassword"`
	SimilarWords   []Meros `json:"similarWords,omitempty"`
}

for SimilarWords you can use an existing reference:

Type: graphql.NewList(dictionary),
		
the model should be named:

convertWordResponseType
and the name itself: convertWordResponse

```

### create unit test

```text

```

### create integration test

```text

```