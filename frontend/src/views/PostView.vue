<template>
  <div class="post-view">
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>Loading post...</p>
    </div>

    <div v-else-if="error" class="error-container">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"/>
        <line x1="12" y1="8" x2="12" y2="12"/>
        <line x1="12" y1="16" x2="12.01" y2="16"/>
      </svg>
      <h2>Post not found</h2>
      <p>{{ error }}</p>
      <router-link to="/feed" class="btn-primary">Back to Feed</router-link>
    </div>

    <div v-else-if="post" class="post-container">
      <!-- Post Card -->
      <div class="post-card">
        <!-- Author Header -->
        <div class="post-header">
          <router-link :to="`/profile?user_id=${post.author?.id}`" class="author-link">
            <img 
              class="author-avatar"
              :src="post.author?.avatar ? getImageUrl(post.author.avatar) : '/images/placeholder.png'"
              :alt="post.author?.username"
            />
          </router-link>
          <div class="author-info">
            <router-link :to="`/profile?user_id=${post.author?.id}`" class="author-name">
              {{ post.author?.username || 'Unknown User' }}
            </router-link>
            <span class="post-time">{{ formatTime(post.created_at) }}</span>
          </div>
          <span v-if="post.privacy=== 'almost_private'" class="privacy-badge">
            Almost Private
          </span>
          <span v-else-if="post.privacy != 'almost_private' && post.privacy != ''" class="privacy-badge">
            {{post.privacy}}
          </span>
        </div>

        <!-- Post Content -->
        <div class="post-body">
          <p v-if="post.content" class="post-text">{{ post.content }}</p>
          
          <!-- Post Image -->
          <div v-if="post.image" class="post-image-container">
            <img 
              :src="getImageUrl(post.image)" 
              alt="Post image"
              class="post-image"
            />
          </div>
        </div>

        <!-- Post Actions -->
        <div class="post-actions">
          <div class="action-stats">
            <span class="stat-item">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/>
              </svg>
              {{ post.comments?.length || 0 }} comments
            </span>
          </div>
        </div>
      </div>

      <!-- Comments Section -->
      <div class="comments-section">
        <h2 class="comments-title">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/>
          </svg>
          Comments
        </h2>

        <!-- Add Comment Form -->
        <div class="add-comment-form">
          <img 
            class="comment-avatar"
            :src="authStore.user?.avatar ? getImageUrl(authStore.user.avatar) : '/images/placeholder.png'"
            :alt="authStore.user?.username"
          />
          <div class="comment-input-wrapper">
            <textarea
              v-model="newComment[post.id]"
              placeholder="Write a comment..."
              rows="2"
              @keydown.enter.ctrl="submitComment(post)"
            ></textarea>
            <div class="comment-actions">
              <label class="attach-image-btn">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                  <circle cx="8.5" cy="8.5" r="1.5"/>
                  <polyline points="21 15 16 10 5 21"/>
                </svg>
                <input type="file" accept=".jpg,.jpeg,.png,.gif" @change="handleImageSelect" hidden />
              </label>
              <button 
                class="btn-submit-comment" 
                @click="submitComment(post)"
                :disabled="!newComment[post.id]?.trim() && !commentImage"
              >
                Post
              </button>
            </div>
          </div>
        </div>

        <!-- Image Preview for Comment -->
        <div v-if="commentImagePreview" class="comment-image-preview">
          <img :src="commentImagePreview" alt="Preview" />
          <button class="remove-preview-btn" @click="removeCommentImage">×</button>
        </div>

        <!-- Comments List -->
        <div v-if="post.comments?.length === 0" class="no-comments">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/>
          </svg>
          <p>No comments yet. Be the first to comment!</p>
        </div>

        <div v-else class="comments-list">
          <div v-for="comment in post.comments" :key="comment.id" class="comment-item">
            <router-link :to="`/profile?user_id=${comment.author?.id}`" class="comment-avatar-link">
              <img 
                class="comment-avatar"
                :src="comment.author?.avatar ? getImageUrl(comment.author.avatar) : '/images/placeholder.png'"
                :alt="comment.author?.username"
              />
            </router-link>
            <div class="comment-content">
              <div class="comment-bubble">
                <router-link :to="`/profile?user_id=${comment.author?.id}`" class="comment-author">
                  {{ comment.author?.username || 'Unknown' }}
                </router-link>
                <p v-if="comment.content" class="comment-text">{{ comment.content }}</p>
                <img 
                  v-if="comment.image" 
                  :src="getImageUrl(comment.image)" 
                  alt="Comment image"
                  class="comment-image"
                />
              </div>
              <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useApi } from '@/composables/useApi'
import { usePostStore } from '@/stores/posts'
import api from '@/services/api'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { getImageUrl, getApiUrl } = useApi()

