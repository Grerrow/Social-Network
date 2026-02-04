<template>
  <div class="main-layout">
    <!-- WebSocket component mounted once for notifications -->
    <WebSocketVue />

     <nav class="navbar">
      <div class="navbar-content">
        <!-- BRUTALIST C&D LOGO -->
        <div class="navbar-brand">
          <router-link to="/feed" class="brutalist-logo">
            <div class="logo-grid">
              <div class="logo-square logo-c">C</div>
              <div class="logo-square logo-ampersand">&</div>
              <div class="logo-square logo-d">D</div>
            </div>
          </router-link>
        </div>

        <!-- SEARCH (CENTERED) -->
        <div class="navbar-center">
          <div class="search-bar-container">
            <SearchBar 
              ref="searchBarRef" 
              @search="handleSearch" 
              @enter="handleSearchEnter"
              @escape="handleSearchEscape"
            />

            <!-- SEARCH RESULTS -->
            <transition name="fade">
              <div v-if="showResults" class="search-results">
                <template v-if="searchStore.users.length || searchStore.groups.length">
                  <!-- USERS -->
                  <div v-if="searchStore.users.length">
                    <div class="search-section-title">Users</div>
                    <div
                      v-for="user in searchStore.users"
                      :key="'user-' + user.id"
                      class="search-item"
                    >
                      <router-link
                        :to="`/profile?user_id=${user.id}`"
                        class="search-link"
                        @click="handleResultClick"
                      >
                        <div class="search-result-avatar">
                          <img 
                            :src="user.avatar ? getImageUrl(user.avatar) : '/images/placeholder.png'" 
                            :alt="user.username || user.name"
                          />
                        </div>
                        <div class="search-result-info">
                          <div class="search-result-name">{{ user.username || user.name }}</div>
                        </div>
                      </router-link>
                    </div>
                  </div>

                  <!-- GROUPS -->
                  <div v-if="searchStore.groups.length">
                    <div class="search-section-title">Groups</div>
                    <div
                      v-for="group in searchStore.groups"
                      :key="'group-' + group.id"
                      class="search-item"
                    >
                      <router-link
                        :to="`/groups/${group.id}`"
                        class="search-link"
                        @click="handleResultClick"
                      >
                        <span class="search-icon">👥</span>
                        <span>{{ group.title || group.name }}</span>
                      </router-link>
                    </div>
                  </div>
                </template>

                <!-- EMPTY -->
                <div v-else class="search-empty">
                  No results found
                </div>
              </div>
            </transition>
          </div>
        </div>

        <!-- MENU WITH USER INFO -->
        <div class="navbar-menu">
          <!-- USER INFO -->
          <div class="user-info">
            <div class="user-avatar">
              <img 
                
                 :src="authStore.user.avatar ? getImageUrl(authStore.user.avatar) : '/images/placeholder.png'" 
                :alt="authStore.user.first_name"
              />
            </div>
            <div class="user-name">
              {{ authStore.user?.first_name }} {{ authStore.user?.last_name }}
            </div>
          </div>

          <NotificationsIndicator />
          
          <!-- Chat Toggle Button (visible on smaller screens) -->
          <button class="btn-chat-toggle" @click="toggleChat" title="Toggle Chat">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
          </button>
          
          <button class="btn-logout" @click="handleLogout">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
              <polyline points="16 17 21 12 16 7"/>
              <line x1="21" y1="12" x2="9" y2="12"/>
            </svg>
            <span>Logout</span>
          </button>
        </div>
      </div>
    </nav>

    <!-- LAYOUT -->
    <div class="page-container">
      <!-- LEFT NAVIGATION -->
      <aside class="left-sidebar">
        <div class="sidebar-content">
          <router-link to="/feed" class="sidebar-link">
            <span class="sidebar-icon"><svg class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="30" height="30" fill="none" viewBox="0 0 24 24">
  <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m4 12 8-8 8 8M6 10.5V19a1 1 0 0 0 1 1h3v-3a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v3h3a1 1 0 0 0 1-1v-8.5"/>
