## Prompts


### Create Method for English To Greek translation

```text
Can you generate a method for me that takes a string word and converts the english word to a greek word following the rules outlined in this mapping:

func createEnglishToGreekDict() (map[string]string, error) {
    const translationMap = `
        {
            "a": "α",  "a/": "ά",  "a\\": "ὰ", "a=": "ᾶ", "a(": "ἀ", "a)": "ἁ", "a|": "ᾳ", 
            "A": "Α",  "A/": "Ά",  "A\\": "Ὰ", "A=": "ᾶ", "A(": "Ἀ", "A)": "Ἁ", "A|": "ᾼ",
            "b": "β",  "B": "Β",
			"c": "κ",  "C": "Κ",
            "d": "δ",  "D": "Δ",
            "e": "ε",  "e/": "έ",  "e\\": "ὲ", "e(": "ἐ", "e)": "ἑ",
            "E": "Ε",  "E/": "Έ",  "E\\": "Ὲ", "E(": "Ἐ", "E)": "Ἑ",
            "f": "φ",  "F": "Φ",
            "g": "γ",  "G": "Γ",
            "h": "η",  "h/": "ή",  "h\\": "ὴ", "h=": "ῆ", "h(": "ἠ", "h)": "ἡ", "h|": "ῃ",
            "H": "Η",  "H/": "Ή",  "H\\": "Ὴ", "H=": "ῆ", "H(": "Ἠ", "H)": "Ἡ", "H|": "ῌ",
            "i": "ι",  "i/": "ί",  "i\\": "ὶ", "i=": "ῖ", "i(": "ἰ", "i)": "ἱ",
            "I": "Ι",  "I/": "Ί",  "I\\": "Ὶ", "I=": "ῖ", "I(": "Ἰ", "I)": "Ἱ",
            "j": "ξ",  "J": "Ξ",
            "k": "κ",  "K": "Κ",
            "l": "λ",  "L": "Λ",
            "m": "μ",  "M": "Μ",
            "n": "ν",  "N": "Ν",
            "o": "ο",  "o/": "ό",  "o\\": "ὸ", "o(": "ὀ", "o)": "ὁ",
            "O": "Ο",  "O/": "Ό",  "O\\": "Ὸ", "O(": "Ὀ", "O)": "Ὁ",
            "p": "π",  "P": "Π",
            "q": "θ",  "Q": "Θ",
            "r": "ρ",  "r(": "ῤ", "r)": "ῥ",
            "R": "Ρ",  "R(": "Ῥ",
            "s": "σ",  "s_end": "ς", "S": "Σ",
            "t": "τ",  "T": "Τ",
            "u": "υ",  "u/": "ύ",  "u\\": "ὺ", "u=": "ῦ", "u(": "ὐ", "u)": "ὑ", "u|": "ῡ",
            "U": "Υ",  "U/": "Ύ",  "U\\": "Ὺ", "U=": "ῦ", "U(": "Ὑ", "U)": "Ὑ", "U|": "Ῡ",
            "w": "ω",  "w/": "ώ",  "w\\": "ὼ", "w=": "ῶ", "w(": "ὠ", "w)": "ὡ", "w|": "ῳ",
            "W": "Ω",  "W/": "Ώ",  "W\\": "Ὼ", "W=": "ῶ", "W(": "Ὠ", "W)": "Ὡ", "W|": "ῼ",
            "x": "χ",  "X": "Χ",
            "y": "ψ",  "Y": "Ψ",
            "z": "ζ",  "Z": "Ζ"
        }`

    var dict map[string]string
    err := json.Unmarshal([]byte(translationMap), &dict)
    if err != nil {
        return nil, err
    }
    
    return dict, nil
}

the method should be part of func (d *DiogenesHandler) and the mapping can be found in that struct:

d.EnglishToGreekDict[]

Please only give me the method without anything else.
```

### Create Method for password generation

```text
can you make me a method that is part of: func (d *DiogenesHandler) in golang

What it should do is take an argument word string and an int to determine the length.
Combine the word with the current timestamp

It then creates from that word a strong password based on the length argument and returns it as a string.
```


### Elastic Query

