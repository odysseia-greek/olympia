<script setup>
import {ref, watch, onMounted, reactive, nextTick, getCurrentInstance, computed} from 'vue';
import { useQuery } from '@vue/apollo-composable';
import {apolloClient} from "@/apollo";
import { useRoute } from 'vue-router';
import {useBouleId} from '@/composables/useBoule';
import AnalyzeResults from "@/components/AnalyzeResults.vue";
import QuizProgress from "@/components/QuizProgress.vue";
import HistoryTable from "@/components/HistoryTable.vue";
import { updateQuizUrl } from '@/utils/sharedQuiz.js';
import {GrammarAnswer, GrammarCreateQuestion, GrammarOptions} from "@/constants/grammarBasedGraphql";

const { proxy } = getCurrentInstance();
const theme = ref('');
const segments = ref('');
const segment = ref('');
const set = ref(1);
const maxSet = ref(1);
const comprehensive = ref(false);
const showHistoryIndicator = ref(true);
const correct = ref(false);
const showNextQuestionIndicator = ref(false);
const nextAnswerSelectable = ref(true);
const numberOfAnswersNeeded = ref(2);
const attemptMade = ref(false);
const analyzeResults = ref([]);
const quizContainerRef = ref(null);
const selectThemeRef = ref(null);

const finished = ref(false);
const quizWord = ref('');
const answers = ref([]);
const numberOfItemsInSet = ref(0);

const headerInfo = ref({});

const themes = ref([]);
const setupStep = ref(0);
const minimized = ref(false);


const answerStates = reactive({});

const totalPlayed = ref(0);
const totalMistakes = ref(0);
const worstThree = ref([]);
const streak = ref(0);
const quizProgress = ref([]);

const newHistoryItemToPush = ref(null);
const route = useRoute();
const boule = useBouleId();


const { result: optionsResult, loading, onResult } = useQuery(GrammarOptions);

onResult(({ data }) => {
  if (data && data.grammarOptions) {
    themes.value = data.grammarOptions.themes;
    initializeFromRoute();
  }
});

watch([comprehensive], ([newComprehensive]) => {
  if (newComprehensive !== comprehensive ) {
    updateQuizUrl(
        proxy.$router,
        proxy.$route.query,
        'QuizGrammar', // or 'QuizMedia'
        { comprehensive: newComprehensive }
    );
  }
})

watch([ numberOfAnswersNeeded], ([newNumberOfAnswersNeeded]) => {
  if (newNumberOfAnswersNeeded) {
    updateQuizUrl(
        proxy.$router,
        proxy.$route.query,
        'QuizGrammar', // or 'QuizMedia'
        { doneAfter: newNumberOfAnswersNeeded, }
    );
  }
})

watch([set], ([newSet]) => {
  if (newSet) {
    getGrammarBasedQuiz()
    updateQuizUrl(
        proxy.$router,
        proxy.$route.query,
        'QuizGrammar', // or 'QuizMedia'
        { theme: theme.value, set: newSet }
    );
  }
})

