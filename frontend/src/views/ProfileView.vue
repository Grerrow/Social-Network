<!-- filepath: /home/kali/Desktop/social-network/frontend/src/views/ProfileView.vue -->
<template>
  <div class="profile-view">
    <div v-if="!profile" class="loading">
      <div class="loading-spinner"></div>
      <p>Loading profile...</p>
    </div>

    <div v-else class="profile-container">
      <!-- PROFILE HEADER -->
      <div class="profile-header">
        <div class="profile-avatar">
          <img
            :src="profile.avatar ? getImageUrl(profile.avatar) : '/images/placeholder.png'"
            :alt="profile.username"
          />
        </div>

        <div class="profile-info">
          <h2 v-if="profile.public && profile.username">{{ profile.username }} ({{ profile.first_name }} {{ profile.last_name }})</h2>
          <h2 v-else-if="profile.public && !profile.username">{{ profile.first_name }} {{ profile.last_name }}</h2>
          <h2 v-else-if="!profile.public && profile.username && profile.follow_status !== 'following'">{{ profile.username }}</h2>
          <h2 v-else-if="!profile.public && !profile.username && profile.follow_status !== 'following'">{{ profile.first_name }} {{ profile.last_name }}</h2>
          <h2 v-else-if="!profile.public && profile.username && profile.follow_status === 'following'">{{ profile.username }} ({{ profile.first_name }} {{ profile.last_name }})</h2>

          <h2 v-else>{{ profile.first_name }} {{ profile.last_name }}</h2>

          <p class="bio" v-if="profile.about_me">{{ profile.about_me }}</p>

          <div class="info-row">
            <div v-if="profile.email" class="info-item">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="info-icon">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                <polyline points="22,6 12,13 2,6"/>
              </svg>
              <span>{{ profile.email }}</span>
            </div>
            <div v-if="profile.date_of_birth" class="info-item">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="info-icon">
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                <line x1="16" y1="2" x2="16" y2="6"/>
                <line x1="8" y1="2" x2="8" y2="6"/>
                <line x1="3" y1="10" x2="21" y2="10"/>
              </svg>
              <span>{{ formatDateOnly(profile.date_of_birth) }}</span>
            </div>
          </div>

            <div class="stats">
            <div class="stat">
              <span class="stat-value">{{ profile.posts?.length || profile.posts_count || 0 }}</span>
              <span class="stat-label">Posts</span>
            </div>
            <div class="stat">
              <span class="stat-value">{{ profile.followers?.length || profile.followers_count || 0 }}</span>
              <router-link v-if="profile.follow_status === 'following' || profile.owner || profile.public"
              class="stat-label"
              :to="{ path: '/friends', query: { user_id: profile.id, tab: 'followers' } }"
              >
              Followers
              </router-link>
              <span v-else class="stat-label">Followers</span>
            </div>
            <div class="stat">
              <span class="stat-value">{{ profile.following?.length || profile.following_count || 0 }}</span>
              <router-link v-if="profile.follow_status === 'following' || profile.owner || profile.public"
              class="stat-label"
              :to="{ path: '/friends', query: { user_id: profile.id, tab: 'following' } }"
              >
              Following
              </router-link>
              <span v-else class="stat-label">Following</span>
            </div>
            </div>

          <FollowRequestButton
            v-if="!profile.owner"
            :user-id="profile.id"
            :initial-status="profile.follow_status"
          />
        </div>
      </div>

      <!-- POSTS MASONRY -->
      <div class="profile-posts" v-if="profile.posts">
        <h3>
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="section-icon">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
            <circle cx="8.5" cy="8.5" r="1.5"/>
            <polyline points="21 15 16 10 5 21"/>
          </svg>
          Posts
        </h3>

        <div v-if="lastPosts.length === 0" class="empty-posts">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="empty-icon">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="7 10 12 15 17 10"/>
            <line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          <p>No posts yet</p>
          <span class="empty-hint">Share your first post to get started</span>
        </div>

        <div v-else class="posts-masonry">
          <!-- Image-only posts: Large photo cards -->
          <div
            v-for="post in lastPosts"
            :key="post.id"
            class="post-card"
            :class="getPostCardClass(post)"
          >
            <!-- Image-only: Clean photo showcase -->
            <template v-if="post.image && !post.content">
              <div class="post-image-wrapper">
                <img
                  :src="getImageUrl(post.image)"
                  alt="Post Image"
                  loading="lazy"
                />
                <div class="post-overlay">
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2">
                    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                    <circle cx="8.5" cy="8.5" r="1.5"/>
                    <polyline points="21 15 16 10 5 21"/>
                  </svg>
                </div>
              </div>
            </template>

            <!-- Text-only card -->
            <template v-else-if="!post.image && post.content">
              <div class="post-quote">
                <p class="quote-text">{{ post.content }}</p>
              </div>
            </template>

            <!-- Image + Text: Story card -->
            <template v-else>
              <div class="post-story">
                <img
                  :src="getImageUrl(post.image)"
                  alt="Post Image"
                  loading="lazy"
                  class="story-image"
                />
                <div class="story-content">
                  <p class="story-text">{{ post.content }}</p>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import FollowRequestButton from '@/components/FollowRequestButton.vue'
