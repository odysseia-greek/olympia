<script setup>
import {ref, watch, onMounted, reactive, nextTick, getCurrentInstance, computed} from 'vue';
import { useQuery } from '@vue/apollo-composable';
import {apolloClient} from "@/apollo";
import {MediaAnswer, MediaCreateQuestion, MediaOptions} from '@/constants/mediaGraphql';
import { useRoute } from 'vue-router';
import {useBouleId} from '@/composables/useBoule';
import AnalyzeResults from "@/components/AnalyzeResults.vue";
import QuizProgress from "@/components/QuizProgress.vue";
import HistoryTable from "@/components/HistoryTable.vue";
import { updateQuizUrl } from '@/utils/sharedQuiz.js';

const { proxy } = getCurrentInstance();
const theme = ref('');
const segment = ref('');
const set = ref(1);
const maxSet = ref(1);
const comprehensive = ref(false);
const showHistoryIndicator = ref(false);
const correct = ref(false);
const showTextWithImage = ref(false);
const showNextQuestionIndicator = ref(false);
const nextAnswerSelectable = ref(true);
const numberOfAnswersNeeded = ref(2);
const attemptMade = ref(false);
const analyzeResults = ref([]);
const quizContainerRef = ref(null);
const finished = ref(false);

const themes = ref([]);
const segments = ref([]);
const setupStep = ref(0);
const minimized = ref(false);

const quizWord = ref('');
const answers = ref([]);
const numberOfItemsInSet = ref(0);
const answerStates = reactive({});
const loadedImages = reactive({});
const imageMap = import.meta.glob('/src/assets/icons/*.webp');

const totalPlayed = ref(0);
const totalMistakes = ref(0);
const worstThree = ref([]);
const streak = ref(0);
const quizProgress = ref([]);


const route = useRoute();
const boule = useBouleId();
const newHistoryItemToPush = ref(null);

const { result: optionsResult, loading, onResult } = useQuery(MediaOptions);

onResult(({ data }) => {
  if (data && data.mediaOptions) {
    themes.value = data.mediaOptions.themes;
    initializeFromRoute();
  }
});

watch([theme, segment], ([newTheme, newSegment]) => {
  if (newTheme && newSegment) {
    getMediaQuiz();
  }
});

watch([comprehensive], ([newComprehensive]) => {
  if (newComprehensive ) {
    updateQuizUrl(
        proxy.$router,
        proxy.$route.query,
        'QuizMedia',
        { comprehensive: newComprehensive }
    );
  }
})


watch([ numberOfAnswersNeeded], ([newNumberOfAnswersNeeded]) => {
  if (newNumberOfAnswersNeeded) {
    updateQuizUrl(
        proxy.$router,
        proxy.$route.query,
        'QuizMedia',
        { doneAfter: newNumberOfAnswersNeeded }
    );
  }
})

// Fetch quiz from backend
const getMediaQuiz = async () => {
  try {
    const { data } = await apolloClient.query({
      query: MediaCreateQuestion,
      variables: {
        input: {
          theme: theme.value,
          segment: segment.value,
          set: String(set.value),
          doneAfter: numberOfAnswersNeeded.value,
          resetProgress: false,
          archiveProgress: false,
        },
      },
      context: {
        headers: {
          'boule': boule,
        },
      },
      fetchPolicy: 'no-cache',
    });

    const result = data.mediaQuiz;
    quizWord.value = result.quizItem;
    numberOfItemsInSet.value = result.numberOfItems;
    answers.value = result.options;

    for (const item of result.options) {
      const imgPath = `/src/assets/icons/${item.imageUrl}`;
      const importer = imageMap[imgPath];
      if (importer) {
        const module = await importer();
        loadedImages[item.imageUrl] = module.default;
      } else {
        console.warn('Image not found for', item.imageUrl);
      }
    }

    if (result.progress) {
      updateProgressStats(result.progress);
    }

    scrollMeTo('quiz')

  } catch (err) {
    console.error('Error fetching media quiz:', err);
  }
};

