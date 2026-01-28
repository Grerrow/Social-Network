<!-- FollowRequestButton.vue -->
<template>
  <button @click="handleFollow" class="btn-follow" :class="buttonClass">
    {{ follow.buttonLabel }}
  </button>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useFollowStatus } from '@/composables/useFollowStatus'

const props = defineProps({
  userId: { type: Number, required: true },
  initialStatus: { type: String, default: 'none' }
})

const follow = useFollowStatus(props.initialStatus)

// Update status when prop changes
watch(
  () => props.initialStatus,
  (val) => {
    follow.status.value = val
  }
)

// Dynamic button styling based on state
const buttonClass = computed(() => {
  if (follow.isFollowing.value) return 'btn-following'
  if (follow.isRequested.value) return 'btn-requested'
  return 'btn-primary'
})

// Actions
const handleFollow = async () => {
  if (follow.canFollow.value) await follow.sendFollowRequest(props.userId)
  else if (follow.isFollowing.value) await follow.unfollow(props.userId)
  else if (follow.isRequested.value) await follow.cancelFollowRequest(props.userId)
}
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.btn-follow {
  padding: 12px 24px;
  border-radius: 12px;
  border: none;
  font-weight: 700;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
  letter-spacing: 0.01em;
  font-family: inherit;
}

.btn-primary {
  background: var(--honey-bronze);
  color: var(--ink-black);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(246, 189, 96, 0.35);
  background: #f7c570;
}

.btn-following {
  background: var(--muted-teal);
  color: white;
}

.btn-following:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(146, 191, 177, 0.35);
  background: #7da99c;
}

.btn-requested {
  background: #f7c570;
  color: var(--ink-black);
  border: 2px solid rgba(13, 19, 33, 0.12);
}

.btn-requested:hover:not(:disabled) {
  background: #f7c570;
}

.btn-follow:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-follow:active:not(:disabled) {
  transform: translateY(0);
}

@media (max-width: 768px) {
  .btn-follow {
    padding: 10px 20px;
    font-size: 0.8125rem;
  }
}
</style>
