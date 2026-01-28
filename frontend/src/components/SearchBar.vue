<!-- Search Bar Component -->
<template>
  <div class="search-bar">
    <div class="search-input-wrapper">
      <span class="search-icon">
        <svg width="18" height="18" viewBox="0 0 20 20" fill="none">
          <circle cx="9" cy="9" r="7" stroke="currentColor" stroke-width="2"/>
          <line x1="14.4142" y1="14.7071" x2="18" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </span>
      <input
        v-model="query"
        @input="onInput"
        placeholder="Search users or groups..."
        class="search-input"
      />
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'SearchBar',
  emits: ['search'],
  setup(_, { emit }) {
    const query = ref('')

    const onInput = () => {
      emit('search', query.value)
    }

    return { query, onInput }
  }
}
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.search-bar {
  display: flex;
  justify-content: center;
  width: 100%;
}

.search-input-wrapper {
  position: relative;
  width: 100%;
  max-width: 500px;
}

.search-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  display: flex;
  align-items: center;
  color: rgba(13, 19, 33, 0.5);
  transition: color 0.2s ease;
}

.search-input {
  width: 100%;
  padding: 12px 16px 12px 48px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  outline: none;
  font-size: 0.9375rem;
  font-family: inherit;
  background: white;
  color: var(--ink-black);
  transition: all 0.2s ease;
  box-sizing: border-box;
  letter-spacing: 0.01em;
}

.search-input::placeholder {
  color: rgba(13, 19, 33, 0.5);
}

.search-input:focus {
  border-color: var(--honey-bronze);
  background: white;
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

.search-input:focus + .search-icon,
.search-input-wrapper:focus-within .search-icon {
  color: var(--honey-bronze);
}

@media (max-width: 768px) {
  .search-input {
    font-size: 0.875rem;
  }
}
</style>
