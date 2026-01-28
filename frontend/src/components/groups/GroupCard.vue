<!-- filepath: /home/kali/Desktop/social-network/frontend/src/components/groups/GroupCard.vue -->
<template>
  <div class="group-card" @click="$emit('click')">
    <div class="group-info">
      <h3 class="group-title">{{ group.title }}</h3>
      <span class="group-handle">@{{ group.group_name }}</span>
      <p class="group-description">{{ truncatedDescription }}</p>
      <div class="group-meta">
        <span class="meta-item">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
          {{ group.member_count || 0 }} members
        </span>
      </div>
      
      <!-- Invite Actions -->
      <div v-if="showInviteActions" class="invite-actions" @click.stop>
        <button class="btn-accept" @click="$emit('accept')">
          Accept
        </button>
        <button class="btn-decline" @click="$emit('decline')">
          Decline
        </button>
      </div>
      
      <!-- Status Badges -->
      <div class="group-badges">
        <span v-if="group.is_member" class="badge badge-member">Member</span>
        <span v-else-if="group.has_requested" class="badge badge-pending">Pending</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  group: { type: Object, required: true },
  showInviteActions: { type: Boolean, default: false }
})

defineEmits(['click', 'accept', 'decline'])

const truncatedDescription = computed(() => {
  if (!props.group.description) return ''
  return props.group.description.length > 80
    ? props.group.description.slice(0, 80) + '...'
    : props.group.description
})
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.group-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.08);
  cursor: pointer;
  transition: all 0.2s ease;
  border: 2px solid transparent;
}

.group-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(13, 19, 33, 0.12);
  border-color: var(--honey-bronze);
}

.group-info {
  padding: 16px;
}

.group-title {
  margin: 0 0 4px 0;
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.01em;
}

.group-handle {
  display: block;
  font-size: 0.8125rem;
  color: rgba(13, 19, 33, 0.5);
  margin-bottom: 8px;
  font-weight: 500;
}

.group-description {
  margin: 0 0 12px 0;
  font-size: 0.875rem;
  color: rgba(13, 19, 33, 0.7);
  line-height: 1.5;
}

.group-meta {
  font-size: 0.8125rem;
  color: rgba(13, 19, 33, 0.6);
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-weight: 500;
}

.invite-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.btn-accept,
.btn-decline {
  flex: 1;
  padding: 8px 16px;
  border: none;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  letter-spacing: 0.01em;
}

.btn-accept {
  background: var(--muted-teal);
  color: white;
}

.btn-accept:hover {
  background: #7da99c;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(146, 191, 177, 0.3);
}

.btn-decline {
  background: rgba(13, 19, 33, 0.08);
  color: var(--ink-black);
  border: 2px solid rgba(13, 19, 33, 0.12);
}

.btn-decline:hover {
  background: rgba(13, 19, 33, 0.12);
}

.group-badges {
  margin-top: 12px;
}

.badge {
  display: inline-block;
  padding: 6px 12px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.02em;
}

.badge-member {
  background: rgba(146, 191, 177, 0.15);
  color: var(--muted-teal);
}

.badge-pending {
  background: rgba(246, 189, 96, 0.15);
  color: #cc8c1f;
}

@media (max-width: 600px) {
  .group-info {
    padding: 12px;
  }

  .group-title {
    font-size: 1rem;
  }

  .invite-actions {
    flex-direction: column;
  }
}
</style>