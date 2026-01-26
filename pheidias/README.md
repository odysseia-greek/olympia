# Pheidias - Φειδίας

Pheidias (or Phidias) is one of the great artists of the Greek world, most famous for his work on the Athenian Akropolis. An apt choice for the frontend of the Odysseia-Greek application.

## About Odysseia-Greek

Odysseia-Greek is an interactive platform for learning Ancient (Attic) Greek. While there are many resources available for studying Ancient Greek, Odysseia-Greek aims to complement them by offering a unique, interactive experience.

### Why choose Odysseia-Greek?

- **Engaging Content:** All tools and exercises are designed to keep you motivated and entertained as you learn.
- **Easy to Use:** Each component is designed to be interactive and user-friendly, with examples and instructions.
- **Community Driven:** We value your input! If you spot any mistakes or have suggestions, don't hesitate to reach out.

### Learning Paths

Odysseia-Greek caters to learners of all levels:

- **Philosopher (Advanced)** - For experienced learners
- **Sophist (Intermediate)** - For those with some knowledge
- **Demagogue (Beginner)** - For newcomers to Ancient Greek

## Features

### Quiz (Sokrates)

The Quiz section offers various types of quizzes to practice your Ancient Greek:

- **Media Quizzes** - Image-based quizzes for learning vocabulary
- **Multiple Choice Quizzes** - Test your knowledge with multiple options
- **Author-based Quizzes** - Learn from famous Greek works
- **Dialogue Quizzes** - Practice through interactive dialogues
- **Grammar Quizzes** - Test your understanding of Greek grammar
- **Journey Quizzes** - Follow a structured learning path

### Texts (Herodotos)

The Texts section allows you to translate Greek texts and compare your translations with official ones. Named after Herodotos, the "father of history," this section helps you practice your translation skills.

### Grammar (Dionysios)

The Grammar section helps you with Greek grammar, including declensions and word forms. Named after Dionysios, a grammarian from Alexandria, this section supports nouns, verbs, and participia.

### Dictionary (Alexandros)

The Dictionary section allows you to search for Greek words in multiple languages (English, Ancient Greek, Dutch). Each keystroke sends a query to the backend, and accents are ignored for easier searching.

## Technologies

Pheidias is built with modern web technologies:

- **Vue 3** - Progressive JavaScript framework
- **Vite** - Next generation frontend tooling
- **Vuetify 3** - Material Design component framework
- **Apollo Client** - GraphQL client for data fetching
- **Pinia** - State management for Vue applications
- **Vue Router 4** - Official router for Vue.js

## Development

### Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run dev
```

### Compiles and minifies for production
```
npm run build
```

### Preview the production build
```
npm run serve
```

### Customize configuration
See [Vite Configuration Reference](https://vitejs.dev/config/).

## Deployment

### Docker

The application can be built and deployed using Docker:

```bash
# Build the Docker image
docker build -t odysseia-greek/pheidias .

# Run the container
docker run -p 80:80 odysseia-greek/pheidias
```

This will build the application and serve it using Nginx on port 80.
