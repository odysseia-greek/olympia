<template>
  <div>
    <v-app id="homepage">
      <v-container background="primary" style="max-width: 100%">

        <v-parallax :src="odysseus" height="1200px">
          <div class="d-flex flex-column justify-center align-center text-white">
            <h1 class="text-h2 mb-2" style="margin: 2em; text-align: center">
              Welcome to your Odysseia!
            </h1>
            <h2 class="text-h4 mb-2" style="margin: 1em; text-align: center">
              Start learning Ancient Greek Today
            </h2>
            <h3 class="text-h5 mb-2" style="margin-top: 5em; margin-left: 1em; margin-right: 1em; text-align: center">
              Choose your career in Ancient Greek
            </h3>
              <v-row class="d-flex justify-center text-center">
                <v-col cols="12" md="4">
                  <v-btn color="#ffd700" class="mx-4" width="10em"
                         :href="'/quiz/dialogue?theme=Plato+-+Euthyphro'">
                    {{ buttonTexts.philosopher }}
                  </v-btn>
                </v-col>
                <v-col cols="12" md="4">
                  <v-btn color="#c0c0c0" class="mx-4" width="10em"
                         :href="'/quiz/authorbased?theme=Herodotos+-+Histories&segment=1.1.0'">
                    {{ buttonTexts.sophist }}
                  </v-btn>
                </v-col>
                <v-col cols="12" md="4">
                  <v-btn color="#cd7f32" class="mx-4" width="10em"
                         :href="'quiz/media?theme=Basic&segment=First+Words+1'">
                    {{ buttonTexts.demagogue }}
                  </v-btn>
                </v-col>
              </v-row>
            <v-icon class="mt-4 scroll-icon" size="36" @click="scrollMeTo('exampleSectionRef')">
              mdi-chevron-down
            </v-icon>
            <p class="scroll-indicator" style="margin-top: 1em; text-align: center">
              Scroll down for more info about Odysseia-Greek
            </p>
          </div>
        </v-parallax>

        <v-container fluid style="max-width: 1400px">
          <div ref="exampleSectionRef"></div>
          <v-card color="triadic" style="margin-bottom: 1em; padding: 2em;" elevation="10" rounded="lg">
            <v-card-text>
              <v-row>
                <v-col cols="12" md="6">
                  <div class="text-h4 text--primary my-text-block" style="margin-top:1em;">
                    Welcome to Odysseia-Greek, your interactive gateway to learning Ancient (Attic) Greek!
                  </div>
                  <div :class="$vuetify.theme.textsecondary + ' text-h5 my-text-block'">
                    Dive into the world of Ancient Greece with our app designed to make learning fun and educational. While there's a wealth of resources available for studying Ancient Greek, Odysseia-Greek aims to complement them by offering a unique, interactive experience.
                  </div>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="12" md="6">
                  <v-card class="custom-card" color="text" rounded="lg" style="margin-top:2em;">
                    <v-card-text>
                      <TypingText :texts="introTexts" />
                    </v-card-text>
                  </v-card>
                </v-col>
                <v-col cols="12" md="6">
                  <div :class="$vuetify.theme.textsecondary + ' text-h5 my-text-block'" style="margin-top:2em;">
                    <strong>Why choose Odysseia-Greek?</strong>
                    <br><br>
                    <strong>Engaging Content:</strong> All tools and exercises are designed to keep you motivated and entertained as you learn.
                    <br>
                    <strong>Easy to Use:</strong> Each component is designed to be interactive and easy to use. You can find a howto with examples at each component.
                    <br>
                    <strong>Community Driven:</strong> We value your input! If you spot any mistakes or have suggestions, don’t hesitate to reach out.
                  </div>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="12" md="6">
                  <div :class="$vuetify.theme.textsecondary + ' text-h5 my-text-block'" style="margin-top:2em;">
                    <strong>New to Ancient Greek?</strong>
                    <br>
                    Not sure where to begin? Scroll up and start your journey as a philosopher (advanced), sophist (intermediate) or demagogue (beginner). This will take you to a quiz appropriate for your level. Or scroll down to see each component explained with a brief introduction.
                    <br><br>
                    Embark on your learning adventure with Odysseia-Greek and explore the beauty of Ancient Greek language and culture.
                  </div>
                </v-col>
                <v-col cols="12" md="6">
                  <v-card class="custom-card" color="text" rounded="lg" style="margin-top:2em;">
                    <v-card-text style="padding: 0;" width="100%">
                      <v-card-title>Example Dictionary</v-card-title>
                      <SearchComponent />
                    </v-card-text>
                  </v-card>
                </v-col>
                <v-col cols="12" md="6">
                  <div :class="$vuetify.theme.textsecondary + ' text-h5 my-text-block'" style="margin-top:2em;">
                    <strong>Ancient Greek can be challenging</strong>
                    <br>
                    Learning Ancient Greek, like any new language, can be difficult. The complexity often depends on your knowledge of the alphabet and familiarity with declensions.
                    <br><br>
                    Odysseia-Greek caters to both experienced learners and beginners. For those with more experience, we offer a dialogue mode in quizzes, along with dictionary and grammar components. For beginners or absolute newcomers, we provide media quiz modes and a texts component to facilitate learning.
                  </div>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
          <v-row dense>
            <v-col v-for="card in cards" :key="card.title" :cols="flex">
              <v-card :color="card.color" height="100%">
                <a :href="card.link" style="color: white;">
                  <v-img
                      :src="card.src"
                      height="20em"
                      cover
                      gradient="to bottom, rgba(0,0,0,.1), rgba(0,0,0,.5)"
                  >
                    <v-row class="fill-height white--text align-end" style="margin: 0;">
                      <v-card-title v-text="card.title" class="white--text"></v-card-title>
                    </v-row>
                  </v-img>
                </a>
                <br><br>
                <v-card-subtitle class="pb-0">{{ card.subTitle }}</v-card-subtitle>
                <v-card-text class="text--primary">
                  <div class="text--primary">{{ card.shortText }}</div>
                  <br>
                  <div class="text--primary">{{ card.longText }}</div>
                </v-card-text>
                <v-card-actions>
                  <v-btn variant="text" style="cursor: pointer" :to="card.link">Link</v-btn>
                </v-card-actions>
              </v-card>
            </v-col>
          </v-row>
        </v-container>
      </v-container>
    </v-app>
  </div>
