import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const useEventStore = defineStore('events', () => {
  const events = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  // const setEvents = (newEvents) => {
  //   events.value = newEvents
  // }

  const setLoading = (loading) => {
    isLoading.value = loading
  }

  const setError = (err) => {
    error.value = err
  }

  const clearError = () => {
    error.value = null
  }

  const fetchEvents = async (groupId) => {
    setLoading(true)
    clearError()
    try {
      const response = await api.getEvents(groupId)
      setEvents(response.events || response || [])
    } catch (err) {
      setError(err.message || 'Failed to fetch events')
    } finally {
      setLoading(false)
    }
  }

  const createEvent = async (eventData) => {
    setLoading(true)
    clearError()
    try {
     await api.createEvent(eventData)
      // if (response.event) {
      //   events.value.push(response.event)
      // }
      // return response
      return
    } catch (err) {
      setError(err.message || 'Failed to create event')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const submitResponse = async (responseData) => {
    clearError()
    try {
      const response = await api.submitEventResponse(responseData)
      return response
    } catch (err) {
      setError(err.message || 'Failed to submit response')
      throw err
    }
  }

  return {
    events,
    isLoading,
    error,
    setEvents,
    setLoading,
    setError,
    clearError,
    fetchEvents,
    createEvent,
    submitResponse
  }
})
