<script setup>
import { useRoute } from 'vue-router';
import { useBouleId } from '@/composables/useBoule';
import { useQuery } from '@vue/apollo-composable';
import {
  AuthorBasedAnswer,
  AuthorBasedCreateQuestion, AuthorBasedForms,
  AuthorBasedOptions
} from '@/constants/authorBasedGraphql';
import {computed, getCurrentInstance, nextTick, reactive, ref, watch} from 'vue';
import { updateQuizUrl } from '@/utils/sharedQuiz';
import QuizProgress from "@/components/QuizProgress.vue";
import HistoryTable from "@/components/HistoryTable.vue";
import {apolloClient} from "@/apollo";
import GrammarDetails from "@/components/GrammarDetails.vue";

const { proxy } = getCurrentInstance();
const theme = ref('');
const segment = ref('');
const set = ref(1);
const showHistoryIndicator = ref(false);
const comprehensive = ref(false);
const correct = ref(false);
const showNextQuestionIndicator = ref(false);
const nextAnswerSelectable = ref(true);
const numberOfAnswersNeeded = ref(2);
const attemptMade = ref(false);
const newHistoryItemToPush = ref(null);
const finished = ref(false);
const finishedAfterGrammar = ref(false);
const answers = ref([]);
const numberOfItemsInSet = ref(0);
const answerStates = reactive({});

const searchQuery = ref('');
const expandedPanels = ref([]);
const expandAll = ref(false);
const selectedAuthor = ref('');
const selectedBook = ref('');
const selectedSegment = ref('');
const selectedSection = ref('');
const selectedSectionIndex = ref(-1);
const selectedBookReferences = ref([]);
const selectedReferenceSections = ref([]);
const quizContainerRef = ref(null);
const selectThemeRef = ref(null);
const translation = ref('');
const splitAuthorSentence = ref({});
const clickedWord = ref('');
const forceUpdate = ref(0);
const grammarQuizzes = ref ([])
const grammarQuizMode = ref(false)
const currentGrammarIndex = ref(0)

const quizWord = ref('');
const themes = ref([]);
const segments = ref([]);
const setupStep = ref(0);
const minimized = ref(false);

const totalPlayed = ref(0);
const totalMistakes = ref(0);
const worstThree = ref([]);
const streak = ref(0);
const quizProgress = ref([]);
const dictToForms = ref({});
const route = useRoute();
const boule = useBouleId();

const { result: optionsResult, loading, onResult } = useQuery(AuthorBasedOptions);

onResult(({ data }) => {
  if (data?.authorBasedOptions?.themes) {
    themes.value = data.authorBasedOptions.themes;
    initializeFromRoute();
  }
});

watch([ numberOfAnswersNeeded], ([newNumberOfAnswersNeeded]) => {
  if (newNumberOfAnswersNeeded) {
    updateQuizUrl(
        proxy.$router,
        proxy.$route.query,
        'QuizAuthorBased',
        { doneAfter: newNumberOfAnswersNeeded }
    );
  }
})

const authors = computed(() => {
  const result = [];

  themes.value.forEach(theme => {
    const [authorName, bookTitleRaw] = theme.name.split(' - ');
    const bookTitle = bookTitleRaw?.trim() || 'Unknown Book';

    let author = result.find(a => a.key === authorName);

    if (!author) {
      author = {
        key: authorName,
        books: [],
      };
      result.push(author);
    }

    author.books.push({
      key: bookTitle,
      references: theme.segments.map(seg => ({
        key: seg.name,
        maxSet: seg.maxSet
      }))
    });
  });

  return result;
});

const filteredAuthors = computed(() => {
  const query = (searchQuery.value || '').toLowerCase();
  const results = authors.value
      .map(author => {
        const filteredBooks = author.books.filter(book =>
            book.key.toLowerCase().includes(query)
        );
        if (author.key.toLowerCase().includes(query) || filteredBooks.length > 0) {
          return {
            ...author,
            filteredBooks,
          };
        }
        return null;
      })
      .filter(Boolean);

  if (query !== '') {
    expandedPanels.value = Array.from({ length: results.length }, (_, index) => index);
  }

  return results;
});

const autocompleteAuthorsAndBooks = computed(() => {
  return authors.value.flatMap(author => [
    author.key,
    ...author.books.map(book => book.key)
  ]);
});

const toggleExpandAll = () => {
  expandAll.value = !expandAll.value;
  expandedPanels.value = expandAll.value
      ? Array.from({ length: filteredAuthors.value.length }, (_, i) => i)
      : [];
};

