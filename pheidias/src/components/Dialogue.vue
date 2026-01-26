<template>
  <v-container v-if="dialogueOptions">
    <!-- Introduction -->
    <v-card class="paper-card mb-4" elevation="4">
      <v-card-title>Introduction</v-card-title>
      <v-card-text>{{ dialogueOptions.introduction }}</v-card-text>

      <v-card-title>Text Information</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6">
            <strong>Section:</strong> {{ dialogueOptions.section }}
          </v-col>
          <v-col cols="12" sm="6">
            <strong>Perseus:</strong>
            <a :href="dialogueOptions.linkToPerseus" target="_blank" rel="noopener">
              Link
            </a>
          </v-col>
        </v-row>
      </v-card-text>

      <v-card-title>What role do you want to play?</v-card-title>
      <v-card-text>
        <v-row>
          <v-col v-for="(speaker, index) in dialogueOptions.speakers" :key="index" cols="12" sm="6">
            <v-btn
                class="ma-1"
                :color="selectedSpeaker === speaker.shorthand ? 'primary' : 'triadic'"
                block
                @click="setSpeaker(speaker);scrollMeTo('dialogueRef');"
            >
              {{ speaker.shorthand }} {{ speaker.translation }}
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
    <v-card v-if="selectedSpeaker !== ''" class="mb-3 paper-card">
      <v-card-title ref="dialogueRef">Dialogue</v-card-title>
      <v-switch
          v-model="showDialogueTranslation"
          color="primary"
          label="Show Translation"
      ></v-switch>
      <v-card-text>
        <div
            v-for="(line, index) in dialogueText"
            :key="index"
            class="dialogue-line"
            :class="{
              'user-speaker': line.speaker === selectedSpeaker,
              'other-speaker': line.speaker !== selectedSpeaker,
              'wrongly-placed': line.isWronglyPlaced,
              'correctly-placed': line.isCorrectlyPlaced
            }"
        >
          <div class="dialogue-bubble" :ref="el => setDialogueOptionRef(el, index)">
            <strong>{{ line.speaker }}:</strong>
            {{ line.greek }}
            <div v-if="line.speaker === selectedSpeaker" class="move-buttons">
              <v-btn icon="mdi-chevron-up" small @click="moveBubbleUp(index)">
              </v-btn>
              <v-btn icon="mdi-chevron-down" small @click="moveBubbleDown(index)">
              </v-btn>
            </div>
          </div>
          <div v-if="showDialogueTranslation && !hideTranslation" class="translation-text">
            <strong>{{ line.speaker }}:</strong>
            {{ line.translation }}
          </div>
        </div>
        <v-btn v-if="!hideTranslation" @click="checkDialogueAnswer();" color="primary">Check</v-btn>
      </v-card-text>
    </v-card>
    <v-card  class="mb-3 paper-card">
      <v-card-text>
      <v-card-title v-if="responseOptions.length > 0">Possible Responses</v-card-title>
        <v-row>
          <v-col v-for="(response, index) in responseOptions" :key="index" cols="12" sm="6">
            <v-card class="ma-2 response-card" @click="setDialogue(response)">
              <v-card-text>
                {{ response.greek }}
                <span v-if="showDialogueTranslation && !hideTranslation" class="translation-text"> <em><br>{{ response.translation }}</em></span>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script>
import { ref, watch, nextTick } from 'vue';
import {DialogueBasedAnswer} from "@/constants/dialogueBasedGraphql";
import {apolloClient} from "@/apollo";

