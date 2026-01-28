<!-- filepath: /home/kali/Desktop/social-network/frontend/src/views/GroupDetails.vue -->
<template>
  <div class="group-details-page" v-if="group">
    <!-- Group Header -->
    <div class="group-header">
      <div class="group-cover">
        <div class="group-avatar">
          {{ group.title?.charAt(0).toUpperCase() }}
        </div>
      </div>
      <div class="group-info">
        <h1>{{ group.title }}</h1>
        <p class="group-description">{{ group.description }}</p>
        <div class="group-meta">
          <span>{{ group.member_count || 0 }} members</span>
          <span>•</span>
          <span>Created by {{ group.creator?.username }}</span>
        </div>
        
        <!-- Action Buttons based on membership status -->
        <div class="group-actions">
          <template v-if="!group.is_member && !group.is_invited && !group.has_requested">
            <button class="btn-primary" @click="handleRequestJoin">
              Request to Join
            </button>
          </template>
          <template v-else-if="group.has_requested">
            <button class="btn-secondary" disabled>
              Request Pending
            </button>
          </template>
          <template v-else-if="group.is_invited">
            <button class="btn-primary" @click="handleAcceptInvite">
              Accept Invitation
            </button>
            <button class="btn-secondary" @click="handleDeclineInvite">
              Decline
            </button>
          </template>
          <template v-else-if="group.is_member">
            <button class="btn-secondary" @click="showInviteModal = true">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="8.5" cy="7" r="4"/>
                <line x1="20" y1="8" x2="20" y2="14"/>
                <line x1="23" y1="11" x2="17" y2="11"/>
              </svg>
              Invite
            </button>
            <!-- if user not creator -->
            <button v-if="!group.is_creator" class="btn-danger" @click="handleLeaveGroup">
              Leave Group
            </button>
          </template>
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="group-tabs" v-if="group.is_member">
      <button 
        :class="['tab', { active: activeTab === 'posts' }]"
        @click="activeTab = 'posts'"
      >
        Posts
      </button>
      <button 
        :class="['tab', { active: activeTab === 'events' }]"
        @click="activeTab = 'events'"
      >
        Events
      </button>
      <button 
        :class="['tab', { active: activeTab === 'members' }]"
        @click="activeTab = 'members'"
      >
        Members
      </button>
      <button 
        v-if="group.is_creator"
        :class="['tab', { active: activeTab === 'requests' }]"
        @click="activeTab = 'requests'"
      >
        Requests
        <span v-if="group.pending_requests?.length" class="badge">
          {{ group.pending_requests.length }}
        </span>
      </button>
    </div>

    <!-- Tab Content -->
    <div class="tab-content" v-if="group.is_member">
      <!-- Posts Tab -->
      <div v-if="activeTab === 'posts'" class="posts-section">
        <!-- Create Post -->
        <div class="create-post-card">
          <div class="create-post-input" @click="showPostModal = true">
            <img 
              class="avatar-placeholder"
              :src="authStore.user?.avatar ? getImageUrl(authStore.user.avatar) : '/images/placeholder.png'" 
              :alt="authStore.user?.username" 
            />
            <span>Write something to the group...</span>
          </div>
        </div>

        <!-- Posts List -->
        <div v-if="groupsStore.groupPosts.length === 0" class="empty-state">
          <p>No posts yet. Be the first to post!</p>
        </div>
        <GroupPost 
          v-for="post in groupsStore.groupPosts" 
          :key="post.id" 
          :post="post"
          :group-id="group.id"
        />
      </div>

      <!-- Events Tab -->
      <div v-if="activeTab === 'events'" class="events-section">
        <button class="btn-primary create-event-btn" @click="showEventModal = true">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          Create Event
        </button>

        <div v-if="groupsStore.groupEvents.length === 0" class="empty-state">
          <p>No events scheduled</p>
        </div>
        <GroupEvent 
          v-for="event in groupsStore.groupEvents" 
          :key="event.id" 
          :event="event"
          :group-id="group.id"
        />
      </div>

      <!-- Members Tab -->
      <div v-if="activeTab === 'members'" class="members-section">
        <div class="members-list">
          <div v-for="member in group.members" :key="member.id" class="member-item">
            <div class="member-avatar">
              <img :src="member.avatar ? getImageUrl(member.avatar) : '/images/placeholder.png'" :alt="member.username" />
            </div>
            <div class="member-info">
              <span class="member-name">{{ member.username }}</span>
              <span v-if="member.id === group.creator_id" class="member-role">Creator</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Requests Tab (Creator Only) -->
      <div v-if="activeTab === 'requests' && group.is_creator" class="requests-section">
        <div v-if="!group.pending_requests?.length" class="empty-state">
          <p>No pending requests</p>
        </div>
        <div v-for="request in group.pending_requests" :key="request.id" class="request-item">
          <div class="member-avatar">
            <img :src="request.avatar ? getImageUrl(request.avatar) : '/images/placeholder.png'" :alt="request.username" />
          </div>
          <div class="request-info">
            <span class="member-name">{{ request.username }}</span>
            <span class="request-text">wants to join</span>
          </div>
          <div class="request-actions">
            <button class="btn-primary btn-sm" @click="handleApproveRequest(request.user_id)">
              Approve
            </button>
            <button class="btn-secondary btn-sm" @click="handleRejectRequest(request.user_id)">
              Reject
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Non-member view -->
    <div v-else class="non-member-view">
      <div class="locked-content">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
          <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
        </svg>
        <p>Join this group to see posts and events</p>
      </div>
    </div>

    <!-- Create Post Modal -->
    <div v-if="showPostModal" class="modal-overlay" @click.self="showPostModal = false">
      <div class="modal">
        <div class="modal-header">
          <h2>Create Post</h2>
          <button class="close-btn" @click="showPostModal = false">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
            </svg>
          </button>
        </div>
        <form @submit.prevent="handleCreatePost" class="modal-form">
          <div class="form-group">
            <textarea 
              v-model="newPost.content" 
              placeholder="What's on your mind?"
              rows="4"
            ></textarea>
          </div>
          
          <!-- Image Preview -->
          <div v-if="newPost.imagePreview" class="image-preview">
            <img :src="newPost.imagePreview" alt="Preview" />
            <button type="button" class="remove-image-btn" @click="removeImage">×</button>
          </div>
          
          <div class="form-group">
            <label class="file-label">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                <circle cx="8.5" cy="8.5" r="1.5"/>
                <polyline points="21 15 16 10 5 21"/>
              </svg>
              {{ newPost.imageFile ? 'Change Photo' : 'Add Photo' }}
              <input type="file" accept="image/jpeg,image/png,image/gif" @change="handleImageSelect" hidden />
            </label>
            <span v-if="newPost.imageFile" class="file-name">{{ newPost.imageFile.name }}</span>
          </div>
          <div class="modal-actions">
            <button type="submit" class="btn-primary" :disabled="!newPost.content.trim() && !newPost.imageFile">
              Post
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Create Event Modal - PROFESSIONAL DATE & TIME PICKER -->
    <div v-if="showEventModal" class="modal-overlay" @click.self="showEventModal = false">
      <div class="modal">
        <div class="modal-header">
          <h2>Create Event</h2>
          <button class="close-btn" @click="showEventModal = false">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
            </svg>
          </button>
        </div>
        <form @submit.prevent="handleCreateEvent" class="modal-form">
          <div class="form-group">
            <label>Event Title</label>
            <input 
              v-model="newEvent.title" 
              type="text" 
              placeholder="Event name"
              required
            />
          </div>
          <div class="form-group">
            <label>Description</label>
            <textarea 
              v-model="newEvent.description" 
              placeholder="What's the event about?"
              rows="3"
              required
            ></textarea>
          </div>
          
          <!-- PROFESSIONAL DATE PICKER -->
          <div class="form-group">
            <label>Event Date</label>
            <div class="date-picker-row">
              <select v-model="eventDateForm.month" required class="date-select">
                <option value="" disabled>Month</option>
                <option v-for="(month, index) in months" :key="index" :value="index + 1">
                  {{ month }}
                </option>
              </select>
              
              <select v-model="eventDateForm.day" required class="date-select">
                <option value="" disabled>Day</option>
                <option v-for="day in daysInEventMonth" :key="day" :value="day">
                  {{ day }}
                </option>
              </select>
              
              <select v-model="eventDateForm.year" required class="date-select">
                <option value="" disabled>Year</option>
                <option v-for="year in eventYears" :key="year" :value="year">
                  {{ year }}
                </option>
              </select>
            </div>
          </div>

          <!-- PROFESSIONAL TIME PICKER -->
          <div class="form-group">
            <label>Event Time</label>
            <div class="time-picker-row">
              <select v-model="eventDateForm.hour" required class="time-select">
                <option value="" disabled>Hour</option>
                <option v-for="hour in 12" :key="hour" :value="hour">
                  {{ hour }}
                </option>
              </select>
              
              <select v-model="eventDateForm.minute" required class="time-select">
                <option value="" disabled>Min</option>
                <option value="00">00</option>
                <option value="15">15</option>
                <option value="30">30</option>
                <option value="45">45</option>
              </select>
              
              <select v-model="eventDateForm.period" required class="time-select">
                <option value="" disabled>AM/PM</option>
                <option value="AM">AM</option>
                <option value="PM">PM</option>
              </select>
            </div>
          </div>

          <div class="modal-actions">
            <button type="button" class="btn-secondary" @click="showEventModal = false">
              Cancel
            </button>
            <button type="submit" class="btn-primary">
              Create Event
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Invite Modal -->
    <div v-if="showInviteModal" class="modal-overlay" @click.self="showInviteModal = false">
      <div class="modal">
        <div class="modal-header">
          <h2>Invite to Group</h2>
          <button class="close-btn" @click="showInviteModal = false">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
            </svg>
          </button>
        </div>
        <div class="modal-form">
          <div class="form-group">
            <input 
              v-model="inviteSearch" 
              type="text" 
              placeholder="Search users..."
              class="search-input"
            />
          </div>
          <div class="invite-list">
            <div 
              v-for="user in filteredInviteUsers" 
              :key="user.id" 
              class="invite-item"
            >
              <div class="member-avatar">
                <img  :src="user.avatar ? getImageUrl(user.avatar) : '/images/placeholder.png'" :alt="user.username" />
              </div>
              <span class="member-name">{{ user.username }}</span>
              <button 
                class="btn-primary btn-sm" 
                @click="handleInviteUser(user.id)"
                :disabled="user.invited"
              >
                {{ user.invited ? 'Invited' : 'Invite' }}
              </button>
            </div>
            <div v-if="filteredInviteUsers.length === 0" class="empty-state">
              <p>No users to invite</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div v-else-if="isLoading" class="loading-page">
    <div class="spinner"></div>
    <p>Loading group...</p>
  </div>

  <div v-else class="error-page">
    <p>Group not found</p>
    <router-link to="/groups" class="btn-primary">Back to Groups</router-link>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useGroupsStore } from '@/stores/groups'
