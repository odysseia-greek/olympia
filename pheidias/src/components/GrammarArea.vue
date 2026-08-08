<template>
  <v-app id="grammar" :style="{ background: $vuetify.theme.themes[theme].background }">
    <v-main>
      <header class="grammar-hero" :style="{ backgroundImage: `url(${grammarHeroImage})` }">
        <div class="grammar-hero-shade">
          <v-container class="grammar-hero-content">
            <section class="grammar-search-panel" aria-labelledby="grammar-heading">
              <div class="panel-heading">
                <div>
                  <div class="section-label">Grammar</div>
                  <h1 id="grammar-heading">Dionysios</h1>
                  <p>
                    Enter a Greek form to inspect possible rules, roots, and translations.
                    Useful when reading a passage and a declined form needs unpacking.
                  </p>
                </div>
                <v-btn
                    aria-label="Grammar information"
                    color="primary"
                    icon="mdi-information"
                    variant="tonal"
                    @click="infoDialogVisible = true"
                ></v-btn>
              </div>

              <v-btn-toggle v-model="grammarMode" class="grammar-mode-toggle" color="primary" mandatory variant="outlined">
                <v-btn value="word" prepend-icon="mdi-text-short" size="large">Single word</v-btn>
                <v-btn value="sentence" prepend-icon="mdi-text-long" size="large">Sentence</v-btn>
              </v-btn-toggle>

              <v-switch
                  v-model="auditMode"
                  class="audit-switch"
                  color="secondary"
                  label="Audit available texts for occurrences and sources"
                  hide-details
                  inset
              ></v-switch>

              <v-autocomplete
                  v-if="grammarMode === 'word'"
                  :loading="loading"
                  v-model="queryWord"
                  @update:search="updateFromHistory"
                  :items="historyOfWords"
                  color="primary"
                  label="What Greek word are you trying to decline?"
                  placeholder="Try λόγων, ἔβαλλε, πόλεως..."
                  prepend-inner-icon="mdi-magnify"
                  variant="outlined"
                  @keyup.enter="grammarWord($event.target.value)"
                  clearable
              ></v-autocomplete>

              <div v-if="grammarMode === 'word'" class="example-row">
                <span>Examples</span>
                <v-chip
                    v-for="word in historyOfWords.slice(0, 6)"
                    :key="word"
                    color="secondary"
                    variant="tonal"
                    @click="grammarWord(word)"
                >
                  {{ word }}
                </v-chip>
              </div>

              <div v-else class="sentence-input">
                <v-select
                    :items="sentenceExamples"
                    color="secondary"
                    label="Try an example sentence"
                    prepend-inner-icon="mdi-format-quote-open"
                    variant="outlined"
                    @update:model-value="sentenceText = $event; sentenceError = ''"
                ></v-select>
                <v-textarea
                    v-model="sentenceText"
                    :disabled="sentenceLoading"
                    :rules="sentenceRules"
                    auto-grow
                    color="primary"
                    :hint="`${sentenceWordCount} of 50 words · minimum 3`"
                    label="Ancient Greek sentence"
                    placeholder="Paste or type a Greek sentence to analyze..."
                    persistent-hint
                    rows="3"
                    variant="outlined"
                    @update:model-value="sentenceError = ''"
                    @keydown.ctrl.enter="analyzeSentence"
                    @keydown.meta.enter="analyzeSentence"
                ></v-textarea>
                <div class="sentence-actions">
                  <span>Ctrl/⌘ + Enter to analyze</span>
                  <v-btn
                      :disabled="!sentenceIsValid"
                      :loading="sentenceLoading"
                      color="primary"
                      prepend-icon="mdi-text-search"
                      @click="analyzeSentence"
                  >
                    Analyze sentence
                  </v-btn>
                </div>
                <v-alert v-if="sentenceError" class="mt-4" type="error" variant="tonal">
                  {{ sentenceError }}
                </v-alert>
              </div>

            </section>
          </v-container>
        </div>
      </header>

      <main class="grammar-content">
        <v-container class="content-container">
          <section v-if="grammarMode === 'word'" class="results-section">
            <div class="section-heading">
              <div>
                <v-chip color="secondary" variant="tonal">Results</v-chip>
                <h2>Grammar analysis</h2>
              </div>
              <p>
                Results show the queried form, matching rule, likely root, and available translations.
              </p>
            </div>

            <v-sheet class="grammar-panel results-panel" color="secondaryPapyrus">
              <v-data-table
                  dense
                  :headers="headers"
                  :items="grammarResults"
                  :items-per-page="10"
                  item-key="name"
              >
                <template v-slot:item.translations="{ item }">
                  <ol class="translation-list">
                    <li v-for="(trans, index) in item.translations" :key="index">
                      {{ trans }}
                    </li>
                  </ol>
                </template>
                <template v-slot:item.actions="{ item }">
                  <div class="dictionary-actions">
                    <v-btn
                        :href="dictionaryUrl(item.rootWord || item.word, false)"
                        color="primary"
                        prepend-icon="mdi-book-open-variant"
                        rel="noopener noreferrer"
                        size="small"
                        target="_blank"
                        variant="tonal"
                    >
                      Dictionary
                    </v-btn>
                    <v-btn
                        :href="dictionaryUrl(item.rootWord || item.word, true)"
                        color="secondary"
                        prepend-icon="mdi-text-box-search-outline"
                        rel="noopener noreferrer"
                        size="small"
                        target="_blank"
                        variant="tonal"
                    >
                      Find in texts
                    </v-btn>
                  </div>
                </template>
              </v-data-table>
            </v-sheet>
          </section>

          <section v-if="grammarMode === 'word' && grammarAudit" class="text-analysis-section">
            <div class="section-heading">
              <div>
                <v-chip color="triadic" variant="tonal">Grammar audit</v-chip>
                <h2>How Dionysios resolved the word</h2>
              </div>
              <p>{{ grammarAudit.reason }}</p>
            </div>
            <v-sheet class="grammar-panel" color="secondaryPapyrus">
              <div class="grammar-audit-summary">
                <div>
                  <span>Word</span>
                  <strong class="audit-word">{{ grammarAudit.word }}</strong>
                </div>
                <div>
                  <span>Outcome</span>
                  <v-chip color="secondary" size="small" variant="tonal">{{ grammarAudit.outcome }}</v-chip>
                </div>
                <div>
                  <span>Source</span>
                  <v-chip color="primary" size="small" variant="tonal">{{ grammarAudit.source || 'unknown' }}</v-chip>
                </div>
              </div>

              <v-timeline class="grammar-audit-events" density="compact" side="end" truncate-line="both">
                <v-timeline-item
                    v-for="(event, index) in grammarAudit.events"
                    :key="`${event.step}-${index}`"
                    :dot-color="event.status === 'ok' ? 'secondary' : 'warning'"
                    size="small"
                >
                  <v-card class="audit-event-card" variant="outlined">
                    <v-card-title>{{ event.step }}</v-card-title>
                    <v-card-subtitle>
                      {{ event.source || 'Dionysios' }}
                      <template v-if="event.resultCount"> · {{ event.resultCount }} result(s)</template>
                      <template v-if="event.candidateCount"> · {{ event.candidateCount }} candidate(s)</template>
                    </v-card-subtitle>
                    <v-card-text>
                      <p>{{ event.reason }}</p>
                      <p v-if="event.rootWord || event.rule" class="audit-event-meta">
                        {{ event.rootWord }}<template v-if="event.rootWord && event.rule"> · </template>{{ event.rule }}
                      </p>
                      <ul v-if="event.details?.length" class="audit-detail-list">
                        <li v-for="detail in event.details" :key="detail">{{ detail }}</li>
                      </ul>
                    </v-card-text>
                  </v-card>
                </v-timeline-item>
              </v-timeline>
            </v-sheet>
          </section>

          <section v-if="grammarMode === 'sentence' && sentenceResult" class="results-section sentence-results">
            <div class="section-heading">
              <div>
                <v-chip color="secondary" variant="tonal">TextMode</v-chip>
                <h2>Sentence analysis</h2>
              </div>
              <p>Select a word to inspect every grammatical interpretation Dionysios found.</p>
            </div>

            <v-sheet class="grammar-panel" color="secondaryPapyrus">
              <div class="interlinear-sentence" aria-label="Interlinear sentence analysis">
                <button
                    v-for="token in sentenceResult.tokens"
                    :key="`${token.position}-${token.token}`"
                    class="interlinear-token"
                    :class="{
                      unresolved: !token.resolved,
                      selected: selectedSentenceToken?.position === token.position,
                    }"
                    type="button"
                    @click="selectedTokenPosition = token.position"
                >
                  <span class="greek-token">{{ token.token }}</span>
                  <span class="token-gloss">{{ token.gloss || token.token }}</span>
                </button>
              </div>

              <v-divider class="my-5"></v-divider>

              <div class="literal-translation">
                <span>Literal translation</span>
                <p>{{ sentenceResult.literalTranslation }}</p>
              </div>

              <div v-if="selectedSentenceToken" class="token-details">
                <div class="token-detail-heading">
                  <div>
                    <span class="greek-token">{{ selectedSentenceToken.token }}</span>
                    <span class="selected-gloss">{{ selectedSentenceToken.gloss }}</span>
                  </div>
                  <v-chip :color="selectedSentenceToken.resolved ? 'secondary' : 'error'" size="small" variant="tonal">
                    {{ selectedSentenceToken.resolved ? `${selectedSentenceToken.results.length} interpretation(s)` : 'Unresolved' }}
                  </v-chip>
                </div>

                <v-alert
                    v-if="!selectedSentenceToken.resolved"
                    class="mt-4"
                    type="warning"
                    variant="tonal"
                >
                  {{ selectedSentenceToken.message || 'No grammatical options found.' }}
                </v-alert>

                <div v-else class="interpretation-grid">
                  <v-card
                      v-for="(result, index) in selectedSentenceToken.results"
                      :key="`${result.rootWord}-${result.rule}-${index}`"
                      class="interpretation-card"
                      variant="outlined"
                  >
                    <v-card-title class="interpretation-root">{{ result.rootWord }}</v-card-title>
                    <v-card-subtitle>{{ result.rule }}</v-card-subtitle>
                    <v-card-text>{{ result.translations?.join(' · ') || 'No translation found' }}</v-card-text>
                  </v-card>
                </div>
              </div>
            </v-sheet>
          </section>

          <section v-if="grammarMode === 'sentence' && auditMode && activeTextSearch?.searched" class="text-analysis-section">
            <div class="section-heading">
              <div>
                <v-chip color="triadic" variant="tonal">Audit</v-chip>
                <h2>Sources in available texts</h2>
              </div>
              <p>{{ activeTextSearch.matchCount }} matching passage(s) found for “{{ activeTextSearch.query }}”.</p>
            </div>
            <v-alert v-if="!activeTextSearch.found" type="info" variant="tonal">
              {{ activeTextSearch.message || 'No matching passages were found.' }}
            </v-alert>
            <div v-else class="audit-results">
              <v-card
                  v-for="(match, index) in activeTextSearch.matches"
                  :key="`${match.author}-${match.book}-${match.reference}-${index}`"
                  class="audit-card"
                  variant="outlined"
              >
                <v-card-title>{{ match.author }} · {{ match.book }}</v-card-title>
                <v-card-subtitle>Reference {{ match.reference }} · section {{ match.text?.section }}</v-card-subtitle>
                <v-card-text>
                  <p class="audit-greek" v-html="highlightWord(match.text?.greek)"></p>
                  <p v-for="translation in (match.text?.translations || [])" :key="translation" class="audit-translation">
                    {{ translation }}
                  </p>
                </v-card-text>
                <v-card-actions v-if="match.referenceLink">
                  <v-btn :href="match.referenceLink" color="primary" target="_blank" variant="text">
                    Open interactive text
                  </v-btn>
                </v-card-actions>
              </v-card>
            </div>
          </section>
        </v-container>
      </main>

      <v-dialog v-model="infoDialogVisible" max-width="860">
        <v-card class="info-card">
          <v-card-title class="headline">Grammar Conjugation</v-card-title>
          <v-card-text>
            <v-list>
              <v-list-item>
                <v-list-item-title class="subtitle-1">
                  Enter a Greek form to see possible declensions, rule matches, and root words.
                </v-list-item-title>
              </v-list-item>
              <v-divider></v-divider>
              <v-list-item>
                <v-list-item-title><strong>Search Input:</strong></v-list-item-title>
                <v-list-item-subtitle>
                  Try forms such as λόγων, ἔβαλλε, φέροντος, ἀληθῆ, Ἀθηναῖος, or πόλεως.
                </v-list-item-subtitle>
              </v-list-item>
              <v-divider></v-divider>
              <v-list-item>
                <v-list-item-title><strong>Text audit:</strong></v-list-item-title>
                <v-list-item-subtitle>
                  Enable audit mode to search available texts and show the source of every matching passage.
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" variant="text" @click="infoDialogVisible = false">Close</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      </v-main>
  </v-app>
