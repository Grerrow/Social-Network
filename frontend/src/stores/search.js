import { defineStore } from 'pinia'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export const useSearchStore = defineStore('search', {
  state: () => ({
    query: '',
    users: [],
    groups: [],
    isLoading: false
  }),

  actions: {
    async search(query) {
      this.query = query.trim()

      if (!this.query) {
        this.users = []
        this.groups = []
        return
      }

      this.isLoading = true

      try {
        const res = await fetch(
          `${API_BASE_URL}/search?q=${encodeURIComponent(this.query)}`,
          { credentials: 'include' }
        )

        const data = await res.json()
        this.users = data.users || []
        this.groups = data.groups || []
      } catch (err) {
        console.error('Search failed:', err)
      } finally {
        this.isLoading = false
      }
    },

    clear() {
      this.query = ''
      this.users = []
      this.groups = []
    }
  }
})
