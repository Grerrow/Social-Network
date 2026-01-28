<template>
  <div class="settings-view">
    <h1>Settings</h1>

    <div class="settings-container">
      <!-- Sidebar Navigation -->
      <nav class="settings-nav">
        <button
          v-for="tab in tabs"
          :key="tab"
          @click="activeTab = tab"
          class="tab-button"
          :class="{ active: activeTab === tab }"
        >
          {{ tab }}
        </button>
      </nav>

      <!-- Content Area -->
      <div class="settings-content">
        <!-- Account Tab -->
        <div v-if="activeTab === 'Account'" class="settings-section">
          <h2>Account Settings</h2>
          <form @submit.prevent="saveSettings">
            <div class="form-group">
              <label for="email">Email</label>
              <input
                id="email"
                v-model="settings.email"
                type="email"
                placeholder="your@email.com"
              />
            </div>

            <div class="form-group">
              <label for="username">Username</label>
              <input
                id="username"
                v-model="settings.username"
                type="text"
                placeholder="Your username"
              />
            </div>

            <button type="submit" class="btn-save">Save Changes</button>
          </form>
        </div>

        <!-- Privacy Tab -->
        <div v-if="activeTab === 'Privacy'" class="settings-section">
          <h2>Privacy Settings</h2>
          <ChangeProfilePrivacy />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import ChangeProfilePrivacy from '@/components/changeProfilePrivacy.vue'
import { useApi } from '@/composables/useApi'

const { getApiUrl } = useApi()

const tabs = ['Account', 'Privacy']
const activeTab = ref('Account')
const settings = ref({
  email: 'user@example.com',
  username: 'username',
})

const saveSettings = async () => {
  try {
    const response = await fetch(getApiUrl('/profile/update'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(settings.value)
    })
    if (response.ok) {
      alert('Settings saved successfully')
    } else {
      alert('Failed to save settings')
    }
  } catch (error) {
    console.error('Error saving settings:', error)
    alert('Error saving settings')
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

/* ===== Layout ===== */
.settings-view {
  max-width: 960px;
  margin: 0 auto;
  padding: 32px 24px;
}

.settings-view h1 {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--ink-black);
  margin-bottom: 32px;
  letter-spacing: -0.02em;
}

/* ===== Container ===== */
.settings-container {
  display: flex;
  gap: 0;
  background-color: white;
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  overflow: hidden;
}

/* ===== Sidebar ===== */
.settings-nav {
  width: 220px;
  background-color: var(--lavender-mist);
  border-right: 1px solid rgba(13, 19, 33, 0.08);
  display: flex;
  flex-direction: column;
  padding: 16px 0;
}

/* ===== Tabs ===== */
.tab-button {
  padding: 12px 24px;
  border: none;
  background: transparent;
  text-align: left;
  cursor: pointer;
  font-size: 0.9375rem;
  font-weight: 600;
  color: rgba(13, 19, 33, 0.7);
  transition: all 0.15s ease;
  border-left: 3px solid transparent;
  letter-spacing: 0.01em;
}

.tab-button:hover {
  background: rgba(13, 19, 33, 0.05);
  color: var(--ink-black);
}

.tab-button.active {
  background: white;
  color: var(--ink-black);
  border-left-color: var(--honey-bronze);
  font-weight: 700;
}

/* ===== Content ===== */
.settings-content {
  flex: 1;
  padding: 32px;
}

.settings-section h2 {
  margin: 0 0 24px 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.01em;
}

/* ===== Forms ===== */
.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: var(--ink-black);
  font-size: 0.875rem;
  letter-spacing: 0.01em;
}

.form-group input[type='text'],
.form-group input[type='email'] {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.9375rem;
  transition: all 0.2s ease;
  box-sizing: border-box;
  font-family: inherit;
  color: var(--ink-black);
  background: white;
}

.form-group input[type='text']:focus,
.form-group input[type='email']:focus {
  outline: none;
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

/* ===== Primary Button ===== */
.btn-save {
  margin-top: 8px;
  padding: 12px 32px;
  background: var(--honey-bronze);
  color: var(--ink-black);
  border: none;
  border-radius: 12px;
  font-size: 0.9375rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  letter-spacing: 0.02em;
}

.btn-save:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(246, 189, 96, 0.35);
  background: #f7c570;
}

.btn-save:active {
  transform: translateY(0);
}

/* ===== Responsive ===== */
@media (max-width: 768px) {
  .settings-view {
    padding: 24px 16px;
  }

  .settings-container {
    flex-direction: column;
  }

  .settings-nav {
    width: 100%;
    flex-direction: row;
    overflow-x: auto;
    border-right: none;
    border-bottom: 1px solid rgba(13, 19, 33, 0.08);
  }

  .tab-button {
    white-space: nowrap;
    border-left: none;
    border-bottom: 3px solid transparent;
  }

  .tab-button.active {
    border-left-color: transparent;
    border-bottom-color: var(--honey-bronze);
  }

  .settings-content {
    padding: 24px;
  }
}
</style>
