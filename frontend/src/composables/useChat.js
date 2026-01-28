import { useChatStore } from '../stores/chat'
import { computed } from 'vue'

export default function useChat() {
  const chatStore = useChatStore()

  const openChats = computed(() => chatStore.openChats)
  const openChat = chatStore.openChat
  const closeChat = chatStore.closeChat

  return {
    openChats,
    openChat,
    closeChat
  }
}