<template>
  <div>
    <MultipleChoiceQuiz />
  </div>
</template>

<script>
import { defineComponent } from 'vue';
import MultipleChoiceQuiz from '@/components/MultipleChoiceQuiz.vue';

export default defineComponent({
  name: 'SokratesMultipleChoice',
  components: {
    MultipleChoiceQuiz,
  },
});
</script>
