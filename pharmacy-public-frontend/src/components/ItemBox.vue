<script setup lang="ts">
import { ref } from "vue";
import messagebox from "@/components/Message.vue";
import { useCartStore } from "@/stores/CartStore";
import type { Goods } from "@/types";
import { imagesAPI } from "@/api/";


defineProps<{
    product: Goods
}>();

const cartStore = useCartStore();
const addToCart = (id: number) => {
    cartStore.addToCart(id);
}
const tags = ref<string[]>([]);

</script>

<template>
    <div class="content">
        <div v-if="product" class="content_placement">
            <div class="main_content">
                <img :src="imagesAPI.getImageSRC() + product.Image" alt="No image" class="item_picture" />
                <div class="main_content_info">
                    <div class="main_content_text">
                        <p class="item_name">{{ product.Name }}</p>
                        <router-link :to="{ path: '/catalog', query: { q: product.Producer } }"
                            class="item_producer">
                            {{ product.Producer }}
                        </router-link>
                        <p v-if="product.IsPrescriptionNeeded" class="prescription_text">Prescription needed!</p>
                        <div>
                            <p v-if="!product.IsInStock" class="prescription_text">Out of stock!</p>
                            <p v-else class="prescription_text">In stock</p>
                        </div>
                    </div>
                    <div class="main_content_func">
                        <p class="item_price">{{ product.Price }} Rub.</p>
                        <button class="item_button" @click="addToCart(product.ID)" v-if="product.IsInStock">Add to cart</button>
                        <button class="item_button_disabled" v-else>Add to cart</button>
                    </div>
                </div>
            </div>
            <div v-if="tags.length > 0" class="item_tags_position">
                <div class="item_tags_box" v-for="tag in tags" :key="tag">
                    <router-link :to="{ path: '/catalog', query: { q: tag } }" class="item_tags_text">{{ tag
                    }}</router-link>
                </div>
            </div>
            <div>
                <p class="item_description_title">Description:</p>
                <p>{{ product.Description || "No description available" }}</p>
            </div>
            <div>
                <p class="item_instruction_title">Instruction:</p>
                <p>{{ product.Instruction || "No instruction available" }}</p>
            </div>
        </div>
        <div v-else>
            <messagebox :is-error="true" />
        </div>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.content {
    @apply w-full h-fit;
}

.content_placement {
    @apply flex flex-col w-full gap-5;
}

.main_content {
    @apply flex flex-row justify-between w-full gap-5;
}

.main_content_info {
    @apply flex flex-col justify-center gap-[20%] w-2/3;
}

.item_picture {
    @apply w-1/3 object-cover rounded-lg shadow-md border-2 border-black;
}

.item_name {
    @apply text-2xl font-bold w-fit;
}

.item_producer {
    @apply text-gray-400 font-light hover:underline w-fit;
}

.prescription_text {
    @apply w-fit;
}

.main_content_text {
    @apply flex flex-col;
}

.main_content_func {
    @apply flex flex-col;
}

.item_price {
    @apply text-xl font-normal w-fit;
}

.item_button {
    @apply bg-blue-500 text-white font-extrabold px-4 py-2 rounded hover:bg-blue-700;
}

.item_button_disabled {
    @apply bg-neutral-500 text-white font-extrabold px-4 py-2 rounded;
}

.item_tags_position {
    @apply flex flex-row justify-start gap-2;
}

.item_tags_box {
    @apply bg-blue-500 rounded-full shadow-md p-2 w-fit;
}

.item_tags_text {
    @apply text-white hover:underline;
}

.item_description_title {
    @apply font-semibold;
}

.item_instruction_title {
    @apply font-semibold;
}
</style>
