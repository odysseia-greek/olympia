<template>
  <v-app id="dictionary" :style="{ background: $vuetify.theme.themes[theme].background }">
    <v-main>
      <header
          ref="dictionaryHeroRef"
          class="dictionary-hero"
          :style="{ backgroundImage: `url(${dictionaryHeroImage})` }"
      >
        <div class="dictionary-hero-shade">
          <v-container class="dictionary-hero-content">
            <section ref="searchPanelRef" class="search-panel" aria-labelledby="dictionary-search-heading">
              <div class="panel-heading">
                <div>
                  <div class="section-label">Dictionary</div>
                  <h1 id="dictionary-search-heading">Alexandros</h1>
                  <p>Search Ancient Greek, English, or Dutch entries with lexical detail, glosses, and text references.</p>
                </div>
                <v-btn
                    aria-label="Dictionary information"
                    color="primary"
                    icon="mdi-information"
                    variant="tonal"
                    @click="infoDialogVisible = true"
                ></v-btn>
              </div>

              <v-autocomplete
                  :loading="loading"
                  :model-value="selectedSearchItem"
                  :search="search"
                  @update:model-value="onSearchSelect"
                  @update:search="onSearchInput"
                  @update:focused="searchInputFocused = $event"
                  :items="searchHistory"
                  hide-no-data
                  hide-selected
                  color="primary"
                  label="Enter a word to search"
                  placeholder="Try λόγος, γράφω, wisdom..."
                  prepend-inner-icon="mdi-magnify"
                  @keyup.enter="commitSearch($event.target.value, { scrollToResults: true })"
                  @click:clear="clearSearch"
                  clearable
              ></v-autocomplete>

              <v-row dense>
                <v-col cols="12" md="5">
                  <div class="control-label">Language</div>
                  <v-radio-group v-model="selectedLanguage" inline density="compact">
                    <v-radio color="secondary" label="Greek" value="greek"></v-radio>
                    <v-radio color="secondary" label="English" value="english"></v-radio>
                    <v-radio color="secondary" label="Nederlands" value="dutch"></v-radio>
                  </v-radio-group>
                </v-col>

                <v-col cols="12" md="7">
                  <div class="control-label">Search mode</div>
                  <v-radio-group v-model="dictionaryMode" inline density="compact">
                    <v-radio color="secondary" label="Partial" value="partial"></v-radio>
                    <v-radio color="secondary" label="Exact" value="exact"></v-radio>
                    <v-radio color="secondary" label="Extended" value="extended"></v-radio>
                    <v-radio color="secondary" label="Fuzzy" value="fuzzy"></v-radio>
                  </v-radio-group>
                </v-col>
              </v-row>

              <v-expand-transition>
                <v-switch
                    v-if="canSearchInText"
                    v-model="extendedMode"
                    label="Search word in available texts"
                    color="secondary"
                    inset
                    hide-details
                ></v-switch>
              </v-expand-transition>
            </section>
          </v-container>
        </div>
      </header>

      <main class="dictionary-content">
        <v-container class="content-container">
          <section ref="resultsSectionRef" class="results-section">
            <div class="section-heading">
              <div>
                <v-chip color="secondary" variant="tonal">Results</v-chip>
                <h2 ref="resultsContainerRef">Dictionary results</h2>
              </div>
              <p>
                Results adapt to the selected language and search mode. Exact Greek search can also inspect available texts.
              </p>
            </div>

            <v-sheet class="dictionary-panel results-panel" color="text">
              <v-data-table
                  dense
                  :headers="headers"
                  :items="searchResults"
                  :items-per-page="10"
                  item-key="id"
              >

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

                  <template v-slot:item.linkedWord="{ item }">
                    <v-btn
                        v-if="item.linkedWord"
                        class="linked-word-action"
                        size="small"
                        color="secondary"
                        variant="tonal"
                        prepend-icon="mdi-arrow-right-circle"
                        @click="searchLinkedWord(item.linkedWord)"
                    >
                      {{ item.linkedWord }}
                    </v-btn>
                    <span v-else class="italic-text">—</span>
                  </template>
              </v-data-table>
            </v-sheet>
          </section>

          <section v-if="richResults.length" ref="featuredRef" class="featured-section">
            <div class="section-heading">
              <div>
                <v-chip color="triadic" variant="tonal">Featured entry</v-chip>
                <h2>Expanded lexical detail</h2>
              </div>
              <v-switch
                  v-model="cycleRich"
                  color="secondary"
                  inset
                  label="Cycle"
                  hide-details
              />
            </div>

            <v-carousel
                class="featured-carousel"
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
                <v-sheet class="featured-entry" rounded="lg" color="secondaryPapyrus">
                  <div class="entry-heading">
                    <div>
                      <div class="entry-headword">{{ r.headword }}</div>
                      <div class="text-subtitle-1">
                        <strong>{{ r.partOfSpeech }}</strong>
                        <span class="ml-2 italic-text" v-if="r.normalized">({{ r.normalized }})</span>
                      </div>
                    </div>
                    <div class="entry-tags">
                      <v-chip v-if="r.noun" class="ma-1" color="primary" variant="flat">
                        {{ r.noun.declension }} decl · {{ r.noun.genitive }}
                      </v-chip>
                      <v-chip v-if="r.verb" class="ma-1" color="primary" variant="flat">
                        verb
                      </v-chip>
                    </div>
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
                    <v-btn
                        class="linked-word-action"
                        color="secondary"
                        variant="tonal"
                        prepend-icon="mdi-link-variant"
                        @click="searchLinkedWord(r.linkedWord)"
                    >
                      Search linked: {{ r.linkedWord }}
                    </v-btn>
                  </div>
                </v-sheet>
              </v-carousel-item>
            </v-carousel>
          </section>

          <section
              v-if="extendedMode && selectedLanguage.toLowerCase() === 'greek' && dictionaryMode.toLowerCase() === 'exact'"
              class="text-search-section"
          >
            <div class="section-heading">
              <div>
                <v-chip color="primary" variant="tonal">Text search</v-chip>
                <h2>Found in available texts</h2>
              </div>
              <p>
                Exact Greek searches can surface passages that contain matching forms, keeping text evidence aligned with the dictionary result.
              </p>
            </div>
            <v-sheet class="dictionary-panel text-search-panel" color="secondaryPapyrus">
              <AnalyzeResults :analyzeResults="analyzeResults" />
            </v-sheet>
          </section>

          <section ref="popularRef" class="popular-section">
            <div class="section-heading">
              <div>
                <v-chip color="primary" variant="tonal">Recent searches</v-chip>
                <h2>Words learners are checking</h2>
              </div>
              <p>
                The list updates as searches are performed, giving you a quick route back into common entries.
              </p>
            </div>
            <v-sheet class="dictionary-panel top-search-panel" color="secondaryPapyrus">
              <DictionaryTopFive :refresh-token="topFiveRefreshToken" />
            </v-sheet>
          </section>
        </v-container>
      </main>

      <v-dialog v-model="infoDialogVisible" max-width="860">
        <v-card class="info-card">
          <v-card-title class="headline">Dictionary</v-card-title>
          <v-card-text>
            <v-list>
              <v-list-item>
                <v-list-item-title class="subtitle-1">
                  This section provides information about the different controls.
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
                  Partial matches fragments, Exact matches a specific word, Extended searches phrases, and Fuzzy tolerates typos.
                </v-list-item-subtitle>
              </v-list-item>
              <v-divider></v-divider>

              <v-list-item>
                <v-list-item-title><strong>Search Input:</strong></v-list-item-title>
                <v-list-item-subtitle>
                  Enter the word you are looking for. The search happens as you type.
                </v-list-item-subtitle>
              </v-list-item>
              <v-divider></v-divider>

              <v-list-item>
                <v-list-item-title><strong>Extended Search:</strong></v-list-item-title>
                <v-list-item-subtitle>
                  Exact + Greek can additionally search in available texts.
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" @click="infoDialogVisible = false">Close</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-main>
  </v-app>
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
  const debounced = (...args) => {
    if (t) clearTimeout(t);
    t = setTimeout(() => fn(...args), waitMs);
  };
  debounced.cancel = () => {
    if (t) clearTimeout(t);
    t = null;
  };
  return debounced;
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
    const dictionaryHeroImage = ref('');
    const dictionaryHeroRef = ref();
    const searchPanelRef = ref();
    const popularRef = ref();
    const resultsSectionRef = ref();
    const featuredRef = ref();

    const search = ref('');
    const selectedSearchItem = ref(null);
    const searchInputFocused = ref(false);
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
    const activeRequestId = ref(0);
    const suppressAutoSearch = ref(false);

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
        if (resultsSectionRef.value) {
          resultsSectionRef.value.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
      });
    }

    function scrollMeTo(refName) {
      nextTick(() => {
        if (refName === 'dictionaryHeroRef' && dictionaryHeroRef.value) {
          dictionaryHeroRef.value.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
        if (refName === 'searchPanelRef' && searchPanelRef.value) {
          searchPanelRef.value.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
        if (refName === 'popularRef' && popularRef.value) {
          popularRef.value.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
        if (refName === 'resultsSectionRef' && resultsSectionRef.value) {
          resultsSectionRef.value.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
        if (refName === 'featuredRef' && featuredRef.value) {
          featuredRef.value.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
      });
    }

    function loadHeroImage() {
      import('@/assets/alexander.webp').then((module) => {
        dictionaryHeroImage.value = module.default;
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
      if (!value) {
        loading.value = false;
        return;
      }

      if (!canSearchInText.value) extendedMode.value = false;

      const requestId = activeRequestId.value + 1;
      activeRequestId.value = requestId;
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
        if (requestId === activeRequestId.value) {
          loading.value = false;
        }
      }
    }

    const debouncedFetch = debounce(fetchDictionary, 500);

    function commitSearch(value, options = {}) {
      const v = (value || '').trim();
      if (!v) {
        clearSearch();
        return;
      }

      if (!searchHistory.value.includes(v)) searchHistory.value.push(v);
      search.value = v;
      selectedSearchItem.value = v;

      updateUrl({
        mode: dictionaryMode.value,
        language: selectedLanguage.value,
        extended: extendedMode.value,
        word: v,
      });

      debouncedFetch(v);
      if (options.scrollToResults) scrollToResults();
    }

    function onSearchInput(value) {
      const v = (value || '').trim();
      if (!v) {
        if (!searchInputFocused.value) return;

        clearSearch();
        return;
      }

      search.value = v;
      selectedSearchItem.value = null;

      updateUrl({
        mode: dictionaryMode.value,
        language: selectedLanguage.value,
        extended: extendedMode.value,
        word: v,
      });

      debouncedFetch(v);
    }

    function onSearchSelect(value) {
      if (!value) return;

      commitSearch(value);
    }

    async function searchLinkedWord(value) {
      const v = (value || '').trim();
      if (!v) return;

      suppressAutoSearch.value = true;
      selectedLanguage.value = 'greek';
      dictionaryMode.value = 'exact';
      extendedMode.value = false;

      await nextTick();
      suppressAutoSearch.value = false;
      commitSearch(v, { scrollToResults: true });
    }

    function clearSearch() {
      debouncedFetch.cancel();
      activeRequestId.value += 1;
      loading.value = false;
      search.value = '';
      selectedSearchItem.value = null;
      rawResults.value = [];
      searchResults.value = [];
      analyzeResults.value = [];
    }

    watch(dictionaryMode, () => {
      if (suppressAutoSearch.value) return;

      if (!canSearchInText.value) extendedMode.value = false;
      if (search.value) commitSearch(search.value);
    });

    watch(selectedLanguage, () => {
      if (suppressAutoSearch.value) return;

      if (!canSearchInText.value) extendedMode.value = false;
      if (search.value) commitSearch(search.value);
    });

    watch(extendedMode, () => {
      if (suppressAutoSearch.value) return;

      if (search.value) commitSearch(search.value);
    });

    async function initializeFromURL() {
      const { language, mode, word, extended } = proxy.$route.query;

      if (language) selectedLanguage.value = language;
      if (mode) dictionaryMode.value = mode;
      if (extended) extendedMode.value = String(extended).toLowerCase() === 'true';

      if (word) {
        search.value = word;
        selectedSearchItem.value = word;
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
      loadHeroImage();
      initializeFromURL();
    });

    return {
      theme,
      selectedLanguage,
      dictionaryMode,
      extendedMode,
      dictionaryHeroImage,
      dictionaryHeroRef,
      searchPanelRef,
      popularRef,
      resultsSectionRef,
      featuredRef,
      canSearchInText,
      search,
      selectedSearchItem,
      searchInputFocused,
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
      onSearchSelect,
      searchLinkedWord,
      clearSearch,
      commitSearch,
      scrollMeTo,
    };
  },
};
</script>

<style scoped>
#dictionary {
  --dictionary-primary: #1c61d1;
  --dictionary-secondary: #1cd18c;
  --dictionary-triadic: #1cbcd1;
  --dictionary-ink: #20334f;
  --dictionary-muted: #536987;
  color: var(--dictionary-ink);
}

* {
  box-sizing: border-box;
}

a {
  cursor: pointer;
}

.italic-text {
  font-style: italic;
}

.dictionary-hero {
  min-height: 470px;
  background-position: center 34%;
  background-size: cover;
}

.dictionary-hero-shade {
  min-height: 470px;
  display: flex;
  align-items: center;
  background:
      linear-gradient(90deg, rgba(22, 20, 15, 0.66) 0%, rgba(38, 31, 21, 0.42) 48%, rgba(26, 22, 16, 0.14) 100%),
      linear-gradient(180deg, rgba(0, 0, 0, 0.04), rgba(26, 21, 14, 0.32));
}

.dictionary-hero-content {
  padding-top: 56px;
  padding-bottom: 40px;
}

.search-panel {
  max-width: 980px;
  padding: 22px;
  border: 1px solid rgba(28, 188, 209, 0.32);
  border-radius: 8px;
  background: rgba(253, 246, 227, 0.94);
  box-shadow: 0 18px 48px rgba(11, 39, 85, 0.3);
  backdrop-filter: blur(8px);
}

.panel-heading,
.section-heading,
.entry-heading {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 28px;
  margin-bottom: 22px;
}

.section-heading p {
  max-width: 520px;
  margin: 0;
  color: #344765;
  line-height: 1.65;
}

.section-label,
.control-label {
  color: #64789e;
  font-size: 0.82rem;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.control-label {
  margin-top: 4px;
  margin-bottom: 4px;
}

.panel-heading h1,
.section-heading h2 {
  margin: 8px 0 0;
  font-size: clamp(1.6rem, 3vw, 2.35rem);
  line-height: 1.15;
  letter-spacing: 0;
}

.panel-heading p {
  max-width: 680px;
  margin: 10px 0 0;
  color: #344765;
  line-height: 1.55;
}

.dictionary-content {
  background:
      linear-gradient(150deg, rgba(28, 97, 209, 0.16) 0%, rgba(28, 188, 209, 0.12) 34%, rgba(28, 209, 140, 0.1) 62%, rgba(254, 252, 245, 0.98) 100%),
      linear-gradient(180deg, #d5eff7 0%, #f2fbf7 46%, #fefcf5 100%);
}

.content-container {
  max-width: 1240px;
  padding-top: 54px;
  padding-bottom: 58px;
}

.dictionary-hero,
.search-panel,
.popular-section,
.results-section,
.featured-section {
  scroll-margin-top: 80px;
}

.popular-section,
.results-section,
.featured-section,
.text-search-section {
  margin-bottom: 46px;
}

.dictionary-panel,
.featured-entry {
  border: 1px solid rgba(28, 97, 209, 0.16);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 8px 24px rgba(28, 97, 209, 0.1);
}

.dictionary-panel {
  padding: 18px;
  overflow: hidden;
}

.results-panel {
  padding: 0;
}

.featured-carousel {
  border-radius: 8px;
}

.featured-entry {
  min-height: 420px;
  padding: 24px;
}

.entry-heading {
  align-items: start;
}

.entry-headword {
  color: var(--dictionary-ink);
  font-size: clamp(2rem, 4vw, 3rem);
  line-height: 1;
}

.entry-tags {
  display: flex;
  flex-wrap: wrap;
  justify-content: end;
}

.linked-word-action {
  text-transform: none;
}

.text-search-panel {
  padding: 20px;
}

.top-search-panel {
  color: var(--dictionary-ink);
}

@media (max-width: 900px) {
  .dictionary-hero,
  .dictionary-hero-shade {
    min-height: auto;
  }

  .dictionary-hero-content {
    padding-top: 34px;
    padding-bottom: 28px;
  }

  .panel-heading,
  .section-heading,
  .entry-heading {
    display: block;
  }

  .section-heading p,
  .entry-tags {
    margin-top: 10px;
  }

  .entry-tags {
    justify-content: start;
  }
}
</style>
