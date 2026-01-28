import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const usePostStore = defineStore('posts', () => {
  const posts = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  const setPosts = (newPosts) => {
    posts.value = newPosts
  }

  const setLoading = (loading) => {
    isLoading.value = loading
  }

  const setError = (err) => {
    error.value = err
  }

  const clearError = () => {
    error.value = null
  }

  const fetchPosts = async () => {
    setLoading(true)
    clearError()
    try {
      //get feed posts from api
      const response = await api.getPosts()
      if (response.posts) {
        setPosts(response.posts)
      }
    } catch (err) {
      setError(err.message || 'Failed to fetch posts')
    } finally {
      setLoading(false)
    }
  }
  const getFollowers = async () => {
    setLoading(true)
    clearError()
    try {
      const response = await api.getFollowers()
      if (response.followers) {
        return response.followers
      }
    } catch (err) {
      setError(err.message || 'Failed to fetch followers')
      return []
    } finally {
      setLoading(false)
    }
  }
  const createPost = async (postData) => {
    setLoading(true)
    clearError()
    try {
      const response = await api.createPost(postData)
      if (response.post) {
        // Instead of just unshifting, fetch the latest posts from backend
        await fetchPosts()
      }
      return response
    } catch (err) {
      setError(err.message || 'Failed to create post')
      throw err
    } finally {
      setLoading(false)
    }
  }

  return {
    posts,
    isLoading,
    error,
    setPosts,
    setLoading,
    setError,
    clearError,
    fetchPosts,
    getFollowers,
    createPost
  }
})
