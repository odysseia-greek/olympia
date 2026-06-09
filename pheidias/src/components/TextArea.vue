<template>
  <v-app id="text" :style="{ background: $vuetify.theme.themes[theme].background }">
    <v-main>
      <header class="texts-hero" :style="{ backgroundImage: `url(${textsHeroImage})` }">
        <div class="texts-hero-shade">
          <v-container class="texts-hero-content">
          <section class="texts-search-panel" aria-labelledby="texts-heading">
            <div class="panel-heading">
              <div>
                <div class="section-label">Texts</div>
                <h1 id="texts-heading">Herodotos</h1>
                <p>Select an author, choose a reference, translate the Greek, and compare your work with available translations.</p>
              </div>
              <v-btn
                  color="primary"
                  variant="tonal"
                  prepend-icon="mdi-information"
                  @click="showInfoBar = !showInfoBar"
              >
                Help
              </v-btn>
            </div>

            <v-autocomplete
                v-model="searchQuery"
                :items="autocompleteAuthorsAndBooks"
                color="primary"
                label="Search authors and books"
                placeholder="Try Herodotos, Plato, Histories..."
                prepend-inner-icon="mdi-magnify"
                variant="outlined"
                @keyup.enter="handleKeyPress"
                clearable
            ></v-autocomplete>

            <div class="search-actions">
              <v-btn
                  color="secondary"
                  variant="flat"
                  prepend-icon="mdi-bookshelf"
                  @click="toggleExpandAll"
              >
                {{ expandAll ? 'Collapse library' : 'Expand library' }}
              </v-btn>
              <div v-if="selectedAuthor && selectedBook" class="selected-path">
                {{ selectedAuthor }} · {{ selectedBook }}
                <span v-if="selectedReference"> · {{ selectedReference }}</span>
              </div>
            </div>

            <v-expand-transition>
              <v-sheet v-if="showInfoBar" class="help-panel" color="secondaryPapyrus">
                <div class="help-title">How text translation works</div>
                <ol>
                  <li>Select an author and book, then choose a reference and optional section.</li>
                  <li>Translate the Greek passages, then press <strong>Check</strong> to compare your translation.</li>
                  <li>Greek words are clickable and open lexical/grammar detail without leaving the passage.</li>
                  <li>Results show scores, official translations, possible typos, and quick copy/fullscreen actions.</li>
                </ol>
              </v-sheet>
            </v-expand-transition>
          </section>
          </v-container>
        </div>
      </header>

      <main class="texts-content">
        <v-container class="content-container">
          <section v-if="filteredAuthors.length" class="library-section">
            <div class="section-heading">
              <div>
                <v-chip color="secondary" variant="tonal">Library</v-chip>
                <h2>Choose author and book</h2>
              </div>
              <p>Use search for a fast path, or browse the available authors and books directly.</p>
            </div>

            <v-sheet class="text-panel library-panel" color="secondaryPapyrus">
              <v-expansion-panels v-model="expandedPanels" multiple>
                <v-expansion-panel
                    v-for="author in filteredAuthors"
                    :key="author.key"
                    class="author-panel"
                >
                  <v-expansion-panel-title>
                    <strong>{{ author.key }}</strong>
                  </v-expansion-panel-title>
                  <v-expansion-panel-text>
                    <div class="book-grid">
                      <v-btn
                          v-for="book in author.filteredBooks.length ? author.filteredBooks : author.books"
                          :key="book.key"
                          class="book-button"
                          :class="{ 'is-selected': selectedAuthor === author.key && selectedBook === book.key }"
                          :color="selectedAuthor === author.key && selectedBook === book.key ? 'primary' : 'triadic'"
                          :variant="selectedAuthor === author.key && selectedBook === book.key ? 'flat' : 'tonal'"
                          @click="onBookSelected(author.key, book.key)"
                      >
                        {{ book.key }}
                      </v-btn>
                    </div>
                  </v-expansion-panel-text>
                </v-expansion-panel>
              </v-expansion-panels>
            </v-sheet>
          </section>

          <section v-if="selectedBookReferences.length" ref="referenceSectionRef" class="reference-section">
            <div class="section-heading">
              <div>
                <v-chip color="triadic" variant="tonal">Passage</v-chip>
                <h2>Choose reference</h2>
              </div>
              <p>References point to the text passages. Sections let you narrow a longer reference before translating.</p>
            </div>

            <v-sheet class="text-panel" color="secondaryPapyrus">
              <div class="chip-group">
                <v-btn
                    v-for="reference in selectedBookReferences"
                    :key="reference.key"
                    class="reference-button"
                    :class="{ 'pulsate': selectedReference === reference.key }"
                    :color="selectedReference === reference.key ? 'primary' : 'triadic'"
                    :variant="selectedReference === reference.key ? 'flat' : 'tonal'"
                    @click="onReferenceSelected(selectedBook, reference.key)"
                >
                  {{ reference.key }}
                </v-btn>
              </div>

              <div
                  v-if="selectedReferenceSections.length && selectedReferenceSections[0].key !==''"
                  ref="sectionPickerRef"
                  class="section-picker"
              >
                <h3 ref="loadingResultsRef">Sections</h3>
                <div class="chip-group">
                  <v-btn
                      v-for="section in selectedReferenceSections"
                      :key="section.key"
                      class="reference-button"
                      :class="{ 'pulsate': selectedSection === section.key }"
                      :color="selectedSection === section.key ? 'primary' : 'triadic'"
                      :variant="selectedSection === section.key ? 'flat' : 'tonal'"
                      @click="onSectionSelected(section.key)"
                  >
                    {{ section.key }}
                  </v-btn>
                  <v-btn
                      class="reference-button"
                      :color="selectedSectionIndex === -1 ? 'primary' : 'triadic'"
                      :variant="selectedSectionIndex === -1 ? 'flat' : 'tonal'"
                      @click="onSectionSelected('')"
                  >
                    All sections
                  </v-btn>
                </div>
              </div>
            </v-sheet>
          </section>

          <section v-if="resultData" ref="translationSectionRef" class="translation-section">
            <div class="section-heading">
              <div>
                <v-chip color="primary" variant="tonal">Translate</v-chip>
                <h2>Work through the Greek</h2>
              </div>
              <p>Click a Greek word for detail, then write your translation beside each section.</p>
            </div>

            <v-card class="translation-card paper-card" width="100%">
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
                        @update:model-value="hideTranslationError"
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
                    <v-card class="results-card paper-card" width="100%">
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
                            <div v-if="possibleTypos.length">
                              <p style="margin-top:2em"><strong>Possible Typos:</strong></p>
                              <v-list dense class="section-content">
                                <v-list-item v-for="typo in possibleTypos" :key="typo.provided" class="list-item">
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
          </section>
        </v-container>
      </main>

      <v-snackbar
          v-model="translationErrorVisible"
          color="error"
          location="bottom"
          timeout="3600"
      >
        {{ translationError?.message || 'Unable to check your translation.' }}
        <template v-slot:actions>
          <v-btn variant="text" @click="hideTranslationError">Close</v-btn>
        </template>
      </v-snackbar>
    </v-main>
  </v-app>
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
    const textsHeroImage = ref('');
    const referenceSectionRef = ref();
    const sectionPickerRef = ref();
    const translationSectionRef = ref();
    const resultsContainerRef = ref();
    const loadingResultsRef = ref();
    const pendingScrollTarget = ref('');

    const resultData = ref(null);
    const queryLoading = ref(false);
    const queryError = ref(null);

    const clickedWord = ref('');
    const grammarResults = ref([]);
    const grammarError = ref(null);
    const translations = ref({});
    const translationResults = ref(null);
    const translationError = ref(null);
    const translationErrorVisible = ref(false);
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
            expandedPanels.value = [index];
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
              if (pendingScrollTarget.value) {
                scrollMeTo(pendingScrollTarget.value);
                pendingScrollTarget.value = '';
              }
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

    const loadHeroImage = () => {
      import('@/assets/alexandria.webp').then((module) => {
        textsHeroImage.value = module.default;
      });
    };

    const scrollMeTo = (target) => {
      nextTick(() => {
        const targets = {
          references: referenceSectionRef.value,
          sections: sectionPickerRef.value,
          translation: translationSectionRef.value,
        };

        if (targets[target]) {
          targets[target].scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
      });
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
      hideTranslationError();
    };

    const showTranslationError = (message) => {
      translationError.value = { message };
      translationErrorVisible.value = true;
    };

    const hideTranslationError = () => {
      translationErrorVisible.value = false;
      translationError.value = null;
    };

    const totalGreekTextLength = computed(() => {
      return (resultData.value?.create?.rhemai || [])
          .map((rhema) => rhema.greek || '')
          .join(' ')
          .replace(/\s+/g, '')
          .length;
    });

    const totalTranslationLength = computed(() => {
      return Object.values(translations.value || {})
          .join(' ')
          .trim()
          .replace(/\s+/g, '')
          .length;
    });

    const possibleTypos = computed(() => {
      return Array.isArray(translationResults.value?.possibleTypos)
          ? translationResults.value.possibleTypos
          : [];
    });

    const translationIsTooShort = computed(() => {
      const greekLength = totalGreekTextLength.value;
      if (!greekLength) return false;

      return totalTranslationLength.value < Math.ceil(greekLength * 0.1);
    });

    const checkTranslations = async () => {
      translationResults.value = null;
      hideTranslationError();

      if (translationIsTooShort.value) {
        showLoading.value = false;
        loadingPercentage.value = 0;
        showTranslationError('Your entry is too short. Add a fuller attempt before checking.');
        return;
      }

      scrollToResults('loadingResults');
      showLoading.value = true;
      loadingPercentage.value = 0;

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
          }
        });
        onError((error) => {
          showLoading.value = false;
          translationResults.value = null;
          showTranslationError(error.message || 'Unable to check your translation.');
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
        showTranslationError(error.message || 'Unable to check your translation.');
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
      hideTranslationError();
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

      scrollMeTo('references');
    };

    const onReferenceSelected = (bookKey, referenceKey) => {
      translationResults.value = null;
      hideTranslationError();
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

      if (selectedReferenceSections.value.length && selectedReferenceSections.value[0].key !== '') {
        scrollMeTo('sections');
      } else {
        pendingScrollTarget.value = 'translation';
      }
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
      hideTranslationError();
      selectedSection.value = sectionKey;
      selectedSectionIndex.value = selectedReferenceSections.value.findIndex((s) => s.key === sectionKey);
      resultData.value = null;
      translations.value = {};
      pendingScrollTarget.value = 'translation';
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
      loadHeroImage();
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
      textsHeroImage,
      referenceSectionRef,
      sectionPickerRef,
      translationSectionRef,
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
      translationErrorVisible,
      possibleTypos,
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
      scrollMeTo,
      handleKeyPress,
      checkTranslations,
      clearTranslations,
      hideTranslationError,
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
#text {
  --texts-primary: #1c61d1;
  --texts-secondary: #1cd18c;
  --texts-triadic: #1cbcd1;
  --texts-ink: #20334f;
  --texts-muted: #536987;
  color: var(--texts-ink);
}

* {
  box-sizing: border-box;
}

.texts-hero {
  min-height: 470px;
  background-position: center 42%;
  background-size: cover;
}

.texts-hero-shade {
  min-height: 470px;
  display: flex;
  align-items: center;
  background:
      linear-gradient(90deg, rgba(22, 20, 15, 0.72) 0%, rgba(38, 31, 21, 0.46) 48%, rgba(26, 22, 16, 0.16) 100%),
      linear-gradient(180deg, rgba(0, 0, 0, 0.06), rgba(26, 21, 14, 0.36));
}

.texts-hero-content {
  padding-top: 58px;
  padding-bottom: 44px;
}

.texts-search-panel,
.text-panel,
.translation-card {
  border: 1px solid rgba(28, 97, 209, 0.16);
  border-radius: 8px;
  box-shadow: 0 14px 36px rgba(11, 39, 85, 0.16);
}

.texts-search-panel {
  max-width: 1060px;
  padding: 24px;
  background: rgba(253, 246, 227, 0.96);
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

.section-heading p,
.panel-heading p {
  max-width: 560px;
  margin: 0;
  color: #344765;
  line-height: 1.65;
}

.panel-heading p {
  margin-top: 10px;
}

.section-label {
  color: #64789e;
  font-size: 0.82rem;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.panel-heading h1,
.section-heading h2 {
  margin: 8px 0 0;
  font-size: clamp(1.65rem, 3vw, 2.45rem);
  line-height: 1.12;
}

.search-actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 14px;
  margin-top: 8px;
}

.selected-path {
  color: var(--texts-muted);
  font-weight: 800;
}

.help-panel {
  margin-top: 18px;
  padding: 18px;
  color: var(--texts-ink);
  border: 1px solid rgba(28, 188, 209, 0.2);
  border-radius: 8px;
}

.help-title {
  margin-bottom: 8px;
  font-weight: 800;
}

.help-panel ol {
  margin: 0;
  padding-left: 20px;
  line-height: 1.65;
}

.texts-content {
  background:
      linear-gradient(150deg, rgba(28, 97, 209, 0.16) 0%, rgba(28, 188, 209, 0.12) 34%, rgba(28, 209, 140, 0.1) 62%, rgba(254, 252, 245, 0.98) 100%),
      linear-gradient(180deg, #d5eff7 0%, #f2fbf7 46%, #fefcf5 100%);
}

.content-container {
  max-width: 1240px;
  padding-top: 54px;
  padding-bottom: 58px;
}

.library-section,
.reference-section,
.translation-section {
  margin-bottom: 46px;
}

.reference-section,
.section-picker,
.translation-section {
  scroll-margin-top: 80px;
}

.text-panel {
  padding: 18px;
  color: var(--texts-ink);
  overflow: hidden;
}

.library-panel :deep(.v-expansion-panel) {
  color: var(--texts-ink);
  background: rgba(255, 255, 255, 0.56);
}

.library-panel :deep(.v-expansion-panel-title) {
  color: var(--texts-ink);
}

.book-grid,
.chip-group {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.book-button,
.reference-button {
  text-transform: none;
}

.book-button.is-selected {
  box-shadow: 0 8px 22px rgba(28, 97, 209, 0.24);
}

.section-picker {
  margin-top: 22px;
  padding-top: 20px;
  border-top: 1px solid rgba(28, 97, 209, 0.14);
}

.section-picker h3 {
  margin: 0 0 14px;
}

.paper-card,
.translation-card {
  background: #fdf6e3;
  color: var(--texts-ink);
  border-radius: 8px;
  padding: 24px;
  font-family: 'Roboto', serif;
}

.results-card {
  margin-top: 18px;
  background: #fefcf5;
}

.section-content {
  background: #fefcf5;
  color: var(--texts-ink);
  padding: 10px;
  border-radius: 4px;
}

.rhema-section {
  margin-bottom: 20px;
  padding: 18px 0;
  border-bottom: 1px solid rgba(28, 97, 209, 0.12);
}

.rhema-section p {
  color: var(--texts-ink);
  line-height: 1.8;
}

.clickable-word {
  cursor: pointer;
  color: #10284b;
  font-family: 'Noto Sans Coptic', serif;
  transition: color 0.16s ease, background-color 0.16s ease;
}

.clickable-word:hover {
  color: var(--texts-primary);
  background: rgba(28, 209, 140, 0.2);
  text-decoration: underline;
}

.v-card .v-card-text {
  max-width: 100%;
  overflow: hidden;
}

.text-right {
  text-align: right;
}

.list-item {
  color: inherit;
  white-space: normal;
}

@media (max-width: 900px) {
  .texts-hero,
  .texts-hero-shade {
    min-height: auto;
  }

  .texts-hero-content {
    padding-top: 34px;
    padding-bottom: 28px;
  }

  .panel-heading,
  .section-heading {
    display: block;
  }

  .panel-heading .v-btn,
  .section-heading p {
    margin-top: 12px;
  }
}

</style>
