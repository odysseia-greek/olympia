Feature: Odysseia
  In order to use functions that cross over different apis
  As a greek enthusiast
  We need to be able to validate the workings of these cross api functions

  @odysseia
  Scenario Outline: Creating a new grammar entry and then using comprehensive mode will return a result set
    Given a grammar entry is made for the word "<word>"
    When a quiz is played in comprehensive mode for the word "<fullWord>" and the correct answer "<grammarResponse>" with type "<quizType>" set "<set>" and theme "<theme>" and segment "<segment>"
    Then the options returned from the grammar api should include "<grammarResponse>"
    And the quizresponse is expanded with text and similar words
    Examples:
      | word  | fullWord | grammarResponse | set | theme  | quizType | segment        |
      | ἦν    | εἰμί     | to be           | 1   | Basic  | media    | Common Verbs 1 |
      | ἱερόν | τό ἱερόν | temple          | 1   | Social | media    | City Life      |

  @odysseia
  Scenario: A word can be expanded and used
    Given the grammar is checked for a random word in the list
    When a response with a rootword is returned
    And that rootword is queried in Alexandros with "true"
    Then the query result has texts included

  @wip
  Scenario: A word played in the quiz mode authorbased can be found even if the declension is Ionic for example
    Given an authorbased quiz is played that includes grammarOptions
    When a word is found that would normally not be easy to decline
    And that word is searched for in the grammar component
    Then a result should be returned
