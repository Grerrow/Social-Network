<template>
  <div class="register-view">
    <div class="register-card">
      <div class="card-header">
        <h2>Create Account</h2>
        <p class="subtitle">Join our community today</p>
      </div>
      
      <form @submit.prevent="handleRegister">

        <!-- Username -->
        <div class="form-group">
          <label for="username">Nickname (Optional)</label>
          <div class="input-wrapper">
            <span class="input-icon">👤</span>
            <input
              id="username"
              v-model="form.username"
              type="text"
              placeholder="Choose a username"
            />
          </div>
        </div>

        <!-- Email -->
        <div class="form-group">
          <label for="email">Email</label>
          <div class="input-wrapper">
            <span class="input-icon">📧</span>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              placeholder="Enter your email"
            />
          </div>
        </div>

        <!-- Name Row -->
        <div class="form-row">
          <div class="form-group">
            <label for="firstName">First Name</label>
            <input
              id="firstName"
              v-model="form.first_name"
              type="text"
              required
              placeholder="First name"
            />
          </div>

          <div class="form-group">
            <label for="lastName">Last Name</label>
            <input
              id="lastName"
              v-model="form.last_name"
              type="text"
              required
              placeholder="Last name"
            />
          </div>
        </div>

        <!-- Date of Birth - Professional Date Picker -->
        <div class="form-group">
          <label>Date of Birth</label>
          <div class="date-picker-row">
            <select v-model="dateForm.month" required class="date-select">
              <option value="" disabled>Month</option>
              <option v-for="(month, index) in months" :key="index" :value="index + 1">
                {{ month }}
              </option>
            </select>
            
            <select v-model="dateForm.day" required class="date-select">
              <option value="" disabled>Day</option>
              <option v-for="day in daysInMonth" :key="day" :value="day">
                {{ day }}
              </option>
            </select>
            
            <select v-model="dateForm.year" required class="date-select">
              <option value="" disabled>Year</option>
              <option v-for="year in years" :key="year" :value="year">
                {{ year }}
              </option>
            </select>
          </div>
        </div>

        <!-- Avatar Upload - Professional Card Style -->
        <div class="form-group">
          <label>Profile Picture (Optional)</label>
          <div class="avatar-upload-container">
            <input
              ref="avatarInput"
              type="file"
              accept="image/*"
              @change="handleAvatarUpload"
              hidden
            />
            
            <div class="avatar-preview-area">
              <div v-if="!form.avatar" class="avatar-empty">
                <div class="avatar-icon">
                  <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                  </svg>
                </div>
                <span class="avatar-label">No photo</span>
              </div>
              
              <div v-else class="avatar-preview-img">
                <img :src="form.avatar" alt="Avatar preview" />
              </div>
            </div>
            
            <div class="avatar-actions">
              <button 
                v-if="!form.avatar"
                type="button" 
                class="btn-upload-avatar"
                @click="triggerAvatarInput"
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                  <polyline points="17 8 12 3 7 8"/>
                  <line x1="12" y1="3" x2="12" y2="15"/>
                </svg>
                Upload Photo
              </button>
              
              <button 
                v-else
                type="button" 
                class="btn-change-avatar"
                @click="triggerAvatarInput"
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                  <polyline points="17 8 12 3 7 8"/>
                  <line x1="12" y1="3" x2="12" y2="15"/>
                </svg>
                Change
              </button>
              
              <button 
                v-if="form.avatar"
                type="button" 
                class="btn-remove-avatar"
                @click="clearAvatar"
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
                Remove
              </button>
            </div>
          </div>
        </div>

        <!-- About Me -->
        <div class="form-group">
          <label for="aboutMe">About Me (Optional)</label>
          <textarea
            id="aboutMe"
            v-model="form.about_me"
            rows="3"
            placeholder="Tell us about yourself"
          ></textarea>
        </div>

        <!-- Password -->
        <div class="form-group">
          <label for="password">Password</label>
          <div class="input-wrapper">
            <span class="input-icon">🔒</span>
            <input
              id="password"
              v-model="form.password"
              type="password"
              required
              placeholder="Create a password"
            />
          </div>
        </div>

        <!-- Confirm Password -->
        <div class="form-group">
          <label for="confirmPassword">Confirm Password</label>
          <div class="input-wrapper">
            <span class="input-icon">🔒</span>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              required
              placeholder="Confirm your password"
            />
          </div>
        </div>

        <!-- Submit -->
        <button type="submit" :disabled="authStore.isLoading" class="btn-submit">
          <span v-if="authStore.isLoading" class="loading-spinner"></span>
          {{ authStore.isLoading ? 'Creating account...' : 'Create Account' }}
        </button>

      </form>

      <!-- Error Message -->
      <p v-if="authStore.error" class="error-message">{{ authStore.error }}</p>

      <!-- Login link -->
      <p class="login-link">
        Already have an account? <router-link to="/login">Sign in</router-link>
      </p>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { isValidDate, isValidEmail } from '../utils/helpers'