const handleKeyPress = (event) => {
  if (event.key === 'Enter') {
    searchQuery.value = event.target.value;
  }
};

const onBookSelected = (authorKey, bookKey) => {
  selectedAuthor.value = authorKey;
  selectedBook.value = bookKey;

  const author = authors.value.find(a => a.key === authorKey);
  const book = author?.books.find(b => b.key === bookKey);
  selectedBookReferences.value = book?.references || [];

  selectedSegment.value = '';
  selectedReferenceSections.value = [];
  selectedSection.value = '';
  selectedSectionIndex.value = -1;

  // Optional: sync theme and segments if needed
  const themeEntry = themes.value.find(t => t.name === `${authorKey} - ${bookKey}`);
  if (themeEntry) {
    theme.value = themeEntry.name;
    segment.value = ''; // reset
    setupStep.value = 2;
    finished.value = false;
    segments.value = themeEntry.segments;
  }

  updateQuizUrl(
      proxy.$router,
      proxy.$route.query,
      'QuizAuthorBased',
      {
        theme: theme.value,
      }
  );
};

const onReferenceSelected = (bookKey, referenceKey) => {
  selectedBook.value = bookKey;
  selectedSegment.value = referenceKey;
  finished.value = false;

  // You may enhance this later if sections are available per reference
  selectedReferenceSections.value = [];
  selectedSectionIndex.value = -1;
  selectedSection.value = '';

  getAuthorBasedQuiz()
  fetchWordForms()

  updateQuizUrl(
      proxy.$router,
      proxy.$route.query,
      'QuizAuthorBased',
      {
        theme: theme.value,
        segment: referenceKey,
      }
  );
};


