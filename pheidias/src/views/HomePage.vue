<template>
  <v-app id="homepage">
    <div class="homepage-shell">
      <v-parallax ref="interactiveSectionRef" :src="odysseus" class="home-hero">
        <div class="hero-shade">
          <v-container class="hero-content">
            <div class="hero-copy">
              <h1>Welcome to your Odysseia</h1>
              <p>
                Practise vocabulary, translate texts, explore grammar, and search Greek words from one place.
              </p>
            </div>

            <nav class="quick-nav" aria-label="Quick navigation">
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

            <section class="journey-panel" aria-labelledby="journey-heading">
              <div class="journey-heading">
                <h2 id="journey-heading">Choose a starting point</h2>
                <span>Pick the path that matches your current level.</span>
              </div>

              <v-row class="journey-grid" dense>
                <v-col v-for="option in journeyOptions" :key="option.key" cols="12" md="4">
                  <v-card class="journey-card">
                    <v-card-text>
                      <div class="journey-card-top">
                        <v-avatar :color="option.color" size="42">
                          <v-icon color="white">{{ option.icon }}</v-icon>
                        </v-avatar>
                        <div>
                          <div class="journey-level">{{ option.level }}</div>
                          <h3>{{ option.title }}</h3>
                        </div>
                      </div>
                      <p>{{ option.description }}</p>
                    </v-card-text>
                    <v-card-actions>
                      <v-btn :color="option.color" :to="option.to" block variant="flat">
                        Start {{ option.level }}
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

      <main ref="exampleSectionRef" class="home-content">
        <v-container class="content-container">
          <section ref="whySectionRef" class="intro-section">
            <div class="intro-copy">
              <v-chip class="section-chip" color="primary" variant="tonal">Why Odysseia</v-chip>
              <h2>Your gateway to Ancient Attic Greek</h2>
              <p>
                Odysseia-Greek complements traditional books and classes with short, interactive practice:
                quizzes, translation checks, grammar helpers, and dictionary lookup.
              </p>
            </div>

            <v-row class="feature-row" dense>
              <v-col v-for="feature in featureHighlights" :key="feature.title" cols="12" md="4">
                <div class="feature-tile">
                  <v-icon :color="feature.color" size="30">{{ feature.icon }}</v-icon>
                  <h3>{{ feature.title }}</h3>
                  <p>{{ feature.text }}</p>
                </div>
              </v-col>
            </v-row>

            <v-row class="demo-row" dense>
              <v-col cols="12" md="6">
                <v-sheet class="demo-panel quote-panel" color="text">
                  <div class="panel-label">
                    <v-icon size="20">mdi-format-quote-open</v-icon>
                    Greek Quotes
                  </div>
                  <TypingText :texts="introTexts" controls />
                </v-sheet>
              </v-col>
              <v-col cols="12" md="6">
                <v-sheet class="demo-panel" color="text">
                  <div class="panel-label">
                    <v-icon size="20">mdi-magnify</v-icon>
                    Try the dictionary
                  </div>
                  <SearchComponent />
                </v-sheet>
              </v-col>
            </v-row>
          </section>

          <section ref="toolsSectionRef" class="tools-section" aria-labelledby="tools-heading">
            <div class="section-heading">
              <div>
                <v-chip color="secondary" variant="tonal">Components</v-chip>
                <h2 id="tools-heading">Explore the learning components</h2>
              </div>
              <p>
                Move between guided quiz practice, real text translation, grammar analysis, and dictionary search.
              </p>
            </div>

            <v-row dense>
              <v-col v-for="card in cards" :key="card.title" cols="12" md="6">
                <v-card class="tool-card" height="100%">
                  <v-img
                      :src="card.src"
                      class="tool-image"
                      cover
                      gradient="to bottom, rgba(0,0,0,.05), rgba(0,0,0,.72)"
                  >
                    <div class="tool-image-content">
                      <v-icon color="white" size="34">{{ card.icon }}</v-icon>
                      <h3>{{ card.title }}</h3>
                    </div>
                  </v-img>

                  <v-card-text>
                    <div class="tool-subtitle">{{ card.subTitle }}</div>
                    <p>{{ card.shortText }}</p>
                    <p>{{ card.longText }}</p>
                  </v-card-text>
                  <v-card-actions>
                    <v-btn color="primary" variant="text" :to="card.link">
                      Open {{ card.title }}
                      <v-icon end>mdi-arrow-right</v-icon>
                    </v-btn>
                  </v-card-actions>
                </v-card>
              </v-col>
            </v-row>
          </section>
        </v-container>
      </main>
    </div>
  </v-app>
