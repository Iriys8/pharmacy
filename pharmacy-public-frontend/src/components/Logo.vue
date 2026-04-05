<template>
    <div class="logo_placement">
        <img class="logo_img" src="@/assets/Logo.png" alt="Logo.png" />
        <div class="logo_text">
            <slot name="logo"></slot>
        </div>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.logo_placement {
    @apply flex items-start gap-2;
}

.logo_img {
    @apply max-h-12 max-w-12 place-self-center;
}

.logo_text {
    @apply text-4xl font-bold place-self-center dark:text-white;
}
</style>