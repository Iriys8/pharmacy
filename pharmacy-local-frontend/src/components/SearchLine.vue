<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';

const searchQuery = ref<string>('');
const router = useRouter();
const route = useRoute();

const props = defineProps<{
    page: string; 
}>();


onMounted(() => {
    searchQuery.value = route.query.q as string ?? '';
})

const search = (): void => {
    const trimmedQuery: string = searchQuery.value.trim();
    router.push({ path: props.page, query: trimmedQuery ? { q: trimmedQuery } : { q: "" } });
};
</script>

<template>
    <div class="search_position">
        <input v-model="searchQuery" type="text" placeholder="Search in catalog..." class="search_field" @keyup.enter="search" />
        <button class="search_button" @click="search">🔍</button>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.search_position {
    @apply flex flex-row w-full border rounded-lg bg-neutral-100 focus:outline-none focus:ring shadow-md;
}

.search_field {
    @apply w-full px-2 py-2 focus:outline-none;
}

.search_button {
    @apply bg-neutral-400 px-2 border hover:outline-none rounded-lg hover:bg-blue-500;
}
</style>