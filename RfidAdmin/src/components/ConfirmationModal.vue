<script setup>
import { ref, defineProps, defineEmits } from 'vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true
  },
  title: {
    type: String,
    default: 'Delete Confirmation'
  },
  itemName: {
    type: String,
    required: true
  },
  itemInfo: {
    type: [String, Object],
    default: null
  }
})

const emit = defineEmits(['close', 'confirm'])

const handleClose = () => {
  emit('close')
}

const handleConfirm = () => {
  emit('confirm', props.itemInfo)
  handleClose()
}

const handleModalClick = (e) => {
  if (e.target.classList.contains('modal')) {
    handleClose()
  }
}
</script>

<template>
  <div 
    class="modal" 
    :class="{ active: isOpen }"
    @click="handleModalClick"
  >
    <div class="modal-content confirmation-content" @click.stop>
      <div class="confirmation-icon">
        <i class="fa-solid fa-triangle-exclamation"></i>
      </div>
      <div class="confirmation-message">
        <h3>{{ title }}</h3>
        <p>Are you sure you want to delete <span>{{ itemName }}</span>?</p>
        <p class="warning-text">This action cannot be undone.</p>
      </div>
      <div class="confirmation-actions">
        <button 
          class="cancel-btn" 
          @click="handleClose"
        >
          Cancel
        </button>
        <button 
          class="delete-confirm-btn" 
          @click="handleConfirm"
        >
          Delete
        </button>
      </div>
    </div>
  </div>
</template>