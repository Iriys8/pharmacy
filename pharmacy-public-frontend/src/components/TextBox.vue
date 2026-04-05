<template>
    <div class="content">
        <div class="box_title">
            <slot name="title"></slot>
        </div>
        <div class="message_text_box_inner">
            <div class="message_text">
                <slot name="text"></slot>
            </div>
            <slot></slot>
        </div>
    </div>
</template>


<style scoped lang="postcss">
@import "tailwindcss";

.box_title {
    @apply text-2xl font-bold place-self-center dark:text-white;
}

.content {
    @apply flex flex-col h-full;
}

.message_text_box_inner {
    @apply flex p-8 border rounded-lg dark:border-white;
}

.message_text {
    @apply place-self-center text-xl font-medium dark:text-white;
}
</style>