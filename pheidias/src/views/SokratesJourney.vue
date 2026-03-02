<template>
  <div>
    <JourneyQuiz />
  </div>
</template>

<script>
import { defineComponent } from 'vue';
import JourneyQuiz from '@/components/JourneyQuiz.vue';

export default defineComponent({
  name: 'SokratesJourney',
  components: {
    JourneyQuiz,
  },
});
</script>