import { useApi } from '@/composables/useApi'

const { getImageUrl, getApiUrl } = useApi()


export default {
  name: 'ProfileView',
  components: { FollowRequestButton },

  setup() {
    const route = useRoute()
    const profile = ref(null)
    const routeUserID = ref(route.query.user_id || null)

    const lastPosts = computed(() => {
      if (!profile.value?.posts) return []
      return [...profile.value.posts].slice(-6).reverse()
    })

    const formatDateOnly = (dateString) => {
      if (!dateString) return ''
      // Extract just the date part (YYYY-MM-DD)
      const datePart = dateString.split('T')[0]
      const date = new Date(datePart + 'T00:00:00')
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }

    const getPostCardClass = (post) => {
      if (post.image && !post.content) return 'card-image-only'
      if (!post.image && post.content) return 'card-text-only'
      return 'card-story'
    }

    const fetchProfile = async () => {
      try {
        const url = routeUserID.value
          ? getApiUrl(`/profile?user_id=${routeUserID.value}`)
          : getApiUrl('/profile')

        const res = await fetch(url, { credentials: 'include' })
        if (!res.ok) {
          profile.value = null
          return
        }

        profile.value = await res.json()
      } catch (err) {
        console.error('Error fetching profile:', err)
        profile.value = null
      }
    }

    onMounted(fetchProfile)

    watch(
      () => route.query.user_id,
      (newID) => {
        routeUserID.value = newID || null
        fetchProfile()
      }
    )

    return {
      profile,
      lastPosts,
      formatDateOnly,
      getPostCardClass,
      getApiUrl,
      getImageUrl
    }
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

.profile-view {
  max-width: 900px;
  margin: 0 auto;
  padding: 32px 16px;
  font-family: 'Inter', 'Segoe UI', system-ui, -apple-system, sans-serif;
}

.profile-container {
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  overflow: hidden;
}

.profile-header {
  display: flex;
  flex-wrap: wrap;
  gap: 32px;
  padding: 32px;
  background: linear-gradient(135deg, var(--ink-black) 0%, #1a2332 100%);
  color: white;
}

.profile-avatar {
  flex-shrink: 0;
  width: 150px;
  height: 150px;
  border-radius: 50%;
  overflow: hidden;
  background: var(--lavender-mist);
  border: 4px solid white;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.profile-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 3rem;
  font-weight: 700;
  color: var(--ink-black);
  user-select: none;
}

.profile-info {
  flex: 1;
  min-width: 280px;
}

.profile-info h2 {
  margin: 0 0 8px 0;
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.bio {
  margin: 0 0 16px 0;
  font-style: italic;
  opacity: 0.9;
  line-height: 1.5;
  max-width: 500px;
  font-size: 0.9375rem;
}

.info-row {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  margin-bottom: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9375rem;
}

.info-icon {
  flex-shrink: 0;
  opacity: 0.9;
}

.stats {
  display: flex;
  gap: 32px;
  margin-bottom: 16px;
}

.stat {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  letter-spacing: -0.02em;
}
.stat-label {
  font-size: 0.8125rem;
  opacity: 0.8;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  text-decoration: none; /* Remove underline */
  color: inherit; /* Inherit white color from parent */
  cursor: pointer;
  transition: all 0.2s ease;
  display: inline-block;
}

/* When it's a router-link */
a.stat-label {
  text-decoration: none;
  color: inherit;
}

/* Subtle hover effect */
a.stat-label:hover {
  opacity: 1;
  transform: translateY(-1px);
}

/* Remove visited link color */
a.stat-label:visited {
  color: inherit;
}

/* Active state */
a.stat-label:active {
  transform: translateY(0);
}

.profile-posts {
  padding: 32px;
}

.profile-posts h3 {
  margin: 0 0 24px 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--ink-black);
  display: flex;
  align-items: center;
  gap: 12px;
  letter-spacing: -0.01em;
}

.section-icon {
  color: var(--muted-teal);
  flex-shrink: 0;
}

/* MASONRY LAYOUT */
.posts-masonry {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 12px;
  align-items: start;
}

/* POST CARDS */
.post-card {
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1.5px solid rgba(13, 19, 33, 0.12);
  background: white;
}

.post-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(13, 19, 33, 0.12);
  border-color: rgba(13, 19, 33, 0.2);
}

/* IMAGE-ONLY CARDS */
.card-image-only {
  position: relative;
  aspect-ratio: 1 / 1;
}

.post-image-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
}

.post-image-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.post-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(180deg, transparent 0%, rgba(13, 19, 33, 0.4) 100%);
  display: flex;
  align-items: flex-end;
  justify-content: flex-end;
  padding: 16px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.post-card:hover .post-overlay {
  opacity: 1;
}

/* TEXT-ONLY CARDS */
.card-text-only {
  background: white;
  min-height: 140px;
  display: flex;
  align-items: center;
}

.post-quote {
  padding: 16px;
  display: flex;
  flex-direction: column;
  width: 100%;
}

.quote-icon {
  display: none;
}

.quote-text {
  margin: 0;
  font-size: 0.9375rem;
  line-height: 1.5;
  color: var(--ink-black);
  font-weight: 400;
  letter-spacing: 0.005em;
  word-wrap: break-word;
  overflow-wrap: break-word;
  display: -webkit-box;
  -webkit-line-clamp: 6;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.quote-text::before,
.quote-text::after {
  content: none;
}

/* IMAGE + TEXT CARDS (STORY STYLE) */
.card-story {
  background: white;
}

.post-story {
  display: flex;
  flex-direction: column;
}

.story-image {
  width: 100%;
  height: 160px;
  object-fit: cover;
  display: block;
}

.story-content {
  padding: 12px;
  background: white;
  border-top: 1px solid rgba(13, 19, 33, 0.08);
}

.story-text {
  margin: 0;
  font-size: 0.875rem;
  line-height: 1.4;
  color: var(--ink-black);
  word-wrap: break-word;
  overflow-wrap: break-word;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* EMPTY STATE */
.empty-posts {
  text-align: center;
  padding: 64px 32px;
  background: var(--lavender-mist);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.empty-icon {
  color: rgba(13, 19, 33, 0.3);
  margin-bottom: 8px;
}

.empty-posts p {
  margin: 0;
  color: rgba(13, 19, 33, 0.6);
  font-size: 1.125rem;
  font-weight: 600;
}

.empty-hint {
  font-size: 0.875rem;
  color: rgba(13, 19, 33, 0.5);
}

/* LOADING */
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 32px;
  color: rgba(13, 19, 33, 0.6);
}

.loading p {
  margin: 0;
  font-size: 0.9375rem;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid rgba(13, 19, 33, 0.1);
  border-top-color: var(--honey-bronze);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* RESPONSIVE */
@media (max-width: 768px) {
  .posts-masonry {
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 10px;
  }
}

@media (max-width: 600px) {
  .profile-view {
    padding: 16px 12px;
  }

  .profile-header {
    padding: 24px;
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .profile-avatar {
    width: 120px;
    height: 120px;
  }

  .avatar-placeholder {
    font-size: 2.5rem;
  }

  .profile-info h2 {
    font-size: 1.5rem;
  }

  .info-row {
    flex-direction: column;
    gap: 12px;
    align-items: center;
  }

  .stats {
    justify-content: center;
    gap: 24px;
  }

  .profile-posts {
    padding: 24px 16px;
  }

  .posts-masonry {
    grid-template-columns: 1fr;
  }

  .empty-posts {
    padding: 48px 24px;
  }

  .quote-text {
    font-size: 1rem;
  }
}
</style>