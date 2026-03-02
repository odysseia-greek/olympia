import gql from 'graphql-tag';

export const MediaOptions = gql`
    query {
        mediaOptions {
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

export const MediaCreateQuestion = gql`
    query mediaQuiz($input: MediaQuizInput!) {
        mediaQuiz(input: $input) {
            numberOfItems
            quizItem
            options {
                option
                imageUrl
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

export const MediaAnswer = gql`
    query mediaAnswer($input: MediaAnswerInput!) {
        mediaAnswer(input: $input) {
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
