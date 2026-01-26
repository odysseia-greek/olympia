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


const { result: optionsResult, loading, onResult } = useQuery(JourneyOptions)

onResult(({ data }) => {
  if (data && data.journeyOptions) {
    journeyThemes.value = data.journeyOptions.themes
  }
})

const quizProgressPercent = computed(() =>
    ((currentStep.value + 1) / currentQuiz.value.quiz.length) * 100
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
    mapMode.value = false

  } catch (err) {
    console.error('Error fetching media quiz:', err);
  }
};

function handleNextStep() {
  if (currentStep.value < currentQuiz.value.quiz.length - 1) {
    currentStep.value++
  } else {
    currentStep.value = 0
  }
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

  mapMode.value = true
}


</script>

<template>
  <v-container fluid class="pa-0">
    <!-- Header -->
    <v-card
        class="paper-card pa-6 mb-4"
        elevation="4"
        v-show="mapMode"
    >
      <div class="d-flex justify-space-between align-center mb-3">
        <v-card-title class="text-h5" v-if="!selectedJourney">Journey Mode</v-card-title>
        <v-card-title class="text-h5"  v-if="selectedJourney">{{selectedJourney.name}}</v-card-title>

        <v-btn
            icon
            variant="text"
            @click="showIntro = !showIntro"
        >
          <v-icon>{{ showIntro ? 'mdi-minus' : 'mdi-plus' }}</v-icon>
        </v-btn>
      </div>

      <v-expand-transition>
        <div v-show="showIntro">
          <p class="mb-4">
            In Journey Mode, you explore the world of Ancient Greece by selecting a <strong>theme</strong> on the map. Each theme consists of <strong>segments</strong> — historical or literary moments tied to real-world locations.
          </p>
          <p class="mb-4">
            In each segment, you'll read authentic Greek texts and deepen your understanding through <strong>matching</strong>, <strong>grammar</strong>, and <strong>trivia</strong> questions. Your path unfolds as you progress — each completed segment unlocking the next.
          </p>
          <v-alert type="info" variant="tonal" class="mt-4">
            <strong>“Πάντα ῥεῖ καὶ οὐδὲν μένει.”</strong><br />
            <em>“Everything flows, nothing stands still.” — Herakleitos</em>
          </v-alert>
        </div>
      </v-expand-transition>
      <v-btn
          v-if="selectedJourney"
          @click="selectedJourney = null; mapMode = true; showIntro = true;"
          class="ma-2"
          color="white"
      >
        Back to Overview
      </v-btn>
    </v-card>

    <!-- Map with Hotspots -->
        <div class="map-container" v-if="mapMode">
          <img src="@/assets/ancient_greece.svg" class="map-image" />
          <div
              v-if="!selectedJourney"
              v-for="point in hotspots"
              :key="point.id"
              class="hotspot"
              :style="getStyle(point.coordinates)"
              @click="selectThemePoint(point)"
          >
            <v-icon color="primary" size="40">mdi-map-marker</v-icon>
            <!-- Show label only for the first segment -->
            <div
                class="hotspot-label"
            >
              {{ point.title }}
            </div>
          </div>
    <div
        v-if="selectedJourney"
    >
      <div
          v-for="segment in selectedJourney.segments"
          :key="segment.name"
          class="hotspot"
          :style="getStyle(segment.coordinates)"
      >
        <v-icon
            size="40"
            :color="getColorForSegment(segment)"
            :class="{ 'clickable': isSegmentSelectable(segment) }"
            @click="isSegmentSelectable(segment) && getJourneyQuiz(segment)"
        >
          mdi-map-marker
        </v-icon>
        <div class="hotspot-label">{{segment.number}}. {{ segment.name }}</div>
      </div>
    </div>
  </div>
    <v-container class="pa-4" v-if="!mapMode">
      <v-card class="paper-card" elevation="4">
        <!-- Header -->
        <v-row>
          <v-col cols="12">
            <h3>{{ currentQuiz.theme }} — {{ currentQuiz.segment }}</h3>
            <p class="text-subtitle-2 font-italic">{{ currentQuiz.contextNote }}</p>
          </v-col>
        </v-row>

        <!-- Sentence & Translation -->
        <v-row class="mt-2">
          <v-col cols="12">
            <div class="greek-sentence">{{ currentQuiz.sentence }}</div>
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
            <div class="d-flex justify-space-between mt-1 text-caption">
              <span>Step {{ currentStep + 1 }} of {{ currentQuiz.quiz.length }}</span>
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
.map-image {
  width: 100%;
  display: block;
}
.hotspot {
  position: absolute;
  cursor: pointer;
  z-index: 2;
}

.map-container {
  position: relative;
  width: 95%; /* Default for mobile/smaller screens */
  max-width: 1500px; /* Maximum width cap */
  margin: auto;
}

.hotspot-label {
  position: absolute;
  top: -1.5em;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(255, 255, 255, 0.85);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 16px;
  white-space: nowrap;
  pointer-events: none;
}

.paper-card {
  background: #fdf6e3; /* A light, papyrus-like color */
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
  font-family: "Roboto", serif;
}

.clickable {
  cursor: pointer;
}
</style>