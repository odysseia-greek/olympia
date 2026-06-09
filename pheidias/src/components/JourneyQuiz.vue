<script setup>
import {computed, onMounted, ref} from 'vue'

import {JourneyCreateQuestion, JourneyOptions} from '@/constants/journeyBasedGraphql'
import {useQuery} from "@vue/apollo-composable";
import {getJourneyProgress, updateJourneyProgress, useBouleId} from "@/composables/useBoule";
import {apolloClient} from "@/apollo";
import JourneyQuizArea from "@/components/journey/JourneyQuizArea.vue";


const mapMode = ref(true)
const journeyThemes = ref([])
const selectedJourney = ref(null)
const selectedPoint = ref(null)
const lastUnlocked = ref(1)
const completedSegments = ref({})

const currentQuiz = ref(null)
const boule = useBouleId();
const currentStep = ref(0)
const showIntro = ref(true)
const completedQuizSteps = ref([])


const { result: optionsResult, loading, onResult } = useQuery(JourneyOptions)

onResult(({ data }) => {
  if (data && data.journeyOptions) {
    journeyThemes.value = data.journeyOptions.themes
  }
})

const quizProgressPercent = computed(() =>
    ((currentStep.value + 1) / currentQuiz.value.quiz.length) * 100
)

const currentSection = computed(() => currentQuiz.value?.quiz?.[currentStep.value])

const isFinalStep = computed(() => currentSection.value?.__typename === 'FinalTranslationQuiz')

const canGoNext = computed(() =>
    currentQuiz.value
    && currentStep.value < currentQuiz.value.quiz.length - 1
    && isStepCompleted(currentStep.value)
)

const hotspots = computed(() => {
  return journeyThemes.value.map(theme => {
    const firstSegment = theme.segments.find(s => s.number === 1)
    return {
      id: `theme-${theme.name}`,
      title: theme.name,
      summary: `Begin your journey at ${firstSegment.location}`,
      coordinates: firstSegment.coordinates,
      theme,
      completed: isJourneyComplete(theme),
    }
  })
})

function getStyle({ x, y }) {
  return {
    position: 'absolute',
    left: `${x * 100}%`,
    top: `${y * 100}%`,
    transform: 'translate(-50%, -100%)'
  }
}

function getHotspotLabelClass({ x, y }) {
  return {
    'label-left-edge': x < 0.16,
    'label-right-edge': x > 0.84,
    'label-top-edge': y < 0.12,
  }
}

function selectThemePoint(point) {
  const progress = getJourneyProgress(point.theme.name)
  completedSegments.value[point.theme.name] = progress.completed
  lastUnlocked.value = progress.current
  selectedJourney.value = point.theme
  showIntro.value = false
  selectedPoint.value = point
}

function getCompletionStatus(segment) {
  const completedUpTo = getCompletedNumberForTheme(segment.theme || selectedJourney.value.name)

  if (segment.number <= completedUpTo) return 'finished'
  if (segment.number === completedUpTo + 1) return 'current'
  return 'locked'
}

function getCompletedNumberForTheme(themeName) {
  const segments = completedSegments.value[themeName] || []
  return segments.length > 0 ? Math.max(...segments) : 0
}

function isJourneyComplete(theme) {
  const progress = completedSegments.value[theme.name] || getJourneyProgress(theme.name).completed || []
  return theme.segments.every(segment => progress.includes(segment.number))
}

function getColorForSegment(segment) {
  const status = getCompletionStatus(segment)
  if (status === 'finished') return 'secondary'
  if (status === 'current') return 'triadic'
  return 'grey lighten-1'
}

function isSegmentSelectable(segment) {
  const status = getCompletionStatus(segment)
  return status === 'current' || status === 'finished'
}

const getJourneyQuiz = async (segment) => {
  try {
    const { data } = await apolloClient.query({
      query: JourneyCreateQuestion,
      variables: {
        input: {
          theme: selectedJourney.value.name,
          segment: segment.name,
        },
      },
      context: {
        headers: {
          'boule': boule,
        },
      },
      fetchPolicy: 'no-cache',
    });

    currentQuiz.value = data.journeyQuiz;
    currentStep.value = 0
    completedQuizSteps.value = []
    mapMode.value = false

  } catch (err) {
    console.error('Error fetching media quiz:', err);
  }
};