</template>

<script>
import { computed, getCurrentInstance, onMounted, ref, watch } from 'vue';
import { useApolloClient, useQuery } from '@vue/apollo-composable';
import { AnalyzeSentence, CheckGrammar } from '@/constants/grammarGraphql';
import { useBouleId } from '@/composables/useBoule';

export default {
  name: 'GrammarArea',
  setup() {
    const { proxy } = getCurrentInstance();
    const { client } = useApolloClient();
    const theme = ref('light');
    const stepper = ref(1);
    const infoDialogVisible = ref(false);
    const grammarHeroImage = ref('');
    const grammarResults = ref([]);
    const errors = ref([]);
    const historyOfWords = ref([
      "ἔβαλλε",
      "φέροντος",
      "ἀληθῆ",
      "λόγων",
      "Ἀθηναῖος",
      "πόλεως"
    ]);
    const loading = ref(false);
    const queryWord = ref('');
    const grammarMode = ref('word');
    const sentenceExamples = [
      'μήτε ἔργα μεγάλα τε καὶ θωμαστά, τὰ μὲν Ἕλλησι τὰ δὲ βαρβάροισι ἀποδεχθέντα, ἀκλεᾶ γένηται',
      'καὶ ὁ λόγος ἦν πρὸς τὸν θεόν',
      'τὰ ἴδια διάφορα πᾶσι τὸ ἴσον, κατὰ δὲ τὴν ἀξίωσιν',
      'νόμος οὕτω κελεύει νομοθετεῖν, γράφεσθαι μέν, ἄν τίς τινα τῶν ὑπαρχόντων νόμων μὴ καλῶς ἔχειν ἡγῆται',
    ];
    const sentenceText = ref(sentenceExamples[0]);
    const sentenceWordCount = computed(() => {
      const text = sentenceText.value.trim();
      return text ? text.split(/\s+/u).length : 0;
    });
    const sentenceIsValid = computed(() => sentenceWordCount.value >= 3 && sentenceWordCount.value <= 50);
    const sentenceRules = [
      () => sentenceWordCount.value >= 3 || 'Enter at least 3 words.',
      () => sentenceWordCount.value <= 50 || 'Enter no more than 50 words.',
    ];
    const sentenceLoading = ref(false);
    const sentenceError = ref('');
    const sentenceResult = ref(null);
    const sentenceSessionId = ref(useBouleId());
    const selectedTokenPosition = ref(null);
    const selectedSentenceToken = computed(() => {
      const tokens = sentenceResult.value?.tokens || [];
      return tokens.find(token => token.position === selectedTokenPosition.value) || tokens[0] || null;
    });
    const auditMode = ref(false);
    const grammarAudit = ref(null);
    const activeTextSearch = computed(() => sentenceResult.value?.textSearch);

    const headers = [
      {title: 'Queried', align: 'start', sortable: true, value: 'word'},
      {title: 'Rule', value: 'rule'},
      {title: 'Root', value: 'rootWord'},
      {title: 'Translation(s)', value: 'translations' },
      {title: 'Continue', value: 'actions', sortable: false},
    ];

    const loadHeroImage = () => {
      import('@/assets/grammar.webp').then((module) => {
        grammarHeroImage.value = module.default;
      });
    };

    const analyzeSentence = async (syncUrl = true) => {
      const text = sentenceText.value.trim();
      if (sentenceLoading.value) return;
      if (!sentenceIsValid.value) {
        sentenceError.value = sentenceWordCount.value < 3
            ? 'Enter at least 3 words for sentence analysis.'
            : 'Sentence analysis accepts no more than 50 words.';
        return;
      }

      sentenceLoading.value = true;
      sentenceError.value = '';
      sentenceResult.value = null;
      selectedTokenPosition.value = null;

      if (syncUrl !== false) {
        updateUrl({
          mode: 'sentence',
          text,
          word: null,
          audit: String(auditMode.value),
        });
      }

      try {
        const context = sentenceSessionId.value
            ? { headers: { boule: sentenceSessionId.value } }
            : undefined;
        const { data } = await client.query({
          query: AnalyzeSentence,
          variables: { text, includeAudit: auditMode.value },
          fetchPolicy: 'no-cache',
          context,
        });

        sentenceResult.value = data?.sentence || null;
        sentenceSessionId.value = data?.sentence?.sessionId || sentenceSessionId.value;
        selectedTokenPosition.value = data?.sentence?.tokens?.[0]?.position ?? null;
      } catch (error) {
        sentenceError.value = error.message || 'The sentence could not be analyzed.';
      } finally {
        sentenceLoading.value = false;
      }
    };

    const updateFromHistory = async () => {
      if (historyOfWords.value.includes(queryWord.value)) {
        await grammarWord(queryWord.value)
      }
    }

    const grammarWord = async (word, syncUrl = true) => {
      loading.value = true;
      grammarAudit.value = null;
      if (syncUrl !== false) {
        updateUrl({
          mode: 'word',
          word,
          text: null,
          audit: String(auditMode.value),
        });
      }

      try {
        const { onResult, onError } = useQuery(CheckGrammar, {
          word: word,
          includeAudit: auditMode.value,
        });

        if (!historyOfWords.value.includes(word)) {
          historyOfWords.value.push(word);
        }

        onResult((response) => {
          loading.value = false;
          grammarAudit.value = response?.data?.grammar?.audit || null;
          if (response?.data?.grammar?.results) {
            grammarResults.value = response.data.grammar.results.map(item => {
              if (!item.translations || item.translations.length === 0) {
                item.translations = ['No translation found'];
              }

              return item;
            });

          } else {
            grammarResults.value = [];
          }
        });


        onError((error) => {
          setTimeout(() => {
            loading.value = false;
          }, 1500);

          grammarResults.value = [{
            word: word,
            translation: ['No translation found'],
            rootWord: word,
            rule: 'No rule found'
          }];
          errors.value.push(error);
        });
      } catch (error) {
        console.error('Unexpected error:', error); // Log unexpected errors
        setTimeout(() => {
          loading.value = false;
        }, 1500);

        grammarResults.value = [{
          word: word,
          translation: ['No translation found'],
          rootWord: word,
          rule: 'No rule found'
        }];
        errors.value.push(error);
      }
    };

    let suppressNextRouteRestore = false;

    const restoreFromURL = async (query) => {
      const mode = query.mode === 'sentence' ? 'sentence' : 'word';
      grammarMode.value = mode;
      auditMode.value = String(query.audit).toLowerCase() === 'true';

      if (mode === 'sentence' && query.text) {
        sentenceText.value = String(query.text);
        await analyzeSentence(false);
      } else if (query.word) {
        queryWord.value = String(query.word);
        await grammarWord(queryWord.value, false);
      } else {
        queryWord.value = '';
        grammarResults.value = [];
        grammarAudit.value = null;
        sentenceResult.value = null;
        sentenceError.value = '';
      }
    };

    const initializeFromURL = () => restoreFromURL(proxy.$route.query);

    const updateUrl = (query) => {
      const currentQuery = proxy.$route.query;
      const newQuery = { ...currentQuery, ...query };
      Object.keys(newQuery).forEach((key) => {
        if (newQuery[key] === null || newQuery[key] === undefined || newQuery[key] === '') {
          delete newQuery[key];
        }
      });
      const queryKeys = new Set([...Object.keys(currentQuery), ...Object.keys(newQuery)]);
      const queryChanged = [...queryKeys].some(key => currentQuery[key] !== newQuery[key]);

      if (queryChanged) {
        suppressNextRouteRestore = true;
        proxy.$router.push({ name: 'Dionysios', query: newQuery });
      }
    };

    const highlightWord = (text) => {
      if (!text) return '';
      const regex = /&&&(.*?)&&&/g;
      return text.replace(regex, '<span style="background-color: yellow;">$1</span>');
    };

    const dictionaryUrl = (word, extended) => {
      const params = new URLSearchParams({
        mode: 'exact',
        language: 'greek',
        extended: String(extended),
        word: word || '',
      });
      return `/dictionary?${params.toString()}`;
    };

    onMounted(() => {
      loadHeroImage();
      initializeFromURL();
    });

    watch(
      () => proxy.$route.query,
      (query) => {
        if (suppressNextRouteRestore) {
          suppressNextRouteRestore = false;
          return;
        }
        restoreFromURL(query);
      },
      { deep: true },
    );

    return {
      theme,
      stepper,
      grammarResults,
      grammarHeroImage,
      errors,
      headers,
      infoDialogVisible,
      historyOfWords,
      loading,
      queryWord,
      grammarMode,
      sentenceExamples,
      sentenceText,
      sentenceWordCount,
      sentenceIsValid,
      sentenceRules,
      sentenceLoading,
      sentenceError,
      sentenceResult,
      selectedTokenPosition,
      selectedSentenceToken,
      analyzeSentence,
      auditMode,
      grammarAudit,
      activeTextSearch,
      grammarWord,
      highlightWord,
      dictionaryUrl,
      updateFromHistory,
      initializeFromURL,
      updateUrl,
      loadHeroImage,

    };
  }
};
</script>

