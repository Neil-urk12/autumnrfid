<script setup>
import { defineAsyncComponent } from "vue";
const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"));
</script>

<template>
    <main>
        <div class="sidebar">
            <Sidebar />
        </div>
        <section>
        </section>
    </main>
</template>
