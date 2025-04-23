<script setup lang="ts">
import { defineAsyncComponent } from 'vue'
import { useRoute } from 'vue-router'

const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"))
const route = useRoute()
</script>

<template>
  <div class="app-layout">
    <Sidebar v-if="route.path !== '/login'" />
    <div class="main-content">
      <router-view></router-view>
    </div>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
}

.main-content {
  flex-grow: 1;
  /* Add padding or other styles as needed to prevent content overlap */
  padding: 1rem; /* Example padding */
  height: 100vh;
  overflow-y: auto; /* Allow scrolling for main content if it overflows */
}

/* Ensure router-view content takes up space */
.main-content > :deep(main) {
  display: flex;
  flex-direction: column;
  min-height: 100%;
}
</style>
