<script setup lang="ts">
import { RouterView } from 'vue-router';
import siteheader from '@/components/Header.vue'
import sitefooter from '@/components/Footer.vue'
</script>

<template>
    <div class="base">
        <siteheader />
        <div class="main">
            <RouterView />
        </div>
        <sitefooter />
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.base {
    @apply font-[Inter] min-h-screen flex flex-col;
}

.main {
    @apply flex-1;
}

</style>