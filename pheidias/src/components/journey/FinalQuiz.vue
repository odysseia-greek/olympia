<template>
  <div class="final-quiz-card">
    <div class="journey-section-kicker">Final step</div>
    <h3 class="mb-2">Final Translation</h3>
    <p class="instruction-text mb-4">{{ section.instruction }}</p>

    <v-list class="final-options-list">
      <v-list-item
          v-for="(opt, i) in section.options"
          :key="i"
          :title="opt"
          :color="getColor(opt)"
          @click="selectAnswer(opt)"
          class="final-option mb-2"
          :disabled="selected === section.answer"
      />
    </v-list>

    <v-alert
        v-if="wrongAnswer"
        type="error"
        class="mt-4"
    >
      Not quite. Remember the tone and intent of the original.
    </v-alert>
    <v-alert
        v-if="selected === section.answer"
        type="success"
        class="mt-4"
    >
      Excellent! You’ve captured the essence of the passage.
    </v-alert>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  section: Object
})

const emit = defineEmits(['finishSegment'])

const selected = ref(null)
const wrongAnswer = ref(false)

function selectAnswer(opt) {
  if (selected.value === props.section.answer) return

  wrongAnswer.value = false
  selected.value = opt
  if (opt === props.section.answer) {
    setTimeout(() => emit('finishSegment'), 1200)
  } else {
    wrongAnswer.value = true
    setTimeout(() => {
      if (selected.value === opt) {
        selected.value = null
      }
      wrongAnswer.value = false
    }, 1200)
  }
}

function getColor(opt) {
  if (!selected.value) return 'primary'
  if (opt === props.section.answer) return 'success'
  if (opt === selected.value) return 'error'
  return ''
}
</script>

<style scoped>
.final-quiz-card {
  max-width: 860px;
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

.final-quiz-card h3 {
  color: #10284b;
  font-weight: 900;
}

.instruction-text {
  color: #536987;
  line-height: 1.65;
}

.final-options-list {
  background: transparent;
}

.final-option {
  border: 1px solid rgba(28, 188, 209, 0.16);
  border-radius: 12px;
  background: rgba(253, 246, 227, 0.72);
  color: #20334f;
  transition: transform 0.18s ease, box-shadow 0.18s ease, background 0.18s ease;
}

.final-option:hover {
  transform: translateY(-2px);
  background: rgba(28, 188, 209, 0.12);
  box-shadow: 0 12px 24px rgba(28, 97, 209, 0.12);
}
</style>