const checkAnswer = async (selectedAnswer) => {
  if (!nextAnswerSelectable.value) return;

  attemptMade.value = true;

  try {
    const { data } = await apolloClient.query({
      query: MediaAnswer,
      variables: {
        input: {
          theme: theme.value,
          segment: segment.value,
          set: String(set.value),
          quizWord: quizWord.value,
          answer: selectedAnswer.option,
          comprehensive: comprehensive.value,
          doneAfter: numberOfAnswersNeeded.value,
        },
      },
      context: {
        headers: {
          'boule': boule, // full: "boule": sessionId
        },
      },
      fetchPolicy: 'no-cache',
    });

    const result = data.mediaAnswer;
    showNextQuestionIndicator.value = true;
    nextAnswerSelectable.value = false;

    correct.value = result.correct
    finished.value = result.finished;

    // update UI based on result
    answerStates[selectedAnswer.option] = {
      selected: true,
      isCorrect: result.correct,
    };

    if (correct.value) {
      streak.value++;
    } else {
      totalMistakes.value++
      streak.value = 0;
    }

    totalPlayed.value++;

    if (result.foundInText?.texts?.length > 0) {
      analyzeResults.value = [{
        rootword: result.foundInText.rootword || quizWord.value,
        conjugations: result.foundInText.conjugations || [],
        similarWords: result.similarWords || [],
        results: result.foundInText.texts || [],
      }];
    }

    if (result.progress) {
      updateProgressStats(result.progress);
    }

    newHistoryItemToPush.value = {
      greek: quizWord.value,
      input: selectedAnswer.option,
      correct: correct.value,
    };

    if (correct.value) {
      setTimeout(() => {
        showNextQuestionIndicator.value = false;
        nextAnswerSelectable.value = true
        Object.keys(answerStates).forEach(k => delete answerStates[k]);
        if (correct.value) {
          getMediaQuiz();
        }
      }, 1500);
    } else {
      showNextQuestionIndicator.value = false;
      setTimeout(() => {
        Object.keys(answerStates).forEach(k => delete answerStates[k]);
        nextAnswerSelectable.value = true
      }, 1000);
    }

  } catch (err) {
    console.error('Error checking answer:', err);
  }
};

const updateProgressStats = (progressArray) => {
  quizProgress.value = progressArray;
};


// Handle theme selection
const onThemeChange = (selected) => {
  const themeData = themes.value.find((t) => t.name.toLowerCase() === selected.toLowerCase());
  if (themeData) {
    theme.value = themeData.name;
    quizWord.value = ''
    segments.value = themeData.segments;
    segment.value = ''; // reset segment
    setupStep.value = 2
    finished.value = false
  }

  updateQuizUrl(
      proxy.$router,
      proxy.$route.query,
      'QuizMedia',
      {     theme: themeData.name, }
  );
};

const selectSegment = (s) => {
  segment.value = s.name;
  maxSet.value = s.maxSet;
  setupStep.value = 10;

  if (maxSet.value === 1) {
    set.value = 1;
  }

  finished.value = false

  updateQuizUrl(
      proxy.$router,
      proxy.$route.query,
      'QuizMedia',
      {
        theme: theme.value,
        segment: segment.value,
      }
  );
};

const randomize = () => {
  const themeIndex = Math.floor(Math.random() * themes.value.length);
  const randomTheme = themes.value[themeIndex];

  theme.value = randomTheme.name;
  finished.value = false;

  const segmentIndex = Math.floor(Math.random() * randomTheme.segments.length);
  const randomSegment = randomTheme.segments[segmentIndex];

  segments.value = randomTheme.segments;
  segment.value = randomSegment.name;
  maxSet.value = randomSegment.maxSet;
  if (maxSet.value === 1) {
    set.value = 1;
  }

  updateQuizUrl(
      proxy.$router,
      proxy.$route.query,
      'QuizMedia',
      {
        theme: theme.value,
        segment: segment.value,
      }
  );

  // Show full form and start quiz
  setupStep.value = 10;
};

const scrollMeTo = (refName) => {
  nextTick(() => {
    if (refName === 'quiz' && quizContainerRef.value) {
      quizContainerRef.value.scrollIntoView({ behavior: 'smooth' });
    }
  });
};

const initializeFromRoute = () => {
  const { theme: qTheme, segment: qSegment, comprehensive: qComprehensive, doneAfter: qDoneAfter } = route.query;
  if (qTheme) {
    onThemeChange(qTheme);
  }

  if (qComprehensive !== undefined) {
    comprehensive.value = qComprehensive === 'true';
  }

  if (qDoneAfter !== undefined) {
    numberOfAnswersNeeded.value = parseInt(qDoneAfter, 10);
  }

  if (qTheme) {
    const themeObj = themes.value.find(t => t.name.toLowerCase() === qTheme.toLowerCase());

    if (themeObj) {
      if (qSegment) {
        const segmentObj = themeObj.segments.find(
            s => s.name.toLowerCase() === qSegment.toLowerCase()
        );
        if (segmentObj) {
          selectSegment(segmentObj);
        }
      }
    }
  }
};

</script>

