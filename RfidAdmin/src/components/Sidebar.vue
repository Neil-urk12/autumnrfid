<script setup>
import { ref, computed } from "vue"
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const currentPath = computed(() => route.path)

const isActive = (path) => {
  return currentPath.value === path
}

const goTo = (path) => {
  router.push(path)
}
</script>

<template>
  <ul>
    <div class="menulist">
      <div class="sidebar-section">
        <div class="sidebar-label">Core</div>
        <li :class="{ active: isActive('/') }">
          <a @click="goTo('/')">
            <div class="icon">
              <i class="fa-solid fa-gauge"></i>
            </div>
            <div class="text">Dashboard</div>
          </a>
        </li>
      </div>

      <div class="sidebar-section">
        <div class="sidebar-label">Management</div>
        <li :class="{ active: isActive('/manage') }">
          <a @click="goTo('/manage')">
            <div class="icon">
              <i class="fa-solid fa-book"></i>
            </div>
            <div class="text">Manage</div>
          </a>
        </li>
        <li :class="{ active: isActive('/students') }">
          <a @click="goTo('/students')">
            <div class="icon">
              <i class="fa-solid fa-user-graduate"></i>
            </div>
            <div class="text">Students</div>
          </a>
        </li>
        <li :class="{ active: isActive('/courses') }">
          <a @click="goTo('/courses')">
            <div class="icon">
              <i class="fa-solid fa-file-invoice-dollar"></i>
            </div>
            <div class="text">Course</div>
          </a>
        </li>
        <li :class="{ active: isActive('/grades') }">
          <a @click="goTo('/grades')">
            <div class="icon">
              <i class="fa-solid fa-star"></i>
            </div>
            <div class="text">Grades</div>
          </a>
        </li>
        <li :class="{ active: isActive('/bills') }">
          <a @click="goTo('/bills')">
            <div class="icon">
              <i class="fa-solid fa-file-invoice-dollar"></i>
            </div>
            <div class="text">Bills</div>
          </a>
        </li>
      </div>

      <div class="sidebar-section">
        <div class="sidebar-label">Settings</div>
        <li :class="{ active: isActive('/admin') }">
          <a @click="goTo('/admin')">
            <div class="icon">
              <i class="fa-solid fa-user-shield"></i>
            </div>
            <div class="text">Admin</div>
          </a>
        </li>
        <li :class="{ active: isActive('/settings') }">
          <a @click="goTo('/settings')">
            <div class="icon">
              <i class="fa-solid fa-gear"></i>
            </div>
            <div class="text">Settings</div>
          </a>
        </li>
        <li class="logout">
          <a @click="goTo('/logout')">
            <div class="icon">
              <i class="fa-solid fa-right-from-bracket"></i>
            </div>
            <div class="text">Logout</div>
          </a>
        </li>
      </div>
    </div>
  </ul>
</template>