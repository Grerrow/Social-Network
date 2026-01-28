import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const isLoading = ref(false)
  const error = ref(null)
  const checked = ref(false) // new: backend check done

  const isAuthenticated = computed(() => !!user.value)

  const setUser = (userData) => (user.value = userData)
  const setLoading = (loading) => (isLoading.value = loading)
  const setError = (err) => (error.value = err)
  const clearError = () => (error.value = null)

  // New: check backend session on app start / refresh
  const checkAuth = async () => {
    setLoading(true)
    try {
      const res = await fetch(`${api.API_BASE_URL}/me`, {
        credentials: 'include',
        cache: 'no-cache' // Prevent caching
      })
      if (!res.ok) {
        user.value = null
      } else {
        const data = await res.json()
        // Handle nested user object from /me endpoint
        if (data.user && data.user.id) {
          user.value = { ...data.user } // Create new object to ensure reactivity
        } else if (data.user_id) {
          user.value = { id: data.user_id }
        } else {
          user.value = null
        }
      }
    } catch {
      user.value = null
    } finally {
      checked.value = true
      setLoading(false)
    }
  }

  const login = async (credentials) => {
    setLoading(true)
    clearError()
    try {
      const response = await api.login(credentials)
      // After successful login, fetch the full user profile
      await checkAuth()
      return response
    } catch (err) {
      setError(err.message || 'Login failed')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const register = async (userData) => {
  setLoading(true)
  clearError()
  try {
    const { confirmPassword, ...registrationData } = userData
    await api.register(registrationData)
    checked.value = true
  } catch (err) {
    setError(err.message || 'Registration failed')
    throw err
  } finally {
    setLoading(false)
  }
}


  const logout = async () => {
    setLoading(true)
    clearError()
    try {
      await api.logout()
      user.value = null
      checked.value = false // Reset checked so it will re-check on next auth
    } catch (err) {
      setError(err.message || 'Logout failed')
      throw err
    } finally {
      setLoading(false)
    }
  }

  return {
    user,
    isLoading,
    error,
    checked,
    isAuthenticated,
    setUser,
    setLoading,
    setError,
    clearError,
    login,
    register,
    logout,
    checkAuth
  }
})
