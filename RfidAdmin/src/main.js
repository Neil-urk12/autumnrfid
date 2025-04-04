import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import Admin from '@/views/Admin.vue'
import Manage from '@/views/Manage.vue'
import Students from '@/views/Students.vue'
import Courses from '@/views/Courses.vue'
import Grades from '@/views/Grades.vue'


const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'Admin',
            component: Admin
        },
        {
            path: '/manage',
            name: 'Manage',
            component: Manage
        },
        {
            path: '/students',
            name: 'Students',
            component: Students
        }, 
        {
            path: '/courses',
            name: 'Courses',
            component: Courses
        },
        {
            path: '/grades',
            name: 'Grades',
            component: Grades
        }

    ]
})

const app = createApp(App)
app.use(router)
app.mount('#app')