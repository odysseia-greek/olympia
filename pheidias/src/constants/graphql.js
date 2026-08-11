import gql from "graphql-tag";

export const HerodotosOptions = gql`
    query CorpusOptions {
        corpusOptions {
            authors {
                name
                books {
                    name
                    references {
                        name
                        sections
                    }
                }
            }
        }
    }
`

export const HerodotosCreate = gql`
    query CorpusText($input: TextInput!) {
        corpusText(input: $input) {
            author
            book
            type
            reference
            perseusTextLink
            passages {
                greek
                section
                translations
            }
        }
    }
`

export const HerodotosCheck = gql`
    query CheckText($input: CheckTextInput!) {
        checkText(input: $input) {
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

export const HerodotosChapterOptions = gql`
    query ChapterOptions {
        chapterOptions {
            chapters {
                chapter
                title
                order
                level
            }
        }
    }
`

export const HerodotosChapter = gql`
    query Chapter($chapter: String!) {
        chapter(chapter: $chapter) {
            chapter
            title
            description
            context
            order
            level
            grammar {
                grammar
                title
                explanation
                example { greek translation note }
            }
            vocabulary { greek translation }
            texts {
                text
                title
                type
                greek
                readingHints
                source { author work reference dialect }
            }
        }
    }
`

export const HerodotosCheckChapter = gql`
    query CheckChapter($input: CheckChapterInput!) {
        checkChapter(input: $input) {
            chapter
            texts { text sourceText actualText learnerText }
        }
    }
`