import api from '../services/api'
import { useApi } from '@/composables/useApi'

const { getImageUrl, getApiUrl } = useApi()

export default {
  name: 'RegisterView',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const avatarInput = ref(null)

    const form = ref({
      username: '',
      email: '',
      first_name: '',
      last_name: '',
      date_of_birth: '',
      password: '',
      confirmPassword: '',
      avatar: '',
      about_me: '',
      is_private: false
    })

    // Date picker state
    const dateForm = ref({
      month: '',
      day: '',
      year: ''
    })

    const months = [
      'January', 'February', 'March', 'April', 'May', 'June',
      'July', 'August', 'September', 'October', 'November', 'December'
    ]

    // Generate years (100 years back from 18 years ago)
    const currentYear = new Date().getFullYear()
    const years = computed(() => {
      const arr = []
      for (let i = currentYear - 18; i >= currentYear - 100; i--) {
        arr.push(i)
      }
      return arr
    })

    // Calculate days in selected month/year
    const daysInMonth = computed(() => {
      if (!dateForm.value.month || !dateForm.value.year) return 31
      const month = parseInt(dateForm.value.month)
      const year = parseInt(dateForm.value.year)
      return new Date(year, month, 0).getDate()
    })

    // Watch for date changes and update form.date_of_birth
    watch(dateForm, (newDate) => {
      if (newDate.year && newDate.month && newDate.day) {
        const month = String(newDate.month).padStart(2, '0')
        const day = String(newDate.day).padStart(2, '0')
        form.value.date_of_birth = `${newDate.year}-${month}-${day}`
      }
    }, { deep: true })

    const triggerAvatarInput = () => {
      avatarInput.value.click()
    }

    const resizeImage = (file, maxWidth = 300, maxHeight = 300) => {
      return new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.onload = (event) => {
          const img = new Image()
          img.onload = () => {
            const canvas = document.createElement('canvas')
            let width = img.width
            let height = img.height

            if (width > height) {
              if (width > maxWidth) {
                height = (height * maxWidth) / width
                width = maxWidth
              }
            } else {
              if (height > maxHeight) {
                width = (width * maxHeight) / height
                height = maxHeight
              }
            }

            canvas.width = width
            canvas.height = height
            const ctx = canvas.getContext('2d')
            ctx.drawImage(img, 0, 0, width, height)

            canvas.toBlob(
              (blob) => {
                if (blob) resolve(new File([blob], file.name, { type: file.type }))
                else reject(new Error('Failed to convert image'))
              },
              file.type,
              0.9
            )
          }
          img.src = event.target.result
        }
        reader.onerror = () => reject(new Error('Failed to read file'))
        reader.readAsDataURL(file)
      })
    }

    const handleAvatarUpload = async (e) => {
      const file = e.target.files[0]
      if (!file) return

      authStore.clearError?.()

      try {
        const resizedFile = await resizeImage(file)
        const res = await api.uploadImage(resizedFile, 'avatar')

        if (!res.image_url) {
          throw new Error('No image URL returned')
        }

        form.value.avatar = getApiUrl(res.image_url)
      } catch (err) {
        console.error(err)
        form.value.avatar = ''
        authStore.setError('Failed to upload avatar')
      }
    }

    const clearAvatar = () => {
      form.value.avatar = ''
      if (avatarInput.value) {
        avatarInput.value.value = ''
      }
    }

    const handleRegister = async () => {
      authStore.clearError?.()

      if (form.value.password !== form.value.confirmPassword) {
        authStore.setError('Passwords do not match')
        return
      }

      if (!isValidEmail(form.value.email)) {
        authStore.setError('Invalid email format')
        return
      }

      if (!isValidDate(form.value.date_of_birth)) {
        authStore.setError('Invalid date of birth')
        return
      }

      try {
        await authStore.register({
          ...form.value,
          avatar: form.value.avatar ? form.value.avatar.replace(getApiUrl(''), '') : ''
        })

        router.push('/login')
      } catch (err) {
        console.error('Registration error:', err)
      }
    }

    return {
      form,
      dateForm,
      months,
      years,
      daysInMonth,
      authStore,
      avatarInput,
      triggerAvatarInput,
      handleAvatarUpload,
      clearAvatar,
      handleRegister,
      getApiUrl
    }
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

