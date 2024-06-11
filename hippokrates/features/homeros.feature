Feature: Validate Homeros GraphQL Gateway Functionality
  As a Greek enthusiast
  I want to ensure the proper functioning of the Homeros GraphQL gateway

  @homeros
  Scenario Outline: Using the gateway it should be possible to create and answer a quiz
    Given the gateway is up
    When I create a new quiz with quizType "<quizType>"
    And I answer the quiz through the gateway
    Then the gateway should respond with a correctness
    And other possibilities should be included in the response
    Examples:
      | quizType |
      | media    |
      | authorbased |

  @homeros
  Scenario Outline: Alexandros search word
    Given the gateway is up
    When the word "<word>" is queried using "<mode>" and "<language>" through the gateway
    Then a Greek translation should be included in the response
    Examples:
      | word     | mode     | language |
      | ἰδιώτης  | exact    | greek    |
      | ἄλλος    | extended | greek    |
      | ομαι     | partial  | greek    |
      | αλλας    | fuzzy    | greek    |
      | house    | exact    | english  |
      | round    | extended | english  |
      | so       | partial  | english  |

  @homeros
  Scenario: Alexandros search word with expanded results
    Given the gateway is up
    When the word "λογός" is queried using "exact" and "greek" and searchInText through the gateway
    Then a foundInText response should include results

  @homeros
  Scenario Outline: Dionysios search grammar results
    Given the gateway is up
    When the grammar is checked for word "<word>" through the gateway
    Then the declension "<declension>" should be included in the response as a gateway struct
    Examples:
      | declension                     | word    |
      | 2nd plural - pres - mid - ind  | μάχεσθε |
      | 2nd sing - pres - mid - ind    | μάχει   |
      | inf - aorist - pas             | ἀγορευθῆναι  |

  @homeros
  Scenario: Using the gateway it should be possible to get all the options for text creation, create and answer a text
    Given the gateway is up
    When a query is made for all text options
    And that response is used to create a new text
    And the text is checked against the official translation
    Then the average levenshtein should be perfect

  @homeros
  Scenario Outline: Using the gateway words can be analysed
    Given the gateway is up
    When the word "<word>" is analyzed through the gateway
    Then the response has a complete analyzes included
    Examples:
      | word  |
      | λόγος    |
      | Ἀθηναῖος |