</template>

<script>
import {ref, onMounted, onBeforeUnmount, nextTick} from 'vue';
import TypingText from "@/components/TypingText.vue";
import SearchComponent from "@/components/SearchBar.vue";

export default {
  name: "HomePage",
  components: {SearchComponent, TypingText},
  setup() {
    const odysseus = ref('');
    const exampleSectionRef = ref();
    const interactiveSectionRef = ref();
    const whySectionRef = ref();
    const toolsSectionRef = ref();
    const images = import.meta.glob('/src/assets/*.webp');

    const journeyOptions = ref([
      {
        key: 'demagogue',
        level: 'Beginner',
        title: 'Demagogue',
        icon: 'mdi-sprout',
        color: '#cd7f32',
        to: {
          name: 'QuizMedia',
          query: {theme: 'Basic', segment: 'First Words 1'},
        },
        description: 'Start with visual prompts and first words before moving into longer texts.',
      },
      {
        key: 'sophist',
        level: 'Intermediate',
        title: 'Sophist',
        icon: 'mdi-bookshelf',
        color: '#9ea7ad',
        to: {
          name: 'QuizAuthorBased',
          query: {theme: 'Herodotos - Histories', segment: '1.1.0'},
        },
        description: 'Practise recognisable passages and build confidence with author-based context.',
      },
      {
        key: 'philosopher',
        level: 'Advanced',
        title: 'Philosopher',
        icon: 'mdi-account-question',
        color: '#d4af37',
        to: {
          name: 'QuizDialogue',
          query: {theme: 'Plato - Euthyphro'},
        },
        description: 'Step into dialogue mode and test meaning through back-and-forth Greek responses.',
      },
    ]);

    const featureHighlights = ref([
      {
        title: 'Guided starts',
        icon: 'mdi-compass-outline',
        color: 'primary',
        text: 'The first choice sends you straight to a useful exercise instead of a generic menu.',
      },
      {
        title: 'Practical feedback',
        icon: 'mdi-check-decagram-outline',
        color: 'secondary',
        text: 'Translation and quiz modes help you compare, correct, and repeat small chunks.',
      },
      {
        title: 'Reference nearby',
        icon: 'mdi-text-search',
        color: 'triadic',
        text: 'Dictionary and grammar tools stay close to the reading and quiz workflow.',
      },
    ]);

    const quickLinks = ref([
      {
        title: 'Why Odysseia',
        icon: 'mdi-compass-outline',
        refName: 'whySectionRef',
      },
      {
        title: 'Components',
        icon: 'mdi-view-dashboard-outline',
        refName: 'toolsSectionRef',
      },
    ]);

    const loadHeroImage = () => {
      const heroImport = window.innerWidth <= 900
          ? import('@/assets/mobile_odysseus.webp')
          : import('@/assets/odysseus.webp');

      heroImport.then((module) => {
        odysseus.value = module.default;
      });
    };

    const loadCardImages = () => {
      cards.value.forEach(card => {
        const importPath = `/src/assets/${card.image}`;
        if (images[importPath]) {
          images[importPath]().then((module) => {
            card.src = module.default;
          });
        }
      });
    };

    const scrollMeTo = (refName) => {
      nextTick(() => {
        if (refName === 'exampleSectionRef' && exampleSectionRef.value) {
          exampleSectionRef.value.scrollIntoView({behavior: 'smooth', block: 'start'});
        }
        if (refName === 'interactiveSectionRef' && interactiveSectionRef.value) {
          interactiveSectionRef.value.$el.scrollIntoView({behavior: 'smooth', block: 'start'});
        }
        if (refName === 'whySectionRef' && whySectionRef.value) {
          whySectionRef.value.scrollIntoView({behavior: 'smooth', block: 'start'});
        }
        if (refName === 'toolsSectionRef' && toolsSectionRef.value) {
          toolsSectionRef.value.scrollIntoView({behavior: 'smooth', block: 'start'});
        }
      });
    };

    const introTexts = ref([
      {
        greek: "ἓν οἶδα ὅτι οὐδὲν οἶδα.",
        translation: "I know one thing that I know nothing."
      },
      {
        greek: "νόμωι (γάρ φησι) γλυκὺ καὶ νόμωι πικρόν, νόμωι θερμόν, νόμωι ψυχρόν, νόμωι χροιή, ἐτεῆι δὲ ἄτομα καὶ κενόν",
        translation: "By convention sweet is sweet, bitter is bitter, hot is hot, cold is cold, color is color; but in truth there are only atoms and the void."
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
        greek: "οἷον ἡ ψυχή ἡ ἡμετέρα ἀὴρ οὖσα συγκρατεῖ ἡμᾶς, καὶ ὅλον τὸν κόσμον πνεῦμα καὶ ἀὴρ περιέχει",
        translation: "Just as our soul, being air, constrains us, so breath and air envelops the whole kosmos."
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
      {
        greek: "κακῶς φρονοῦντες· ὡς τρὶς ἂν παρ’ ἀσπίδα στῆναι θέλοιμ’ ἂν μᾶλλον ἢ τεκεῖν ἅπαξ.",
        translation: "How wrong they are! I would rather stand three times with a shield in battle than give birth once."
      },
    ]);

    const cards = ref([
      {
        title: "Quiz",
        icon: "mdi-alphabet-greek",
        src: '',
        image: 'school_of_athens.webp',
        link: "/quiz",
        subTitle: "Sokrates",
        shortText: "Imagine standing in the Athenian Agora, just selling your wares. A man, an ugly man at that, walks up and strikes up a conversation. Before you know it you need to answer questions about morality and how you know you are correct.",
        longText: "That insistent questioning is why Sokrates is the name given to the Quiz section. Practise multiple choice options from Greek works, image-based quizzes, and dialogue. It is the best place to start if you are new.",
      },
      {
        title: "Texts",
        icon: "mdi-bookshelf",
        src: '',
        image: 'alexandria.webp',
        link: "/texts",
        subTitle: "Herodotos",
        shortText: "Known as the father of history, Herodotos came from Halicarnassus in what is now Western Turkey, but was then a thriving Greek city.",
        longText: "In the Herodotos or Text section of Odysseia you can translate texts and compare your entry with an official translation. The score helps you see where your wording is close and where it needs work.",
      },
      {
        title: "Grammar",
        icon: "mdi-feather",
        src: '',
        image: 'grammar.webp',
        link: "/grammar",
        subTitle: "Dionysios",
        shortText: "As a grammarian working in Alexandria, Dionysios influenced how generations learned Greek grammar.",
        longText: "Dionysios helps you decline words and search for possible meanings. The backend does the heavier analysis, with current support for nouns, verbs, and participia.",
      },
      {
        title: "Dictionary",
        icon: "mdi-magnify",
        src: '',
        image: 'alexander.webp',
        link: "/dictionary",
        subTitle: "Alexandros",
        shortText: "Few people loom as large as Alexander the Great in world history, and his reach makes him a fitting guide for search.",
        longText: "Search partial words in English, Ancient Greek, or Dutch. Accents are ignored to keep lookup fast when you are reading or checking an exercise.",
      },
    ]);

    onMounted(() => {
      loadHeroImage();
      loadCardImages();
      window.addEventListener('resize', loadHeroImage);
    });

    onBeforeUnmount(() => {
      window.removeEventListener('resize', loadHeroImage);
    });

    return {
      odysseus,
      cards,
      featureHighlights,
      introTexts,
      journeyOptions,
      quickLinks,
      exampleSectionRef,
      interactiveSectionRef,
      whySectionRef,
      toolsSectionRef,
      scrollMeTo,
    };
  },
};
</script>

<style>
#homepage {
  --og-primary: #1c61d1;
  --og-secondary: #1cd18c;
  --og-triadic: #1cbcd1;
  --og-ink: #20334f;
  --og-muted: #536987;
  color: #394b6b;
}

