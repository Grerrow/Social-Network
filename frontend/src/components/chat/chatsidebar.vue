<template>
  <div class="chat-sidebar">
    <div class="chat-search">
      <input 
        type="text" 
        placeholder="Search contacts or groups..." 
        class="search-input"
        v-model="searchQuery"
      />
    </div>

    <div class="chat-header">
      <h3>Contacts</h3>
    </div>

    <!-- Private chats (top) -->
    <ul class="chat-list">
      <li 
        v-for="conv in filteredPrivate" 
        :key="conv.id" 
        @click="chatStore.openChat?.(conv.id)"
        class="chat-item"
        :class="{ 'has-unread': conv.unreadCount > 0 }"
      >
        <div class="avatar-wrapper">
          <img 
            v-if="conv.avatar" 
            :src="getImageUrl(conv.avatar)" 
            class="avatar"
          />
          <img 
            v-else
            src="/images/placeholder.png"
            alt="Placeholder Avatar"
            class="avatar"
          />
          <span v-if="conv.online" class="online-indicator"></span>
          <!-- Unread badge with count -->
          <span v-if="conv.unreadCount > 0" class="unread-badge">
            {{ conv.unreadCount > 99 ? '99+' : conv.unreadCount }}
          </span>
        </div>
        <div class="chat-info">
          <span class="username">{{ conv.username }}</span>
          <span v-if="conv.online" class="status online">Active now</span>
          <span v-else class="status">Offline</span>
        </div>
      </li>
      <li v-if="filteredPrivate.length === 0" class="no-results">
        <span v-if="searchQuery">No contacts found</span>
        <span v-else>No contacts yet</span>
      </li>
    </ul>

    <!-- Group chats (bottom) -->
    <div class="chat-header" style="border-top:1px solid rgba(13, 19, 33, 0.08);">
      <h3>Groups</h3>
    </div>
    <ul class="chat-list">
      <li 
        v-for="group in filteredGroups" 
        :key="group.id" 
        @click="chatStore.openGroupChat?.(group.id)"
        class="chat-item"
      >
        <div class="avatar-wrapper">
          <div class="avatar avatar-placeholder group-avatar">
            {{ group.group_name?.charAt(0).toUpperCase() }}
          </div>
        </div>
        <div class="chat-info">
          <span class="username">{{ group.group_name }}</span>
          <span class="status">Group chat</span>
        </div>
      </li>
      <li v-if="filteredGroups.length === 0" class="no-results">
        <span v-if="searchQuery">No groups found</span>
        <span v-else>No groups yet</span>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useChatStore } from '../../stores/chat'
import { useApi } from '../../composables/useApi'
const { getImageUrl } = useApi()

const chatStore = useChatStore()
const searchQuery = ref('')

const filteredPrivate = computed(() => {
  const list = chatStore.conversations || []
  if (!searchQuery.value.trim()) return list
  const q = searchQuery.value.toLowerCase()
  return list.filter(c => c.username?.toLowerCase().includes(q))
})

const filteredGroups = computed(() => {
  const list = chatStore.groups || []
  if (!searchQuery.value.trim()) return list
  const q = searchQuery.value.toLowerCase()
  // Filter by group_name property, not the object itself
  return list.filter(g => g.group_name?.toLowerCase().includes(q))
})

onMounted(() => chatStore.fetchConversations())
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.chat-sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.08);
  overflow: hidden;
}

.chat-header {
  padding: 16px;
  border-bottom: 1px solid rgba(13, 19, 33, 0.08);
  flex-shrink: 0;
}

.chat-header h3 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.01em;
}

.chat-search {
  padding: 12px;
  border-bottom: 1px solid rgba(13, 19, 33, 0.08);
  flex-shrink: 0;
}

.search-input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  background: white;
  border-radius: 12px;
  font-size: 0.875rem;
  outline: none;
  box-sizing: border-box;
  transition: all 0.2s ease;
  color: var(--ink-black);
  font-family: inherit;
}

.search-input:focus {
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

.search-input::placeholder {
  color: rgba(13, 19, 33, 0.5);
}

.chat-list {
  list-style: none;
  margin: 0;
  padding: 8px 0;
  overflow-y: auto;
  flex: 1;
}

.chat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  cursor: pointer;
  transition: background-color 0.15s;
}

.chat-item:hover {
  background: var(--lavender-mist);
}

.chat-item.has-unread {
  background: rgba(246, 189, 96, 0.08);
}

.avatar-wrapper {
  position: relative;
  flex-shrink: 0;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 40px;
  height: 40px;
  background: var(--ink-black);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 1rem;
  color: white;
}

.online-indicator {
  position: absolute;
  bottom: 1px;
  right: 1px;
  width: 10px;
  height: 10px;
  background: #31a24c;
  border: 2px solid white;
  border-radius: 50%;
}

.unread-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 12px;
  background: #d32f2f;
  color: white;
  font-size: 0.6875rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 0 2px white;
  letter-spacing: -0.02em;
}

.chat-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
  flex: 1;
  gap: 2px;
}

.username {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--ink-black);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  letter-spacing: 0.01em;
}

.status {
  font-size: 0.75rem;
  color: rgba(13, 19, 33, 0.5);
}

.status.online {
  color: #31a24c;
  font-weight: 600;
}

.no-results {
  padding: 32px 24px;
  text-align: center;
  color: rgba(13, 19, 33, 0.5);
  font-size: 0.875rem;
}
</style>