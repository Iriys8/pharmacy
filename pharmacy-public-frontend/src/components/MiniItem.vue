<script setup lang="ts" generic="T extends Goods | Announce | OrderedItem">
import { imagesAPI } from "@/api"

import type { Announce, Goods, OrderedItem } from "@/types";
import { useCartStore } from "@/stores/CartStore";

const props = defineProps<{
    item: T;
    is_in_window: boolean;
    is_add_action: boolean;
    is_in_advert_box: boolean;
}>();

const emit = defineEmits<{
    (e: 'item', value: T): void
}>();

const selectItem = () => {
    emit('item', props.item);
}

const cartStore = useCartStore();
const addToCart = (id: number) => {
    cartStore.addToCart(id);
}
const removeFromCart = (id: number) => cartStore.removeFromCart(id);
    
</script>

<template>
    <!-- Goods -->
    <div class="item_box" v-if="item.Type === 'Goods'">
        <img :src="imagesAPI.getImageSRC() + item.Image" alt="No image" class="item_picture"/>
        <div class="item_text_box">
            <p class="item_name">
                <router-link :to="'/item/' + item.ID" class="item_name_link">
                    {{ item.Name }}
                </router-link>
            </p>
            <p class="item_discription">
                {{ item.Description && item.Description.length > 100 ? item.Description.slice(0, 100) + '...' :
                    item.Description || "No description available" }}
            </p>
        </div>
        <div class="item_price_box">
            <p class="item_price_numbers">{{ item.Price }}</p>
            <p class="item_price_text">Rub.</p>
        </div>
        <div>
            <button class="item_button" @click="is_add_action ? addToCart(item.ID) : removeFromCart(item.ID)" v-if="item.IsInStock">{{ is_add_action ? "Add to cart" : "Remove from cart" }}</button>
            <button class="item_button_disabled" v-else>Out of stock</button>
        </div>
    </div>

    <!-- Announce -->
    <div class="item_box" v-if="item.Type === 'Announce' && is_in_advert_box">
        <div class="announce_box">
            <p class="announce_name ">
                {{ item.DateTime }}
            </p>
            <p class="item_discription">
                From: {{ item.From }}
            </p>
            <p class="announce_text">
                {{ item.Announce }}
            </p>
        </div>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.item_box {
    @apply flex items-center bg-white dark:bg-cyan-800 shadow-md p-4 rounded-lg w-full gap-5;
}

.item_picture {
    @apply w-24 h-24 object-cover rounded-md;
}

.item_text_box {
    @apply ml-4 mr-4 w-11/12 flex flex-col justify-start self-start h-full;
}

.item_name_link {
    @apply text-blue-500 hover:underline;
}

.item_name_button {
    @apply hover:text-blue-500 hover:drop-shadow-lg hover:drop-shadow-blue-400;
}

.item_name {
    @apply text-2xl font-bold;
}

.item_discription {
    @apply text-gray-600 dark:text-white text-sm h-full;
}

.item_price_box {
    @apply flex flex-row gap-1 h-full w-1/12 items-center justify-end;
}

.item_price_numbers {
    @apply font-bold dark:text-white text-xl w-fit h-fit;
}

.item_price_text {
    @apply font-normal dark:text-white text-sm w-fit h-fit;
}

.item_button {
    @apply bg-blue-500 text-white font-extrabold px-4 py-2 rounded-md hover:bg-blue-700 w-full h-fit;
}

.item_button_disabled {
    @apply bg-neutral-500 text-white font-extrabold px-4 py-2 rounded-md w-full h-fit;
}

.announce_name {
    @apply text-2xl font-bold self-center dark:text-white;
}

.announce_box {
    @apply ml-4 mr-4 w-full flex flex-col justify-start self-start h-full;
}

.announce_text {
    @apply dark:text-white text-xl w-fit h-fit;
}
</style>