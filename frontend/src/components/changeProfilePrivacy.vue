<template>
  <div class="change-profile-privacy">
    <h2>Profile Privacy</h2>
    <p>Your current profile is: <strong>{{ isPrivate ? 'Private' : 'Public' }}</strong></p>
    <p class="privacy-description">
      {{ isPrivate 
        ? 'Only approved followers can see your posts and profile information.' 
        : 'Anyone can see your posts and profile information.' 
      }}
    </p>
    <button @click="togglePrivacy" class="btn-toggle" :disabled="loading">
      {{ loading ? 'Updating...' : `Make Profile ${isPrivate ? 'Public' : 'Private'}` }}
    </button>
    <p v-if="message" class="message" :class="messageType">{{ message }}</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useApi } from '@/composables/useApi'

const { getApiUrl } = useApi()

const isPrivate = ref(false)
const loading = ref(false)
const message = ref('')
const messageType = ref('')

const fetchPrivacyStatus = async () => {
  try {
    const res = await fetch(getApiUrl('/profile'), {
      credentials: 'include'
    })

    if (!res.ok) throw new Error()

    const data = await res.json()
    isPrivate.value = Boolean(data.is_private)
  } catch {
    message.value = 'Error loading privacy settings'
    messageType.value = 'error'
  }
}

const togglePrivacy = async () => {
  loading.value = true
  message.value = ''

  try {
    const res = await fetch(getApiUrl('/profile/privacy'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        is_private: !isPrivate.value
      })
    })

    if (!res.ok) throw new Error()

    isPrivate.value = !isPrivate.value
    message.value = `Profile is now ${isPrivate.value ? 'Private' : 'Public'}`
    messageType.value = 'success'
  } catch {
    message.value = 'Failed to update privacy settings'
    messageType.value = 'error'
  } finally {
    loading.value = false
  }
}

onMounted(fetchPrivacyStatus)
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

/* ===== Section Wrapper ===== */
.change-profile-privacy {
  max-width: 560px;
}

/* ===== Title & Text ===== */
.change-profile-privacy h2 {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--ink-black);
  margin: 0 0 16px 0;
  letter-spacing: -0.01em;
}

.change-profile-privacy p {
  margin: 0 0 16px 0;
  font-size: 0.9375rem;
  color: rgba(13, 19, 33, 0.8);
  line-height: 1.5;
}

.privacy-description {
  padding: 12px 16px;
  background: var(--lavender-mist);
  border-radius: 12px;
  border-left: 3px solid var(--honey-bronze);
  font-size: 0.875rem;
  color: rgba(13, 19, 33, 0.7);
}

/* ===== Button ===== */
.btn-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 12px 24px;
  border-radius: 12px;
  border: none;
  background: var(--honey-bronze);
  color: var(--ink-black);
  font-size: 0.9375rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: inherit;
  letter-spacing: 0.02em;
}

.btn-toggle:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(246, 189, 96, 0.35);
  background: #f7c570;
}

.btn-toggle:active:not(:disabled) {
  transform: translateY(0);
}

.btn-toggle:disabled {
  background: rgba(13, 19, 33, 0.1);
  color: rgba(13, 19, 33, 0.4);
  cursor: not-allowed;
}

/* ===== Feedback Messages ===== */
.message {
  margin-top: 16px;
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 600;
  letter-spacing: 0.01em;
}

.message.success {
  color: #16a34a;
  background: rgba(22, 163, 74, 0.08);
  border: 1px solid rgba(22, 163, 74, 0.2);
}

.message.error {
  color: #d32f2f;
  background: rgba(211, 47, 47, 0.08);
  border: 1px solid rgba(211, 47, 47, 0.2);
}

@media (max-width: 600px) {
  .change-profile-privacy h2 {
    font-size: 1.125rem;
  }

  .btn-toggle {
    width: 100%;
  }
}
</style>
