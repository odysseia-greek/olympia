import gql from 'graphql-tag';

export const DialogueBasedOptions = gql`
    query {
        dialogueOptions {
            themes {
                name
                maxSet
            }
        }
    }
`;

export const DialogueBasedQuestion = gql`
    query dialogueQuiz($input: DialogueQuizInput!) {
        dialogueQuiz(input: $input) {
            dialogue {
                introduction
                section
                linkToPerseus
                speakers {
                    shorthand
                    translation
                }
            }
            content {
                translation
                greek
                place
                speaker
            }
        }
    }
`;

export const DialogueBasedAnswer = gql`
    query dialogueAnswer($input: DialogueAnswerInput!) {
        dialogueAnswer(input: $input) {
        percentage
        input {
            greek
            place
            speaker
            translation
        }
        answer {
            greek
            place
            speaker
            translation
        }
        wronglyPlaced {
            greek
            translation
            speaker
            place
        }
    }
    }
`
