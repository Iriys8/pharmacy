<script setup lang="ts">
import { RouterLink } from 'vue-router';
</script>

<template>
    <footer class="footer">
        <div class="footer_box">
            <div class="footer_objects_placement">
                <p class="footer_text">Citadel corp. No one right protected.</p>
            </div>
            <div class="footer_objects_placement">
                <router-link to="#top" class="footer_to_the_top">To the top</router-link>
            </div>
        </div>
    </footer>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.footer {
    @apply flex bg-neutral-100 dark:bg-cyan-950 mr-40 ml-40 mt-5 pl-10 pr-10 shadow-md;
}

.footer_box {
    @apply flex h-12 w-full justify-between;
}

.footer_objects_placement {
    @apply place-self-center;
}

.footer_text {
    @apply text-neutral-500 dark:text-neutral-400 font-light;
}

.footer_to_the_top {
    @apply text-neutral-500 dark:text-neutral-400 font-light hover:underline;
}

</style>