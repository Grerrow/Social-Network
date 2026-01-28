<!-- filepath: /home/kali/Desktop/social-network/frontend/src/components/chat/chatinput.vue -->
<template>
  <div class="chat-input-container">
    <div class="input-row">
      <!-- Image upload button -->
      <button 
        class="action-btn" 
        title="Attach image" 
        type="button"
        @click="triggerImageUpload"
      >
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
          <circle cx="8.5" cy="8.5" r="1.5"/>
          <polyline points="21 15 16 10 5 21"/>
        </svg>
      </button>

      <input 
        ref="imageInput" 
        type="file" 
        accept="image/jpeg,image/png,image/gif" 
        style="display: none"
        @change="handleImageUpload"
      />
      
      <div class="input-wrapper">
        <input
          v-model="message"
          @keydown.enter.prevent="send"
          placeholder="Type a message..."
          class="chat-input"
          ref="inputRef"
        />
        
        <button 
          class="emoji-btn" 
          title="Emoji" 
          type="button"
          @click="showEmojiPicker = !showEmojiPicker"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <path d="M8 14s1.5 2 4 2 4-2 4-2"/>
            <line x1="9" y1="9" x2="9.01" y2="9"/>
            <line x1="15" y1="9" x2="15.01" y2="9"/>
          </svg>
        </button>
      </div>
      
      <button 
        class="send-btn" 
        @click="send"
        :disabled="!message.trim() && !pendingImage"
        :class="{ 'active': message.trim() || pendingImage }"
        title="Send message"
        type="button"
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
          <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"/>
        </svg>
      </button>
    </div>

    <!-- Emoji picker dropdown -->
    <div v-if="showEmojiPicker" class="emoji-picker">
      <div class="emoji-grid">
        <button 
          v-for="emoji in emojis" 
          :key="emoji"
          type="button"
          class="emoji-btn-item"
          @click="insertEmoji(emoji)"
        >
          {{ emoji }}
        </button>
      </div>
    </div>

    <!-- Image preview -->
    <div v-if="pendingImage" class="image-preview">
      <img :src="pendingImage.preview" :alt="pendingImage.name" />
      <button 
        type="button" 
        class="remove-image-btn"
        @click="removePendingImage"
      >
        ×
      </button>
    </div>

    <!-- Upload progress -->
    <div v-if="isUploading" class="upload-progress">
      <div class="spinner"></div>
      <span>Uploading...</span>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useChatStore } from '../../stores/chat'
import api from '../../services/api'

const props = defineProps({ 
  chatId: Number,
  groupId: Number
})

const chatStore = useChatStore()
const message = ref('')
const inputRef = ref(null)
const imageInput = ref(null)
const showEmojiPicker = ref(false)
const pendingImage = ref(null)
const isUploading = ref(false)

const emojis = [
  '😀', '😂', '😍', '🥰', '😘', '😜', '😎',
  '👍', '👏', '🙌', '🤝', '❤️', '🔥', '💯',
  '😢', '😡', '😱', '😴', '🤔', '🙏', '✨'
]

const triggerImageUpload = () => {
  imageInput.value?.click()
}

const handleImageUpload = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return

  if (!['image/jpeg', 'image/png', 'image/gif'].includes(file.type)) {
    alert('Only JPEG, PNG, and GIF images are allowed')
    imageInput.value.value = ''
    return
  }

  if (file.size > 5 * 1024 * 1024) {
    alert('Image must be less than 5MB')
    imageInput.value.value = ''
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    pendingImage.value = {
      file,
      name: file.name,
      preview: e.target?.result
    }
  }
  reader.readAsDataURL(file)
}

const removePendingImage = () => {
  pendingImage.value = null
  if (imageInput.value) imageInput.value.value = ''
}

const insertEmoji = (emoji) => {
  message.value += emoji
  showEmojiPicker.value = false
  inputRef.value?.focus()
}

const send = async () => {
  if (!message.value.trim() && !pendingImage.value) return

  try {
    let imageUrl = null

    if (pendingImage.value) {
      isUploading.value = true
      const uploadResult = await api.uploadImage(pendingImage.value.file, 'message')
      imageUrl = uploadResult.image_url
      isUploading.value = false
    }

    const content = message.value.trim() || imageUrl
    
    // Send to group or private
    if (props.groupId) {
      await chatStore.sendGroupMessage(props.groupId, content)
    } else {
      await chatStore.sendMessage(props.chatId, content)
    }

    message.value = ''
    removePendingImage()
    showEmojiPicker.value = false
    inputRef.value?.focus()
  } catch (error) {
    console.error('Failed to send message:', error)
    isUploading.value = false
    alert('Failed to send message')
  }
}
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.chat-input-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  background: white;
  border-top: 1px solid rgba(13, 19, 33, 0.08);
  flex-shrink: 0;
}

.input-row {
  display: flex;
  align-items: flex-end;
  gap: 8px;
}

.action-btn,
.emoji-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--muted-teal);
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.action-btn:hover,
.emoji-btn:hover {
  background: var(--lavender-mist);
  color: #7da99c;
}

.input-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  background: var(--lavender-mist);
  border-radius: 12px;
  padding-right: 4px;
  min-width: 0;
  transition: background-color 0.2s ease;
}

.input-wrapper:focus-within {
  background: rgba(246, 240, 249, 0.7);
}

.chat-input {
  flex: 1;
  padding: 8px 12px;
  border: none;
  background: transparent;
  font-size: 0.875rem;
  outline: none;
  color: var(--ink-black);
  min-width: 0;
  font-family: inherit;
}

.chat-input::placeholder {
  color: rgba(13, 19, 33, 0.5);
}

.send-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: rgba(13, 19, 33, 0.1);
  border-radius: 50%;
  cursor: not-allowed;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(13, 19, 33, 0.4);
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.send-btn.active {
  background: var(--honey-bronze);
  color: var(--ink-black);
  cursor: pointer;
}

.send-btn.active:hover {
  background: #f7c570;
  transform: scale(1.05);
}

.send-btn svg {
  margin-left: 2px;
}

/* Emoji Picker */
.emoji-picker {
  position: absolute;
  bottom: 80px;
  right: 12px;
  background: white;
  border: 2px solid rgba(13, 19, 33, 0.08);
  border-radius: 12px;
  padding: 8px;
  box-shadow: 0 4px 16px rgba(13, 19, 33, 0.15);
  z-index: 1000;
}

.emoji-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
}

.emoji-btn-item {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  font-size: 1.2rem;
  cursor: pointer;
  border-radius: 8px;
  transition: background-color 0.2s;
}

.emoji-btn-item:hover {
  background: var(--lavender-mist);
}

/* Image Preview */
.image-preview {
  position: relative;
  width: 100%;
  max-width: 150px;
}

.image-preview img {
  width: 100%;
  height: auto;
  border-radius: 12px;
  border: 2px solid rgba(13, 19, 33, 0.08);
  max-height: 150px;
  object-fit: cover;
}

.remove-image-btn {
  position: absolute;
  top: -8px;
  right: -8px;
  width: 24px;
  height: 24px;
  background: #d32f2f;
  color: white;
  border: none;
  border-radius: 50%;
  font-size: 1.2rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  font-weight: 700;
  border: 2px solid white;
}

.remove-image-btn:hover {
  background: #b71c1c;
  transform: scale(1.1);
}

/* Upload Progress */
.upload-progress {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.875rem;
  color: rgba(13, 19, 33, 0.6);
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(13, 19, 33, 0.1);
  border-top-color: var(--honey-bronze);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>