import { useAuthStore } from '@/stores/auth'
import GroupPost from '@/components/groups/GroupPost.vue'
import GroupEvent from '@/components/groups/GroupEvent.vue'
import api from '@/services/api'
import { useApi } from '@/composables/useApi'

const { getImageUrl } = useApi()


const route = useRoute()
const groupsStore = useGroupsStore()
const authStore = useAuthStore()

const activeTab = ref('posts')
const showPostModal = ref(false)
const showEventModal = ref(false)
const showInviteModal = ref(false)
const inviteSearch = ref('')
const inviteUsers = ref([])

const newPost = ref({
  content: '',
  imageFile: null,
  imagePreview: null
})

const newEvent = ref({
  title: '',
  description: '',
  event_date: ''
})

// Professional Date & Time Picker State
const eventDateForm = ref({
  month: '',
  day: '',
  year: '',
  hour: '',
  minute: '',
  period: ''
})

const months = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December'
]

// Generate years (current year to 5 years ahead)
const currentYear = new Date().getFullYear()
const eventYears = computed(() => {
  const arr = []
  for (let i = currentYear; i <= currentYear + 5; i++) {
    arr.push(i)
  }
  return arr
})

// Calculate days in selected month/year
const daysInEventMonth = computed(() => {
  if (!eventDateForm.value.month || !eventDateForm.value.year) return 31
  const month = parseInt(eventDateForm.value.month)
  const year = parseInt(eventDateForm.value.year)
  return new Date(year, month, 0).getDate()
})

