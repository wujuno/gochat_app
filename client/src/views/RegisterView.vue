<template>
  <div>
    <router-link v-if="redirect" :to="redirectTo" replace></router-link>
    <v-container class="mt-5">
      <v-card class="mx-auto" width="440px">
        <v-card-text>
          <v-form
            :class="{ 'is-invalid': isInvalid }"
            @submit.prevent="onSubmit"
          >
            <v-text-field
              v-model="username"
              :error-messages="isInvalid ? [message] : []"
              label="Username"
              name="username"
              type="text"
              :counter="8"
              required
              >{{ username }}</v-text-field
            >
            <v-text-field
              v-model="password"
              label="Password"
              name="password"
              type="password"
              :counter="14"
              required
              >{{ password }}</v-text-field
            >
            <v-btn block color="success" type="submit" class="mt-5">
              Register
            </v-btn>
          </v-form>
        </v-card-text>
      </v-card>
    </v-container>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      username: "",
      password: "",
      message: "",
      isInvalid: "",
      endpoint: 'http://localhost:8080/"register',
      redirect: false,
      redirectTo: "/chat?u=",
    };
  },
  methods: {
    async onSubmit() {
      try {
        const res = await axios.post(this.endpoint, {
          username: this.username,
          password: this.password,
        });

        console.log("register", res);
      } catch (error) {
        console.log(error);
      }
    },
  },
};
</script>

<style></style>
