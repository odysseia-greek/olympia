<template>
  <div>
    <MediaQuiz />
  </div>
</template>

<script>
import { defineComponent } from 'vue';
import MediaQuiz from '@/components/MediaQuiz.vue';

export default defineComponent({
  name: 'SokratesMedia',
  components: {
    MediaQuiz,
  },
});
</script>
