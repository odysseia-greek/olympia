<template>
  <div>
    <GrammarQuiz />
  </div>
</template>

<script>
import { defineComponent } from 'vue';
import GrammarQuiz from '@/components/GrammarQuiz.vue';

export default defineComponent({
  name: 'SokratesGrammar',
  components: {
    GrammarQuiz,
  },
});
</script>
