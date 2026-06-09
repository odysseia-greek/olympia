<script setup>
import {ref, watch, onMounted, reactive, nextTick, getCurrentInstance, computed} from 'vue';
import { useQuery } from '@vue/apollo-composable';
import {apolloClient} from "@/apollo";
import { useRoute } from 'vue-router';
import {useBouleId} from '@/composables/useBoule';
import { updateQuizUrl } from '@/utils/sharedQuiz.js';
import {DialogueBasedOptions, DialogueBasedQuestion} from "@/constants/dialogueBasedGraphql";
import Dialogue from "@/components/Dialogue.vue";

const { proxy } = getCurrentInstance();
const theme = ref('');
const set = ref(1);
const maxSet = ref(1);

const themes = ref([]);
const minimized = ref(false);
const dialogueOptions = ref({});
const dialogueContent = ref({});
const selectThemeRef = ref(null);

const route = useRoute();
const boule = useBouleId();
const hasDialogueContent = computed(() => Array.isArray(dialogueContent.value) && dialogueContent.value.length > 0);

const { result: optionsResult, loading, onResult } = useQuery(DialogueBasedOptions);

onResult(({ data }) => {
  if (data && data.dialogueOptions) {
    themes.value = data.dialogueOptions.themes;
    initializeFromRoute();
  }
});


watch([theme], ([newTheme]) => {
  if (newTheme) {
    getDialogueQuiz();
    updateQuizUrl(
        proxy.$router,
        proxy.$route.query,
        'QuizDialogue',
        {
          theme: theme.value,
          set: String(set.value),
        }
    );
  }
});

watch([set], ([newSet]) => {
  if (newSet) {
    getDialogueQuiz()
    updateQuizUrl(
        proxy.$router,
        proxy.$route.query,
        'QuizDialogue', // or 'QuizMedia'
        { theme: theme.value, set: newSet }
    );
  }
})

// Handle theme selection
const onThemeChange = (selected) => {
  const themeData = themes.value.find((t) => t.name.toLowerCase() === selected.toLowerCase());
  if (themeData) {
    theme.value = themeData.name;
    maxSet.value = themeData.maxSet;
    set.value = 1;
    dialogueOptions.value = {}
    dialogueContent.value = {}
  }

  updateQuizUrl(
      proxy.$router,
      proxy.$route.query,
      'QuizDialogue',
      {
        theme: theme.value,
        set: String(set.value),
      }
  );
}

// Fetch quiz from backend
const getDialogueQuiz = async () => {
  try {
    const {data} = await apolloClient.query({
      query: DialogueBasedQuestion,
      variables: {
        input: {
          theme: theme.value,
          set: String(set.value),
        },
      },
      context: {
        headers: {
          'boule': boule,
        },
      },
      fetchPolicy: 'no-cache',
    });

    const result = data.dialogueQuiz;
    dialogueOptions.value = result.dialogue;
    dialogueContent.value = result.content;
  }
  catch (error) {
    console.error(error);
  }
};

const decrementSelectedSet = () => {
  if (set.value > 1) {
    set.value--;
  }
};
const incrementSelectedSet = () => {
  if (set.value < maxSet.value) {
    set.value++;
  }
};


const initializeFromRoute = () => {
  const { theme: qTheme, set: qSet } = route.query;
  if (qTheme) {
    onThemeChange(qTheme);
  }

  if (qSet) {
    set.value = parseInt(qSet, 10);
  }
};

</script>

