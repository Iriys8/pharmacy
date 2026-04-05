<script setup lang="ts">
import SearchWindow from '@/components/SearchWindow.vue';
</script>

<template>
    <div class="main_box">
        <div class="content_placement">
            <SearchWindow :is_search_window=false :is_advert_box=false page_addres="/catalog" type="Goods" />
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
