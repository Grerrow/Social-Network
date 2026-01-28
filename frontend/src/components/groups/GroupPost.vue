<!-- filepath: /home/kali/Desktop/social-network/frontend/src/components/groups/GroupPost.vue -->
<template>
  <div class="group-post">
    <div class="post-header">
      <div class="author-avatar">
        <img 
          :src="post.author?.avatar ? getImageUrl(post.author.avatar) : '/images/placeholder.png'" 
          :alt="post.author?.username" 
        />
      </div>
      <div class="post-meta">
        <span class="author-name">{{ post.author?.username || 'Unknown' }}</span>
        <span class="post-time">{{ formatTime(post.created_at) }}</span>
      </div>
    </div>
    
    <div class="post-content">
      <p v-if="post.content">{{ post.content }}</p>
    </div>

    <!-- Post Image (outside content, like FeedView) -->
    <div v-if="post.image" class="post-media">
      <img 
        :src="getImageUrl(post.image)" 
        alt="Post image" 
        class="post-image"
        @click="showPostImageModal = true"
        loading="lazy"
      />
    </div>
    
    <div class="post-actions">
      <button class="btn-action" @click="toggleComments">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
        </svg>
        <span>{{ showComments ? 'Hide' : 'Show' }} Comments</span>
      </button>
    </div>
    
    <!-- Comments Section -->
    <div v-if="showComments" class="comments-section">
      <!-- Comments List FIRST (like FeedView) -->
      <GroupComment 
        v-for="comment in comments" 
        :key="comment.id" 
        :comment="comment" 
      />

      <!-- Comment Form AFTER comments (like FeedView) -->
      <div class="comment-form">
        <input 
          v-model="newComment" 
          @keyup.enter="submitComment"
          placeholder="Write a comment..."
          class="comment-input"
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
          accept="image/*"
          @change="handleImageSelect"
          hidden
        />
        
        <button class="btn-submit-comment" @click="submitComment">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
            <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"/>
          </svg>
        </button>
      </div>

      <!-- Image Preview AFTER form (like FeedView) -->
      <div v-if="commentImagePreview" class="image-preview">
        <img :src="commentImagePreview" alt="Preview" />
        <button class="btn-clear-image" @click="clearImage">×</button>
      </div>
    </div>

    <!-- Post Image Modal -->
    <div v-if="showPostImageModal" class="image-modal" @click="showPostImageModal = false">
      <div class="modal-content" @click.stop>
        <button class="modal-close" @click="showPostImageModal = false">×</button>
        <img :src="getImageUrl(post.image)" alt="Full size image" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useGroupsStore } from '@/stores/groups'
import GroupComment from './GroupComment.vue'
import api from '@/services/api'
import { useApi } from '@/composables/useApi'


const props = defineProps({
  post: {
    type: Object,
    required: true
  },
  groupId: {
    type: [Number, String],
    required: true
  }
})

const { getImageUrl } = useApi()
const groupsStore = useGroupsStore()
const showComments = ref(false)
const comments = ref([])
const newComment = ref('')
const commentImage = ref(null)
const commentImagePreview = ref(null)
const showPostImageModal = ref(false)

const formatTime = (timestamp) => {
  if (!timestamp) return ''
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

const handleImageSelect = (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  if (!file.type.startsWith('image/')) {
    alert('Please select an image file')
    return
  }
  
  if (file.size > 5 * 1024 * 1024) {
    alert('Image size must be less than 5MB')
    return
  }
  
  commentImage.value = file
  commentImagePreview.value = URL.createObjectURL(file)
}

const clearImage = () => {
  commentImage.value = null
  commentImagePreview.value = null
  const input = document.getElementById(`comment-image-${props.post.id}`)
  if (input) input.value = ''
}

const toggleComments = async () => {
  showComments.value = !showComments.value
  if (showComments.value && comments.value.length === 0) {
    comments.value = await groupsStore.fetchPostComments(props.groupId, props.post.id) || []
  }
}

const submitComment = async () => {
  if (!newComment.value.trim() && !commentImage.value) return
  
  try {
    let imagePath = null
    
    // Upload image first if present
    if (commentImage.value) {
      const uploadRes = await api.uploadImage(commentImage.value, 'comment')
      imagePath = uploadRes.image_url
    }
    
    // Use store method with image support
    const comment = await groupsStore.createPostComment(
      props.groupId, 
      props.post.id, 
      newComment.value.trim(),
      imagePath  // Pass image path as 4th parameter
    )
    
    if (comment) {
      comments.value.push(comment)
      newComment.value = ''
      clearImage()
    }
  } catch (error) {
    console.error('Failed to submit comment:', error)
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

.group-post {
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(13, 19, 33, 0.08);
  margin-bottom: 24px;
  overflow: hidden;
  transition: box-shadow 0.2s;
}

.group-post:hover {
  box-shadow: 0 4px 24px rgba(13, 19, 33, 0.13);
}

.post-header {
  padding: 24px 32px 16px 32px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.author-avatar {
  flex-shrink: 0;
}

.author-avatar img,
.author-avatar .avatar-placeholder {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  object-fit: cover;
  display: block;
}

.avatar-placeholder {
  background: var(--lavender-mist);
  color: var(--ink-black);
  font-weight: 700;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 1.25rem;
  border: 2px solid rgba(13, 19, 33, 0.08);
  user-select: none;
}

.post-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.author-name {
  font-size: 1rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: 0.01em;
}

.post-time {
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
}

.post-content p {
  margin: 0;
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
  .post-header {
    padding: 20px 24px 16px;
  }

  .post-content {
    padding: 0 24px 16px;
  }

  .post-media {
    padding: 0 24px;
    margin-bottom: 16px;
  }

  .post-actions,
  .comments-section {
    padding: 12px 24px 16px;
  }
}

@media (max-width: 600px) {
  .post-header {
    padding: 16px 20px 12px;
  }

  .author-avatar img,
  .author-avatar .avatar-placeholder {
    width: 48px;
    height: 48px;
    font-size: 1.125rem;
  }

  .post-content {
    padding: 0 20px 16px;
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
    padding: 12px 20px 16px;
  }

  .modal-close {
    top: -48px;
    width: 40px;
    height: 40px;
    font-size: 1.5rem;
  }
}
</style>