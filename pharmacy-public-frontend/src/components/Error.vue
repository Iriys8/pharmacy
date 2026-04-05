<script setup lang="ts">
import { useDark } from '@vueuse/core'

const isDark = useDark();

</script>

<template>
    <div class="error_box">
        <img src="@/assets/ErrorDark.png" alt="error" class="error_img" v-if="isDark">
        <img src="@/assets/Error.png" alt="error" class="error_img" v-else>
        <p class="error_text">Something went wrong, try reload page.</p>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.error_box {
    @apply w-full h-full flex flex-col gap-5 bg-neutral-300 dark:bg-cyan-800 p-10 border justify-center dark:border-white;
}

.error_img {
    @apply place-self-center h-24 w-24;
}

.error_text {
    @apply place-self-center font-bold dark:text-white text-xl;
}
</style>