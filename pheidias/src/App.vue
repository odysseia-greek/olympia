<template>
  <v-app id="odysseia">
    <v-app-bar floating color="footer">

      <template v-slot:prepend >
        <v-app-bar-nav-icon v-if="hamburgerMenu"></v-app-bar-nav-icon>
      </template>

      <v-app-bar-title>{{ appTitle }}</v-app-bar-title>
      <v-list-item v-if="!hamburgerMenu" v-for="(item, i) in menuItems" :key="i" :to="item.path" link @click="toTop" color="papyrus">
        <v-icon>{{ item.icon }}</v-icon>
        <v-list-item-title>{{ item.title }}</v-list-item-title>
      </v-list-item>
      <v-menu v-if="hamburgerMenu" v-model="menu" activator="parent" :close-on-content-click="true" location="left">
        <v-list nav min-width="30em" style="background-color: #e0f7fa">
          <v-list-item v-for="(item, i) in menuItems" :key="i" :to="item.path" link @click="toTop" color="primary">
            <v-icon>{{ item.icon }}</v-icon>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-app-bar>

    <v-main>
      <router-view></router-view>
      <v-btn
          v-show="fab"
          @click="toTop"
          color="triadic"
          size="large"
          icon="mdi-arrow-up"
          style="position: fixed; bottom: 16px; right: 16px; z-index: 1000;"
      >
      </v-btn>
      <v-footer color="background">
        <v-card flat width="100%" class="footer lighten-1 text-center" color="footer">
          <v-card-text>
            <v-btn
                class="ma-2"
                :href="item.path"
                v-for="item in footerItems"
                :key="item.icon"
                :icon="item.icon"
                variant="text"
            >
            </v-btn>
          </v-card-text>

          <v-divider></v-divider>

          <v-card-text class="white--text">
            {{ new Date().getFullYear() }} — <strong>Odysseia-greek</strong>
          </v-card-text>

          <v-card-text class="white--text mt-4">
            <v-row justify="center" align="center">
              <v-col cols="12" md="6">
                <v-row no-gutters>
                  <v-col class="center">
                    <strong>Pheidias Version: {{ frontendVersion }}</strong>
                  </v-col>
                </v-row>
              </v-col>
              <v-col cols="12" md="6">
                <v-row no-gutters>
                  <v-col class="center">
                    <strong>Homeros Version: {{ gatewayVersion }}</strong>
                    <v-icon
                        :color="gatewayHealth ? 'green' : 'red'"
                        class="ml-2"
                        size="16"
                    >
                      mdi-circle
                    </v-icon>
                    <small v-if="gatewayHealth"> (Checked: {{ lastCheckedTime }})</small>
                  </v-col>
                </v-row>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-footer>
    </v-main>
  </v-app>
</template>

<script>
import { ref, onMounted, onBeforeUnmount } from 'vue';

export default {
  name: 'App',
  setup() {
    const fab = ref(false);
    const appTitle = ref('Odysseia-Greek');
    const closeOnClick = ref(true);
    const menu = ref(false);
    const footerItems = ref([
      { icon: 'mdi-github', path: 'https://github.com/odysseia-greek' },
      { icon: 'mdi-linkedin', path: 'https://nl.linkedin.com/in/joeri-vrijaldenhoven-22713a80' },
      { icon: 'mdi-book-open-outline', path: '/ploutarchos/' },
    ]);
    const menuItems = ref([
      { title: 'Home', path: '/', icon: 'mdi-home-variant' },
      { title: 'Quiz', path: '/quiz', icon: 'mdi-alphabet-greek' },
      { title: 'Texts', path: '/texts', icon: 'mdi-bookshelf' },
      { title: 'Grammar', path: '/grammar', icon: 'mdi-feather' },
      { title: 'Dictionary', path: '/dictionary', icon: 'mdi-magnify' },
    ]);

    const frontendVersion = ref(process.env.VERSION || 'unknown');
    const gatewayVersion = ref('loading...');
    const gatewayHealth = ref(false);
    const lastCheckedTime = ref('');
    const hamburgerMenu = ref(true);

    const fetchGatewayHealth = async () => {
      let url = document.location.origin + '/homeros/v1/health';
      if (process.env.NODE_ENV === 'k3d') {
        url = 'http://k3d-odysseia.greek:8080/homeros/v1/health';
      }

      try {
        const response = await fetch(url);
        if (response.ok) {
          const data = await response.json();
          gatewayVersion.value = data.version || 'unknown';
          gatewayHealth.value = data.healthy;
          lastCheckedTime.value = new Date().toLocaleTimeString();
        } else {
          gatewayVersion.value = 'unknown';
          gatewayHealth.value = false;
          lastCheckedTime.value = new Date().toLocaleTimeString();
        }
      } catch (error) {
        console.error('Error fetching gateway health:', error);
        gatewayVersion.value = 'unknown';
        gatewayHealth.value = false;
        lastCheckedTime.value = new Date().toLocaleTimeString();
      }
    };

    const onScroll = (e) => {
      if (typeof window === 'undefined') return;
      const top = window.pageYOffset || e.target.scrollTop || 0;
      fab.value = top > 20;
    };

    const toTop = () => {
      window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    const checkScreenSize = () => {
      hamburgerMenu.value = window.innerWidth <= 810;
    };

    // Set up a recurring interval to update the gateway version every 30 seconds
    onMounted(() => {
      checkScreenSize();
      window.addEventListener('scroll', onScroll);
      fetchGatewayHealth();


      const intervalId = setInterval(() => {
        fetchGatewayHealth();
      }, 30000); // 30 seconds

      // Cleanup the interval when the component is unmounted
      onBeforeUnmount(() => {
        clearInterval(intervalId);
      });
    });

    return {
      fab,
      appTitle,
      closeOnClick,
      footerItems,
      menuItems,
      menu,
      frontendVersion,
      gatewayVersion,
      gatewayHealth,
      lastCheckedTime,
      hamburgerMenu,
      onScroll,
      toTop,
    };
  },
};
</script>

<style>
.v-data-table tbody tr:hover:not(.v-data-table__expanded__content):not(.v-data-table__empty-wrapper) {
  background-color: transparent !important;
}
</style>
