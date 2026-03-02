<template>
  <div class="pa-4">
      <div v-if="finished">
        <transition name="fade">
        <v-card class="paper-card text-center pa-6 mb-4">
          <h2 class="mb-2">🎉 Well Done!</h2>
          <p class="mb-4">
            You’ve completed this challenge. Time to press forward.
          </p>

          <v-alert type="success" variant="tonal" class="mb-4">
            <p class="text-greek mb-1">
              ποταμοῖσι τοῖσιν αὐτοῖσιν ἐμβαίνομέν τε καὶ οὐκ ἐμβαίνομεν, εἰμέν τε καὶ οὐκ εἰμέν.
            </p>
            <small>
              “We both step and do not step into the same rivers; we are and we are not.” — Herakleitos
            </small>
          </v-alert>

            <v-btn
                v-if="showNext"
                color="primary"
                size="large"
                @click="$emit('next')"
                class="mt-2"
            >
              Next
            </v-btn>
        </v-card>
        </transition>
      </div>
      <div v-if="finishedSegment">
        <v-card class="paper-card text-center mb-4">
          <h2 class="mb-2">🎉 Section Complete!</h2>
          <p class="mb-4">
            You’ve finished <strong>{{ section.name }}</strong>. Onward, adventurer.
          </p>

          <v-alert type="success" variant="tonal" class="mb-4">
            <p class="text-greek mb-1">
              Θάλαττα! Θάλαττα!
            </p>
            <small>
              "The sea! The sea!" — Xenophon, *Anabasis* 4.7
            </small>
          </v-alert>
          <TypingText :texts="[
        {
          greek: props.text,
          translation: props.translation
        }
      ]" />

          <transition name="fade">
            <v-btn
                v-if="showNext"
                color="primary"
                size="large"
                @click="$emit('finishSegment')"
                class="mt-4"
            >
              Return to Map
            </v-btn>
          </transition>
        </v-card>
      </div>

    <div v-if="!finished">

    <ListedQuiz
        v-if="section.__typename === 'TriviaQuiz'"
        :section="section"
        @finish="setFinishedState"
    />

    <MatchQuiz
        v-else-if="section.__typename === 'MatchQuiz'"
        :section="section"
        @finish="setFinishedState"
    />

    <ListedQuiz
        v-else-if="section.__typename === 'StructureQuiz'"
        :section="section"
        @finish="setFinishedState"
    />

    <MediaQuiz
        v-else-if="section.__typename === 'MediaQuiz'"
        :section="section"
        @finish="setFinishedState"
    />

    <FinalQuiz
        v-else-if="section.__typename === 'FinalTranslationQuiz' && !finishedSegment"
        :section="section"
        @finishSegment="finishedSegment = true"
    />

    </div>

  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import MatchQuiz from "@/components/journey/MatchQuiz.vue";
import MediaQuiz from "@/components/journey/MediaQuiz.vue";
import ListedQuiz from "@/components/journey/ListedQuiz.vue";
import FinalQuiz from "@/components/journey/FinalQuiz.vue";
import TypingText from "@/components/TypingText.vue";

const props = defineProps({
  section: Object,
  text: String,
  translation: String
})
const emit = defineEmits(['next'])
const loadedImages = ref({})
const finished = ref(false)
const finishedSegment = ref(false)
const showNext = ref(false)

async function loadImages(files) {
  for (const item of files) {
    const key = item.answer
    if (!loadedImages.value[key]) {
      try {
        const imgPath = `/src/assets/icons/${key}`;
        const imgModule = await import(/* @vite-ignore */ imgPath);
        loadedImages.value[key] = imgModule.default
      } catch (e) {
        console.warn(`Image ${key} could not be loaded`, e)
      }
    }
  }
}

function setFinishedState() {
  finished.value = true
  setTimeout(() => {
    showNext.value = true
  }, 800)
}

watch(
    () => props.section,
    () => {
      finished.value = false;
      // Also reload images if necessary for media quizzes
      if (props.section && props.section.mediaFiles) {
        loadImages(props.section.mediaFiles);
      }
    },
    { immediate: true }
)
</script>
