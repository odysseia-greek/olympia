<template>
  <div>
    <h3 class="mb-2">Final Translation</h3>
    <p class="mb-4">{{ section.instruction }}</p>

    <v-list>
      <v-list-item
          v-for="(opt, i) in section.options"
          :key="i"
          :title="opt"
          :color="getColor(opt)"
          @click="selectAnswer(opt)"
          class="mb-2"
          :disabled="selected !== null"
      />
    </v-list>

    <v-alert
        v-if="selected && selected !== section.answer"
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

function selectAnswer(opt) {
  if (selected.value) return
  selected.value = opt
  if (opt === props.section.answer) {
    setTimeout(() => emit('finishSegment'), 1200)
  }
}

function getColor(opt) {
  if (!selected.value) return 'primary'
  if (opt === props.section.answer) return 'success'
  if (opt === selected.value) return 'error'
  return ''
}
</script>