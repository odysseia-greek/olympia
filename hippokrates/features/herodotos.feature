Feature: Herodotos
  In order to work with sentences
  As a greek enthusiast
  We need to be able to validate the functioning of the Herodotos api


  @herodotos
  Scenario: Querying options should return all options
    Given the "herodotos" is running
    When a query is made for options
    Then a list of books, authors and references should be returned

  # this test assumes a certain amount of words can always be found in the texts
  @herodotos
  Scenario Outline: A word can be analysed for more detailed information about that word
    Given the "herodotos" is running
    When a the word "<word>" is analyzed
    Then the response has a complete analyzes included
    Examples:
      | word  |
      | λόγος    |
      | Ἀθηναῖος |

  @herodotos
  Scenario: A word can be analysed and from the result a text can be created which can then be used to send in a check
    Given the "herodotos" is running
    When a the word "λόγος" is analyzed
    And the response is used to create a new text
    And the sentence is checked against the official translation
    Then the average levenshtein should be perfect

  @herodotos
  Scenario: A translation with an obvious typo should return that typo and the levhenstein should be affected
    Given the "herodotos" is running
    When the text with author "herodotos" and book "histories" and reference "1.1" and section "0" is checked with typos
    Then the response should include possibleTypos
    And the average levenshtein should be less than perfect
