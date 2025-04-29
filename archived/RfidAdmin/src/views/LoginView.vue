<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const isLoading = ref(false)

const handleLogin = async () => {
  isLoading.value = true
  
  try {
    const success = await authStore.login(email.value, password.value)
    
    if (success) {
      router.push('/')
    }
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-form">
      <h1>RFID Admin Login</h1>
      
      <div v-if="authStore.error" class="error-message">
        {{ authStore.error }}
      </div>
      
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="email">Email</label>
          <input 
            type="email" 
            id="email" 
            v-model="email" 
            placeholder="Enter your email"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">Password</label>
          <input 
            type="password" 
            id="password" 
            v-model="password" 
            placeholder="Enter your password"
            required
          />
        </div>
        
        <button type="submit" :disabled="isLoading">{{ isLoading ? 'Logging in...' : 'Login' }}</button>
      </form>
    </div>
  </div>
</template>

<style scoped>
@import '@/../style.css';

.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: var(--clr);
}

.login-form {
  width: 100%;
  max-width: 400px;
  padding: 2.5rem;
  background: var(--sidebar-gradient);
  border-radius: 15px;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
  box-shadow: 0 8px 32px var(--shadow-color, rgba(0, 0, 0, 0.3));
}

h1 {
  text-align: center;
  margin-bottom: 2rem;
  color: var(--text-primary);
  font-size: 1.8em;
}

.form-group {
  margin-bottom: 1.5rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: var(--text-secondary);
  font-size: 0.9em;
}

input {
  width: 100%;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
  background: var(--input-bg, rgba(255, 255, 255, 0.05));
  color: var(--text-primary);
  font-size: 1rem;
  transition: border-color 0.3s, box-shadow 0.3s;
}

input:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--border-glow);
}

button {
  width: 100%;
  padding: 0.85rem;
  background-color: var(--accent);
  color: var(--text-primary);
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover:not(:disabled) {
  background-color: var(--hover-color);
}

button:disabled {
  background-color: var(--accent-disabled, rgba(79, 70, 229, 0.5));
  cursor: not-allowed;
  opacity: 0.7;
}

.error-message {
  background-color: var(--error-bg, #f8d7da); 
  color: var(--error-text, #721c24); 
  padding: 0.75rem 1.25rem;
  margin-bottom: 1rem;
  border: 1px solid var(--error-border, #f5c6cb); 
  border-radius: 0.25rem;
  font-size: 0.9em;
}
</style>