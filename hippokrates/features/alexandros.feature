Feature: Alexandros
  In order to use the dictionary
  As a greek enthusiast
  We need to be able to validate the functioning of the Alexandros api

  @alexandros
  Scenario Outline: Searching for a word in the dictionary that word should be included in the response
    Given the "<service>" is running
    When the word "<word>" is queried
    Then the word "<word>" should be included in the response
    Examples:
      | service    | word     |
      | alexandros | ἀγαθός   |
      | alexandros | ἡσσάομαι |

  @alexandros
  Scenario Outline: Searching for a word stripped of accents the result should contain an original version of that word
    Given the "<service>" is running
    When the word "<word>" is stripped of accents
    Then the word "<word>" should be included in the response
    Examples:
      | service    | word    |
      | alexandros | ὕδατος  |
      | alexandros | ἰδιώτης |

  @alexandros
  Scenario Outline: Searching for the beginning of a word a response with a full set of words should be returned
    Given the "<service>" is running
    When the partial "<partial>" is queried
    Then the word "<word>" should be included in the response
    Examples:
      | service    | partial | word    |
      | alexandros | αγα     | ἀγαθός  |
      | alexandros | ἱστ     | ἱστορία |

  @alexandros
  Scenario Outline: Different modes and languages are supported
    Given the "<service>" is running
    When the word "<word>" is queried using "<mode>" and "<language>" and "<extend>"
    Then a Greek translation should be included in the response
    Examples:
      | service    | word     | mode   | language | extend |
      | alexandros | ἰδιώτης  | exact  | greek    | false  |
      | alexandros | λόγος  | exact  | greek    | true  |
      | alexandros | Ἀθηναῖος  | exact  | greek    | true  |
      | alexandros | ἄλλος    | extended | greek    |  false  |
      | alexandros | ομαι     | partial  | greek    | false  |
      | alexandros | αγαθ     | fuzzy  | greek    | false  |
      | alexandros | house    | exact  | english  | false  |
      | alexandros | round    | extended | english  | false  |
      | alexandros | so       | partial  | english  | false  |

