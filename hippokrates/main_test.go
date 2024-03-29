package hippokrates

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	uuid2 "github.com/google/uuid"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"os"
	"strings"
	"testing"
)

const (
	sokratesApi   = "sokrates"
	herodotosApi  = "herodotos"
	alexandrosApi = "alexandros"
	dionysiosApi  = "dionysios"
	ResponseBody  = "responseBody"
	ErrorBody     = "errorBody"
	ContextAuthor = "contextAuthor"
	ContextId     = "contextId"
	UpdatedAnswer = "updateAnswer"
	AnswerBody    = "answerBody"
	Rootword      = "rootWord"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

//go:embed features/*.feature
var featureFiles embed.FS

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func (l *OdysseiaFixture) theIsRunning(service string) error {
	var healthy *models.Health
	uuid := uuid2.New().String()

	switch service {
	case alexandrosApi:
		response, err := l.client.Alexandros().Health(uuid)
		if err != nil {
			return err
		}

		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&healthy)
	case sokratesApi:
		response, err := l.client.Sokrates().Health(uuid)
		if err != nil {
			return err
		}

		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&healthy)

	case herodotosApi:
		response, err := l.client.Herodotos().Health(uuid)
		if err != nil {
			return err
		}

		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&healthy)
	case dionysiosApi:
		response, err := l.client.Dionysios().Health(uuid)
		if err != nil {
			return err
		}

		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&healthy)
	default:
	}

	if !healthy.Healthy {
		return fmt.Errorf("service was %v were a healthy status was expected", healthy.Healthy)
	}

	return nil
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {

		//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=HIPPOKRATES
		logging.System(`
 __ __  ____  ____  ____   ___   __  _  ____    ____  ______    ___  _____
|  |  ||    ||    \|    \ /   \ |  |/ ]|    \  /    ||      |  /  _]/ ___/
|  |  | |  | |  o  )  o  )     ||  ' / |  D  )|  o  ||      | /  [_(   \_ 
|  _  | |  | |   _/|   _/|  O  ||    \ |    / |     ||_|  |_||    _]\__  |
|  |  | |  | |  |  |  |  |     ||     ||    \ |  _  |  |  |  |   [_ /  \ |
|  |  | |  | |  |  |  |  |     ||  .  ||  .  \|  |  |  |  |  |     |\    |
|__|__||____||__|  |__|   \___/ |__|\_||__|\_||__|__|  |__|  |_____| \___|
                                                                          
`)
		logging.System("\"ὄμνυμι Ἀπόλλωνα ἰητρὸν καὶ Ἀσκληπιὸν καὶ Ὑγείαν καὶ Πανάκειαν καὶ θεοὺς πάντας τε καὶ πάσας, ἵστορας ποιεύμενος, ἐπιτελέα ποιήσειν κατὰ δύναμιν καὶ κρίσιν ἐμὴν ὅρκον τόνδε καὶ συγγραφὴν τήνδε:\"")
		logging.System("\"I swear by Apollo Healer, by Asclepius, by Hygieia, by Panacea, and by all the gods and goddesses, making them my witnesses, that I will carry out, according to my ability and judgment, this oath and this indenture.\"")
		logging.System("starting test suite setup.....")

		logging.System("getting env variables and creating config")

	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
	})

	odysseia, err := New()
	if err != nil {
		os.Exit(1)
	}

	//general
	ctx.Step(`^the "([^"]*)" is running$`, odysseia.theIsRunning)

	//alexandros
	ctx.Step(`^the word "([^"]*)" is queried$`, odysseia.theWordIsQueried)
	ctx.Step(`^the word "([^"]*)" is stripped of accents$`, odysseia.theWordIsStrippedOfAccents)
	ctx.Step(`^the partial "([^"]*)" is queried$`, odysseia.thePartialIsQueried)
	ctx.Step(`^the word "([^"]*)" is queried and not found$`, odysseia.theWordIsQueriedWithAndNotFound)
	ctx.Step(`^the word "([^"]*)" should be included in the response$`, odysseia.theWordShouldBeIncludedInTheResponse)
	ctx.Step(`^an error containing "([^"]*)" is returned$`, odysseia.anErrorContainingIsReturned)
	ctx.Step(`^the word "([^"]*)" is queried using "([^"]*)" and "([^"]*)" and "([^"]*)"$`, odysseia.theWordIsQueriedUsingAndAnd)
	ctx.Step(`^a Greek translation should be included in the response$`, odysseia.aGreekTranslationShouldBeIncludedInTheResponse)

	//herodotos
	ctx.Step(`^a query is made for all authors$`, odysseia.aQueryIsMadeForAllAuthors)
	ctx.Step(`^the author "([^"]*)" should be included$`, odysseia.theAuthorShouldBeIncluded)
	ctx.Step(`^a query is made for all books by author "([^"]*)"$`, odysseia.aQueryIsMadeForAllBooksByAuthor)
	ctx.Step(`^a translation is returned$`, odysseia.aTranslationIsReturned)
	ctx.Step(`^an author and book combination is queried$`, odysseia.anAuthorAndBookCombinationIsQueried)
	ctx.Step(`^the number of authors should exceed "([^"]*)"$`, odysseia.theNumberOfAuthorsShouldExceed)
	ctx.Step(`^the book "([^"]*)" should be included$`, odysseia.theBookShouldBeIncluded)
	ctx.Step(`^the sentenceId should be longer than "([^"]*)"$`, odysseia.theSentenceIdShouldBeLongerThan)
	ctx.Step(`^the sentence should include non-ASCII \(Greek\) characters$`, odysseia.theSentenceShouldIncludeNonASCIIGreekCharacters)
	ctx.Step(`^a correctness percentage$`, odysseia.aCorrectnessPercentage)
	ctx.Step(`^a sentence with a translation$`, odysseia.aSentenceWithATranslation)

	//sokrates
	ctx.Step(`^a list of themes with the highest set should be returned$`, odysseia.aListOfThemesWithTheHighestSetShouldBeReturned)
	ctx.Step(`^a query is made for the options for the quizType "([^"]*)"$`, odysseia.aQueryIsMadeForTheOptionsForTheQuizType)
	ctx.Step(`^a new quiz question is made with the quizType "([^"]*)"$`, odysseia.aNewQuizQuestionIsMadeWithTheQuizType)
	ctx.Step(`^the question can be answered from the response$`, odysseia.theQuestionCanBeAnsweredFromTheResponse)

	//dionysios
	ctx.Step(`^the grammar is checked for word "([^"]*)"$`, odysseia.theGrammarIsCheckedForWord)
	ctx.Step(`^the grammar for word "([^"]*)" is queried with an error$`, odysseia.theGrammarForWordIsQueriedWithAnError)
	ctx.Step(`^the declension "([^"]*)" should be included in the response$`, odysseia.theDeclensionShouldBeIncludedInTheResponse)
	ctx.Step(`^the number of results should be equal to or exceed "([^"]*)"$`, odysseia.theNumberOfResultsShouldBeEqualToOrExceed)
	ctx.Step(`^the number of translations should be equal to er exceed "([^"]*)"$`, odysseia.theNumberOfTranslationsShouldBeEqualToErExceed)
	ctx.Step(`^the number of declensions should be equal to or exceed "([^"]*)"$`, odysseia.theNumberOfDeclensionsShouldBeEqualToOrExceed)

	//homeros
	ctx.Step(`^I send a status GraphQL query$`, odysseia.iSendAStatusGraphQLQuery)
	ctx.Step(`^the gateway is up$`, odysseia.theGatewayIsUp)
	ctx.Step(`^authors and books should be returned in a single response$`, odysseia.authorsAndBooksShouldBeReturnedInASingleResponse)
	ctx.Step(`^the grammar is checked for word "([^"]*)" through the gateway$`, odysseia.theGrammarIsCheckedForWordThroughTheGateway)
	ctx.Step(`^the declension "([^"]*)" should be included in the response as a gateway struct$`, odysseia.theDeclensionShouldBeIncludedInTheResponseAsAGatewayStruct)
	ctx.Step(`^the word "([^"]*)" is queried using "([^"]*)" and "([^"]*)" through the gateway$`, odysseia.theWordIsQueriedUsingAndThroughTheGateway)
	ctx.Step(`^I create a new quiz with quizType "([^"]*)"$`, odysseia.iCreateANewQuizWithQuizType)
	ctx.Step(`^other possibilities should be included in the response$`, odysseia.otherPossibilitiesShouldBeIncludedInTheResponse)
	ctx.Step(`^I answer the quiz through the gateway$`, odysseia.iAnswerTheQuizThroughTheGateway)
	ctx.Step(`^the gateway should respond with a correctness$`, odysseia.theGatewayShouldRespondWithACorrectness)
	ctx.Step(`^I answer the sentence through the gateway$`, odysseia.iAnswerTheSentenceThroughTheGateway)
	ctx.Step(`^I create a new sentence response from those methods$`, odysseia.iCreateANewSentenceResponseFromThoseMethods)
	ctx.Step(`^I update my answer using the verified translation$`, odysseia.iUpdateMyAnswerUsingTheVerifiedTranslation)
	ctx.Step(`^I create a new sentence response from those methods with author "([^"]*)"$`, odysseia.iCreateANewSentenceResponseFromThoseMethodsWithAuthor)
	ctx.Step(`^all APIs should be healthy$`, odysseia.allAPIsShouldBeHealthy)
	ctx.Step(`^I query for a tree of Herodotos authors$`, odysseia.iQueryForATreeOfHerodotosAuthors)
	ctx.Step(`^that response should include a Levenshtein distance$`, odysseia.thatResponseShouldIncludeALevenshteinDistance)
	ctx.Step(`^that response should include the number of mistakes with a percentage$`, odysseia.thatResponseShouldIncludeTheNumberOfMistakesWithAPercentage)
	ctx.Step(`^the Levenshtein score should be (\d+)$`, odysseia.theLevenshteinScoreShouldBe)

	//odysseia
	ctx.Step(`^a grammar entry is made for the word "([^"]*)"$`, odysseia.aGrammarEntryIsMadeForTheWord)
	ctx.Step(`^a quiz is played in comprehensive mode for the word "([^"]*)" and the correct answer "([^"]*)" with type "([^"]*)" set "([^"]*)" and theme "([^"]*)"$`, odysseia.aQuizIsPlayedInComprehensiveModeForTheWordAndTheCorrectAnswerWithTypeSetAndTheme)
	ctx.Step(`^the quizresponse is expanded with text and similar words$`, odysseia.theQuizresponseIsExpandedWithTextAndSimilarWords)
	ctx.Step(`^the options returned from the grammar api should include "([^"]*)"$`, odysseia.theOptionsReturnedFromTheGrammarApiShouldInclude)

	//flow
	ctx.Step(`^a response with a rootword is returned$`, odysseia.aResponseWithARootwordIsReturned)
	ctx.Step(`^that rootword is queried in Alexandros with "([^"]*)"$`, odysseia.thatRootwordIsQueriedInAlexandrosWith)
	ctx.Step(`^the grammar is checked for a random word in the list$`, odysseia.theGrammarIsCheckedForARandomWordInTheList)
	ctx.Step(`^the query result has texts included$`, odysseia.theQueryResultHasTextsIncluded)
}

func TestMain(m *testing.M) {
	format := "pretty"
	var tag string // Initialize an empty slice to store the tags

	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" {
			format = "progress"
		} else if strings.HasPrefix(arg, "-tags=") {
			tagsString := strings.TrimPrefix(arg, "-tags=")
			tag = strings.Split(tagsString, ",")[0]
		}
	}

	opts := godog.Options{
		Format:          format,
		FeatureContents: getFeatureContents(), // Get the embedded feature files
	}

	if tag != "" {
		opts.Tags = tag
	}

	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	os.Exit(status)
}

func getFeatureContents() []godog.Feature {
	features := []godog.Feature{}
	featureFileNames, _ := featureFiles.ReadDir("features")
	for _, file := range featureFileNames {
		if !file.IsDir() && file.Name() != "README.md" { // Skip directories and README.md if any
			filePath := fmt.Sprintf("features/%s", file.Name())
			fileContent, _ := featureFiles.ReadFile(filePath)
			features = append(features, godog.Feature{Name: file.Name(), Contents: fileContent})
		}
	}
	return features
}
