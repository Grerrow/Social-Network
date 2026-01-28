<!-- filepath: /home/kali/Desktop/social-network/frontend/src/components/Comment.vue -->
<template>
  <div class="comment">
    <div class="comment-header">
      <router-link :to="`/profile?user_id=${comment.user_id}`" class="avatar-link">
        <div class="avatar">
          <img 
            :src="comment.avatar ? getImageUrl(comment.avatar) : '/images/placeholder.png'" 
            :alt="comment.username"
          />
        </div>
      </router-link>
      <div class="comment-bubble">
        <div class="comment-user">
          <router-link :to="`/profile?user_id=${comment.user_id}`" class="username">
            {{ comment.username }}
          </router-link>
        </div>
        <p class="comment-text">{{ comment.content }}</p>
        
        <div v-if="comment.image" class="comment-image">
          <img 
            :src="getImageUrl(comment.image)" 
            alt="Comment image"
          />
        </div>
      </div>
    </div>
    <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
  </div>
</template>

<script>

import { useApi } from '@/composables/useApi'

export default {
  name: "Comment",
  props: {
    comment: { type: Object, required: true }
  },


  setup() {
    const { getImageUrl } = useApi()
    return { getImageUrl }
  },
  methods: {
    formatDate(date) {
      if (!date) return "just now"
      const d = new Date(date)
      const now = new Date()
      const diff = Math.floor((now - d) / 1000)
      
      if (diff < 60) return 'just now'
      if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
      if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
      if (diff < 604800) return `${Math.floor(diff / 86400)}d ago`
      return d.toLocaleDateString()
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
}

.comment-image {
  margin-top: 8px;
}

.comment-image img {
  max-width: 100%;
  max-height: 200px;
  border-radius: 8px;
  object-fit: cover;
  display: block;
}

.comment-time {
  display: block;
  font-size: 0.75rem;
  color: rgba(13, 19, 33, 0.5);
  margin-top: 4px;
  margin-left: 48px;
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
}
</style>