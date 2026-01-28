<!-- filepath: /home/kali/Desktop/social-network/frontend/src/components/chat/chatmessage.vue -->
<template>
  <div class="chat-messages" :data-chat-id="chatId || groupId">
    <div v-if="chatMessages.length === 0" class="no-messages">
      <div class="no-messages-icon">💬</div>
      <span class="no-messages-text">No messages yet</span>
      <span class="no-messages-hint">Send a message to start the conversation</span>
    </div>
    
    <div
      v-for="msg in chatMessages"
      :key="msg.id"
      class="message-wrapper"
      :class="{ 'own-message': isOwnMessage(msg) }"
    >
      <div class="message-bubble">
        <!-- Image message -->
        <img 
          v-if="isImageUrl(msg.content)"
          :src="getImageUrl(msg.content)"
          :alt="msg.content"
          class="message-image"
          loading="lazy"
        />
        <!-- Text message -->
        <span v-else class="message-text">{{ msg.content }}</span>
      </div>
      <div class="message-time">
        {{ formatTime(msg.created_at || msg.timestamp) }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, watch, nextTick } from 'vue'
import { useChatStore } from '../../stores/chat'
import { useApi } from '../../composables/useApi'
const { getImageUrl } = useApi()

const props = defineProps({ 
  chatId: Number,
  groupId: Number
})

const chatStore = useChatStore()

const chatMessages = computed(() => {
  if (props.groupId) {
    return chatStore.groupMessages[props.groupId] || []
  }
  return chatStore.messages[props.chatId] || []
})

const isOwnMessage = (msg) => {
  return msg.sender_id === chatStore.currentUser?.id
}

const isImageUrl = (content) => {
  return content?.startsWith('/images/')
}

const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(typeof timestamp === 'number' ? timestamp * 1000 : timestamp)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const scrollToBottom = () => {
  nextTick(() => {
    const container = document.querySelector(`[data-chat-id="${props.chatId || props.groupId}"]`)
    if (container) {
      container.scrollTop = container.scrollHeight
    }
  })
}

onMounted(() => {
  if (props.groupId) {
    chatStore.fetchGroupMessages(props.groupId)
  } else if (props.chatId) {
    chatStore.fetchMessages(props.chatId)
  }
  scrollToBottom()
})

watch(chatMessages, () => {
  scrollToBottom()
}, { deep: true })
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: linear-gradient(180deg, #fafbfc 0%, #f7f8f9 100%);
  min-height: 0;
}

/* Custom Scrollbar */
.chat-messages::-webkit-scrollbar {
  width: 5px;
}

.chat-messages::-webkit-scrollbar-track {
  background: transparent;
}

.chat-messages::-webkit-scrollbar-thumb {
  background: rgba(13, 19, 33, 0.15);
  border-radius: 3px;
}

.chat-messages::-webkit-scrollbar-thumb:hover {
  background: rgba(13, 19, 33, 0.25);
}

/* No Messages State */
.no-messages {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  gap: 8px;
  padding: 32px 24px;
  text-align: center;
}

.no-messages-icon {
  font-size: 2.5rem;
  margin-bottom: 8px;
  opacity: 0.4;
  filter: grayscale(0.3);
}

.no-messages-text {
  font-size: 0.9375rem;
  font-weight: 600;
  color: rgba(13, 19, 33, 0.6);
  letter-spacing: 0.01em;
}

.no-messages-hint {
  font-size: 0.8125rem;
  font-weight: 400;
  color: rgba(13, 19, 33, 0.4);
  max-width: 180px;
  line-height: 1.4;
}

/* Message Wrapper */
.message-wrapper {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  max-width: 75%;
  animation: fadeInSlide 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes fadeInSlide {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-wrapper.own-message {
  align-self: flex-end;
  align-items: flex-end;
}

/* Message Bubble */
.message-bubble {
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(13, 19, 33, 0.08);
  transition: box-shadow 0.2s ease;
}

.message-bubble:hover {
  box-shadow: 0 2px 6px rgba(13, 19, 33, 0.12);
}

/* Received Messages */
.message-wrapper:not(.own-message) .message-bubble {
  background: white;
  border-bottom-left-radius: 4px;
}

/* Sent Messages */
.own-message .message-bubble {
  background: var(--muted-teal);
  border-bottom-right-radius: 4px;
}

/* Text Content */
.message-text {
  display: block;
  padding: 10px 14px;
  font-size: 0.875rem;
  line-height: 1.45;
  word-wrap: break-word;
  word-break: break-word;
  white-space: pre-wrap;
  color: var(--ink-black);
  letter-spacing: 0.01em;
}

.own-message .message-text {
  color: white;
}

/* Image Messages */
.message-image {
  display: block;
  max-width: 220px;
  max-height: 180px;
  width: 100%;
  height: auto;
  object-fit: cover;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.message-image:hover {
  transform: scale(1.02);
}

/* Adjust bubble padding for images */
.message-bubble:has(.message-image) {
  padding: 0;
  background: transparent;
  box-shadow: none;
}

.message-bubble:has(.message-image):hover {
  box-shadow: none;
}

.message-bubble .message-image {
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(13, 19, 33, 0.12);
}

.own-message .message-bubble .message-image {
  border-bottom-right-radius: 4px;
}

.message-wrapper:not(.own-message) .message-bubble .message-image {
  border-bottom-left-radius: 4px;
}

/* Message Time */
.message-time {
  font-size: 0.6875rem;
  color: rgba(13, 19, 33, 0.45);
  margin-top: 4px;
  padding: 0 6px;
  font-weight: 500;
  opacity: 0;
  transition: opacity 0.2s ease;
  letter-spacing: 0.02em;
}

.message-wrapper:hover .message-time {
  opacity: 1;
}

/* Always show time on last message */
.message-wrapper:last-child .message-time {
  opacity: 0.8;
}

/* Mobile Adjustments */
@media (max-width: 480px) {
  .chat-messages {
    padding: 12px 8px;
  }

  .message-wrapper {
    max-width: 85%;
  }

  .message-text {
    padding: 9px 12px;
    font-size: 0.8125rem;
  }

  .message-image {
    max-width: 180px;
    max-height: 150px;
  }

  .no-messages {
    padding: 24px 16px;
  }

  .no-messages-icon {
    font-size: 2rem;
  }
}

/* Compact Mode for Multiple Messages */
.message-wrapper + .message-wrapper {
  margin-top: 2px;
}

.message-wrapper + .message-wrapper:not(.own-message):not(:has(.message-bubble + .message-bubble)) {
  margin-top: 1px;
}

.message-wrapper.own-message + .message-wrapper.own-message {
  margin-top: 1px;
}

/* Different senders = more space */
.message-wrapper:not(.own-message) + .message-wrapper.own-message,
.message-wrapper.own-message + .message-wrapper:not(.own-message) {
  margin-top: 12px;
}
</style>