<template>
  <v-data-table
      v-if="history.length"
      :headers="headers"
      :items="history"
      :items-per-page="5"
      :items-per-page-options="[5, 10, 25, -1]"
      :disable-sort="true"
      class="elevation-1 mt-5"
  >
    <template v-slot:item.input="{ item }">
      <v-chip
          :color="item.correct ? 'success' : 'error'"
          variant="tonal"
      >
        {{ item.input }}
      </v-chip>
    </template>
  </v-data-table>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { saveHistory, getHistory } from '@/composables/useBoule';

const headers = [
  { text: 'Greek', value: 'greek', align: 'center' },
  { text: 'Provided Answer', value: 'input', align: 'center' },
];

const props = defineProps({
  newEntry: Object, // single object like { greek, input, correct }
  maxLength: {
    type: Number,
    default: 50
  }
});

const emit = defineEmits(['update']);
const history = ref([]);

function persist() {
  saveHistory(history.value);
}

function restore() {
  history.value = getHistory();
  emit('update', history.value);
}

function addEntry(entry) {
  if (!entry || !entry.greek || !entry.input) return;

  const normalized = {
    greek: entry.greek,
    input: entry.input,
    correct: !!entry.correct,
  };

  history.value = [normalized, ...history.value].slice(0, props.maxLength);
  persist();
  emit('update', history.value);
}

onMounted(() => {
  restore();
  if (props.newEntry) {
    addEntry(props.newEntry);
  }
});

watch(() => props.newEntry, (entry) => {
  if (entry) addEntry(entry);
});
</script>