```text
I want to build a query in elastic search. I need the raw json that works with kibana and also a golang map[string]interface{} the index is called "dictionary" and the rootword for a test can be φερω

it should be a fuzzy query that matches the value rootWord with a fuzzines of 2. rootWord is the argument given to the function so it should not be in "" 
the fuzzy query itself should contain the term "greek".

A max of 5 documents should be returned

Please make it a golang function that takes as argument rootWord string and returns the map[string]interface{} and only return the actual method
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
		
then add a "convert" method for my root.go taking the "grammar" example above and use the convertWordResponseType. It takes one argument "rootword" which should be converted to r := models.EdgecaseRequest{Rootword: rootword}
and passed to the method on the handler (handler.Convert()) which needs []byte so set it to json

```

### create unit test

```text
the following is the input needed for the tests:
name - in - out
		{"Empty string", "", ""},
		{"All uppercase word", "UPPERCASE", "ΥΠΠΕΡΚΑΣΕ"},
		{"All lowercase word", "lowercase", "λοωερκασε"},
		{"Word ending in sigma", "sometimes", "σομετιμες"},
		{"Word with a sigma", "some", "σομε"},
		{"Complex word with accents", "A)qh=naios", "Ἁθῆναιος"},
		{"Simple word with accents", "lo/gos", "λόγος"},
		
	
		the createEnglishToGreekDict() is part of config.go and should just be taken as is and added to the handler like this example:

	eToGDict, err := createEnglishToGreekDict()
	if err != nil {
		t.Fatalf("Failed to create English to Greek dictionary: %v", err)
	}
	d := DiogenesHandler{
		EnglishToGreekDict: eToGDict,
	}
```

### create integration test

