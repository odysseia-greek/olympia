package hippokrates

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
	graphql "github.com/odysseia-greek/olympia/homeros/models"
	"net/http"
	"strconv"
)

func (l *OdysseiaFixture) theGatewayIsUp() error {
	// Send the GraphQL query as an HTTP GET request
	response, err := http.Get(l.homeros.health)
	if err != nil {
		return err
	}

	var health models.Health
	err = json.NewDecoder(response.Body).Decode(&health)

	if !health.Healthy {
		return fmt.Errorf("service was %v were a healthy status was expected", health.Healthy)
	}

	return nil
}

func (l *OdysseiaFixture) allAPIsShouldBeHealthy() error {
	resp := l.ctx.Value(ResponseBody).(graphql.Health)
	if !resp.Overall {
		return fmt.Errorf("service was %v were a healthy status was expected", resp.Overall)
	}

	return nil
}

func (l *OdysseiaFixture) iSendAStatusGraphQLQuery() error {
	// Define your GraphQL query
	query := `query healthQuery {
	status {
		overallHealth
		herodotos {
			healthy
		}
		sokrates {
			healthy
		}
		alexandros {
			healthy
		}
		dionysios {
			healthy
		}
	}
}
`

	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var healthResponse struct {
		Data struct {
			Status graphql.Health `json:"status"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&healthResponse)

	l.ctx = context.WithValue(l.ctx, ResponseBody, healthResponse.Data.Status)

	return nil
}

func (l *OdysseiaFixture) iQueryForATreeOfHerodotosAuthors() error {
	// Define your GraphQL query
	query := `query authors {
	authors {
		name
		books {
			book
		}
	}
}
`
	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var herodotosResponse struct {
		Data struct {
			Authors []graphql.AuthorTree `json:"authors"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&herodotosResponse)

	l.ctx = context.WithValue(l.ctx, ResponseBody, herodotosResponse.Data.Authors)

	return nil
}

func (l *OdysseiaFixture) iQueryForATreeOfSokratesMethods() error {
	// Define your GraphQL query
	query := `query methods {
	methods {
		name
		categories {
			name
			highestChapter
		}
	}
}
`
	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var sokratesResponse struct {
		Data struct {
			Methods graphql.MethodGraph `json:"methods"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&sokratesResponse)

	l.ctx = context.WithValue(l.ctx, ResponseBody, sokratesResponse.Data.Methods)

	return nil
}

func (l *OdysseiaFixture) theGrammarIsCheckedForWordThroughTheGateway(word string) error {
	// Define your GraphQL query
	query := fmt.Sprintf(`query grammar {
	grammar(word: "%s") {
		translation
		word
		rule
		rootWord
	}
}`, word)
	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var dionysiosResponse struct {
		Data struct {
			Response []models.Result `json:"grammar"`
		} `json:"data"`
	}

	err = json.NewDecoder(response.Body).Decode(&dionysiosResponse)

	l.ctx = context.WithValue(l.ctx, ResponseBody, dionysiosResponse.Data.Response)

	return nil
}

func (l *OdysseiaFixture) theWordIsQueriedUsingAndThroughTheGateway(word, mode, language string) error {
	// Define your GraphQL query
	query := fmt.Sprintf(`query dictionary {
	dictionary(word: "%s", language: "%s", mode: "%s") {
		greek
		english
		original
		dutch
		linkedWord
	}
}
`, word, language, mode)
	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var alexandrosResponse struct {
		Data struct {
			Response []models.Meros `json:"dictionary"`
		} `json:"data"`
	}

	err = json.NewDecoder(response.Body).Decode(&alexandrosResponse)

	l.ctx = context.WithValue(l.ctx, ResponseBody, alexandrosResponse.Data.Response)

	return nil
}

func (l *OdysseiaFixture) authorsAndBooksShouldBeReturnedInASingleResponse() error {
	authors := l.ctx.Value(ResponseBody).([]graphql.AuthorTree)

	if authors == nil {
		return errors.New("empty tree returned")
	}
	for _, author := range authors {
		if author.Name == "" {
			return errors.New("empty name for an author")
		}
		if len(author.Books) == 0 {
			return errors.New("no books for author " + author.Name)
		}
	}
	return nil
}

func (l *OdysseiaFixture) methodsAndCategoriesShouldBeReturnedInASingleResponse() error {
	methods := l.ctx.Value(ResponseBody).(graphql.MethodGraph)

	for _, method := range methods.MethodTree {
		if method.Name == "" {
			return errors.New("empty name for a method")
		}
		if len(method.Categories) == 0 {
			return errors.New("no categories for method " + method.Name)
		}
	}
	return nil
}

func (l *OdysseiaFixture) theDeclensionShouldBeIncludedInTheResponseAsAGatewayStruct(declension string) error {
	declensions := l.ctx.Value(ResponseBody).([]models.Result)

	found := false
	for _, decResult := range declensions {
		if decResult.Rule == declension {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("could not find declension %v in slice", declension)
	}
	return nil
}

func (l *OdysseiaFixture) iAnswerTheQuizWithAAnswerThroughTheGateway(answer string) error {
	quiz := l.ctx.Value(ResponseBody).(models.QuizResponse)
	correct, _ := strconv.ParseBool(answer)
	provided := quiz.Answer

	if !correct {
		for _, quizQuestion := range quiz.QuizQuestions {
			if quizQuestion != quiz.Answer {
				provided = quizQuestion
				break
			}
		}
	}

	query := fmt.Sprintf(`query answer {
	answer(quizWord: "%s", answerProvided: "%s") {
		correct
		quizWord
		possibilities{
			greek
			category
			translation
		}
	}
}
`, quiz.Question, provided)

	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var sokratesResponse struct {
		Data struct {
			Response models.CheckAnswerResponse `json:"answer"`
		} `json:"data"`
	}

	err = json.NewDecoder(response.Body).Decode(&sokratesResponse)

	l.ctx = context.WithValue(l.ctx, ResponseBody, sokratesResponse.Data.Response)

	return nil
}

func (l *OdysseiaFixture) iCreateANewQuizFromThoseMethods() error {
	// Define your GraphQL query
	query := `query quiz {
	quiz(category: "nomina", chapter: "1", method: "mouseion") {
		question
		answer
		quiz
	}
}
`
	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var sokratesResponse struct {
		Data struct {
			Response models.QuizResponse `json:"quiz"`
		} `json:"data"`
	}

	err = json.NewDecoder(response.Body).Decode(&sokratesResponse)

	l.ctx = context.WithValue(l.ctx, ResponseBody, sokratesResponse.Data.Response)

	return nil
}

func (l *OdysseiaFixture) otherPossibilitiesShouldBeIncludedInTheResponse() error {
	response := l.ctx.Value(ResponseBody).(models.CheckAnswerResponse)

	if len(response.Possibilities) == 0 {
		return fmt.Errorf("expected possibilties to be greater than zero: %v", response.Possibilities)
	}

	return nil
}

func (l *OdysseiaFixture) theGatewayShouldRespondWithACorrect(answer string) error {
	parsedAnswer, err := strconv.ParseBool(answer)
	if err != nil {
		return err
	}
	response := l.ctx.Value(ResponseBody).(models.CheckAnswerResponse)

	if parsedAnswer != response.Correct {
		return fmt.Errorf("expected answer %v to be equal to correctness %v", parsedAnswer, response.Correct)
	}
	return nil
}

func (l *OdysseiaFixture) iAnswerTheSentenceThroughTheGateway() error {
	sentence := l.ctx.Value(ResponseBody).(graphql.SentenceGraph)

	query := fmt.Sprintf(`query text {
	text(
		author: "%s"
		sentenceId: "%s"
		answer: "a het the many yasm was is si aws e ui opp ti up pu small lsa er tha thet no ona oe"
	) {
		levenshtein
		input
		quiz
		splitQuiz {
			word
		}
		splitAnswer {
			word
		}
		matches {
			word
			index
		}
		mistakes {
			word
			index
			nonMatches {
				match
				levenshtein
				index
				percentage
			}
		}
	}
}

`, sentence.Author, sentence.Id)

	l.ctx = context.WithValue(l.ctx, ContextAuthor, sentence.Author)
	l.ctx = context.WithValue(l.ctx, ContextId, sentence.Id)

	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var herodotosResponse struct {
		Data struct {
			Text graphql.Answer `json:"text"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&herodotosResponse)

	l.ctx = context.WithValue(l.ctx, ResponseBody, herodotosResponse.Data.Text)

	return nil
}

func (l *OdysseiaFixture) iCreateANewSentenceResponseFromThoseMethodsWithAuthor(author string) error {
	authorTree := l.ctx.Value(ResponseBody).([]graphql.AuthorTree)
	var book models.Book
	for _, auth := range authorTree {
		if auth.Name == author {
			book = auth.Books[GenerateRandomNumber(len(auth.Books))]
		}
	}

	query := fmt.Sprintf(`query sentence {
	sentence(author: "%s", book: "%v") {
		id
		greek
		author
		book
	}
}
`, author, book.Book)

	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var herodotosResponse struct {
		Data struct {
			Sentence graphql.SentenceGraph `json:"sentence"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&herodotosResponse)

	l.ctx = context.WithValue(l.ctx, ResponseBody, herodotosResponse.Data.Sentence)

	return nil
}

func (l *OdysseiaFixture) iCreateANewSentenceResponseFromThoseMethods() error {
	authorTree := l.ctx.Value(ResponseBody).([]graphql.AuthorTree)
	author := authorTree[GenerateRandomNumber(len(authorTree))]
	book := author.Books[GenerateRandomNumber(len(author.Books))]

	query := fmt.Sprintf(`query sentence {
	sentence(author: "%s", book: "%v") {
		id
		greek
		author
		book
	}
}
`, author.Name, book.Book)

	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var herodotosResponse struct {
		Data struct {
			Sentence graphql.SentenceGraph `json:"sentence"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&herodotosResponse)

	l.ctx = context.WithValue(l.ctx, ResponseBody, herodotosResponse.Data.Sentence)

	return nil
}

func (l *OdysseiaFixture) thatResponseShouldIncludeALevenshteinDistance() error {
	answer := l.ctx.Value(ResponseBody).(graphql.Answer)
	levenshtein, err := strconv.ParseFloat(answer.LevenshteinPercentage, 64)
	if err != nil {
		return err
	}

	if levenshtein < 1 {
		return fmt.Errorf("expected a levenshtein higher than 1 but was: %v", levenshtein)
	}

	return nil
}

func (l *OdysseiaFixture) thatResponseShouldIncludeTheNumberOfMistakesWithAPercentage() error {
	answer := l.ctx.Value(ResponseBody).(graphql.Answer)
	misMatches := answer.NonMatchingWords
	var typos []string
	for _, misMatch := range misMatches {
		for _, mistake := range misMatch.Matches {
			percentage, err := strconv.ParseFloat(mistake.Percentage, 32)
			if err != nil {
				return err
			}

			if percentage > 50 {
				typos = append(typos, mistake.Percentage)
			}
		}
	}

	if len(typos) == 0 {
		return fmt.Errorf("expected mismatches to be found but found: %v", len(typos))
	}

	return nil
}

func (l *OdysseiaFixture) iUpdateMyAnswerUsingTheVerifiedTranslation() error {
	answer := l.ctx.Value(ResponseBody).(graphql.Answer)
	author := l.ctx.Value(ContextAuthor).(string)
	id := l.ctx.Value(ContextId).(string)

	query := fmt.Sprintf(`query text {
	text(
		author: "%s"
		sentenceId: "%s"
		answer: "%s"
	) {
		levenshtein
		input
		quiz
		splitQuiz {
			word
		}
		splitAnswer {
			word
		}
		matches {
			word
			index
		}
		mistakes {
			word
			index
			nonMatches {
				match
				levenshtein
				index
				percentage
			}
		}
	}
}

`, author, id, answer.Quiz)
	response, err := l.graphqlHelper(query)
	if err != nil {
		return err
	}

	var herodotosResponse struct {
		Data struct {
			Text graphql.Answer `json:"text"`
		} `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&herodotosResponse)

	l.ctx = context.WithValue(l.ctx, UpdatedAnswer, herodotosResponse.Data.Text)

	return nil
}

func (l *OdysseiaFixture) theLevenshteinScoreShouldBe(expected int) error {
	answer := l.ctx.Value(UpdatedAnswer).(graphql.Answer)
	levenshtein, err := strconv.ParseFloat(answer.LevenshteinPercentage, 32)
	if err != nil {
		return err
	}

	if expected != int(levenshtein) {
		return fmt.Errorf("expected levenshtein to be 100 but was %v", levenshtein)
	}

	return nil
}

func (l *OdysseiaFixture) graphqlHelper(query string) (*http.Response, error) {
	// Define any required variables (if applicable)
	variables := map[string]interface{}{
		// If your query requires variables, define them here
	}

	// Create the JSON payload for the GraphQL query
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	// Convert the payload to JSON
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil
	}

	// Send the GraphQL query as an HTTP POST request
	return http.Post(l.homeros.graphql, "application/json", bytes.NewBuffer(requestBodyBytes))
}
