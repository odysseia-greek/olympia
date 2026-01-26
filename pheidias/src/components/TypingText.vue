<template>
  <div id="typing-text">
    <span>{{ displayedText }}</span>
    <p v-if="displayedTranslation">------</p>
    <span v-if="displayedTranslation">{{ displayedTranslation }}</span>
  </div>
</template>

<script>
import { ref, watch, onMounted } from 'vue';

export default {
  name: 'TypingText',
  props: {
    texts: {
      type: Array,
      required: true,
    },
  },
  setup(props) {
    const displayedText = ref('');
    const displayedTranslation = ref('');
    const currentIndex = ref(0);

    const typeText = (text, callback) => {
      let typedText = '';
      const typingSpeed = 50;
      for (let i = 0; i < text.length; i++) {
        setTimeout(() => {
          typedText += text[i];
          displayedText.value = typedText;
          if (i === text.length - 1 && callback) {
            callback();
          }
        }, i * typingSpeed);
      }
    };

    const typeTranslation = (text, callback) => {
      let typedText = '';
      const typingSpeed = 50;
      for (let i = 0; i < text.length; i++) {
        setTimeout(() => {
          typedText += text[i];
          displayedTranslation.value = typedText;
          if (i === text.length - 1 && callback) {
            callback();
          }
        }, i * typingSpeed);
      }
    };

    const typeDialogue = () => {
      const { greek, translation } = props.texts[currentIndex.value];
      typeText(greek, () => {
        setTimeout(() => {
          typeTranslation(translation, () => {
            setTimeout(() => {
              currentIndex.value = (currentIndex.value + 1) % props.texts.length;
              displayedText.value = '';
              displayedTranslation.value = '';
              typeDialogue();
            }, 1000);
          });
        }, 1000);
      });
    };

    watch(
        () => props.texts,
        () => {
          displayedText.value = '';
          displayedTranslation.value = '';
          currentIndex.value = 0;
          typeDialogue();
        },
        { immediate: true }
    );

    return {
      displayedText,
      displayedTranslation,
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

</style>
