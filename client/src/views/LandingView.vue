<template>
  <v-container>
    <div class="landing-alert-box mb-5">
      <div class="text-h6">
        실시간 챗 어플리케이션입니다. Vue.js v2 로 만들어졌으며 Vuetify를
        사용하고 있습니다. 커뮤니케이션을 위해 웹소켓을 사용하고 있습니다. 지금
        바로 회원으로 등록하고 사용하세요!
      </div>
    </div>
    <v-container>
      <v-row justify="center" align="center">
        <router-link :to="{ name: 'REGISTER' }"
          ><v-col align="center" col="auto">
            <v-btn color="success"> register </v-btn></v-col
          ></router-link
        >
        <router-link :to="{ name: 'LOGIN' }"
          ><v-col align="center" col="auto">
            <v-btn> login </v-btn></v-col
          ></router-link
        >
      </v-row>
    </v-container>
  </v-container>
</template>

<script></script>

<style></style>