const post = ref(null)
const loading = ref(true)
const error = ref(null)
const newComment = ref({})
const commentImage = ref(null)
const commentImagePreview = ref(null)

const postsStore = usePostStore()

const fetchPost = async () => {
  try {
    const postId = parseInt(route.params.id)
    const url = getApiUrl(`/post?post_id=${postId}`)
    
    const res = await fetch(url, { credentials: 'include' })
    if (!res.ok) {
      error.value = 'The post you are looking for does not exist or has been removed.'
      loading.value = false
      return
    }
    
    const data = await res.json()
    
    // Transform backend response to match frontend expectations
    post.value = {
      ...data.post,
      author: {
        id: data.post.user_id,
        username: data.post.username,
        avatar: data.post.avatar
      },
      comments: data.comments.map(comment => ({
        ...comment,
        author: {
          id: comment.user_id,
          username: comment.username,
          avatar: comment.avatar
        }
      }))
    }
    
    loading.value = false
  } catch (err) {
    console.error('Error fetching post:', err)
    error.value = 'Failed to load post'
    loading.value = false
  }
}

const handleImageSelect = (e) => {
  const file = e.target.files?.[0]
  if (!file) return
  
  const validTypes = ['image/jpeg', 'image/png', 'image/gif']
      if (!validTypes.includes(file.type)) {
        alert('Only JPG, PNG, GIF allowed')
        return
      }
  if (file.size > 5 * 1024 * 1024) {
    alert('Image must be less than 5MB')
    return
  }
  
  commentImage.value = file
  commentImagePreview.value = URL.createObjectURL(file)
}

const removeCommentImage = () => {
  if (commentImagePreview.value) {
    URL.revokeObjectURL(commentImagePreview.value)
  }
  commentImage.value = null
  commentImagePreview.value = null
}

const submitComment = async (postData) => {
  if (!newComment.value[postData.id]?.trim() && !commentImage.value) {
    return
  }
  
  try {
    let imagePath = null
    
    // Upload image first if present
    if (commentImage.value) {
      const uploadRes = await api.uploadImage(commentImage.value, 'comment')
      imagePath = uploadRes.image_url
    }

    // Send comment data
    const commentRes = await api.createComment({
      post_id: postData.id,
      content: newComment.value[postData.id] || '',
      image: imagePath || ''
    })
    
    // Handle different response structures
    const commentData = commentRes.comment || commentRes
    
    // Add new comment to the list with author structure
    const newCommentData = {
      ...commentData,
      image: imagePath || commentData.image || '',
      author: {
        id: commentData.user_id,
        username: commentData.username || 'Unknown',
        avatar: commentData.avatar
      }
    }
    
    post.value.comments.push(newCommentData)
    
    // Clear form
    newComment.value[postData.id] = ''
    removeCommentImage()
  } catch (err) {
    console.error('Error submitting comment:', err)
    alert('Failed to post comment')
  }
}

const formatTime = (timestamp) => {
  if (!timestamp) return 'just now'
  
  let date
  if (typeof timestamp === 'number') {
    // Unix timestamp
    date = new Date(timestamp * 1000)
  } else if (typeof timestamp === 'string') {
    // SQLite format: "2024-01-28 12:34:56" - treat as UTC
    if (timestamp.includes(' ')) {
      date = new Date(timestamp.replace(' ', 'T') + 'Z')
    } else {
      date = new Date(timestamp)
    }
  } else {
    return 'just now'
  }
  
  // Check if date is valid
  if (isNaN(date.getTime())) {
    console.error('Invalid timestamp:', timestamp)
    return 'just now'
  }
  
  const now = new Date()
  const diff = Math.floor((now - date) / 1000)
  
  if (diff < 0) return 'just now' // Future dates
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

onMounted(() => {
  fetchPost()
})
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.post-view {
  max-width: 800px;
  margin: 0 auto;
  padding: 24px 16px;
}

.loading-container,
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  gap: 16px;
  text-align: center;
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

.error-container svg {
  color: rgba(13, 19, 33, 0.3);
}

.error-container h2 {
  margin: 0;
  color: var(--ink-black);
  font-size: 1.5rem;
}

.error-container p {
  margin: 0;
  color: rgba(13, 19, 33, 0.6);
}

.post-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* Post Card */
.post-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  overflow: hidden;
}

.post-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 24px;
  border-bottom: 1px solid rgba(13, 19, 33, 0.08);
}

.author-link {
  text-decoration: none;
  flex-shrink: 0;
}

.author-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--lavender-mist);
}

.author-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.author-name {
  font-weight: 700;
  font-size: 1.125rem;
  color: var(--ink-black);
  text-decoration: none;
  letter-spacing: -0.01em;
}

.author-name:hover {
  color: var(--muted-teal);
}

.post-time {
  font-size: 0.875rem;
  color: rgba(13, 19, 33, 0.6);
}