const group = computed(() => groupsStore.currentGroup)
const isLoading = computed(() => groupsStore.isLoading)


const filteredInviteUsers = computed(() => {
  if (!inviteSearch.value.trim()) return inviteUsers.value
  const query = inviteSearch.value.toLowerCase()
  return inviteUsers.value.filter(u => 
    u.username.toLowerCase().includes(query)
  )
})

const loadGroupData = async () => {
  const groupId = route.params.id
  await groupsStore.fetchGroup(groupId)
  
  if (group.value?.is_member) {
    await groupsStore.fetchGroupPosts(groupId)
    await groupsStore.fetchGroupEvents(groupId)
  }
}

const handleLeaveGroup = async () => {
  const confirmLeave = confirm('Are you sure you want to leave this group?')
  if (confirmLeave) {
    await groupsStore.leaveGroup(route.params.id)
    await loadGroupData()
  }
}

const handleRequestJoin = async () => {
  await groupsStore.requestJoin(route.params.id)
}

const handleAcceptInvite = async () => {
  await groupsStore.acceptInvite(route.params.id)
  await loadGroupData()
}

const handleDeclineInvite = async () => {
  await groupsStore.declineInvite(route.params.id)
}

const handleApproveRequest = async (userId) => {
  await groupsStore.approveJoinRequest(route.params.id, userId)
}

