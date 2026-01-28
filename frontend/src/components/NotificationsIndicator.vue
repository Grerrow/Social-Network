<!-- filepath: /home/kali/Desktop/social-network/frontend/src/components/NotificationsIndicator.vue -->
<template>
  <div class="notification-bell" @click="openNotifications">
    <svg class="bell-icon" width="30" height="30" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
      <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
    </svg>
    <span v-if="unreadCount > 0" class="badge">
      {{ unreadCount > 99 ? '99+' : unreadCount }}
    </span>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useNotificationsStore } from '@/stores/notifications'
import { useRouter } from 'vue-router'

const store = useNotificationsStore()
const router = useRouter()

const unreadCount = computed(() => store.unreadCount)

onMounted(() => store.fetchNotifications())

const openNotifications = () => {
  router.push('/notifications')
}
</script>

<style scoped>
.notification-bell {
  position: relative;
  cursor: pointer;
  padding: 10px;
  border-radius: 12px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--lavender-mist);
}



.notification-bell:active {
  transform: translateY(0);
}

.bell-icon {
  color: var(--ink-black);
}

.badge {
  position: absolute;
  top: 4px;
  right: 4px;
  background: linear-gradient(135deg, #ef5350, #d32f2f);
  color: white;
  font-size: 0.6875rem;
  padding: 2px 6px;
  border-radius: 10px;
  font-weight: 800;
  min-width: 18px;
  text-align: center;
  border: 2px solid white;
  letter-spacing: -0.02em;
  line-height: 1.2;
  /* box-shadow: 0 2px 6px rgba(211, 47, 47, 0.4); */
}

@media (max-width: 768px) {
  .notification-bell {
    padding: 8px;
  }

  .bell-icon {
    width: 18px;
    height: 18px;
  }

  .badge {
    font-size: 0.625rem;
    padding: 2px 5px;
    min-width: 16px;
  }
}
</style>