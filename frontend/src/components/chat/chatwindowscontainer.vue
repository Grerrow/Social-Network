<template>
  <div class="chat-windows-container">
    <!-- Private chats -->
    <chatwindow
      v-for="chatId in chatStore.openChats"
      :key="'private-' + chatId"
      :chat-id="chatId"
      @close="chatStore.closeChat(chatId)"
    />
    
    <!-- Group chats -->
    <chatwindow
      v-for="groupId in chatStore.openGroupChats"
      :key="'group-' + groupId"
      :group-id="groupId"
      @close="chatStore.closeGroupChat(groupId)"
    />
  </div>
</template>

<script setup>
import { useChatStore } from '../../stores/chat'
import chatwindow from './chatwindow.vue'

const chatStore = useChatStore()
</script>

<style scoped>
.chat-windows-container {
  position: fixed;
  bottom: 0;
  right: 380px;
  display: flex;
  flex-direction: row-reverse;
  align-items: flex-end;
  gap: 8px;
  z-index: 999;
  pointer-events: none;
  padding-right: 12px;
  max-width: calc(100vw - 400px);
}

.chat-windows-container > * {
  pointer-events: all;
  flex-shrink: 0;
}

@media (max-width: 1280px) {
  .chat-windows-container {
    right: 340px;
    max-width: calc(100vw - 360px);
  }
}

@media (max-width: 1024px) {
  .chat-windows-container {
    right: 24px;
    max-width: calc(100vw - 48px);
  }
}

@media (max-width: 768px) {
  .chat-windows-container {
    right: 16px;
    gap: 8px;
    padding-right: 8px;
    max-width: calc(100vw - 32px);
  }
}

@media (max-width: 480px) {
  .chat-windows-container {
    right: 8px;
    gap: 8px;
    padding-right: 8px;
  }
}
</style>