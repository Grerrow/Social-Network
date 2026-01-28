<!-- filepath: /home/kali/Desktop/social-network/frontend/src/views/FriendsView.vue -->
<template>
  <div class="friends-view">
    <div class="friends-container">
      <!-- Header -->
      <div class="friends-header">
        <div v-if="viewingUser" class="header-with-user">
          <div class="header-avatar">
            <img 
              v-if="viewingUser.avatar" 
              :src="getImageUrl(viewingUser.avatar) || '/images/placeholder.png'" 
              :alt="viewingUser.username"
            />
          </div>
          <div>
            <h1>@{{ viewingUser.username }}</h1>
            <p class="subtitle">{{ isOwnProfile ? 'Your' : viewingUser.username + "'s" }} connections</p>
          </div>
        </div>
        <div v-else>
          <h1>Friends & Connections</h1>
          <p class="subtitle">Manage your followers and following</p>
        </div>
      </div>

      <!-- Search Bar -->
      <div class="search-section">
        <div class="search-input-wrapper">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="search-icon">
            <circle cx="11" cy="11" r="8"/>
            <path d="m21 21-4.35-4.35"/>
          </svg>
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Search by username or name..."
            class="search-input"
          />
          <button 
            v-if="searchQuery" 
            @click="searchQuery = ''"
            class="clear-search-btn"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Tabs -->
      <div class="friends-tabs">
        <button 
          :class="['tab', { active: activeTab === 'followers' }]"
          @click="activeTab = 'followers'"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
          Followers
          <span class="tab-count">{{ followers.length }}</span>
        </button>
        <button 
          :class="['tab', { active: activeTab === 'following' }]"
          @click="activeTab = 'following'"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="8.5" cy="7" r="4"/>
            <line x1="20" y1="8" x2="20" y2="14"/>
            <line x1="23" y1="11" x2="17" y2="11"/>
          </svg>
          Following
          <span class="tab-count">{{ following.length }}</span>
        </button>
      </div>

      <!-- User Lists -->
      <div class="users-list">
        <!-- Followers Tab -->
        <div v-if="activeTab === 'followers'">
          <div v-if="filteredFollowers.length === 0" class="empty-state">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
            </svg>
            <p>{{ searchQuery ? 'No followers found' : 'No followers yet' }}</p>
            <span class="empty-hint">{{ searchQuery ? 'Try a different search term' : 'People who follow will appear here' }}</span>
          </div>
          
          <div 
            v-for="user in filteredFollowers" 
            :key="user.id" 
            class="user-card"
          >
            <router-link :to="`/profile?user_id=${user.id}`" class="user-info">
              <div class="user-avatar">
                <img 
                  :src="user.avatar ? getImageUrl(user.avatar) : '/images/placeholder.png'"
                  :alt="user.username"
                />
              </div>
              <div class="user-details">
                <span class="user-name">
                  {{ user.first_name && user.last_name ? `${user.first_name} ${user.last_name}` : user.username }}
                </span>
                <span class="user-username">@{{ user.username }}</span>
              </div>
            </router-link>
            
            <!-- Only show action buttons if viewing own profile -->
            <template v-if="isOwnProfile">
              <button 
                v-if="!isFollowing(user.id)"
                @click="handleFollow(user.id)"
                class="btn-action btn-follow"
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"/>
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Follow Back
              </button>
              <button 
                v-else
                @click="handleUnfollow(user.id)"
                class="btn-action btn-following"
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
                Following
              </button>
            </template>
          </div>
        </div>

        <!-- Following Tab -->
        <div v-if="activeTab === 'following'">
          <div v-if="filteredFollowing.length === 0" class="empty-state">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="8.5" cy="7" r="4"/>
              <line x1="20" y1="8" x2="20" y2="14"/>
              <line x1="23" y1="11" x2="17" y2="11"/>
            </svg>
            <p>{{ searchQuery ? 'No users found' : 'Not following anyone yet' }}</p>
            <span class="empty-hint">{{ searchQuery ? 'Try a different search term' : 'Start following people to see them here' }}</span>
          </div>
          
          <div 
            v-for="user in filteredFollowing" 
            :key="user.id" 
            class="user-card"
          >
            <router-link :to="`/profile?user_id=${user.id}`" class="user-info">
              <div class="user-avatar">
                <img 
                  v-if="user.avatar" 
                  :src="getImageUrl(user.avatar) || '/images/placeholder.png'" 
                  :alt="user.username"
                />
              </div>
              <div class="user-details">
                <span class="user-name">
                  {{ user.first_name && user.last_name ? `${user.first_name} ${user.last_name}` : user.username }}
                </span>
                <span class="user-username">@{{ user.username }}</span>
              </div>
            </router-link>
            
            <!-- Only show unfollow button if viewing own profile -->
            <button 
              v-if="isOwnProfile"
              @click="handleUnfollow(user.id)"
              class="btn-action btn-unfollow"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
              Unfollow
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useFollowStatus } from '@/composables/useFollowStatus'
import { useApi } from '@/composables/useApi'

