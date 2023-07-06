<template>
  <v-app>
    <v-container>
      <v-btn color="primary">Hello Vuetify</v-btn>
    </v-container>
    <router-view />
  </v-app>
</template>

<script>
import { VApp, VContainer, VBtn } from "vuetify/lib";

export default {
  components: {
    VApp,
    VContainer,
    VBtn,
  },
};
</script>
