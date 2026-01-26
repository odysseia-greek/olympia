<template>
  <div>
    <TextArea></TextArea>
  </div>
</template>

<script>
import { defineComponent } from 'vue';
import TextArea from '@/components/TextArea.vue';

export default defineComponent({
  name: 'Herodotos',
  components: {
    TextArea,
  },
  setup() {
    // If you have any state or logic, it will go here
    return {};
  },
});
</script>