export default {
  name: 'Dialogue',
  props: {
    dialogueOptions: Object,
    dialogueContent: Object,
    selectedTheme: String,
    selectedSet: Number,
    boule: String,
  },
  setup(props) {
    const selectedSpeaker = ref('');
    const dialogueText = ref([]);
    const responseOptions = ref([]);
    const dialogueOptionsRefs = ref([]);
    const wronglyPlaced = ref([]);
    const showDialogueTranslation = ref(true);
    const hideTranslation = ref(false);

    const setSpeaker = (value) => {
      selectedSpeaker.value = value.shorthand;
      dialogueText.value = [];
      scrollToLatestLine();
      initializeDialogue();
    };

    const initializeDialogue = async () => {
      const firstSpeaker = props.dialogueContent[0].speaker;

      if (selectedSpeaker.value === firstSpeaker) {
        dialogueText.value = [props.dialogueContent[0], props.dialogueContent[1]];
      } else {
        dialogueText.value = [props.dialogueContent[0]];
      }

      const responses = props.dialogueContent.filter((line) => line.speaker === selectedSpeaker.value && line.place !== 1);

      responseOptions.value = await createNewArray(responses);
    };

    const setDialogue = (selectedDialogue) => {
      let setTranslationBack = false;
      if (showDialogueTranslation.value) {
        hideTranslation.value = true;
        setTranslationBack = true;
      }

      updateResponseOptions(selectedDialogue);
      scrollToLatestLine();

      typeDialogue(selectedDialogue, () => {
        const nextIndex = dialogueText.value.length;
        if (nextIndex < props.dialogueContent.length && props.dialogueContent[nextIndex].speaker !== selectedSpeaker.value) {
          setTimeout(() => {
            typeDialogue(props.dialogueContent[nextIndex], () => {
              if (setTranslationBack) {
                setTimeout(() => {
                  hideTranslation.value = false
                }, 500);
              }
            });
          }, 500);
        } else {
          if (setTranslationBack) {
            setTimeout(() => {
              hideTranslation.value = false
            }, 500);
          }
        }
      });
    };

    const updateResponseOptions = (selectedDialogue) => {
      responseOptions.value = responseOptions.value.filter((option) => option.place !== selectedDialogue.place);
    };

    const typeDialogue = (line, callback = null) => {
      let typedText = '';
      const typingSpeed = 50;
      const newLine = { ...line, greek: '' };
      dialogueText.value.push(newLine);

      const index = dialogueText.value.length - 1;

      for (let i = 0; i < line.greek.length; i++) {
        setTimeout(() => {
          typedText += line.greek[i];
          dialogueText.value[index].greek = typedText; // Directly modify the array element
          if (i === line.greek.length - 1 && callback) {
            callback();
          }
        }, i * typingSpeed);
      }
    };


    const scrollToLatestLine = () => {
      nextTick(() => {
        const targetIndex = Math.max(dialogueText.value.length - 5, 0);
        const targetElement = dialogueOptionsRefs.value[targetIndex];

        if (targetElement) {
          targetElement.scrollIntoView({ behavior: 'smooth' });
        }
      });
    };

    const setDialogueOptionRef = (el, index) => {
      if (el) {
        dialogueOptionsRefs.value[index] = el;
      } else {
        dialogueOptionsRefs.value.splice(index, 1);
      }
    };

    const checkDialogueAnswer = async () => {
      const dialogueData = dialogueText.value.map(({ isWronglyPlaced, isCorrectlyPlaced, __typename, ...rest }, index) => {
        return { ...rest, place: index + 1 };
      });

      dialogueText.value = dialogueText.value.map((text) => {
        return {
          ...text,
        };
      });

      const {data} = await apolloClient.query({
        query: DialogueBasedAnswer,
        variables: {
          input: {
            theme: props.selectedTheme,
            set: String(props.selectedSet),
            content: dialogueData,
          },
        },
        context: {
          headers: {
            'boule': props.boule,
          },
        },
        fetchPolicy: 'no-cache',
      });

      const result = data.dialogueAnswer;

      if (result.wronglyPlaced) {
        wronglyPlaced.value = result.wronglyPlaced;
        wronglyPlaced.value.forEach((wrongItem) => {
          const index = wrongItem.place - 1;
          if (dialogueText.value[index]) {
            if (dialogueText.value[index].speaker === selectedSpeaker.value) {
              dialogueText.value[index].isWronglyPlaced = true;
            }
          }
        });
      }

      // Set isCorrectlyPlaced to true for correctly placed items
      dialogueText.value.forEach((text, index) => {
        if (!wronglyPlaced.value.some((wrongItem) => wrongItem.place - 1 === index)) {
          if (dialogueText.value[index].speaker === selectedSpeaker.value) {
            dialogueText.value[index].isCorrectlyPlaced = true;
            setTimeout(() => {
              dialogueText.value[index].isCorrectlyPlaced = false;
            }, 5000); // Remove the flashing effect after 5 seconds
          }
        }
      });
    }

    const createNewArray = async (shuffledArray) => {
      for (let i = shuffledArray.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [shuffledArray[i], shuffledArray[j]] = [shuffledArray[j], shuffledArray[i]];
      }

      return shuffledArray;
    };

    const scrollMeTo = (refName) => {
      nextTick(() => {
        if (refName === 'quiz') {
          quizContainerRef.value.scrollIntoView({ behavior: 'smooth' });
        }
      });
    };

    const moveBubbleUp = (index) => {
      if (index > 0) {
        let newIndex = index - 1;
        // Skip non-speaker bubbles
        while (newIndex > 0 && dialogueText.value[newIndex].speaker !== selectedSpeaker.value) {
          newIndex--;
        }
        if (dialogueText.value[newIndex].speaker === selectedSpeaker.value) {
            if (dialogueText.value[newIndex].isWronglyPlaced && dialogueText.value[newIndex].isWronglyPlaced !== '') {
              dialogueText.value[newIndex].isWronglyPlaced = false;
             }

          if (dialogueText.value[index].isWronglyPlaced && dialogueText.value[index].isWronglyPlaced !== '') {
            dialogueText.value[index].isWronglyPlaced = false;
          }
          const temp = dialogueText.value[index];
          dialogueText.value[index] = dialogueText.value[newIndex];
          dialogueText.value[newIndex] = temp;
        }
      }
    };

    const moveBubbleDown = (index) => {
      if (index < dialogueText.value.length - 1) {
        let newIndex = index + 1;
        // Skip non-speaker bubbles
        while (newIndex < dialogueText.value.length - 1 && dialogueText.value[newIndex].speaker !== selectedSpeaker.value) {
          newIndex++;
        }
        if (dialogueText.value[newIndex].speaker === selectedSpeaker.value) {
          if (dialogueText.value[newIndex].isWronglyPlaced && dialogueText.value[newIndex].isWronglyPlaced !== '') {
            dialogueText.value[newIndex].isWronglyPlaced = false;
          }

          if (dialogueText.value[index].isWronglyPlaced && dialogueText.value[index].isWronglyPlaced !== '') {
            dialogueText.value[index].isWronglyPlaced = false;
          }
          const temp = dialogueText.value[index];
          dialogueText.value[index] = dialogueText.value[newIndex];
          dialogueText.value[newIndex] = temp;
        }
      }
    };

    return {
      selectedSpeaker,
      dialogueText,
      responseOptions,
      dialogueOptionsRefs,
      wronglyPlaced,
      showDialogueTranslation,
      hideTranslation,
      setSpeaker,
      setDialogue,
      initializeDialogue,
      scrollToLatestLine,
      setDialogueOptionRef,
      checkDialogueAnswer,
      scrollMeTo,
      moveBubbleUp,
      moveBubbleDown,
    };
  },
};
</script>

