<!-- filepath: /home/kali/Desktop/social-network/frontend/src/views/FeedView.vue -->
<template>
  <div class="feed-view">
    <div class="feed-header">
      <h1>Feed</h1>
      <router-link to="/create-post" class="btn-primary">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        Create Post
      </router-link>
    </div>

    <div v-if="!postStore.isLoading && postStore.posts.length === 0" class="empty-state">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="empty-icon">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
        <polyline points="7 10 12 15 17 10"/>
        <line x1="12" y1="15" x2="12" y2="3"/>
      </svg>
      <p>No posts yet. Be the first to post!</p>
    </div>

    <div class="posts-list">
      <div v-for="post in postStore.posts" :key="post.id" class="post-card">
        <div class="post-header">
          <router-link :to="`/profile?user_id=${post.user_id}`" class="user-link">
            <div class="user-avatar">
              <img 
                v-if="post.avatar" 
                :src="getImageUrl(post.avatar) || '/images/placeholder.png'" 
                :alt="post.username"
              />
            </div>
            <div class="user-info">
              <h4>{{ post.username }}</h4>
              <span class="post-time">{{ formatDate(post.created_at) }}</span>
            </div>
          </router-link>
        </div>

        <p v-if="post.content" class="post-content">{{ post.content }}</p>

        <div v-if="post.image" class="post-media">
          <img 
            :src="getImageUrl(post.image)" 
            alt="Post image" 
            class="post-image"
            @click="openImageModal(post.image)"
            loading="lazy"
          />
        </div>

        <div class="post-actions">
          <button @click="toggleComments(post)" class="btn-action">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
            <span>{{ showComments[post.id] ? 'Hide' : 'Show' }} Comments</span>
          </button>
        </div>

        <div v-if="showComments[post.id]" class="comments-section">
          <Comment v-for="comment in post.comments" :key="comment.id" :comment="comment" />

          <div class="comment-form">
            <input
              v-model="newComment[post.id]"
              placeholder="Write a comment..."
              class="comment-input"
              @keyup.enter="submitComment(post)"
            />
            
            <label :for="`comment-image-${post.id}`" class="btn-upload" title="Attach image">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                <circle cx="8.5" cy="8.5" r="1.5"/>
                <polyline points="21 15 16 10 5 21"/>
              </svg>
            </label>
            <input
              :id="`comment-image-${post.id}`"
              type="file"
              accept=".jpg,.jpeg,.png,.gif"
              @change="handleCommentImageChange($event, post.id)"
              hidden
            />
            
            <button class="btn-submit-comment" @click="submitComment(post)">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"/>
              </svg>
            </button>
          </div>

          <div v-if="commentImagePreview[post.id]" class="image-preview">
            <img :src="commentImagePreview[post.id]" alt="Preview" />
            <button class="btn-clear-image" @click="clearCommentImage(post.id)">×</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="modalImage" class="image-modal" @click="modalImage = null">
      <div class="modal-content" @click.stop>
        <button class="modal-close" @click="modalImage = null">×</button>
        <img :src="getImageUrl(modalImage)" alt="Full image" />
      </div>
    </div>
  </div>
</template>

<script>
import { onMounted, reactive, ref } from 'vue'
import { usePostStore } from '../stores/posts'
import Comment from '../components/Comment.vue'
import api from '../services/api'
import { useApi } from '@/composables/useApi'

