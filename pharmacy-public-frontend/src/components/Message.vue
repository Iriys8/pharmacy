<script setup lang="ts">
import { useDark } from '@vueuse/core'

const isDark = useDark();

const props = defineProps<{
    isError: boolean,
}>();

</script>

<template>
    <div class="message_box">
        <img src="@/assets/ErrorDark.png" alt="error" class="error_img" v-if="isDark && isError">
        <img src="@/assets/Error.png" alt="error" class="error_img" v-else-if="isError">
        <div class="message_text">
            <slot name="text"></slot>
        </div>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.message_box {
    @apply w-full h-full flex flex-col gap-5 bg-neutral-300 dark:bg-cyan-800 p-10 border justify-center dark:border-white;
}

.error_img {
    @apply place-self-center h-24 w-24;
}

.message_text {
    @apply place-self-center font-bold dark:text-white text-xl;
}
</style>