<template>
  <router-link :to="{ name: 'LANDING' }" class="link-reset">
    <div class="chat-icon">
      <v-icon size="40">mdi-chat</v-icon>
      <span class="chat-text">Go Chat</span>
    </div>
  </router-link>
</template>

<script>
export default {};
</script>

<style>
.chat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.chat-text {
  margin-left: 0.5rem;
  font-size: 2rem;
  font-weight: 700;
  color: inherit;
}

.link-reset {
  text-decoration: none;
}
</style>