.register-view {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 80px);
  padding: 32px 16px;
  font-family: 'Inter', 'Segoe UI', system-ui, -apple-system, sans-serif;
  background: var(--lavender-mist);
}

.register-card {
  width: 100%;
  max-width: 480px;
  padding: 32px;
  border-radius: 16px;
  background: white;
  box-shadow: 0 4px 24px rgba(13, 19, 33, 0.1);
}

.card-header {
  text-align: center;
  margin-bottom: 32px;
}

.card-header h2 {
  margin: 0 0 8px 0;
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.02em;
}

.subtitle {
  margin: 0;
  color: rgba(13, 19, 33, 0.6);
  font-size: 0.9375rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

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

.input-wrapper {
  position: relative;
}

.input-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 1.125rem;
  pointer-events: none;
  opacity: 0.7;
}

.form-group input,
.form-group textarea {
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

.input-wrapper input {
  padding-left: 48px;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
  line-height: 1.5;
}

/* Professional Date Picker */
.date-picker-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1.5fr;
  gap: 12px;
}

.date-select {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.9375rem;
  font-family: inherit;
  color: var(--ink-black);
  background: white;
  cursor: pointer;
  transition: all 0.2s ease;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg width='12' height='8' viewBox='0 0 12 8' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M1 1L6 6L11 1' stroke='%230d1321' stroke-width='2' stroke-linecap='round'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 36px;
}

.date-select:focus {
  outline: none;
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

.date-select option {
  padding: 8px;
}

/* Professional Avatar Upload */
.avatar-upload-container {
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  padding: 20px;
  background: white;
  display: flex;
  align-items: center;
  gap: 20px;
}

.avatar-preview-area {
  flex-shrink: 0;
}

.avatar-empty {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  background: var(--lavender-mist);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.avatar-icon {
  color: rgba(13, 19, 33, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-label {
  font-size: 0.75rem;
  color: rgba(13, 19, 33, 0.5);
  font-weight: 600;
}

.avatar-preview-img {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  overflow: hidden;
  border: 3px solid var(--lavender-mist);
  box-shadow: 0 2px 8px rgba(13, 19, 33, 0.1);
}

.avatar-preview-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.avatar-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
}

.btn-upload-avatar,
.btn-change-avatar,
.btn-remove-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 16px;
  border-radius: 10px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
  letter-spacing: 0.01em;
}

.btn-upload-avatar,
.btn-change-avatar {
  background: var(--honey-bronze);
  color: var(--ink-black);
}

.btn-upload-avatar:hover,
.btn-change-avatar:hover {
  background: #f7c570;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(246, 189, 96, 0.3);
}

.btn-remove-avatar {
  background: rgba(211, 47, 47, 0.1);
  color: #d32f2f;
}

.btn-remove-avatar:hover {
  background: rgba(211, 47, 47, 0.15);
  transform: translateY(-1px);
}

.btn-upload-avatar svg,
.btn-change-avatar svg,
.btn-remove-avatar svg {
  flex-shrink: 0;
}

.btn-submit {
  width: 100%;
  padding: 16px;
  margin-top: 8px;
  background: var(--honey-bronze);
  color: var(--ink-black);
  border: none;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.2s ease;
  letter-spacing: 0.02em;
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(246, 189, 96, 0.35);
  background: #f7c570;
}

.btn-submit:active:not(:disabled) {
  transform: translateY(0);
}

.btn-submit:disabled {
  background: rgba(13, 19, 33, 0.1);
  color: rgba(13, 19, 33, 0.4);
  cursor: not-allowed;
}

.loading-spinner {
  width: 18px;
  height: 18px;
  border: 2px solid transparent;
  border-top-color: var(--ink-black);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-message {
  color: #d32f2f;
  font-size: 0.875rem;
  margin-top: 16px;
  padding: 12px 16px;
  background: rgba(211, 47, 47, 0.08);
  border: 1px solid rgba(211, 47, 47, 0.2);
  border-radius: 12px;
  text-align: center;
  font-weight: 500;
}

.login-link {
  text-align: center;
  margin-top: 24px;
  color: rgba(13, 19, 33, 0.6);
  font-size: 0.9375rem;
}

.login-link a {
  color: var(--muted-teal);
  text-decoration: none;
  font-weight: 600;
  transition: color 0.2s ease;
}

.login-link a:hover {
  color: #7da99c;
  text-decoration: underline;
}

@media (max-width: 500px) {
  .register-card {
    padding: 24px;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .date-picker-row {
    grid-template-columns: 1fr;
  }

  .avatar-upload-container {
    flex-direction: column;
    text-align: center;
  }
}
</style>