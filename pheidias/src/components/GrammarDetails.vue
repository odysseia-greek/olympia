<script>
import { useQuery } from "@vue/apollo-composable";
import { CheckGrammar } from "@/constants/graphql";
import { ref, watch } from "vue";

export default {
  name: "GrammarDetails",
  props: {
    clickedWord: {
      type: String,
      required: true,
    },
    forceUpdate: {
      type: Number,
      required: true,
    },
  },
  setup(props) {
    const grammarResults = ref([]);
    const grammarError = ref(null);
    const cycle = ref(false);
    const dialog = ref(false);

    const fetchGrammarDetails = async (word) => {
      grammarError.value = null;
      grammarResults.value = null;
      try {
        const { onResult } = useQuery(CheckGrammar, { word });
        onResult((response) => {
          if (response.data) {
            grammarResults.value = response.data.grammar.results;
          } else {
            grammarResults.value = [{
              word,
              rootWord: word,
              translations: ['No declension found.', 'Unfortunately the grammar component is not complete (yet)', 'Try searching the dictionary with partial or fuzzy results.'],
              rule: 'No rule found',
            }];
          }
        });
      } catch (error) {
        grammarError.value = error;
      }
    };

    const goToDictionaryEntry = (rootWord) => {
      const url = `/dictionary?mode=exact&language=greek&extended=true&word=${encodeURIComponent(rootWord)}`;
      window.open(url, '_blank');
    }

    const goToGrammarEntry = (rootWord) => {
      if (rootWord.includes('–')) {
        rootWord = rootWord.split('–')[0].trim();
      }

      if (rootWord.includes(',')) {
        rootWord = rootWord.split(',')[0].trim();
      }

      const url = `/grammar?word=${encodeURIComponent(rootWord)}`;
      window.open(url, '_blank');
    }

    const onWordClick = (word) => {
      // Define a regular expression to match common punctuation and reading signs
      const punctuationRegex = /[.,;!?(){}[\]"'<>·]/g;

      // Remove the punctuation from the word
      word = word.replace(punctuationRegex, '').trim();
      dialog.value = true;
      fetchGrammarDetails(word);
    };


    watch([() => props.clickedWord, () => props.forceUpdate], ([newWord]) => {
      if (newWord) {
        onWordClick(newWord);
      }
    });

    return {
      grammarResults,
      grammarError,
      cycle,
      dialog,
      goToDictionaryEntry,
      goToGrammarEntry,
    };
  },
}
</script>

<template>
  <v-dialog v-model="dialog" max-width="600px">
    <v-card color="papyrus">
      <v-card-title class="headline">Word Details</v-card-title>
      <v-card-text>
        <div v-if="grammarError">{{ grammarError.message }}</div>
        <div v-if="grammarResults">
          <v-carousel
                :continuous="false"
                :cycle="cycle"
                show-arrows="hover"
                hide-delimiters
            >
              <v-carousel-item v-for="(result, i) in grammarResults" :key="i" color="papyrus">
                <v-sheet
                    height="100%"
                    class="text-center ma-5"
                    color="secondaryPapyrus"
                >
                  <v-list style="background: #fefcf5">
                    <v-list-item-title class="text-h5 ma-1">{{ result.word }}</v-list-item-title>
                    <v-list-item-subtitle class="ma-1">
                      <strong>Translations:</strong>
                      <v-list-item
                          v-for="(translation, j) in result.translations"
                          :key="j"
                          :title="j+1 + '. ' + translation"
                          class="ma-0"
                      ></v-list-item>
                    </v-list-item-subtitle>
                    <v-list-item-subtitle class="ma-1"><strong>Root:</strong> {{ result.rootWord }}</v-list-item-subtitle>
                    <v-list-item-subtitle class="ma-1"><strong>Rule:</strong> {{ result.rule }}</v-list-item-subtitle>
                    <br>
                    <br>
                    <v-divider></v-divider>
                    <v-btn
                        class="ma-5"
                        color="primary"
                        @click="goToDictionaryEntry(result.rootWord)"
                    >
                      Dictionary
                      <v-icon
                          end
                      >
                        mdi-text-box-search
                      </v-icon>
                    </v-btn>
                    <br>
                    <v-btn
                        class="ma-5"
                        color="primary"
                        @click="goToGrammarEntry(result.rootWord)"
                    >
                      Grammar
                      <v-icon
                          end
                      >
                        mdi-book-search
                      </v-icon>
                    </v-btn>
                  </v-list>
                </v-sheet>
              </v-carousel-item>
            </v-carousel>
              <v-list two-line style="background: #fdf6e3">
                <v-list-item>
                  <v-list-item-action>
                    <v-switch
                        color="primary"
                        v-model="cycle"
                        label="Cycle Translations"
                        inset
                    ></v-switch>
                  </v-list-item-action>
            </v-list-item>
          </v-list>
        </div>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<style scoped>
</style>
