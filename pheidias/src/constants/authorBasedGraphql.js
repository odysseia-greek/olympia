import gql from 'graphql-tag';

export const AuthorBasedOptions = gql`
    query {
        authorBasedOptions {
            themes {
                name
                segments {
                    name
                    maxSet
                }
            }
        }
    }
`;

export const AuthorBasedCreateQuestion = gql`
    query authorBasedQuiz($input: AuthorBasedInput!) {
        authorBasedQuiz(input: $input) {
            reference
            translation
            fullSentence
            quiz {
                numberOfItems
                quizItem
                options {
                    quizWord
                }
            }
            grammarQuiz {
                correctAnswer
                extraInformation
                wordInText
                options {
                    quizWord
                }
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

export const AuthorBasedAnswer = gql`
    query authorBasedAnswer($input: AuthorBasedAnswerInput!) {
        authorBasedAnswer(input: $input) {
            quizWord
            correct
            finished
            wordsInText
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
`

export const AuthorBasedForms = gql`
    query authorBasedAnswer($input: AuthorBasedWordFormsInput!) {
        authorBasedWordForms(input: $input) {
            forms {
                dictionaryForm
                wordsInText
            }
        }
    }
`