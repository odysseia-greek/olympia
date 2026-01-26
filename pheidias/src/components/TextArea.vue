<template>
  <div id="text">
    <v-app :style="{ background: $vuetify.theme.themes[theme].background }">
      <v-main>
        <v-expansion-panels color="secondary" v-if="showInfoBar">
          <v-expansion-panel>
            <v-expansion-panel-title>
              <v-icon left>mdi-information</v-icon>
              Don't know what to do?
              <v-spacer></v-spacer>
              <v-btn
                  icon="mdi-close"
                  variant="text"
                  @click="showInfoBar = !showInfoBar"
                  class="mx-4"
              >
              </v-btn>
            </v-expansion-panel-title>
            <v-expansion-panel-text>
              <v-card-title class="headline">Text Translation</v-card-title>
              <v-card-text>
                <ol>
                  <li><strong>Select Author and Book</strong> once selected you will get the "references" these point to passages in the text and you can see an example here:<a href="https://scaife.perseus.org/reader/urn:cts:greekLit:tlg0016.tlg001.perseus-grc2:1.1.0-1.1.4/"> Scaife</a>.</li>
                  <li>You can use either the <strong>Search Input</strong>, which is autocompleted and can search for both authors and books. Or select an author and a book from the list below.</li>
                  <li><strong>Select Reference and optional section</strong> and start translating! Once you are done click the <strong>CHECK</strong> button in the bottom to see how you did. Your entry is checked against official translations and a levenshtein score is calculated to see how correct you are. This however does not mean your translation is wrong - it might just be a different word order or you used a different word.</li>
                  <li>Each word is <strong>clickable!</strong> This will open a dialog where a declension is attempted as well as links to the grammar and dictionary component to really get all the info about a word that can be found.</li>
                  <li>Once you have translated a text a <strong>Results</strong> section will appear where you can see possible typos and expand each section to see the official translation.</li>
                  <li><v-icon>mdi-arrow-top-right-thick</v-icon> will set the official translation in your translation box.</li>
                  <li><v-icon>mdi-fullscreen</v-icon> will open a full screen, this is especially handy on a mobile phone.</li>
                </ol>
                <h2 class="md-4">Extra Information</h2>
                <ul>
                  <li>
                    <strong>Levenshtein</strong> refers to the similarity between two sentences or words (called strings). It measures the number of steps needed to change one word or sentence into another. Each step can be an insertion, deletion, or substitution of a single character. For example, the words "kitten" and "sitting" have a Levenshtein distance of 3 because you can transform "kitten" into "sitting" by substituting 'k' with 's', substituting 'e' with 'i', and inserting a 'g' at the end. The fewer the steps, the more similar the strings are.
                  </li>
                </ul>
              </v-card-text>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
        <v-container>
          <v-card class="mx-auto" color="primary" dark>
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
                  ></v-autocomplete>
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
          <v-row v-if="filteredAuthors.length">
            <v-col>
              <v-expansion-panels v-model="expandedPanels" multiple color="secondaryPapyrus">
                <v-expansion-panel v-for="author in filteredAuthors" :key="author.key">
                  <v-expansion-panel-title>
                    {{ author.key }}
                  </v-expansion-panel-title>
                  <v-expansion-panel-text style="background-color: #fefcf5">
                    <v-list style="background-color: #fefcf5">
                      <v-list-item
                          v-for="book in author.filteredBooks.length ? author.filteredBooks : author.books"
                          :key="book.key"
                          :style="selectedAuthor === author.key && selectedBook === book.key ? 'color: white; background-color: #1c61d1' : 'color: black; background-color: #fefcf5'"
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
          <v-row v-if="selectedBookReferences.length">
            <v-col>
              <v-row>
                <v-col cols="12">
                  <h3>References</h3>
                </v-col>
                <v-btn
                    v-for="reference in selectedBookReferences"
                    :key="reference.key"
                    class="quiz-button"
                    :class="{ 'pulsate': selectedReference === reference.key }"
                    :color="selectedReference === reference.key ? 'primary' : 'triadic'"
                    @click="onReferenceSelected(selectedBook, reference.key)"
                >
                  {{ reference.key }}
                </v-btn>
              </v-row>
              <v-row v-if="selectedReferenceSections.length && selectedReferenceSections[0].key !==''">
                <v-col cols="12">
                  <h3 ref="loadingResultsRef">Sections</h3>
                </v-col>
                <v-btn
                    v-for="section in selectedReferenceSections"
                    :key="section.key"
                    class="quiz-button"
                    :class="{ 'pulsate': selectedSection === section.key }"
                    :color="selectedSection === section.key ? 'primary' : 'triadic'"
                    @click="onSectionSelected(section.key)"
                >
                  {{ section.key }}
                </v-btn>
                <v-btn
                    class="quiz-button"
                    :color="selectedSectionIndex === -1 ? 'primary' : 'triadic'"
                    @click="onSectionSelected('')"
                >
                  All Sections
                </v-btn>
              </v-row>
            </v-col>
          </v-row>
          <v-row v-if="resultData">
            <v-col>
              <v-card class="mx-auto paper-card" width="100%">
                <div v-if="showLoading" class="text-center mb-4" >
                  <v-progress-circular
                      :model-value="loadingPercentage"
                      :rotate="360"
                      color="primary"
                      width="8"
                      size="75"
                  >
                    {{ loadingPercentage }}
                  </v-progress-circular>
                  <v-card-title>
                    Checking Your Translations...
                  </v-card-title>
                  <v-card-subtitle><v-icon>mdi-information</v-icon>Did you know? You can click each Greek word!</v-card-subtitle>
                </div>
                <GrammarDetails :clickedWord="clickedWord" :forceUpdate="forceUpdate" />
                <v-row v-for="rhema in resultData.create.rhemai" :key="rhema.section" class="rhema-section" v-bind:align="mobileView ? 'center' : undefined">
                  <v-col :cols="12" :md="6">
                    <p><strong>Section {{ rhema.section }}</strong></p>
                    <p>
                    <span v-for="(word, index) in rhema.greek.split(' ')" :key="index">
                      <span class="clickable-word" @click="setClickedWord(word)">
                        {{ word }}
                      </span>
                      <span v-if="index < rhema.greek.split(' ').length - 1">&nbsp;</span>
                    </span>
                    </p>
                  </v-col>
                  <v-col :cols="12" :md="6">
                    <v-textarea
                        append-icon="mdi-fountain-pen-tip"
                        v-model="translations[rhema.section]"
                        v-if="!showLoading"
                        label="Enter your Translation here"
                        auto-grow
                        clearable
                        variant="outlined"
                        rows="1"
                    ></v-textarea>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col cols="12" class="text-right">
                    <v-btn color="primary" @click="checkTranslations">Check</v-btn>
                  </v-col>
                </v-row>
                <v-row v-if="translationResults && !showLoading">
                  <v-col>
                    <v-card class="mx-auto paper-card" width="100%">
                      <h2 ref="resultsContainerRef">Results</h2>
                      <v-card-text>
                        <p class="ma-3"><strong>Average Levenshtein Percentage:</strong> {{ translationResults.averageLevenshteinPercentage }}%</p>
                        <p class="ma-3"><strong>Sections:</strong></p>
                            <div v-if="groupedTranslations.length">
                              <v-expansion-panels>
                                <v-expansion-panel v-for="(section, index) in groupedTranslations" :key="index">
                                  <v-expansion-panel-title class="section-content">Section {{ section.section }}</v-expansion-panel-title>
                                  <v-expansion-panel-text class="section-content">
                                    <v-list dense class="section-content">
                                      <v-list-item v-for="(translation, index) in section.translations" :key="index">
                                          <v-list-item-subtitle>Levenshtein Percentage:
                                          </v-list-item-subtitle>
                                          <v-list-item class="list-item">{{ translation.levenshteinPercentage }}%</v-list-item>
                                          <v-list-item-subtitle>Official Translation:</v-list-item-subtitle>
                                          <v-list-item class="list-item">{{ translation.quizSentence }}</v-list-item>
                                          <v-list-item-subtitle>Provided Translation:</v-list-item-subtitle>
                                          <v-list-item class="list-item">{{ translation.answerSentence }}</v-list-item>
                                        <v-list-item-subtitle>Further Options:</v-list-item-subtitle>
                                        <v-btn @click="sectionFullscreen = true" icon="mdi-fullscreen" variant="text">
                                        </v-btn>
                                        <v-btn @click="setSectionText(translation.quizSentence, section.section)" icon="mdi-arrow-top-right-thick" variant="text">
                                        </v-btn>
                                          <v-dialog v-model="sectionFullscreen" fullscreen hide-overlay transition="dialog-bottom-transition">
                                            <v-card class="mx-auto paper-card">
                                              <v-toolbar dark color="primary">
                                                <v-btn dark @click="sectionFullscreen = false" variant="text">
                                                  <v-icon>mdi-close</v-icon>
                                                </v-btn>
                                                <v-toolbar-title>Section {{ section.section }}</v-toolbar-title>
                                                <v-spacer></v-spacer>
                                              </v-toolbar>
                                              <v-card-text >
                                                <h2>Official Translation:</h2>
                                                <v-text-field readonly>{{ translation.quizSentence }}</v-text-field>
                                                <h2>Provided Translation:</h2>
                                                <v-text-field readonly>{{ translation.answerSentence }}</v-text-field>
                                              </v-card-text>
                                            </v-card>
                                          </v-dialog>
                                      </v-list-item>
                                    </v-list>
                                  </v-expansion-panel-text>
                                </v-expansion-panel>
                              </v-expansion-panels>
                            </div>
                            <div v-if="translationResults.possibleTypos.length">
                              <p style="margin-top:2em"><strong>Possible Typos:</strong></p>
                              <v-list dense class="section-content">
                                <v-list-item v-for="typo in translationResults.possibleTypos" :key="typo.provided" class="list-item">
                                    <v-list-item>Provided: {{ typo.provided }}</v-list-item>
                                    <v-list-item>Correction: {{ typo.source }}</v-list-item>
                                </v-list-item>
                              </v-list>
                            </div>
                          </v-card-text>
                        </v-card>
                      <v-card-actions>
                        <v-btn color="primary" @click="clearTranslations">Clear</v-btn>
                      </v-card-actions>
                  </v-col>
                </v-row>
              </v-card>
            </v-col>
          </v-row>
        </v-container>
      </v-main>
    </v-app>
  </div>
