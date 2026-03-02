<template>
  <div id="dictionary">
    <v-app :style="{ background: $vuetify.theme.themes[theme].background }">
      <v-main>
        <v-card color="primary">
          <v-card-text>
            Dictionary provides words in Ancient Greek, English and Dutch.
            <v-btn @click="infoDialogVisible = true" variant="text" icon="mdi-information"></v-btn>
          </v-card-text>

          <v-dialog v-model="infoDialogVisible" max-width="80%">
            <v-card class="ma-5">
              <v-card-title class="headline">Dictionary</v-card-title>
              <v-card-text>
                <v-list>
                  <v-list-item>
                    <v-list-item-title class="subtitle-1">
                      This section provides information about the different components:
                    </v-list-item-title>
                  </v-list-item>
                  <v-divider></v-divider>

                  <v-list-item>
                    <v-list-item-title><strong>Selected Language:</strong></v-list-item-title>
                    <v-list-item-subtitle>Allows you to choose the language.</v-list-item-subtitle>
                  </v-list-item>
                  <v-divider></v-divider>

                  <v-list-item>
                    <v-list-item-title><strong>Search Mode:</strong></v-list-item-title>
                    <v-list-item-subtitle>
                      The mode you want to use for searching:
                      <v-list>
                        <v-list-item><strong>Partial:</strong> "ouse" matches "house," "mouse," and "trousers".</v-list-item>
                        <v-list-item><strong>Exact:</strong> "house" matches only "house".</v-list-item>
                        <v-list-item><strong>Extended:</strong> matches phrases/expressions (now uses <code>phrase</code> query).</v-list-item>
                        <v-list-item><strong>Fuzzy:</strong> "hiuse" matches "house" based on Levenshtein distance.</v-list-item>
                      </v-list>
                    </v-list-item-subtitle>
                  </v-list-item>
                  <v-divider></v-divider>

                  <v-list-item>
                    <v-list-item-title><strong>Search Input:</strong></v-list-item-title>
                    <v-list-item-subtitle>
                      Enter the word you are looking for. The search happens as you type.
                      <v-list-item><strong>Examples:</strong> όφο, δοτος, Ἀθῆ.</v-list-item>
                    </v-list-item-subtitle>
                  </v-list-item>
                  <v-divider></v-divider>

                  <v-list-item>
                    <v-list-item-title><strong>Results Table:</strong></v-list-item-title>
                    <v-list-item-subtitle>
                      Displays the search results plus extra lexical detail (POS, glosses, definitions).
                    </v-list-item-subtitle>
                  </v-list-item>

                  <v-divider></v-divider>

                  <v-list-item>
                    <v-list-item-title><strong>Extended Search:</strong></v-list-item-title>
                    <v-list-item-subtitle>
                      Exact + Greek can additionally search in texts (foundInText).
                    </v-list-item-subtitle>
                  </v-list-item>

                  <v-divider></v-divider>

                  <v-list-item>
                    <v-list-item-title><strong>Please mind:</strong></v-list-item-title>
                    <v-list-item-subtitle>
                      Search-in-text works only with language set to Greek and Exact mode.
                    </v-list-item-subtitle>
                  </v-list-item>
                </v-list>
              </v-card-text>
              <v-card-actions>
                <v-btn color="primary" @click="infoDialogVisible = false">Close</v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>

          <h3 class="mx-4">Language</h3>
          <v-radio-group v-model="selectedLanguage" class="mx-4">
            <v-radio color="secondary" label="Greek (default)" value="greek"></v-radio>
            <v-radio color="secondary" label="English" value="english"></v-radio>
            <v-radio color="secondary" label="Nederlands" value="dutch"></v-radio>
          </v-radio-group>

          <h3 class="mx-4">Search Mode</h3>
          <v-radio-group v-model="dictionaryMode" class="mx-4">
            <v-radio color="secondary" label="Partial (default)" value="partial"></v-radio>
            <v-radio color="secondary" label="Exact" value="exact"></v-radio>
            <v-radio color="secondary" label="Extended" value="extended"></v-radio>
            <v-radio color="secondary" label="Fuzzy" value="fuzzy"></v-radio>
          </v-radio-group>

          <h3 class="mx-4" v-if="canSearchInText">Extended Search</h3>
          <v-switch
              class="mx-4"
              v-if="canSearchInText"
              v-model="extendedMode"
              label="Search word in available texts"
              color="secondary"
          ></v-switch>

          <v-card-text>
            <v-autocomplete
                :loading="loading"
                v-model="search"
                @update:search="onSearchInput"
                :items="searchHistory"
                hide-no-data
                color="white"
                label="Enter a word to search"
                placeholder="Start typing..."
                prepend-icon="mdi-magnify"
                @keyup.enter="commitSearch($event.target.value)"
                clearable
            ></v-autocomplete>
          </v-card-text>

          <v-divider></v-divider>

          <DictionaryTopFive :refresh-token="topFiveRefreshToken" />

          <v-expand-transition>
            <v-card light color="background">
              <v-card-text>
                <v-data-table
                    dense
                    :headers="headers"
                    :items="searchResults"
                    :items-per-page="10"
                    item-key="id"
                    class="elevation-1"
                >
                  <template v-slot:top>
                    <v-toolbar flat>
                      <h2 ref="resultsContainerRef" class="mx-4">Dictionary Results</h2>
                      <v-spacer></v-spacer>
                    </v-toolbar>
                  </template>

                  <!-- nicer rendering for multi-line / arrays -->
                  <template v-slot:item.quickGlosses="{ item }">
                    <div v-if="item.quickGlosses?.length">
                      <div v-for="(g, i) in item.quickGlosses" :key="i">
                        <strong>{{ g.language }}:</strong> {{ g.gloss }}
                      </div>
                    </div>
                    <span v-else class="italic-text">—</span>
                  </template>

                  <template v-slot:item.principalParts="{ item }">
                    <div v-if="item.principalParts?.length">
                      {{ item.principalParts.join(' · ') }}
                    </div>
                    <span v-else class="italic-text">—</span>
                  </template>

                  <template v-slot:item.definitions="{ item }">
                    <div v-if="item.definitionsText">
                      {{ item.definitionsText }}
                    </div>
                    <span v-else class="italic-text">—</span>
                  </template>

                  <template v-slot:item.modernConnections="{ item }">
                    <div v-if="item.modernConnections?.length">
                      <div v-for="(mc, i) in item.modernConnections" :key="i">
                        <strong>{{ mc.term }}</strong><span v-if="mc.note"> — {{ mc.note }}</span>
                      </div>
                    </div>
                    <span v-else class="italic-text">—</span>
                  </template>

                  <template v-slot:item.nounInfo="{ item }">
                    <div v-if="item.nounInfo">
                      {{ item.nounInfo }}
                    </div>
                    <span v-else class="italic-text">—</span>
                  </template>
                </v-data-table>
                <!-- Rich Results -->
                <div v-if="richResults.length" class="ma-10 d-flex justify-center">
                  <div style="width: 100%; max-width: 980px;">
                    <div class="d-flex align-center">
                    <h3 class="mr-4">Featured entry</h3>
                    <v-spacer />
                    <v-switch
                        v-model="cycleRich"
                        color="secondary"
                        inset
                        label="Cycle"
                        class="ml-4"
                    />
                  </div>

                  <v-carousel
                      :continuous="false"
                      :cycle="cycleRich"
                      :show-arrows="richResults.length > 1 ? 'hover' : false"
                      hide-delimiters
                      height="100%"
                  >
                    <v-carousel-item
                        v-for="(r, i) in richResults"
                        :key="`${r.headword}-${i}`"

                    >

                      <v-sheet class="pa-4" rounded="lg" color="secondaryPapyrus">
                        <div class="d-flex align-center">
                          <div>
                            <div class="text-h4">{{ r.headword }}</div>
                            <div class="text-subtitle-1">
                              <strong>{{ r.partOfSpeech }}</strong>
                              <span class="ml-2 italic-text" v-if="r.normalized">({{ r.normalized }})</span>
                            </div>
                          </div>
                          <v-spacer />
                          <v-chip v-if="r.noun" class="ma-1" color="primary" variant="flat">
                            {{ r.noun.declension }} decl · {{ r.noun.genitive }}
                          </v-chip>
                          <v-chip v-if="r.verb" class="ma-1" color="primary" variant="flat">
                            verb
                          </v-chip>
                        </div>

                        <v-divider class="my-3" />

                        <!-- Glosses -->
                        <div v-if="r.quickGlosses?.length" class="mb-3">
                          <div class="text-subtitle-2 mb-1"><strong>Quick glosses</strong></div>
                          <div class="d-flex flex-wrap">
                            <v-chip
                                v-for="(g, j) in r.quickGlosses"
                                :key="j"
                                class="ma-1"
                                color="primary"
                                variant="outlined"
                            >
                              {{ g.language }}: {{ g.gloss }}
                            </v-chip>
                          </div>
                        </div>

                        <!-- Principal parts -->
                        <div v-if="r.verb?.principalParts?.length" class="mb-3">
                          <div class="text-subtitle-2 mb-1"><strong>Principal parts</strong></div>
                          <div class="text-body-1">
                            {{ r.verb.principalParts.join(' · ') }}
                          </div>
                        </div>

                        <!-- Definitions (group by grade) -->
                        <div v-if="r.definitions?.length" class="mb-3">
                          <div class="text-subtitle-2 mb-1"><strong>Definitions</strong></div>
                          <v-list density="compact" style="background: transparent">
                            <v-list-item
                                v-for="(d, di) in r.definitions"
                                :key="di"
                                class="px-0"
                            >
                              <v-list-item-title>
                                <strong>Grade {{ d.grade }}</strong>
                              </v-list-item-title>
                              <v-list-item-subtitle>
                                <div v-for="(m, mi) in d.meanings" :key="mi">
                                  <strong>{{ m.language }}:</strong> {{ m.definition }}
                                </div>
                              </v-list-item-subtitle>
                            </v-list-item>
                          </v-list>
                        </div>

                        <!-- Modern connections -->
                        <div v-if="r.modernConnections?.length" class="mb-1">
                          <div class="text-subtitle-2 mb-1"><strong>Modern connections</strong></div>
                          <v-list density="compact" style="background: transparent">
                            <v-list-item
                                v-for="(c, ci) in r.modernConnections"
                                :key="ci"
                                class="px-0"
                            >
                              <v-list-item-title>
                                <strong>{{ c.term }}</strong>
                              </v-list-item-title>
                              <v-list-item-subtitle v-if="c.note">{{ c.note }}</v-list-item-subtitle>
                            </v-list-item>
                          </v-list>
                        </div>

                        <!-- Linked word -->
                        <div v-if="r.linkedWord" class="mt-2">
                          <v-chip color="secondary" variant="tonal">
                            linked: {{ r.linkedWord }}
                          </v-chip>
                        </div>
                      </v-sheet>
                    </v-carousel-item>
                  </v-carousel>
                </div>
                </div>

                <AnalyzeResults
                    v-if="extendedMode && selectedLanguage.toLowerCase() === 'greek' && dictionaryMode.toLowerCase() === 'exact'"
                    :analyzeResults="analyzeResults"
                />
              </v-card-text>
            </v-card>
          </v-expand-transition>
        </v-card>
      </v-main>
    </v-app>
  </div>