<style scoped>
#grammar {
  --grammar-primary: #1c61d1;
  --grammar-secondary: #1cd18c;
  --grammar-triadic: #1cbcd1;
  --grammar-ink: #20334f;
  --grammar-muted: #536987;
  color: var(--grammar-ink);
}

* {
  box-sizing: border-box;
}

a {
  cursor: pointer;
}

.grammar-hero {
  min-height: 470px;
  background-position: center 36%;
  background-size: cover;
}

.grammar-hero-shade {
  min-height: 470px;
  display: flex;
  align-items: center;
  background:
      linear-gradient(90deg, rgba(22, 20, 15, 0.72) 0%, rgba(38, 31, 21, 0.46) 48%, rgba(26, 22, 16, 0.16) 100%),
      linear-gradient(180deg, rgba(0, 0, 0, 0.06), rgba(26, 21, 14, 0.36));
}

.grammar-hero-content {
  padding-top: 58px;
  padding-bottom: 44px;
}

.grammar-search-panel,
.grammar-panel {
  border: 1px solid rgba(28, 97, 209, 0.16);
  border-radius: 8px;
  box-shadow: 0 14px 36px rgba(11, 39, 85, 0.16);
}

.grammar-search-panel {
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

.example-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 10px;
  color: var(--grammar-muted);
  font-weight: 800;
}

