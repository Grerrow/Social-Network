<!-- filepath: /home/kali/Desktop/social-network/frontend/src/views/NotificationsView.vue -->
<template>
  <div class="notifications-view">
    <div class="header">
      <h1>
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="header-icon">
          <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
          <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
        </svg>
        Notifications
      </h1>

      <button
        v-if="store.unreadCount > 0"
        class="btn-mark-read"
        @click="markAllRead"
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="20 6 9 17 4 12"/>
        </svg>
        Mark all as read
      </button>
    </div>

    <!-- Loading -->
    <div v-if="store.isLoading" class="loading">
      <div class="spinner"></div>
      <span>Loading notifications...</span>
    </div>

    <!-- EMPTY -->
    <div v-else-if="store.notifications.length === 0" class="empty">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="empty-icon">
        <circle cx="12" cy="12" r="10"/>
        <path d="M8 14s1.5 2 4 2 4-2 4-2"/>
        <line x1="9" y1="9" x2="9.01" y2="9"/>
        <line x1="15" y1="9" x2="15.01" y2="9"/>
      </svg>
      <h3>You're all caught up!</h3>
      <p>No new notifications at the moment</p>
    </div>

    <!-- LIST -->
    <div v-else class="notifications-list">
      <TransitionGroup name="notification">
        <div
          v-for="n in store.notifications"
          :key="n.id"
          class="notification-card"
          :class="{ unread: !n.is_read }"
        >
          <!-- ICON -->
          <div class="icon" :class="getIconClass(n.type)">
            <!-- New Comment -->
            <template v-if="n.type === 'new_comment'">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              </svg>
            </template>

            <!-- Follow Request -->
            <template v-else-if="n.type === 'follow_request'">
              <img
                v-if="n.sender_avatar"
                :src="getImageUrl(n.sender_avatar)"
                :alt="n.sender_name"
                class="avatar"
              />
              <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="8.5" cy="7" r="4"/>
                <line x1="20" y1="8" x2="20" y2="14"/>
                <line x1="23" y1="11" x2="17" y2="11"/>
              </svg>
            </template>

            <!-- New Follower -->
            <template v-else-if="n.type === 'new_follower'">
              <img
                v-if="n.sender_avatar"
                :src="getImageUrl(n.sender_avatar)"
                :alt="n.sender_name"
                class="avatar"
              />
              <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              </svg>
            </template>

            <!-- User Accepted Follow -->
            <template v-else-if="n.type === 'user_accepted_follow'">
              <img
                v-if="n.sender_avatar"
                :src="getImageUrl(n.sender_avatar)"
                :alt="n.sender_name"
                class="avatar"
              />
              <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="8.5" cy="7" r="4"/>
                <polyline points="17 11 19 13 23 9"/>
              </svg>
            </template>

            <!-- Message -->
            <template v-else-if="n.type === 'message'">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                <polyline points="22,6 12,13 2,6"/>
              </svg>
            </template>

            <!-- Group Invite / Invitation -->
            <template v-else-if="n.type === 'group_invite' || n.type === 'group_invitation'">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              </svg>
            </template>

            <!-- Group Join Request -->
            <template v-else-if="n.type === 'group_join_request'">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <line x1="19" y1="8" x2="19" y2="14"/>
                <line x1="22" y1="11" x2="16" y2="11"/>
              </svg>
            </template>

            <!-- Group Join Request Approved -->
            <template v-else-if="n.type === 'group_join_request_approved'">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <polyline points="16 11 18 13 22 9"/>
              </svg>
            </template>

            <!-- Group Event -->
            <template v-else-if="n.type === 'group_event'">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                <line x1="16" y1="2" x2="16" y2="6"/>
                <line x1="8" y1="2" x2="8" y2="6"/>
                <line x1="3" y1="10" x2="21" y2="10"/>
              </svg>
            </template>

            <!-- Event Created -->
            <template v-else-if="n.type === 'event_created'">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                <line x1="16" y1="2" x2="16" y2="6"/>
                <line x1="8" y1="2" x2="8" y2="6"/>
                <line x1="3" y1="10" x2="21" y2="10"/>
                <path d="M9 16l2 2 4-4"/>
              </svg>
            </template>

            <!-- Default Notification Bell -->
            <template v-else>
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
                <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
              </svg>
            </template>
          </div>

          <!-- CONTENT -->
          <div class="content">
            <p class="message">
              <!-- Group Join Request Approved - Special format -->
              <template v-if="n.type === 'group_join_request_approved'">
                <span class="action-text">Your request to join </span>
                <span class="group-name">{{ n.group_name || 'the group' }}</span>
                <span class="action-text"> was approved!</span>
              </template>

              <!-- Event Created - Special format -->
              <template v-else-if="n.type === 'event_created'">
                <span class="sender-name" v-if="n.sender_name" @click="viewProfile(n.sender_id)">
                  {{ n.sender_name }}
                </span>
                <span class="action-text"> created a new event</span>
                <span v-if="n.event_title" class="event-title"> "{{ n.event_title }}"</span>
                <span class="action-text"> in </span>
                <span class="group-name">{{ n.group_name || 'a group' }}</span>
              </template>

              <!-- Standard format -->
              <template v-else>
                <span class="sender-name" v-if="n.sender_name" @click="viewProfile(n.sender_id)">
                  {{ n.sender_name }}
                </span>
                <span class="action-text">{{ getActionText(n) }}</span>
                <span v-if="n.group_name && n.type !== 'group_join_request_approved'" class="group-name">
                  {{ n.group_name }}
                </span>
              </template>
            </p>

            <div class="meta">
              <span class="time">{{ formatTime(n.created_at) }}</span>
              <span v-if="!n.is_read" class="unread-dot"></span>
            </div>

            <!-- ACTIONS -->

            <!-- Follow Request -->
            <div v-if="n.type === 'follow_request'" class="actions">
              <button class="btn-primary" @click="acceptFollow(n)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
                Accept
              </button>
              <button class="btn-secondary" @click="rejectFollow(n)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"/>
                  <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
                Decline
              </button>
            </div>

            <!-- New Follower / Accepted Follow -->
            <div v-else-if="n.type === 'new_follower' || n.type === 'user_accepted_follow'" class="actions">
              <button class="btn-outline" @click="viewProfile(n.sender_id)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
                View Profile
              </button>
            </div>

            <!-- New Comment -->
            <div v-else-if="n.type === 'new_comment'" class="actions">
              <button class="btn-outline" @click="viewPost(n.post_id)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
                View Post
              </button>
            </div>

            <!-- Group Invite -->
            <div v-else-if="n.type === 'group_invite' || n.type === 'group_invitation'" class="actions">
              <button class="btn-primary" @click="acceptGroupInvite(n)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
                Join Group
              </button>
              <button class="btn-secondary" @click="declineGroupInvite(n)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"/>
                  <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
                Decline
              </button>
            </div>

            <!-- Group Join Request -->
            <div v-else-if="n.type === 'group_join_request'" class="actions">
              <button class="btn-primary" @click="approveJoinRequest(n)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
                Approve
              </button>
              <button class="btn-secondary" @click="rejectJoinRequest(n)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"/>
                  <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
                Reject
              </button>
            </div>

            <!-- Group Join Request Approved -->
            <div v-else-if="n.type === 'group_join_request_approved'" class="actions">
              <button class="btn-primary" @click="viewGroup(n.group_id)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                  <circle cx="9" cy="7" r="4"/>
                  <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                </svg>
                View Group
              </button>
            </div>

          

            <!-- Event Created -->
            <div v-else-if="n.type === 'event_created'" class="actions">
              <button class="btn-primary" @click="viewGroup(n.group_id)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                  <line x1="16" y1="2" x2="16" y2="6"/>
                  <line x1="8" y1="2" x2="8" y2="6"/>
                  <line x1="3" y1="10" x2="21" y2="10"/>
                </svg>
                View Event
              </button>
            </div>
          </div>

          <!-- Dismiss button -->
          <button class="dismiss-btn" @click="dismissNotification(n)" title="Dismiss">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
      </TransitionGroup>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationsStore } from '@/stores/notifications'
