<script setup lang="ts">
import { ref, watch } from "vue";
import { useRoute } from "vue-router";
import errorbox from "@/components/Error.vue";
import itembox from "@/components/ItemBox.vue";
import advertbox from "@/components/AdvertBox.vue";
import { goodsAPI } from "@/api";
import type { Goods } from "@/types";

const route = useRoute();
const item = ref<Goods>({
    ID: 0,
    Name: "",
    Image: "",
    Instruction: "",
    Description: "",
    IsPrescriptionNeeded: false,
	IsInStock: true,
    Price: 0,
    Producer: "",
    Tags: undefined,
	Type: "Goods"
});

const fetchProduct = async () => {
    try {
        item.value = await goodsAPI.getGood(Number(route.params.id)) as Goods
    } catch (error) {
        console.error(error);
    }
};

watch(() => route.params.id, fetchProduct, { immediate: true });
</script>

<template>
    <div class="main_box">
        <div class="content_placement">
            <div class="main_content">
                <div class="item_box">
                    <itembox v-if="item" :product="item" />
                    <errorbox v-else />
                </div>
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
    @apply flex flex-col gap-5 w-2/3 h-fit dark:text-white;
}

.item_box {
    @apply flex flex-col bg-neutral-100 dark:bg-cyan-950 pt-2 p-5 shadow-md rounded-lg;
}

.additional_content {
    @apply flex flex-col w-1/3 bg-neutral-100 dark:bg-cyan-950 pt-2 p-5 shadow-md rounded-lg overflow-y-auto max-h-dvh;
}

</style>