.grammar-mode-toggle {
  height: auto;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

.grammar-mode-toggle :deep(.v-btn) {
  min-width: 190px;
  min-height: 54px;
  padding-inline: 26px;
  font-size: 1rem;
  font-weight: 750;
}

.audit-switch {
  margin-top: -8px;
  margin-bottom: 18px;
}

.sentence-input {
  margin-top: 4px;
}

.sentence-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-top: -6px;
  color: var(--grammar-muted);
  font-size: 0.85rem;
}

.grammar-content {
  background:
      linear-gradient(150deg, rgba(28, 97, 209, 0.16) 0%, rgba(28, 188, 209, 0.12) 34%, rgba(28, 209, 140, 0.1) 62%, rgba(254, 252, 245, 0.98) 100%),
      linear-gradient(180deg, #d5eff7 0%, #f2fbf7 46%, #fefcf5 100%);
}

.content-container {
  max-width: 1240px;
  padding-top: 54px;
  padding-bottom: 58px;
}

.results-section,
.text-analysis-section {
  margin-bottom: 46px;
  scroll-margin-top: 80px;
}

.grammar-panel {
  padding: 18px;
  color: var(--grammar-ink);
  overflow: hidden;
}

.results-panel {
  padding: 0;
}

.translation-list {
  margin: 0;
  padding-left: 18px;
}

.dictionary-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 6px 0;
}