const handleRejectRequest = async (userId) => {
  await groupsStore.rejectJoinRequest(route.params.id, userId)
}

const handleImageSelect = (e) => {
  const file = e.target.files?.[0]
  if (file) {
    if (!['image/jpeg', 'image/png', 'image/gif'].includes(file.type)) {
      alert('Only JPEG, PNG, and GIF images are allowed')
      return
    }
    
    if (file.size > 5 * 1024 * 1024) {
      alert('Image must be less than 5MB')
      return
    }
    
    newPost.value.imageFile = file
    newPost.value.imagePreview = URL.createObjectURL(file)
  }
}

const removeImage = () => {
  if (newPost.value.imagePreview) {
    URL.revokeObjectURL(newPost.value.imagePreview)
  }
  newPost.value.imageFile = null
  newPost.value.imagePreview = null
}

const handleCreatePost = async () => {
  if (!newPost.value.content.trim() && !newPost.value.imageFile) return
  
  await groupsStore.createGroupPost(
    route.params.id,
    newPost.value.content.trim() || '',
    newPost.value.imageFile
  )
  
  removeImage()
  newPost.value.content = ''
  showPostModal.value = false
}


const handleCreateEvent = async () => {
  if (!newEvent.value.title.trim() || 
      !eventDateForm.value.month || 
      !eventDateForm.value.day || 
      !eventDateForm.value.year ||
      !eventDateForm.value.hour ||
      !eventDateForm.value.minute ||
      !eventDateForm.value.period) {
    alert('Please fill in all date and time fields')
    return
  }

  // Convert 12-hour to 24-hour format
  let hour24 = parseInt(eventDateForm.value.hour)
  if (eventDateForm.value.period === 'PM' && hour24 !== 12) {
    hour24 += 12
  } else if (eventDateForm.value.period === 'AM' && hour24 === 12) {
    hour24 = 0
  }

  // Build ISO date string
  const month = String(eventDateForm.value.month).padStart(2, '0')
  const day = String(eventDateForm.value.day).padStart(2, '0')
  const hour = String(hour24).padStart(2, '0')
  const minute = eventDateForm.value.minute

  const dateString = `${eventDateForm.value.year}-${month}-${day}T${hour}:${minute}:00`
  const eventTimestamp = Math.floor(new Date(dateString).getTime() / 1000)
  
  await groupsStore.createGroupEvent(route.params.id, {
    title: newEvent.value.title.trim(),
    description: newEvent.value.description.trim(),
    event_date: eventTimestamp
  })
  
  newEvent.value = { title: '', description: '', event_date: '' }
  eventDateForm.value = { month: '', day: '', year: '', hour: '', minute: '', period: '' }
  showEventModal.value = false
}