<template>
  <v-container class="quiz-container dialogue-container text-center">
    <v-card class="paper-card dialogue-setup-card pa-6" elevation="4">
      <div class="dialogue-toolbar">
        <div class="dialogue-status">
          <v-chip
              v-if="theme"
              class="status-chip"
              color="primary"
              variant="flat"
              size="small"
          >
            {{ theme }}
          </v-chip>
          <v-chip
              v-if="theme && maxSet > 1"
              class="status-chip"
              color="secondary"
              variant="flat"
              size="small"
          >
            Set {{ set }} / {{ maxSet }}
          </v-chip>
        </div>
        <v-btn
            icon="mdi-minus"
            variant="text"
            color="primary"
            @click="minimized = !minimized"
        >
          <v-icon>{{ minimized ? 'mdi-plus' : 'mdi-minus' }}</v-icon>
        </v-btn>
      </div>

      <div v-if="!minimized" ref="selectThemeRef">
        <v-card-title class="text-h5 dialogue-title">Dialogue Quiz</v-card-title>
        <p class="dialogue-quote">
          <span>Ἀρχὴ πάσης πράξεως ἐστὶν ἡ τοῦ αἱρεῖσθαι ἀρχή.</span>
          <small>The beginning of every action is the choice.</small>
        </p>
        <p class="dialogue-copy">
          Pick a theme, choose the role you want to play, then build the conversation by selecting and ordering your lines.
        </p>
          <!-- Theme -->
          <v-combobox
              class="dialogue-select mt-5"
              :class="{ 'pulsate': theme === '' }"
              v-model="theme"
              :items="themes.map(t => t.name)"
              item-title="name"
              item-value="name"
              label="Select a Theme"
              color="primary"
              @update:modelValue="onThemeChange"
              style="margin-top: 2em"
          />
          <!-- Toggles -->
        <div v-if="maxSet > 1" class="dialogue-set-panel">
          <v-row class="mt-4 align-center">
            <v-col class="text-left">
                          <span class="subheading font-weight-light me-1"
                          >Set
                          </span>
              <span
                  class="text-h4 font-weight-light"
                  v-text="set"
              ></span>
              <span class="subheading font-weight-light me-1">
                            of
                          </span>
              <span
                  class="text-h4 font-weight-light"
                  v-text="maxSet"
              ></span>
            </v-col>
          </v-row>
          <v-slider
              class="my-5"
              v-model="set"
              :max="maxSet"
              :min="1"
              step="1"
              color="primary"
              track-color="accent"
              thumb-color="primary"
          >
            <template v-slot:prepend>
              <v-btn
                  icon="mdi-minus"
                  size="small"
                  variant="text"
                  @click="decrementSelectedSet"
              ></v-btn>
            </template>

            <template v-slot:append>
              <v-btn
                  icon="mdi-plus"
                  size="small"
                  variant="text"
                  @click="incrementSelectedSet"
              ></v-btn>
            </template>
          </v-slider>
        </div>
      </div>
    </v-card>
  <Dialogue
      v-if="hasDialogueContent"
      :dialogueOptions="dialogueOptions"
      :dialogueContent="dialogueContent"
      :selectedTheme="theme"
      :selectedSet="set"
      :boule="boule"
  />
    <v-card v-else-if="theme" class="paper-card dialogue-loading-card" elevation="2">
      <v-progress-circular indeterminate color="primary" size="28" />
      <span>Preparing dialogue...</span>
    </v-card>
  </v-container>
</template>

<style scoped>
.dialogue-container {
  max-width: 1180px;
}

.dialogue-setup-card {
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(28, 97, 209, 0.16);
  border-radius: 28px;
  background:
      radial-gradient(circle at top left, rgba(28, 209, 140, 0.18), transparent 32%),
      linear-gradient(145deg, #fff7df 0%, #f1e1b7 100%);
  color: #263149;
}

.dialogue-setup-card::before {
  position: absolute;
  inset: 0;
  pointer-events: none;
  content: "";
  background:
      linear-gradient(120deg, rgba(28, 97, 209, 0.08), transparent 35%),
      repeating-linear-gradient(90deg, rgba(128, 90, 35, 0.04) 0 1px, transparent 1px 18px);
}

.dialogue-toolbar {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 36px;
}

.dialogue-status {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.status-chip {
  color: #fff;
  font-weight: 800;
  letter-spacing: 0.02em;
}

.dialogue-title {
  position: relative;
  z-index: 1;
  justify-content: center;
  color: #17345f;
  font-family: Georgia, "Times New Roman", serif;
  font-size: clamp(2rem, 5vw, 4rem) !important;
  font-weight: 900;
  letter-spacing: -0.04em;
}

.dialogue-quote {
  position: relative;
  z-index: 1;
  display: grid;
  gap: 8px;
  max-width: 720px;
  margin: 16px auto 0;
  color: #42506b;
}

.dialogue-quote span {
  color: #1c61d1;
  font-family: Georgia, "Times New Roman", serif;
  font-size: clamp(1.2rem, 3vw, 1.85rem);
  font-weight: 800;
}

.dialogue-quote small {
  color: #69593a;
  font-size: 0.98rem;
  font-weight: 700;
}

.dialogue-copy {
  position: relative;
  z-index: 1;
  max-width: 660px;
  margin: 18px auto 0;
  color: #36435c;
  font-size: 1.02rem;
}

.dialogue-select,
.dialogue-set-panel {
  position: relative;
  z-index: 1;
}

.dialogue-set-panel {
  margin-top: 18px;
  padding: 16px 18px 4px;
  border: 1px solid rgba(28, 188, 209, 0.2);
  border-radius: 20px;
  background: rgba(255, 252, 242, 0.72);
}

.dialogue-loading-card {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  margin-top: 20px;
  padding: 18px 22px;
  border-radius: 999px;
  background: #fff8e7;
  color: #17345f;
  font-weight: 800;
}

@media (max-width: 680px) {
  .dialogue-setup-card {
    border-radius: 22px;
    padding: 22px !important;
  }

  .dialogue-toolbar {
    align-items: flex-start;
  }
}
</style>
