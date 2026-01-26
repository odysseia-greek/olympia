<template>
  <v-container fluid>
    <v-parallax :src="quizHomeImage" height="75em" class="darkened-parallax">
      <div class="overlay"></div>
      <div class="d-flex flex-column justify-center align-center text-white parallax-content">
        <h1 class="text-h2 mb-2">Welcome to Quiz Home!</h1>

          <TypingText :texts="introTexts" />
      <v-icon class="mt-4 scroll-icon" size="36" @click="scrollMeTo('quizModeRef')">
        mdi-chevron-down
      </v-icon>
      <p class="scroll-indicator" style="text-align: center">
        Scroll down for Quiz Modes and additional info.
      </p>
      </div>
    </v-parallax>

    <v-container class="mt-10">
      <div ref="quizModeRef"></div>
      <v-row dense>
        <v-col v-for="mode in quizModes" :key="mode.id" cols="12" sm="6" md="4">
          <v-card :color="mode.color" height="100%" class="d-flex flex-column justify-space-between">
            <a :href="mode.route" style="color: white;">
              <v-img
                  :src="mode.src"
                  height="20em"
                  cover
                  gradient="to bottom, rgba(0,0,0,.1), rgba(0,0,0,.5)"
              >
                <v-row class="fill-height white--text align-end" style="margin: 0;">
                  <v-card-title v-text="mode.name" class="white--text"></v-card-title>
                </v-row>
              </v-img>
            </a>

            <v-card-subtitle class="text-subtitle-2 mt-2">
              <strong>{{ mode.greekName }}</strong>
            </v-card-subtitle>

            <v-card-text class="flex-grow-1">
              <p>{{ mode.description }}</p>
            </v-card-text>

            <v-card-actions>
              <v-btn color="primary" :to="mode.route" variant="text" v-if="mode.route">Start Quiz</v-btn>
              <v-chip v-else color="grey" variant="flat" label>Coming soon</v-chip>
            </v-card-actions>
          </v-card>

        </v-col>

      </v-row>
    </v-container>
  </v-container>
</template>
<script>
import {ref, onMounted, nextTick} from 'vue';
import { useRouter } from 'vue-router';
import SearchComponent from "@/components/SearchBar.vue";
import TypingText from "@/components/TypingText.vue";

