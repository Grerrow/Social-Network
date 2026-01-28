import { ref, computed } from 'vue'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export function useFollowStatus(initialStatus = 'none') {
  const status = ref(initialStatus)
  const isLoading = ref(false)
  const error = ref(null)

  /* ---------- computed ---------- */

  const isFollowing = computed(() => status.value === 'following')
  const isRequested = computed(() => status.value === 'requested')
  const canFollow = computed(() => status.value === 'none')

  const buttonLabel = computed(() => {
    if (isFollowing.value) return 'Following'
    if (isRequested.value) return 'Requested'
    return 'Follow'
  })

  /* ---------- actions ---------- */

  async function sendFollowRequest(userId) {
    await safe(async () => {
      const res = await fetch(
        `${API_BASE_URL}/follow/request?user_id=${userId}`,
        { method: 'POST', credentials: 'include' }
      )
      
      if (!res.ok) {
        throw new Error('Failed to send follow request')
      }
      
      const data = await res.json()
      status.value = data.message.includes('now following')
        ? 'following'
        : 'requested'
    })
  }

  async function unfollow(userId) {
    await safe(async () => {
      const res = await fetch(
        `${API_BASE_URL}/follow/unfollow?user_id=${userId}`,
        { method: 'POST', credentials: 'include' }
      )
      
      if (!res.ok) {
        throw new Error('Failed to unfollow')
      }
      
      status.value = 'none'
    })
  }

  async function cancelFollowRequest(userId) {
    await safe(async () => {
      const res = await fetch(
        `${API_BASE_URL}/follow/request/cancel?user_id=${userId}`,
        { method: 'POST', credentials: 'include' }
      )
      
      if (!res.ok) {
        throw new Error('Failed to cancel request')
      }
      
      status.value = 'none'
    })
  }

  async function acceptFollowRequest(id) {
    await safe(async () => {
      const success = await updateRequestStatus(id, 'approve')
      if (success) {
        status.value = 'following'
      }
    })
  }

  async function rejectFollowRequest(id) {
    await safe(async () => {
      await updateRequestStatus(id, 'reject')
      // Always set to none, even if request was already cancelled
      status.value = 'none'
    })
  }

  async function updateRequestStatus(follow_request_id, action) {
    const res = await fetch(
      `${API_BASE_URL}/follow/request/status?follow_request_id=${follow_request_id}&action=${action}`,
      { method: 'POST', credentials: 'include' }
    )
    
    if (!res.ok) {
      if (res.status === 404) {
        // Follow request was cancelled or doesn't exist
        console.log('Follow request not found - may have been cancelled')
        return false
      }
      throw new Error('Failed to update request status')
    }
    
    return true
  }

  async function safe(fn) {
    if (isLoading.value) {
      // Prevent multiple simultaneous requests
      return
    }
    
    try {
      isLoading.value = true
      error.value = null
      await fn()
    } catch (err) {
      error.value = err.message
      console.error('Follow action error:', err.message)
    } finally {
      isLoading.value = false
    }
  }

  return {
    status,
    isLoading,
    error,
    isFollowing,
    isRequested,
    canFollow,
    buttonLabel,
    sendFollowRequest,
    unfollow,
    cancelFollowRequest,
    acceptFollowRequest,
    rejectFollowRequest
  }
}
