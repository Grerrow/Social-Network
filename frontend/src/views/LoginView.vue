<template>
  <div class="login-view">
    <div class="login-card">
      <div class="card-header">
        <h2>Welcome Back</h2>
        <p class="subtitle">Sign in to continue</p>
      </div>
      
      <form @submit.prevent="handleLogin">
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
        
        <div class="form-group">
          <label for="password">Password</label>
          <div class="input-wrapper">
            <span class="input-icon">🔒</span>
            <input
              id="password"
              v-model="form.password"
              type="password"
              required
              placeholder="Enter your password"
            />
          </div>
        </div>
        
        <button type="submit" :disabled="authStore.isLoading" class="btn-submit">
          <span v-if="authStore.isLoading" class="loading-spinner"></span>
          {{ authStore.isLoading ? 'Signing in...' : 'Sign In' }}
        </button>
      </form>
      
      <p v-if="authStore.error" class="error-message">Wrong email or password</p>
      
      <p class="register-link">
        Don't have an account? <router-link to="/register">Create one</router-link>
      </p>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

export default {
  name: 'LoginView',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const form = ref({
      email: '',
      password: ''
    })

    const handleLogin = async () => {
      try {
        await authStore.login(form.value)
        router.push('/feed')
      } catch (err) {
        console.error('Login error:', err)
      }
    }

    return {
      form,
      authStore,
      handleLogin
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

.login-view {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 80px);
  padding: 32px 16px;
  font-family: 'Inter', 'Segoe UI', system-ui, -apple-system, sans-serif;
  background: var(--lavender-mist);
}

.login-card {
  width: 100%;
  max-width: 420px;
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

.form-group input {
  width: 100%;
  padding: 12px 16px;
  padding-left: 48px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.9375rem;
  transition: all 0.2s ease;
  box-sizing: border-box;
  font-family: inherit;
  color: var(--ink-black);
  background: white;
}

.form-group input:focus {
  outline: none;
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
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

.register-link {
  text-align: center;
  margin-top: 24px;
  color: rgba(13, 19, 33, 0.6);
  font-size: 0.9375rem;
}

.register-link a {
  color: var(--muted-teal);
  text-decoration: none;
  font-weight: 600;
  transition: color 0.2s ease;
}

.register-link a:hover {
  color: #7da99c;
  text-decoration: underline;
}

@media (max-width: 500px) {
  .login-card {
    padding: 24px;
  }
}
</style>