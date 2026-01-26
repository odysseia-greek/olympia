<template>
  <div id="search">
    <v-app id="notfound">
      <v-main style="background:#fdf6e3">
        <div class="text-center">
          <h1 class="error-title">Oops! Page Not Found</h1>
          <img src="../assets/pagenotfound.webp" alt="Page not found" class="error-image">
          <p class="error-message">
            It looks like the page you are looking for doesn't exist.
            Maybe it was moved or you mistyped the URL.
            Don't worry, let's get you back on track!
          </p>
          <v-btn
              v-for="item in menuItems"
              :key="item.path"
              color="primary"
              class="ma-2"
              :to="item.path"
          >
            {{ item.title }}
            <v-icon end>{{ item.icon }}</v-icon>
          </v-btn>
        </div>
      </v-main>
    </v-app>
  </div>
</template>

<script>
import { defineComponent, ref } from 'vue';

export default defineComponent({
  name: 'App',
  setup() {
    const appTitle = ref('Odysseia');
    const sidebar = ref(false);
    const menuItems = ref([
      { title: 'ΠΑΝΤΑΡΕΙ', path: '/', icon: 'mdi-home-variant' },
    ]);

    return {
      appTitle,
      sidebar,
      menuItems,
    };
  },
});
</script>

<style>
.error-image {
  width: 100%;
  max-width: 400px;
  margin: 20px auto;
}

.error-title {
  font-size: 2.5em;
  margin: 20px 0;
}

.error-message {
  font-size: 1.2em;
  margin: 10px 0 30px 0;
  color: #757575;
}
</style>
