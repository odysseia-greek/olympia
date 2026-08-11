<template>
  <div class="herodotos-page">
    <section class="mode-hero" :style="{ backgroundImage: `url(${heroImage})` }">
      <div class="mode-shade">
        <v-container class="mode-container">
          <div class="mode-panel">
            <div class="section-label">Texts</div>
            <h1>Herodotos</h1>
            <p>Read complete works from the library, or learn through a guided corpus chapter.</p>
            <v-btn-toggle v-model="mode" class="mode-toggle" color="primary" mandatory variant="outlined" @update:model-value="changeMode">
              <v-btn value="full" prepend-icon="mdi-bookshelf" size="large">Full text</v-btn>
              <v-btn value="corpus" prepend-icon="mdi-map-marker-path" size="large">Corpus mode</v-btn>
            </v-btn-toggle>
          </div>
        </v-container>
      </div>
    </section>
    <TextArea v-if="mode === 'full'" class="embedded-text-area" />
    <CorpusMode v-else />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import TextArea from '@/components/TextArea.vue';
import CorpusMode from '@/components/CorpusMode.vue';

const route = useRoute();
const router = useRouter();
const mode = ref(route.query.mode === 'corpus' ? 'corpus' : 'full');
const heroImage = ref('');

function changeMode(nextMode) {
  const query = { ...route.query, mode: nextMode };
  if (nextMode === 'full') delete query.chapter;
  router.replace({ query });
}

onMounted(() => import('@/assets/alexandria.webp').then((module) => { heroImage.value = module.default; }));
</script>

<style scoped>
.mode-hero { min-height: 410px; background-position: center 42%; background-size: cover; }
.mode-shade { min-height: 410px; display: flex; align-items: center; background: linear-gradient(90deg, rgba(22,20,15,.74), rgba(38,31,21,.42) 55%, rgba(26,22,16,.18)); }
.mode-container { padding-top: 54px; padding-bottom: 42px; }
.mode-panel { max-width: 820px; padding: 26px; color: #20334f; background: rgba(253,246,227,.96); border: 1px solid rgba(28,97,209,.16); border-radius: 8px; box-shadow: 0 14px 36px rgba(11,39,85,.18); }
.section-label { color: #64789e; font-size: .82rem; font-weight: 800; letter-spacing: .08em; text-transform: uppercase; }
.mode-panel h1 { margin: 7px 0; font-size: clamp(2rem, 5vw, 3.2rem); }
.mode-panel p { max-width: 610px; margin: 0 0 22px; color: #344765; line-height: 1.65; }
.mode-toggle :deep(.v-btn) { min-width: 180px; text-transform: none; }
.embedded-text-area :deep(.texts-hero) { display: none; }
@media (max-width: 600px) { .mode-hero, .mode-shade { min-height: auto; } .mode-toggle { display: grid; width: 100%; } .mode-toggle :deep(.v-btn) { width: 100%; } }
</style>