.dictionary-actions .v-btn {
  width: auto;
}

.grammar-audit-summary {
  display: flex;
  flex-wrap: wrap;
  gap: 28px;
  padding: 8px 8px 22px;
  border-bottom: 1px solid rgba(28, 97, 209, 0.16);
}

.grammar-audit-summary > div {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 7px;
}

.grammar-audit-summary span:first-child {
  color: var(--grammar-muted);
  font-size: 0.76rem;
  font-weight: 800;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.audit-word {
  font-family: Georgia, 'Times New Roman', serif;
  font-size: 1.25rem;
}

.grammar-audit-events {
  margin-top: 20px;
}

.audit-event-card {
  width: min(720px, 100%);
  border-color: rgba(28, 97, 209, 0.2);
  background: rgba(255, 255, 255, 0.5);
}

.audit-event-meta {
  margin-top: 8px;
  color: var(--grammar-muted);
}

.audit-detail-list {
  margin-top: 8px;
  padding-left: 20px;
}

.audit-results {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 16px;
}

.audit-card {
  border-color: rgba(28, 97, 209, 0.2);
  background: rgba(253, 246, 227, 0.92);
  color: var(--grammar-ink);
}

.audit-greek {
  font-family: Georgia, 'Times New Roman', serif;
  font-size: 1.1rem;
  line-height: 1.75;
}

.audit-translation {
  margin-top: 12px;
  color: var(--grammar-muted);
  line-height: 1.55;
}

.interlinear-sentence {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  gap: 10px 8px;
  padding: 10px;
}

.interlinear-token {
  min-width: 72px;
  padding: 10px 12px;
  border: 1px solid rgba(28, 97, 209, 0.22);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.58);
  color: var(--grammar-ink);
  cursor: pointer;
  text-align: center;
  transition: border-color 140ms ease, box-shadow 140ms ease, transform 140ms ease;
}