import { useFollowStatus } from '@/composables/useFollowStatus'
import { useGroupsStore } from '@/stores/groups'
import { useApi } from '@/composables/useApi'

const { getImageUrl } = useApi()

const store = useNotificationsStore()
const groupsStore = useGroupsStore()
const router = useRouter()

const { acceptFollowRequest, rejectFollowRequest } = useFollowStatus()

onMounted(() => {
  store.fetchNotifications()
})

const getIconClass = (type) => {
  const classes = {
    'new_comment': 'icon-comment',
    'follow_request': 'icon-follow',
    'new_follower': 'icon-follower',
    'user_accepted_follow': 'icon-accepted',
    'message': 'icon-message',
    'group_invite': 'icon-group',
    'group_invitation': 'icon-group',
    'group_join_request': 'icon-group-request',
    'group_join_request_approved': 'icon-group-approved',
    'group_event': 'icon-event',
    'event_created': 'icon-event-created'
  }
  return classes[type] || 'icon-default'
}

const getActionText = (n) => {
  const texts = {
    'new_comment': ' commented on your post',
    'follow_request': ' wants to follow you',
    'new_follower': ' started following you',
    'user_accepted_follow': ' accepted your follow request',
    'message': 'You received a new message',
    'group_invite': ' invited you to join ',
    'group_invitation': ' invited you to join ',
    'group_join_request': ' requested to join ',
    'group_join_request_approved': '',
    'group_event': ' created a new event in ',
    'event_created': ' created a new event in '
  }
  return texts[n.type] || 'You have a new notification'
}