const { getImageUrl, getApiUrl } = useApi()

const route = useRoute()

const activeTab = ref('followers')
const searchQuery = ref('')
const followers = ref([])
const following = ref([])
const viewingUser = ref(null)
const currentUserId = ref(null)
const isLoading = ref(true)

// Initialize useFollowStatus composable
const { sendFollowRequest, unfollow } = useFollowStatus()

// Check if viewing own profile
const isOwnProfile = computed(() => {
  if (!route.query.user_id) return true
  return currentUserId.value && parseInt(route.query.user_id) === currentUserId.value
})

// Filtered lists based on search
const filteredFollowers = computed(() => {
  if (!searchQuery.value.trim()) return followers.value
  
  const query = searchQuery.value.toLowerCase()
  return followers.value.filter(user => 
    user.username?.toLowerCase().includes(query) ||
    user.first_name?.toLowerCase().includes(query) ||
    user.last_name?.toLowerCase().includes(query) ||
    `${user.first_name} ${user.last_name}`.toLowerCase().includes(query)
  )
})

const filteredFollowing = computed(() => {
  if (!searchQuery.value.trim()) return following.value
  
  const query = searchQuery.value.toLowerCase()
  return following.value.filter(user => 
    user.username?.toLowerCase().includes(query) ||
    user.first_name?.toLowerCase().includes(query) ||
    user.last_name?.toLowerCase().includes(query) ||
    `${user.first_name} ${user.last_name}`.toLowerCase().includes(query)
  )
})

// Check if currently following a user
const isFollowing = (userId) => {
  return following.value.some(user => user.id === userId)
}

// Fetch profile data (includes followers & following)
const fetchProfileData = async () => {
  isLoading.value = true
  try {
    const userId = route.query.user_id
    const url = userId 
      ? getApiUrl(`/profile?user_id=${userId}`)
      : getApiUrl('/profile')
    
    const res = await fetch(url, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch profile')
    
    const profile = await res.json()
    
    // Store current user's ID if viewing own profile
    if (profile.owner) {
      currentUserId.value = profile.id
    }
    
    // Populate data
    followers.value = profile.followers || []
    following.value = profile.following || []
    viewingUser.value = profile
    
    // Set active tab from query param
    if (route.query.tab === 'followers' || route.query.tab === 'following') {
      activeTab.value = route.query.tab
    }
  } catch (error) {
    console.error('Failed to fetch profile data:', error)
    followers.value = []
    following.value = []
  } finally {
    isLoading.value = false
  }
}

// Handle follow action - CONNECTED to useFollowStatus
const handleFollow = async (userId) => {
  try {
    await sendFollowRequest(userId)
    // Refresh data to update UI
    await fetchProfileData()
  } catch (error) {
    console.error('Failed to follow user:', error)
  }
}

// Handle unfollow action - CONNECTED to useFollowStatus
const handleUnfollow = async (userId) => {
  try {
    await unfollow(userId)
    // Refresh data to update UI
    await fetchProfileData()
  } catch (error) {
    console.error('Failed to unfollow user:', error)
  }
}

// Watch for route changes
watch(() => [route.query.user_id, route.query.tab], () => {
  fetchProfileData()
}, { deep: true })

onMounted(() => {
  fetchProfileData()
})
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.friends-view {
  max-width: 800px;
  margin: 0 auto;
  padding: 32px 16px;
  font-family: 'Inter', 'Segoe UI', system-ui, -apple-system, sans-serif;
}

.friends-container {
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  overflow: hidden;
}

/* Header */
.friends-header {
  background: linear-gradient(135deg, var(--ink-black) 0%, #1a2332 100%);
  color: white;
  padding: 32px;
  text-align: center;
}

.header-with-user {
  display: flex;
  align-items: center;
  gap: 20px;
  justify-content: center;
}

.header-avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  overflow: hidden;
  background: var(--lavender-mist);
  border: 3px solid white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  flex-shrink: 0;
}

.header-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.header-avatar .avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--ink-black);
}

