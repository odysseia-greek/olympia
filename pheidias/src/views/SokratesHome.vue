<template>
  <div>
    <QuizAreaHome />
  </div>
</template>

<script>
import { defineComponent } from 'vue';
import QuizAreaHome from '@/components/QuizHome.vue';

export default defineComponent({
  name: 'SokratesHome',
  components: {
    QuizAreaHome,
  },
});
</script>