<style scoped>
.paper-card {
  background: #fdf6e3; /* A light, papyrus-like color */
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
  font-family: 'Roboto', serif;
}

.dialogue-line {
  margin-bottom: 10px;
  transition: all 0.3s ease;
}

.response-card {
  background-color: #fefcf5;
}

.dialogue-bubble {
  background-color: #fefcf5;
  padding: 10px;
  border-radius: 10px;
  margin-bottom: 10px;
  max-width: 80%;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* Add a subtle shadow */
  border: 1px solid #e0e0e0; /* Add a light border for contrast */
  transition: transform 0.2s, box-shadow 0.2s; /* Add transitions for lift effect */
}

.dialogue-bubble:hover {
  transform: translateY(-4px); /* Lift the bubble on hover */
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2); /* Enhance the shadow on hover */
}

.dialogue-bubble:hover .move-buttons {
  display: flex;
  flex-direction: row;
}

.move-buttons {
  display: none;
  position: absolute;
  top: 5px;
  right: 5px;
}

@keyframes flash-green {
  0% { border-color: green; }
  50% { border-color: transparent; }
  100% { border-color: green; }
}

.correctly-placed .dialogue-bubble {
  animation: flash-green 2s infinite;
  border-width: 4px; /* Increase the border width */
  border-style: solid; /* Ensure the border style is solid */
  padding: 10px; /* Optional: add padding to make space for the border */
  box-sizing: border-box; /* Ensure the border is included in the element's dimensions */
}

.wrongly-placed .dialogue-bubble {
  background-color: #f6624c; /* Light red background */
}

.user-speaker .dialogue-bubble {
  margin-left: auto; /* Aligns to the right */
  border-top-right-radius: 0; /* Optional: for speech bubble effect */
}

.other-speaker .dialogue-bubble {
  margin-right: auto; /* Aligns to the left */
  border-top-left-radius: 0; /* Optional: for speech bubble effect */
}

.translation-text {
  color: #666;
  font-style: italic;
  margin-top: 5px;
}
</style>