// Fetch quiz from backend
const getAuthorBasedQuiz = async () => {
  try {
    const { data } = await apolloClient.query({
      query: AuthorBasedCreateQuestion,
      variables: {
        input: {
          theme: theme.value,
          segment: selectedSegment.value,
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

    const result = data.authorBasedQuiz;
    if (!result) {
      return;
    }
    quizWord.value = result.quiz.quizItem;
    numberOfItemsInSet.value = result.quiz.numberOfItems;
    answers.value = result.quiz.options;
    splitAuthorSentence.value = result.fullSentence.split(/(\s+|[.,;:!?·])/).filter(Boolean);
    translation.value = result.translation;
    if (result.grammarQuiz) {
      grammarQuizzes.value = result.grammarQuiz
    } else {
      grammarQuizzes.value = []
    }

    if (result.progress) {
      updateProgressStats(result.progress);
    }

    scrollMeTo('quiz')

  } catch (err) {
    console.error('Error fetching media quiz:', err);
  }
};

const fetchWordForms = async () => {
  try {
    const { data } = await apolloClient.query({
      query: AuthorBasedForms,
      variables: {
        input: {
          theme: theme.value,
          segment: selectedSegment.value,
          set: String(set.value),
        },
      },
      fetchPolicy: 'no-cache',
    });

    // Rebuild dictToForms from the response
    dictToForms.value = {};
    data.authorBasedWordForms.forms.forEach(entry => {
      dictToForms.value[entry.dictionaryForm] = entry.wordsInText;
    });

  } catch (err) {
    console.error('Error fetching word forms:', err);
  }
};

const checkAnswer = async (selectedAnswer) => {
  if (!nextAnswerSelectable.value) return;

  attemptMade.value = true;

  try {
    const { data } = await apolloClient.query({
      query: AuthorBasedAnswer,
      variables: {
        input: {
          theme: theme.value,
          segment: selectedSegment.value,
          set: String(set.value),
          quizWord: quizWord.value,
          answer: selectedAnswer.quizWord,
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

    const result = data.authorBasedAnswer;
    showNextQuestionIndicator.value = true;
    nextAnswerSelectable.value = false;

    correct.value = result.correct

    // update UI based on result
    answerStates[selectedAnswer.quizWord] = {
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

    if (result.progress) {
      updateProgressStats(result.progress);
    }

    newHistoryItemToPush.value = {
      greek: quizWord.value,
      input: selectedAnswer.quizWord,
      correct: correct.value,
    };
    if (correct.value && grammarQuizzes.value.length > 0) {
      finishedAfterGrammar.value = result.finished;
      setTimeout(() => {
        showNextQuestionIndicator.value = false;
        Object.keys(answerStates).forEach(k => delete answerStates[k]);
        nextAnswerSelectable.value = true
        grammarQuizMode.value = true
      }, 1000);
      return
    }

    if (correct.value) {
      setTimeout(() => {
        showNextQuestionIndicator.value = false;
        nextAnswerSelectable.value = true
        finished.value = result.finished;
        Object.keys(answerStates).forEach(k => delete answerStates[k]);
        if (correct.value) {
          getAuthorBasedQuiz();
        }
      }, 1500);
    } else {
      finished.value = false
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

const checkGrammarAnswer = async (answer) => {
  if (!nextAnswerSelectable.value) {
    return
  }

  nextAnswerSelectable.value = false
  totalPlayed.value++;

  const grammarCorrect = answer.quizWord === grammarQuizzes.value[currentGrammarIndex.value].correctAnswer
  newHistoryItemToPush.value = {
    greek: quizWord.value,
    input: answer.quizWord,
    correct: grammarCorrect,
  };

  if (grammarCorrect) {
    streak.value++;
    answerStates[answer.quizWord] = {
      selected: true,
      isCorrect: true,
    };
    showNextQuestionIndicator.value = true
    nextAnswerSelectable.value = true
    setTimeout(() => {
      if (grammarQuizzes.value.length > currentGrammarIndex.value+1) {
        currentGrammarIndex.value++
        quizWord.value = grammarQuizzes.value[currentGrammarIndex.value].wordInText;
        showNextQuestionIndicator.value = false
        answerStates[answer.quizWord] = {
          selected: false,
          isCorrect: true,
        };
      } else {
        showNextQuestionIndicator.value = false
        grammarQuizMode.value = false
        grammarQuizzes.value = [];
        currentGrammarIndex.value = 0;
        answerStates[answer.quizWord] = {
          selected: false,
          isCorrect: true,
        };

        getAuthorBasedQuiz();
      }

      Object.keys(answerStates).forEach(k => delete answerStates[k]);
    }, 1500);
  } else {
    totalMistakes.value++
    streak.value = 0;
    answerStates[answer.quizWord] = {
      selected: true,
      isCorrect: false,
    };
    setTimeout(() => {
      nextAnswerSelectable.value = true
      answerStates[answer.quizWord] = {
        selected: false,
        isCorrect: false,
      };
    }, 1000)
  }
  if (finishedAfterGrammar.value) {
    finished.value = true;
  }
}

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

const updateProgressStats = (progressArray) => {
  quizProgress.value = progressArray;
};

const goToTextEntry = () => {
  window.open(authorSpecificContent.value.reference, '_blank');
}

const wordOpacities = computed(() => {
  const result = {};
  const allOpacities = [];

  for (const [dictForm, surfaceForms] of Object.entries(dictToForms.value)) {
    const progressEntry = quizProgress.value.find(entry => entry.greek === dictForm);
    const correct = progressEntry?.correctCount || 0;
    let opacity = 3
    if (correct !== 0) {
      opacity = Math.min(100, (correct / numberOfAnswersNeeded.value) * 100);
    }

    surfaceForms.forEach(surface => {
      result[surface] = opacity; // override default
      allOpacities.push(opacity);
    });
  }

  // Set punctuation based on max or min opacity
  const punctuationOpacity = allOpacities.length
      ? Math.max(...allOpacities) // or Math.min(...)
      : 3;

  splitAuthorSentence.value.forEach(part => {
    if (/^[.,;:!?·]$/.test(part)) {
      result[part] = punctuationOpacity;
    }
  });


  return result;
});

const truncateText = (text) => {
  const maxLength = 35;
  if (text.length > maxLength) {
    return text.substring(0, 32) + '...';
  }
  return text;
};

const setClickedWord = (word) => {
  if (wordOpacities.value[word] >= 50) {
    clickedWord.value = word;
    forceUpdate.value++
  }
};

const initializeFromRoute = () => {
  const { theme: qTheme, segment: qSegment, comprehensive: qComprehensive, doneAfter: qDoneAfter } = route.query;

  if (qComprehensive !== undefined) {
    comprehensive.value = qComprehensive === 'true';
  }

  if (qDoneAfter !== undefined) {
    numberOfAnswersNeeded.value = parseInt(qDoneAfter, 10);
  }
  if (qTheme) {
    const themeEntry = themes.value.find(t => t.name.toLowerCase() === qTheme.toLowerCase());
    if (themeEntry) {
      theme.value = themeEntry.name;
      const [authorName, bookTitle] = themeEntry.name.split(' - ');
      onBookSelected(authorName, bookTitle);

      if (qSegment) {
        const segmentObj = themeEntry.segments.find(
            s => s.name.toLowerCase() === qSegment.toLowerCase()
        );
        if (segmentObj) {
          onReferenceSelected(bookTitle, segmentObj.name);
        }
      }
    }
  }
};
</script>

<template>
  <v-container class="quiz-container author-setup-container text-center">
    <v-card class="paper-card author-setup-card pa-6" elevation="4">
      <div class="setup-toolbar">
        <div class="setup-status" v-if="selectedAuthor || selectedBook || selectedSegment">
          <v-chip v-if="selectedAuthor" color="primary" variant="tonal">{{ selectedAuthor }}</v-chip>
          <v-chip v-if="selectedBook" color="secondary" variant="tonal">{{ selectedBook }}</v-chip>
          <v-chip v-if="selectedSegment" color="triadic" variant="tonal">{{ selectedSegment }}</v-chip>
        </div>
        <div>
        <v-btn
            icon="mdi-minus"
            variant="text"
            @click="minimized = !minimized"
        >
          <v-icon>{{ minimized ? 'mdi-plus' : 'mdi-minus' }}</v-icon>
        </v-btn>
        </div>
      </div>

      <div ref="selectThemeRef" v-if="!minimized">
        <v-card-title class="text-h5 author-title">Author Based Quiz</v-card-title>
        <!-- Thematic quote always visible -->
        <p class="author-quote">
          <span>Ἀρχὴ πάσης πράξεως ἐστὶν ἡ τοῦ αἱρεῖσθαι ἀρχή.</span>
          <v-divider class="my-4" />
          <span>The beginning of every action is the choice.</span>
        </p>

        <!-- Show quiz setup options only after choosing/shuffling -->
          <v-container class="author-library-panel">
            <!-- Autocomplete Search -->
            <v-card class="author-search-card mx-auto" variant="flat">
              <v-container>
                <v-row>
                  <v-col>
                    <v-autocomplete
                        v-model="searchQuery"
                        :items="autocompleteAuthorsAndBooks"
                        prepend-icon="mdi-magnify"
                        placeholder="Start Typing"
                        label="Search Authors and Books"
                        variant="underlined"
                        @keyup.enter="handleKeyPress"
                        clearable
                    />
                  </v-col>
                  <v-col cols="auto">
                    <v-btn
                        icon="mdi-minus"
                        variant="text"
                        @click="toggleExpandAll"
                    >
                      <v-icon>{{ expandAll ? 'mdi-minus' : 'mdi-plus' }}</v-icon>
                    </v-btn>
                  </v-col>
                </v-row>
              </v-container>
            </v-card>

            <!-- Author & Book Selection -->
            <v-row v-if="filteredAuthors.length">
              <v-col>
                <v-expansion-panels v-model="expandedPanels" multiple class="author-panels">
                  <v-expansion-panel
                      v-for="author in filteredAuthors"
                      :key="author.key"
                  >
                    <v-expansion-panel-title>
                      {{ author.key }}
                    </v-expansion-panel-title>
                    <v-expansion-panel-text>
                      <v-list class="author-book-list">
                        <v-list-item
                            v-for="book in author.filteredBooks.length ? author.filteredBooks : author.books"
                            :key="book.key"
                            class="author-book-item"
                            :class="{ 'is-selected': selectedAuthor === author.key && selectedBook === book.key }"
                            @click="onBookSelected(author.key, book.key)"
                        >
                          {{ book.key }}
                        </v-list-item>
                      </v-list>
                    </v-expansion-panel-text>
                  </v-expansion-panel>
                </v-expansion-panels>
              </v-col>
            </v-row>

            <!-- References -->
            <v-row v-if="selectedBookReferences.length" class="segment-selection">
              <v-col>
                <v-row>
                  <v-col
                      cols="12"
                  >
                    <h3>Segments</h3>
                  </v-col>
                  <v-btn
                      v-for="reference in selectedBookReferences"
                      :key="reference.key"
                      class="quiz-button segment-button"
                      :class="{ 'pulsate': selectedSegment === reference.key }"
                      :color="selectedSegment === reference.key ? 'primary' : 'triadic'"
                      @click="onReferenceSelected(selectedBook, reference.key)"
                  >
                    {{ reference.key }}
                  </v-btn>
                </v-row>
              </v-col>
            </v-row>
            <v-row class="quiz-options-row mt-4">
              <v-col cols="12" sm="6">
                <v-switch
                    v-model="showHistoryIndicator"
                    color="primary"
                    label="History Table"
                ></v-switch>
              </v-col>
            </v-row>
          </v-container>
          <div class="completion-target-panel">
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

            <QuizProgress
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
          </div>
      </div>
    </v-card>
    <v-container v-if="quizWord !== ''" class="quiz-container text-center">
      <v-container
          v-if="finished"
      >
        <v-card flat class="mb-3 paper-card completion-card">
          <v-card-text>
            <h3> Well Done! You have finished this section! </h3>
            <p>
              Select a new text or search in the bar above
            </p>
          </v-card-text>
          <!-- Arrow up button -->
          <v-btn
              icon="mdi-arrow-up"
              variant="text"
              @click="scrollMeTo('selectTheme')"
          >
          </v-btn>

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

      <v-card class="quiz-word-container author-word-card" v-if="!finished">
        <h2 class="quiz-word" v-if="!showNextQuestionIndicator && !grammarQuizMode">
          {{ quizWord }}
        </h2>
        <h2 class="quiz-word" v-if="!showNextQuestionIndicator && grammarQuizMode">
          {{ grammarQuizzes[currentGrammarIndex].wordInText }}
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
            v-if="!grammarQuizMode"
        >
          Select the correct meaning.
        </p>
        <p
            class="quiz-instructions"
            v-else-if="grammarQuizMode"
        >
          Select the correct grammatical form.
        </p>
      </v-card>
    </v-container>
    <GrammarDetails
        :clickedWord="clickedWord"
        :forceUpdate="forceUpdate"
    />
      <v-row v-if="!finished && !grammarQuizMode" class="answer-grid">
        <v-col
            v-for="item in answers"
            :key="item.quizWord"
            cols="12"
            sm="6"
        >
          <v-btn
              @click="checkAnswer(item);"
              class="author-answer-button ma-1"
              :class="{
            'answer-correct': answerStates[item.quizWord]?.isCorrect,
            'answer-incorrect': !answerStates[item.quizWord]?.isCorrect && answerStates[item.quizWord]?.selected
          }"
              :color="answerStates[item.quizWord]?.selected
            ? answerStates[item.quizWord]?.isCorrect
              ? '#1de9b6'
              : '#e9501d'
            : 'triadic'"
              block
          >
            <span>{{ truncateText(item.quizWord) }}</span>
          </v-btn>
        </v-col>
      </v-row>
      <v-row v-if="!finished && grammarQuizMode" class="answer-grid">
        <v-col
            v-for="item in grammarQuizzes[currentGrammarIndex].options"
            :key="item.quizWord"
            cols="12"
            sm="6"
        >
          <v-btn
              @click="checkGrammarAnswer(item);"
              class="author-answer-button ma-1"
              :class="{
            'answer-correct': answerStates[item.quizWord]?.isCorrect,
            'answer-incorrect': !answerStates[item.quizWord]?.isCorrect && answerStates[item.quizWord]?.selected
          }"
              :color="answerStates[item.quizWord]?.selected
            ? answerStates[item.quizWord]?.isCorrect
              ? '#1de9b6'
              : '#e9501d'
            : 'triadic'"
              block
          >
            <span>{{ truncateText(item.quizWord) }}</span>
          </v-btn>
        </v-col>
      </v-row>
      <v-card class="paper-card author-passage-card ma-3" v-if="quizWord">
        <v-card-title
            v-if="!finished"
            class="text-wrap white-space-normal break-words"
        >
          Reveal words by answering correctly. Click the word for a
          declension attempt after {{ numberOfAnswersNeeded }} correct answers!
        </v-card-title>
        <h2 class="mb-4" v-if="!finished">
        <span
            v-for="(word, index) in splitAuthorSentence"
            :key="index"
            @click="setClickedWord(word)"
            :style="{
            opacity: wordOpacities[word] + '%',
            cursor: (wordOpacities[word] ?? 3) >= 50 ? 'pointer' : 'default'
          }"
        >
          {{ word }}
        </span>
        </h2>
        <div v-if="finished">
          <h2 class="mb-4">
          <span
              v-for="(word, index) in splitAuthorSentence"
              :key="index"
          >
          {{ word }}
        </span>
          </h2>
          <v-divider v-if="finished"></v-divider>
          <h3 v-if="finished" class="mb-3 mt-3">
            {{ translation }}
          </h3>
        </div>
      </v-card>
    <HistoryTable v-if="showHistoryIndicator" :new-entry="newHistoryItemToPush" />
    </v-container>
</template>

<style scoped>
.author-setup-container,
.quiz-container,
.answer-grid {
  scroll-margin-top: 80px;
}

.author-setup-card {
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

.author-title {
  color: #20334f;
  font-size: clamp(1.5rem, 3vw, 2.2rem);
  font-weight: 800;
}

.author-quote {
  max-width: 680px;
  margin: 16px auto 26px;
  color: #536987;
  line-height: 1.7;
}

.author-quote span:first-child {
  color: #20334f;
  font-size: 1.12rem;
  font-weight: 800;
}

.author-library-panel {
  padding: 0;
}

.author-search-card {
  border: 1px solid rgba(28, 188, 209, 0.18);
  background:
      radial-gradient(circle at top left, rgba(28, 209, 140, 0.14), transparent 30%),
      linear-gradient(135deg, rgba(28, 97, 209, 0.08), rgba(255, 255, 255, 0.58));
  color: #20334f;
}

.author-search-card :deep(.v-field),
.author-search-card :deep(.v-label),
.author-search-card :deep(.v-icon) {
  color: #20334f;
}

.author-panels {
  margin-top: 18px;
  text-align: left;
}

.author-panels :deep(.v-expansion-panel) {
  border: 1px solid rgba(28, 97, 209, 0.12);
  background: rgba(253, 246, 227, 0.9);
  color: #20334f;
}

.author-panels :deep(.v-expansion-panel-title) {
  font-weight: 800;
}

.author-book-list {
  background: #fefcf5;
  color: #20334f;
}

.author-book-item {
  border-radius: 10px;
  margin: 4px 0;
  color: #20334f;
  transition: background 0.18s ease, color 0.18s ease, transform 0.18s ease;
}

.author-book-item:hover {
  background: rgba(28, 188, 209, 0.12);
  transform: translateX(2px);
}

.author-book-item.is-selected {
  background: #1c61d1;
  color: #fff;
}

.segment-selection {
  margin-top: 18px;
  text-align: left;
}

.segment-button {
  min-height: 44px;
  margin: 6px;
  text-transform: none;
  white-space: normal;
}

.quiz-options-row {
  justify-content: center;
  text-align: left;
}

.completion-target-panel {
  max-width: 760px;
  margin: 20px auto 0;
}

.completion-card {
  color: #20334f;
}

.author-word-card {
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

.author-word-card .quiz-word {
  margin: 0 0 8px;
  color: #10284b;
  letter-spacing: 0.01em;
}

.author-word-card .quiz-instructions {
  margin: 0;
  color: #536987;
}

.answer-grid {
  row-gap: 14px;
}

.author-answer-button {
  min-height: 56px;
  border-radius: 10px;
  box-shadow: 0 10px 24px rgba(28, 97, 209, 0.1);
  text-transform: none;
  white-space: normal;
  transition: transform 0.18s ease, box-shadow 0.18s ease, filter 0.18s ease;
}

.author-answer-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 16px 32px rgba(28, 97, 209, 0.16);
}

.author-answer-button.answer-correct {
  animation: none;
  border: 3px solid #1cd18c;
  box-shadow: 0 0 0 6px rgba(28, 209, 140, 0.16);
}

.author-answer-button.answer-incorrect {
  animation: none;
  border: 3px solid #d1311c;
  box-shadow: 0 0 0 6px rgba(209, 49, 28, 0.14);
}

.author-passage-card {
  border: 1px solid rgba(28, 97, 209, 0.14);
  background:
      linear-gradient(145deg, rgba(254, 252, 245, 0.98), rgba(253, 246, 227, 0.92)),
      #fefcf5;
  color: #20334f;
  line-height: 1.9;
}

.author-passage-card h2 {
  color: #10284b;
  font-size: clamp(1.4rem, 3vw, 2.1rem);
  line-height: 1.8;
}

.author-passage-card span {
  transition: opacity 0.2s ease, color 0.2s ease, text-shadow 0.2s ease;
}

.author-passage-card span:hover {
  color: #1c61d1;
  text-shadow: 0 0 18px rgba(28, 188, 209, 0.22);
}

@media (max-width: 700px) {
  .author-setup-card {
    padding: 18px !important;
  }

  .setup-toolbar {
    align-items: flex-start;
  }

  .segment-button {
    width: 100%;
  }
}
</style>
