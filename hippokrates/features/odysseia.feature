Feature: Odysseia
  In order to use functions that cross over different apis
  As a greek enthusiast
  We need to be able to validate the workings of these cross api functions

  @odysseia
  Scenario: A word can be expanded and used
    Given the grammar is checked for a random word in the list
    When a response with a rootword is returned
    And that rootword is queried in Alexandros with "true"
    Then the query result has texts included
