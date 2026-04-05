<script setup lang="ts">
import usercart from "@/components/Cart.vue";
import advertbox from "@/components/AdvertBox.vue";
</script>

<template>
    <div class="main_box">
        <div class="content_placement">
            <div class="main_content">
                <usercart>
                    <template #title>
                        Cart
                    </template>
                </usercart>
            </div>
            <div class="additional_content">
                <advertbox>
                    <template #title>
                        May be interesting
                    </template>
                </advertbox>
            </div>
        </div>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.main_box {
    @apply flex flex-col mr-40 ml-40 mb-10;
}

.content_placement {
    @apply flex w-full gap-5;
}

.main_content {
    @apply flex flex-col gap-5 w-2/3 h-full bg-neutral-100 dark:bg-cyan-950 pt-2 p-5 shadow-md rounded-lg;
}

.additional_content {
    @apply flex flex-col w-1/3 bg-neutral-100 overflow-y-auto max-h-dvh pt-2 p-5 shadow-md rounded-lg dark:bg-cyan-950;
}

</style>