<script setup>
import { ref, onMounted } from "vue"
import { useRouter } from 'vue-router'

const router = useRouter()
const menuItems = ref([])

onMounted(() => {
    menuItems.value = document.querySelectorAll(".menulist li")
    menuItems.value.forEach(item => item.addEventListener("click", activeLink))
})

const activeLink = function () {
    menuItems.value.forEach(item => item.classList.remove("active"))
    this.classList.add("active")
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
                <li class="active">
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
                <li>
                    <a @click="goTo('/manage')">
                        <div class="icon">
                            <i class="fa-solid fa-book"></i>
                        </div>
                        <div class="text">Manage</div>
                    </a>
                </li>
                <li>
                    <a @click="goTo('/students')">
                        <div class="icon">
                            <i class="fa-solid fa-user-graduate"></i>
                        </div>
                        <div class="text">Students</div>
                    </a>
                </li>
                <li>
                    <a @click="goTo('/courses')">
                        <div class="icon">
                            <i class="fa-solid fa-file-invoice-dollar"></i>
                        </div>
                        <div class="text">Course</div>
                    </a>
                </li>
                <li>
                    <a @click="goTo('/grades')">
                        <div class="icon">
                            <i class="fa-solid fa-star"></i>
                        </div>
                        <div class="text">Grades</div>
                    </a>
                </li>
            </div>

            <div class="sidebar-section">
                <div class="sidebar-label">Settings</div>
                <li>
                    <a href="#admin">
                        <div class="icon">
                            <i class="fa-solid fa-user-shield"></i>
                        </div>
                        <div class="text">Admin</div>
                    </a>
                </li>
                <li>
                    <a href="#settings">
                        <div class="icon">
                            <i class="fa-solid fa-gear"></i>
                        </div>
                        <div class="text">Settings</div>
                    </a>
                </li>
                <li class="logout">
                    <a href="#logout">
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