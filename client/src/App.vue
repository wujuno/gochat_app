<template>
  <v-app>
    <v-container class="container">
      <Header />
      <router-view />
      <Footer />
    </v-container>
  </v-app>
</template>

<script>
import Header from "./components/Header.vue";
import Footer from "./components/Footer.vue";
export default {
  components: {
    Header,
    Footer,
  },
};
</script>

<style>
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 400px;
  max-width: 680px;
}
</style>
