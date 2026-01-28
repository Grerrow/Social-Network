<!-- filepath: /home/kali/Desktop/social-network/frontend/src/components/groups/GroupComment.vue -->
<template>
  <div class="comment">
    <div class="comment-header">
      <router-link 
        :to="`/profile?user_id=${comment.author?.id || comment.user_id}`" 
        class="avatar-link"
      >
        <div class="avatar">
          <img 
            :src="comment.author?.avatar ? getImageUrl(comment.author.avatar) : '/images/placeholder.png'" 
            :alt="comment.author?.username || ''"
          />
        </div>
      </router-link>
      
      <div class="comment-bubble">
        <div class="comment-user">
          <router-link 
            :to="`/profile?user_id=${comment.author?.id || comment.user_id}`" 
            class="username"
          >
            {{ comment.author?.username || 'Unknown' }}
          </router-link>
        </div>
        <p v-if="comment.content" class="comment-text">{{ comment.content }}</p>
        
        <!-- Comment Image -->
        <div v-if="comment.image" class="comment-image" @click="showImageModal = true">
          <img 
            :src="getImageUrl(comment.image)" 
            :alt="comment.author?.username + '\'s image'"
            loading="lazy"
          />
        </div>
      </div>
    </div>
    
    <span class="comment-time">{{ formatTime(comment.created_at) }}</span>

    <!-- Image Modal -->
    <teleport to="body">
      <div v-if="showImageModal" class="image-modal" @click="showImageModal = false">
        <div class="modal-backdrop"></div>
        <div class="modal-content" @click.stop>
          <button class="modal-close" @click="showImageModal = false">×</button>
          <img 
            :src="getImageUrl(comment.image)" 
            :alt="comment.author?.username + '\'s image'"
          />
        </div>
      </div>
    </teleport>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useApi } from '../../composables/useApi'

const { getImageUrl } = useApi()

defineProps({
  comment: {
    type: Object,
    required: true
  }
})

const showImageModal = ref(false)

const formatTime = (timestamp) => {
  if (!timestamp) return 'just now'
  const date = new Date(timestamp * 1000)
  const now = new Date()
  const diff = Math.floor((now - date) / 1000)
  
  if (diff < 60) return 'just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  if (diff < 604800) return `${Math.floor(diff / 86400)}d ago`
  return date.toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric',
    year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined
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

.comment {
  padding: 8px 0;
}

.comment-header {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.avatar-link {
  text-decoration: none;
  flex-shrink: 0;
}

.avatar {
  width: 36px;
  height: 36px;
  background: var(--lavender-mist);
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  flex-shrink: 0;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-weight: 700;
  font-size: 0.875rem;
  color: var(--ink-black);
}

.comment-bubble {
  flex: 1;
  background: var(--lavender-mist);
  border-radius: 12px;
  padding: 12px;
  min-width: 0;
}

.comment-user .username {
  font-weight: 700;
  font-size: 0.875rem;
  color: var(--ink-black);
  text-decoration: none;
  letter-spacing: 0.01em;
}

.comment-user .username:hover {
  text-decoration: underline;
  color: var(--muted-teal);
}

.comment-text {
  margin: 4px 0 0;
  font-size: 0.9375rem;
  line-height: 1.5;
  color: var(--ink-black);
  word-wrap: break-word;
  white-space: pre-wrap;
}

.comment-image {
  margin-top: 8px;
  cursor: pointer;
  transition: opacity 0.2s ease;
}

.comment-image:hover {
  opacity: 0.95;
}

.comment-image img {
  max-width: 100%;
  max-height: 200px;
  border-radius: 8px;
  object-fit: cover;
  display: block;
  box-shadow: 0 2px 8px rgba(13, 19, 33, 0.1);
}

.comment-time {
  display: block;
  font-size: 0.75rem;
  color: rgba(13, 19, 33, 0.5);
  margin-top: 4px;
  margin-left: 48px;
}

/* Image Modal */
.image-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-backdrop {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(13, 19, 33, 0.95);
  backdrop-filter: blur(8px);
  cursor: pointer;
}

.modal-content {
  position: relative;
  max-width: 90vw;
  max-height: 90vh;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
}

.modal-close {
  position: absolute;
  top: -64px;
  right: 0;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: white;
  color: var(--ink-black);
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  font-weight: 700;
  line-height: 1;
  transition: all 0.2s ease;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.modal-close:hover {
  background: var(--honey-bronze);
  transform: scale(1.1);
}

.modal-content img {
  max-width: 100%;
  max-height: 90vh;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
  background: white;
}

@media (max-width: 600px) {
  .comment-header {
    gap: 8px;
  }

  .avatar {
    width: 32px;
    height: 32px;
  }

  .avatar-placeholder {
    font-size: 0.8125rem;
  }

  .comment-bubble {
    padding: 10px;
  }

  .comment-text {
    font-size: 0.875rem;
  }

  .comment-time {
    margin-left: 40px;
  }

  .comment-image img {
    max-height: 180px;
  }

  .modal-close {
    top: -56px;
    width: 40px;
    height: 40px;
    font-size: 1.5rem;
  }
}
</style>