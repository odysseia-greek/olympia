<template>
  <v-card-text class="custom-card-text">
    <v-card-subtitle>Language</v-card-subtitle>
    <v-radio-group v-model="selectedLanguage"class="mx-4">
      <v-radio color="secondary" label="Greek" value="greek"></v-radio>
      <v-radio color="secondary" label="English" value="english"></v-radio>
    </v-radio-group>
    <v-autocomplete
        width="100%"
        :loading="loading"
        v-model="selectedWord"
        :items="words"
        hide-no-data
        item-text="greek"
        item-value="translation"
        label="Pick a word from the dropdown"
        prepend-icon="mdi-magnify"
        auto-select-first
    ></v-autocomplete>
    <div v-if="selectedWord" class="translation">
      <p><strong>Translation:</strong> {{ selectedWord }}</p>
      <v-divider></v-divider>
      <p>More like this? Time to search the dictionary! <a href="/dictionary?mode=exact&language=greek&extended=true&word=λόγος">link</a></p>
    </div>
  </v-card-text>
</template>

<script>
import {ref, watch} from 'vue';

export default {
  name: 'SearchComponent',
  setup() {
    const loading = ref(false);
    const selectedWord = ref(null);
    const selectedLanguage = ref('greek');
    const words = ref([
      { title: 'λόγος', translation: 'word, reason' },
      { title: 'ἀγάπη', translation: 'love' },
      { title: 'ψυχή', translation: 'soul' },
      { title: 'ἄνθρωπος', translation: 'man, human' },
      { title: 'θεός', translation: 'god' },
    ]);

    watch(selectedLanguage, (newLanguage) => {
      if (newLanguage === 'greek') {
        words.value = [
          {title: 'λόγος', translation: 'word, reason'},
          {title: 'ἀγάπη', translation: 'love'},
          {title: 'ψυχή', translation: 'soul'},
          {title: 'ἄνθρωπος', translation: 'man, human'},
          {title: 'θεός', translation: 'god'},
        ];
      } else if (newLanguage === 'english') {
        words.value = [
          {title: 'word, reason', translation: 'λόγος'},
          {title: 'love', translation: 'ἀγάπη'},
          {title: 'soul', translation: 'ψυχή'},
          {title: 'man, human', translation: 'ἄνθρωπος'},
          {title: 'god', translation: 'θεός'},
        ];
      }
    });

    return {
      loading,
      selectedWord,
      selectedLanguage,
      words,
    };
  },
};
</script>

<style>
.translation {
  margin-top: 1em;
}

</style>
