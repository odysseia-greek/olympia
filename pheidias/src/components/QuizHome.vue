<template>
  <v-container id="quiz-home" fluid>
    <v-parallax ref="quizHeroRef" :src="quizHomeImage" class="quiz-hero">
      <div class="quiz-hero-shade">
        <v-container class="quiz-hero-content">
          <div class="quiz-hero-copy">
            <h1>Sokrates</h1>
            <p>
              Choose a quiz mode for vocabulary, grammar, dialogue, image prompts, or guided reading.
            </p>
          </div>

          <nav class="quiz-quick-nav" aria-label="Quiz quick navigation">
            <v-btn
                v-for="link in quickLinks"
                :key="link.title"
                color="papyrus"
                variant="flat"
                @click="scrollMeTo(link.refName)"
            >
              <v-icon start>{{ link.icon }}</v-icon>
              {{ link.title }}
            </v-btn>
          </nav>

          <section class="quiz-hero-panel" aria-labelledby="quiz-start-heading">
            <div class="panel-heading">
              <div>
                <div class="section-label">Start by level</div>
                <h2 id="quiz-start-heading">Pick the pressure you want</h2>
              </div>
              <p>
                Begin with image or choice-based recall, move into author passages, or test yourself with dialogue.
              </p>
            </div>

            <v-row dense>
              <v-col v-for="start in startingPoints" :key="start.level" cols="12" md="4">
                <v-card class="start-card" :style="{'--start-color': start.color}" height="100%">
                  <v-card-text>
                    <v-avatar :color="start.color" size="42">
                      <v-icon color="white">{{ start.icon }}</v-icon>
                    </v-avatar>
                    <div class="start-level">{{ start.level }}</div>
                    <h3>{{ start.title }}</h3>
                    <p>{{ start.text }}</p>
                  </v-card-text>
                  <v-card-actions>
                    <v-btn :to="start.route" :color="start.color" variant="flat" block>
                      Start {{ start.mode }}
                      <v-icon end>mdi-arrow-right</v-icon>
                    </v-btn>
                  </v-card-actions>
                </v-card>
              </v-col>
            </v-row>
          </section>
        </v-container>
      </div>
    </v-parallax>

    <main class="quiz-content">
      <v-container class="content-container">
        <section ref="snippetsRef" class="snippet-section">
          <v-row dense>
            <v-col cols="12" md="5">
              <div class="section-heading">
                <h2>Warm up before choosing a mode</h2>
              </div>
            </v-col>
            <v-col cols="12" md="7">
              <v-sheet class="snippet-panel" color="text">
                <TypingText :texts="introTexts" controls />
              </v-sheet>
            </v-col>
          </v-row>
        </section>

        <section ref="modesRef" class="modes-section" aria-labelledby="quiz-modes-heading">
          <div class="section-heading modes-heading">
            <div>
              <v-chip color="secondary" variant="tonal">Quiz modes</v-chip>
              <h2 id="quiz-modes-heading">Choose how you want to practise</h2>
            </div>
            <p>
              Each mode has a different rhythm: quick recall, contextual reading, grammar drills, or structured dialogue.
            </p>
          </div>

          <v-row dense>
            <v-col v-for="mode in quizModes" :key="mode.id" cols="12" md="6" lg="4">
              <v-card class="mode-card" height="100%">
                <v-img
                    :src="mode.src"
                    class="mode-image"
                    cover
                    gradient="to bottom, rgba(0,0,0,.02), rgba(0,0,0,.72)"
                >
                  <div class="mode-image-content">
                    <v-chip :color="mode.levelColor" size="small" variant="flat">{{ mode.level }}</v-chip>
                    <h3>{{ mode.name }}</h3>
                  </div>
                </v-img>

                <v-card-text>
                  <div class="mode-greek-name">{{ mode.greekName }}</div>
                  <p>{{ mode.description }}</p>
                </v-card-text>

                <v-card-actions>
                  <v-btn color="primary" :to="mode.route" variant="text" v-if="mode.route">
                    Start {{ mode.name }}
                    <v-icon end>mdi-arrow-right</v-icon>
                  </v-btn>
                  <v-chip v-else color="grey" variant="flat" label>Coming soon</v-chip>
                </v-card-actions>
              </v-card>
            </v-col>
          </v-row>
        </section>
      </v-container>
    </main>
  </v-container>
</template>

<script>
import {ref, onMounted, nextTick} from 'vue';
import TypingText from "@/components/TypingText.vue";

