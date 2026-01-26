import gql from 'graphql-tag';

export const MultipleChoiceOptions = gql`
    query {
        multipleChoiceOptions {
            themes {
                name
                maxSet
            }
        }
    }
`;

export const MultipleChoiceCreateQuiz = gql`
    query multipleChoiceQuiz($input: MultipleQuizInput!) {
        multipleChoiceQuiz(input: $input) {
            numberOfItems
            quizItem
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

export const MultipleChoiceAnswer = gql`
    query multipleChoiceAnswer($input: MultipleChoiceAnswerInput!) {
        multipleChoiceAnswer(input: $input) {
            correct
            quizWord
            finished
            similarWords {
                greek
                english
            }
            foundInText {
                rootword
                conjugations {
                    word
                    rule
                }
                texts {
                    author
                    reference
                    referenceLink
                    book
                    text {
                        translations
                        greek
                    }
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
