import { createRouter, createWebHistory } from "vue-router";
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
    // {
    //     path: "/login",
    //     name: "Login",
    //     component: () => import("@/views/Login.vue"),
    // },
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

export default router
