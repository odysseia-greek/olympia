<template>
  <div id="grammar">
    <v-app :style="{ background: $vuetify.theme.themes[theme].background }">
      <v-main>
        <v-card color="primary">
          <v-card-text class="d-flex justify-space-between align-center">
            <div>
              <strong>Grammar Declension Tool</strong><br>
              Enter a word in its undeclined form (e.g., λόγων) to see its declensions and root word (e.g., λόγος).<br>
              The root word is then searched within available texts in all the currently known conjugations.
            </div>
            <v-btn @click="infoDialogVisible = true" variant="text" icon="mdi-information"></v-btn>
          </v-card-text>

      <v-dialog v-model="infoDialogVisible" max-width="80%">
        <v-card class="ma-5">
          <v-card-title class="headline">Grammar Conjugation</v-card-title>
          <v-card-text>
            <v-list>
              <v-list-item>
                <v-list-item-title class="subtitle-1">
                  This section provides information about the different components:
                </v-list-item-title>
              </v-list-item>
              <v-divider></v-divider>
              <v-list-item>
                <v-list-item-title>
                  <strong>Search Input:</strong>
                </v-list-item-title>
                <v-list-item>
                  Enter a word in its undeclined form to see its declensions and root word. For example, entering <em>λόγων</em> will show its declension and root word <em>λόγος</em>.
                  <br>
                  The root word is then searched within available texts in all the currently known conjugations.
                  <br><br>
                  <ul>
                    <li>ἔβαλλε</li>
                    <li>φέροντος</li>
                    <li>ἀληθῆ</li>
                    <li>λόγων</li>
                    <li>Ἀθηναῖος</li>
                    <li>πόλεως</li>
                  </ul>
                </v-list-item>
              </v-list-item>
              <v-divider></v-divider>
              <v-list-item>
                <v-list-item-title>
                  <strong>Text Analyze:</strong>
                </v-list-item-title>
                <v-list-item>
                  If any conjugations with translations are found, an attempt is made to find these words in all their forms in the text. You can scroll through these texts and attempt to translate them using the provided link.
                </v-list-item>
              </v-list-item>
              <v-divider></v-divider>
            </v-list>
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" variant="text" @click="infoDialogVisible = false">Close</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>


              <v-card-text>
                <v-autocomplete
                    :loading="loading"
                    v-model="queryWord"
                    @update:search="updateFromHistory"
                    :items="historyOfWords"
                    prepend-icon="mdi-magnify"
                    placeholder="What Greek word are you trying to decline?"
                    variant="underlined"
                    @keyup.enter="grammarWord($event.target.value)"
                    clearable
                ></v-autocomplete>
              </v-card-text>


            <v-expand-transition>
              <v-card light color="background">
                <v-card-text>
                  <h2>Results</h2>
                  <br>
                  <v-data-table
                      dense
                      :headers="headers"
                      :items="grammarResults"
                      :items-per-page="10"
                      item-key="name"
                      class="elevation-1"
                  >
                    <template v-slot:item.translations="{ item }">
                      <ol>
                        <li v-for="(trans, index) in item.translations" :key="index">
                          {{ trans }}
                        </li>
                      </ol>
                    </template>
                  </v-data-table>
                </v-card-text>
              </v-card>
            </v-expand-transition>
          </v-card>
          <AnalyzeResults :analyzeResults="analyzeResults" />
      </v-main>
    </v-app>
  </div>
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
      initializeFromURL();
    });

    return {
      theme,
      stepper,
      grammarResults,
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

    };
  }
};
</script>

<style scoped>
h4 {
  margin-top: 2em;
}

h3 {
  margin-top: 0.5em;
}

a {
  cursor: pointer;
}

* {
  box-sizing: border-box;
}

ul {
  padding-left: 20px;
  list-style-type: disc;
}

/* Style the search field */
form.livesearch input[type="text"] {
  padding: 10px;
  font-size: 17px;
  border: 1px solid grey;
  text-align: center;
  width: 40%;
  background: #f1f1f1;
}

</style>
