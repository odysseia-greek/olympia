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
const selectThemeRef = ref(null);
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
      quizContainerRef.value.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }

    if (refName === 'selectTheme' && selectThemeRef.value) {
      selectThemeRef.value.scrollIntoView({ behavior: 'smooth', block: 'start' });
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
  <v-container class="quiz-container media-setup-container text-center">
    <v-card class="paper-card media-setup-card pa-6" elevation="4">
      <div class="setup-toolbar">
        <div class="setup-status" v-if="theme || segment">
          <v-chip v-if="theme" color="primary" variant="tonal">{{ theme }}</v-chip>
          <v-chip v-if="segment" color="secondary" variant="tonal">{{ segment }}</v-chip>
        </div>
        <div>
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
      </div>

      <div v-if="!minimized" ref="selectThemeRef">
        <v-card-title class="text-h5 media-title">Media Quiz</v-card-title>
        <!-- Thematic quote always visible -->
        <p class="media-quote">
          <span>Ἀρχὴ πάσης πράξεως ἐστὶν ἡ τοῦ αἱρεῖσθαι ἀρχή.</span>
          <v-divider class="my-4" />
          <span>The beginning of every action is the choice.</span>
        </p>

      <!-- Only show this part if setupStep === 0 -->
      <template v-if="setupStep === 0">
        <v-card class="media-intro-card ma-5" variant="flat">
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
          <div class="intro-actions">
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
          </div>
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
        <v-row v-if="segments.length > 0 && setupStep === 2 || setupStep === 10" class="segment-grid mt-4">
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
        <v-row class="quiz-options-row mt-4" v-if="setupStep === 10">
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

    <v-card class="quiz-word-container media-word-card" v-if="!finished">
      <h2 class="quiz-word" v-if="!showNextQuestionIndicator">
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
      <v-row class="answer-grid">
        <v-col v-for="item in answers" :key="item.option" cols="6">
          <v-card
              class="media-answer-card"
              flat
              @click="checkAnswer(item)"
              :class="{
              'card-correct': answerStates[item.option]?.isCorrect,
              'card-incorrect': !answerStates[item.option]?.isCorrect && answerStates[item.option]?.selected,
              'is-waiting': !nextAnswerSelectable,
            }"
          >
            <v-img
                :src="loadedImages[item.imageUrl] || ''"
                class="media-answer-image mb-2"
                aspect-ratio="1"
                cover
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

<style scoped>
.media-setup-container,
.inner-quiz-area {
  scroll-margin-top: 80px;
}

.media-setup-card {
  border: 1px solid rgba(28, 97, 209, 0.16);
  color: #20334f;
  background:
      linear-gradient(145deg, rgba(254, 252, 245, 0.98), rgba(253, 246, 227, 0.92)),
      #fdf6e3;
}

.setup-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 40px;
  gap: 16px;
}

.setup-status {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.media-title {
  color: #20334f;
  font-size: clamp(1.5rem, 3vw, 2.2rem);
  font-weight: 800;
}

.media-quote {
  max-width: 680px;
  margin: 16px auto 26px;
  color: #536987;
  line-height: 1.7;
}

.media-quote span:first-child {
  color: #20334f;
  font-size: 1.12rem;
  font-weight: 800;
}

.media-intro-card {
  border: 1px solid rgba(28, 188, 209, 0.18);
  background: rgba(255, 255, 255, 0.58);
  color: #20334f;
  text-align: left;
}

.intro-actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}

.segment-grid .v-btn {
  min-height: 44px;
  text-transform: none;
  white-space: normal;
}

.quiz-options-row {
  justify-content: center;
  text-align: left;
}

.media-word-card {
  min-height: 11rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(28, 97, 209, 0.22);
  background:
      radial-gradient(circle at 18% 24%, rgba(28, 209, 140, 0.2), transparent 30%),
      linear-gradient(135deg, rgba(28, 188, 209, 0.18), rgba(253, 246, 227, 0.96));
  box-shadow: 0 14px 34px rgba(28, 97, 209, 0.14);
}

.media-word-card .quiz-word {
  margin: 0 0 8px;
  color: #10284b;
  letter-spacing: 0.01em;
}

.media-word-card .quiz-instructions {
  margin: 0;
  color: #536987;
}

.answer-grid {
  row-gap: 14px;
}

.media-answer-card {
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(28, 97, 209, 0.14);
  border-radius: 10px;
  background: #fefcf5;
  box-shadow: 0 10px 24px rgba(28, 97, 209, 0.1);
  cursor: pointer;
  transition: transform 0.18s ease, box-shadow 0.18s ease, filter 0.18s ease;
}

.media-answer-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 16px 32px rgba(28, 97, 209, 0.16);
}

.media-answer-card.is-waiting {
  cursor: default;
}

.media-answer-card.is-waiting:not(.card-correct):not(.card-incorrect) {
  filter: saturate(0.78) opacity(0.72);
}

.media-answer-image {
  background: #fdf6e3;
}

.media-answer-card :deep(.v-card-text) {
  color: #20334f;
  font-weight: 800;
}

.card-correct {
  animation: none;
  border: 3px solid #1cd18c;
  box-shadow: 0 0 0 6px rgba(28, 209, 140, 0.16);
}

.card-incorrect {
  animation: none;
  border: 3px solid #d1311c;
  box-shadow: 0 0 0 6px rgba(209, 49, 28, 0.14);
}

@media (max-width: 700px) {
  .media-setup-card {
    padding: 18px !important;
  }

  .setup-toolbar {
    align-items: flex-start;
  }

  .media-intro-card {
    margin-left: 0 !important;
    margin-right: 0 !important;
  }
}
</style>