// Fetch quiz from backend
const getGrammarBasedQuiz = async () => {
  try {
    const { data } = await apolloClient.query({
      query: GrammarCreateQuestion,
      variables: {
        input: {
          theme: theme.value,
          set: String(set.value),
          segment: segment.value,
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

    const result = data.grammarQuiz;
    quizWord.value = result.quizItem;
    numberOfItemsInSet.value = result.numberOfItems;
    answers.value = result.options;

    headerInfo.value = {
      stem: result.stem,
      difficulty: result.difficulty,
      description: result.description,
      contractionRule: result.contractionRule,
      dictionaryForm: result.dictionaryForm,
      translation: result.translation,
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
      query: GrammarAnswer,
      variables: {
        input: {
          theme: theme.value,
          set: String(set.value),
          segment: segment.value,
          quizWord: quizWord.value,
          answer: selectedAnswer.option,
          comprehensive: comprehensive.value,
          doneAfter: numberOfAnswersNeeded.value,
          dictionaryForm: headerInfo.value.dictionaryForm,
        },
      },
      context: {
        headers: {
          'boule': boule, // full: "boule": sessionId
        },
      },
      fetchPolicy: 'no-cache',
    });

    const result = data.grammarAnswer;
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
          getGrammarBasedQuiz();
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
    set.value = Math.floor(Math.random() * themeData.maxSet) + 1;
    finished.value = false
  }

  updateQuizUrl(
      proxy.$router,
      proxy.$route.query,
      'QuizGrammar', // or 'QuizMedia'
      {
        theme: theme.value,
        set: String(set.value),
      }
  );
}

const selectSegment = (s) => {
  segment.value = s.name;
  maxSet.value = s.maxSet;
  setupStep.value = 10;
  set.value = 1;

  finished.value = false
  getGrammarBasedQuiz()

  updateQuizUrl(
      proxy.$router,
      proxy.$route.query,
      'QuizGrammar',
      {
        theme: theme.value,
        segment: segment.value,
        set: String(set.value),
      }
  );
};

const decrementSelectedSet = () => {
  if (set.value > 1) {
    set.value--;
  }
};
const incrementSelectedSet = () => {
  if (set.value < maxSet.value) {
    set.value++;
  }
};


const randomize = () => {
  const themeIndex = Math.floor(Math.random() * themes.value.length);
  const randomTheme = themes.value[themeIndex];

  theme.value = randomTheme.name;
  finished.value = false;

  const segmentIndex = Math.floor(Math.random() * randomTheme.segments.length);
  const randomSegment = randomTheme.segments[segmentIndex];

  segments.value = randomTheme.segments;
  selectSegment(randomSegment);

  updateQuizUrl(
      proxy.$router,
      proxy.$route.query,
      'QuizGrammar',
      {
        theme: theme.value,
        segment: segment.value,
        set: String(set.value),
      }
  );

  // Show full form and start quiz
  setupStep.value = 10;
};

const truncateText = (text) => {
  const maxLength = 35;
  if (text.length > maxLength) {
    return text.substring(0, 32) + '...';
  }
  return text;
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
  const { theme: qTheme, segment: qSegment, set: qSet, comprehensive: qComprehensive, doneAfter: qDoneAfter } = route.query;
  if (qTheme) {
    const themeData = themes.value.find((t) => t.name.toLowerCase() === qTheme.toLowerCase());
    if (themeData) {
      theme.value = themeData.name;
      quizWord.value = ''
      setupStep.value = 10
      finished.value = false
    }
  }

  if (qSet) {
    set.value = parseInt(qSet, 10);
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
        segments.value = themeObj.segments;
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
  <v-container class="quiz-container grammar-setup-container text-center">
    <v-card class="paper-card grammar-setup-card pa-6" elevation="4">
      <div class="setup-toolbar">
        <div class="setup-status" v-if="theme || segment || setupStep >= 10">
          <v-chip v-if="theme" color="primary" variant="tonal">{{ theme }}</v-chip>
          <v-chip v-if="segment" color="secondary" variant="tonal">{{ segment }}</v-chip>
          <v-chip v-if="setupStep >= 10" color="triadic" variant="tonal">Set {{ set }} / {{ maxSet }}</v-chip>
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
        <v-card-title class="text-h5 grammar-title">Grammar Quiz</v-card-title>
        <!-- Thematic quote always visible -->
        <p class="grammar-quote">
          <span>Ἀρχὴ πάσης πράξεως ἐστὶν ἡ τοῦ αἱρεῖσθαι ἀρχή.</span>
          <v-divider class="my-4" />
          <span>The beginning of every action is the choice.</span>
        </p>

        <!-- Only show this part if setupStep === 0 -->
        <template v-if="setupStep === 0">
          <v-card class="grammar-intro-card ma-5" variant="flat">
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
                    A segment is a subset of quiz items within a theme. Examples are: Aorist - Present - Imperfect.
                  </v-list-item>
                </v-list-item>
                <v-divider></v-divider>
                <br>

                <v-list-item>
                  <v-list-item-title>
                    <strong>Set:</strong>
                  </v-list-item-title>
                  <v-list-item>
                    A set is a logical grouping of within a segment. This might mean uncontracted vs contracted for example.
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
                  style="text-transform: none"
              >
                {{ s.name }}

              </v-btn>
            </v-col>
          </v-row>

          <!-- Toggles -->
          <v-row class="quiz-options-row mt-4" v-if="setupStep === 10">
              <v-col class="text-left">
                          <span class="subheading font-weight-light me-1"
                          >Set
                          </span>
                <span
                    class="text-h4 font-weight-light"
                    v-text="set"
                ></span>
                <span class="subheading font-weight-light me-1">
                            of
                          </span>
                <span
                    class="text-h4 font-weight-light"
                    v-text="maxSet"
                ></span>
              </v-col>
            </v-row>
            <v-slider
                class="my-5"
                v-model="set"
                :max="maxSet"
                :min="1"
                step="1"
                color="primary"
                track-color="accent"
                thumb-color="primary"
                v-if="setupStep === 10"
            >
              <template v-slot:prepend>
                <v-btn
                    icon="mdi-minus"
                    size="small"
                    variant="text"
                    @click="decrementSelectedSet"
                ></v-btn>
              </template>

              <template v-slot:append>
                <v-btn
                    icon="mdi-plus"
                    size="small"
                    variant="text"
                    @click="incrementSelectedSet"
                ></v-btn>
              </template>
            </v-slider>

            <v-col cols="12" sm="6" v-if="setupStep === 3 || setupStep === 10">
              <v-switch
                  v-model="comprehensive"
                  label="Extended Results (Comprehensive)"
                  color="primary"
              />
              <v-switch
                  v-model="showHistoryIndicator"
                  color="primary"
                  label="History Table"
              ></v-switch>
            </v-col>
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
              v-if="setupStep === 10 && quizWord !== ''"
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
            <v-expansion-panels flat class="grammar-details-panel my-6" variant="accordion" v-if="setupStep === 10 && quizWord !== ''">
              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon start icon="mdi-book-open-page-variant"></v-icon>
                  {{ headerInfo.description }}
                </v-expansion-panel-title>

                <v-expansion-panel-text>
                  <v-row dense>
                    <v-col cols="12" sm="4">
                      <v-card flat class="grammar-detail-card">
                        <v-card-subtitle>Dictionary Form</v-card-subtitle>
                        <v-card-text class="text-h6">{{ headerInfo.dictionaryForm }}</v-card-text>
                      </v-card>
                    </v-col>

                    <v-col cols="12" sm="4">
                      <v-card flat class="grammar-detail-card">
                        <v-card-subtitle>Translation</v-card-subtitle>
                        <v-card-text class="text-h6">{{ headerInfo.translation }}</v-card-text>
                      </v-card>
                    </v-col>


                    <v-col cols="12" sm="4" v-if="headerInfo.stem">
                      <v-card flat class="grammar-detail-card">
                        <v-card-subtitle>Stem</v-card-subtitle>
                        <v-card-text class="text-h6">{{ headerInfo.stem }}</v-card-text>
                      </v-card>
                    </v-col>

                    <v-col cols="12" sm="4" v-if="headerInfo.contractionRule">
                      <v-card flat class="grammar-detail-card">
                        <v-card-subtitle>Rule</v-card-subtitle>
                        <v-card-text class="text-body-2">
                          {{ headerInfo.contractionRule }}
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
        </template>
      </div>
    </v-card>
    <div ref="quizContainerRef"></div>
  </v-container>
  <v-container v-if="setupStep === 10 && quizWord !== ''" class="quiz-container text-center">
    <v-container
        v-if="finished"
    >
      <v-card flat class="mb-3 paper-card completion-card">
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

    <v-card class="quiz-word-container grammar-word-card" v-if="!finished">
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
        Choose the correct grammatical answer.
      </p>
    </v-card>
  </v-container>
  <v-container v-if="theme && quizWord && setupStep >= 10 && !finished" class="inner-quiz-area">
    <v-row class="answer-grid">
      <v-col
          v-for="item in answers"
          :key="item.option"
          cols="12"
          sm="6"
      >
        <v-btn
            @click="checkAnswer(item);"
            class="grammar-answer-button ma-1"
            :class="{
            'answer-correct': answerStates[item.option]?.isCorrect,
            'answer-incorrect': !answerStates[item.option]?.isCorrect && answerStates[item.option]?.selected
          }"
            :color="answerStates[item.option]?.selected
            ? answerStates[item.option]?.isCorrect
              ? '#1de9b6'
              : '#e9501d'
            : 'triadic'"
            block
        >
          <span>{{ truncateText(item.option) }}</span>
        </v-btn>
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
.grammar-setup-container,
.inner-quiz-area {
  scroll-margin-top: 80px;
}

.grammar-setup-card {
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

.grammar-title {
  color: #20334f;
  font-size: clamp(1.5rem, 3vw, 2.2rem);
  font-weight: 800;
}

.grammar-quote {
  max-width: 680px;
  margin: 16px auto 26px;
  color: #536987;
  line-height: 1.7;
}

.grammar-quote span:first-child {
  color: #20334f;
  font-size: 1.12rem;
  font-weight: 800;
}

.grammar-intro-card {
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

.grammar-details-panel :deep(.v-expansion-panel) {
  border: 1px solid rgba(28, 97, 209, 0.12);
  background: rgba(253, 246, 227, 0.86);
  color: #20334f;
}

.grammar-detail-card {
  min-height: 100%;
  border: 1px solid rgba(28, 188, 209, 0.14);
  background: #fefcf5;
  color: #20334f;
}

.completion-card {
  color: #20334f;
}

.grammar-word-card {
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

.grammar-word-card .quiz-word {
  margin: 0 0 8px;
  color: #10284b;
  letter-spacing: 0.01em;
}

.grammar-word-card .quiz-instructions {
  margin: 0;
  color: #536987;
}

.answer-grid {
  row-gap: 14px;
}

.grammar-answer-button {
  min-height: 56px;
  border-radius: 10px;
  box-shadow: 0 10px 24px rgba(28, 97, 209, 0.1);
  text-transform: none;
  white-space: normal;
  transition: transform 0.18s ease, box-shadow 0.18s ease, filter 0.18s ease;
}

.grammar-answer-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 16px 32px rgba(28, 97, 209, 0.16);
}

.grammar-answer-button.answer-correct {
  animation: none;
  border: 3px solid #1cd18c;
  box-shadow: 0 0 0 6px rgba(28, 209, 140, 0.16);
}

.grammar-answer-button.answer-incorrect {
  animation: none;
  border: 3px solid #d1311c;
  box-shadow: 0 0 0 6px rgba(209, 49, 28, 0.14);
}

@media (max-width: 700px) {
  .grammar-setup-card {
    padding: 18px !important;
  }

  .setup-toolbar {
    align-items: flex-start;
  }

  .grammar-intro-card {
    margin-left: 0 !important;
    margin-right: 0 !important;
  }
}
</style>