.interlinear-token:hover,
.interlinear-token.selected {
  border-color: var(--grammar-primary);
  box-shadow: 0 5px 14px rgba(28, 97, 209, 0.16);
  transform: translateY(-2px);
}

.interlinear-token.unresolved {
  border-color: rgba(190, 55, 55, 0.48);
  background: rgba(255, 235, 235, 0.72);
}

.greek-token {
  display: block;
  font-family: Georgia, 'Times New Roman', serif;
  font-size: 1.3rem;
  font-weight: 700;
}

.token-gloss {
  display: block;
  max-width: 180px;
  margin-top: 5px;
  color: var(--grammar-muted);
  font-size: 0.78rem;
  line-height: 1.25;
}

.literal-translation > span {
  color: var(--grammar-muted);
  font-size: 0.78rem;
  font-weight: 800;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.literal-translation p {
  margin: 6px 0 22px;
  font-size: 1.05rem;
  line-height: 1.6;
}

.token-details {
  padding-top: 18px;
  border-top: 1px solid rgba(28, 97, 209, 0.16);
}

.token-detail-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.selected-gloss {
  display: block;
  margin-top: 3px;
  color: var(--grammar-muted);
}

.interpretation-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 12px;
  margin-top: 18px;
}

.interpretation-card {
  border-color: rgba(28, 97, 209, 0.2);
  background: rgba(255, 255, 255, 0.48);
}

.interpretation-root {
  font-family: Georgia, 'Times New Roman', serif;
}

.info-card {
  background: #fefcf5;
  color: var(--grammar-ink);
}

@media (max-width: 900px) {
  .grammar-hero,
  .grammar-hero-shade {
    min-height: auto;
  }

  .grammar-hero-content {
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

  .sentence-actions,
  .token-detail-heading {
    align-items: stretch;
    flex-direction: column;
  }

  .sentence-actions .v-btn {
    width: 100%;
  }

  .grammar-mode-toggle {
    display: flex;
    width: 100%;
  }

  .grammar-mode-toggle :deep(.v-btn) {
    flex: 1 1 auto;
  }
}

</style>
