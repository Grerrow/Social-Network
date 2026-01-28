import { computed } from 'vue'

export function useApi() {
  // Get base URLs from environment variables
  const apiBaseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'
  
  // Helper to build full image URL
  const getImageUrl = (path) => {
    if (!path) return ''
    if (path.startsWith('http')) return path // Already full URL
    return `${apiBaseUrl}${path}`
  }
  
  // Helper to build API endpoint URL
  const getApiUrl = (endpoint) => {
    if (!endpoint) return apiBaseUrl
    if (endpoint.startsWith('http')) return endpoint
    return `${apiBaseUrl}${endpoint}`
  }

  return {
    apiBaseUrl: computed(() => apiBaseUrl),
    wsUrl: computed(() => wsUrl),
    getImageUrl,
    getApiUrl
  }
}