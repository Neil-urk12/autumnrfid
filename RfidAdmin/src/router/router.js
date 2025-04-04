import { createRouter, createWebHistory } from "vue-router";
import Manage from "@/views/Manage.vue";
const routes = [
    {
        path: "/manage",
        name: "Manage",
        component: Manage,
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router