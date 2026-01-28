import { defineStore } from 'pinia'
import api from '../services/api'

export const useNotificationsStore = defineStore('notifications', {
  state: () => ({
    notifications: [],
    isLoading: false
  }),

  getters: {
    unreadCount: (state) =>
      Array.isArray(state.notifications)
        ? state.notifications.filter(n => !n.is_read).length
        : 0
  },

  actions: {
    // Fetch all notifications from backend
    async fetchNotifications() {
      this.isLoading = true
      try {
        const res = await api.getNotifications()
        console.log('[NOTIFICATIONS] Fetched:', res)
        // Handle both { notifications: [] } and direct array response
        if (Array.isArray(res)) {
          //sort by created_at desc
          const sorted = res.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
          this.notifications = sorted
        } else if (res && Array.isArray(res.notifications)) {
          const sorted = res.notifications.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
          this.notifications = sorted
        } else {
          this.notifications = []
        }
      } catch (err) {
        console.error('Failed to fetch notifications', err)
        this.notifications = []
      } finally {
        this.isLoading = false
      }
    },

    // Mark all notifications as read
    async markAllAsRead() {
      try {
        await api.markAllNotificationsAsRead()
        this.notifications.forEach(n => (n.is_read = true))
      } catch (err) {
        console.error('Failed to mark notifications as read', err)
      }
    },

    

    // Remove a notification by ID
    async removeNotification(id) {
      try {
        // i will handle the deletion at frontend only for now
        this.notifications = this.notifications.filter(n => n.id !== id)
        const res = await api.removeNotification(id)
        console.log('[NOTIFICATIONS] Removed:', res)

      } catch (err) {
        console.error('Failed to remove notification', err)
      }
    },

    // Add a new notification (from WebSocket or backend push)
    addNotificationSignal(notification) {
      if (!notification) return
      console.log('[NOTIFICATIONS] Adding signal:', notification)
      // Avoid duplicates
      if (this.notifications.find(n => n.id === notification.id)) return
      this.notifications.unshift(notification)
    }
  }
})