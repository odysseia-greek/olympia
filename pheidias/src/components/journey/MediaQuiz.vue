<template>
  <div class="media-match-area">
    <h3 class="mb-4">{{ section.instruction }}</h3>
    <v-row dense>
      <v-col
          v-for="(item, i) in shuffledWords"
          :key="i"
          cols="6"
      >
        <v-card
            class="pa-2 text-center text-greek"
            :color="getCardColor(item.word, 'left')"
            @click="selectItem(item.word, 'left')"
            v-show="!matched[item.word]"
        >
          {{ item.word }}
        </v-card>
      </v-col>

    <v-col
        v-for="(item, i) in shuffledImages"
        :key="i"
        cols="6"
    >
      <v-card
          class="pa-2 text-center"
          :color="getCardColor(item.answer, 'right')"
          @click="selectItem(item.answer, 'right')"
          v-show="!matchedByAnswer[item.answer]"
      >
        <v-img
            :src="loadedImages[item.answer]"
            aspect-ratio="1"
            class="mx-auto"
        />
      </v-card>
    </v-col>
    </v-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import {shuffle} from "@/utils/sharedQuiz";

const props = defineProps({
  section: Object
})

const emit = defineEmits(['finish'])

const selected = ref({ left: null, right: null })
const feedback = ref({})
const matched = ref({})
const loadedImages = ref({})
const shuffledWords = ref([])
const shuffledImages = ref([])
const matchedByAnswer = ref({})
const imageMap = import.meta.glob('/src/assets/icons/*.webp');

async function loadImages(files) {
  for (const item of files) {
    const key = item.answer
    if (!loadedImages.value[key]) {
      const path = `/src/assets/icons/${key}`;
      const importer = imageMap[path];
      if (importer) {
        try {
          const module = await importer();
          loadedImages.value[key] = module.default;
        } catch (e) {
          console.warn(`Image ${key} failed to load`, e);
        }
      } else {
        console.warn(`No image found for ${path}`);
      }
    }
  }
}


function selectItem(value, side) {
  if ((side === 'left' && matched.value[value]) || (side === 'right' && matchedByAnswer.value[value])) return

  selected.value[side] = value

  if (selected.value.left && selected.value.right) {
    checkMatch()
  }
}

function checkMatch() {
  const left = selected.value.left
  const right = selected.value.right
  const isCorrect = props.section.mediaFiles.some(item => item.word === left && item.answer === right)

  if (isCorrect) {
    feedback.value[left] = 'correct'
    feedback.value[right] = 'correct'
    setTimeout(() => {
      matched.value[left] = true
      matchedByAnswer.value[right] = true
      delete feedback.value[left]
      delete feedback.value[right]
      selected.value.left = null
      selected.value.right = null

      // Check if all items are matched
      const allMatched = props.section.mediaFiles.every(item =>
          matched.value[item.word] && matchedByAnswer.value[item.answer]
      )

      if (allMatched) {
        setTimeout(() => {
          emit('finish')
        }, 1200)
      }
    }, 1000)
  } else {
    feedback.value[left] = 'error'
    feedback.value[right] = 'error'
    setTimeout(() => {
      delete feedback.value[left]
      delete feedback.value[right]
      selected.value.left = null
      selected.value.right = null
    }, 700)
  }
}

function getCardColor(word, side) {
  if (feedback.value[word] === 'correct') return 'secondary'
  if (feedback.value[word] === 'error') return 'red'
  if (selected.value[side] === word) return 'triadic'
  return '#fefcf5'
}


onMounted(() => {
  const copy = [...props.section.mediaFiles]
  shuffledWords.value = shuffle(copy)
  shuffledImages.value = shuffle(copy)
  loadImages(copy)
})

</script>

<style scoped>
.media-match-area {
  min-height: 30em;
  position: relative;
  max-width: 820px;
  margin: auto;
}
.text-greek {
  font-family: "EB Garamond", serif;
  font-size: 1.25rem;
}
</style>