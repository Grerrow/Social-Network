<!-- filepath: /home/kali/Desktop/social-network/frontend/src/views/CreatePostView.vue -->
<template>
  <div class="create-post-view">
    <div class="create-post-card">
      <h1>Create a Post</h1>
      
      <form @submit.prevent="handleCreatePost">
        <div class="form-group">
          <label for="content">What's on your mind?</label>
          <textarea
            id="content"
            v-model="form.content"
            placeholder="Share your thoughts..."
            rows="5"
            maxlength="400"
          ></textarea>
          <span class="char-count">{{ form.content.length }}/400</span>
        </div>

        <div class="form-group">
          <label for="privacy">Privacy</label>
          <select id="privacy" v-model="form.privacy" required>
            <option disabled value="">Select privacy</option>
            <option value="public">Public</option>
            <option value="private">Almost Private</option>
            <option value="almost_private">Private</option>
          </select>
        </div>

        <div v-if="form.privacy === 'almost_private'" class="form-group">
          <label for="allowedFollowers">Select Followers to Share With</label>
          <select id="allowedFollowers" multiple v-model="form.allowedFollowers">
            <option v-for="follower in followersList" :key="follower.id" :value="follower.id">
              {{ follower.username || follower.first_name + ' ' + follower.last_name }}
            </option>
          </select>
        </div>

        <!-- Professional Image Upload -->
        <div class="form-group">
          <label>Attach Image (Optional)</label>
          <div class="image-upload-container">
            <input 
              type="file" 
              ref="fileInput"
              @change="handleImageChange" 
              accept=".jpg,.jpeg,.png,.gif" 
              hidden
            />
            
            <div v-if="!form.imagePreview" class="image-upload-empty">
              <div class="upload-area" @click="triggerFileInput">
                <div class="upload-icon-wrapper">
                  <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                    <circle cx="8.5" cy="8.5" r="1.5"/>
                    <polyline points="21 15 16 10 5 21"/>
                  </svg>
                </div>
                <div class="upload-text">
                  <p class="upload-title">Click to upload image</p>
                  <p class="upload-subtitle">JPG, PNG, GIF up to 5MB</p>
                </div>
              </div>
            </div>
            
            <div v-else class="image-upload-preview">
              <div class="preview-wrapper">
                <img :src="form.imagePreview" alt="preview" />
                <div class="preview-overlay">
                  <button type="button" class="btn-preview-change" @click="triggerFileInput">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                    Change
                  </button>
                  <button type="button" class="btn-preview-remove" @click="clearImage">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="3 6 5 6 21 6"/>
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                    </svg>
                    Remove
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="form-actions">
          <router-link to="/feed" class="btn-cancel">Cancel</router-link>
          <button
            type="submit"
            :disabled="postStore.isLoading || (!form.content.trim() && !form.image)"
            class="btn-submit"
          >
            {{ postStore.isLoading ? 'Posting...' : 'Post' }}
          </button>
        </div>

        <p v-if="postStore.error" class="error-message">{{ postStore.error }}</p>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { usePostStore } from '../stores/posts'
import api from '../services/api'

export default {
  setup() {
    const router = useRouter()
    const postStore = usePostStore()
    const fileInput = ref(null)
    const followersList = ref([])

    const form = ref({
      content: '',
      privacy: '',
      image: null,
      imagePreview: null,
      allowedFollowers: []
    })

    // Watch for privacy changes to fetch followers when "almost_private" is selected
    watch(() => form.value.privacy, async (newPrivacy) => {
      if (newPrivacy === 'almost_private' && followersList.value.length === 0) {
        try {
          const response = await api.getFollowers()
          followersList.value = response.followers || []
        } catch (err) {
          console.error('Failed to fetch followers', err)
        }
      }
    })

    const triggerFileInput = () => {
      fileInput.value.click()
    }

    const handleImageChange = async (e) => {
      const file = e.target.files[0]
      if (!file) return

      const validTypes = ['image/jpeg', 'image/png', 'image/gif']
      if (!validTypes.includes(file.type)) {
        alert('Only JPG, PNG, GIF allowed')
        return
      }

      if (file.size > 5 * 1024 * 1024) {
        alert('File size must be less than 5MB')
        return
      }

      form.value.image = file
      form.value.imagePreview = URL.createObjectURL(file)
    }

    const clearImage = () => {
      form.value.image = null
      form.value.imagePreview = null
      if (fileInput.value) {
        fileInput.value.value = ''
      }
    }

    const handleCreatePost = async () => {
      if (!form.value.content.trim() && !form.value.image) {
        alert('Post must have text or image')
        return
      }

      if (!form.value.privacy) {
        alert('Select privacy')
        return
      }

      let imagePath = null
      if (form.value.image) {
        try {
          const res = await api.uploadImage(form.value.image, 'post')
          imagePath = res.image_url
        } catch (err) {
          console.error('Upload failed', err)
          alert('Image upload failed')
          return
        }
      }

      try {
        const postData = {
          content: form.value.content.trim(),
          privacy: form.value.privacy,
          image: imagePath
        }
        
        // Include allowed_followers for almost_private posts
        if (form.value.privacy === 'almost_private' && form.value.allowedFollowers.length > 0) {
          postData.allowed_followers = form.value.allowedFollowers
        }

        await postStore.createPost(postData)

        form.value = { content: '', privacy: '', image: null, imagePreview: null, allowedFollowers: [] }
        router.push('/feed')
      } catch (err) {
        console.error('Post creation failed', err)
      }
    }

    return { form, postStore, fileInput, followersList, triggerFileInput, handleImageChange, clearImage, handleCreatePost }
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

.create-post-view {
  max-width: 560px;
  margin: 0 auto;
}

.create-post-card {
  background: white;
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 2px 12px rgba(13, 19, 33, 0.08);
}

.create-post-card h1 {
  margin: 0 0 24px;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--ink-black);
  letter-spacing: -0.02em;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: var(--ink-black);
  font-size: 0.875rem;
  letter-spacing: 0.01em;
}

.form-group textarea,
.form-group select {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  font-size: 0.9375rem;
  font-family: inherit;
  transition: all 0.2s ease;
  box-sizing: border-box;
  color: var(--ink-black);
  background: white;
}

.form-group textarea {
  resize: vertical;
  line-height: 1.5;
}

.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: var(--honey-bronze);
  box-shadow: 0 0 0 3px rgba(246, 189, 96, 0.15);
}

