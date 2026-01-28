<!-- filepath: /home/kali/Desktop/social-network/frontend/src/components/groups/GroupEvent.vue -->
<template>
  <div class="group-event">
    <!-- Professional Date Badge -->
    <div class="event-date-badge">
      <div class="badge-month">{{ eventMonth }}</div>
      <div class="badge-day">{{ eventDay }}</div>
      <div class="badge-year">{{ eventYear }}</div>
    </div>
    
    <div class="event-content">
      <h3 class="event-title">{{ event.title }}</h3>
      <p class="event-description">{{ event.description }}</p>
      
      <!-- Professional Date & Time Display -->
      <div class="event-datetime">
        <div class="datetime-item">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
            <line x1="16" y1="2" x2="16" y2="6"/>
            <line x1="8" y1="2" x2="8" y2="6"/>
            <line x1="3" y1="10" x2="21" y2="10"/>
          </svg>
          <span>{{ eventDateFormatted }}</span>
        </div>
        <div class="datetime-item">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
          </svg>
          <span>{{ eventTimeFormatted }}</span>
        </div>
      </div>
      
      <!-- Response Counts -->
      <div class="event-responses">
        <div class="response-badge going">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <polyline points="20 6 9 17 4 12"/>
          </svg>
          <span><strong>{{ event.going_count || 0 }}</strong> Going</span>
        </div>
        <div class="response-badge not-going">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
          <span><strong>{{ event.not_going_count || 0 }}</strong> Not Going</span>
        </div>
      </div>
      
      <!-- Action Buttons -->
      <div class="event-actions">
        <button 
          :class="['response-btn', 'btn-going', { active: event.user_response === 'going' }]"
          @click="respond('going')"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <polyline points="20 6 9 17 4 12"/>
          </svg>
          I'm Going
        </button>
        <button 
          :class="['response-btn', 'btn-not-going', { active: event.user_response === 'not_going' }]"
          @click="respond('not_going')"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
          Can't Go
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useGroupsStore } from '@/stores/groups'

const props = defineProps({
  event: {
    type: Object,
    required: true
  }
})

const groupsStore = useGroupsStore()

// FIX: Convert Unix timestamp (seconds) to milliseconds
const eventDate = computed(() => {
  // Multiply by 1000 to convert seconds to milliseconds
  return new Date(props.event.event_date * 1000)
})

const eventMonth = computed(() => {
  return eventDate.value.toLocaleDateString('en-US', { month: 'short' }).toUpperCase()
})

const eventDay = computed(() => {
  return eventDate.value.getDate()
})

const eventYear = computed(() => {
  return eventDate.value.getFullYear()
})

const eventDateFormatted = computed(() => {
  return eventDate.value.toLocaleDateString('en-US', {
    weekday: 'long',
    month: 'long',
    day: 'numeric',
    year: 'numeric'
  })
})

const eventTimeFormatted = computed(() => {
  return eventDate.value.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
  })
})

const respond = async (response) => {
  await groupsStore.respondToEvent(props.event.group_id, props.event.id, response)
}
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.group-event {
  display: flex;
  gap: 20px;
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.08);
  margin-bottom: 16px;
  border: 1px solid rgba(13, 19, 33, 0.06);
  transition: all 0.2s ease;
}

.group-event:hover {
  box-shadow: 0 4px 16px rgba(13, 19, 33, 0.12);
  transform: translateY(-2px);
}

/* Professional Date Badge */
.event-date-badge {
  width: 80px;
  min-height: 90px;
  background: linear-gradient(135deg, var(--lavender-mist) 0%, rgba(246, 240, 249, 0.5) 100%);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 2px solid rgba(13, 19, 33, 0.08);
  padding: 8px;
  gap: 2px;
}

.badge-month {
  font-size: 0.6875rem;
  font-weight: 800;
  color: var(--honey-bronze);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.badge-day {
  font-size: 2rem;
  font-weight: 800;
  color: var(--ink-black);
  line-height: 1;
  letter-spacing: -0.03em;
}

.badge-year {
  font-size: 0.75rem;
  font-weight: 600;
  color: rgba(13, 19, 33, 0.6);
  letter-spacing: 0.02em;
}

/* Event Content */
.event-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.event-title {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.02em;
  line-height: 1.3;
}

.event-description {
  margin: 0;
  font-size: 0.9375rem;
  color: rgba(13, 19, 33, 0.7);
  line-height: 1.5;
}

/* Professional Date & Time Display */
.event-datetime {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  padding: 12px 0;
  border-top: 1px solid rgba(13, 19, 33, 0.06);
  border-bottom: 1px solid rgba(13, 19, 33, 0.06);
}

.datetime-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.875rem;
  color: rgba(13, 19, 33, 0.7);
  font-weight: 600;
}

.datetime-item svg {
  color: var(--muted-teal);
  flex-shrink: 0;
}

/* Response Counts */
.event-responses {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.response-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 0.8125rem;
  font-weight: 500;
}

.response-badge.going {
  background: rgba(146, 191, 177, 0.1);
  color: var(--muted-teal);
}

.response-badge.going svg {
  color: var(--muted-teal);
}

.response-badge.not-going {
  background: rgba(211, 47, 47, 0.1);
  color: #d32f2f;
}

.response-badge.not-going svg {
  color: #d32f2f;
}

.response-badge strong {
  font-weight: 700;
}

/* Action Buttons */
.event-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-top: 4px;
}

.response-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  background: white;
  color: var(--ink-black);
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  letter-spacing: 0.01em;
  flex: 1;
  min-width: 120px;
  justify-content: center;
}

.response-btn:hover {
  background: rgba(13, 19, 33, 0.04);
  border-color: rgba(13, 19, 33, 0.2);
  transform: translateY(-1px);
}

.response-btn:active {
  transform: translateY(0);
}

.response-btn svg {
  flex-shrink: 0;
}

.btn-going.active {
  background: var(--muted-teal);
  color: white;
  border-color: var(--muted-teal);
  box-shadow: 0 2px 8px rgba(146, 191, 177, 0.3);
}

.btn-going.active:hover {
  background: #7da99c;
  border-color: #7da99c;
  box-shadow: 0 4px 12px rgba(146, 191, 177, 0.4);
}

.btn-not-going.active {
  background: #d32f2f;
  color: white;
  border-color: #d32f2f;
  box-shadow: 0 2px 8px rgba(211, 47, 47, 0.3);
}

.btn-not-going.active:hover {
  background: #b71c1c;
  border-color: #b71c1c;
  box-shadow: 0 4px 12px rgba(211, 47, 47, 0.4);
}

/* Responsive Design */
@media (max-width: 640px) {
  .group-event {
    flex-direction: column;
    gap: 16px;
    padding: 16px;
  }

  .event-date-badge {
    width: 100%;
    min-height: auto;
    flex-direction: row;
    justify-content: flex-start;
    padding: 12px 16px;
    gap: 12px;
  }

  .badge-month {
    font-size: 0.75rem;
  }

  .badge-day {
    font-size: 1.75rem;
  }

  .badge-year {
    font-size: 0.8125rem;
    margin-left: auto;
  }

  .event-title {
    font-size: 1.125rem;
  }

  .event-datetime {
    flex-direction: column;
    gap: 10px;
  }

  .event-actions {
    flex-direction: column;
    gap: 8px;
  }

  .response-btn {
    width: 100%;
    min-width: auto;
  }
}

@media (max-width: 480px) {
  .event-responses {
    flex-direction: column;
    gap: 8px;
  }

  .response-badge {
    width: 100%;
  }
}
</style>