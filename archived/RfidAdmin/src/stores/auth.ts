import { defineStore } from 'pinia'
import { ref } from 'vue'

interface User {
  email: string
  token?: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isAuthenticated = ref(false)
  const error = ref('')
  
  const initializeAuth = () => {
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      user.value = JSON.parse(storedUser)
      isAuthenticated.value = true
    }
  }
  
  initializeAuth()
  
  const login = async (email: string, password: string) => {
    try {
      error.value = ''
      
      if (!email || !password) {
        throw new Error('Email and password are required')
      }
      
      if (email === 'admin@example.com' && password === 'password') {
        const loggedInUser = {
          email,
          token: 'simulated-jwt-token'
        }
        
        user.value = loggedInUser
        isAuthenticated.value = true
        
        localStorage.setItem('user', JSON.stringify(loggedInUser))
        
        return true
      } else {
        throw new Error('Invalid email or password')
      }
    } catch (err: any) {
      error.value = err.message || 'Login failed'
      return false
    }
  }
  
  const logout = () => {
    user.value = null
    isAuthenticated.value = false
    localStorage.removeItem('user')
  }
  
  return {
    user,
    isAuthenticated,
    error,
    login,
    logout
  }
})