</template>

<script>
import { ref, computed, watch, onMounted, getCurrentInstance, nextTick } from 'vue';
import { useApolloClient } from '@vue/apollo-composable';

import AnalyzeResults from '../components/AnalyzeResults.vue';
import DictionaryTopFive from '../components/DictionaryTopFive.vue';

import {
  DictionaryExact,
  DictionaryPartial,
  DictionaryFuzzy,
  DictionaryPhrase,
} from '../constants/dictionaryGraphql';

function debounce(fn, waitMs) {
  let t = null;
  return (...args) => {
    if (t) clearTimeout(t);
    t = setTimeout(() => fn(...args), waitMs);
  };
}

function formatDefinitions(definitions, preferredLang /* 'greek'|'english'|'dutch' */) {
  if (!definitions?.length) return '';
  const langMap = { greek: 'gr', english: 'en', dutch: 'nl' };
  const want = langMap[preferredLang] || 'en';

  const sorted = [...definitions].sort((a, b) => (b.grade ?? 0) - (a.grade ?? 0));
  const parts = [];

  for (const def of sorted) {
    const meanings = def.meanings || [];
    const preferred = meanings.find(m => m.language === want) || meanings[0];
    if (preferred?.definition) parts.push(preferred.definition);
    if (parts.length >= 2) break;
  }

  return parts.join(' / ');
}

