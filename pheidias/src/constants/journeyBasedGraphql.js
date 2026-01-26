import gql from 'graphql-tag';

export const JourneyOptions = gql`
    query {
        journeyOptions {
            themes {
                name
                segments {
                    name
                    location
                    number
                    coordinates{
                        x
                        y
                    }
                }
            }
        }
    }
`
export const GrammarCreateQuestion = gql`
    query grammarQuiz($input: GrammarQuizInput!) {
        grammarQuiz(input: $input) {
            numberOfItems
            quizItem
            stem
            difficulty
            description
            contractionRule
            dictionaryForm
            translation
            options {
                option
            }
            progress {
                greek
                translation
                playCount
                correctCount
                incorrectCount
                lastPlayed
            }
        }
    }
`;

export const JourneyCreateQuestion = gql`
    query journeyQuiz($input: JourneyQuizInput) {
        journeyQuiz(input: $input) {
            segment
            theme
            number
            translation
            sentence
            contextNote
            intro {
                author
                background
                work
            }
            quiz {
                ... on MatchQuiz {
                    __typename
                    instruction
                    pairs {
                        greek
                        answer
                    }
                }
                ... on TriviaQuiz {
                    __typename
                    question
                    options
                    answer
                    note
                }
                ... on StructureQuiz {
                    __typename
                    title
                    text
                    question
                    options
                    answer
                    note
                }
                ... on MediaQuiz {
                    __typename
                    instruction
                    mediaFiles {
                        word
                        answer
                    }
                }
                ... on FinalTranslationQuiz {
                    __typename
                    instruction
                    options
                    answer
                }
            }
        }
    }
`