</svg>
</span>
            <span class="sidebar-text">Feed</span>
          </router-link>
          <router-link to="/profile" class="sidebar-link">
            <span class="sidebar-icon"><svg class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="33" height="33" fill="none" viewBox="0 0 24 24">
  <path stroke="currentColor" stroke-width="2" d="M7 17v1a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1v-1a3 3 0 0 0-3-3h-4a3 3 0 0 0-3 3Zm8-9a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"/>
</svg>
</span>
            <span class="sidebar-text">Profile</span>
          </router-link>
          <router-link to="/notifications" class="sidebar-link">
            <span class="sidebar-icon"><svg class="bell-icon" width="27" height="27" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
      <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
    </svg></span>
            <span class="sidebar-text">Notifications</span>
          </router-link>
          <router-link to="/groups" class="sidebar-link">
            <span class="sidebar-icon"><svg class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="30" height="30" fill="none" viewBox="0 0 24 24">
  <path stroke="currentColor" stroke-linecap="round" stroke-width="2" d="M16 19h4a1 1 0 0 0 1-1v-1a3 3 0 0 0-3-3h-2m-2.236-4a3 3 0 1 0 0-4M3 18v-1a3 3 0 0 1 3-3h4a3 3 0 0 1 3 3v1a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1Zm8-10a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"/>
</svg>
</span>
            <span class="sidebar-text">Groups</span>
          </router-link>
          <router-link to="/settings" class="sidebar-link">
            <span class="sidebar-icon"><svg class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="33" height="33" fill="none" viewBox="0 0 24 24">
  <path stroke="currentColor" stroke-linecap="square" stroke-linejoin="round" stroke-width="2" d="M10 19H5a1 1 0 0 1-1-1v-1a3 3 0 0 1 3-3h2m10 1a3 3 0 0 1-3 3m3-3a3 3 0 0 0-3-3m3 3h1m-4 3a3 3 0 0 1-3-3m3 3v1m-3-4a3 3 0 0 1 3-3m-3 3h-1m4-3v-1m-2.121 1.879-.707-.707m5.656 5.656-.707-.707m-4.242 0-.707.707m5.656-5.656-.707.707M12 8a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"/>
</svg>
</span>
            <span class="sidebar-text">Settings</span>
          </router-link>
        </div>
      </aside>

      <!-- MAIN CONTENT -->
      <main class="main-content">
        <slot />
      </main>

      <!-- RIGHT CHAT SIDEBAR -->
      <aside class="right-sidebar" :class="{ 'chat-open': isChatOpen }">
        <!-- Close button for mobile -->
        <button v-if="isChatOpen" class="btn-close-chat" @click="toggleChat">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
        <chatsidebar />
      </aside>
    </div>

    <!-- Chat Overlay for mobile -->
    <div v-if="isChatOpen" class="chat-overlay" @click="toggleChat"></div>

    <!-- CHAT WINDOWS CONTAINER (Fixed at bottom) -->
    <chatwindowscontainer />
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSearchStore } from '@/stores/search'
import { useChatStore } from '@/stores/chat'

import WebSocketVue from '@/components/WebSocket.vue'
import chatsidebar from '@/components/chat/chatsidebar.vue'
import chatwindowscontainer from '@/components/chat/chatwindowscontainer.vue'
import SearchBar from '@/components/SearchBar.vue'
import NotificationsIndicator from '@/components/NotificationsIndicator.vue'
import { useApi } from '@/composables/useApi'