<template>
  <v-container class="quiz-container text-center">
    <v-card class="paper-card pa-6" elevation="4">
      <div class="text-right">
      <v-btn
          v-if="setupStep >= 10"
          icon="mdi-minus"
          variant="text"
          @click="minimized = !minimized"
      >
        <v-icon>{{ minimized ? 'mdi-plus' : 'mdi-minus' }}</v-icon>
      </v-btn>
      <v-btn
          v-if="setupStep >= 10"
          icon="mdi-shuffle-variant"
          variant="text"
          @click="randomize"
      ></v-btn>
      </div>

      <div v-if="!minimized" ref="selectThemeRef">
        <v-card-title class="text-h5">Media Quiz</v-card-title>
        <!-- Thematic quote always visible -->
        <p class="ma-4">
          Ἀρχὴ πάσης πράξεως ἐστὶν ἡ τοῦ αἱρεῖσθαι ἀρχή.
          <v-divider class="my-4" />
          The beginning of every action is the choice.
        </p>

      <!-- Only show this part if setupStep === 0 -->
      <template v-if="setupStep === 0">
        <v-card class="ma-5">
          <v-card-title class="headline">How would you like to begin?</v-card-title>
          <v-card-text>
            <v-list>
              <v-list-item>
                  This section provides information about the different concepts used
              </v-list-item>
              <v-divider></v-divider>
              <br>

              <v-list-item>
                <v-list-item-title>
                  <strong>Theme:</strong>
                </v-list-item-title>
                <v-list-item>
                   A theme represents a major topic that contains subsets (segments).
                </v-list-item>
              </v-list-item>
              <v-divider></v-divider>
              <br>

              <v-list-item>
                <v-list-item-title>
                  <strong>Segment:</strong>
                </v-list-item-title>
                <v-list-item>
                  A segment is a subset of (20) quiz items within a theme. Examples are: Basic Words 1, The Market, War.
                </v-list-item>
              </v-list-item>
              <v-divider></v-divider>
              <br>

              <v-list-item>
                <v-list-item-title>
                  <strong>Extended mode:</strong>
                </v-list-item-title>
                <v-list-item>
                  Searches the quiz word in all it's declined forms in texts from the text component. It was also show similar words coming from the Dictionary.

                  This can be toggled on and off.
                </v-list-item>
              </v-list-item>
              <v-divider></v-divider>
              <br>

              <v-list-item>
                <v-list-item-title>
                  <strong>Images:</strong>
                </v-list-item-title>
                <v-list-item>
                  Currently all images are generated using AI. This can be bothersome to some people but it is the only way it is feasaibile for odysseia-greek at this point.
                </v-list-item>
              </v-list-item>
              <v-divider></v-divider>
              <br>
            </v-list>
          </v-card-text>
            <v-btn
                class="ma-5"
                color="secondary"
                variant="elevated"
                prepend-icon="mdi-format-list-bulleted"
                @click="setupStep = 1"
            >
              Choose Options
            </v-btn>
          <v-btn
              class="ma-5"
              color="primary"
              variant="elevated"
              prepend-icon="mdi-shuffle-variant"
              @click="randomize"
          >
            Random Quiz
          </v-btn>
        </v-card>

      </template>

      <!-- Show quiz setup options only after choosing/shuffling -->
      <template v-else>
        <!-- Theme -->
        <v-combobox
            v-if="setupStep === 1 || setupStep === 10"
            class="mt-5"
            :class="{ 'pulsate': theme === '' }"
            v-model="theme"
            :items="themes.map(t => t.name)"
            item-title="name"
            item-value="name"
            label="Select a Theme"
            color="primary"
            @update:modelValue="onThemeChange"
            style="margin-top: 2em"
        />

        <!-- Segment Buttons -->
        <v-row v-if="segments.length > 0 && setupStep === 2 || setupStep === 10" class="mt-4">
          <v-col cols="12">
            <h4>Choose a Segment</h4>
          </v-col>
          <v-col
              v-for="s in segments"
              :key="s.name"
              cols="6"
              sm="4"
              md="3"
          >
            <v-btn
                block
                color="triadic"
                :variant="segment === s.name ? 'flat' : 'outlined'"
                @click="selectSegment(s)"
            >
              {{ s.name }}

            </v-btn>
          </v-col>
        </v-row>
        <!-- Toggles -->
        <v-row class="mt-4" v-if="setupStep === 10">
          <v-col cols="12" sm="6">
            <v-switch
                v-model="comprehensive"
                label="Extended Results (Comprehensive)"
                color="primary"
            />
            <v-switch
                v-model="showTextWithImage"
                label="Show Text With Images"
                color="primary"
            />
            <v-switch
                v-model="showHistoryIndicator"
                color="primary"
                label="History Table"
            ></v-switch>
          </v-col>
        </v-row>
        <div v-if="setupStep === 3 || setupStep === 10">
        <p>Correct answers needed before marked complete</p>
        <v-slider
            :label="String(numberOfAnswersNeeded)"
            class="my-5"
            v-model="numberOfAnswersNeeded"
            :min="1"
            :max="4"
            step="1"
            color="primary"
            track-color="accent"
            thumb-color="primary"
            show-ticks
        ></v-slider>
        </div>

        <QuizProgress
            v-if="quizWord !== ''"
            :progress="quizProgress"
            :neededCorrectCount="numberOfItemsInSet * numberOfAnswersNeeded"
            :streak="streak"
            :totalPlayed="totalPlayed"
            :totalMistakes="totalMistakes"
            :itemsInThisSet="numberOfItemsInSet"
            @updateWorst="(val) => worstThree = val"
            @initProgress="val => {
              streak = val.streak;
              totalPlayed = val.totalPlayed;
              totalMistakes = val.totalMistakes;
            }"
        />
      </template>
      </div>
    </v-card>
    <div ref="quizContainerRef"></div>
  </v-container>
  <v-container v-if="setupStep === 10 && quizWord !== ''" class="quiz-container text-center">
    <v-container
        v-if="finished"
    >
      <v-card flat class="mb-3 paper-card">
        <v-card-text>
          <h3> Well Done! You have finished this section! </h3>
          <p>
            Select a new quiz by going up or click the randomize button
          </p>
        </v-card-text>
        <!-- Arrow up button -->
        <v-btn
            icon="mdi-arrow-up"
            variant="text"
            @click="scrollMeTo('selectTheme')"
        >
        </v-btn>
        <!-- Randomize button -->
        <v-btn
            icon="mdi-shuffle-variant"
            variant="text"
            @click="randomize"
        ></v-btn>
        <div class="mt-3">
          <template v-if="worstThree.length > 0">
            <v-alert type="info" variant="tonal" class="mt-4">
              <strong>Mistakes are part of the journey — endure like Odysseus.</strong><br />
              <em>τέτλαθι δή, κραδίη· καὶ κύντερον ἄλλο ποτ’ ἔτλης</em><br />
              <small>“Endure, my heart, you have endured worse than this before.” — Homer, <i>Odyssey</i> 20.18</small>
            </v-alert>
            <strong>Most difficult words:</strong>
            <ul class="no-bullets">
              <li v-for="word in worstThree" :key="word.greek">
                <strong>{{ word.greek }} — {{ word.translation }}</strong> {{ word.incorrectCount }} mistake(s)
              </li>
            </ul>
          </template>

          <template v-else>
            <v-alert type="success" variant="tonal" class="mt-2">
              <strong>No mistakes so far — you're certainly acquiring knowledge! 🎉</strong><br />
              <em>τὸ γὰρ γνῶναι ἐπιστήμην που λαβεῖν ἐστιν</em><br />
              <small>“For learning to know is acquiring knowledge.” — Plato, <i>Theaetetus</i></small>
            </v-alert>
          </template>
        </div>


      </v-card>
    </v-container>

    <v-card class="quiz-word-container" height="10em" v-if="!finished">
      <h2 class="quiz-word" style="margin-top: 1em" v-if="!showNextQuestionIndicator">
        {{ quizWord }}
      </h2>
      <div
          v-if="showNextQuestionIndicator"
          class="text-center mb-4"
      >
        <v-progress-circular
            style="margin-bottom: 2em; margin-top: 2em"
            indeterminate
            color="primary"
            width="8"
            size="72"
        ></v-progress-circular>
      </div>
      <p
          class="quiz-instructions"
      >
        Match the word to the image.
      </p>
    </v-card>
  </v-container>
  <v-container v-if="theme && segment && quizWord && setupStep >= 10 && !finished" class="inner-quiz-area">
      <v-row>
        <v-col v-for="item in answers" :key="item.option" cols="6">
          <v-card
              flat
              @click="checkAnswer(item)"
              :class="{
              'card-correct': answerStates[item.option]?.isCorrect,
              'card-incorrect': !answerStates[item.option]?.isCorrect && answerStates[item.option]?.selected,
            }"
          >
            <v-img
                :src="loadedImages[item.imageUrl] || ''"
                class="mb-2"
                aspect-ratio="1"
            />
            <v-card-text v-if="showTextWithImage" class="text-center">
              {{ item.option }}
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    <AnalyzeResults
        v-if="comprehensive && correct"
        :analyzeResults="analyzeResults"
    />
    <HistoryTable v-if="showHistoryIndicator" :new-entry="newHistoryItemToPush" />
    </v-container>

</template>