/* Multi-select followers list */
.form-group select[multiple] {
  min-height: 200px;
  padding: 8px;
  background: var(--lavender-mist);
  border: 2px solid rgba(13, 19, 33, 0.12);
}

.form-group select[multiple]:focus {
  background: white;
  border-color: var(--honey-bronze);
}

.form-group select[multiple] option {
  padding: 12px 16px;
  margin: 4px 0;
  border-radius: 8px;
  background: white;
  color: var(--ink-black);
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
  border: 1.5px solid transparent;
}

.form-group select[multiple] option:hover {
  background: linear-gradient(135deg, var(--honey-bronze) 0%, #f7c570 100%);
  color: var(--ink-black);
  border-color: rgba(13, 19, 33, 0.1);
  transform: translateX(4px);
}

.form-group select[multiple] option:checked {
  background: var(--muted-teal);
  color: white;
  font-weight: 600;
  border-color: rgba(13, 19, 33, 0.15);
}

.form-group select[multiple] option:checked::before {
  content: "✓ ";
  margin-right: 8px;
  font-weight: 700;
}

.char-count {
  display: block;
  margin-top: 8px;
  font-size: 0.8125rem;
  color: rgba(13, 19, 33, 0.6);
  text-align: right;
}

/* Professional Image Upload Container */
.image-upload-container {
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  overflow: hidden;
  background: white;
}

.image-upload-empty {
  padding: 16px;
}

.upload-area {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  border: 2px dashed rgba(13, 19, 33, 0.2);
  border-radius: 10px;
  background: var(--lavender-mist);
  cursor: pointer;
  transition: all 0.2s ease;
}

.upload-area:hover {
  border-color: var(--honey-bronze);
  background: rgba(246, 189, 96, 0.1);
}

.upload-icon-wrapper {
  flex-shrink: 0;
  width: 56px;
  height: 56px;
  border-radius: 12px;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(13, 19, 33, 0.5);
}

.upload-text {
  flex: 1;
}

.upload-title {
  margin: 0 0 4px 0;
  font-weight: 600;
  font-size: 0.9375rem;
  color: var(--ink-black);
}

.upload-subtitle {
  margin: 0;
  font-size: 0.8125rem;
  color: rgba(13, 19, 33, 0.5);
}

.image-upload-preview {
  padding: 16px;
}

.preview-wrapper {
  position: relative;
  border-radius: 10px;
  overflow: hidden;
  background: var(--lavender-mist);
}

.preview-wrapper img {
  width: 100%;
  max-height: 320px;
  object-fit: cover;
  display: block;
}

.preview-overlay {
  position: absolute;
  bottom: 16px;
  right: 16px;
  display: flex;
  gap: 10px;
}

.btn-preview-change,
.btn-preview-remove {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
  backdrop-filter: blur(8px);
  letter-spacing: 0.01em;
}

.btn-preview-change {
  background: rgba(255, 255, 255, 0.95);
  color: var(--ink-black);
  box-shadow: 0 2px 8px rgba(13, 19, 33, 0.15);
}

.btn-preview-change:hover {
  background: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(13, 19, 33, 0.2);
}

.btn-preview-remove {
  background: rgba(211, 47, 47, 0.95);
  color: white;
  box-shadow: 0 2px 8px rgba(211, 47, 47, 0.3);
}

.btn-preview-remove:hover {
  background: #d32f2f;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(211, 47, 47, 0.4);
}

.btn-preview-change svg,
.btn-preview-remove svg {
  flex-shrink: 0;
}

.form-actions {
  display: flex;
  gap: 16px;
  margin-top: 24px;
}

.btn-cancel {
  flex: 1;
  padding: 12px;
  border: 2px solid rgba(13, 19, 33, 0.12);
  border-radius: 12px;
  background: white;
  color: var(--ink-black);
  font-weight: 600;
  text-align: center;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.9375rem;
  letter-spacing: 0.01em;
}

.btn-cancel:hover {
  border-color: var(--muted-teal);
  color: var(--muted-teal);
  background: rgba(146, 191, 177, 0.05);
}

.btn-submit {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: 12px;
  background: var(--honey-bronze);
  color: var(--ink-black);
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.9375rem;
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

.error-message {
  color: #d32f2f;
  font-size: 0.875rem;
  margin-top: 16px;
  padding: 12px 16px;
  background: rgba(211, 47, 47, 0.08);
  border: 1px solid rgba(211, 47, 47, 0.2);
  border-radius: 12px;
  text-align: center;
  font-weight: 500;
}

@media (max-width: 600px) {
  .create-post-card {
    padding: 24px;
  }

  .form-actions {
    flex-direction: column;
  }

  .upload-area {
    flex-direction: column;
    text-align: center;
  }

  .preview-overlay {
    position: static;
    margin-top: 12px;
    justify-content: center;
  }
}
</style>