<!-- filepath: /home/kali/Desktop/social-network/frontend/src/views/GroupsBrowse.vue -->
<template>
  <div class="groups-page">
    <!-- Header -->
    <div class="groups-header">
      <h1>Groups</h1>
      <button @click="showCreateModal = true" class="btn-create">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        Create Group
      </button>
    </div>

    <!-- Main Layout: Sidebar + Content -->
    <div class="groups-layout">
      <!-- Left Sidebar: My Groups -->
      <aside class="groups-sidebar">
        <div class="sidebar-section">
          <h2>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
            </svg>
            My Groups
          </h2>
          
          <div v-if="myGroups.length === 0" class="sidebar-empty">
            <p>You haven't joined any groups yet</p>
          </div>
          
          <div v-else class="sidebar-list">
            <div
              v-for="group in myGroups"
              :key="group.id"
              @click="goToGroup(group.id)"
              class="sidebar-group-item"
            >
              <div class="group-icon">
                {{ group.title?.charAt(0).toUpperCase() }}
              </div>
              <div class="group-info">
                <h4>{{ group.title }}</h4>
                <span>{{ group.member_count || 0 }} members</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Pending Invites in Sidebar -->
        <div v-if="pendingInvites.length > 0" class="sidebar-section invites-section">
          <h2>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
              <polyline points="22 4 12 14.01 9 11.01"/>
            </svg>
            Invitations
            <span class="badge">{{ pendingInvites.length }}</span>
          </h2>
          
          <div class="sidebar-list">
            <div
              v-for="group in pendingInvites"
              :key="group.id"
              class="sidebar-group-item invite-item"
            >
              <div class="group-icon invite-icon">
                {{ group.title?.charAt(0).toUpperCase() }}
              </div>
              <div class="group-info">
                <h4>{{ group.title }}</h4>
                <div class="invite-actions">
                  <button @click.stop="handleAcceptInvite(group.id)" class="btn-accept">
                    Accept
                  </button>
                  <button @click.stop="handleDeclineInvite(group.id)" class="btn-decline">
                    Decline
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </aside>

      <!-- Main Content: Discover Groups -->
      <main class="groups-main">
        <div class="discover-header">
          <h2>Discover Groups</h2>
          <div class="search-wrapper">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="search-icon">
              <circle cx="11" cy="11" r="8"/>
              <path d="m21 21-4.35-4.35"/>
            </svg>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search groups..."
              class="search-input"
            />
          </div>
        </div>

        <div v-if="isLoading" class="loading">
          <div class="spinner"></div>
          <p>Loading groups...</p>
        </div>

        <div v-else-if="discoverGroups.length === 0" class="empty-state">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
          <p>No groups found</p>
          <span>Try adjusting your search</span>
        </div>

        <div v-else class="groups-grid">
          <div
            v-for="group in discoverGroups"
            :key="group.id"
            @click="goToGroup(group.id)"
            class="group-card"
          >
            <div class="card-header">
              <div class="group-avatar">
                {{ group.title?.charAt(0).toUpperCase() }}
              </div>
              <div class="card-badge" v-if="group.has_requested">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                </svg>
                Requested
              </div>
            </div>
            
            <h3>{{ group.title }}</h3>
            <p class="card-description">{{ group.description }}</p>
            
            <div class="card-footer">
              <div class="member-count">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                  <circle cx="9" cy="7" r="4"/>
                  <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                </svg>
                {{ group.member_count || 0 }} members
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>

    <!-- Create Group Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal">
        <div class="modal-header">
          <h2>Create New Group</h2>
          <button @click="showCreateModal = false" class="close-btn">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
            </svg>
          </button>
        </div>

        <form @submit.prevent="handleCreateGroup" class="modal-form">
          <div class="form-group">
            <label for="groupTitle">Group Name</label>
            <input
              id="groupTitle"
              v-model="newGroup.title"
              type="text"
              required
              placeholder="Enter group name"
            />
          </div>

          <div class="form-group">
            <label for="groupDescription">Description</label>
            <textarea
              id="groupDescription"
              v-model="newGroup.description"
              rows="4"
              placeholder="Describe your group..."
              required
            ></textarea>
          </div>

          <div class="modal-actions">
            <button type="button" @click="showCreateModal = false" class="btn-cancel">
              Cancel
            </button>
            <button type="submit" :disabled="isCreating" class="btn-submit">
              {{ isCreating ? 'Creating...' : 'Create Group' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGroupsStore } from '@/stores/groups'

const router = useRouter()
const groupsStore = useGroupsStore()

const searchQuery = ref('')
const showCreateModal = ref(false)
const isCreating = ref(false)
const newGroup = ref({
  group_name: '',
  title: '',
  description: ''
})

const isLoading = computed(() => groupsStore.isLoading)

const myGroups = computed(() => {
  return groupsStore.groups.filter(g => g.is_member)
})

const pendingInvites = computed(() => {
  return groupsStore.groups.filter(g => g.is_invited)
})

const discoverGroups = computed(() => {
  let filtered = groupsStore.groups.filter(g => !g.is_member && !g.is_invited)
  
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(g => 
      g.title?.toLowerCase().includes(query) ||
      g.group_name?.toLowerCase().includes(query) ||
      g.description?.toLowerCase().includes(query)
    )
  }
  
  return filtered
})

