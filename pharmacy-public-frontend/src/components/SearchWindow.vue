<script setup lang="ts" generic="T extends Goods | Announce">
import { ref, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import messagebox from "@/components/Message.vue";
import type{ Announce, Goods} from "@/types";
import miniitem from "@/components/MiniItem.vue";
import { announcesAPI, goodsAPI } from "@/api/"
import { useDark } from '@vueuse/core'

const isDark = useDark();
type itemsType = 'Goods' | 'Announce';

const props = defineProps<{
    is_advert_box: boolean;
    is_search_window: boolean;
    page_addres: string;
    type: itemsType;
}>();

const emit = defineEmits<{
    (e: 'item', value: T): void
}>();

const selectItem = (value: T) => {
    emit('item', value);
}

const items = ref<T[]>([]);
const route = useRoute();
const router = useRouter();
const error = ref<boolean>(false);
const noResults = ref<boolean>(false);
const isLoaded = ref<boolean>(false);
const totalPages = ref<number>(1);
const currentPage = ref<number>(1);
const perPage = ref<number>(10);

const fetchData = async (): Promise<void> => {
    isLoaded.value = false;
    const searchQuery = route.query.q as string;
    const page = (route.query.page as string | undefined) ?? "1";
    const limit = (route.query.limit as string | undefined) ?? perPage.value.toString();

    try {
        var response;
        switch(props.type){
            case 'Goods':
                response = await goodsAPI.getGoods(searchQuery, page, limit);
                break;
            case 'Announce':
                response = await announcesAPI.getAnnounces(searchQuery, page, limit);
                break;
        }
        items.value = response.Items as T[];
        totalPages.value = response.TotalPages as number;
        currentPage.value = response.CurrentPage as number;
        noResults.value = response.Items === null;
        isLoaded.value = true;
    } catch (err) {
        console.log(err);
        error.value = true;
    }
};

const changePage = (newPage: number): void => {
    router.push({
        path: props.page_addres,
        query: {
            q: route.query.q,
            page: newPage.toString(),
            limit: perPage.value.toString(),
        },
    });
};

const changeLimit = (): void => {
    router.push({
        path: props.page_addres,
        query: { q: route.query.q, page: "1", limit: perPage.value.toString() },
    });
};

onMounted(()=>{
    fetchData();
});
watch(() => [route.query.q, route.query.page, route.query.limit], fetchData);    
</script>

<template>
    <div class="main_content" v-if="isLoaded">
        <div v-if="error">
            <messagebox :is-error="true" #text>
                An error occurred, please reload the page
            </messagebox>
        </div>
        <div v-else-if="noResults">
            <div class="no_results">
                <img src="@/assets/ErrorDark.png" alt="error" class="no_results_picture" v-if="isDark">
                <img src="@/assets/Error.png" alt="error" class="no_results_picture" v-else>
                <p class="no_results_text">Nothing found</p>
            </div>
        </div>
        <div v-else>
            <div class="items">
                <miniitem v-for="item in items as T[]" :item="item" :is_in_window="is_search_window" :is_add_action="true" :is_in_advert_box="is_advert_box" @item="value => selectItem(value as T)"/>
            </div>
            <div class="items_limit_box">
                <label for="perPage" class="items_limit_text">Products on page:</label>
                <select id="perPage" v-model="perPage" @change="changeLimit" class="items_limit_shoose">
                    <option value="10">10</option>
                    <option value="20">20</option>
                    <option value="30">30</option>
                    <option value="40">40</option>
                </select>
            </div>
            <div class="items_pages_box">
                <button v-if="currentPage > 1" @click="changePage(currentPage - 1)" class="items_pages_buttons">
                    ⬅️
                </button>
                <button v-if="currentPage > 2" @click="changePage(currentPage - 1)" class="items_pages_buttons">
                    {{ currentPage - 1 }}
                </button>
                <button class="items_pages_current">{{ currentPage }}</button>
                <button v-if="currentPage < totalPages - 1" @click="changePage(currentPage + 1)"
                    class="items_pages_buttons">
                    {{ currentPage + 1 }}
                </button>
                <span v-if="currentPage < totalPages - 2" class="items_page_a_lot_of_pages">...</span>
                <button v-if="currentPage < totalPages" @click="changePage(totalPages)"
                    class="items_pages_buttons">
                    {{ totalPages }}
                </button>
                <button v-if="currentPage < totalPages" @click="changePage(currentPage + 1)"
                    class="items_pages_buttons">
                    ➡️
                </button>
            </div>
        </div>
    </div>
    <div v-else>
        <messagebox :is-error="false" #text>
            Loading...
        </messagebox>
    </div>
</template>

<style lang="postcss" scoped>
@import "tailwindcss";

.main_content {
    @apply flex flex-col gap-5 w-full h-fit bg-neutral-100 dark:bg-cyan-950 pt-2 p-5 shadow-md rounded-lg;
}

.items {
    @apply grid gap-4;
}

.no_results {
    @apply flex flex-col border rounded-lg pt-2 p-5 bg-neutral-300 dark:bg-cyan-800;
}

.no_results_picture {
    @apply place-self-center h-24 w-24;
}

.no_results_text {
    @apply place-self-center text-2xl dark:text-white;
}

.items_limit_box {
    @apply flex justify-center mt-4 dark:text-white;
}

.items_limit_text {
    @apply mr-2;
}

.items_limit_choose {
    @apply border rounded px-2 py-1;
}

.items_pages_box {
    @apply flex justify-center mt-4 space-x-2;
}

.items_pages_buttons {
    @apply px-3 py-1 border rounded hover:bg-blue-500 hover:text-white dark:text-white dark:hover:border-cyan-800;
}

.items_pages_current {
    @apply px-3 py-1 bg-blue-500 text-white rounded;
}

.items_page_a_lot_of_pages {
    @apply px-3 py-1 text-gray-500 dark:text-gray-400;
}
</style>