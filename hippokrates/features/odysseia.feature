Feature: Odysseia
  In order to use functions that cross over different apis
  As a greek enthusiast
  We need to be able to validate the workings of these cross api functions

  @odysseia
  Scenario Outline: Creating a new grammar entry and then using comprehensive mode will return a result set
    Given a grammar entry is made for the word "<word>"
    When a quiz is played in comprehensive mode for the word "<fullWord>" and the correct answer "<grammarResponse>" with type "<quizType>" set "<set>" and theme "<theme>"
    Then the options returned from the grammar api should include "<grammarResponse>"
    And the quizresponse is expanded with text and similar words
    Examples:
      | word  | fullWord | grammarResponse | set | theme            | quizType    |
      | ἦν    | εἰμί     | to be           | 20  | Plato - Republic | authorbased |
      | ἱερόν | τό ἱερόν | temple          | 1   | Daily Life       | media       |

  @odysseia
  Scenario: A word can be expanded and used
    Given the grammar is checked for a random word in the list
    When a response with a rootword is returned
    And that rootword is queried in Alexandros with "true"
    Then the query result has texts included