const handleInviteUser = async (userId) => {
  const result = await groupsStore.inviteUser(route.params.id, userId)
  if (result) {
    const user = inviteUsers.value.find(u => u.id === userId)
    if (user) user.invited = true
  }
}

const fetchInviteUsers = async () => {
  try {
    const users = await api.getUsers()
    if (users) {
      const memberIds = new Set(group.value?.members?.map(m => m.id) || [])
      inviteUsers.value = users
        .filter(u => !memberIds.has(u.id) && u.id !== authStore.user?.id)
        .map(u => ({ ...u, invited: false }))
    }
  } catch (e) {
    console.error('Failed to fetch users:', e)
  }
}

watch(showInviteModal, (val) => {
  if (val) {
    fetchInviteUsers()
  }
})



watch(() => route.params.id, () => {
  loadGroupData()
})

onMounted(() => {
  loadGroupData()
})
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.group-details-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 24px 24px;
}

.group-header {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  margin-bottom: 16px;
}

.group-cover {
  height: 200px;
  background: linear-gradient(135deg, var(--muted-teal) 0%, var(--ink-black) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.group-avatar {
  width: 120px;
  height: 120px;
  background: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 3rem;
  font-weight: 700;
  color: var(--muted-teal);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.group-info {
  padding: 24px;
}

.group-info h1 {
  margin: 0 0 12px 0;
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.02em;
}

.group-description {
  margin: 0 0 16px 0;
  color: rgba(13, 19, 33, 0.7);
  line-height: 1.6;
  font-size: 0.9375rem;
}

.group-meta {
  display: flex;
  gap: 12px;
  color: rgba(13, 19, 33, 0.6);
  font-size: 0.875rem;
  margin-bottom: 20px;
  font-weight: 500;
}

.group-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.group-tabs {
  display: flex;
  background: white;
  border-radius: 16px;
  padding: 6px;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.08);
  margin-bottom: 16px;
  gap: 6px;
}

.tab {
  flex: 1;
  padding: 12px 20px;
  background: transparent;
  border: none;
  border-radius: 12px;
  font-size: 0.9375rem;
  font-weight: 600;
  color: rgba(13, 19, 33, 0.6);
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  letter-spacing: 0.01em;
}

.tab:hover {
  background: var(--lavender-mist);
  color: var(--ink-black);
}

.tab.active {
  background: var(--honey-bronze);
  color: var(--ink-black);
  font-weight: 700;
}

.badge {
  background: #e41e3f;
  color: white;
  font-size: 0.75rem;
  padding: 3px 8px;
  border-radius: 12px;
  font-weight: 700;
  min-width: 20px;
  text-align: center;
}

.tab-content {
  min-height: 300px;
}

.create-post-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.08);
  margin-bottom: 16px;
}