.privacy-badge {
  padding: 6px 12px;
  background: var(--lavender-mist);
  border-radius: 8px;
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--ink-black);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.post-body {
  padding: 24px;
}

.post-text {
  margin: 0 0 20px 0;
  font-size: 1.125rem;
  line-height: 1.7;
  color: var(--ink-black);
  white-space: pre-wrap;
  word-wrap: break-word;
}

.post-image-container {
  border-radius: 12px;
  overflow: hidden;
  background: var(--lavender-mist);
}

.post-image {
  width: 100%;
  max-height: 600px;
  object-fit: contain;
  display: block;
}

.post-actions {
  padding: 16px 24px;
  border-top: 1px solid rgba(13, 19, 33, 0.08);
}

.action-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9375rem;
  font-weight: 600;
  color: rgba(13, 19, 33, 0.7);
}

.stat-item svg {
  color: rgba(13, 19, 33, 0.5);
}

/* Comments Section */
.comments-section {
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  padding: 24px;
}

.comments-title {
  margin: 0 0 24px 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--ink-black);
  display: flex;
  align-items: center;
  gap: 12px;
}

.comments-title svg {
  color: var(--muted-teal);
}

.add-comment-form {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 2px solid rgba(13, 19, 33, 0.08);
}

.comment-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.comment-avatar-link {
  text-decoration: none;
  flex-shrink: 0;
}

.comment-input-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.comment-input-wrapper textarea {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.9375rem;
  font-family: inherit;
  resize: vertical;
  outline: none;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.comment-input-wrapper textarea:focus {
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

.comment-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.attach-image-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: var(--lavender-mist);
  cursor: pointer;
  transition: all 0.2s ease;
}

.attach-image-btn:hover {
  background: rgba(246, 240, 249, 0.7);
}

.attach-image-btn svg {
  color: rgba(13, 19, 33, 0.6);
}

.btn-submit-comment {
  padding: 8px 20px;
  background: var(--honey-bronze);
  color: var(--ink-black);
  border: none;
  border-radius: 8px;
  font-weight: 700;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-submit-comment:hover:not(:disabled) {
  background: #f7c570;
  transform: translateY(-1px);
}

.btn-submit-comment:disabled {
  background: rgba(13, 19, 33, 0.1);
  color: rgba(13, 19, 33, 0.4);
  cursor: not-allowed;
}

.comment-image-preview {
  position: relative;
  margin-bottom: 16px;
  border-radius: 12px;
  overflow: hidden;
  max-width: 300px;
}

.comment-image-preview img {
  width: 100%;
  display: block;
}

.remove-preview-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--ink-black);
  color: white;
  border: none;
  cursor: pointer;
  font-size: 1.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.remove-preview-btn:hover {
  transform: scale(1.1);
}

.no-comments {
  text-align: center;
  padding: 48px 24px;
  color: rgba(13, 19, 33, 0.5);
}

.no-comments svg {
  margin-bottom: 16px;
}

.no-comments p {
  margin: 0;
  font-size: 0.9375rem;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.comment-item {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.comment-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.comment-bubble {
  background: var(--lavender-mist);
  border-radius: 12px;
  padding: 12px 16px;
}

.comment-author {
  font-weight: 700;
  font-size: 0.875rem;
  color: var(--ink-black);
  text-decoration: none;
  margin-bottom: 4px;
  display: block;
}

.comment-author:hover {
  color: var(--muted-teal);
}

.comment-text {
  margin: 0;
  font-size: 0.9375rem;
  line-height: 1.5;
  color: var(--ink-black);
  white-space: pre-wrap;
  word-wrap: break-word;
}

.comment-image {
  margin-top: 8px;
  max-width: 100%;
  max-height: 300px;
  border-radius: 8px;
  display: block;
}

.comment-time {
  font-size: 0.75rem;
  color: rgba(13, 19, 33, 0.5);
  margin-left: 16px;
}

.btn-primary {
  padding: 12px 24px;
  background: var(--honey-bronze);
  color: var(--ink-black);
  border: none;
  border-radius: 12px;
  font-size: 0.9375rem;
  font-weight: 700;
  text-decoration: none;
  display: inline-block;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-primary:hover {
  background: #f7c570;
  transform: translateY(-2px);
}

@media (max-width: 768px) {
  .post-view {
    padding: 16px 12px;
  }

  .post-header {
    padding: 16px;
  }

  .author-avatar {
    width: 48px;
    height: 48px;
  }

  .author-name {
    font-size: 1rem;
  }

  .post-body {
    padding: 16px;
  }

  .post-text {
    font-size: 1rem;
  }

  .comments-section {
    padding: 16px;
  }

  .add-comment-form {
    flex-direction: column;
  }

  .comment-avatar {
    width: 36px;
    height: 36px;
  }
}
</style>