export default {
  name: 'MainLayout',
  components: {
    chatsidebar,
    chatwindowscontainer,
    SearchBar,
    NotificationsIndicator,
    WebSocketVue
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const searchStore = useSearchStore()
    const chatStore = useChatStore()

    const showResults = ref(false)
    const searchBarRef = ref(null)
    const isChatOpen = ref(false)
    const {getImageUrl} = useApi()
    let timeout = null

    const handleLogout = async () => {
      //close chat windows
      isChatOpen.value = false
      chatStore.reset() // Clear all chat state including open windows
      await authStore.logout()

      router.push('/login')
    }

    const toggleChat = () => {
      isChatOpen.value = !isChatOpen.value
    }

    const handleSearch = (query) => {
      clearTimeout(timeout)

      query = query.trim()
      if (!query) {
        showResults.value = false
        searchStore.clear()
        return
      }

      timeout = setTimeout(async () => {
        await searchStore.search(query)
        showResults.value = true
      }, 300)
    }

    const handleSearchEnter = async (query) => {
      clearTimeout(timeout)
      await searchStore.search(query.trim())
      showResults.value = true
    }

    const handleSearchEscape = () => {
      showResults.value = false
      searchStore.clear()
    }

    const handleResultClick = () => {
      showResults.value = false
      searchStore.clear()
      searchBarRef.value?.clear()
    }

    const handleClickOutside = (e) => {
      if (!e.target.closest('.search-bar-container')) {
        showResults.value = false
      }
    }

    onMounted(() => {
      document.addEventListener('mousedown', handleClickOutside)
    })

    onBeforeUnmount(() => {
      document.removeEventListener('mousedown', handleClickOutside)
    })

    return {
      handleLogout,
      handleSearch,
      handleSearchEnter,
      handleSearchEscape,
      handleResultClick,
      toggleChat,
      isChatOpen,
      showResults,
      searchStore,
      searchBarRef,
      authStore,
      getImageUrl
    }
  }
}
</script>



<style>
/* =========================== GLOBAL =========================== */
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
  --brutalist-black: #000;
  --brutalist-white: #fff;
  --brutalist-blue: #00a4ef;
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  background: var(--lavender-mist);
  overflow-x: hidden;
}

/* =========================== MAIN LAYOUT =========================== */
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* =========================== NAVBAR =========================== */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  width: 100%;
  z-index: 1000;
  background: white;
  border-bottom: 1px solid rgba(13, 19, 33, 0.08);
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  padding: 12px 32px;
}

.navbar-content {
  max-width: 100%;
  margin: 0 auto;
  display: grid;
  grid-template-columns: minmax(120px, 280px) minmax(300px, 1fr) minmax(280px, 380px);
  align-items: center;
  height: 50px;
  gap: 16px;
}

@media (max-width: 1200px) {
  .navbar-content {
    grid-template-columns: auto 1fr auto;
    gap: 12px;
  }
}

/* =========================== BRUTALIST C&D LOGO =========================== */
.navbar-brand {
  display: flex;
  align-items: center;
}

.brutalist-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  width: 110px;
  height: 50px;
  background-color: var(--brutalist-black);
  color: var(--brutalist-white);
  text-decoration: none;
  font-family: Arial, sans-serif;
  font-weight: bold;
  border: 3px solid var(--brutalist-white);
  outline: 3px solid var(--brutalist-black);
  box-shadow: 6px 6px 0 var(--brutalist-blue);
  transition: all 0.1s ease-out;
  padding: 0;
  box-sizing: border-box;
  position: relative;
  overflow: hidden;
}

/* Shine effect on hover */
.brutalist-logo::before {
  content: "";
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.8),
    transparent
  );
  z-index: 1;
  transition: none;
  opacity: 0;
}

@keyframes slide {
  0% {
    left: -100%;
  }
  100% {
    left: 100%;
  }
}

.brutalist-logo:hover::before {
  opacity: 1;
  animation: slide 2s infinite;
}

.brutalist-logo:hover {
  transform: translate(-4px, -4px);
  box-shadow: 10px 10px 0 var(--brutalist-black);
  background-color: var(--brutalist-black);
  color: var(--brutalist-white);
}

.brutalist-logo:active {
  transform: translate(4px, 4px);
  box-shadow: 0px 0px 0 var(--brutalist-blue);
  background-color: var(--brutalist-white);
  color: var(--brutalist-black);
  border-color: var(--brutalist-black);
}

/* Logo Grid - 3 squares for C & D */
.logo-grid {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  position: relative;
  z-index: 2;
  transition: transform 0.2s ease-out;
}

