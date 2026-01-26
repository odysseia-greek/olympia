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

const route = useRoute();
const boule = useBouleId();

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
  <v-container class="quiz-container text-center">
    <v-card class="paper-card pa-6" elevation="4">
      <div class="text-right">
        <v-btn
            icon="mdi-minus"
            variant="text"
            @click="minimized = !minimized"
        >
          <v-icon>{{ minimized ? 'mdi-plus' : 'mdi-minus' }}</v-icon>
        </v-btn>
      </div>

      <div v-if="!minimized" ref="selectThemeRef">
        <v-card-title class="text-h5">Dialogue Quiz</v-card-title>
        <!-- Thematic quote always visible -->
        <p class="ma-4">
          Ἀρχὴ πάσης πράξεως ἐστὶν ἡ τοῦ αἱρεῖσθαι ἀρχή.
          <v-divider class="my-4" />
          The beginning of every action is the choice.
        </p>
          <!-- Theme -->
          <v-combobox
              class="mt-5"
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
        <div v-if="maxSet > 1">
          <v-row class="mt-4">
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
      v-if="theme"
      :dialogueOptions="dialogueOptions"
      :dialogueContent="dialogueContent"
      :selectedTheme="theme"
      :selectedSet="set"
      :boule="boule"
  />
  </v-container>
</template>