const goToGroup = (groupId) => {
  router.push(`/groups/${groupId}`)
}

const handleCreateGroup = async () => {
  if (!newGroup.value.title.trim() || !newGroup.value.description.trim()) return
  
  isCreating.value = true
  const group = await groupsStore.createGroup(
    newGroup.value.title.trim(),
    newGroup.value.title.trim(),
    newGroup.value.description.trim()
  )
  isCreating.value = false
  
  if (group) {
    showCreateModal.value = false
    newGroup.value = { group_name: '', title: '', description: '' }
    router.push(`/groups/${group.id}`)
  }
}

const handleAcceptInvite = async (groupId) => {
  await groupsStore.acceptInvite(groupId)
  await groupsStore.fetchGroups()
}

const handleDeclineInvite = async (groupId) => {
  await groupsStore.declineInvite(groupId)
  await groupsStore.fetchGroups()
}

onMounted(() => {
  groupsStore.fetchGroups()
})
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

/* ===== Main Layout ===== */
.groups-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 32px 16px;
}

.groups-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.groups-header h1 {
  margin: 0;
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.02em;
}

.btn-create {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--honey-bronze);
  color: var(--ink-black);
  font-weight: 700;
  border: none;
  padding: 12px 20px;
  border-radius: 12px;
  cursor: pointer;
  text-decoration: none;
  font-size: 0.9375rem;
  transition: all 0.2s ease;
  letter-spacing: 0.01em;
}

.btn-create:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(246, 189, 96, 0.35);
  background: #f7c570;
}

/* ===== Two Column Layout ===== */
.groups-layout {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 24px;
  align-items: start;
}

/* ===== Sidebar ===== */
.groups-sidebar {
  position: sticky;
  top: 80px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.sidebar-section {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
}

.sidebar-section h2 {
  margin: 0 0 16px 0;
  font-size: 1rem;
  font-weight: 700;
  color: var(--ink-black);
  display: flex;
  align-items: center;
  gap: 10px;
  letter-spacing: 0.01em;
}

.sidebar-section h2 svg {
  color: var(--muted-teal);
}

.sidebar-empty {
  padding: 32px 16px;
  text-align: center;
}

.sidebar-empty p {
  margin: 0;
  font-size: 0.875rem;
  color: rgba(13, 19, 33, 0.5);
  line-height: 1.5;
}

.sidebar-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.sidebar-group-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  background: var(--lavender-mist);
}

.sidebar-group-item:hover {
  background: rgba(246, 240, 249, 0.6);
  transform: translateX(4px);
}

.group-icon {
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, var(--muted-teal), var(--ink-black));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 1rem;
  flex-shrink: 0;
}

.invite-icon {
  background: linear-gradient(135deg, var(--honey-bronze), #e8a84d);
}

.group-info {
  flex: 1;
  min-width: 0;
}

.group-info h4 {
  margin: 0 0 4px 0;
  font-size: 0.9375rem;
  font-weight: 700;
  color: var(--ink-black);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  letter-spacing: 0.01em;
}

.group-info span {
  font-size: 0.8125rem;
  color: rgba(13, 19, 33, 0.6);
  font-weight: 500;
}

/* Invites Section */
.invites-section {
  background: linear-gradient(135deg, rgba(246, 189, 96, 0.1), rgba(146, 191, 177, 0.1));
  border: 2px solid rgba(246, 189, 96, 0.3);
}

.invites-section h2 {
  position: relative;
}

.badge {
  position: absolute;
  right: 0;
  background: #e41e3f;
  color: white;
  font-size: 0.75rem;
  padding: 3px 8px;
  border-radius: 12px;
  font-weight: 700;
}

.invite-item {
  background: white;
  border: 2px solid rgba(246, 189, 96, 0.2);
}

.invite-item:hover {
  background: white;
  border-color: var(--honey-bronze);
}

.invite-actions {
  display: flex;
  gap: 6px;
  margin-top: 8px;
}

.btn-accept,
.btn-decline {
  flex: 1;
  padding: 6px 12px;
  border: none;
  border-radius: 8px;
  font-size: 0.8125rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  letter-spacing: 0.01em;
}

.btn-accept {
  background: var(--muted-teal);
  color: white;
}

.btn-accept:hover {
  background: #7fa89d;
  transform: translateY(-1px);
}

.btn-decline {
  background: rgba(13, 19, 33, 0.08);
  color: var(--ink-black);
}

.btn-decline:hover {
  background: rgba(13, 19, 33, 0.15);
}

/* ===== Main Content ===== */
.groups-main {
  min-height: 500px;
}

.discover-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.discover-header h2 {
  margin: 0;
  font-size: 1.375rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.01em;
}

.search-wrapper {
  position: relative;
  flex: 1;
  max-width: 400px;
}

.search-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: rgba(13, 19, 33, 0.4);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 12px 16px 12px 44px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 24px;
  font-size: 0.9375rem;
  transition: all 0.2s ease;
  font-family: inherit;
  color: var(--ink-black);
  background: white;
}

