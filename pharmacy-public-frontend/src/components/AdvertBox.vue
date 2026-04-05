<script setup lang="ts">
import errorbox from '@/components/Error.vue';
import { RouterLink } from 'vue-router';
import { ref, onMounted } from "vue";
import type { PromoItem } from '@/types';
import { goodsAPI, imagesAPI } from '@/api/';


const promoItems = ref<PromoItem[]>([]);

const loadPromoItems = async () => {
    try {
        promoItems.value = await goodsAPI.getPromoItem()
    } catch (error) {
        console.log(error);
    }
};

onMounted(loadPromoItems);
</script>

<template>
    <div v-if="promoItems.length" class="content_placement">
        <div class="title">
            <slot name="title"></slot>
        </div>
        <div v-for="item in promoItems" :key="item.id" class="advert_item_box">
            <RouterLink :to="'/item/' + item.id" class="advert_item_title_box">
                <img :src="imagesAPI.getImageSRC() + item.image" alt="No image" class="advert_item_picture" />
                <p class="advert_item_title">{{ item.name }}</p>
            </RouterLink>
            <p class="advert_item_discription"> {{ item.description.length > 50 ? item.description.slice(0, 50) + '...'
                : item.description }} </p>
            <p class="advert_item_price">{{ item.price }} Rub.</p>
        </div>
    </div>
    <div v-else>
        <errorbox />
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.content_placement {
    @apply flex flex-col;
}

.advert_item_box {
    @apply bg-white dark:bg-cyan-800 p-2 rounded-lg shadow-md w-full text-center;
}

.advert_item_title_box {
    @apply flex flex-col h-fit;
}

.advert_item_picture {
    @apply w-full object-contain dark:text-white;
}

.advert_item_title {
    @apply text-blue-500 text-xl font-semibold mt-2 hover:underline;
}

.advert_item_discription {
    @apply text-sm font-light h-fit dark:text-white;
}

.advert_item_price {
    @apply text-xs text-gray-600 dark:text-gray-400;
}

.title {
    @apply text-2xl font-bold place-self-center pb-2 dark:text-white;
}
</style>