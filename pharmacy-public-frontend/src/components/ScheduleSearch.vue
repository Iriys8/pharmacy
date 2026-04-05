<script setup lang="ts">
import { ref, watch } from "vue";
import scheduletable from "@/components/ScheduleTable.vue";

const startDate = ref<string>("");
const endDate = ref<string>("");

watch([startDate, endDate], () => {
    if (startDate.value && endDate.value) {
        const start = new Date(startDate.value);
        const end = new Date(endDate.value);
        const diff = (end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24);

        if (diff < 0) {
            alert("Start date cannot be later than the end date!");
            endDate.value = "";
        } else if (diff > 10) {
            alert("Date range cannot exceed 10 days!");
            endDate.value = "";
        }
    }
});
</script>

<template>
    <div class="interface">
        <div class="box_title">
            <slot name="title"></slot>
        </div>
        <div class="input">
            <label class="input_field_label">Start Date:</label>
            <input type="date" v-model="startDate" onkeydown="return false" class="input_field" />
        </div>
        <div class="input">
            <label class="input_field_label">End Date:</label>
            <input type="date" v-model="endDate" onkeydown="return false" class="input_field" />
        </div>
    </div>
    <div>
        <scheduletable :startDate="startDate" :endDate="endDate" />
    </div>
</template>


<style scoped lang="postcss">
@import "tailwindcss";

.interface {
    @apply flex flex-col;
}

.box_title {
    @apply text-2xl font-bold place-self-center pb-2 dark:text-white;
}

.input {
    @apply mb-4 dark:text-white;
}

.input_field_label {
    @apply block text-gray-600 dark:text-gray-500 mb-2;
}

.input_field {
    @apply border border-gray-300 dark:border-gray-600 rounded-md p-2 w-full;
}

.input_field_postlabel {
    @apply text-sm text-gray-400 dark:text-gray-500 mt-2;
}
</style>