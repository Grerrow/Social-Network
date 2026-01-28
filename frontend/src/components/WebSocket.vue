<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { useNotificationsStore } from '@/stores/notifications'
import { useChatStore } from '@/stores/chat'
import { useAuthStore } from '@/stores/auth'

const store = useNotificationsStore()
const chatStore = useChatStore()
const authStore = useAuthStore()
const socket = ref(null)
let pingInterval = null
let reconnectTimeout = null

// Watch for logout
watch(() => authStore.user, (newUser, oldUser) => {
  if (oldUser && !newUser) {
    // User logged out
    console.log('[WS] User logged out, resetting chat')
    chatStore.reset()
    disconnect()
  } else if (!oldUser && newUser) {
    // User logged in
    console.log('[WS] User logged in, connecting')
    connect()
  }
}, { immediate: false })
const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'



const connect = () => {
  if (!authStore.user) return
  
  if (socket.value) {
    socket.value.close()
    socket.value = null
  }

  console.log('[WS] Connecting...')
  socket.value = new WebSocket(WS_URL)

  socket.value.onopen = () => {
    console.log('[WS] Connected')
    chatStore.initCurrentUser()
    chatStore.loadOpenChats()
    startPing()
  }

  socket.value.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data)
      console.log('[WS] Received:', message.type)

      switch (message.type) {
        case 'pong':
          break
        case 'notification':
          if (window.location.pathname === '/notifications') {
            store.fetchNotifications()
          } else {
            store.addNotificationSignal(message.data)
          }
          break
        case 'private_message':
          chatStore.addMessage(message.data)
          break
          case 'group_message':
          chatStore.addGroupMessage(message.data)
          break
        case 'new_message_notification':
          // update the chat in chat sidebar notification indicator
          chatStore.notifyChat(message.data)
          break
        case 'presence:init':
          chatStore.setOnlineUsers(message.data)
          break
        case 'presence':
          chatStore.updateUserPresence(message.data.user_id, message.data.online)
          break
        case 'group_invitation':
          store.addNotificationSignal(message.data)
          break
        case 'group_join_request':
          store.addNotificationSignal(message.data)
          break
          case 'group_join_request_approved':
          store.addNotificationSignal(message.data)
          break
        case 'event_created':
          store.addNotificationSignal(message.data)
          break
      }
    } catch (e) {
      console.error('[WS] Parse error:', e)
    }
  }

  socket.value.onclose = (event) => {
    console.log('[WS] Closed:', event.code)
    stopPing()
    chatStore.clearOnlineUsers()
    chatStore.setWebSocket(null)
    
    // Only reconnect if user is still logged in
    if (authStore.user) {
      reconnectTimeout = setTimeout(() => {
        console.log('[WS] Reconnecting...')
        connect()
      }, 2000)
    }
  }

  socket.value.onerror = (err) => {
    console.error('[WS] Error:', err)
  }
}

const disconnect = () => {
  stopPing()
  if (reconnectTimeout) {
    clearTimeout(reconnectTimeout)
    reconnectTimeout = null
  }
  if (socket.value) {
    socket.value.close()
    socket.value = null
  }
}

const startPing = () => {
  stopPing()
  pingInterval = setInterval(() => {
    if (socket.value && socket.value.readyState === WebSocket.OPEN) {
      socket.value.send(JSON.stringify({ type: 'ping' }))
    }
  }, 25000)
}

const stopPing = () => {
  if (pingInterval) {
    clearInterval(pingInterval)
    pingInterval = null
  }
}

onMounted(() => {
  if (authStore.user) {
    connect()
  }
})

onBeforeUnmount(() => {
  disconnect()
})
</script>