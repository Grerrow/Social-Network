import { defineStore } from 'pinia'
import { ref, reactive, watch } from 'vue'
import api from '../services/api'
import { useAuthStore } from './auth'


const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const OPEN_CHATS_KEY = 'social_network_open_chats'
const OPEN_GROUP_CHATS_KEY = 'social_network_open_group_chats'

export const useChatStore = defineStore('chat', () => {
  const conversations = ref([])     // private contacts
  const groups = ref([])            // group chats
  const openChats = ref([])         // array of user IDs with open DM windows
  const openGroupChats = ref([])    // array of group IDs with open group windows
  const messages = reactive({})     // { [chatId]: Message[] }
  const groupMessages = reactive({}) // { [groupId]: Message[] }
  const currentUser = ref(null)
  const wsInstance = ref(null)

  function initCurrentUser() {
    const authStore = useAuthStore()
    if (authStore.user) {
      currentUser.value = { id: authStore.user.id || authStore.user.user_id }
    }
  }

  function setWebSocket(ws) {
    wsInstance.value = ws
  }

  function loadOpenChats() {
    const authStore = useAuthStore()
    if (!authStore.user) return
    try {
      const userId = authStore.user.id || authStore.user.user_id
      const key = `${OPEN_CHATS_KEY}_${userId}`
      const saved = localStorage.getItem(key)
      if (saved) {
        const parsed = JSON.parse(saved)
        if (Array.isArray(parsed)) {
          openChats.value = parsed
          parsed.forEach(chatId => {
            if (!messages[chatId]) {
              fetchMessages(chatId)
            }
          })
        }
      }

      // Load open group chats
      const groupKey = `${OPEN_GROUP_CHATS_KEY}_${userId}`
      const savedGroups = localStorage.getItem(groupKey)
      if (savedGroups) {
        const parsedGroups = JSON.parse(savedGroups)
        if (Array.isArray(parsedGroups)) {
          openGroupChats.value = parsedGroups
          parsedGroups.forEach(groupId => {
            if (!groupMessages[groupId]) {
              fetchGroupMessages(groupId)
            }
          })
        }
      }
    } catch (e) {
      console.error('[CHAT] Load open chats error:', e)
    }
  }

  function saveOpenChats() {
    const authStore = useAuthStore()
    if (!authStore.user) return
    try {
      const userId = authStore.user.id || authStore.user.user_id
      const key = `${OPEN_CHATS_KEY}_${userId}`
      localStorage.setItem(key, JSON.stringify(openChats.value))

      // Save open group chats
      const groupKey = `${OPEN_GROUP_CHATS_KEY}_${userId}`
      localStorage.setItem(groupKey, JSON.stringify(openGroupChats.value))
    } catch (e) {
      console.error('[CHAT] Save open chats error:', e)
    }
  }

  watch(openChats, () => saveOpenChats(), { deep: true })
  watch(openGroupChats, () => saveOpenChats(), { deep: true })

  function updateUserPresence(userId, online) {
    const id = typeof userId === 'string' ? parseInt(userId, 10) : userId
    const user = conversations.value.find(u => u.id === id)
    if (user) user.online = online
  }

  function setOnlineUsers(userIds) {
    const onlineSet = new Set(userIds.map(id => parseInt(id, 10)))
    conversations.value.forEach(user => { user.online = onlineSet.has(user.id) })
  }

  function clearOnlineUsers() {
    conversations.value.forEach(user => { user.online = false })
  }

  async function fetchConversations() {
    try {
      const res = await fetch(`${API_BASE_URL}/chat/conversations`, { credentials: 'include' })
      const data = await res.json()

      const all = [
        ...(Array.isArray(data.followers) ? data.followers : []),
        ...(Array.isArray(data.following) ? data.following : [])
      ]

      const seen = new Set()
      conversations.value = all
        .filter(u => u?.id && !seen.has(u.id) && (seen.add(u.id), true))
        .map(u => ({ ...u, unreadCount: u.unread_count || 0, online: !!u.online }))

      groups.value = Array.isArray(data.groups) ? data.groups : []
    } catch (e) {
      console.error('[CHAT] fetchConversations failed', e)
      conversations.value = []
      groups.value = []
    }
  }

  async function fetchMessages(chatId) {
    try {
      const fetched = await api.getMessages(chatId)
      messages[chatId] = (fetched || []).sort((a, b) => a.created_at - b.created_at)
      setTimeout(() => scrollToBottom(chatId), 100)
    } catch (e) {
      console.error('[CHAT] Fetch messages error:', e)
      messages[chatId] = []
    }
  }

  async function fetchGroupMessages(groupId) {
    try {
      const fetched = await api.getGroupMessages(groupId)
      groupMessages[groupId] = (fetched || []).sort((a, b) => a.created_at - b.created_at)
      setTimeout(() => scrollToBottom(groupId), 100)
    } catch (e) {
      console.error('[CHAT] Fetch group messages error:', e)
      groupMessages[groupId] = []
    }
  }

  async function sendMessage(chatId, content) {
    try {
      await api.sendMessage({ receiver_id: chatId, content })
      setTimeout(() => scrollToBottom(chatId), 100)
    } catch (e) {
      console.error('[CHAT] Send message error:', e)
    }
  }

  async function sendGroupMessage(groupId, content) {
    try {
      await api.sendGroupMessage({ group_id: groupId, content })
      setTimeout(() => scrollToBottom(groupId), 100)
    } catch (e) {
      console.error('[CHAT] Send group message error:', e)
    }
  }

  function addMessage(msg) {
    if (!currentUser.value) initCurrentUser()
    const chatId = msg.sender_id === currentUser.value?.id ? msg.receiver_id : msg.sender_id
    if (!messages[chatId]) messages[chatId] = []
    if (!messages[chatId].some(m => m.id === msg.id)) {
      messages[chatId].push(msg)
      messages[chatId].sort((a, b) => a.created_at - b.created_at)
      setTimeout(() => scrollToBottom(chatId), 100)
    }
  }

  function addGroupMessage(msg) {
    if (!currentUser.value) initCurrentUser()
    const groupId = msg.group_id
    if (!groupMessages[groupId]) groupMessages[groupId] = []
    if (!groupMessages[groupId].some(m => m.id === msg.id)) {
      groupMessages[groupId].push(msg)
      groupMessages[groupId].sort((a, b) => a.created_at - b.created_at)
      setTimeout(() => scrollToBottom(groupId), 100)
    }
  }

  // Increment unread when receiving a new private message notification
  function notifyChat(data) {
    if (!currentUser.value) initCurrentUser()
    const senderId = typeof data === 'string' ? parseInt(data, 10) : data
    if (!senderId) return
    if (openChats.value.includes(senderId)) return
    const conv = conversations.value.find(u => u.id === senderId)
    if (conv) conv.unreadCount = (conv.unreadCount || 0) + 1
  }

  function scrollToBottom(chatId) {
    const container = document.querySelector(`[data-chat-id="${chatId}"]`)
    if (container) container.scrollTop = container.scrollHeight
  }

  function openChat(chatId) {
    if (!currentUser.value) initCurrentUser()
    if (!openChats.value.includes(chatId)) openChats.value.push(chatId)
    const conv = conversations.value.find(u => u.id === chatId)
    if (conv) conv.unreadCount = 0
    if (!messages[chatId]) fetchMessages(chatId)
    else setTimeout(() => scrollToBottom(chatId), 100)
  }

  function closeChat(chatId) {
    openChats.value = openChats.value.filter(id => id !== chatId)
  }

  function openGroupChat(groupId) {
    if (!currentUser.value) initCurrentUser()
    if (!openGroupChats.value.includes(groupId)) openGroupChats.value.push(groupId)
    if (!groupMessages[groupId]) fetchGroupMessages(groupId)
    else setTimeout(() => scrollToBottom(groupId), 100)
  }

  function closeGroupChat(groupId) {
    openGroupChats.value = openGroupChats.value.filter(id => id !== groupId)
  }

  function reset() {
    conversations.value = []
    groups.value = []
    openChats.value = []
    openGroupChats.value = []
    currentUser.value = null
    wsInstance.value = null
    Object.keys(messages).forEach(k => delete messages[k])
    Object.keys(groupMessages).forEach(k => delete groupMessages[k])
  }

  return {
    conversations,
    groups,
    openChats,
    openGroupChats,
    messages,
    groupMessages,
    currentUser,
    fetchConversations,
    fetchMessages,
    fetchGroupMessages,
    sendMessage,
    sendGroupMessage,
    openChat,
    closeChat,
    openGroupChat,
    closeGroupChat,
    addMessage,
    addGroupMessage,
    notifyChat,
    initCurrentUser,
    scrollToBottom,
    updateUserPresence,
    setOnlineUsers,
    clearOnlineUsers,
    loadOpenChats,
    setWebSocket,
    reset
  }
})