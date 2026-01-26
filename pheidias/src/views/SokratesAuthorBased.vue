<template>
  <div>
    <AuthorBasedQuiz />
  </div>
</template>

<script>
import { defineComponent } from 'vue';
import AuthorBasedQuiz from '@/components/AuthorBasedQuiz.vue';

export default defineComponent({
  name: 'SokratesAuthorBased',
  components: {
    AuthorBasedQuiz,
  },
});
</script>
