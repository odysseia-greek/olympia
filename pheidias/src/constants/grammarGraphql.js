import gql from 'graphql-tag';

export const Analyze = gql`
  query Analyze($rootword: String) {
    analyze(rootword: $rootword) {
      rootword
      conjugations {
        word
        rule
      }
      texts {
        text {
          greek
          translations
          section
        }
        referenceLink
        author
        book
        reference
      }
    }
  }
`;

export const CheckGrammar = gql`
  query Grammar($word: String, $includeAudit: Boolean) {
    grammar(word: $word, includeAudit: $includeAudit) {
      results {
        translations
        word
        rule
        rootWord
      }
      audit {
        requestId
        word
        outcome
        source
        reason
        events {
          step
          status
          reason
          source
          rule
          rootWord
          searchTerm
          resultCount
          candidateCount
          details
        }
      }
    }
  }
`;

export const AnalyzeSentence = gql`
  query Sentence($text: String!, $includeAudit: Boolean) {
    sentence(text: $text, includeAudit: $includeAudit) {
      sessionId
      originalText
      literalTranslation
      rateLimit {
        upstreamIp
        windowSeconds
        nextAllowedTime
      }
      tokens {
        token
        position
        gloss
        resolved
        message
        results {
          rootWord
          rule
          translations
          word
        }
      }
      textSearch {
        searched
        found
        status
        message
        matchCount
        query
        matches {
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
      }
    }
  }
`;