```text
I have a system to test create my integration tests:

const (
	dialogueModel = `{
  "took": 0,
  "timed_out": false,
  "_shards": {
    "total": 2,
    "successful": 2,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": {
      "value": 3,
      "relation": "eq"
    },
    "max_score": 6.2132735,
    "hits": [
      {
        "_index": "quiz",
        "_id": "ZnugVY0BlxdSsJ_QaT5p",
        "_score": 6.2132735,
        "_ignored": [
          "dialogue.introduction.keyword",
          "content.translation.keyword",
          "content.greek.keyword"
        ],
        "_source": {
          "quizMetadata": {
            "language": "English"
          },
          "theme": "Euripides - Medea",
          "quizType": "dialogue",
          "set": 1,
          "dialogue": {
            "introduction": "Euripides' 'Medea,' an ancient Greek tragedy written in 431 BCE, stands as a seminal work in the canon of classical literature. This play delves into the complexities of passion, betrayal, and revenge through its central character, Medea, a woman of formidable intelligence and emotion. Euripides challenges the societal norms of his time by portraying a female protagonist who defies the submissive role traditionally assigned to women in Greek society.  The narrative unfolds in Corinth, where Medea, a foreigner and a sorceress, grapples with the betrayal of her husband, Jason. Having forsaken her homeland and committed unspeakable acts for his sake, Medea is devastated when Jason abandons her and their children to marry Glauce, the daughter of Creon, the Corinthian king. The play examines Medea’s psychological turmoil as she oscillates between love, hatred, and the pursuit of justice.  Euripides' masterful use of dramatic tension, combined with his exploration of themes such as the plight of the outsider, the consequences of hubris, and the complexities of the human psyche, make 'Medea' a timeless tragedy. It not only reflects the mores of ancient Greek society but also resonates with contemporary audiences, inviting reflection on the nature of justice, loyalty, and the human condition.",
            "speakers": [
              {
                "name": "ΚΡΕΩΝ",
                "shorthand": "Κρέων",
                "translation": "Creon"
              },
              {
                "name": "ΜΗΔΕΙΑ",
                "shorthand": "Μήδεια",
                "translation": "Medea"
              }
            ],
            "section": "315-356",
            "linkToPerseus": "https://scaife.perseus.org/reader/urn:cts:greekLit:tlg0006.tlg003.perseus-grc2:300-360?right=perseus-eng2"
          },
          "content": [
            {
              "translation": "Thy words are soft to hear, but much I dread lest thou art devising some mischief in thy heart, and less than ever do I trust thee now; for a cunning woman, and man likewise,is easier to guard against when quick-tempered than when taciturn. Nay, begone at once! speak me no speeches, for this is decreed, nor hast thou any art whereby thou shalt abide amongst us, since thou hatest me.",
              "greek": "λέγεις ἀκοῦσαι μαλθάκ’, ἀλλ’ ἔσω φρενῶν ὀρρωδία μοι μή τι βουλεύσῃς κακόν, τόσῳ δέ γ’ ἧσσον ἢ πάρος πέποιθά σοι· γυνὴ γὰρ ὀξύθυμος, ὡς δ’ αὔτως ἀνήρ, ῥᾴων φυλάσσειν ἢ σιωπηλὸς σοφός.ἀλλ’ ἔξιθ’ ὡς τάχιστα, μὴ λόγους λέγε· ὡς ταῦτ’ ἄραρε, κοὐκ ἔχεις τέχνην ὅπως μενεῖς παρ’ ἡμῖν οὖσα δυσμενὴς ἐμοί.",
              "place": 1,
              "speaker": "Κρέων"
            }
          ]
        }
      }
    ]
  }
}`

func TestCreateQuizEndpoint(t *testing.T) {
	ticker := time.NewTicker(1 * time.Hour)
	quizAttempts := make(chan models.QuizAttempt)
	aggregatedResult := make(map[string]models.QuizAttempt)
	randomizer, err := config.CreateNewRandomizer()
	assert.Nil(t, err)

	t.Run("Dialogue", func(t *testing.T) {
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient([][]byte{[]byte(dialogueModel)}, mockCode)
		assert.Nil(t, err)

		request := models.CreationRequest{
			Theme:    "sometheme",
			Set:      "1",
			QuizType: models.DIALOGUE,
		}

		jsonBody, err := json.Marshal(request)
		bodyInBytes := bytes.NewReader(jsonBody)
		assert.Nil(t, err)

		testConfig := SokratesHandler{
			Elastic:            mockElasticClient,
			QuizAttempts:       quizAttempts,
			AggregatedAttempts: aggregatedResult,
			Ticker:             ticker,
		}

		router := InitRoutes(&testConfig)
		response := performPostRequest(router, "/sokrates/v1/quiz/create", bodyInBytes)

		var dialogue models.DialogueQuiz
		err = json.NewDecoder(response.Body).Decode(&dialogue)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, len(dialogue.Content) > 1)
	})
	
	func performPostRequest(r http.Handler, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
	
	now i want to create a similar test. Everywhere where it says sokrates/Sokrates it should say Diogenes and the endpoint is:
	
	/diogenes/v1/words/_convert
	
	The DiogenesHandler takes Elastic and a create EnglishToGreekDict that can be created by calling the method createEnglishToGreekDict() this can be set for all the tests since it doesnt change, it return the dict and an error
	
	
	Also these are the test cases that should be included:
	
		expectedWord := "φερω"
		expectedStrongPasswordLength := 16
		expectedSimilarWordsLength := 2

    and the model for the response is:
    
type EdgecaseResponse struct {
	OriginalWord   string  `json:"originalWord"`
	GreekWord      string  `json:"greekWord"`
	StrongPassword string  `json:"strongPassword"`
	SimilarWords   []Meros `json:"similarWords,omitempty"`
}

    and the request is simply
    
    type EdgecaseRequest struct {
	// example: ferw
	// required: true
	Rootword string `json:"rootword"`
}
    
    This is the json input that goes into the mockElasticClient instead of the dialogueModel
	
{  "took":4,
  "timed_out":false,
  "_shards":{
    "total":1,
    "successful":1,
    "skipped":0,
    "failed":0
  },
  "hits":{
    "total":{
      "value":176,
      "relation":"eq"
    },
    "max_score":21.65853,
    "hits":[
      {
        "_index":"dictionary",
        "_type":"",
        "_id":"ky0c1pEBX4KIXZ5B04dm",
        "_score":21.65853,
        "_source":{
          "english":"carry out",
          "greek":"ἐκφέρω"
        }
      },
      {
        "_index":"dictionary",
        "_type":"",
        "_id":"hy0c1pEBX4KIXZ5B04pq",
        "_score":19.556658,
        "_source":{
          "english":"carry on, make a difference",
          "greek":"διαφέρω"
        }
      }
    ]
  }
}

```