<template>
  <div class="media-match-area">
    <div class="journey-section-kicker">Image match step</div>
    <h3 class="mb-4">{{ section.instruction }}</h3>
    <v-row class="media-grid" dense>
      <v-col
          v-for="(item, i) in shuffledWords"
          :key="i"
          cols="6"
      >
        <v-card
            class="media-choice-card pa-2 text-center text-greek"
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
          class="media-choice-card image-choice-card pa-2 text-center"
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
  if (feedback.value[word] === 'error') return 'error'
  if (selected.value[side] === word) return 'primary'
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

.media-match-area h3 {
  color: #10284b;
  font-weight: 900;
}

.media-grid {
  row-gap: 12px;
}

.media-choice-card {
  min-height: 56px;
  border: 1px solid rgba(28, 188, 209, 0.16);
  border-radius: 12px;
  color: #20334f;
  cursor: pointer;
  box-shadow: 0 10px 22px rgba(28, 97, 209, 0.08);
  transition: transform 0.18s ease, box-shadow 0.18s ease;
}

.media-choice-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 14px 28px rgba(28, 97, 209, 0.14);
}

.image-choice-card {
  overflow: hidden;
}

.text-greek {
  font-family: "EB Garamond", serif;
  font-size: 1.25rem;
}
</style>
