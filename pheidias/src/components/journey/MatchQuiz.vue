<template>
  <div class="match-game-area">
    <h3 class="ma-5">{{ section.instruction }}</h3>
      <br>
    <v-row class="ma-5">
      <!-- Greek words -->
        <v-col cols="6">
          <div v-for="(item, i) in leftItems" :key="i">
              <v-card
                  class="pa-3 text-center text-greek mb-10"
                  :color="getCardColor(item.greek, 'left')"
                  @click="selectItem(item.greek, 'left')"
                  v-show="!matched[item.greek]"
              >
                {{ item.greek }}
              </v-card>
          </div>
        </v-col>

      <!-- Answers -->
      <v-col cols="6">
        <div v-for="(item, i) in rightItems" :key="i">
            <v-card
                class="pa-3 text-center text-greek mb-10"
                :color="getCardColor(item.answer, 'right')"
                @click="selectItem(item.answer, 'right')"
                v-show="!matchedByAnswer[item.answer]"
            >
              {{ item.answer }}
            </v-card>
        </div>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
import {ref, computed, watch, onMounted} from 'vue'
import {shuffle} from "@/utils/sharedQuiz";

const props = defineProps({
  section: Object
})

const emit = defineEmits(['finish'])

const leftItems = ref([])
const rightItems =ref([])

const selected = ref({ left: null, right: null })
const feedback = ref({}) // { word: 'error' | 'correct' }

const matched = ref({})
const matchedByAnswer = ref({})


const remainingPairs = computed(() => props.section.pairs.filter(
    p => !matched.value[p.greek] && !matchedByAnswer.value[p.answer]
))

function selectItem(word, side) {
  // Don't allow reselect of matched items
  if ((side === 'left' && matched.value[word]) || (side === 'right' && matchedByAnswer.value[word])) return

  selected.value[side] = word

  if (selected.value.left && selected.value.right) {
    checkMatch()
  }
}

function checkMatch() {
  const left = selected.value.left
  const right = selected.value.right
  const isCorrect = props.section.pairs.some(p => p.greek === left && p.answer === right)

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

      if (remainingPairs.value.length === 0) {
        emit('finish')
      }
    }, 1200) // delay before hiding matched cards
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
  if (feedback.value[word] === 'correct') return '#1de9b6'
  if (feedback.value[word] === 'error') return 'red'
  if (selected.value[side] === word) return 'triadic'
  return '#fefcf5'
}

onMounted(() => {
  leftItems.value = shuffle(props.section.pairs)
  rightItems.value = shuffle([...props.section.pairs]) // independent shuffle
})
</script>

<style scoped>
.text-greek {
  font-family: "EB Garamond", serif;
  font-size: 1.25rem;
}

.match-game-area {
  min-height: 35em; /* or 400px, based on your tallest layout */
  position: relative;
}
</style>