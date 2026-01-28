<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// Get the layout from route meta, default to empty layout
const layoutComponent = computed(() => route.meta.layout || null)
</script>

<template>
  <div id="app">
    <!-- Dynamic layout based on current route -->
    <component :is="layoutComponent">
      <router-view />
    </component>
  </div>
</template>

<!-- <style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f5f5f5;
  color: #212529;
}

#app {
  min-height: 100vh;
}

html,
body,
#app {
  height: 100%;
}
</style> -->
