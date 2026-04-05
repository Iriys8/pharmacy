<script setup lang="ts">
import { RouterLink } from 'vue-router';
import searchline from "@/components/SearchLine.vue";
import logo from "@/components/Logo.vue";
import cartico from "@/components/CartIco.vue"
</script>

<template>
    <header class="header">
        <div class="header_box">
            <RouterLink class="logo_placement" to="/">
                <logo>
                    <template #logo>
                        Apteka-17
                    </template>
                </logo>
            </RouterLink>
            <div class="nav_box">
                <nav class="nav_links">
                    <RouterLink to="/" class="nav_link">Home</RouterLink>
                    <RouterLink to="/catalog?q=" class="nav_link">Catalog</RouterLink>
                    <RouterLink to="/schedule" class="nav_link">Schedule</RouterLink>
                    <RouterLink to="/about" class="nav_link">About</RouterLink>
                </nav>
                <nav class="nav_cart_placement">
                    <RouterLink to="/shoppingcart">
                        <cartico />
                    </RouterLink>
                </nav>
            </div>
        </div>
        <div class="search_box">
            <searchline />
        </div>
    </header>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.header {
    @apply flex flex-col mr-40 ml-40;
}

.header_box {
    @apply flex h-24 w-full justify-between bg-neutral-100 dark:bg-cyan-950 pl-10 pr-10 mb-5 shadow-md;
}

.logo_placement {
    @apply place-self-center;
}

.nav_box {
    @apply w-1/2 flex justify-between;
}

.nav_links {
    @apply flex justify-center gap-6 place-self-center dark:text-white;
}

.nav_link {
    @apply text-xl font-medium hover:text-blue-500 hover:drop-shadow-lg hover:drop-shadow-blue-200/50;
}

.nav_cart_placement {
    @apply place-self-center;
}

.search_box {
    @apply flex w-full mb-5;
}

</style>