function pickBestTextForLang(r, lang /* 'en'|'nl' */) {
  // 1) quick gloss
  const g = r.quickGlosses?.find(x => x.language === lang)?.gloss;
  if (g) return g;

  // 2) first definition meaning in that language (highest grade first)
  const defs = r.definitions || [];
  const sorted = [...defs].sort((a, b) => (b.grade ?? 0) - (a.grade ?? 0));
  for (const d of sorted) {
    const m = d.meanings?.find(x => x.language === lang);
    if (m?.definition) return m.definition;
  }

  return '';
}

function isRichRaw(r) {
  return Boolean(r.partOfSpeech);
}

export default {
  name: 'DictionaryArea',
  components: { AnalyzeResults, DictionaryTopFive },
  setup() {
    const { proxy } = getCurrentInstance();
    const { client } = useApolloClient();

    const theme = ref('light');
    const selectedLanguage = ref('greek');
    const dictionaryMode = ref('partial');
    const extendedMode = ref(false);

    const search = ref('');
    const searchHistory = ref([
      'Λακεδαιμονιος',
      'λόγος',
      'ποταμός',
      'Ἀθηναῖος',
      'ναυτικός',
      'ἀγάπη',
      'εἰρήνη',
      'σοφία',
      'γίγνομαι',
      'καί',
      'λέγω',
      'γράφω',
      'ποιέω',
    ]);

    const loading = ref(false);

    const rawResults = ref([]);

    const searchResults = ref([]);

    const analyzeResults = ref([]);
    const infoDialogVisible = ref(false);
    const resultsContainerRef = ref();
    const topFiveRefreshToken = ref(0);

    const canSearchInText = computed(() =>
        selectedLanguage.value.toLowerCase() === 'greek' &&
        dictionaryMode.value.toLowerCase() === 'exact'
    );

    const headers = computed(() => {
      const base = [];

      if (selectedLanguage.value === 'greek') {
        base.push({ title: 'Greek', value: 'greek' });
      } else if (selectedLanguage.value === 'english') {
        base.push({ title: 'English', value: 'english' });
        base.push({ title: 'Greek', value: 'greek' });
      } else {
        base.push({ title: 'Nederlands', value: 'dutch' });
        base.push({ title: 'Grieks', value: 'greek' });
      }

      base.push({ title: 'Glosses', value: 'quickGlosses' });
      base.push({ title: 'Linked', value: 'linkedWord' });
      base.push({ title: 'Normalized', value: 'normalized' });

      return base;
    });

    const cycleRich = ref(false);

    const richResults = computed(() => {
      return (rawResults.value || []).filter(isRichRaw).slice(0, 10);
    });

    function scrollToResults() {
      nextTick(() => {
        if (resultsContainerRef.value) {
          resultsContainerRef.value.scrollIntoView({ behavior: 'smooth' });
        }
      });
    }

    function updateUrl(query) {
      const currentQuery = proxy.$route.query;
      const newQuery = { ...currentQuery, ...query };
      const queryChanged = Object.keys(newQuery).some((key) => currentQuery[key] !== newQuery[key]);
      if (queryChanged) {
        proxy.$router.replace({ name: 'Alexandros', query: newQuery });
      }
    }

    function normalizeMode(mode) {
      const m = (mode || '').toLowerCase();
      if (m === 'extended') return 'phrase';
      if (m === 'partial' || m === 'exact' || m === 'fuzzy' || m === 'phrase') return m;
      return 'partial';
    }

    function pickQuery(mode) {
      switch (normalizeMode(mode)) {
        case 'exact':
          return DictionaryExact;
        case 'fuzzy':
          return DictionaryFuzzy;
        case 'phrase':
          return DictionaryPhrase;
        case 'partial':
        default:
          return DictionaryPartial;
      }
    }

    async function fetchDictionary(word) {
      const value = (word || '').trim();
      if (!value) return;

      if (!canSearchInText.value) extendedMode.value = false;

      loading.value = true;
      analyzeResults.value = [];
      rawResults.value = [];
      searchResults.value = [];

      try {
        const mode = normalizeMode(dictionaryMode.value);
        const query = pickQuery(mode);

        const languageEnum = toGraphqlLanguageEnum(selectedLanguage.value);

        const input =
            mode === 'exact'
                ? { word: value, expand: true, size: 10, language: languageEnum }
                : { word: value, size: 10, language: languageEnum };

        const { data } = await client.query({
          query,
          variables: { input },
          fetchPolicy: 'no-cache',
        });

        const payload = data?.[mode];
        const results = payload?.results || [];

        rawResults.value = results;
        searchResults.value = results.map((r, index) => {
          const principalParts = r.verb?.principalParts || [];
          const nounInfo =
              r.noun?.declension || r.noun?.genitive
                  ? [r.noun?.declension, r.noun?.genitive].filter(Boolean).join(' ')
                  : '';

          return {
            id: `${r.headword || r.normalized || value}-${index}`,
            greek: r.headword,
            english: pickBestTextForLang(r, 'en'),
            dutch: pickBestTextForLang(r, 'nl'),

            headword: r.headword,
            normalized: r.normalized,
            partOfSpeech: r.partOfSpeech,
            quickGlosses: r.quickGlosses || [],
            principalParts,
            nounInfo,
            definitionsText: formatDefinitions(r.definitions, selectedLanguage.value),
            modernConnections: r.modernConnections || [],
            linkedWord: r.linkedWord || '',
          };
        });

        // foundInText only on exact
        if (
            mode === 'exact' &&
            extendedMode.value &&
            selectedLanguage.value.toLowerCase() === 'greek'
        ) {
          const fit = payload?.foundInText;
          if (fit?.texts?.length || fit?.conjugations?.length) {
            analyzeResults.value = [
              {
                rootword: fit.rootword || value,
                conjugations: fit.conjugations || [],
                results: (fit.texts || []).map((t) => ({
                  author: t.author,
                  book: t.book,
                  text: t.text,
                  reference: t.reference,
                  referenceLink: t.referenceLink,
                })),
              },
            ];
          }
        }

        topFiveRefreshToken.value += 1;
      } catch (e) {
        console.log(e);
      } finally {
        setTimeout(() => {
          loading.value = false;
        }, 300);
      }
    }

    const debouncedFetch = debounce(fetchDictionary, 250);

    function commitSearch(value) {
      const v = (value || '').trim();
      if (!v) return;

      if (!searchHistory.value.includes(v)) searchHistory.value.push(v);
      search.value = v;

      updateUrl({
        mode: dictionaryMode.value,
        language: selectedLanguage.value,
        extended: extendedMode.value,
        word: v,
      });

      debouncedFetch(v);
      scrollToResults();
    }

    function onSearchInput(value) {
      const v = (value || '').trim();
      if (!v) return;

      search.value = v;

      updateUrl({
        mode: dictionaryMode.value,
        language: selectedLanguage.value,
        extended: extendedMode.value,
        word: v,
      });

      debouncedFetch(v);
    }

    watch(dictionaryMode, () => {
      if (!canSearchInText.value) extendedMode.value = false;
      if (search.value) commitSearch(search.value);
    });

    watch(selectedLanguage, () => {
      if (!canSearchInText.value) extendedMode.value = false;
      if (search.value) commitSearch(search.value);
    });

    watch(extendedMode, () => {
      if (search.value) commitSearch(search.value);
    });

    async function initializeFromURL() {
      const { language, mode, word, extended } = proxy.$route.query;

      if (language) selectedLanguage.value = language;
      if (mode) dictionaryMode.value = mode;
      if (extended) extendedMode.value = String(extended).toLowerCase() === 'true';

      if (word) {
        search.value = word;
        commitSearch(word);
      }
    }

    function toGraphqlLanguageEnum(uiLang) {
      switch ((uiLang || '').toLowerCase()) {
        case 'greek':
          return 'LANG_GREEK';
        case 'english':
          return 'LANG_ENGLISH';
        case 'dutch':
          return 'LANG_DUTCH';
        default:
          return 'LANGUAGE_UNSPECIFIED';
      }
    }

    onMounted(() => {
      initializeFromURL();
    });

    return {
      theme,
      selectedLanguage,
      dictionaryMode,
      extendedMode,
      canSearchInText,
      search,
      searchHistory,
      loading,
      rawResults,     // optional (debug)
      searchResults,
      richResults,
      cycleRich,
      analyzeResults,
      infoDialogVisible,
      headers,
      resultsContainerRef,
      topFiveRefreshToken,
      onSearchInput,
      commitSearch,
    };
  },
};
</script>

  <style scoped>
h4 { margin-top: 2em; }
h3 { margin-top: 0.5em; }
a { cursor: pointer; }
* { box-sizing: border-box; }
.italic-text { font-style: italic; }
</style>
