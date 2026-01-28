<template>
  <div class="chat-header" @click="$emit('toggle-minimize')">
    <div class="user-info">
      <div class="avatar-wrapper">
        <!-- Group avatar -->
        <div v-if="groupId" class="avatar avatar-placeholder">
          {{ groupInfo?.group_name?.charAt(0).toUpperCase() || 'G' }}
        </div>
        <!-- User avatar -->
        <template v-else>
          <img 
            class="avatar"
            :src="conversation.avatar ? getImageUrl(conversation.avatar) : '/images/placeholder.png'"
            :alt="conversation?.username || ''"
          />
          <span v-if="conversation?.online" class="online-dot"></span>
        </template>
      </div>
      <div class="user-details">
        <span class="username">{{ displayName }}</span>
        <span v-if="!groupId && conversation?.online" class="status">Active now</span>
        <span v-if="groupId" class="status">Group chat</span>
      </div>
    </div>
    <div class="header-actions">
      <button 
        class="header-btn" 
        @click.stop="$emit('toggle-minimize')" 
        :title="isMinimized ? 'Expand' : 'Minimize'"
      >
        <svg v-if="isMinimized" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
          <path d="M7 14l5-5 5 5H7z"/>
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
          <path d="M7 10l5 5 5-5H7z"/>
        </svg>
      </button>
      <button class="header-btn close-btn" @click.stop="$emit('close')" title="Close">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
          <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useChatStore } from '../../stores/chat'
import { useApi } from '@/composables/useApi'
const { getImageUrl } = useApi()
const props = defineProps({ 
  chatId: Number,
  groupId: Number,
  isMinimized: Boolean
})

defineEmits(['close', 'toggle-minimize'])

const chatStore = useChatStore()

const conversation = computed(() => 
  props.chatId ? chatStore.conversations.find(c => c.id === props.chatId) : null
)

const groupInfo = computed(() => 
  props.groupId ? chatStore.groups.find(g => g.id === props.groupId) : null
)

const displayName = computed(() => {
  if (props.groupId) return groupInfo.value?.group_name || 'Group'
  return conversation.value?.username || 'Chat'
})
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: white;
  border-bottom: 1px solid rgba(13, 19, 33, 0.08);
  cursor: pointer;
  min-height: 48px;
  box-sizing: border-box;
  flex-shrink: 0;
}

.chat-header:hover {
  background: var(--lavender-mist);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
  flex: 1;
}

.avatar-wrapper {
  position: relative;
  flex-shrink: 0;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 32px;
  height: 32px;
  background: var(--ink-black);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.875rem;
  color: white;
}

.online-dot {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 10px;
  height: 10px;
  background: #31a24c;
  border: 2px solid white;
  border-radius: 50%;
}

.user-details {
  display: flex;
  flex-direction: column;
  min-width: 0;
  gap: 1px;
}

.username {
  font-size: 0.875rem;
  font-weight: 700;
  color: var(--ink-black);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  letter-spacing: 0.01em;
}

.status {
  font-size: 0.75rem;
  color: #31a24a;
  font-weight: 600;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 2px;
  flex-shrink: 0;
}

.header-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(13, 19, 33, 0.6);
  transition: all 0.2s ease;
}

.header-btn:hover {
  background: rgba(13, 19, 33, 0.08);
  color: var(--ink-black);
}

.close-btn:hover {
  background: rgba(211, 47, 47, 0.08);
  color: #d32f2f;
}
</style>