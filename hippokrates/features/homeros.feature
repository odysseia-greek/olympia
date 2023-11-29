Feature: Validate Homeros GraphQL Gateway Functionality
  As a Greek enthusiast
  I want to ensure the proper functioning of the Homeros GraphQL gateway

  @homeros
  Scenario: Status check with the gateway
    Given the gateway is up
    When I send a status GraphQL query
    Then all APIs should be healthy

  @homeros
  Scenario: Herodotos aggregate
    Given the gateway is up
    When I query for a tree of Herodotos authors
    Then authors and books should be returned in a single response

  @homeros
  Scenario: Using the gateway it should be possible to create and answer a sentence
    Given the gateway is up
    When I query for a tree of Herodotos authors
    And I create a new sentence response from those methods
    And I answer the sentence through the gateway
    Then that response should include a Levenshtein distance
    And that response should include the number of mistakes with a percentage

  @homeros
  Scenario Outline: Using the gateway it should be possible to correctly answer a sentence
    Given the gateway is up
    When I query for a tree of Herodotos authors
    And I create a new sentence response from those methods with author "<author>"
    And I answer the sentence through the gateway
    And I update my answer using the verified translation
    Then the Levenshtein score should be 100
    Examples:
      | author       |
      | thucydides   |
      | ploutarchos  |
      | plato        |

  @homeros
  Scenario: Sokrates aggregate
    Given the gateway is up
    When I query for a tree of Sokrates methods
    Then methods and categories should be returned in a single response

  @homeros
  Scenario Outline: Using the gateway it should be possible to create and answer a quiz
    Given the gateway is up
    When I query for a tree of Sokrates methods
    And I create a new quiz from those methods
    And I answer the quiz with a "<answer>" answer through the gateway
    Then the gateway should respond with a correct "<answer>"
    And other possibilities should be included in the response
    Examples:
      | answer |
      | true   |
      | false  |

  @homeros
  Scenario Outline: Alexandros search word
    Given the gateway is up
    When the word "<word>" is queried using "<mode>" and "<language>" through the gateway
    Then a Greek translation should be included in the response
    Examples:
      | word     | mode   | language |
      | ἰδιώτης  | exact  | Greek    |
      | ἄλλος    | phrase | Greek    |
      | ομαι     | fuzzy  | Greek    |
      | house    | exact  | English  |
      | round    | phrase | English  |
      | so       | fuzzy  | English  |

  @homeros
  Scenario Outline: Dionysios search grammar results
    Given the gateway is up
    When the grammar is checked for word "<word>" through the gateway
    Then the declension "<declension>" should be included in the response as a gateway struct
    Examples:
      | declension                     | word    |
      | 2nd plural - pres - mid - ind  | μάχεσθε |
      | 2nd sing - pres - mid - ind    | μάχει   |
