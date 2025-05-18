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
	sokratesApi             = "sokrates"
	herodotosApi            = "herodotos"
	alexandrosApi           = "alexandros"
	dionysiosApi            = "dionysios"
	ResponseBody            = "responseBody"
	ErrorBody               = "errorBody"
	Rootword                = "rootWord"
	TraceId          string = "hippokrates-traceid"
	QuizContext      string = "quizContext"
	BodyContext      string = "bodyContext"
	GrammarContext   string = "grammarContext"
	AggregateContext string = "aggregateContext"
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
	ctx.Step(`^a query is made for options$`, odysseia.aQueryIsMadeForOptions)
	ctx.Step(`^a the word "([^"]*)" is analyzed$`, odysseia.aTheWordIsAnalyzed)
	ctx.Step(`^the response has a complete analyzes included$`, odysseia.theResponseHasACompleteAnalyzesIncluded)
	ctx.Step(`^a list of books, authors and references should be returned$`, odysseia.aListOfBooksAuthorsAndReferencesShouldBeReturned)
	ctx.Step(`^the average levenshtein should be perfect$`, odysseia.theAverageLevenshteinShouldBePerfect)
	ctx.Step(`^the response is used to create a new text$`, odysseia.theResponseIsUsedToCreateANewText)
	ctx.Step(`^the sentence is checked against the official translation$`, odysseia.theSentenceIsCheckedAgainstTheOfficialTranslation)
	ctx.Step(`^the average levenshtein should be less than perfect$`, odysseia.theAverageLevenshteinShouldBeLessThanPerfect)
	ctx.Step(`^the response should include possibleTypos$`, odysseia.theResponseShouldIncludePossibleTypos)
	ctx.Step(`^the text with author "([^"]*)" and book "([^"]*)" and reference "([^"]*)" and section "([^"]*)" is checked with typos$`, odysseia.theTextWithAuthorAndBookAndReferenceAndSectionIsCheckedWithTypos)

	//dionysios
	ctx.Step(`^the grammar is checked for word "([^"]*)"$`, odysseia.theGrammarIsCheckedForWord)
	ctx.Step(`^the grammar for word "([^"]*)" is queried with an error$`, odysseia.theGrammarForWordIsQueriedWithAnError)
	ctx.Step(`^the declension "([^"]*)" should be included in the response$`, odysseia.theDeclensionShouldBeIncludedInTheResponse)
	ctx.Step(`^the number of results should be equal to or exceed "([^"]*)"$`, odysseia.theNumberOfResultsShouldBeEqualToOrExceed)
	ctx.Step(`^the number of translations should be equal to er exceed "([^"]*)"$`, odysseia.theNumberOfTranslationsShouldBeEqualToErExceed)
	ctx.Step(`^the number of declensions should be equal to or exceed "([^"]*)"$`, odysseia.theNumberOfDeclensionsShouldBeEqualToOrExceed)
	ctx.Step(`^the result should just be "([^"]*)" as a rule$`, odysseia.theResultShouldJustBeAsARule)
	ctx.Step(`^the rootWord should include "([^"]*)"$`, odysseia.theRootWordShouldInclude)
	ctx.Step(`^the rootWord should not include "([^"]*)"$`, odysseia.theRootWordShouldNotInclude)
	ctx.Step(`^not "([^"]*)" as a rule$`, odysseia.notAsARule)

	//homeros
	ctx.Step(`^the gateway is up$`, odysseia.theGatewayIsUp)
	ctx.Step(`^the grammar is checked for word "([^"]*)" through the gateway$`, odysseia.theGrammarIsCheckedForWordThroughTheGateway)
	ctx.Step(`^the declension "([^"]*)" should be included in the response as a gateway struct$`, odysseia.theDeclensionShouldBeIncludedInTheResponseAsAGatewayStruct)
	ctx.Step(`^the word "([^"]*)" is queried using "([^"]*)" and "([^"]*)" through the gateway$`, odysseia.theWordIsQueriedUsingAndThroughTheGateway)
	ctx.Step(`^the gateway should respond with a correctness$`, odysseia.theGatewayShouldRespondWithACorrectness)
	ctx.Step(`^a query is made for all text options$`, odysseia.aQueryIsMadeForAllTextOptions)
	ctx.Step(`^that response is used to create a new text$`, odysseia.thatResponseIsUsedToCreateANewText)
	ctx.Step(`^the text is checked against the official translation$`, odysseia.theTextIsCheckedAgainstTheOfficialTranslation)
	ctx.Step(`^the word "([^"]*)" is analyzed through the gateway$`, odysseia.theWordIsAnalyzedThroughTheGateway)
	ctx.Step(`^the word "([^"]*)" is queried using "([^"]*)" and "([^"]*)" and searchInText through the gateway$`, odysseia.theWordIsQueriedUsingAndAndSearchInTextThroughTheGateway)
	ctx.Step(`^a foundInText response should include results$`, odysseia.aFoundInTextResponseShouldIncludeResults)

	//odysseia
	ctx.Step(`^a grammar entry is made for the word "([^"]*)"$`, odysseia.aGrammarEntryIsMadeForTheWord)
	ctx.Step(`^a result should be returned$`, odysseia.aResultShouldBeReturned)
	ctx.Step(`^that word is searched for in the grammar component$`, odysseia.thatWordIsSearchedForInTheGrammarComponent)
	ctx.Step(`^the options returned from the grammar api should include "([^"]*)"$`, odysseia.theOptionsReturnedFromTheGrammarApiShouldInclude)
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
