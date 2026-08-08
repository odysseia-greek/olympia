Feature: Smoke test the Homeros GraphQL gateway
  As a client of Homeros
  I want one small end-to-end check for each area of the application
  So that broken gateway wiring is detected without duplicating service integration tests

  Background:
    Given the gateway is up

  @homeros @text
  Scenario: Create and check a text
    When a query is made for all text options
    And that response is used to create a new text
    And the text is checked against the official translation
    Then the average levenshtein should be perfect

  @homeros @dictionary
  Scenario Outline: Search the dictionary through every Homeros search endpoint
    When the dictionary endpoint "<endpoint>" is searched for "<word>" in "<language>"
    Then the dictionary response should contain results
    Examples:
      | endpoint | word    | language |
      | exact    | λόγος   | Greek    |
      | text     | λόγος   | Greek    |
      | fuzzy    | λογος   | Greek    |
      | partial  | λογ     | Greek    |
      | phrase   | λόγος   | Greek    |
      | exact    | word    | English  |

  @homeros @dictionary
  Scenario: An exact dictionary match can include occurrences in texts
    When the dictionary endpoint "exact" is searched for "λόγος" in "Greek" with text expansion
    Then the dictionary response should contain text occurrences

  @homeros @grammar
  Scenario Outline: Check grammar through Homeros
    When the grammar is checked for word "<word>" through the gateway
    Then the declension "<declension>" should be included in the response as a gateway struct
    Examples:
      | declension                    | word         |
      | 2nd plural - pres - mid - ind | μάχεσθε      |
      | 2nd sing - pres - mid - ind   | μάχει        |
      | inf - aorist - pas            | ἀγορευθῆναι  |

  @homeros @quiz
  Scenario: Create and answer a media quiz
    When the media quiz options are requested
    And a media quiz is created from those options
    And the first media quiz option is answered
    Then the gateway should respond with a correctness
