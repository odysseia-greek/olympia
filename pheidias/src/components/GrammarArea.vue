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

              <v-autocomplete
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

              <div class="example-row">
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
            </section>
          </v-container>
        </div>
      </header>

      <main class="grammar-content">
        <v-container class="content-container">
          <section class="results-section">
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
              </v-data-table>
            </v-sheet>
          </section>

          <section v-if="analyzeResults.length" class="text-analysis-section">
            <div class="section-heading">
              <div>
                <v-chip color="triadic" variant="tonal">Text analysis</v-chip>
                <h2>Found in available texts</h2>
              </div>
              <p>When root-word text analysis is enabled, matching passages appear here.</p>
            </div>
            <v-sheet class="grammar-panel" color="secondaryPapyrus">
              <AnalyzeResults :analyzeResults="analyzeResults" />
            </v-sheet>
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
                <v-list-item-title><strong>Text Analyze:</strong></v-list-item-title>
                <v-list-item-subtitle>
                  If enabled later, matching root forms can be searched in available texts.
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
import {getCurrentInstance, onMounted, ref} from 'vue';
import {useQuery} from '@vue/apollo-composable';
import {Analyze, CheckGrammar} from '@/constants/graphql';
import AnalyzeResults from "@/components/AnalyzeResults.vue";

export default {
  name: 'GrammarArea',
  components: {
    AnalyzeResults, // Register the new component
  },
  setup() {
    const { proxy } = getCurrentInstance();
    const theme = ref('light');
    const stepper = ref(1);
    const infoDialogVisible = ref(false);
    const grammarHeroImage = ref('');
    const grammarResults = ref([]);
    const analyzeResults = ref([]);
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

    const headers = [
      {title: 'Queried', align: 'start', sortable: true, value: 'word'},
      {title: 'Rule', value: 'rule'},
      {title: 'Root', value: 'rootWord'},
      {title: 'Translation(s)', value: 'translations' },
    ];

    const loadHeroImage = () => {
      import('@/assets/grammar.webp').then((module) => {
        grammarHeroImage.value = module.default;
      });
    };

    const analyzeRootWord = (grammarResults) => {
      analyzeResults.value = [];
      let foundRootWords = new Set();
      let foundTranslations = new Set();
      grammarResults.forEach(result => {
        let parsedrootWord = result.rootWord
        if (parsedrootWord.includes('–')) {
          parsedrootWord = parsedrootWord.split('–')[0].trim();
        }

        if (parsedrootWord.includes(',')) {
          parsedrootWord = parsedrootWord.split(',')[0].trim();
        }

        if (foundRootWords.has(parsedrootWord) || foundTranslations.has(result.translation)) {
          return; // Skip if the root word has already been queried
        }
        foundRootWords.add(parsedrootWord); // Add the root word to the set
        foundTranslations.add(result.translation);
        const { onResult, onError } = useQuery(Analyze, { rootword: parsedrootWord });

        onResult(({ data }) => {
          if (data) {
            analyzeResults.value.push(data.analyze);
          }
        });

        onError((error) => {
          errors.value.push(error);
        });
      });
    };

    const updateFromHistory = async () => {
      if (historyOfWords.value.includes(queryWord.value)) {
        await grammarWord(queryWord.value)
      }
    }

    const grammarWord = async (word) => {
      loading.value = true;
      updateUrl({
        word: word,
      });

      try {
        const { onResult, onError } = useQuery(CheckGrammar, {
          word: word
        });

        if (!historyOfWords.value.includes(word)) {
          historyOfWords.value.push(word);
        }

        onResult((response) => {
          loading.value = false;
          if (response?.data?.grammar?.results) {
            grammarResults.value = response.data.grammar.results.map(item => {
              if (!item.translations || item.translations.length === 0) {
                item.translations = ['No translation found'];
              }

              return item;
            });

            if (response.data.grammar.results.length > 0) {
              // analyzeRootWord(response.data.grammar.results);
            }
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

    const initializeFromURL = async () => {
      const { word } = proxy.$route.query;

      if (word) {
        queryWord.value = word
        await grammarWord(word)
      }
    };

    const updateUrl = (query) => {
      const currentQuery = proxy.$route.query;
      const newQuery = { ...currentQuery, ...query };
      const queryChanged = Object.keys(newQuery).some(key => currentQuery[key] !== newQuery[key]);

      if (queryChanged) {
        proxy.$router.replace({ name: 'Dionysios', query: newQuery });
      }
    };

    const highlightWord = (text) => {
      const regex = /&&&(.*?)&&&/g;
      return text.replace(regex, '<span style="background-color: yellow;">$1</span>');
    };

    onMounted(() => {
      loadHeroImage();
      initializeFromURL();
    });

    return {
      theme,
      stepper,
      grammarResults,
      grammarHeroImage,
      errors,
      headers,
      analyzeResults,
      infoDialogVisible,
      historyOfWords,
      loading,
      queryWord,
      grammarWord,
      analyzeRootWord,
      highlightWord,
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
}

</style>
