<template>
  <v-container v-if="dialogueOptions" class="dialogue-stage">
    <!-- Introduction -->
    <v-card class="paper-card dialogue-intro-card mb-4" elevation="4">
      <div class="intro-grid">
        <div>
          <v-card-title class="intro-title">Set The Scene</v-card-title>
          <v-card-text class="intro-text">{{ dialogueOptions.introduction }}</v-card-text>
        </div>

        <div class="text-facts">
          <span>Section</span>
          <strong>{{ dialogueOptions.section }}</strong>
          <a :href="dialogueOptions.linkToPerseus" target="_blank" rel="noopener">
            Open Perseus
            <v-icon size="16">mdi-open-in-new</v-icon>
          </a>
        </div>
      </div>

      <div class="role-panel">
        <div class="role-copy">
          <span>Step 1</span>
          <strong>Choose your speaker</strong>
          <small>Your lines become selectable responses; the other speaker is typed in automatically.</small>
        </div>
        <v-row class="role-grid">
          <v-col v-for="(speaker, index) in dialogueOptions.speakers" :key="index" cols="12" sm="6">
            <v-btn
                class="role-button"
                :color="selectedSpeaker === speaker.shorthand ? 'primary' : 'triadic'"
                :variant="selectedSpeaker === speaker.shorthand ? 'flat' : 'tonal'"
                block
                rounded="xl"
                @click="setSpeaker(speaker);scrollMeTo('dialogueRef');"
            >
              <span class="role-shorthand">{{ speaker.shorthand }}</span>
              <span>{{ speaker.translation }}</span>
            </v-btn>
          </v-col>
        </v-row>
      </div>
    </v-card>
    <v-card v-if="selectedSpeaker !== ''" class="mb-3 paper-card dialogue-play-card">
      <div class="dialogue-card-header" ref="dialogueRef">
        <div>
          <v-card-title class="dialogue-title">Dialogue</v-card-title>
          <p class="dialogue-help">
            Select responses below, then reorder only your own lines before checking the conversation.
          </p>
        </div>
        <div class="dialogue-controls">
          <v-chip color="secondary" variant="flat" size="small">Playing {{ selectedSpeaker }}</v-chip>
          <v-switch
              v-model="showDialogueTranslation"
              color="primary"
              label="Translation"
              density="compact"
              hide-details
          ></v-switch>
        </div>
      </div>
      <v-card-text>
        <div class="dialogue-thread">
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
              <span class="speaker-pill">{{ line.speaker }}</span>
              <span class="greek-line">{{ line.greek }}</span>
              <div v-if="line.speaker === selectedSpeaker" class="move-buttons">
                <v-btn icon="mdi-chevron-up" size="x-small" variant="text" @click="moveBubbleUp(index)">
                </v-btn>
                <v-btn icon="mdi-chevron-down" size="x-small" variant="text" @click="moveBubbleDown(index)">
                </v-btn>
              </div>
            </div>
            <div v-if="showDialogueTranslation && !hideTranslation" class="translation-text">
              <strong>{{ line.speaker }}:</strong>
              {{ line.translation }}
            </div>
          </div>
        </div>
        <div class="response-dock" ref="responseDockRef">
          <div class="response-dock-header">
            <div>
              <v-card-title class="response-title">
                {{ responseOptions.length > 0 ? 'Choose Your Next Line' : 'Ready To Check' }}
              </v-card-title>
              <p class="response-help">
                {{ responseOptions.length > 0 ? 'Pick the reply that should come next. This list shrinks as the conversation grows.' : 'All of your responses have been placed. Review the order above.' }}
              </p>
            </div>
            <v-chip color="triadic" variant="flat" size="small">
              {{ responseOptions.length }} left
            </v-chip>
          </div>
          <v-row v-if="responseOptions.length > 0" class="response-grid">
            <v-col v-for="(response, index) in responseOptions" :key="index" cols="12" md="6">
              <v-card class="response-card" @click="setDialogue(response)">
                <v-card-text>
                  <span class="response-greek">{{ response.greek }}</span>
                  <span v-if="showDialogueTranslation && !hideTranslation" class="translation-text"> <em><br>{{ response.translation }}</em></span>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
          <div class="dialogue-check-row">
            <v-btn
                v-if="!hideTranslation"
                class="check-button"
                color="primary"
                prepend-icon="mdi-check-circle-outline"
                rounded="xl"
                @click="checkDialogueAnswer();"
            >
              Check Order
            </v-btn>
            <span v-if="wronglyPlaced.length > 0" class="check-hint">
              Move the marked lines and check again.
            </span>
          </div>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script>