</template>

<script>
import {ref, computed, watch, watchEffect, onMounted, getCurrentInstance, nextTick} from 'vue';
import { useQuery } from '@vue/apollo-composable';
import { HerodotosOptions, HerodotosCreate, HerodotosCheck } from '@/constants/graphql';
import GrammarDetails from "@/components/GrammarDetails.vue";

export default {
  components: {
    GrammarDetails,
  },
  setup() {
    const { proxy } = getCurrentInstance();
    const theme = ref('light');
    const expandedPanels = ref([]);
    const expandAll = ref(false);
    const authors = ref([]);
    const searchQuery = ref('');
    const selectedAuthor = ref('');
    const selectedBook = ref('');
    const selectedReference = ref('');
    const selectedSection = ref('');
    const selectedBookReferences = ref([]);
    const selectedReferenceSections = ref([]);
    const selectedSectionIndex = ref(-1);
    const { result, error } = useQuery(HerodotosOptions);
    const mobileView = ref(window.innerWidth <= 600);
    const resultsContainerRef = ref();
    const loadingResultsRef = ref();

    const resultData = ref(null);
    const queryLoading = ref(false);
    const queryError = ref(null);

    const clickedWord = ref('');
    const grammarResults = ref([]);
    const grammarError = ref(null);
    const translations = ref({});
    const translationResults = ref(null);
    const translationError = ref(null);
    const autocompleteAuthorsAndBooks = ref([]);

    const sectionFullscreen = ref(false);
    const showInfoBar = ref(true);
    const showLoading = ref(false);
    const loadingPercentage = ref(0);
    const forceUpdate = ref(0);

    watchEffect(() => {
      if (result.value) {
        authors.value = result.value.textOptions.authors;

        if (selectedAuthor.value) {
          const index = authors.value.findIndex((authorList) =>
              authorList.key === selectedAuthor.value
          );

          if (index !== -1) {
            expandedPanels.value = index;
          }
        }

        const authorsAndBooks = authors.value.flatMap((author) => {
          const books = author.books.map((book) => book.key);
          return [author.key, ...books];
        });
        autocompleteAuthorsAndBooks.value = authorsAndBooks;
        if (selectedBook.value) {
          const author = authors.value.find((a) => a.key === selectedAuthor.value);
          const book = author.books.find((b) => b.key === selectedBook.value);
          selectedBookReferences.value = sortedReferences(book.references);
        }

        if (selectedReference.value) {
          const book = selectedBookReferences.value.find((b) => b.key === selectedReference.value);
          selectedReferenceSections.value = sortedSections(book.sections);
        }
      }
    });

    const updateData = async () => {
      if (selectedReference.value || selectedSection.value) {
        queryLoading.value = true;
        queryError.value = null;
        try {
          const { onResult } = useQuery(HerodotosCreate, {
            input: {
              author: selectedAuthor.value,
              book: selectedBook.value,
              reference: selectedReference.value,
              section: selectedSection.value || null,
            },
          });

          onResult((response) => {
            if (response.data && response.data.create) {
              resultData.value = response.data;
              translations.value = {};
              response.data.create.rhemai.forEach((rhema) => {
                translations.value[rhema.section] = '';
              });
            } else {
              resultData.value = null;
            }
            queryLoading.value = false;
          });
        } catch (error) {
          queryError.value = error;
          queryLoading.value = false;
        }
      }
    };

    watch(
        [selectedAuthor, selectedBook, selectedReference, selectedSection],
        () => {
          updateData();
        }
    );

    const setSectionText = (text, section) => {
      translations.value[section] = text
    }

    const toggleExpandAll = () => {
      expandAll.value = !expandAll.value;
      if (expandAll.value) {
        expandedPanels.value = Array.from({ length: filteredAuthors.value.length }, (_, index) => index);
      } else {
        expandedPanels.value = [];
      }
    };

    const handleKeyPress = (event) => {
      if (event.key === 'Enter') {
        searchQuery.value = event.target.value;
      }
    };

    const filteredAuthors = computed(() => {
      if (searchQuery.value === null) {
        searchQuery.value = ''
      }
      const query = searchQuery.value.toLowerCase();
      const results = authors.value
          .map((author) => {
            const filteredBooks = author.books.filter((book) =>
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
          .filter((author) => author !== null);

      if (query !== '') {
        expandedPanels.value = Array.from({ length: results.length }, (_, index) => index);
      }

      return results;
    });

    const clearTranslations = () => {
      translationResults.value = null;
      translationError.value = null;
    };

    const checkTranslations = async () => {
      scrollToResults('loadingResults');
      showLoading.value = true;
      loadingPercentage.value = 0;
      translationResults.value = null;
      translationError.value = null;

      const inputTranslations = Object.keys(translations.value).map((section) => ({
        section,
        translation: translations.value[section],
      }));

      const variables = {
        input: {
          author: selectedAuthor.value,
          book: selectedBook.value,
          reference: selectedReference.value,
          translations: inputTranslations,
        },
      };

      try {
        const { onResult, onError } = useQuery(HerodotosCheck, variables);
        onResult((response) => {
          if (response.data) {
            translationResults.value = response.data.check;
          } else {
            translationResults.value = null;
            translationError.value = { message: 'No results found' };
          }
        });
        onError((error) => {
          showLoading.value = false;
          translationResults.value = null;
          translationError.value = error;
        });

        const incrementLoading = () => {
          if (loadingPercentage.value < 100) {
            const increment = Math.floor(Math.random() * 10) + 1;
            loadingPercentage.value = Math.min(loadingPercentage.value + increment, 100);

            setTimeout(incrementLoading, Math.floor(Math.random() * 200) + 100);
          } else {
            showLoading.value = false;
            scrollToResults('results');
            loadingPercentage.value = 0;
          }
        };

        incrementLoading();
      } catch (error) {
        showLoading.value = false;
        translationResults.value = null;
        translationError.value = error;
      }
    };


    const groupedTranslations = computed(() => {
      if (!translationResults.value) return [];

      const grouped = translationResults.value.sections.reduce((acc, section) => {
        if (!acc[section.section]) {
          acc[section.section] = {
            section: section.section,
            translations: [],
          };
        }
        acc[section.section].translations.push({
          levenshteinPercentage: section.levenshteinPercentage,
          answerSentence: section.answerSentence,
          quizSentence: section.quizSentence,
        });
        return acc;
      }, {});

      return Object.values(grouped);
    });

    const onBookSelected = (authorKey, bookKey) => {
      translationResults.value = null;
      translationError.value = null;
      selectedAuthor.value = authorKey;
      selectedBook.value = bookKey;
      const author = authors.value.find((a) => a.key === authorKey);
      const book = author.books.find((b) => b.key === bookKey);
      selectedBookReferences.value = sortedReferences(book.references);
      selectedReference.value = '';
      selectedSection.value = '';
      selectedReferenceSections.value = [];
      resultData.value = null;
      translations.value = {};

      updateURL({
        author: authorKey,
        book: bookKey,
      });
    };

    const onReferenceSelected = (bookKey, referenceKey) => {
      translationResults.value = null;
      translationError.value = null;
      selectedBook.value = bookKey;
      selectedReference.value = referenceKey;
      const book = selectedBookReferences.value.find((b) => b.key === referenceKey);
      selectedReferenceSections.value = sortedSections(book.sections);
      selectedSectionIndex.value = -1;
      selectedSection.value = '';
      resultData.value = null;
      translations.value = {};

      updateURL({
        author: selectedAuthor.value,
        book: bookKey,
        reference: referenceKey,
      });
    };

    const onSectionChanged = (index) => {
      if (index >= 0) {
        selectedSection.value = selectedReferenceSections.value[index].key;
      } else {
        selectedSection.value = '';
      }
      resultData.value = null;
      translations.value = {};
    };

    const onSectionSelected = (sectionKey) => {
      translationResults.value = null;
      translationError.value = null;
      selectedSection.value = sectionKey;
      selectedSectionIndex.value = selectedReferenceSections.value.findIndex((s) => s.key === sectionKey);
      resultData.value = null;
      translations.value = {};
    };

    const sortedSections = (sections) => {
      return sections.slice().sort((a, b) => {
        const aNum = parseFloat(a.key);
        const bNum = parseFloat(b.key);
        if (!isNaN(aNum) && !isNaN(bNum)) {
          return aNum - bNum;
        } else {
          return a.key.localeCompare(b.key);
        }
      });
    };

    const sortedReferences = (references) => {
      return references.slice().sort((a, b) => {
        const parseReference = (ref) => {
          const parts = ref.split('.');
          return parts.map((part) => parseFloat(part));
        };

        const [aMain, aSub] = parseReference(a.key);
        const [bMain, bSub] = parseReference(b.key);

        if (aMain !== bMain) {
          return aMain - bMain;
        } else {
          if (aSub !== undefined && bSub !== undefined) {
            return aSub - bSub;
          } else if (aSub === undefined) {
            return -1;
          } else if (bSub === undefined) {
            return 1;
          }
          return 0;
        }
      });
    };

    const updateURL = (query) => {
      const currentQuery = proxy.$route.query;
      const newQuery = { ...currentQuery, ...query };

      if (!query.reference) {
        delete newQuery.reference;
      }

      const queryChanged = Object.keys(newQuery).length !== Object.keys(currentQuery).length ||
          Object.keys(newQuery).some(key => currentQuery[key] !== newQuery[key]);


      if (queryChanged) {
        proxy.$router.replace({ name: 'Herodotos', query: newQuery });
      }
    };

    const initializeFromRoute = () => {
      const { author, book, reference } = proxy.$route.query;
      if (author) {
        selectedAuthor.value = author;
      }
      if (book) {
        selectedBook.value = book;
      }
      if (reference) selectedReference.value = reference;
    };

    const scrollToResults = (refName) => {
      nextTick(() => {
        if (refName === 'results' && resultsContainerRef.value) {
          resultsContainerRef.value.scrollIntoView({ behavior: 'smooth' });
        }

        if (refName === 'loadingResults' && loadingResultsRef.value) {
          loadingResultsRef.value.scrollIntoView({ behavior: 'smooth' });
        }
      });
    };

    const setClickedWord = (word) => {
      clickedWord.value = word;
      forceUpdate.value++
    };

    onMounted(() => {
      initializeFromRoute(); // Initialize state from URL on mount
    });

    return {
      theme,
      searchQuery,
      authors,
      selectedAuthor,
      selectedBook,
      selectedReference,
      selectedSection,
      selectedBookReferences,
      selectedReferenceSections,
      selectedSectionIndex,
      filteredAuthors,
      autocompleteAuthorsAndBooks,
      error,
      resultData,
      queryLoading,
      queryError,
      translations,
      clickedWord,
      grammarResults,
      grammarError,
      translationResults,
      translationError,
      groupedTranslations,
      mobileView,
      expandedPanels,
      expandAll,
      sectionFullscreen,
      showInfoBar,
      showLoading,
      loadingPercentage,
      resultsContainerRef,
      loadingResultsRef,
      forceUpdate,
      scrollToResults,
      handleKeyPress,
      checkTranslations,
      clearTranslations,
      toggleExpandAll,
      onBookSelected,
      onReferenceSelected,
      onSectionChanged,
      onSectionSelected,
      sortedSections,
      sortedReferences,
      setSectionText,
      setClickedWord,
    };
  },
};
</script>


<style scoped>
.clickable-word {
  cursor: pointer;
}

.quiz-button {
  margin: 0.5em;
}

.mx-auto {
  background-color: #1e3c72; /* Dark blueish background */
}

.v-card .v-card-text {
  max-width: 100%;
  overflow: hidden;
}

.paper-card {
  background: #fdf6e3; /* A light, papyrus-like color */
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
  font-family: 'Roboto', serif;
}

.section-content {
  background: #fefcf5; /* A slightly different shade for sections */
  padding: 10px;
  border-radius: 4px;
}

.rhema-section {
  margin-bottom: 20px;
}

.clickable-word {
  cursor: pointer;
  color: #020202;
  font-family: 'Noto Sans Coptic', sans-serif;
}

.clickable-word:hover {
  text-decoration: underline;
}

.text-right {
  text-align: right;
}

.list-item {
  color: inherit; /* Ensure the list item text inherits the color from the parent */
  white-space: normal; /* Allow text to wrap */
}

</style>
