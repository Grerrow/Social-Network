import router from '../router'
import { useAuthStore } from '../stores/auth'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// centralized fetch helper
async function fetchWithAuth(url, options = {}) {
  const authStore = useAuthStore()
  const res = await fetch(`${API_BASE_URL}${url}`, {
    ...options,
    credentials: 'include'
  })

  if (res.status === 401) {
    router.push('/login')
    throw new Error('Unauthorized')
  }

  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'API error')
  }

  const contentType = res.headers.get('content-type')
  if (contentType && contentType.includes('application/json')) {
    return res.json()
  }
  return res.text()
}

const api = {
  API_BASE_URL,

  // Authorization
  register: (userData) =>
    fetchWithAuth('/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(userData)
    }),

  login: (credentials) =>
    fetchWithAuth('/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials)
    }),

  logout: () => fetchWithAuth('/logout', { method: 'POST' }),

  //notifications
  getNotifications: () => fetchWithAuth('/notifications'),

  markAllNotificationsAsRead: () =>
    fetchWithAuth('/notifications/read-all', { method: 'POST' }),

  // Posts
  createPost: (postData) =>
    fetchWithAuth('/posts', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(postData)
    }),

  getPosts: () => fetchWithAuth('/posts'),

  getPost: (postId) => fetchWithAuth(`/post?post_id=${postId}`,{
    method: 'GET',  
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: new URLSearchParams({ post_id: parseInt(postId, 10) })
  }),

  createComment: (commentData) =>
    fetchWithAuth('/comments', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(commentData)
    }),

  // Groups
  getGroups: () => fetchWithAuth('/groups'),

  createGroup: (data) =>
    fetchWithAuth('/groups', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    }),

  getGroup: (groupId) => fetchWithAuth(`/groups/details?id=${groupId}`),

  getGroupMembers: (groupId) => fetchWithAuth(`/groups/members?id=${groupId}`),

  //get users for inviting to group
  getUsers: () => fetchWithAuth('/users'),

  // Group Membership
  inviteToGroup: (groupId, userId) =>
    fetchWithAuth('/groups/invite', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ group_id: parseInt(groupId, 10), user_id: userId })
    }),

  acceptGroupInvite: (groupId) =>
    fetchWithAuth('/groups/invite/accept', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ group_id: groupId })
    }),

  declineGroupInvite: (groupId) =>
    fetchWithAuth('/groups/invite/decline', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ group_id: groupId })
    }),

  // thats the old school way of sending form data, good to know
//just wanted to try it out
  requestJoinGroup: (groupId) =>
    fetchWithAuth('/groups/request', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: new URLSearchParams({ group_id: parseInt(groupId, 10) })
    }),


  getJoinRequests: (groupId) => fetchWithAuth(`/groups/requests?id=${groupId}`),

  approveJoinRequest: (groupId, requesterId) =>
    fetchWithAuth('/groups/request/approve', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ group_id: parseInt(groupId, 10), user_id: requesterId })
    }),

  rejectJoinRequest: (groupId, requesterId) =>
    fetchWithAuth('/groups/request/reject', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ group_id: parseInt(groupId, 10), user_id: requesterId })
    }),

  // Group Posts - send group_id in body as string
  getGroupPosts: (groupId) => fetchWithAuth(`/groups/posts?id=${groupId}`),

  createGroupPost: (groupId, data) =>
    fetchWithAuth('/groups/posts', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        group_id: String(groupId),  // backend expects string
        content: data.content,
        image: data.image || ''
      })
    }),

  // Group Comments
  getGroupComments: (groupId, postId) => 
    fetchWithAuth(`/groups/comments?group_id=${groupId}&post_id=${postId}`),

  createGroupComment: (groupId, postId, content, image) =>
    fetchWithAuth('/groups/comments', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        group_id: groupId,
        post_id: postId, 
        content,
        image: image || ''  // optional image field
      })
    }),

  // Group Events
  getGroupEvents: (groupId) => fetchWithAuth(`/groups/events?id=${groupId}`),

  createGroupEvent: (groupId, data) =>
    fetchWithAuth('/groups/events', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        group_id: parseInt(groupId, 10),  // backend expects int
        title: data.title,
        description: data.description,
        event_date: data.event_date       // already Unix timestamp from frontend
      })
    }),



  respondToEvent: (groupId, eventId, response) =>
    fetchWithAuth('/groups/events/respond', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        group_id: groupId,
        event_id: eventId, 
        response 
      })
    }),

  leaveGroup: (groupId) =>
    fetchWithAuth('/groups/leave', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ group_id: parseInt(groupId, 10) })
    }),

  // Chat
  sendMessage: (messageData) =>
    fetchWithAuth('/chat/messages', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(messageData)
    }),

  sendGroupMessage: (data) =>
    fetchWithAuth('/chat/group-messages', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    }),

  getMessages: (chatId) => fetchWithAuth(`/chat/messages?user_id=${chatId}`),

  getGroupMessages: (groupId) => fetchWithAuth(`/chat/group-messages?group_id=${groupId}`),

  getConversations: () => fetchWithAuth('/chat/conversations'),

  getFollowers: () => fetchWithAuth('/followers'),

    //remove notification by id
  removeNotification: (id) =>
    fetchWithAuth(`/notifications/remove?id=${id}`, {
      method: 'DELETE'
    }),

  // Uploads
  async uploadImage(file, type) {
    const formData = new FormData()
    formData.append('image', file)
    formData.append('image_type', type)

    const res = await fetch(`${API_BASE_URL}/image/upload`, {
      method: 'POST',
      body: formData,
      credentials: 'include'
    })

    if (!res.ok) {
      const text = await res.text()
      throw new Error(text || 'Upload failed')
    }

    return res.json()
  }
}

export default api
