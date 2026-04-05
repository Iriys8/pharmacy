<script setup lang="ts">
import SearchWindow from '@/components/SearchWindow.vue';
</script>

<template>
    <div class="main_box">
        <div class="content_placement">
            <SearchWindow :is_advert_box="false" :is_search_line="true" :is_search_window=false page_addres="/goods" type="Goods" />
        </div>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.main_box {
    @apply flex flex-col mr-40 ml-40 mb-10;
}

.content_placement {
    @apply flex w-full;
}

</style>
