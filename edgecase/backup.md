## Edgecase commands and others


### Quote

"Ἀποσκότησόν μοι"
"Get out of my light."

### Ascii
```ascii
 ___    ____  ___    ____    ___  ____     ___  _____
|   \  |    |/   \  /    |  /  _]|    \   /  _]/ ___/
|    \  |  ||     ||   __| /  [_ |  _  | /  [_(   \_
|  D  | |  ||  O  ||  |  ||    _]|  |  ||    _]\__  |
|     | |  ||     ||  |_ ||   [_ |  |  ||   [_ /  \ |
|     | |  ||     ||     ||     ||  |  ||     |\    |
|_____||____|\___/ |___,_||_____||__|__||_____| \___|
```

                                                     
### Elastic

```go
	var meroi []models.Meros
	
	for _, hit := range elasticResult.Hits.Hits {
		jsonHit, _ := json.Marshal(hit.Source)
		meros, _ := models.UnmarshalMeros(jsonHit)
		if meros.Original != "" {
			meros.Greek = meros.Original
			meros.Original = ""
		}
		meroi = append(meroi, meros)
	}


	if traceCall {
		go a.databaseSpan(elasticResult, query, traceID, spanID)
	}
	
	
	edgecaseResponse := models.EdgecaseResponse{
		OriginalWord:   edgecaseRequest.Rootword,
		GreekWord:      "",
		StrongPassword: "",
		SimilarWords:   meroi,
	}

	middleware.ResponseWithCustomCode(w, http.StatusOK, edgecaseResponse)
```

### URL parsing

```go
	var edgecaseRequest models.EdgecaseRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&edgecaseRequest)
	if err != nil {
		e := models.ValidationError{
			ErrorModel: models.ErrorModel{UniqueCode: traceID},
			Messages: []models.ValidationMessages{
				{
					Field:   "decoding",
					Message: err.Error(),
				},
			},
		}
		middleware.ResponseWithJson(w, e)
		return
	}
```


### mapping Dict

```golang

func createEnglishToGreekDict() (map[string]string, error) {
    const translationMap = `
        {
            "a": "α",  "a/": "ά",  "a\\": "ὰ", "a=": "ᾶ", "a(": "ἀ", "a)": "ἁ", "a|": "ᾳ", 
            "A": "Α",  "A/": "Ά",  "A\\": "Ὰ", "A=": "ᾶ", "A(": "Ἀ", "A)": "Ἁ", "A|": "ᾼ",
            "b": "β",  "B": "Β",
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

```

### method to translate

```golang
package main

import (
	"strings"
)

func (d *DiogenesHandler) translateToGreek(word string) string {
	var translated strings.Builder

	// Iterate over the input word while checking for multi-character keys
	for i := 0; i < len(word); i++ {
		// Check if this is the last character and is 's'
		if i == len(word)-1 && word[i] == 's' {
			translated.WriteString(d.EnglishToGreekDict["s_end"])  // Use final sigma 'ς'
		} else {
			// Attempt to match the longest possible key in the englishToGreekDict
			for j := len(word) - i; j > 0; j-- {
				substr := word[i : i+j]
				if greekChar, ok := d.EnglishToGreekDict[substr]; ok {
					translated.WriteString(greekChar)
					i += j - 1  // Move the index to the end of the matched substring
					break
				} else if j == 1 {  // If no match found, keep the original character
					translated.WriteByte(word[i])
				}
			}
		}
	}

	return translated.String()
}
```


### password

```golang
// generateStrongPassword creates a strong password based on the Greek word and the current timestamp
func (d *DiogenesHandler) generateStrongPassword(greekWord string, passwordLength int) string {
	// Combine the Greek word with the current timestamp
	timestamp := time.Now().UnixNano()
	combined := fmt.Sprintf("%s%d", greekWord, timestamp)

	// Generate a SHA-256 hash of the combined string
	hash := sha256.Sum256([]byte(combined))
	hashStr := hex.EncodeToString(hash[:])

	// Format the hash into a strong password (e.g., first 12 chars, mixed case)
	var password strings.Builder
	for i, char := range hashStr[:passwordLength] {
		if i%2 == 0 {
			password.WriteRune(char)
		} else {
			password.WriteString(strings.ToUpper(string(char)))
		}
	}

	return password.String()
}
```

