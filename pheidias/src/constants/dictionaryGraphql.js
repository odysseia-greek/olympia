import gql from 'graphql-tag';

export const DictionaryExact = gql`
  query DictionaryExact($input: ExpandableSearchQueryInput!) {
    exact(input: $input) {
      results {
        headword
        partOfSpeech
        normalized
        quickGlosses {
          language
          gloss
        }
        noun {
          declension
          genitive
        }
        verb {
          principalParts
        }
        modernConnections {
          term
          note
        }
        definitions {
          grade
          meanings {
            definition
            language
          }
        }
        linkedWord
      }
      pageInfo {
        page
        total
      }
      similarWords {
        english
        greek
        original
      }
      foundInText {
        rootword
        texts {
          author
          book
          reference
          referenceLink
          text {
            greek
            section
            translations
          }
        }
        conjugations {
          rule
          word
        }
      }
    }
  }
`;

export const DictionaryPartial = gql`
  query DictionaryPartial($input: SearchQueryInput!) {
    partial(input: $input) {
      results {
        headword
        partOfSpeech
        normalized
        quickGlosses {
          language
          gloss
        }
        verb {
          principalParts
        }
        modernConnections {
          term
          note
        }
        definitions {
          grade
          meanings {
            definition
            language
          }
        }
        linkedWord
      }
      pageInfo {
        page
        total
      }
    }
  }
`;

export const DictionaryFuzzy = gql`
  query DictionaryFuzzy($input: SearchQueryInput!) {
    fuzzy(input: $input) {
      results {
        headword
        partOfSpeech
        normalized
        quickGlosses {
          language
          gloss
        }
        verb {
          principalParts
        }
        modernConnections {
          term
          note
        }
        definitions {
          grade
          meanings {
            definition
            language
          }
        }
        linkedWord
      }
      pageInfo {
        page
        total
      }
    }
  }
`;

export const DictionaryPhrase = gql`
  query DictionaryPhrase($input: SearchQueryInput!) {
    phrase(input: $input) {
      results {
        headword
        partOfSpeech
        normalized
        quickGlosses {
          language
          gloss
        }
        verb {
          principalParts
        }
        modernConnections {
          term
          note
        }
        definitions {
          grade
          meanings {
            definition
            language
          }
        }
        linkedWord
      }
      pageInfo {
        page
        total
      }
    }
  }
`;

export const CounterTopFive = gql`
  query CounterTopFive {
    counterTopFive {
      topFive {
        lastUsed
        serviceName
        word
        count
      }
    }
  }
`;

export const CounterSession = gql`
  query CounterSession($sessionId: String!) {
    counterSession(sessionId: $sessionId) {
      topFive {
        lastUsed
        serviceName
        word
        count
      }
    }
  }
`;