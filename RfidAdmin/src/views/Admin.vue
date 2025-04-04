<script setup>
import { defineAsyncComponent } from "vue";
const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"));
</script>

<template>
    <div class="sidebar">
        <Sidebar />
    </div>

    <section id="#dashboard">
    </section>
</template>
