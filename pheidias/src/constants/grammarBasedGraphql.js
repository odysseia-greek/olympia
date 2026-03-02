import gql from 'graphql-tag';

export const GrammarOptions = gql`
    query {
        grammarOptions {
            themes {
                name
                segments {
                    name
                    difficulty
                    maxSet
                }
            }
        }
    }
`;

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

export const GrammarAnswer = gql`
    query grammarAnswer($input: GrammarAnswerInput!) {
        grammarAnswer(input: $input) {
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
`
