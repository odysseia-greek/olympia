<template>
  <div>
    <DictionaryArea></DictionaryArea>
  </div>
</template>

<script>
import { defineComponent } from 'vue';
import DictionaryArea from '@/components/DictionaryArea.vue';

export default defineComponent({
  name: 'Alexandros',
  components: {
    DictionaryArea,
  },
  setup() {
    // If you have any state or logic, it will go here
    return {};
  },
});
</script>





