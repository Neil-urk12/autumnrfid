<script setup>
defineProps({
  isOpen: {
    type: Boolean,
    required: true
  }
})

const emit = defineEmits(['close', 'confirm'])

const handleDiscard = () => {
    emit('confirm')
}

const handleContinue = () => {
    emit('close')
}
</script>

<template>
  <div class="modal" :class="{ active: isOpen }">
    <div class="modal-content confirmation-content" @click.stop>
      <div class="confirmation-icon">
        <i class="fa-solid fa-triangle-exclamation"></i>
      </div>
      <div class="confirmation-message">
        <h3>Unsaved Changes</h3>
        <p>You have unsaved changes in your form</p>
        <p class="warning-text">Do you want to save your changes before leaving?</p>
      </div>
      <div class="confirmation-actions">
        <button class="cancel-btn" @click="handleContinue">Continue Editing</button>
        <button class="delete-confirm-btn" @click="handleDiscard">Discard Changes</button>
      </div>
    </div>
  </div>
</template>