export default {
  name: 'QuizAreaHome',
  components: {SearchComponent, TypingText},
  setup() {
    const router = useRouter();
    const quizModeRef = ref();
    const quizHomeImage = ref('');
    const images = import.meta.glob('/src/assets/*.webp');

    const introTexts = ref([
      {
        greek: "ἓν οἶδα ὅτι οὐδὲν οἶδα.",
        translation: "I know one thing that I know nothing."
      },
      {
        greek: "Πάντα ῥεῖ καὶ οὐδὲν μένει.",
        translation: "Everything flows and nothing stays."
      },
      {
        greek: "Ἀρχὴ ἥμισυ παντός.",
        translation: "The beginning is half of the whole."
      },
      {
        greek: "τό γάρ αυτο νοειν έστιν τε καί ειναι.",
        translation: "For it is the same thinking and being."
      },
      {
        greek: "Ου κλέπτω την νίκην.",
        translation: "I will not steal my victory."
      },
      {
        greek: "Αἶψα γὰρ ἐν κακότητι βροτοὶ καταγηράσκουσιν.",
        translation: "Hardship can age a person overnight."
      },
    ]);

    const quizModes = ref([
      {
        id: 1,
        color: "triadic",
        name: 'Media',
        greekName: 'Aristippos – Ἀρίστιππος',
        level: 'Beginner',
        description: `You are presented with a Greek word and need to match it to the correct image.
Audio support will be added in a later stage to help with pronunciation. This mode also offers a
'comprehensive mode', which searches for all declined forms of this word across known texts, helping you
build a deeper connection to its usage.`,
        image: 'playing_a_game.webp',
        src: '',
        route: 'quiz/media'
      },
      {
        id: 2,
        color: "triadic",
        name: 'Multiple Choice',
        greekName: 'Kritias – Κριτίας',
        level: 'Beginner',
        description: `Match the correct Greek word to its translation in English or Dutch. This mode is ideal for
rapid learning and recall. It also includes a comprehensive mode, where words are linked to their variations
in actual texts for an enhanced learning experience.`,
        image: 'alexander.webp',
        src: '',
        route: 'quiz/multiplechoice'
      },
      {
        id: 3,
        color: "triadic",
        name: 'Author Based',
        greekName: 'Xenofon – Ξενοφῶν',
        level: 'Intermediate',
        description: `You’ll work through real Ancient Greek texts by authors like Herodotos, Plato, and Plutarch.
Each sentence starts blurred out. As you correctly answer each word using multiple-choice (based on the dictionary form),
the corresponding Greek word becomes visible. Once all words are revealed, you'll see the full sentence along with its translation.
Each word is clickable, allowing you to explore its grammatical form, such as: λόγων → λόγος.
Some words may also require selecting the correct grammatical form.`,
        image: 'alexandria.webp',
        src: '',
        route: '/quiz/authorbased'
      },
      {
        id: 4,
        color: "triadic",
        name: 'Dialogue',
        greekName: 'Kriton – Κρίτων',
        level: 'Advanced',
        description: `Engage in reconstructed Ancient Greek dialogues between two or more characters.
Each dialogue consists of short, manageable sentences to prevent cognitive overload.
You’ll select a character and then reconstruct the correct order of the conversation. Translations can be toggled on for support,
and feedback is available along the way to guide your learning.`,
        image: 'sokrates.webp',
        src: '',
        route: '/quiz/dialogue'
      },
      {
        id: 5,
        color: "triadic",
        name: 'Grammar',
        greekName: 'Antisthenes – Ἀντισθένης',
        level: 'Intermediate',
        description: `A mode focused on grammar drills. You'll be able to identify, conjugate,
and decline words. This mode will test your understanding of Greek morphology, covering nouns, verbs, and participles.`,
        image: 'grammar.webp',
        src: '',
        route: '/quiz/grammar'
      },
      {
        id: 6,
        color: "triadic",
        name: 'Journey Mode',
        greekName: 'Alkibiades – Ἀλκιβιάδης (preview)',
        level: 'All Levels',
        description: `A guided experience that blends multiple quiz types into a structured path. These can be grammar based, media based,
        or historical background.
Each chapter will help you progress toward reading and understanding real Greek texts that are set in the context of a "theme" such as the Persian wars.
Great for learners who want to follow a more narrative and progressive route. This mode is still in development and will be updated in the future.
It does not work well on mobile devices, so we recommend using a laptop or tablet for best results.`,
        image: 'odysseus.webp',
        src: '',
        route: 'quiz/journey'
      }
    ]);

    const loadCardImages = () => {
      quizModes.value.forEach(card => {
        const importPath = `/src/assets/${card.image}`;
        if (images[importPath]) {
          images[importPath]().then((module) => {
            card.src = module.default;
          });
        }
      });
    };

    const loadImage = () => {
      import('@/assets/school_of_athens.webp').then((module) => {
        quizHomeImage.value = module.default;
      });
    };

    const scrollMeTo = (refName) => {
      nextTick(() => {
        if (refName === 'quizModeRef' && quizModeRef.value) {
          quizModeRef.value.scrollIntoView({ behavior: 'smooth' });
        }
      });
    };

    onMounted(() => {
      loadImage();
      loadCardImages();
    });

    return {
      quizHomeImage,
      quizModes,
      introTexts,
      quizModeRef,
      scrollMeTo
    };
  }
};
</script>

<style scoped>

/* Overlay to enhance contrast */
.overlay {
  position: absolute;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.3);
  top: 0;
  left: 0;
}

/* Ensures content is properly spaced inside the parallax */
.parallax-content {
   margin-top: 10em;
  position: relative;
  z-index: 1;
  text-align: center;
  padding: 5em 2em;
}

</style>
