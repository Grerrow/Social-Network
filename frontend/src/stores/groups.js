import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const useGroupsStore = defineStore('groups', () => {
  const groups = ref([])
  const currentGroup = ref(null)
  const groupPosts = ref([])
  const groupEvents = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  // Fetch all groups
  async function fetchGroups() {
    isLoading.value = true
    error.value = null
    try {
      const data = await api.getGroups()
      groups.value = data || []
    } catch (e) {
      error.value = e.message
      groups.value = []
    } finally {
      isLoading.value = false
    }
  }

  // Fetch single group details
  async function fetchGroup(groupId) {
    isLoading.value = true
    error.value = null
    try {
      currentGroup.value = await api.getGroup(groupId)
    } catch (e) {
      error.value = e.message
      currentGroup.value = null
    } finally {
      isLoading.value = false
    }
  }

  // Create a new group
  async function createGroup(group_name, title, description) {
    isLoading.value = true
    error.value = null
    try {
      const newGroup = await api.createGroup({ group_name, title, description })
      groups.value.unshift(newGroup)
      return newGroup
    } catch (e) {
      error.value = e.message
      return null
    } finally {
      isLoading.value = false
    }
  }

  // Invite user to group
  async function inviteUser(groupId, userId) {
    try {
      return await api.inviteToGroup(groupId, userId)
    } catch (e) {
      error.value = e.message
      return null
    }
  }

  // Accept group invitation
  async function acceptInvite(groupId) {
    try {
      await api.acceptGroupInvite(groupId)
      await fetchGroup(groupId)
      return true
    } catch (e) {
      error.value = e.message
      return false
    }
  }

  // Decline group invitation
  async function declineInvite(groupId) {
    try {
      await api.declineGroupInvite(groupId)
      return true
    } catch (e) {
      error.value = e.message
      return false
    }
  }

  // Request to join group
  async function requestJoin(groupId) {
    try {
      await api.requestJoinGroup(groupId)
      await fetchGroup(groupId)
      return true
    } catch (e) {
      error.value = e.message
      return false
    }
  }

  // Approve join request (creator only)
  async function approveJoinRequest(groupId, requesterId) {
    try {
      await api.approveJoinRequest(groupId, requesterId)
      await fetchGroup(groupId)
      return true
    } catch (e) {
      error.value = e.message
      return false
    }
  }

  // Reject join request (creator only)
  async function rejectJoinRequest(groupId, userId) {
    try {
      await api.rejectJoinRequest(groupId, userId)
      await fetchGroup(groupId)
      return true
    } catch (e) {
      error.value = e.message
      return false
    }
  }

  // Fetch group posts
  async function fetchGroupPosts(groupId) {
    try {
      const data = await api.getGroupPosts(groupId)
      groupPosts.value = data || []
    } catch (e) {
      error.value = e.message
      groupPosts.value = []
    }
  }

  // Create group post
  async function createGroupPost(groupId, content, imageFile = null) {
    try {
      const postData = { content: content || '' }
      
      // Upload image first if provided
      if (imageFile) {
        const uploadResult = await api.uploadImage(imageFile, 'post')
        postData.image = uploadResult.image_url
      }
      
      const newPost = await api.createGroupPost(groupId, postData)
      if (newPost) {
        groupPosts.value.unshift(newPost)
      }
      return newPost
    } catch (e) {
      error.value = e.message
      console.error('[GROUPS] createGroupPost error:', e)
      return null
    }
  }

  // Fetch comments for a post
  async function fetchPostComments(groupId, postId) {
    try {
      const data = await api.getGroupComments(groupId, postId)
      return data || []
    } catch (e) {
      error.value = e.message
      return []
    }
  }

  // Create comment on post
  async function createPostComment(groupId, postId, content, image) {
    try {
      return await api.createGroupComment(groupId, postId, content, image)
    } catch (e) {
      error.value = e.message
      return null
    }
  }

  // Fetch group events
  async function fetchGroupEvents(groupId) {
    try {
      const data = await api.getGroupEvents(groupId)
      groupEvents.value = data || []
    } catch (e) {
      error.value = e.message
      groupEvents.value = []
    }
  }

  // Create group event
  async function createGroupEvent(groupId, eventData) {
    try {
      await api.createGroupEvent(groupId, eventData)
      await fetchGroupEvents(groupId)
    } catch (e) {
      error.value = e.message
      return null
    }
  }

  async function leaveGroup(groupId) {
    try {
      await api.leaveGroup(groupId)
      // Remove group from list
      groups.value = groups.value.filter(g => g.id !== groupId)
      // Clear current group if it was the one left
      if (currentGroup.value && currentGroup.value.id === groupId) {
        currentGroup.value = null
        groupPosts.value = []
        groupEvents.value = []
      }
      return true
    } catch (e) {
      error.value = e.message
      return false
    }
  }

  // Respond to event
  async function respondToEvent(groupId, eventId, response) {
    try {
      await api.respondToEvent(groupId, eventId, response)
      await fetchGroupEvents(groupId)
      return true
    } catch (e) {
      error.value = e.message
      return false
    }
  }

  function reset() {
    groups.value = []
    currentGroup.value = null
    groupPosts.value = []
    groupEvents.value = []
    error.value = null
  }

  return {
    groups,
    currentGroup,
    groupPosts,
    groupEvents,
    isLoading,
    error,
    fetchGroups,
    fetchGroup,
    createGroup,
    inviteUser,
    acceptInvite,
    declineInvite,
    requestJoin,
    approveJoinRequest,
    rejectJoinRequest,
    fetchGroupPosts,
    createGroupPost,
    fetchPostComments,
    createPostComment,
    fetchGroupEvents,
    createGroupEvent,
    leaveGroup,
    respondToEvent,
    reset
  }
})