function isStepCompleted(step) {
  return completedQuizSteps.value.includes(step)
}

function markCurrentStepComplete() {
  if (!isStepCompleted(currentStep.value)) {
    completedQuizSteps.value = [...completedQuizSteps.value, currentStep.value].sort((a, b) => a - b)
  }
}

function handleNextStep() {
  markCurrentStepComplete()

  if (currentStep.value < currentQuiz.value.quiz.length - 1) {
    currentStep.value++
  }
}

function handlePreviousStep() {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

function returnToMap() {
  currentStep.value = 0
  mapMode.value = true
}

function markSegmentAsComplete() {
  const theme = selectedJourney.value.name
  const number = currentQuiz.value.number
  currentStep.value = 0

  if (!completedSegments.value[theme]) {
    completedSegments.value[theme] = []
  }

  if (!completedSegments.value[theme].includes(number)) {
    completedSegments.value[theme].push(number)
    completedSegments.value[theme].sort((a, b) => a - b)

    updateJourneyProgress(theme, number)
  }

  const selectedTheme = journeyThemes.value.find(journey => journey.name === theme)
  if (selectedTheme && isJourneyComplete(selectedTheme)) {
    completedSegments.value[theme] = selectedTheme.segments.map(segment => segment.number)
    selectedJourney.value = null
    showIntro.value = true
  }

  mapMode.value = true
}


</script>

<template>
  <v-container fluid class="journey-page pa-0">
    <!-- Header -->
    <v-card
        class="paper-card journey-hero pa-6 mb-4"
        elevation="4"
        v-show="mapMode"
    >
      <div class="journey-toolbar">
        <div>
          <div class="journey-kicker">Preview path</div>
          <v-card-title class="journey-title text-h5" v-if="!selectedJourney">Journey Mode</v-card-title>
          <v-card-title class="journey-title text-h5"  v-if="selectedJourney">{{selectedJourney.name}}</v-card-title>
        </div>

        <v-btn
            icon
            variant="text"
            color="primary"
            @click="showIntro = !showIntro"
        >
          <v-icon>{{ showIntro ? 'mdi-minus' : 'mdi-plus' }}</v-icon>
        </v-btn>
      </div>

      <v-expand-transition>
        <div v-show="showIntro" class="journey-intro">
          <p>
            In Journey Mode, you explore the world of Ancient Greece by selecting a <strong>theme</strong> on the map. Each theme consists of <strong>segments</strong> — historical or literary moments tied to real-world locations.
          </p>
          <p>
            In each segment, you'll read authentic Greek texts and deepen your understanding through <strong>matching</strong>, <strong>grammar</strong>, and <strong>trivia</strong> questions. Your path unfolds as you progress — each completed segment unlocking the next.
          </p>
          <v-alert type="info" variant="tonal" class="journey-quote mt-4">
            <strong>“Πάντα ῥεῖ καὶ οὐδὲν μένει.”</strong><br />
            <em>“Everything flows, nothing stands still.” — Herakleitos</em>
          </v-alert>
        </div>
      </v-expand-transition>
      <v-btn
          v-if="selectedJourney"
          @click="selectedJourney = null; mapMode = true; showIntro = true;"
          class="ma-2"
          color="primary"
          variant="tonal"
          prepend-icon="mdi-map-outline"
      >
        Back to Overview
      </v-btn>
    </v-card>

    <!-- Map with Hotspots -->
        <div class="map-shell" v-if="mapMode">
        <div class="map-container">
          <img src="@/assets/ancient_greece.svg" class="map-image" alt="Map of Ancient Greece" />
              <div
                  v-if="!selectedJourney"
                  v-for="point in hotspots"
                  :key="point.id"
                  class="hotspot theme-hotspot"
                  :class="{ 'is-complete': point.completed }"
                  :style="getStyle(point.coordinates)"
                  @click="selectThemePoint(point)"
          >
            <v-icon color="primary" size="44">mdi-map-marker</v-icon>
            <!-- Show label only for the first segment -->
            <div
                class="hotspot-label"
                :class="getHotspotLabelClass(point.coordinates)"
            >
              {{ point.title }}
              <span v-if="point.completed">Complete</span>
            </div>
          </div>
    <div
        v-if="selectedJourney"
    >
      <div
          v-for="segment in selectedJourney.segments"
          :key="segment.name"
          class="hotspot segment-hotspot"
          :class="`is-${getCompletionStatus(segment)}`"
          :style="getStyle(segment.coordinates)"
      >
        <v-icon
            size="44"
            :color="getColorForSegment(segment)"
            :class="{ 'clickable': isSegmentSelectable(segment) }"
            @click="isSegmentSelectable(segment) && getJourneyQuiz(segment)"
        >
          mdi-map-marker
        </v-icon>
        <div class="hotspot-label" :class="getHotspotLabelClass(segment.coordinates)">{{segment.number}}. {{ segment.name }}</div>
      </div>
    </div>
  </div>
        </div>
    <v-container class="journey-quiz-container pa-4" v-if="!mapMode">
      <v-card class="paper-card journey-quiz-card" elevation="4">
        <!-- Header -->
        <v-row class="journey-quiz-header">
          <v-col cols="12">
            <div class="journey-kicker">Selected segment</div>
            <h3>{{ currentQuiz.theme }} — {{ currentQuiz.segment }}</h3>
            <p>{{ currentQuiz.contextNote }}</p>
          </v-col>
        </v-row>

        <!-- Sentence & Translation -->
        <v-row class="mt-2">
          <v-col cols="12">
            <div class="journey-sentence greek-sentence">{{ currentQuiz.sentence }}</div>
          </v-col>
        </v-row>

        <!-- Progress Bar for Quiz Sections -->
        <v-row class="mt-4">
          <v-col cols="12">
            <v-progress-linear
                :model-value="quizProgressPercent"
                height="8"
                color="primary"
                rounded
            ></v-progress-linear>
            <div class="journey-step-label d-flex justify-space-between mt-1 text-caption">
              <span>Step {{ currentStep + 1 }} of {{ currentQuiz.quiz.length }}</span>
            </div>
            <div class="journey-step-controls">
              <v-btn
                  color="primary"
                  variant="tonal"
                  prepend-icon="mdi-arrow-left"
                  :disabled="currentStep === 0"
                  @click="handlePreviousStep"
              >
                Previous
              </v-btn>
              <v-btn
                  v-if="!isFinalStep"
                  color="secondary"
                  variant="tonal"
                  prepend-icon="mdi-map-outline"
                  @click="returnToMap"
              >
                Back to Map
              </v-btn>
              <v-btn
                  color="primary"
                  variant="tonal"
                  append-icon="mdi-arrow-right"
                  :disabled="!canGoNext"
                  @click="handleNextStep"
              >
                Next
              </v-btn>
            </div>
          </v-col>
        </v-row>

        <!-- Quiz Section Render -->
        <v-row class="mt-4">
          <v-col cols="12">
            <JourneyQuizArea
                :section="currentQuiz.quiz[currentStep]"
                :text="currentQuiz.sentence"
                :translation="currentQuiz.translation"
                :is-last-step="currentStep >= currentQuiz.quiz.length - 1"
                @next="handleNextStep"
                @finishSegment="markSegmentAsComplete"
            />
          </v-col>
        </v-row>
      </v-card>
    </v-container>
  </v-container>
</template>


<style scoped>
.journey-page {
  color: #20334f;
}

.map-image {
  width: 100%;
  display: block;
  filter: sepia(0.12) saturate(0.96) contrast(0.98);
}

.hotspot {
  position: absolute;
  cursor: pointer;
  z-index: 2;
}

.hotspot :deep(.v-icon) {
  filter: drop-shadow(0 8px 12px rgba(16, 40, 75, 0.28));
  transition: transform 0.18s ease, filter 0.18s ease;
}

.hotspot:hover :deep(.v-icon),
.hotspot .clickable:hover {
  transform: translateY(-4px) scale(1.06);
  filter: drop-shadow(0 12px 18px rgba(28, 97, 209, 0.34));
}

.map-container {
  position: relative;
  width: 100%;
  max-width: 1500px;
  margin: auto;
}

.map-shell {
  width: min(96%, 1540px);
  margin: 0 auto 28px;
  padding: clamp(10px, 2vw, 24px);
  border: 1px solid rgba(28, 97, 209, 0.16);
  border-radius: 24px;
  background:
      radial-gradient(circle at 14% 18%, rgba(28, 209, 140, 0.16), transparent 30%),
      linear-gradient(145deg, rgba(254, 252, 245, 0.98), rgba(253, 246, 227, 0.9));
  box-shadow: 0 24px 60px rgba(28, 97, 209, 0.12);
  overflow: visible;
}

.hotspot-label {
  position: absolute;
  top: -2.1em;
  left: 50%;
  transform: translateX(-50%);
  border: 1px solid rgba(28, 97, 209, 0.14);
  background: rgba(254, 252, 245, 0.94);
  color: #20334f;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 14px;
  font-weight: 800;
  white-space: nowrap;
  pointer-events: none;
  box-shadow: 0 10px 24px rgba(16, 40, 75, 0.14);
}

.hotspot-label.label-left-edge {
  left: 0;
  transform: translateX(0);
}

.hotspot-label.label-right-edge {
  left: auto;
  right: 0;
  transform: translateX(0);
}

.hotspot-label.label-top-edge {
  top: 2.6em;
}

.hotspot-label span {
  display: block;
  color: #1cd18c;
  font-size: 0.68rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.theme-hotspot.is-complete .hotspot-label {
  border-color: rgba(28, 209, 140, 0.34);
  box-shadow: 0 0 0 5px rgba(28, 209, 140, 0.12), 0 10px 24px rgba(16, 40, 75, 0.14);
}

.paper-card {
  background:
      linear-gradient(145deg, rgba(254, 252, 245, 0.98), rgba(253, 246, 227, 0.92)),
      #fdf6e3;
  border: 1px solid rgba(28, 97, 209, 0.14);
  border-radius: 18px;
  box-shadow: 0 14px 34px rgba(28, 97, 209, 0.12);
  padding: 20px;
  color: #20334f;
}

.journey-hero {
  max-width: 1180px;
  margin-inline: auto;
}

.journey-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 14px;
}

.journey-kicker {
  color: #1c61d1;
  font-size: 0.76rem;
  font-weight: 900;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.journey-title {
  padding-left: 0;
  color: #10284b;
  font-size: clamp(1.5rem, 3vw, 2.3rem);
  font-weight: 900;
}

.journey-intro {
  max-width: 860px;
  color: #536987;
  line-height: 1.75;
}

.journey-quote {
  border: 1px solid rgba(28, 188, 209, 0.18);
}

.journey-quiz-container {
  max-width: 1100px;
}

.journey-quiz-card {
  overflow: hidden;
}

.journey-quiz-header h3 {
  color: #10284b;
  font-size: clamp(1.35rem, 3vw, 2rem);
  font-weight: 900;
}

.journey-quiz-header p {
  color: #536987;
  font-style: italic;
}

.journey-sentence {
  border: 1px solid rgba(28, 188, 209, 0.16);
  border-radius: 16px;
  background:
      radial-gradient(circle at 16% 20%, rgba(28, 209, 140, 0.16), transparent 32%),
      #fefcf5;
  color: #10284b;
  font-family: "EB Garamond", serif;
  font-size: clamp(1.45rem, 3vw, 2.25rem);
  line-height: 1.7;
  padding: 24px;
}

.journey-step-label {
  color: #536987;
}

.journey-step-controls {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 12px;
  margin-top: 16px;
}

.clickable {
  cursor: pointer;
}

.segment-hotspot.is-locked {
  cursor: not-allowed;
}

.segment-hotspot.is-locked .hotspot-label {
  opacity: 0.68;
}

@media (max-width: 700px) {
  .hotspot-label {
    font-size: 11px;
    padding: 3px 7px;
  }

  .journey-hero {
    padding: 18px !important;
  }

  .journey-toolbar {
    align-items: flex-start;
  }

  .journey-step-controls .v-btn {
    flex: 1 1 100%;
  }
}
</style>
