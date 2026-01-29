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

