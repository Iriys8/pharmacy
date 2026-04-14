<script setup lang="ts">
import SearchWindow from '@/components/SearchWindow.vue';
import { useAuthStore } from '@/stores';
import { useRouter } from 'vue-router';

const router = useRouter();

const Create = async () => {
	router.push({ path: ("scheduleedit/0"), query: { q: "" }});
}
</script>

<template>
    <div class="main_box">
        <div class="content_placement">
            <button class="create_button" @click="Create" v-if="useAuthStore().user?.permissions?.includes('Create_Schedule')">Create</button>
            <SearchWindow :is_advert_box="false" :is_search_line="true" :is_search_window=false page_addres="/schedule" type="Schedule" />
        </div>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.main_box {
    @apply flex flex-col mr-40 ml-40 mb-10;
}

.content_placement {
    @apply flex flex-col w-full gap-5;
}

.create_button {
	@apply w-full text-xl font-medium place-self-center bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500;
}

</style>