const acceptFollow = async (n) => {
  await acceptFollowRequest(n.follow_request_id)
  store.removeNotification(n.id)
}

const rejectFollow = async (n) => {
  await rejectFollowRequest(n.follow_request_id)
  store.removeNotification(n.id)
}

const acceptGroupInvite = async (n) => {
  const success = await groupsStore.acceptInvite(n.group_id)
  if (success) {
    store.removeNotification(n.id)
    router.push(`/groups/${n.group_id}`)
  }
}

const declineGroupInvite = async (n) => {
  const success = await groupsStore.declineInvite(n.group_id)
  if (success) {
    store.removeNotification(n.id)
  }
}

const approveJoinRequest = async (n) => {
  const success = await groupsStore.approveJoinRequest(n.group_id, n.sender_id)
  if (success) {
    store.removeNotification(n.id)
  }
}

const rejectJoinRequest = async (n) => {
  const success = await groupsStore.rejectJoinRequest(n.group_id, n.sender_id)
  if (success) {
    store.removeNotification(n.id)
  }
}

const dismissNotification = (n) => {
  store.removeNotification(n.id)
}

const markAllRead = () => store.markAllAsRead()

const viewProfile = (userId) => {
  router.push({ path: '/profile', query: { user_id: userId } })
}

const viewGroup = (groupId) => {
  router.push(`/groups/${groupId}`)
}

const viewPost = (postId) => {
  if (postId) {
    router.push(`/post/${postId}`)
  }
}

const formatTime = (date) => {
  if (!date) return ''
  
  let d
  if (typeof date === 'number') {
    d = new Date(date * 1000)
  } else if (typeof date === 'string') {
    d = new Date(date)
  } else {
    return ''
  }
  
  if (isNaN(d.getTime())) return ''
  
  const now = new Date()
  const diff = Math.floor((now - d) / 1000)
  
  if (diff < 60) return 'Just now'
  if (diff < 3600) {
    const mins = Math.floor(diff / 60)
    return `${mins} ${mins === 1 ? 'minute' : 'minutes'} ago`
  }
  if (diff < 86400) {
    const hours = Math.floor(diff / 3600)
    return `${hours} ${hours === 1 ? 'hour' : 'hours'} ago`
  }
  if (diff < 604800) {
    const days = Math.floor(diff / 86400)
    return `${days} ${days === 1 ? 'day' : 'days'} ago`
  }
  
  return d.toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric',
    year: d.getFullYear() !== now.getFullYear() ? 'numeric' : undefined
  })
}
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.notifications-view {
  max-width: 680px;
  margin: 0 auto;
  padding: 24px 16px;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header h1 {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0;
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.02em;
}

.header-icon {
  color: var(--honey-bronze);
}

.btn-mark-read {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  background: var(--muted-teal);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  letter-spacing: 0.01em;
}

.btn-mark-read:hover {
  background: #7da99c;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(146, 191, 177, 0.3);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 48px;
  color: rgba(13, 19, 33, 0.6);
}

.loading span {
  font-size: 0.9375rem;
  font-weight: 500;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(13, 19, 33, 0.1);
  border-top-color: var(--honey-bronze);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty {
  text-align: center;
  padding: 64px 32px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.08);
}

.empty-icon {
  color: var(--honey-bronze);
  margin-bottom: 20px;
}

.empty h3 {
  margin: 0 0 8px;
  font-size: 1.25rem;
  color: var(--ink-black);
  font-weight: 700;
  letter-spacing: -0.01em;
}

.empty p {
  margin: 0;
  color: rgba(13, 19, 33, 0.6);
  font-size: 0.9375rem;
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-card {
  display: flex;
  gap: 14px;
  padding: 16px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.08);
  transition: all 0.2s ease;
  position: relative;
  border: 2px solid transparent;
}

.notification-card:hover {
  box-shadow: 0 4px 16px rgba(13, 19, 33, 0.12);
}

.notification-card.unread {
  background: linear-gradient(to right, rgba(246, 240, 249, 0.6), white);
  border-left: 3px solid var(--honey-bronze);
}

.icon {
  flex-shrink: 0;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--lavender-mist);
}

