<template>
  <div class="match-game-area">
    <div class="journey-section-kicker">Match step</div>
    <h3>{{ section.instruction }}</h3>
    <v-row class="match-grid">
      <!-- Greek words -->
        <v-col cols="6">
          <div v-for="(item, i) in leftItems" :key="i">
              <v-card
                  class="match-card pa-3 text-center text-greek mb-5"
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
                class="match-card pa-3 text-center text-greek mb-5"
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
  if (feedback.value[word] === 'correct') return 'secondary'
  if (feedback.value[word] === 'error') return 'error'
  if (selected.value[side] === word) return 'primary'
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
  min-height: 35em;
  position: relative;
  max-width: 900px;
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

.match-game-area h3 {
  margin: 8px 0 22px;
  color: #10284b;
  font-weight: 900;
}

.match-grid {
  row-gap: 10px;
}

.match-card {
  border: 1px solid rgba(28, 188, 209, 0.16);
  border-radius: 12px;
  color: #20334f;
  cursor: pointer;
  box-shadow: 0 10px 22px rgba(28, 97, 209, 0.08);
  transition: transform 0.18s ease, box-shadow 0.18s ease;
}

.match-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 14px 28px rgba(28, 97, 209, 0.14);
}
</style>
