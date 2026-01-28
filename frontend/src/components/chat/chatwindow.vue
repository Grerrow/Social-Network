<template>
  <div class="chat-window" :class="{ minimized: isMinimized }">
    <chatheader 
      :chat-id="chatId"
      :group-id="groupId"
      :is-minimized="isMinimized"
      @close="$emit('close')" 
      @toggle-minimize="isMinimized = !isMinimized"
    />
    <template v-if="!isMinimized">
      <chatmessage :chat-id="chatId" :group-id="groupId" />
      <chatinput :chat-id="chatId" :group-id="groupId" />
    </template>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import chatheader from './chatheader.vue'
import chatmessage from './chatmessage.vue'
import chatinput from './chatinput.vue'

defineProps({ 
  chatId: Number,
  groupId: Number 
})
defineEmits(['close'])

const isMinimized = ref(false)
</script>

<style scoped>
:root {
  --lavender-mist: #f6f0f9;
  --ink-black: #0d1321;
  --honey-bronze: #f6bd60;
  --muted-teal: #92bfb1;
}

.chat-window {
  background: white;
  border-radius: 16px 16px 0 0;
  box-shadow: 0 4px 16px rgba(13, 19, 33, 0.15);
  display: flex;
  flex-direction: column;
  width: 328px;
  height: 450px;
  overflow: hidden;
  border: 2px solid rgba(13, 19, 33, 0.08);
  border-bottom: none;
  transition: height 0.3s ease, box-shadow 0.2s ease;
  animation: slideUp 0.3s ease;
}

.chat-window.minimized {
  height: 48px;
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.chat-window:hover {
  box-shadow: 0 6px 20px rgba(13, 19, 33, 0.2);
}

@media (max-width: 768px) {
  .chat-window {
    width: 300px;
    height: 420px;
  }
  
  .chat-window.minimized {
    height: 48px;
  }
}

@media (max-width: 480px) {
  .chat-window {
    width: 280px;
    height: 380px;
  }
  
  .chat-window.minimized {
    height: 44px;
  }
}
</style>