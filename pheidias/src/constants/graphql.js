import gql from "graphql-tag";

export const HerodotosOptions = gql`
    query textOptions {
        textOptions {
            authors {
                key
                books {
                    key
                    references {
                        key
                        sections {
                            key
                        }
                    }
                }
            }
        }
    }
`

export const HerodotosCreate = gql`
    query create($input: CreateTextInput!) {
        create(input: $input) {
            author
            book
            reference
            rhemai {
                greek
                section
                translations
            }
        }
    }
`

export const HerodotosCheck = gql`
    query check($input: CheckTextRequestInput!) {
        check(input: $input) {
            averageLevenshteinPercentage
            sections {
                section
                answerSentence
                quizSentence
                levenshteinPercentage
            }
            possibleTypos {
                source
                provided
            }
        }
    }
`
