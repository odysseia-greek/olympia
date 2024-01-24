Feature: Sokrates
  In order to work with multiple choice quizes
  As a greek enthusiast
  We need to be able to validate the functioning of the Sokrates api

  @sokrates
  Scenario Outline: Querying for a last chapter should return a last chapter
    Given the "<service>" is running
    When a query is made for all methods
    And a random method is queried for categories
    And a random category is queried for the last chapter
    Then that chapter should be a number above 0
    Examples:
      | service  |
      | sokrates |

  @sokrates
  Scenario Outline: The flow to create and answer a question should return a right or wrong answer
    Given the "<service>" is running
    When a new quiz question is requested
    And that question is answered with a "<answer>" answer
    Then the result should be "<answer>"
    Examples:
      | service  | answer |
      | sokrates | true   |
      | sokrates | false  |