.search-input:focus {
  outline: none;
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

/* ===== Groups Grid ===== */
.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 20px;
}

.group-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.group-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(13, 19, 33, 0.15);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 4px;
}

.group-avatar {
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, var(--muted-teal), var(--ink-black));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 1.25rem;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(13, 19, 33, 0.15);
}

.card-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: rgba(146, 191, 177, 0.15);
  color: var(--muted-teal);
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.02em;
}

.group-card h3 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.01em;
  line-height: 1.3;
}

.card-description {
  margin: 0;
  font-size: 0.875rem;
  color: rgba(13, 19, 33, 0.7);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid rgba(13, 19, 33, 0.08);
}

.member-count {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.8125rem;
  color: rgba(13, 19, 33, 0.6);
  font-weight: 600;
}

.member-count svg {
  color: var(--muted-teal);
}

/* ===== Loading & Empty States ===== */
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  gap: 16px;
}

.loading p {
  color: rgba(13, 19, 33, 0.6);
  font-size: 0.9375rem;
  margin: 0;
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

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  background: white;
  border-radius: 16px;
  padding: 48px 32px;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  gap: 12px;
}

.empty-state svg {
  color: rgba(13, 19, 33, 0.2);
  margin-bottom: 8px;
}

.empty-state p {
  margin: 0;
  color: rgba(13, 19, 33, 0.7);
  font-size: 1.125rem;
  font-weight: 700;
}

.empty-state span {
  color: rgba(13, 19, 33, 0.5);
  font-size: 0.875rem;
}

/* ===== Modal ===== */
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
  padding: 32px;
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
  padding: 24px 32px;
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
  padding: 24px 32px 32px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-weight: 700;
  color: var(--ink-black);
  font-size: 0.875rem;
  letter-spacing: 0.01em;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.9375rem;
  transition: all 0.2s ease;
  box-sizing: border-box;
  font-family: inherit;
  color: var(--ink-black);
  background: white;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

.form-group textarea {
  resize: vertical;
  line-height: 1.5;
  min-height: 100px;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 8px;
}

.btn-cancel {
  padding: 12px 24px;
  background: white;
  color: var(--ink-black);
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-weight: 700;
  font-size: 0.9375rem;
  cursor: pointer;
  transition: all 0.2s ease;
  letter-spacing: 0.01em;
}

.btn-cancel:hover {
  background: var(--lavender-mist);
  border-color: var(--muted-teal);
}

.btn-submit {
  padding: 12px 24px;
  background: var(--honey-bronze);
  color: var(--ink-black);
  border: none;
  border-radius: 12px;
  font-weight: 700;
  font-size: 0.9375rem;
  cursor: pointer;
  transition: all 0.2s ease;
  letter-spacing: 0.02em;
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(246, 189, 96, 0.35);
  background: #f7c570;
}

.btn-submit:disabled {
  background: rgba(13, 19, 33, 0.1);
  color: rgba(13, 19, 33, 0.4);
  cursor: not-allowed;
}

/* ===== Responsive Design ===== */
@media (max-width: 1100px) {
  .groups-layout {
    grid-template-columns: 280px 1fr;
    gap: 20px;
  }

  .groups-grid {
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  }
}

@media (max-width: 900px) {
  .groups-layout {
    grid-template-columns: 1fr;
  }

  .groups-sidebar {
    position: static;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 16px;
  }

  .discover-header {
    flex-direction: column;
    align-items: stretch;
  }

  .search-wrapper {
    max-width: 100%;
  }
}

@media (max-width: 600px) {
  .groups-page {
    padding: 16px 12px;
  }

  .groups-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
    margin-bottom: 24px;
  }

  .btn-create {
    justify-content: center;
  }

  .groups-sidebar {
    grid-template-columns: 1fr;
  }

  .groups-grid {
    grid-template-columns: 1fr;
  }

  .modal {
    max-width: calc(100vw - 32px);
  }

  .modal-header,
  .modal-form {
    padding: 20px 24px;
  }

  .modal-actions {
    flex-direction: column;
  }

  .btn-cancel,
  .btn-submit {
    width: 100%;
  }
}
</style>