export default {
  name: 'FeedView',
  components: { Comment },
  setup() {
    const postStore = usePostStore()
    const showComments = reactive({})
    const newComment = reactive({})
    const commentImage = reactive({})
    const commentImagePreview = reactive({})
    const modalImage = ref(null)
    const { getApiUrl, getImageUrl } = useApi()
    onMounted(() => {
      postStore.fetchPosts()
    })

    const toggleComments = async (post) => {
      showComments[post.id] = !showComments[post.id]

      if (showComments[post.id] && !post.comments) {
        try {
          const res = await fetch(getApiUrl(`/comments?post_id=${post.id}`), {
            credentials: 'include'
          })
          post.comments = await res.json()
        } catch (err) {
          console.error('Failed to fetch comments:', err)
          post.comments = []
        }
      }
    }

    const handleCommentImageChange = (e, postId) => {
      const file = e.target.files[0]
      if (!file) return

      commentImage[postId] = file
      commentImagePreview[postId] = URL.createObjectURL(file)
    }

    const clearCommentImage = (postId) => {
      commentImage[postId] = null
      commentImagePreview[postId] = null
    }

    const submitComment = async (post) => {
      const content = newComment[post.id]
      if (!content?.trim() && !commentImage[post.id]) return

      try {
        let imagePath = null
        if (commentImage[post.id]) {
          const res = await api.uploadImage(commentImage[post.id], 'comment')
          imagePath = res.image_url
        }

        const res = await fetch(getApiUrl('/comments'), {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({ 
            post_id: post.id, 
            content: content || '',
            image: imagePath
          })
        })

        if (!res.ok) {
          console.error('Failed to post comment')
          return
        }
        const comment = await res.json()
        if (!post.comments) post.comments = []
        post.comments.push(comment)
        newComment[post.id] = ''
        clearCommentImage(post.id)
      } catch (err) {
        console.error('Error posting comment:', err)
      }
    }

    const openImageModal = (image) => {
      modalImage.value = image
    }

    const formatDate = (date) => {
      if (!date) return 'just now'
      const d = new Date(date)
      const now = new Date()
      const diff = Math.floor((now - d) / 1000)
      if (diff < 60) return 'just now'
      if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
      if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
      if (diff < 604800) return `${Math.floor(diff / 86400)}d ago`
      return d.toLocaleDateString()
    }

    return { 
      postStore, 
      showComments, 
      newComment, 
      commentImage,
      commentImagePreview,
      modalImage,
      toggleComments, 
      submitComment, 
      handleCommentImageChange,
      clearCommentImage,
      openImageModal,
      formatDate,
      getImageUrl,
      getApiUrl
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

.feed-view {
  max-width: 900px;
  margin: 0 auto;
  padding: 32px 16px;
  font-family: 'Inter', 'Segoe UI', system-ui, -apple-system, sans-serif;
}

.feed-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.feed-header h1 {
  margin: 0;
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.02em;
}

.btn-primary {
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

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(246, 189, 96, 0.35);
  background: #f7c570;
}

.empty-state {
  text-align: center;
  padding: 64px 32px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.empty-icon {
  color: rgba(13, 19, 33, 0.3);
  margin-bottom: 8px;
}

.empty-state p {
  color: rgba(13, 19, 33, 0.6);
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding-bottom: 48px;
}

.post-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  transition: box-shadow 0.2s;
}

.post-card:hover {
  box-shadow: 0 4px 24px rgba(13, 19, 33, 0.13);
}

.post-header {
  padding: 24px 32px 16px 32px;
}

.user-link {
  display: flex;
  align-items: center;
  gap: 16px;
  text-decoration: none;
  color: inherit;
  transition: opacity 0.2s;
}

.user-link:hover {
  opacity: 0.8;
}

.user-avatar {
  width: 56px;
  height: 56px;
  background: var(--lavender-mist);
  color: var(--ink-black);
  font-weight: 700;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  font-size: 1.25rem;
  flex-shrink: 0;
  border: 2px solid rgba(13, 19, 33, 0.08);
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 1.25rem;
  color: var(--ink-black);
  user-select: none;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.user-info h4 {
  margin: 0;
  font-size: 1rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: 0.01em;
}

.user-info .post-time {
  font-size: 0.8125rem;
  color: rgba(13, 19, 33, 0.5);
  font-weight: 500;
}

.post-content {
  margin: 0;
  padding: 0 32px 20px 32px;
  line-height: 1.6;
  color: var(--ink-black);
  font-size: 1rem;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.post-media {
  padding: 0 32px;
  margin-bottom: 20px;
}

.post-image {
  width: 100%;
  max-height: 500px;
  object-fit: cover;
  border-radius: 12px;
  cursor: pointer;
  transition: opacity 0.2s;
  display: block;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.1);
}

.post-image:hover {
  opacity: 0.95;
}

.post-actions {
  padding: 16px 32px;
  border-top: 1px solid rgba(13, 19, 33, 0.08);
}

.btn-action {
  display: flex;
  align-items: center;
  gap: 10px;
  background: transparent;
  border: none;
  color: rgba(13, 19, 33, 0.7);
  font-weight: 600;
  cursor: pointer;
  padding: 10px 16px;
  border-radius: 10px;
  font-size: 0.9375rem;
  transition: all 0.2s;
  letter-spacing: 0.01em;
}

.btn-action:hover {
  background: var(--lavender-mist);
  color: var(--muted-teal);
}

.btn-action svg {
  flex-shrink: 0;
}

.comments-section {
  padding: 16px 32px 24px 32px;
  border-top: 1px solid rgba(13, 19, 33, 0.08);
  background: rgba(246, 240, 249, 0.3);
}

.comment-form {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 16px;
}

.comment-input {
  flex: 1;
  padding: 12px 20px;
  border-radius: 24px;
  border: 2px solid rgba(13, 19, 33, 0.08);
  outline: none;
  font-size: 0.9375rem;
  background: white;
  transition: border-color 0.2s, box-shadow 0.2s;
  color: var(--ink-black);
}

.comment-input:focus {
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.1);
}

.btn-upload {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: white;
  cursor: pointer;
  transition: all 0.2s;
  color: var(--muted-teal);
  border: 2px solid rgba(13, 19, 33, 0.08);
}

.btn-upload:hover {
  background: var(--lavender-mist);
  border-color: var(--muted-teal);
  transform: scale(1.05);
}

.btn-submit-comment {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background: var(--honey-bronze);
  color: var(--ink-black);
  border: none;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.875rem;
}

.btn-submit-comment:hover {
  background: #f7c570;
  transform: scale(1.05);
}

.image-preview {
  position: relative;
  display: inline-block;
  margin-top: 12px;
}

.image-preview img {
  max-width: 160px;
  max-height: 120px;
  border-radius: 12px;
  object-fit: cover;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.12);
  border: 2px solid white;
}

.btn-clear-image {
  position: absolute;
  top: -10px;
  right: -10px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #d32f2f;
  color: #fff;
  border: 2px solid white;
  cursor: pointer;
  font-size: 1.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  font-weight: 700;
  line-height: 1;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.btn-clear-image:hover {
  background: #b71c1c;
  transform: scale(1.1);
}

.image-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(13, 19, 33, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  cursor: pointer;
  backdrop-filter: blur(4px);
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  position: relative;
  max-width: 90vw;
  max-height: 90vh;
  cursor: default;
}

.modal-close {
  position: absolute;
  top: -56px;
  right: 0;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: white;
  color: var(--ink-black);
  border: none;
  cursor: pointer;
  font-size: 1.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  font-weight: 700;
  line-height: 1;
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
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

@media (max-width: 900px) {
  .feed-view {
    padding: 24px 16px;
  }

  .feed-header {
    margin-bottom: 24px;
  }

  .posts-list {
    gap: 20px;
  }
}

@media (max-width: 600px) {
  .feed-view {
    padding: 16px 12px;
  }

  .feed-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .btn-primary {
    justify-content: center;
  }

  .posts-list {
    gap: 16px;
    padding-bottom: 32px;
  }

  .post-card {
    border-radius: 12px;
  }

  .post-header {
    padding: 16px 20px 12px 20px;
  }

  .user-avatar {
    width: 48px;
    height: 48px;
    font-size: 1.125rem;
  }

  .post-content {
    padding: 0 20px 16px 20px;
    font-size: 0.9375rem;
  }

  .post-media {
    padding: 0 20px;
    margin-bottom: 16px;
  }

  .post-image {
    max-height: 350px;
  }

  .post-actions {
    padding: 12px 20px;
  }

  .comments-section {
    padding: 12px 20px 16px 20px;
  }

  .modal-close {
    top: -48px;
    width: 40px;
    height: 40px;
    font-size: 1.5rem;
  }
}
</style>