export default {
  name: 'QuizAreaHome',
  components: {TypingText},
  setup() {
    const quizHeroRef = ref();
    const snippetsRef = ref();
    const modesRef = ref();
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
        name: 'Media',
        greekName: 'Aristippos - Ἀρίστιππος',
        level: 'Beginner',
        levelColor: '#cd7f32',
        description: `You are presented with a Greek word and need to match it to the correct image. Audio support will be added in a later stage to help with pronunciation. Comprehensive mode searches for all declined forms of this word across known texts, helping you build a deeper connection to its usage.`,
        image: 'playing_a_game.webp',
        src: '',
        route: 'quiz/media'
      },
      {
        id: 2,
        name: 'Multiple Choice',
        greekName: 'Kritias - Κριτίας',
        level: 'Beginner',
        levelColor: '#cd7f32',
        description: `Match the correct Greek word to its translation in English or Dutch. This mode is built for rapid recall and also includes comprehensive mode, where words are linked to variations in actual texts.`,
        image: 'alexander.webp',
        src: '',
        route: 'quiz/multiplechoice'
      },
      {
        id: 3,
        name: 'Author Based',
        greekName: 'Xenofon - Ξενοφῶν',
        level: 'Intermediate',
        levelColor: '#9ea7ad',
        description: `Work through real Ancient Greek texts by authors like Herodotos, Plato, and Plutarch. Each sentence starts blurred out and becomes visible as you answer words correctly. Clickable words let you inspect grammatical forms.`,
        image: 'alexandria.webp',
        src: '',
        route: '/quiz/authorbased'
      },
      {
        id: 4,
        name: 'Dialogue',
        greekName: 'Kriton - Κρίτων',
        level: 'Advanced',
        levelColor: '#d4af37',
        description: `Reconstruct Ancient Greek dialogues between two or more characters. Select a character, order the conversation, toggle translations when needed, and use feedback to tighten your reading.`,
        image: 'sokrates.webp',
        src: '',
        route: '/quiz/dialogue'
      },
      {
        id: 5,
        name: 'Grammar',
        greekName: 'Antisthenes - Ἀντισθένης',
        level: 'Intermediate',
        levelColor: '#9ea7ad',
        description: `Drill morphology by identifying, conjugating, and declining words. This mode focuses on nouns, verbs, and participles so grammar practice stays close to reading practice.`,
        image: 'grammar.webp',
        src: '',
        route: '/quiz/grammar'
      },
      {
        id: 6,
        name: 'Journey Mode',
        greekName: 'Alkibiades - Ἀλκιβιάδης (preview)',
        level: 'All Levels',
        levelColor: '#1cbcd1',
        description: `A guided path that blends multiple quiz types with historical background. Each chapter helps you progress toward reading real Greek texts around a theme. This preview works best on laptop or tablet.`,
        image: 'odysseus.webp',
        src: '',
        route: 'quiz/journey'
      }
    ]);

    const quickLinks = ref([
      {
        title: 'Snippets',
        icon: 'mdi-format-quote-open',
        refName: 'snippetsRef',
      },
      {
        title: 'Modes',
        icon: 'mdi-view-dashboard-outline',
        refName: 'modesRef',
      },
    ]);

    const startingPoints = ref([
      {
        level: 'Beginner',
        title: 'Build recall',
        mode: 'Media',
        icon: 'mdi-image-search-outline',
        color: '#cd7f32',
        route: 'quiz/media',
        text: 'Start with visual prompts or multiple choice when you want quick repetition.',
      },
      {
        level: 'Intermediate',
        title: 'Read in context',
        mode: 'Author Based',
        icon: 'mdi-book-open-page-variant',
        color: '#9ea7ad',
        route: '/quiz/authorbased',
        text: 'Move into author passages when you want vocabulary and grammar tied to real sentences.',
      },
      {
        level: 'Advanced',
        title: 'Reconstruct meaning',
        mode: 'Dialogue',
        icon: 'mdi-account-question-outline',
        color: '#d4af37',
        route: '/quiz/dialogue',
        text: 'Use dialogue mode when you are ready to test order, response, and comprehension.',
      },
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
        if (refName === 'quizHeroRef' && quizHeroRef.value) {
          quizHeroRef.value.$el.scrollIntoView({behavior: 'smooth', block: 'start'});
        }
        if (refName === 'snippetsRef' && snippetsRef.value) {
          snippetsRef.value.scrollIntoView({behavior: 'smooth', block: 'start'});
        }
        if (refName === 'modesRef' && modesRef.value) {
          modesRef.value.scrollIntoView({behavior: 'smooth', block: 'start'});
        }
      });
    };

    onMounted(() => {
      loadImage();
      loadCardImages();
    });

    return {
      quizHeroRef,
      snippetsRef,
      modesRef,
      quizHomeImage,
      quizModes,
      quickLinks,
      startingPoints,
      introTexts,
      scrollMeTo
    };
  }
};
</script>

<style scoped>
#quiz-home {
  --quiz-primary: #1c61d1;
  --quiz-secondary: #1cd18c;
  --quiz-triadic: #1cbcd1;
  --quiz-ink: #20334f;
  --quiz-muted: #536987;
  padding: 0;
  color: var(--quiz-ink);
}

.quiz-hero,
.quiz-hero-shade {
  min-height: calc(100vh - 64px);
}