import { ref, nextTick } from 'vue';
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
    const dialogueRef = ref(null);
    const responseDockRef = ref(null);

    const setSpeaker = (value) => {
      selectedSpeaker.value = value.shorthand;
      dialogueText.value = [];
      wronglyPlaced.value = [];
      scrollToLatestLine();
      initializeDialogue();
    };

    const initializeDialogue = async () => {
      if (!Array.isArray(props.dialogueContent) || props.dialogueContent.length === 0) {
        return;
      }

      const firstSpeaker = props.dialogueContent[0].speaker;

      if (selectedSpeaker.value === firstSpeaker) {
        dialogueText.value = [props.dialogueContent[0], props.dialogueContent[1]];
      } else {
        dialogueText.value = [props.dialogueContent[0]];
      }

      const responses = props.dialogueContent.filter((line) => line.speaker === selectedSpeaker.value && line.place !== 1);

      responseOptions.value = await createNewArray(responses);
      scrollToResponseDock();
    };

    const setDialogue = (selectedDialogue) => {
      let setTranslationBack = false;
      if (showDialogueTranslation.value) {
        hideTranslation.value = true;
        setTranslationBack = true;
      }

      updateResponseOptions(selectedDialogue);

      const restoreTranslationAndFocusResponses = () => {
        const finish = () => {
          scrollToResponseDock();
        };

        if (setTranslationBack) {
          setTimeout(() => {
            hideTranslation.value = false;
            finish();
          }, 500);
          return;
        }

        finish();
      };

      typeDialogue(selectedDialogue, () => {
        const nextIndex = dialogueText.value.length;
        if (nextIndex < props.dialogueContent.length && props.dialogueContent[nextIndex].speaker !== selectedSpeaker.value) {
          setTimeout(() => {
            typeDialogue(props.dialogueContent[nextIndex], restoreTranslationAndFocusResponses);
          }, 500);
        } else {
          restoreTranslationAndFocusResponses();
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
      scrollToLatestLine();

      const index = dialogueText.value.length - 1;

      for (let i = 0; i < line.greek.length; i++) {
        setTimeout(() => {
          typedText += line.greek[i];
          dialogueText.value[index].greek = typedText; // Directly modify the array element
          if (i === line.greek.length - 1 && callback) {
            scrollToLatestLine();
            callback();
          }
        }, i * typingSpeed);
      }
    };


    const scrollToLatestLine = () => {
      nextTick(() => {
        const targetIndex = Math.max(dialogueText.value.length - 1, 0);
        const targetElement = dialogueOptionsRefs.value[targetIndex];

        if (targetElement) {
          targetElement.scrollIntoView({ behavior: 'smooth', block: 'end' });
        }
      });
    };

    const scrollToResponseDock = () => {
      nextTick(() => {
        const target = responseDockRef.value?.$el || responseDockRef.value;
        target?.scrollIntoView?.({ behavior: 'smooth', block: 'center' });
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

      wronglyPlaced.value = [];
      dialogueText.value = dialogueText.value.map((text) => {
        return {
          ...text,
          isWronglyPlaced: false,
          isCorrectlyPlaced: false,
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
        if (refName === 'dialogueRef') {
          const target = dialogueRef.value?.$el || dialogueRef.value;
          target?.scrollIntoView?.({ behavior: 'smooth', block: 'start' });
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
      dialogueRef,
      responseDockRef,
      setSpeaker,
      setDialogue,
      initializeDialogue,
      scrollToLatestLine,
      scrollToResponseDock,
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
.dialogue-stage {
  max-width: 1180px;
  padding-inline: 0;
}

.paper-card {
  border: 1px solid rgba(119, 86, 37, 0.14);
  border-radius: 26px;
  background:
      radial-gradient(circle at top right, rgba(28, 188, 209, 0.1), transparent 30%),
      linear-gradient(145deg, #fff9e9 0%, #eedfb7 100%);
  box-shadow: 0 20px 42px rgba(54, 43, 23, 0.12);
  color: #283650;
  padding: 20px;
}

.dialogue-intro-card {
  overflow: hidden;
}

.intro-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 220px;
  gap: 18px;
  align-items: stretch;
}

.intro-title,
.dialogue-title,
.response-title {
  color: #17345f;
  font-family: Georgia, "Times New Roman", serif;
  font-weight: 900;
  letter-spacing: -0.03em;
}

.intro-text {
  color: #41506b;
  font-size: 1rem;
  line-height: 1.65;
}

.text-facts {
  display: grid;
  align-content: center;
  gap: 8px;
  padding: 18px;
  border: 1px solid rgba(28, 97, 209, 0.16);
  border-radius: 22px;
  background: rgba(255, 252, 242, 0.72);
  text-align: left;
}

.text-facts span {
  color: #69768d;
  font-size: 0.78rem;
  font-weight: 900;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.text-facts strong {
  color: #17345f;
  font-size: 1.25rem;
}

.text-facts a {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  width: fit-content;
  color: #1c61d1;
  font-weight: 900;
  text-decoration: none;
}

.role-panel {
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr);
  gap: 18px;
  margin-top: 22px;
  padding: 18px;
  border-radius: 24px;
  background: rgba(28, 97, 209, 0.08);
}

.role-copy {
  display: grid;
  gap: 5px;
  align-content: center;
  text-align: left;
}

.role-copy span {
  color: #1c61d1;
  font-size: 0.78rem;
  font-weight: 900;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.role-copy strong {
  color: #17345f;
  font-size: 1.3rem;
}

.role-copy small {
  color: #46546c;
  line-height: 1.5;
}

.role-grid {
  align-items: center;
}

.role-button {
  min-height: 56px;
  justify-content: flex-start;
  text-transform: none;
}

.role-shorthand {
  margin-right: 10px;
  font-family: Georgia, "Times New Roman", serif;
  font-size: 1.2rem;
  font-weight: 900;
}

.dialogue-play-card {
  background:
      linear-gradient(160deg, rgba(255, 249, 233, 0.96), rgba(240, 224, 188, 0.96)),
      radial-gradient(circle at top left, rgba(28, 209, 140, 0.14), transparent 30%);
}

.dialogue-card-header {
  display: flex;
  gap: 18px;
  align-items: flex-start;
  justify-content: space-between;
  padding: 4px 4px 0;
  scroll-margin-top: 84px;
}

.dialogue-help,
.response-help {
  margin: 0 16px;
  color: #566177;
  font-size: 0.95rem;
  line-height: 1.45;
}

.dialogue-controls {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
  justify-content: flex-end;
  min-width: 230px;
}

.dialogue-thread {
  display: grid;
  gap: 10px;
  padding: 18px 10px 8px;
  scroll-behavior: smooth;
}

.dialogue-line {
  transition: all 0.3s ease;
}

.dialogue-bubble {
  position: relative;
  display: grid;
  gap: 8px;
  width: fit-content;
  max-width: 80%;
  margin-bottom: 8px;
  padding: 15px 50px 15px 16px;
  border: 1px solid rgba(28, 97, 209, 0.14);
  border-radius: 18px;
  background-color: #fffdfa;
  box-shadow: 0 12px 24px rgba(42, 32, 16, 0.1);
  text-align: left;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.dialogue-bubble:hover {
  transform: translateY(-2px);
  box-shadow: 0 16px 30px rgba(42, 32, 16, 0.16);
}

.speaker-pill {
  width: fit-content;
  padding: 3px 9px;
  border-radius: 999px;
  background: rgba(28, 97, 209, 0.1);
  color: #1c61d1;
  font-family: Georgia, "Times New Roman", serif;
  font-size: 0.82rem;
  font-weight: 900;
}

.greek-line,
.response-greek {
  color: #263149;
  font-family: Georgia, "Times New Roman", serif;
  font-size: 1.08rem;
  font-weight: 800;
  line-height: 1.55;
}

.move-buttons {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  gap: 2px;
  opacity: 0.22;
  transition: opacity 0.2s ease;
}

.dialogue-bubble:hover .move-buttons,
.dialogue-bubble:focus-within .move-buttons {
  opacity: 1;
}

@keyframes flash-green {
  0% { border-color: #1cd18c; box-shadow: 0 0 0 0 rgba(28, 209, 140, 0.32); }
  50% { border-color: rgba(28, 209, 140, 0.26); box-shadow: 0 0 0 8px rgba(28, 209, 140, 0); }
  100% { border-color: #1cd18c; box-shadow: 0 0 0 0 rgba(28, 209, 140, 0.32); }
}

.correctly-placed .dialogue-bubble {
  animation: flash-green 1.8s ease-in-out infinite;
  border-color: #1cd18c;
}

.wrongly-placed .dialogue-bubble {
  border-color: rgba(202, 56, 40, 0.72);
  background: #ffe8df;
}

.user-speaker .dialogue-bubble {
  margin-left: auto;
  border-top-right-radius: 4px;
  background:
      linear-gradient(145deg, rgba(28, 97, 209, 0.12), rgba(28, 188, 209, 0.1)),
      #fffdfa;
}

.other-speaker .dialogue-bubble {
  margin-right: auto;
  border-top-left-radius: 4px;
}

.translation-text {
  max-width: 76%;
  color: #5c6680;
  font-style: italic;
  margin-top: 4px;
  text-align: left;
}

.user-speaker .translation-text {
  margin-left: auto;
  text-align: right;
}

.dialogue-check-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  justify-content: flex-end;
  margin-top: 14px;
}

.check-button {
  min-width: 150px;
  text-transform: none;
}

.check-hint {
  color: #a2432c;
  font-size: 0.92rem;
  font-weight: 800;
}

.response-dock {
  position: relative;
  margin: 18px 6px 4px;
  padding: 18px;
  border: 1px solid rgba(28, 97, 209, 0.16);
  border-radius: 24px;
  background:
      radial-gradient(circle at top right, rgba(28, 209, 140, 0.18), transparent 30%),
      linear-gradient(145deg, rgba(255, 253, 250, 0.95), rgba(242, 229, 197, 0.95));
  box-shadow: 0 18px 34px rgba(42, 32, 16, 0.12);
  scroll-margin-block: 120px;
}

.response-dock::before {
  position: absolute;
  top: -10px;
  left: 50%;
  width: 18px;
  height: 18px;
  border-top: 1px solid rgba(28, 97, 209, 0.16);
  border-left: 1px solid rgba(28, 97, 209, 0.16);
  background: rgba(255, 253, 250, 0.95);
  content: "";
  transform: translateX(-50%) rotate(45deg);
}

.response-dock-header {
  position: relative;
  z-index: 1;
  display: flex;
  gap: 16px;
  align-items: flex-start;
  justify-content: space-between;
}

.response-grid {
  position: relative;
  z-index: 1;
  margin-top: 8px;
}

.response-card {
  height: 100%;
  border: 1px solid rgba(28, 97, 209, 0.12);
  border-radius: 18px;
  background-color: rgba(255, 253, 250, 0.88);
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.response-card:hover {
  border-color: rgba(28, 97, 209, 0.36);
  box-shadow: 0 16px 28px rgba(42, 32, 16, 0.14);
  transform: translateY(-3px);
}

@media (max-width: 760px) {
  .intro-grid,
  .role-panel {
    grid-template-columns: 1fr;
  }

  .dialogue-card-header {
    display: grid;
  }

  .dialogue-controls {
    justify-content: flex-start;
    min-width: 0;
  }

  .response-dock-header {
    display: grid;
  }

  .dialogue-bubble,
  .translation-text {
    max-width: 92%;
  }
}
</style>
