<template>
  <main class="corpus-content">
    <v-container class="corpus-container">
      <section class="section-heading">
        <div>
          <v-chip color="secondary" variant="tonal">Guided corpus</v-chip>
          <h2>Choose a chapter</h2>
        </div>
        <p>Build up the vocabulary and grammar first, then translate an original passage with hints when you need them.</p>
      </section>

      <v-alert v-if="optionsError" type="error" variant="tonal">{{ optionsError }}</v-alert>
      <div v-else class="chapter-grid" aria-label="Available corpus chapters">
        <v-card
            v-for="option in chapters"
            :key="option.chapter"
            :class="['chapter-card', { selected: option.chapter === selectedChapter }]"
            :variant="option.chapter === selectedChapter ? 'elevated' : 'outlined'"
            @click="selectChapter(option.chapter)"
        >
          <v-card-text>
            <div class="chapter-meta">Chapter {{ option.order }} · Level {{ option.level }}</div>
            <h3>{{ option.title }}</h3>
            <v-btn class="mt-4" color="primary" :variant="option.chapter === selectedChapter ? 'flat' : 'tonal'">
              {{ option.chapter === selectedChapter ? 'Selected' : 'Start chapter' }}
            </v-btn>
          </v-card-text>
        </v-card>
      </div>

      <div v-if="chapterLoading" class="loading-state">
        <v-progress-circular color="primary" indeterminate />
        <span>Opening chapter…</span>
      </div>

      <v-alert v-else-if="chapterError" class="mt-6" type="error" variant="tonal">{{ chapterError }}</v-alert>

      <template v-else-if="chapter">
        <section class="chapter-intro">
          <div class="chapter-kicker">Chapter {{ chapter.order }} · Level {{ chapter.level }}</div>
          <h1>{{ chapter.title }}</h1>
          <p class="description">{{ chapter.description }}</p>
          <v-sheet class="context-card" color="secondaryPapyrus">
            <v-icon color="primary">mdi-book-open-page-variant</v-icon>
            <p>{{ chapter.context }}</p>
          </v-sheet>
        </section>

        <section class="learning-grid">
          <v-card class="learning-card" variant="outlined">
            <v-card-title><v-icon class="mr-2" color="secondary">mdi-format-list-bulleted</v-icon>Vocabulary</v-card-title>
            <v-list bg-color="transparent">
              <v-list-item v-for="word in chapter.vocabulary" :key="word.greek">
                <template #prepend><span class="greek-term">{{ word.greek }}</span></template>
                <v-list-item-title>{{ word.translation }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-card>

          <v-card class="learning-card" variant="outlined">
            <v-card-title><v-icon class="mr-2" color="primary">mdi-feather</v-icon>Grammar to notice</v-card-title>
            <v-expansion-panels variant="accordion">
              <v-expansion-panel v-for="item in chapter.grammar" :key="item.grammar">
                <v-expansion-panel-title>{{ item.title }}</v-expansion-panel-title>
                <v-expansion-panel-text>
                  <p>{{ item.explanation }}</p>
                  <div v-if="item.example" class="grammar-example">
                    <strong>{{ item.example.greek }}</strong>
                    <span>{{ item.example.translation }}</span>
                    <small v-if="item.example.note">{{ item.example.note }}</small>
                  </div>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
          </v-card>
        </section>

        <section class="translation-workspace">
          <div class="section-heading">
            <div>
              <v-chip color="primary" variant="tonal">Read & translate</v-chip>
              <h2>Make your own translation</h2>
            </div>
            <p>Try the Greek unaided. Reveal the reading hints one at a time if you get stuck.</p>
          </div>

          <v-card v-for="text in chapter.texts" :key="text.text" class="passage-card">
            <v-card-text>
              <div class="passage-heading">
                <div>
                  <div class="chapter-meta">{{ text.type }} · {{ text.source.dialect }}</div>
                  <h3>{{ text.title }}</h3>
                  <span>{{ text.source.author }}, {{ text.source.work }} {{ text.source.reference }}</span>
                </div>
              </div>
              <p class="greek-passage">{{ text.greek }}</p>

              <div v-if="text.readingHints.length" class="hint-panel">
                <div class="hint-heading">
                  <strong>Reading hints</strong>
                  <span>{{ visibleHintCount(text.text) }} / {{ text.readingHints.length }}</span>
                </div>
                <ol v-if="visibleHintCount(text.text)">
                  <li v-for="hint in text.readingHints.slice(0, visibleHintCount(text.text))" :key="hint">{{ hint }}</li>
                </ol>
                <v-btn
                    v-if="visibleHintCount(text.text) < text.readingHints.length"
                    color="secondary"
                    prepend-icon="mdi-lightbulb-on-outline"
                    variant="tonal"
                    @click="revealHint(text.text)"
                >Reveal next hint</v-btn>
              </div>

              <v-textarea
                  v-model="answers[text.text]"
                  class="mt-5"
                  label="Your translation"
                  placeholder="Write a complete translation of the passage…"
                  rows="4"
                  variant="outlined"
                  auto-grow
                  @update:model-value="checkResult = null"
              />
            </v-card-text>
          </v-card>

          <div class="submit-row">
            <span>{{ completedAnswers }} of {{ chapter.texts.length }} passage(s) attempted</span>
            <v-btn :disabled="!canCheck" :loading="checking" color="primary" prepend-icon="mdi-check-decagram" size="large" @click="checkChapter">
              Check translation
            </v-btn>
          </div>

          <v-alert v-if="checkError" class="mt-5" type="error" variant="tonal">{{ checkError }}</v-alert>

          <section v-if="checkResult" ref="resultsRef" class="results-section">
            <div class="section-heading">
              <div><v-chip color="secondary" variant="tonal">Complete</v-chip><h2>Compare translations</h2></div>
              <p>Read the full translation beside your own, then return to the Greek and notice where your choices differ.</p>
            </div>
            <v-card v-for="result in checkResult.texts" :key="result.text" class="result-card" variant="outlined">
              <v-card-text>
                <div class="comparison-label">Greek source</div>
                <p class="greek-result">{{ result.sourceText }}</p>
                <div class="comparison-grid">
                  <div><span class="comparison-label">Full translation</span><p>{{ result.actualText }}</p></div>
                  <div><span class="comparison-label">Your translation</span><p>{{ result.learnerText }}</p></div>
                </div>
              </v-card-text>
            </v-card>
          </section>
        </section>
      </template>
    </v-container>
  </main>
</template>

<script setup>
import { computed, nextTick, onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { apolloClient } from '@/apollo';
import { HerodotosChapter, HerodotosChapterOptions, HerodotosCheckChapter } from '@/constants/graphql';

const route = useRoute();
const router = useRouter();
const chapters = ref([]);
const selectedChapter = ref('');
const chapter = ref(null);
const answers = reactive({});
const hintCounts = reactive({});
const optionsError = ref('');
const chapterError = ref('');
const checkError = ref('');
const chapterLoading = ref(false);
const checking = ref(false);
const checkResult = ref(null);
const resultsRef = ref(null);

const completedAnswers = computed(() => Object.values(answers).filter((answer) => answer?.trim()).length);
const canCheck = computed(() => chapter.value?.texts?.length && completedAnswers.value === chapter.value.texts.length);

const visibleHintCount = (text) => hintCounts[text] || 0;
const revealHint = (text) => { hintCounts[text] = visibleHintCount(text) + 1; };

async function selectChapter(chapterId, updateRoute = true) {
  selectedChapter.value = chapterId;
  chapter.value = null;
  checkResult.value = null;
  chapterError.value = '';
  chapterLoading.value = true;
  if (updateRoute) await router.replace({ query: { ...route.query, mode: 'corpus', chapter: chapterId } });
  try {
    const { data } = await apolloClient.query({ query: HerodotosChapter, variables: { chapter: chapterId }, fetchPolicy: 'network-only' });
    chapter.value = data.chapter;
    Object.keys(answers).forEach((key) => delete answers[key]);
    Object.keys(hintCounts).forEach((key) => delete hintCounts[key]);
    data.chapter.texts.forEach((text) => { answers[text.text] = ''; hintCounts[text.text] = 0; });
    await nextTick();
    document.querySelector('.chapter-intro')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  } catch (error) {
    chapterError.value = error.message || 'Unable to load this chapter.';
  } finally {
    chapterLoading.value = false;
  }
}

async function checkChapter() {
  checking.value = true;
  checkError.value = '';
  checkResult.value = null;
  try {
    const { data } = await apolloClient.query({
      query: HerodotosCheckChapter,
      variables: { input: { chapter: chapter.value.chapter, answers: chapter.value.texts.map((text) => ({ text: text.text, learnerText: answers[text.text].trim() })) } },
      fetchPolicy: 'network-only',
    });
    checkResult.value = data.checkChapter;
    await nextTick();
    resultsRef.value?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  } catch (error) {
    checkError.value = error.message || 'Unable to check your translation.';
  } finally {
    checking.value = false;
  }
}

onMounted(async () => {
  try {
    const { data } = await apolloClient.query({ query: HerodotosChapterOptions, fetchPolicy: 'network-only' });
    chapters.value = [...data.chapterOptions.chapters].sort((a, b) => a.order - b.order);
    const initial = route.query.chapter || chapters.value[0]?.chapter;
    if (initial) await selectChapter(initial, !route.query.chapter);
  } catch (error) {
    optionsError.value = error.message || 'Unable to load the available chapters.';
  }
});
</script>

<style scoped>
.corpus-content { background: linear-gradient(155deg, #d9f0f7 0%, #f2fbf7 48%, #fffaf0 100%); color: #20334f; min-height: 70vh; }
.corpus-container { max-width: 1240px; padding-top: 46px; padding-bottom: 70px; }
.section-heading { display: flex; align-items: end; justify-content: space-between; gap: 28px; margin-bottom: 22px; }
.section-heading h2 { margin: 8px 0 0; font-size: clamp(1.65rem, 3vw, 2.35rem); }
.section-heading p { max-width: 570px; margin: 0; color: #536987; line-height: 1.65; }
.chapter-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(235px, 1fr)); gap: 16px; margin-bottom: 52px; }
.chapter-card { cursor: pointer; background: rgba(253, 246, 227, .88); border-color: rgba(28, 97, 209, .2); transition: transform .18s ease, box-shadow .18s ease; }
.chapter-card:hover, .chapter-card.selected { transform: translateY(-3px); box-shadow: 0 12px 28px rgba(11, 39, 85, .14); }
.chapter-card h3 { margin-top: 8px; font-size: 1.25rem; }
.chapter-meta, .chapter-kicker { color: #64789e; font-size: .78rem; font-weight: 800; letter-spacing: .06em; text-transform: uppercase; }
.loading-state { display: flex; justify-content: center; align-items: center; gap: 14px; padding: 48px; }
.chapter-intro { scroll-margin-top: 82px; margin: 20px 0 30px; }
.chapter-intro h1 { margin: 8px 0; font-size: clamp(2rem, 5vw, 3.4rem); }
.description { max-width: 760px; font-size: 1.15rem; line-height: 1.65; }
.context-card { display: flex; gap: 16px; margin-top: 22px; padding: 22px; border: 1px solid rgba(28, 97, 209, .16); border-radius: 8px; }
.context-card p { margin: 0; line-height: 1.7; }
.learning-grid { display: grid; grid-template-columns: minmax(280px, .8fr) minmax(360px, 1.2fr); gap: 20px; margin-bottom: 50px; }
.learning-card, .passage-card, .result-card { background: #fdf6e3; border-color: rgba(28, 97, 209, .16); }
.greek-term { display: inline-block; min-width: 130px; color: #10284b; font-weight: 700; }
.grammar-example { display: grid; gap: 3px; margin-top: 14px; padding: 14px; background: rgba(28, 188, 209, .1); border-radius: 6px; }
.grammar-example strong { font-size: 1.1rem; }
.grammar-example small { margin-top: 5px; color: #536987; }
.translation-workspace { scroll-margin-top: 82px; }
.passage-card { margin-bottom: 20px; padding: 10px; }
.passage-heading h3 { margin: 5px 0; font-size: 1.4rem; }
.greek-passage { margin: 28px 0; color: #10284b; font-family: 'Noto Sans Coptic', serif; font-size: clamp(1.25rem, 2.3vw, 1.65rem); line-height: 1.9; }
.hint-panel { padding: 16px; background: rgba(28, 209, 140, .1); border: 1px solid rgba(28, 209, 140, .24); border-radius: 8px; }
.hint-heading { display: flex; justify-content: space-between; margin-bottom: 10px; }
.hint-panel ol { margin: 0 0 14px; padding-left: 22px; line-height: 1.7; }
.submit-row { display: flex; align-items: center; justify-content: space-between; gap: 20px; margin: 28px 0 54px; }
.submit-row span { color: #536987; }
.results-section { scroll-margin-top: 82px; }
.result-card { margin-bottom: 16px; }
.comparison-label { display: block; margin-bottom: 7px; color: #64789e; font-size: .76rem; font-weight: 800; letter-spacing: .06em; text-transform: uppercase; }
.greek-result { margin-bottom: 22px; font-size: 1.2rem; line-height: 1.7; }
.comparison-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.comparison-grid > div { padding: 18px; background: rgba(255,255,255,.52); border-radius: 7px; }
.comparison-grid p { margin: 0; line-height: 1.65; }
@media (max-width: 760px) { .section-heading { display: block; } .section-heading p { margin-top: 12px; } .learning-grid, .comparison-grid { grid-template-columns: 1fr; } .submit-row { align-items: stretch; flex-direction: column; } }
</style>