.quiz-hero-shade {
  display: flex;
  align-items: center;
  background:
      linear-gradient(90deg, rgba(22, 20, 15, 0.82) 0%, rgba(38, 31, 21, 0.55) 48%, rgba(26, 22, 16, 0.18) 100%),
      linear-gradient(180deg, rgba(0, 0, 0, 0.08), rgba(26, 21, 14, 0.44));
}

.quiz-hero-content {
  display: grid;
  gap: 26px;
  padding-top: 72px;
  padding-bottom: 38px;
}

.quiz-hero-copy {
  max-width: 820px;
  color: #fff;
  text-shadow: 0 2px 14px rgba(0, 0, 0, 0.46);
}

.hero-kicker {
  margin-bottom: 18px;
  color: var(--quiz-ink);
  font-weight: 800;
}

.quiz-hero-copy h1 {
  margin: 0;
  font-size: clamp(3rem, 8vw, 6rem);
  line-height: 0.96;
  letter-spacing: 0;
}

.quiz-hero-copy p {
  max-width: 680px;
  margin: 20px 0 0;
  font-size: clamp(1.1rem, 2.4vw, 1.45rem);
  line-height: 1.55;
}

.quiz-quick-nav {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.quiz-quick-nav .v-btn {
  color: var(--quiz-ink);
  min-height: 42px;
}

.quiz-hero-panel {
  max-width: 1120px;
  padding: 18px;
  border: 1px solid rgba(28, 188, 209, 0.32);
  border-radius: 8px;
  background: rgba(253, 246, 227, 0.96);
  box-shadow: 0 18px 48px rgba(11, 39, 85, 0.3);
  backdrop-filter: blur(8px);
}

.panel-heading,
.section-heading {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 28px;
  margin-bottom: 22px;
}

.panel-heading p,
.section-heading p {
  max-width: 520px;
  margin: 0;
  color: #344765;
  line-height: 1.65;
}

.section-label,
.start-level,
.mode-greek-name {
  color: #64789e;
  font-size: 0.82rem;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.panel-heading h2,
.section-heading h2 {
  margin: 8px 0 0;
  font-size: clamp(1.6rem, 3vw, 2.35rem);
  line-height: 1.15;
  letter-spacing: 0;
}

.start-card,
.mode-card,
.snippet-panel {
  border: 1px solid rgba(28, 97, 209, 0.16);
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(28, 97, 209, 0.1);
}

.start-card {
  border-color: color-mix(in srgb, var(--start-color) 34%, transparent);
}

.start-card .v-avatar {
  margin-bottom: 16px;
}

.start-card h3,
.mode-card h3 {
  margin: 7px 0 10px;
  color: var(--quiz-ink);
}

.start-card p,
.mode-card p {
  margin: 0;
  color: #344765;
  line-height: 1.6;
}

.start-card .v-btn,
.mode-card .v-btn {
  min-height: 44px;
  white-space: normal;
}

.quiz-content {
  background:
      linear-gradient(150deg, rgba(28, 97, 209, 0.16) 0%, rgba(28, 188, 209, 0.12) 34%, rgba(28, 209, 140, 0.1) 62%, rgba(254, 252, 245, 0.98) 100%),
      linear-gradient(180deg, #d5eff7 0%, #f2fbf7 46%, #fefcf5 100%);
}

.content-container {
  max-width: 1240px;
  padding-top: 54px;
  padding-bottom: 58px;
}

.snippet-section,
.modes-section,
.quiz-hero {
  scroll-margin-top: 80px;
}

.snippet-section {
  margin-bottom: 46px;
}

.snippet-panel {
  min-height: 260px;
  padding: 24px;
  background: rgba(255, 255, 255, 0.9);
}

.modes-heading {
  margin-bottom: 24px;
}

.mode-card {
  overflow: hidden;
  transition: transform 160ms ease, box-shadow 160ms ease;
}

.mode-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 12px 26px rgba(57, 75, 107, 0.22);
}

.mode-image {
  height: 238px;
}

.mode-image-content {
  display: flex;
  height: 100%;
  flex-direction: column;
  justify-content: end;
  align-items: start;
  gap: 8px;
  padding: 22px;
  color: #fff;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.55);
}

.mode-image-content h3 {
  margin: 0;
  color: #fff;
  font-size: 2rem;
  line-height: 1;
  letter-spacing: 0;
}

.mode-greek-name {
  margin-bottom: 12px;
}

@media (max-width: 900px) {
  .quiz-hero,
  .quiz-hero-shade {
    min-height: auto;
  }

  .quiz-hero-content {
    padding-top: 46px;
    padding-bottom: 28px;
  }

  .quiz-hero-copy h1 {
    font-size: clamp(2.5rem, 13vw, 4.4rem);
  }

  .panel-heading,
  .section-heading {
    display: block;
  }

  .panel-heading p,
  .section-heading p {
    margin-top: 10px;
  }

  .mode-image {
    height: 220px;
  }
}
</style>
