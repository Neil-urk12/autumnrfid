<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  message: {
    type: String,
    required: true
  },
  duration: {
    type: Number,
    default: 3000
  },
  type: {
    type: String,
    default: 'error'
  }
})

const visible = ref(false)
const timeoutId = ref(null)

onMounted(() => {
  visible.value = true
  if (props.duration !== null) {
    timeoutId.value = setTimeout(() => {
      visible.value = false
    }, props.duration)
  }
})

onUnmounted(() => {
  if (timeoutId.value) {
    clearTimeout(timeoutId.value)
  }
})
</script>

<template>
  <Transition name="toast">
    <div v-if="visible" :class="['toast', type]">
      <div class="toast-content">
        <i :class="['toast-icon', type === 'error' ? 'fas fa-exclamation-circle' : 'fas fa-check-circle']"></i>
        <span class="toast-message">{{ message }}</span>
      </div>
    </div>
  </Transition>
</template>
