<script>
import { defineComponent, toRefs } from 'vue';

export default defineComponent({
  name: 'AnalyzeResults',
  props: {
    analyzeResults: {
      type: Array,
      required: true,
    },
  },
  setup(props) {
    const { analyzeResults } = toRefs(props);

    const highlightWord = (text) => {
      const regex = /&&&(.*?)&&&/g;
      return text.replace(regex, '<span style="background-color: yellow;">$1</span>');
    };

    return {
      analyzeResults,
      highlightWord,
    };
  },
});
</script>

<template>
  <v-row v-if="analyzeResults && analyzeResults.length" style="margin-top: 10px">
    <v-col>
      <v-card class="mx-auto paper-card" v-for="(result, index) in analyzeResults" :key="`result-${index}`">
        <v-card-title v-if="result.rootword">Rootword: {{ result.rootword }}</v-card-title>
        <v-card-title v-if="result.conjugations && result.conjugations.length">Conjugations</v-card-title>
        <v-card-text v-if="result.conjugations && result.conjugations.length">
          <ul class="aligned-list">
            <li v-for="(conjugation, index) in result.conjugations" :key="`conjugation-${index}`"><strong>{{ conjugation.word }}</strong> : {{conjugation.rule}}</li>
          </ul>
        </v-card-text>
        <v-card-title v-if="result.similarWords && result.similarWords.length">Similar Words</v-card-title>
        <v-card-text v-if="result.similarWords && result.similarWords.length">
          <ul class="aligned-list">
            <li v-for="(similarWord, index) in result.similarWords" :key="`similarWord-${index}`"><strong>{{ similarWord.greek }}</strong> - {{ similarWord.english }}</li>
          </ul>
        </v-card-text>
        <v-carousel hide-delimiters show-arrows="hover" height="100%">
          <v-carousel-item v-for="(res, resIndex) in result.results" :key="`res-${index}-${resIndex}`">
            <v-card class="section-content">
              <v-card-title>{{ res.author }} - {{res.book}} Reference: {{ res.reference }}</v-card-title>
              <v-card-text>
                <p v-html="highlightWord(res.text.greek)"></p>
                <p> ---- </p>
                <p>{{res.text.translations[0]}}</p>
                <br>
                <a :href="res.referenceLink" target="_blank"> Link to Interactive Text</a>
              </v-card-text>
            </v-card>
          </v-carousel-item>
        </v-carousel>
      </v-card>
    </v-col>
  </v-row>
</template>

<style scoped>
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

.aligned-list {
  text-align: left; /* Ensure left alignment */
  padding-left: 20px; /* Add some padding to the left */
  margin: 0; /* Remove default margin */
}

.aligned-list li {
  list-style-type: none; /* Remove default list style */
  margin-bottom: 5px; /* Add some space between items */
}

.aligned-list li::before {
  content: "•"; /* Add a custom bullet */
  color: #333; /* Color of the bullet */
  display: inline-block; /* Align bullet with text */
  width: 1em; /* Space for the bullet */
  margin-left: -1em; /* Position the bullet */
}

</style>