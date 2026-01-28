/**
 * Date formatting utility
 */
export const formatDate = (date) => {
  if (!date) return 'just now'
  return new Date(date).toLocaleDateString()
}

export const formatTime = (date) => {
  if (!date) return 'just now'
  return new Date(date).toLocaleTimeString()
}

export const formatDateTime = (date) => {
  if (!date) return 'just now'
  return new Date(date).toLocaleString()
}

/**
 * String utilities
 */
export const truncate = (str, length = 100) => {
  if (!str) return ''
  return str.length > length ? str.substring(0, length) + '...' : str
}

export const capitalize = (str) => {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1)
}

/**
 * Validation utilities
 */
export const isValidEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

export const isValidPassword = (password) => {
  // At least 8 characters, 1 uppercase, 1 lowercase, 1 number
  return password.length >= 8 && /[A-Z]/.test(password) && /[a-z]/.test(password) && /[0-9]/.test(password)
}
// check for date validation not older than today or future date
export const isValidDate = (date) => {
  if (!date) return false
  const inputDate = new Date(date)
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const minDate = new Date('1940-01-01')
  return inputDate >= minDate && inputDate <= today
}

/**
 * Storage utilities
 */
export const setLocal = (key, value) => {
  localStorage.setItem(key, JSON.stringify(value))
}

export const getLocal = (key) => {
  const item = localStorage.getItem(key)
  return item ? JSON.parse(item) : null
}

export const removeLocal = (key) => {
  localStorage.removeItem(key)
}

/**
 * API utilities
 */
export const getErrorMessage = (error) => {
  if (error.message) return error.message
  if (error.error) return error.error
  return 'An error occurred'
}
