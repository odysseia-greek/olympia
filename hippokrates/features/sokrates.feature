Feature: Sokrates
  In order to work with different quiz modes
  As a greek enthusiast
  We need to be able to validate the functioning of the Sokrates api

  @sokrates
  Scenario Outline: Querying for options should return options
    Given the "<service>" is running
    When a query is made for the options for the quizType "<quizType>"
    Then a list of themes with the highest set should be returned
    Examples:
      | service  | quizType       |
      | sokrates | media          |
      | sokrates | multiplechoice |
      | sokrates | dialogue       |
      | sokrates | authorbased    |

  @sokrates
  Scenario Outline: The simple flow to create and answer a quiz should be functional
    Given the "<service>" is running
    When a query is made for the options for the quizType "<quizType>"
    And a new quiz question is made with the quizType "<quizType>"
    Then the question can be answered from the response
    Examples:
      | service  | quizType       |
      | sokrates | media          |
      | sokrates | multiplechoice |
      | sokrates | dialogue       |
      | sokrates | authorbased       |
