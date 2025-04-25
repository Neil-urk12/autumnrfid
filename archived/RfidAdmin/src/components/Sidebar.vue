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
  <div class="sidebar">
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
</div>
</template>

<style scoped>
.sidebar {
  position: fixed;
  width: 83px;
  height: 100%;
  background: var(--sidebar-gradient);
  z-index: 100000;
  transition: 0.5s;
  padding-left: 6.5px;
  overflow: hidden;
  border-left: 10px solid #4338ca;
}

.sidebar:hover {
  width: 230px;
}

.sidebar.active {
  width: 230px;
}

.sidebar ul {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding: 30px 0;
  ;
}

.menulist {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  flex: 1;
  position: relative;

}

.menulist li:last-child {
  margin-top: auto;
}

.sidebar ul li {
  position: relative;
  list-style: none;
}

.sidebar ul li.active {
  background: var(--clr);
  border-top-left-radius: 50px;
  border-bottom-left-radius: 50px;
}

.sidebar ul li.active::before {
  content: "";
  position: absolute;
  top: -20px;
  right: 0;
  width: 20px;
  height: 20px;
  border-bottom-right-radius: 20px;
  box-shadow: 5px 5px 0 5px var(--clr);
  background: transparent;
}

.sidebar ul li.active::after {
  content: "";
  position: absolute;
  bottom: -20px;
  right: 0;
  width: 20px;
  height: 20px;
  border-top-right-radius: 20px;
  box-shadow: 5px -5px 0 5px var(--clr);
  background: transparent;
}

.sidebar ul li a {
  position: relative;
  display: flex;
  white-space: nowrap;
  text-decoration: none;
  padding: 3px;
}

.sidebar ul li a .icon {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-width: 40px;
  padding-left: 17px;
  padding-bottom: 1.5px;
  height: 50px;
  font-size: 1em;
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.sidebar ul li a .text {
  position: relative;
  height: 50px;
  display: flex;
  align-items: center;
  font-size: 12px;
  color: var(--text-primary);
  padding-left: 30px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  transition: 0.5s;
}

.sidebar ul li.active a .icon {
  color: #fff;
  z-index: 2;
}

.sidebar ul li.active a .text {
  color: var(--accent)
}

.sidebar ul li:hover a .icon,
.sidebar ul li:hover a .text {
  color: var(--hover-color);

}

.sidebar ul li.active a .icon::before {
  content: "";
  position: absolute;
  inset: 3px;
  width: 45px;
  background: var(--clr);
  border: 3px solid var(--accent);
  border-radius: 50%;
  transition: 0.5s;
  z-index: -1;
  box-shadow: 0 0 15px var(--border-glow),
      inset 0 0 15px var(--border-glow);
  animation: 2s ease-in-out infinite;
}

.sidebar ul li:hover.active a .icon {
  color: #fcfbfb;
  text-shadow: 0 0 5px var(--hover-color),
      0 0 10px var(--hover-color);
}

.sidebar-section {
  margin: 5px 0;
}

.sidebar-label {
  color: var(--text-secondary);
  font-size: 0.7em;
  text-transform: uppercase;
  letter-spacing: 1px;
  padding: 5px 0 10px 15px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.sidebar:hover .sidebar-label {
  opacity: 1;
}

</style>