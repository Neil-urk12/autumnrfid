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
  },
  placeholderText: {
    type: String,
    default: 'Enter ID to Confirm Deletion'
  },
  showStatusSelection: {
    type: Boolean,
    default: false
  },
  statusOptions: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['close', 'confirm'])
const confirmationId = ref('')
const errorMessage = ref('')
const selectedStatus = ref('')

const handleClose = () => {
  confirmationId.value = ''
  errorMessage.value = ''
  emit('close')
}

const handleConfirm = () => {
  const itemId = props.itemInfo?.id || props.itemInfo?.ecode || props.itemInfo?.name
  if (!itemId) return

  if (props.showStatusSelection && !selectedStatus.value) {
    errorMessage.value = 'Please select a status'
    return
  }

  if (confirmationId.value.toLowerCase() === itemId.toLowerCase()) {
    emit('confirm', { ...props.itemInfo, newStatus: selectedStatus.value })
    confirmationId.value = ''
    errorMessage.value = ''
    selectedStatus.value = ''
    handleClose()
  } else {
    errorMessage.value = 'input does not match'
  }
}

const handleModalClick = (e) => {
  if (e.target.classList.contains('modal')) {
    handleClose()
  }
}
</script>

<template>
  <div class="modal" :class="{ active: isOpen }" @click="handleModalClick">
    <div class="modal-content confirmation-content" @click.stop>
      <div class="confirmation-icon">
        <i class="fa-solid fa-triangle-exclamation"></i>
      </div>
      <div class="confirmation-message">
        <h3>{{ title }}</h3>
        <p>Are you sure you want to delete <span>{{ itemName }}</span>?</p>
        <p class="warning-text">This action cannot be undone.</p>
        
        <div v-if="showStatusSelection" class="status-selection">
          <label>Select student's new status:</label>
          <select v-model="selectedStatus" required>
            <option value="" disabled>Select status</option>
            <option v-for="status in statusOptions" :key="status" :value="status">
              {{ status }}
            </option>
          </select>
        </div>
        
        <div class="confirmation-input">
          <input 
            type="text" 
            v-model="confirmationId"
            :placeholder="placeholderText"
            :class="{ 'error': errorMessage }"
          >
          <p class="error-message" v-if="errorMessage">{{ errorMessage }}</p>
        </div>

       
      </div>
      <div class="confirmation-actions">
        <button class="cancel-btn" @click="handleClose">
          Cancel
        </button>
        <button class="delete-confirm-btn" @click="handleConfirm">
          Delete
        </button>
      </div>
    </div>
  </div>
</template>