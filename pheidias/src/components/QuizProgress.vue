<template>
  <v-expansion-panels flat class="my-6" variant="accordion">
    <v-expansion-panel style="background: #fdf6e3;" active>
      <v-expansion-panel-title>
        <v-icon start icon="mdi-chart-bar"></v-icon>
        Session Stats
      </v-expansion-panel-title>

      <v-expansion-panel-text>
        <v-row class="mb-4" dense>
          <v-col cols="6" sm="3">
            <v-card flat style="background: #fefcf5;">
              <v-card-subtitle>Seen</v-card-subtitle>
              <v-card-text class="text-h6">{{ seen }} / {{ itemsInThisSet }} </v-card-text>
            </v-card>
          </v-col>

          <v-col cols="6" sm="3">
            <v-card flat style="background: #fefcf5;">
              <v-card-subtitle>Mistakes</v-card-subtitle>
              <v-card-text class="text-h6 text-error">{{ mistakes }}</v-card-text>
            </v-card>
          </v-col>

          <v-col cols="6" sm="3">
            <v-card flat style="background: #fefcf5;">
              <v-card-subtitle>Correct</v-card-subtitle>
              <v-card-text class="text-h6 text-success">
                {{ correct }} / {{ neededCorrectCount }}
              </v-card-text>
              <v-progress-linear
                  :model-value="progressTowardCompletion"
                  color="success"
                  height="6"
                  class="mt-2"
              />
            </v-card>
          </v-col>

          <v-col cols="6" sm="3">
            <v-card flat style="background: #fefcf5;">
              <v-card-subtitle>Streak</v-card-subtitle>
              <v-card-text class="text-h6">{{ streak }}</v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <v-divider class="my-4"></v-divider>

        <v-row align="center" justify="center">
          <v-col cols="12" class="text-center">
            <v-card-subtitle>Total Played: {{ totalPlayed }}</v-card-subtitle>

            <v-progress-circular
                class="mt-2"
                rotate="360"
                :color="progressColor"
                width="8"
                size="72"
                :model-value="currentCorrectness"
            >
              {{ currentCorrectness }}%
            </v-progress-circular>
          </v-col>
        </v-row>
      </v-expansion-panel-text>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script setup>
import {ref, computed, watch, onMounted} from 'vue';
import {getProgress, saveProgress} from "@/composables/useBoule";

const props = defineProps({
  progress: Array,
  neededCorrectCount: Number,
  streak: Number,
  totalPlayed: Number,
  totalMistakes: Number,
  itemsInThisSet: Number,
});


const green = [29, 233, 182]; // RGB for green
const orange = [255, 165, 0]; // RGB for orange
const red = [233, 29, 29]; // RGB for red

const progressColor = computed(() => {
  // Interpolate between green and orange, then orange and red
  let r, g, b;
  if (currentCorrectness.value > 50) {
    const ratio = (currentCorrectness.value - 50) / 50;
    r = Math.round(orange[0] * (1 - ratio) + green[0] * ratio);
    g = Math.round(orange[1] * (1 - ratio) + green[1] * ratio);
    b = Math.round(orange[2] * (1 - ratio) + green[2] * ratio);
  } else {
    const ratio = currentCorrectness.value / 50;
    r = Math.round(red[0] * (1 - ratio) + orange[0] * ratio);
    g = Math.round(red[1] * (1 - ratio) + orange[1] * ratio);
    b = Math.round(red[2] * (1 - ratio) + orange[2] * ratio);
  }

  return `rgb(${r}, ${g}, ${b})`;
});

const expanded = ref(true);
const emit = defineEmits(['updateWorst', 'initProgress']);

const seen = computed(() => props.progress.filter(p => p.playCount > 0).length);
const mistakes = computed(() => props.progress.reduce((acc, p) => acc + p.incorrectCount, 0));
const correct = computed(() => props.progress.reduce((acc, p) => acc + p.correctCount, 0));
const streak = ref(0);
const totalPlayed = ref(0);
const totalMistakes = ref(0);

const currentCorrectness = computed(() => {
  if (props.neededCorrectCount === 0) return 0;
  return 100-Math.min(100, Math.round(( totalMistakes.value / totalPlayed.value ) * 100));
});

const progressTowardCompletion = computed(() => {
  return Math.min(
      Math.round((correct.value / props.neededCorrectCount) * 100),
      100
  );
});

const worstThree = computed(() =>
    [...props.progress]
        .sort((a, b) => b.incorrectCount - a.incorrectCount)
        .filter(p => p.incorrectCount > 0)
        .slice(0, 3)
);

watch(worstThree, (newWorst) => {
  if (newWorst && newWorst.length > 0) {  // Only emit if there are mistakes
    emit('updateWorst', newWorst);
  }
}, { immediate: true });


watch(
    () => [props.streak, props.totalPlayed, props.totalMistakes],
    ([newStreak, newPlayed, newMistakes]) => {
      streak.value = newStreak;
      totalPlayed.value = newPlayed;
      totalMistakes.value = newMistakes;
      if (newStreak !== 0 && newPlayed !== 0 && newMistakes !== 0) {
        saveProgress({ streak: newStreak, totalPlayed: newPlayed, totalMistakes: newMistakes });
      }

    },
    { immediate: true }
);

onMounted(() => {
  const saved = getProgress();
  streak.value = props.streak || saved.streak || 0;
  totalPlayed.value = props.totalPlayed || saved.totalPlayed || 0;
  totalMistakes.value = props.totalMistakes || saved.totalMistakes || 0;
});

onMounted(() => {
  const saved = getProgress();
  emit('initProgress', {
    streak: saved.streak || 0,
    totalPlayed: saved.totalPlayed || 0,
    totalMistakes: saved.totalMistakes || 0
  });
});

</script>

<style scoped>
</style>