</template>

<script>
import {ref, onMounted, nextTick} from 'vue';
import TypingText from "@/components/TypingText.vue";
import SearchComponent from "@/components/SearchBar.vue";

export default {
  name: "HomePage",
  components: {SearchComponent, TypingText},
  setup() {
    const flex = ref(6);
    const odysseus = ref('');
    const exampleSectionRef = ref();
    const images = import.meta.glob('/src/assets/*.webp');
    const buttonTexts = ref({
      philosopher: "Philosopher",
      sophist: "Sophist",
      demagogue: "Demagogue"
    });

    const loadImage = () => {
      import('@/assets/odysseus.webp').then((module) => {
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
          exampleSectionRef.value.scrollIntoView({ behavior: 'smooth' });
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
        color: "triadic",
        src: '',
        image: 'school_of_athens.webp',
        link: "/quiz",
        subTitle: "Sokrates",
        shortText: "Imagine standing in the Athenian Agora, just selling your wares. A man, an ugly man at that, walks up and strikes up a conversation. Before you know it you need to answer questions about morality and how you know you are correct.",
        longText: "That insisted questioning is why Sokrates is the name given to the Quiz section, it is a place to practise multiple choice word options taken from famous Greek works. Image based quizzes or even dialogue! It's the best place to start if you are new.",
      },
      {
        title: "Texts",
        color: "triadic",
        src: '', // Will be loaded dynamically
        image: 'alexandria.webp',
        link: "/texts",
        subTitle: "Herodotos",
        shortText: "Known as the father of history Herodotos came from the city of Halicarnassus in what is now Western Turkey but then a thriving Greek city.",
        longText: "What can be said with certainty is that his influence over Ancient Greek is vast. He is part of most textbooks that help students take their first steps into Ancient Greek. In the Herodotos or Text section of Odysseia you can translate texts. It will calculate how accurate your translation is compared to an official translation.",
      },
      {
        title: "Grammar",
        color: "triadic",
        src: '', // Will be loaded dynamically
        image: 'grammar.webp',
        link: "/grammar",
        subTitle: "Dionysios",
        shortText: "As a grammarian working in Alexandria, at the time the centre of Greek learning, Dionysios has had a large influence on learning Greek. His grammatical text has been the standard manual for a period up to the Renaissance.",
        longText: "The fact that he wrote about Grammar and his text is known as a grammatical text makes Dionysios the logical name for the Grammar part of Odysseia. So what can you expect here? Dionysios will help you decline words and search for a possible meaning. The backend takes care of the heavy lifting for you. Currently nouns, verbs and participia are supported (but very much incomplete).",
      },
      {
        title: "Dictionary",
        color: "triadic",
        src: '', // Will be loaded dynamically
        image: 'alexander.webp',
        link: "/dictionary",
        subTitle: "Alexandros",
        shortText: "The end of the earth was still not enough in his search and quest for glory. Few people loom as large as Alexander the Great in world history.",
        longText: "His name is therefore linked to the dictionary component of Odysseia which is pretty self-explanatory. You can search for (partial) words here. Each keystroke will send a query to the backend and accents are ignored for easy searching. Support English, Ancient Greek or Dutch (Nederlands) as search parameters.",
      },
    ]);

    const isMobile = () => screen.width <= 900;

    onMounted(() => {
      let showingAdvanced = true;
      loadImage();
      loadCardImages();
      if (isMobile()) {
        flex.value = 12;
        import('@/assets/mobile_odysseus.webp').then((module) => {
          odysseus.value = module.default;
        });
      } else {
        flex.value = 6;
      }

      const interval = setInterval(() => {
          // After blinking, switch all button texts
          showingAdvanced = !showingAdvanced;
          buttonTexts.value.philosopher = showingAdvanced ? "Advanced" : "Philosopher";
          buttonTexts.value.sophist = showingAdvanced ? "Intermediate" : "Sophist";
          buttonTexts.value.demagogue = showingAdvanced ? "Beginner" : "Demagogue";

      }, 4000); // Adjust timing as needed
    });


    return {
      flex,
      buttonTexts,
      odysseus,
      cards,
      introTexts,
      exampleSectionRef,
      scrollMeTo,
    };
  },
};
</script>

<style>
.my-text-block {
  margin-bottom: 1em; /* Adjust this value as needed */
}

.custom-card {
  padding: 1em;
}

#homepage h1.text-lg-h1 {
  color: rgb(248, 248, 248) !important;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
}

</style>