.brutalist-logo:hover .logo-grid {
  transform: rotate(-5deg) scale(1.05);
}

.brutalist-logo:active .logo-grid {
  transform: rotate(5deg) scale(0.95);
}

.logo-square {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 900;
  color: var(--brutalist-white);
  text-transform: uppercase;
  letter-spacing: -0.05em;
}

.logo-c {
  background-color: #f25022; /* Microsoft red */
}

.logo-ampersand {
  background-color: #7fba00; /* Microsoft green */
  font-size: 16px;
}

.logo-d {
  background-color: #00a4ef; /* Microsoft blue */
}

/* Active state - invert colors */
.brutalist-logo:active .logo-square {
  color: var(--brutalist-black);
}

/* Center Search */
.navbar-center {
  display: flex;
  justify-content: center;
  align-items: center;
  min-width: 0;
}

@media (max-width: 900px) {
  .navbar-center {
    max-width: 400px;
  }
}

/* Menu with User Info */
.navbar-menu {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  min-width: 0;
  flex-shrink: 0;
}

@media (max-width: 1200px) {
  .navbar-menu {
    gap: 8px;
  }
}

/* User Info */
.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 14px 6px 6px;
  border-radius: 12px;
  background: var(--lavender-mist);
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  border: 2px solid white;
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
  background: linear-gradient(135deg, var(--muted-teal), var(--honey-bronze));
  color: white;
  font-weight: 700;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-name {
  font-size: 0.875rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: 0.01em;
}

/* Search bar */
.search-bar-container {
  position: relative;
  width: 100%;
  max-width: 600px;
}

.search-results {
  position: absolute;
  top: calc(100% + 12px);
  left: 0;
  right: 0;
  background: white;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(13, 19, 33, 0.15);
  border: 1px solid rgba(13, 19, 33, 0.08);
  max-height: 450px;
  overflow-y: auto;
  z-index: 1001;
  animation: slideDown 0.2s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.search-results::-webkit-scrollbar {
  width: 8px;
}

.search-results::-webkit-scrollbar-track {
  background: rgba(13, 19, 33, 0.03);
  border-radius: 0 16px 16px 0;
}

.search-results::-webkit-scrollbar-thumb {
  background: rgba(13, 19, 33, 0.15);
  border-radius: 4px;
}

.search-results::-webkit-scrollbar-thumb:hover {
  background: rgba(13, 19, 33, 0.25);
}

.search-section-title {
  padding: 16px 20px 12px;
  font-weight: 700;
  font-size: 0.75rem;
  color: rgba(13, 19, 33, 0.6);
  border-bottom: 1px solid rgba(13, 19, 33, 0.08);
  letter-spacing: 0.08em;
  text-transform: uppercase;
  background: rgba(246, 240, 249, 0.5);
}

.search-item {
  border-bottom: 1px solid rgba(13, 19, 33, 0.04);
  transition: all 0.2s ease;
}

.search-item:last-child {
  border-bottom: none;
}

.search-item:hover {
  background: linear-gradient(to right, rgba(246, 240, 249, 0.6), rgba(246, 189, 96, 0.08));
}

.search-link {
  text-decoration: none;
  color: var(--ink-black);
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 20px;
  font-weight: 600;
  font-size: 0.9375rem;
  transition: all 0.2s ease;
}

.search-result-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  border: 2px solid rgba(246, 189, 96, 0.3);
  box-shadow: 0 2px 8px rgba(13, 19, 33, 0.08);
  transition: all 0.2s ease;
}

.search-item:hover .search-result-avatar {
  border-color: var(--honey-bronze);
  box-shadow: 0 4px 12px rgba(246, 189, 96, 0.25);
  transform: scale(1.05);
}

.search-result-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.search-result-info {
  flex: 1;
  min-width: 0;
}