.friends-header h1 {
  margin: 0 0 8px 0;
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.subtitle {
  margin: 0;
  opacity: 0.9;
  font-size: 0.9375rem;
}

/* Search Section */
.search-section {
  padding: 24px 24px 16px;
  border-bottom: 2px solid rgba(13, 19, 33, 0.08);
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 16px;
  color: rgba(13, 19, 33, 0.5);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 14px 48px 14px 48px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.9375rem;
  font-family: inherit;
  color: var(--ink-black);
  background: white;
  transition: all 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

.search-input::placeholder {
  color: rgba(13, 19, 33, 0.5);
}

.clear-search-btn {
  position: absolute;
  right: 12px;
  width: 28px;
  height: 28px;
  border: none;
  background: rgba(13, 19, 33, 0.08);
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--ink-black);
  transition: all 0.2s ease;
}

.clear-search-btn:hover {
  background: rgba(13, 19, 33, 0.15);
  transform: scale(1.1);
}

/* Tabs */
.friends-tabs {
  display: flex;
  background: white;
  padding: 0 24px;
  border-bottom: 2px solid rgba(13, 19, 33, 0.08);
}

.tab {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 16px 20px;
  background: transparent;
  border: none;
  border-bottom: 3px solid transparent;
  font-size: 0.9375rem;
  font-weight: 600;
  color: rgba(13, 19, 33, 0.6);
  cursor: pointer;
  transition: all 0.2s ease;
  letter-spacing: 0.01em;
}

.tab:hover {
  color: var(--ink-black);
  background: rgba(246, 240, 249, 0.5);
}

.tab.active {
  color: var(--ink-black);
  border-bottom-color: var(--honey-bronze);
  font-weight: 700;
}

.tab svg {
  flex-shrink: 0;
}

.tab-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 24px;
  padding: 0 8px;
  background: var(--lavender-mist);
  color: var(--ink-black);
  font-size: 0.75rem;
  font-weight: 700;
  border-radius: 12px;
  letter-spacing: 0.02em;
}

.tab.active .tab-count {
  background: var(--honey-bronze);
}

/* User List */
.users-list {
  padding: 24px;
  min-height: 400px;
}

.user-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  margin-bottom: 12px;
  background: var(--lavender-mist);
  border-radius: 12px;
  transition: all 0.2s ease;
  border: 2px solid transparent;
}

.user-card:hover {
  background: rgba(246, 240, 249, 0.7);
  border-color: rgba(146, 191, 177, 0.3);
  transform: translateX(4px);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
  text-decoration: none;
  color: inherit;
  min-width: 0;
}

.user-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  background: linear-gradient(135deg, var(--muted-teal), var(--honey-bronze));
  border: 3px solid white;
  box-shadow: 0 2px 8px rgba(13, 19, 33, 0.1);
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 1.25rem;
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.user-name {
  font-size: 0.9375rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: 0.01em;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-username {
  font-size: 0.8125rem;
  color: rgba(13, 19, 33, 0.6);
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Action Buttons */
.btn-action {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 18px;
  border: none;
  border-radius: 10px;
  font-size: 0.875rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  letter-spacing: 0.01em;
  flex-shrink: 0;
  white-space: nowrap;
}

.btn-action svg {
  flex-shrink: 0;
}

.btn-follow {
  background: var(--honey-bronze);
  color: var(--ink-black);
}

.btn-follow:hover {
  background: #f7c570;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(246, 189, 96, 0.35);
}

.btn-following {
  background: var(--muted-teal);
  color: white;
}

.btn-following:hover {
  background: #7da99c;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(146, 191, 177, 0.35);
}

.btn-unfollow {
  background: rgba(211, 47, 47, 0.1);
  color: #d32f2f;
  border: 2px solid transparent;
}

.btn-unfollow:hover {
  background: #d32f2f;
  color: white;
  border-color: #d32f2f;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(211, 47, 47, 0.3);
}

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 32px;
  text-align: center;
  color: rgba(13, 19, 33, 0.6);
}

.empty-state svg {
  margin-bottom: 16px;
  opacity: 0.4;
}

.empty-state p {
  margin: 0 0 8px 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: rgba(13, 19, 33, 0.7);
}

.empty-hint {
  font-size: 0.875rem;
  color: rgba(13, 19, 33, 0.5);
}

/* Responsive */
@media (max-width: 768px) {
  .friends-view {
    padding: 16px 12px;
  }

  .friends-header {
    padding: 24px 20px;
  }

  .header-with-user {
    flex-direction: column;
    gap: 16px;
  }

  .friends-header h1 {
    font-size: 1.5rem;
  }

  .search-section {
    padding: 20px 16px 12px;
  }

  .friends-tabs {
    padding: 0 16px;
  }

  .tab {
    font-size: 0.875rem;
    padding: 14px 12px;
    gap: 6px;
  }

  .tab svg {
    width: 16px;
    height: 16px;
  }

  .users-list {
    padding: 16px;
  }

  .user-card {
    padding: 14px 16px;
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .user-info {
    width: 100%;
  }

  .btn-action {
    width: 100%;
    justify-content: center;
  }

  .user-avatar {
    width: 48px;
    height: 48px;
  }

  .avatar-placeholder {
    font-size: 1.125rem;
  }
}

@media (max-width: 480px) {
  .search-input {
    padding: 12px 40px 12px 40px;
    font-size: 0.875rem;
  }

  .tab {
    flex-direction: column;
    gap: 4px;
    font-size: 0.8125rem;
  }

  .tab-count {
    min-width: 20px;
    height: 20px;
    font-size: 0.6875rem;
  }

  .empty-state {
    padding: 48px 24px;
  }

  .empty-state svg {
    width: 48px;
    height: 48px;
  }
}
</style>