.create-post-input {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  background: var(--lavender-mist);
  border-radius: 24px;
  cursor: pointer;
  transition: all 0.2s;
}

.create-post-input:hover {
  background: rgba(246, 240, 249, 0.8);
  transform: translateY(-1px);
}

.create-post-input span {
  color: rgba(13, 19, 33, 0.6);
  font-size: 0.9375rem;
}

.avatar-placeholder {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, var(--muted-teal), var(--ink-black));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  flex-shrink: 0;
  font-size: 1rem;
}

.members-list {
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.08);
  overflow: hidden;
}

.member-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(13, 19, 33, 0.08);
  transition: background 0.2s;
}

.member-item:last-child {
  border-bottom: none;
}

.member-item:hover {
  background: var(--lavender-mist);
}

.member-avatar img,
.member-avatar .avatar-placeholder {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.member-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.member-name {
  font-weight: 700;
  color: var(--ink-black);
  font-size: 0.9375rem;
}

.member-role {
  font-size: 0.8125rem;
  color: var(--muted-teal);
  font-weight: 600;
}

.requests-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.request-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.08);
}

.request-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.request-text {
  font-size: 0.8125rem;
  color: rgba(13, 19, 33, 0.6);
}

.request-actions {
  display: flex;
  gap: 8px;
}

.create-event-btn {
  margin-bottom: 16px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

/* leave button css */
.btn-danger {
  padding: 12px 24px;
  background: #e41e3f;
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 0.9375rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  letter-spacing: 0.01em;
}
.invite-list {
  max-height: 360px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 4px 0;
}

.invite-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
  background: var(--lavender-mist);
  border-radius: 12px;
  transition: background 0.2s;
}

.invite-item:hover {
  background: rgba(246, 240, 249, 0.8);
}

.invite-item .member-name {
  flex: 1;
}

.non-member-view {
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.08);
  overflow: hidden;
}

.locked-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 32px;
  color: rgba(13, 19, 33, 0.6);
}

.locked-content svg {
  margin-bottom: 20px;
  opacity: 0.4;
  color: var(--ink-black);
}

.locked-content p {
  font-size: 0.9375rem;
  margin: 0;
}

.empty-state {
  text-align: center;
  padding: 48px 32px;
  color: rgba(13, 19, 33, 0.6);
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.08);
}

.empty-state p {
  margin: 0;
  font-size: 0.9375rem;
}

.btn-primary {
  padding: 12px 24px;
  background: var(--honey-bronze);
  color: var(--ink-black);
  border: none;
  border-radius: 12px;
  font-size: 0.9375rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  letter-spacing: 0.01em;
}

.btn-primary:hover {
  background: #f7c570;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(246, 189, 96, 0.35);
}

.btn-primary:disabled {
  background: rgba(13, 19, 33, 0.2);
  cursor: not-allowed;
  transform: none;
}

.btn-secondary {
  padding: 12px 24px;
  background: white;
  color: var(--ink-black);
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.9375rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  letter-spacing: 0.01em;
}

.btn-secondary:hover {
  background: var(--lavender-mist);
  border-color: var(--muted-teal);
  transform: translateY(-2px);
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.btn-sm {
  padding: 8px 16px;
  font-size: 0.875rem;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(13, 19, 33, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 3000;
  backdrop-filter: blur(4px);
}

.modal {
  background: white;
  border-radius: 16px;
  width: 100%;
  max-width: 520px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 2px solid rgba(13, 19, 33, 0.08);
}

.modal-header h2 {
  margin: 0;
  font-size: 1.375rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.01em;
}

.close-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: rgba(13, 19, 33, 0.08);
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--ink-black);
  transition: all 0.2s;
}

.close-btn:hover {
  background: rgba(13, 19, 33, 0.15);
  transform: scale(1.1);
}