.search-result-name {
  font-weight: 600;
  font-size: 0.9375rem;
  color: var(--ink-black);
  letter-spacing: 0.01em;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.search-icon {
  font-size: 1.25rem;
  flex-shrink: 0;
}

.search-empty {
  padding: 32px 24px;
  text-align: center;
  color: rgba(13, 19, 33, 0.5);
  font-size: 0.9375rem;
  font-weight: 500;
}

/* Logout button */
.btn-logout {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border-radius: 12px;
  border: none;
  cursor: pointer;
  font-weight: 700;
  font-size: 0.9375rem;
  background: var(--ink-black);
  color: white;
  transition: all 0.2s ease;
  letter-spacing: 0.01em;
  white-space: nowrap;
}

@media (max-width: 1200px) {
  .btn-logout {
    padding: 10px 16px;
    font-size: 0.875rem;
    gap: 6px;
  }
}

@media (max-width: 900px) {
  .btn-logout span {
    display: none;
  }
  .btn-logout {
    padding: 10px;
    border-radius: 10px;
  }
}

.btn-logout:hover {
  background: #1a2332;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(13, 19, 33, 0.25);
}

.btn-logout:active {
  transform: translateY(0);
}

.btn-logout svg {
  flex-shrink: 0;
}

/* Chat Toggle Button */
.btn-chat-toggle {
  display: none;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  padding: 0;
  border-radius: 12px;
  border: none;
  cursor: pointer;
  background: var(--muted-teal);
  color: white;
  transition: all 0.2s ease;
}

@media (max-width: 1200px) {
  .btn-chat-toggle {
    display: flex;
  }
}

.btn-chat-toggle:hover {
  background: #7da89c;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(146, 191, 177, 0.4);
}

.btn-chat-toggle:active {
  transform: translateY(0);
}

/* Close Chat Button */
.btn-close-chat {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: rgba(13, 19, 33, 0.1);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  z-index: 10;
}

.btn-close-chat:hover {
  background: rgba(13, 19, 33, 0.15);
  transform: scale(1.05);
}

@media (min-width: 1201px) {
  .btn-close-chat {
    display: none;
  }
}

/* Chat Overlay */
.chat-overlay {
  display: none;
  position: fixed;
  top: 74px;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(13, 19, 33, 0.5);
  z-index: 998;
  backdrop-filter: blur(4px);
}

@media (max-width: 1200px) {
  .chat-overlay {
    display: block;
  }
}

/* Fade transition */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

/* =========================== PAGE CONTAINER =========================== */
.page-container {
  margin-top: 74px;
  display: grid;
  grid-template-columns: minmax(280px, 320px) minmax(400px, 1fr) minmax(320px, 360px);
  gap: 32px;
  width: 100%;
  max-width: 1600px;
  margin-left: auto;
  margin-right: auto;
  padding: 0 32px 48px;
  min-height: calc(100vh - 74px);
}

@media (max-width: 1400px) {
  .page-container {
    grid-template-columns: minmax(260px, 280px) minmax(350px, 1fr) minmax(300px, 340px);
    gap: 24px;
    padding: 0 24px 48px;
  }
}

/* =========================== LEFT SIDEBAR =========================== */
.left-sidebar {
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  padding: 24px 16px;
  height: fit-content;
  position: sticky;
  top: 90px;
}

.sidebar-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.sidebar-link {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 16px;
  border-radius: 12px;
  text-decoration: none;
  color: var(--ink-black);
  font-weight: 700;
  font-size: 1.0625rem;
  transition: all 0.2s ease;
  letter-spacing: 0.01em;
}

.sidebar-link:hover {
  background: var(--lavender-mist);
  transform: translateX(4px);
}

.sidebar-link.router-link-active {
  background: var(--honey-bronze);
  color: var(--ink-black);
  box-shadow: 0 2px 12px rgba(246, 189, 96, 0.3);
}

.sidebar-icon {
  font-size: 1.5rem;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.sidebar-text {
  flex: 1;
}

/* =========================== MAIN CONTENT =========================== */
.main-content {
  background: transparent;
  min-height: calc(100vh - 122px);
}

/* =========================== RIGHT SIDEBAR =========================== */
.right-sidebar {
  height: calc(100vh - 122px);
  position: sticky;
  top: 90px;
  transition: all 0.3s ease;
}

@media (max-width: 1200px) {
  .right-sidebar {
    position: fixed;
    right: 20px !important;
    top: 106px !important;
    width: 360px;
    height: calc(100vh - 138px) !important;
    background: white;
    box-shadow: -4px 0 24px rgba(13, 19, 33, 0.15);
    border-radius: 16px;
    transform: translateX(calc(100% + 20px));
    z-index: 999;
  }

  .right-sidebar.chat-open {
    transform: translateX(0);
  }
}

/* =========================== RESPONSIVE =========================== */
@media (max-width: 1200px) {
  .page-container {
    grid-template-columns: minmax(240px, 260px) minmax(300px, 1fr);
    gap: 20px;
  }
}

@media (max-width: 1024px) {
  .page-container {
    grid-template-columns: minmax(70px, 80px) minmax(280px, 1fr);
    gap: 16px;
    padding: 0 16px 48px;
  }

  .left-sidebar {
    padding: 16px 8px;
    min-width: 70px;
  }

  .sidebar-text {
    display: none;
  }

  .sidebar-link {
    justify-content: center;
    padding: 12px;
  }

  .sidebar-icon {
    width: 40px;
    height: 40px;
  }

  .navbar-center {
    margin: 0 8px;
  }

  .user-name {
    display: none;
  }

  .user-info {
    padding: 6px;
  }

  /* Adjust brutalist logo for tablets */
  .brutalist-logo {
    width: 100px;
    height: 50px;
  }

  .logo-square {
    width: 28px;
    height: 28px;
    font-size: 18px;
  }

  .logo-ampersand {
    font-size: 14px;
  }
}

@media (max-width: 768px) {
  .navbar {
    padding: 12px 16px;
  }

  .navbar-content {
    grid-template-columns: auto 1fr auto;
    height: 56px;
    gap: 12px;
  }

  .navbar-center {
    display: none;
  }

  .right-sidebar {
    width: 100%;
    max-width: 360px;
  }

  /* Smaller brutalist logo for mobile */
  .brutalist-logo {
    width: 90px;
    height: 44px;
    border: 2px solid var(--brutalist-white);
    outline: 2px solid var(--brutalist-black);
    box-shadow: 4px 4px 0 var(--brutalist-blue);
  }

  .brutalist-logo:hover {
    box-shadow: 6px 6px 0 var(--brutalist-black);
  }

  .logo-square {
    width: 24px;
    height: 24px;
    font-size: 16px;
  }

  .logo-ampersand {
    font-size: 12px;
  }

  .page-container {
    grid-template-columns: 1fr;
    margin-top: 68px;
    gap: 0;
    padding: 0 0 48px;
  }

  .left-sidebar {
    position: static;
    width: 100%;
    border-radius: 0;
    padding: 16px;
    top: 0;
    box-shadow: 0 2px 8px rgba(13, 19, 33, 0.08);
  }

  .sidebar-content {
    flex-direction: row;
    overflow-x: auto;
    gap: 12px;
  }

  .sidebar-link {
    flex-direction: column;
    gap: 8px;
    min-width: 80px;
    text-align: center;
    padding: 12px;
    font-size: 0.875rem;
  }

  .sidebar-text {
    display: block;
  }

  .sidebar-icon {
    width: 40px;
    height: 40px;
    font-size: 1.375rem;
  }

  .main-content {
    padding: 24px 16px;
  }
}

@media (max-width: 480px) {
  .navbar-menu {
    gap: 12px;
  }

  .btn-logout {
    padding: 10px 16px;
    font-size: 0.875rem;
  }

  .btn-logout svg {
    width: 14px;
    height: 14px;
  }

  .brutalist-logo {
    width: 80px;
    height: 40px;
  }

  .logo-square {
    width: 20px;
    height: 20px;
    font-size: 14px;
  }

  .logo-ampersand {
    font-size: 11px;
  }

  .main-content {
    padding: 16px 12px;
  }

  .page-container {
    padding: 0 0 32px;
  }
}
</style>