.icon svg {
  color: rgba(13, 19, 33, 0.7);
}

.icon .avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.icon-comment { background: rgba(146, 191, 177, 0.15); }
.icon-comment svg { color: var(--muted-teal); }

.icon-follow, .icon-follower { background: rgba(146, 191, 177, 0.15); }
.icon-follow svg, .icon-follower svg { color: var(--muted-teal); }

.icon-accepted { background: rgba(146, 191, 177, 0.2); }
.icon-accepted svg { color: #31a24c; }

.icon-message { background: rgba(246, 189, 96, 0.15); }
.icon-message svg { color: var(--honey-bronze); }

.icon-group { background: rgba(246, 240, 249, 0.8); }
.icon-group svg { color: var(--ink-black); }

.icon-group-request { background: rgba(246, 189, 96, 0.15); }
.icon-group-request svg { color: var(--honey-bronze); }

.icon-group-approved { background: rgba(146, 191, 177, 0.2); }
.icon-group-approved svg { color: #31a24c; }

.icon-event { background: rgba(246, 189, 96, 0.15); }
.icon-event svg { color: var(--honey-bronze); }

.icon-event-created { background: rgba(146, 191, 177, 0.2); }
.icon-event-created svg { color: var(--muted-teal); }

.icon-default { background: var(--lavender-mist); }
.icon-default svg { color: rgba(13, 19, 33, 0.6); }

.content {
  flex: 1;
  min-width: 0;
}

.message {
  margin: 0 0 6px;
  font-size: 0.9375rem;
  line-height: 1.5;
  color: var(--ink-black);
}

.sender-name {
  font-weight: 700;
  color: var(--ink-black);
  cursor: pointer;
  transition: color 0.2s;
}

.sender-name:hover {
  color: var(--muted-teal);
  text-decoration: underline;
}

.action-text {
  color: rgba(13, 19, 33, 0.7);
}

.group-name {
  font-weight: 700;
  color: var(--ink-black);
}

.event-title {
  font-weight: 700;
  color: var(--muted-teal);
}

.meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.time {
  font-size: 0.8125rem;
  color: rgba(13, 19, 33, 0.5);
  font-weight: 500;
}

.unread-dot {
  width: 8px;
  height: 8px;
  background: var(--honey-bronze);
  border-radius: 50%;
}

.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 9px 16px;
  background: var(--muted-teal);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  letter-spacing: 0.01em;
}

.btn-primary:hover {
  background: #7da99c;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(146, 191, 177, 0.3);
}

.btn-secondary {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 9px 16px;
  background: rgba(13, 19, 33, 0.08);
  color: var(--ink-black);
  border: none;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  letter-spacing: 0.01em;
}

.btn-secondary:hover {
  background: rgba(13, 19, 33, 0.15);
}

.btn-outline {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 9px 16px;
  background: transparent;
  color: var(--muted-teal);
  border: 2px solid var(--muted-teal);
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  letter-spacing: 0.01em;
}

.btn-outline:hover {
  background: rgba(146, 191, 177, 0.1);
}

.dismiss-btn {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  opacity: 0;
  transition: all 0.2s;
  color: rgba(13, 19, 33, 0.6);
}

.notification-card:hover .dismiss-btn {
  opacity: 1;
}

.dismiss-btn:hover {
  background: rgba(13, 19, 33, 0.08);
  color: var(--ink-black);
}

.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.notification-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

@media (max-width: 600px) {
  .notifications-view {
    padding: 16px 12px;
  }

  .header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }

  .btn-mark-read {
    width: 100%;
    justify-content: center;
  }

  .notification-card {
    padding: 14px;
  }

  .icon {
    width: 40px;
    height: 40px;
  }

  .icon .avatar {
    width: 40px;
    height: 40px;
  }

  .icon svg {
    width: 20px;
    height: 20px;
  }

  .actions {
    flex-direction: column;
  }

  .actions button {
    width: 100%;
    justify-content: center;
  }

  .dismiss-btn {
    opacity: 1;
  }
}
</style>