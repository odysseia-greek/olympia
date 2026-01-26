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

export const Analyze = gql`
    query analyze($rootword: String!) {
        analyze(rootword: $rootword) {
            rootword
            conjugations{
                word
                rule
            }
            results {
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

export const CheckGrammar = gql`
query grammar($word: String!) {
    grammar(word: $word) {
        results{
            translations
            word
            rule
            rootWord
        }
    }
}`

export const DictionaryEntry = gql`
    query dictionary($word: String!, $language: String!, $mode: String!, $searchInText: Boolean!) {
        dictionary(word: $word, language: $language, mode: $mode, searchInText: $searchInText) {
            hits {
                hit{
                    greek
                    english
                    dutch
                    original
                    linkedWord
                }
                foundInText{
                    rootword
                    conjugations{
                        word
                        rule
                    }
                    texts{
                        reference
                        referenceLink
                        author
                        book
                        text{
                            translations
                            greek
                        }
                    }
                }
        }
        }
}`
