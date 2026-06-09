<template>
  <div id="typing-text">
    <div class="typing-body">
      <span>{{ displayedText }}</span>
      <p v-if="displayedTranslation"></p>
      <span v-if="displayedTranslation">{{ displayedTranslation }}</span>
    </div>
    <div v-if="controls" class="typing-controls" aria-label="Greek snippet navigation">
      <v-btn
          aria-label="Previous Greek snippet"
          color="primary"
          density="comfortable"
          icon="mdi-chevron-left"
          variant="tonal"
          @click="previousText"
      ></v-btn>
      <span>{{ currentIndex + 1 }} / {{ texts.length }}</span>
      <v-btn
          aria-label="Next Greek snippet"
          color="primary"
          density="comfortable"
          icon="mdi-chevron-right"
          variant="tonal"
          @click="nextText"
      ></v-btn>
    </div>
  </div>
</template>

<script>
import { ref, watch, onBeforeUnmount } from 'vue';

export default {
  name: 'TypingText',
  props: {
    texts: {
      type: Array,
      required: true,
    },
    controls: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const displayedText = ref('');
    const displayedTranslation = ref('');
    const currentIndex = ref(0);
    const timers = [];

    const clearTimers = () => {
      while (timers.length) {
        clearTimeout(timers.pop());
      }
    };

    const schedule = (callback, delay) => {
      const timer = setTimeout(callback, delay);
      timers.push(timer);
    };

    const typeText = (text, callback) => {
      let typedText = '';
      const typingSpeed = 50;
      for (let i = 0; i < text.length; i++) {
        schedule(() => {
          typedText += text[i];
          displayedText.value = typedText;
          if (i === text.length - 1 && callback) {
            callback();
          }
        }, i * typingSpeed);
      }
    };

    const typeCurrentText = () => {
      const currentText = props.texts[currentIndex.value];
      displayedText.value = '';
      displayedTranslation.value = '';

      typeText(currentText?.greek || '', () => {
        schedule(() => {
          typeTranslation(currentText?.translation || '');
        }, 350);
      });
    };

    const typeTranslation = (text, callback) => {
      let typedText = '';
      const typingSpeed = 50;
      for (let i = 0; i < text.length; i++) {
        schedule(() => {
          typedText += text[i];
          displayedTranslation.value = typedText;
          if (i === text.length - 1 && callback) {
            callback();
          }
        }, i * typingSpeed);
      }
    };

    const typeDialogue = () => {
      if (!props.texts.length) {
        displayedText.value = '';
        displayedTranslation.value = '';
        return;
      }

      if (props.controls) {
        typeCurrentText();
        return;
      }

      const { greek, translation } = props.texts[currentIndex.value];
      typeText(greek, () => {
        schedule(() => {
          typeTranslation(translation, () => {
            schedule(() => {
              currentIndex.value = (currentIndex.value + 1) % props.texts.length;
              displayedText.value = '';
              displayedTranslation.value = '';
              typeDialogue();
            }, 1000);
          });
        }, 1000);
      });
    };

    const setCurrentIndex = (index) => {
      if (!props.texts.length) return;

      clearTimers();
      currentIndex.value = (index + props.texts.length) % props.texts.length;
      typeCurrentText();
    };

    const nextText = () => {
      setCurrentIndex(currentIndex.value + 1);
    };

    const previousText = () => {
      setCurrentIndex(currentIndex.value - 1);
    };

    watch(
        () => [props.texts, props.controls],
        () => {
          clearTimers();
          displayedText.value = '';
          displayedTranslation.value = '';
          currentIndex.value = 0;
          typeDialogue();
        },
        { immediate: true }
    );

    onBeforeUnmount(clearTimers);

    return {
      currentIndex,
      displayedText,
      displayedTranslation,
      nextText,
      previousText,
    };
  },
};
</script>

<style>
#typing-text {
  font-size: 1.5em;
  color: var(--v-primary-base);
  white-space: normal;
}

.typing-body {
  min-height: 160px;
}

.typing-body p {
  width: 72px;
  height: 2px;
  margin: 16px 0;
  background: linear-gradient(90deg, #1c61d1, #1cd18c, #1cbcd1);
}

.typing-controls {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 18px;
  color: #536987;
  font-size: 0.9rem;
  font-weight: 700;
}

</style>