.modal-form {
  padding: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 700;
  color: var(--ink-black);
  margin-bottom: 8px;
  letter-spacing: 0.01em;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.9375rem;
  outline: none;
  resize: vertical;
  box-sizing: border-box;
  font-family: inherit;
  transition: all 0.2s;
  color: var(--ink-black);
}

.form-group input:focus,
.form-group textarea:focus {
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

/* Professional Date Picker */
.date-picker-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1.5fr;
  gap: 12px;
}

.date-select {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.9375rem;
  font-family: inherit;
  color: var(--ink-black);
  background: white;
  cursor: pointer;
  transition: all 0.2s ease;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg width='12' height='8' viewBox='0 0 12 8' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M1 1L6 6L11 1' stroke='%230d1321' stroke-width='2' stroke-linecap='round'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 36px;
}

.date-select:focus {
  outline: none;
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

/* Professional Time Picker */
.time-picker-row {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 12px;
}

.time-select {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.9375rem;
  font-family: inherit;
  color: var(--ink-black);
  background: white;
  cursor: pointer;
  transition: all 0.2s ease;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg width='12' height='8' viewBox='0 0 12 8' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M1 1L6 6L11 1' stroke='%230d1321' stroke-width='2' stroke-linecap='round'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 36px;
  text-align: center;
}

.time-select:focus {
  outline: none;
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

.file-label {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: var(--lavender-mist);
  border-radius: 12px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 700;
  color: var(--ink-black);
  transition: all 0.2s;
  letter-spacing: 0.01em;
}

.file-label:hover {
  background: rgba(246, 240, 249, 0.8);
  transform: translateY(-1px);
}

.file-name {
  margin-left: 12px;
  font-size: 0.8125rem;
  color: rgba(13, 19, 33, 0.6);
  font-weight: 500;
}

.image-preview {
  position: relative;
  margin-bottom: 16px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.1);
}

.image-preview img {
  width: 100%;
  max-height: 300px;
  object-fit: cover;
  display: block;
}

.remove-image-btn {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--ink-black);
  color: white;
  border: 2px solid white;
  cursor: pointer;
  font-size: 1.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  font-weight: 700;
  line-height: 1;
}

.remove-image-btn:hover {
  transform: scale(1.1);
  background: #1a2332;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.search-input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.9375rem;
  outline: none;
  box-sizing: border-box;
  transition: all 0.2s;
}

.search-input:focus {
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

.loading-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  gap: 16px;
}

.loading-page p {
  color: rgba(13, 19, 33, 0.6);
  font-size: 0.9375rem;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid rgba(13, 19, 33, 0.1);
  border-top-color: var(--honey-bronze);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-page {
  text-align: center;
  padding: 80px 32px;
  color: rgba(13, 19, 33, 0.6);
}

.error-page p {
  margin: 0 0 24px 0;
  font-size: 1.125rem;
}

@media (max-width: 768px) {
  .group-details-page {
    padding: 0 16px 16px;
  }

  .group-cover {
    height: 150px;
  }

  .group-avatar {
    width: 80px;
    height: 80px;
    font-size: 2rem;
  }

  .group-info {
    padding: 20px;
  }

  .group-info h1 {
    font-size: 1.5rem;
  }

  .group-actions {
    flex-direction: column;
  }

  .btn-primary,
  .btn-secondary {
    width: 100%;
    justify-content: center;
  }

  .group-tabs {
    overflow-x: auto;
    padding: 4px;
  }

  .tab {
    white-space: nowrap;
    padding: 10px 16px;
    font-size: 0.875rem;
  }

  .modal {
    max-width: calc(100vw - 32px);
  }

  .modal-header {
    padding: 16px 20px;
  }

  .modal-form {
    padding: 20px;
  }

  .date-picker-row,
  .time-picker-row {
    grid-template-columns: 1fr;
  }
}
</style>