.homepage-shell {
  background:
      linear-gradient(150deg, rgba(28, 97, 209, 0.2) 0%, rgba(28, 188, 209, 0.14) 28%, rgba(28, 209, 140, 0.12) 54%, rgba(254, 252, 245, 0.96) 100%),
      linear-gradient(180deg, #d5eff7 0%, #f2fbf7 46%, #fefcf5 100%);
}

.home-hero {
  min-height: calc(100vh - 64px);
}

.hero-shade {
  min-height: calc(100vh - 64px);
  display: flex;
  align-items: center;
  background:
      linear-gradient(90deg, rgba(22, 20, 15, 0.78) 0%, rgba(38, 31, 21, 0.5) 48%, rgba(26, 22, 16, 0.18) 100%),
      linear-gradient(180deg, rgba(0, 0, 0, 0.08), rgba(26, 21, 14, 0.42));
}

.hero-content {
  display: grid;
  gap: 28px;
  padding-top: 72px;
  padding-bottom: 36px;
}

.hero-copy {
  max-width: 760px;
  color: #fff;
  text-shadow: 0 2px 14px rgba(0, 0, 0, 0.45);
}

.hero-kicker {
  margin-bottom: 18px;
  color: #394b6b;
  font-weight: 700;
}

.hero-copy h1 {
  max-width: 820px;
  margin: 0;
  font-size: clamp(2.75rem, 7vw, 5.8rem);
  line-height: 0.98;
  letter-spacing: 0;
}

.hero-copy p {
  max-width: 680px;
  margin: 20px 0 0;
  font-size: clamp(1.1rem, 2.4vw, 1.5rem);
  line-height: 1.55;
}

.quick-nav {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.quick-nav .v-btn {
  color: var(--og-ink);
  min-height: 42px;
}

.journey-panel {
  max-width: 1120px;
  padding: 18px;
  border: 1px solid rgba(28, 188, 209, 0.38);
  border-radius: 8px;
  background: rgba(253, 246, 227, 0.96);
  box-shadow: 0 18px 48px rgba(11, 39, 85, 0.32);
  backdrop-filter: blur(8px);
}

.journey-heading {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 14px;
  color: #394b6b;
}

.journey-heading h2,
.section-heading h2,
.intro-copy h2 {
  margin: 8px 0 0;
  font-size: clamp(1.6rem, 3vw, 2.35rem);
  line-height: 1.15;
  letter-spacing: 0;
}

.journey-heading span {
  color: #64789e;
}

.journey-card {
  height: 100%;
  border: 1px solid rgba(89, 70, 37, 0.18);
  border-radius: 8px;
  transition: transform 160ms ease, box-shadow 160ms ease;
}

.journey-card:hover,
.tool-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 12px 26px rgba(57, 75, 107, 0.22);
}

.journey-card-top {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 14px;
}

.journey-card h3 {
  margin: 0;
  font-size: 1.25rem;
  color: var(--og-ink);
}

.journey-level,
.tool-subtitle,
.panel-label {
  color: #64789e;
  font-size: 0.82rem;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.journey-card p,
.feature-tile p,
.tool-card p,
.intro-copy p,
.section-heading p {
  color: #344765;
  line-height: 1.65;
}

.journey-card .v-btn,
.tool-card .v-btn {
  min-height: 44px;
  white-space: normal;
}

.section-chip {
  flex: 0 0 auto;
  max-width: 100%;
}

.scroll-cta {
  justify-self: start;
}

.home-content {
  scroll-margin-top: 80px;
}

.home-hero,
.intro-section,
.tools-section {
  scroll-margin-top: 80px;
}

.content-container {
  max-width: 1240px;
  padding-top: 54px;
  padding-bottom: 56px;
}

.intro-section,
.tools-section {
  margin-bottom: 42px;
}

.intro-copy,
.section-heading {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 28px;
  margin-bottom: 22px;
}

.intro-copy > p,
.section-heading > p {
  max-width: 520px;
  margin: 0;
}

.feature-row,
.demo-row {
  margin-top: 18px;
}

.feature-tile,
.demo-panel {
  height: 100%;
  border: 1px solid rgba(28, 97, 209, 0.18);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.88);
  box-shadow: 0 8px 24px rgba(28, 97, 209, 0.12);
}

.feature-tile {
  padding: 22px;
}

.feature-tile h3 {
  margin: 12px 0 8px;
  color: #394b6b;
}

.feature-tile p {
  margin: 0;
}

.demo-panel {
  padding: 22px;
  overflow: hidden;
}

.quote-panel {
  min-height: 230px;
}

.panel-label {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 18px;
}

.tool-card {
  overflow: hidden;
  border-radius: 8px;
  transition: transform 160ms ease, box-shadow 160ms ease;
}

.tool-image {
  height: 260px;
}

.tool-image-content {
  display: flex;
  height: 100%;
  flex-direction: column;
  justify-content: end;
  gap: 8px;
  padding: 22px;
  color: #fff;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.55);
}

.tool-image-content h3 {
  margin: 0;
  font-size: 2rem;
  line-height: 1;
  letter-spacing: 0;
}

.tool-subtitle {
  margin-bottom: 10px;
}

#typing-text {
  color: #394b6b;
}

@media (max-width: 900px) {
  .home-hero,
  .hero-shade {
    min-height: auto;
  }

  .hero-content {
    padding-top: 46px;
    padding-bottom: 28px;
  }

  .hero-copy h1 {
    font-size: clamp(2.3rem, 12vw, 4.2rem);
  }

  .journey-heading,
  .intro-copy,
  .section-heading {
    display: block;
  }

  .journey-heading span,
  .intro-copy > p,
  .section-heading > p {
    display: block;
    margin-top: 10px;
  }

  .scroll-cta {
    justify-self: stretch;
  }

  .tool-image {
    height: 220px;
  }
}
</style>
