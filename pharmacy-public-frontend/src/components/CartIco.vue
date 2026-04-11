<script setup lang="ts">
import { useCartStore } from "@/stores/CartStore";
import { ref, computed, watch } from 'vue';
import { useDark } from '@vueuse/core'

const cartStore = useCartStore();
const cartItems = ref([]);
const isDark = useDark();

const totalQuantity = computed(() => {
    return Object.values(cartItems.value).reduce((sum, quantity) => sum + quantity, 0);
});

const loadCartItems = async () => {
    const cartData = JSON.parse(localStorage.getItem("cart") || "{}");
    cartItems.value = cartData && Object.keys(cartData).length ? cartData : {};
};

watch(
    () => cartStore.items,
    (newItems) => {
        loadCartItems();
    },
    { deep: true, immediate: true }
);
</script>

<template>
    <div>
        <img class="nav_cart" src="@/assets/CartDark.png" v-if="isDark" />
        <img class="nav_cart" src="@/assets/Cart.png" v-else />
        <div v-if="totalQuantity != 0" class="nav_cart_quantity_box">
            <p class="nav_cart_quantity_text">{{ totalQuantity }}</p>
        </div>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.nav_cart {
    @apply max-h-8 min-h-8 hover:drop-shadow-lg hover:drop-shadow-blue-500/50 relative;
}

.nav_cart_quantity_box {
    @apply flex bg-red-600 p-1 aspect-square w-8 rounded-full shadow-md items-center justify-center absolute top-11 right-47 dark:text-white;
}
</style>