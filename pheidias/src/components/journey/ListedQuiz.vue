<template>
  <div class="listed-quiz-card">
    <div class="journey-section-kicker">Choice step</div>
    <h3 v-if="section.title" class="mb-2">{{ section.title }}</h3>
    <p v-if="section.text" class="section-text mb-3">{{ section.text }}</p>
    <p class="question-text font-weight-medium mb-4">{{ section.question }}</p>

    <v-list class="custom-list">
      <v-list-item
          v-for="(opt, i) in shuffledOptions"
          :key="i"
          :title="opt"
          @click="selectAnswer(opt)"
          class="journey-list-option mb-2"
          :disabled="isCorrect || attempted[opt]"
      />
    </v-list>

    <v-alert
        v-if="lastWrong"
        type="error"
        variant="tonal"
        class="mt-2"
    >
      Not quite — try again.
    </v-alert>
  </div>
</template>

<script setup>
import {onMounted, ref, watch} from 'vue'
import {shuffle} from "@/utils/sharedQuiz";

const props = defineProps({
  section: Object
})
const emit = defineEmits(['finish'])
const shuffledOptions = ref([])

const isCorrect = ref(false)
const lastWrong = ref(false)
const attempted = ref({}) // { option: true }

function selectAnswer(opt) {
  if (isCorrect.value || attempted.value[opt]) return

  if (opt === props.section.answer) {
    isCorrect.value = true
    setTimeout(() => emit('finish'), 1000)
  } else {
    attempted.value[opt] = true
    lastWrong.value = true
    setTimeout(() => (lastWrong.value = false), 1000)
  }
}

onMounted(() => {
  shuffledOptions.value = shuffle(props.section.options)
})

</script>

<style scoped>
.listed-quiz-card {
  max-width: 820px;
  margin: 0 auto;
  border: 1px solid rgba(28, 97, 209, 0.14);
  border-radius: 18px;
  background: #fefcf5;
  color: #20334f;
  padding: clamp(18px, 3vw, 28px);
  box-shadow: 0 12px 28px rgba(28, 97, 209, 0.1);
}

.journey-section-kicker {
  color: #1c61d1;
  font-size: 0.74rem;
  font-weight: 900;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.listed-quiz-card h3 {
  color: #10284b;
  font-weight: 900;
}

.section-text,
.question-text {
  color: #536987;
  line-height: 1.65;
}

.custom-list {
  background: transparent;
}

.journey-list-option {
  border: 1px solid rgba(28, 188, 209, 0.16);
  border-radius: 12px;
  background: rgba(253, 246, 227, 0.72);
  color: #20334f;
  transition: transform 0.18s ease, box-shadow 0.18s ease, background 0.18s ease;
}

.journey-list-option:hover {
  transform: translateY(-2px);
  background: rgba(28, 188, 209, 0.12);
  box-shadow: 0 12px 24px rgba(28, 97, 209, 0.12);
}
</style>
