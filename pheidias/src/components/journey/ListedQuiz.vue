<template>
  <div>
    <h3 v-if="section.title" class="mb-2">{{ section.title }}</h3>
    <p v-if="section.text" class="mb-3">{{ section.text }}</p>
    <p class="font-weight-medium mb-4">{{ section.question }}</p>

    <v-list class="custom-list">
      <v-list-item
          v-for="(opt, i) in shuffledOptions"
          :key="i"
          :title="opt"
          @click="selectAnswer(opt)"
          class="mb-2"
          :disabled="isCorrect || attempted[opt]"
      />
    </v-list>

    <v-alert
        v-if="lastWrong"
        type="error"
        variant="tonal"
        class="mt-2"
    >
      ❌ Not quite — try again.
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

<style>
.custom-list {
  background: #fefcf5;
}

</style>