<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useCartStore } from "@/stores/CartStore";
import orderedItem from '@/components/MiniItem.vue';
import errorbox from '@/components/Error.vue';
import type { Order, OrderedItem } from "@/types";
import { goodsAPI, orderAPI } from "@/api";

const cartStore = useCartStore();
const error = ref<boolean>(false);
const isLoaded = ref<boolean>(false);
const inProgress = ref<boolean>(false);
const totalPrice = computed(() =>
    item.value.Items.reduce((sum, item) => sum + item.Price * item.Quantity, 0)
);

const getItems = async () => {
    const cartData: Record<string, number> = JSON.parse(localStorage.getItem("cart") || "{}");
    if (!cartData || Object.keys(cartData).length === 0) {
        item.value.Items = [];
        return;
    }
    try{
        const itemPromises = Object.keys(cartData).map(async (id) => {
            const response = await goodsAPI.getGood(parseInt(id)) as OrderedItem;
            return { ...response, Quantity: cartData[id] || 1 };
        })
        const loadedItems = await Promise.all(itemPromises);
        item.value.Items = [...loadedItems];
        isLoaded.value = true;
    }
	catch (er) {
        error.value = true;
    }
}

const item = ref<Order>({
    Name: "",
    Email: "",
    Phone: "",
    Items: [],
})

const submitOrder = async () => {
    if (inProgress.value) return;
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!item.value.Phone || item.value.Phone.trim().length === 0 || !item.value.Name || item.value.Name.trim().length === 0 || !item.value.Items || item.value.Items.length === 0)  {
        alert("Error: Missing one of required fields");
        return;
    }
    if (item.value.Email !== undefined && item.value.Email !== "" && (!emailRegex.test(item.value.Email) || item.value.Email !== "")) {
        alert("Error: Wrong email format");
        return;
    }
    inProgress.value = true;
    try {        
        await orderAPI.createOrder(item.value);
        alert("Order submited.");
        cartStore.clearCart();
    }
    catch(err) {
        alert("Error: " + err)
    }
    finally {
        inProgress.value = false;
    }
};

watch(
    () => cartStore.items,
    () => {
        getItems();
    },
    { deep: true, immediate: true }
);
</script>

<template>
    <div v-if="error">
        <errorbox />
    </div>
    <div v-else-if="item.Items.length && isLoaded">
        <div class="products_box">
            <div class="box_title">
                <div class="box_title_text">
                    <slot name="title"></slot>
                </div>
            </div>
            <div v-for="item in item.Items" :key="item.ID" class="product_box">
                <ordered-item :item="item" :is_add_action="false" :is_in_window="true" :is_in_advert_box="false" class="product" />
                <div class="product_info_box">
                    <div class="product_info_box_title">
                        <p>Quantity:</p>
                        <p>Price:</p>
                    </div>
                    <div class="product_info_box_info">
                        <p>{{ item.Quantity }}</p>
                        <p>{{ item.Quantity * item.Price }} Rub.</p>
                    </div>
                </div>
            </div>
            <div class="products_summary_box">
                <p class="products_summary">Summary: {{ totalPrice }} Rub.</p>
            </div>
        </div>

        <div class="submit_form">
            <input v-model="item.Name" class="submit_field" type="text" placeholder="FIO" required />
            <input v-model="item.Email" class="submit_field" type="email" placeholder="Email (optional)" />
            <input v-model="item.Phone" class="submit_field" type="text" placeholder="Phone" required />
            <button :disabled="!isLoaded" @click.prevent="submitOrder()" class="submit_button">Commit order</button>
        </div>
    </div>
    <div v-else class="product_car_is_empty_box">
        <p class="products_car_is_empty">Cart is empty</p>
    </div>
</template>


<style scoped lang="postcss">
@import "tailwindcss";

.box_title_text {
    @apply bg-blue-500 rounded-lg shadow-md p-2;
}

.box_title {
    @apply place-self-center text-4xl font-bold w-fit text-white;
}

.products_box {
    @apply flex flex-col gap-5;
}

.product_box {
    @apply flex flex-row gap-2;
}

.product {
    @apply w-2/3;
}

.product_info_box {
    @apply flex flex-row bg-white dark:bg-cyan-800 dark:text-white rounded-lg shadow-md p-5 w-1/3 justify-center items-center;
}

.product_info_box_title {
    @apply flex flex-col w-fit;
}

.product_info_box_info {
    @apply flex flex-col w-full items-end;
}

.products_summary_box {
    @apply bg-blue-500 rounded-full shadow-md w-fit p-2;
}

.products_summary {
    @apply text-xl font-semibold text-white;
}

.product_car_is_empty_box {
    @apply flex flex-col;
}

.products_car_is_empty {
    @apply text-4xl font-bold pb-10 w-fit place-self-center dark:text-white;
}

.submit_form {
    @apply flex flex-row p-5 gap-5 bg-white dark:bg-cyan-800 rounded-lg shadow-md mt-5;
}

.submit_field {
    @apply border p-2 rounded-lg dark:text-white;
}

.submit_button {
    @apply border bg-blue-500 p-2 rounded-lg w-full hover:bg-blue-700 text-white font-medium;
}
</style>