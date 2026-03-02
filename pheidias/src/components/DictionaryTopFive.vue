<template>
  <v-card class="ma-4" color="background" variant="outlined">
    <v-card-title class="d-flex align-center">
      <span>Recent Top Searches</span>
      <v-spacer />
      <v-btn
        :loading="loading"
        color="secondary"
        icon="mdi-refresh"
        size="small"
        variant="text"
        @click="fetchTopFive"
      />
    </v-card-title>

    <v-card-text>
      <v-alert
        v-if="errorMessage"
        density="compact"
        type="warning"
        variant="tonal"
      >
        {{ errorMessage }}
      </v-alert>

      <v-table v-else-if="entries.length" density="comfortable">
        <thead>
        <tr>
          <th>Word</th>
          <th>Service</th>
          <th>Count</th>
          <th>Last used</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="entry in entries" :key="`${entry.word}-${entry.serviceName}`">
          <td>{{ entry.word }}</td>
          <td>
            <v-chip color="secondary" size="small" variant="tonal">
              {{ entry.serviceName }}
            </v-chip>
          </td>
          <td>{{ entry.count }}</td>
          <td>{{ formatLastUsed(entry.lastUsed) }}</td>
        </tr>
        </tbody>
      </v-table>

      <div v-else class="text-medium-emphasis">
        No search data available yet.
      </div>
    </v-card-text>
  </v-card>
</template>

<script>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { useApolloClient } from '@vue/apollo-composable';

import { CounterTopFive } from '@/constants/dictionaryCounterGraphql';

const POLL_INTERVAL_MS = 30000;

export default {
  name: 'DictionaryTopFive',
  props: {
    refreshToken: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    const { client } = useApolloClient();

    const loading = ref(false);
    const entries = ref([]);
    const errorMessage = ref('');
    let intervalId = null;

    async function fetchTopFive() {
      loading.value = true;
      errorMessage.value = '';

      try {
        const { data } = await client.query({
          query: CounterTopFive,
          fetchPolicy: 'no-cache',
        });

        entries.value = data?.counterTopFive?.topFive || [];
      } catch (error) {
        console.error(error);
        errorMessage.value = 'Unable to load top searches.';
      } finally {
        loading.value = false;
      }
    }

    function formatLastUsed(value) {
      if (!value) return '—';

      const date = new Date(value);
      if (Number.isNaN(date.getTime())) return value;

      return new Intl.DateTimeFormat(undefined, {
        dateStyle: 'medium',
        timeStyle: 'short',
      }).format(date);
    }

    watch(
      () => props.refreshToken,
      () => {
        fetchTopFive();
      }
    );

    onMounted(() => {
      fetchTopFive();
      intervalId = window.setInterval(fetchTopFive, POLL_INTERVAL_MS);
    });

    onBeforeUnmount(() => {
      if (intervalId) {
        window.clearInterval(intervalId);
      }
    });

    return {
      loading,
      entries,
      errorMessage,
      fetchTopFive,
      formatLastUsed,
    };
  },
};
</script>
