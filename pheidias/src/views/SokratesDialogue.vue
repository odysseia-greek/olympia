<template>
  <div>
    <DialogueQuiz />
  </div>
</template>

<script>
import { defineComponent } from 'vue';
import DialogueQuiz from '@/components/DialogueBasedQuiz.vue';

export default defineComponent({
  name: 'SokratesDialogue',
  components: {
    DialogueQuiz,
  },
});
</script>
