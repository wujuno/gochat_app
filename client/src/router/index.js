import Vue from "vue";
import VueRouter from "vue-router";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "LANDING",
    component: () => import("../views/LandingView.vue"),
  },
  {
    path: "/register",
    name: "REGISTER",
    component: () => import("../views/RegisterView.vue"),
  },
  {
    path: "/login",
    name: "LOGIN",
    component: () => import("../views/LoginView.vue"),
  },
  {
    path: "/chat",
    name: "CHAT",
    component: () => import("../views/ChatView.vue"),
  },
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

export default router;
