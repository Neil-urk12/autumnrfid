import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const routes = [
    {
        path: "/manage",
        name: "Manage",
        component: () => import("@/views/Manage.vue"),
    },
    {
        path: "/",
        redirect: "/manage",
    },
    {
        path: "/admin",
        name: "Admin",
        component: () => import("@/views/Admin.vue"),
    },
    {
        path: "/login",
        name: "Login",
        component: () => import("@/views/LoginView.vue"),
    },
    {
        path: "/courses",
        name: "Courses",
        component: () => import("@/views/Courses.vue"),
    },
    {
        path: "/grades",
        name: "Grades",
        component: () => import("@/views/Grades.vue"),
    },
    {
        path: "/students",
        name: "Students",
        component: () => import("@/views/Students.vue"),
    },
    {
        path: "/bills",
        name: "Bills",
        component: () => import("@/views/Bills.vue"),
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes
});

// Navigation guard to check authentication
router.beforeEach((to, from, next) => {
    const authStore = useAuthStore();
    const publicPages = ['/login'];
    const authRequired = !publicPages.includes(to.path);
    
    if (authRequired && !authStore.isAuthenticated) {
        // Redirect to login page if not authenticated
        return next('/login');
    }
    
    // If already logged in and trying to access login page, redirect to home
    if (authStore.isAuthenticated && to.path === '/login') {
        return next('/');
    }
    
    next();
});

export default router
