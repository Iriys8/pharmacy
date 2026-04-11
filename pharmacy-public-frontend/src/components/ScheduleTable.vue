<script setup lang="ts">
import { ref, watch } from "vue";
import type { shedule } from "@/types";
import { scheduleAPI } from "@/api";
import messagebox from "@/components/Message.vue"

const scheduleData = ref<shedule[]>([]);

const loaded = ref<boolean>(false)
const err = ref<string>("");

const props = defineProps<{
    startDate: string,
    endDate:string,
}>();

const fetchData = async () => {
    if (props.startDate != "" && props.endDate != "") {
        try {
            scheduleData.value =  await scheduleAPI.getShedule(props.startDate, props.endDate)
            loaded.value = true
        } catch(error){
            err.value = String(error)
        }
    }
}

watch([() => props.startDate, () => props.endDate], () => fetchData(), { immediate: true });
</script>

<template>
    <div class="content_placement">
        <div class="box_title">
            <slot name="title"></slot>
        </div>
        <table class="schedule_table" v-if="loaded">
            <caption class="schedule_table_name">Work schedule for {{ startDate }} - {{ endDate }}</caption>
            <thead>
                <tr class="schedule_table_title_box">
                    <th class="schedule_table_col_title">Date</th>
                    <th class="schedule_table_col_title">Open time</th>
                    <th class="schedule_table_col_title">Close time</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="date in scheduleData">
                    <td class="schedule_table_col_element">{{ date.Date }}</td>
                    <td class="schedule_table_col_element" v-if="date.IsOpened">{{ date.TimeStart }}</td>
                    <td class="schedule_table_col_element" v-if="date.IsOpened">{{ date.TimeEnd }}</td>
                    <td class="schedule_table_col_element" colspan="2" v-if="!date.IsOpened">Closed</td>
                </tr>
            </tbody>
        </table>
        <div v-else-if="err">
            <messagebox :is-error="true" #text>
                An error occurred, please reload the page
            </messagebox>
        </div>
        <div v-else>
            <messagebox :is-error="false" #text>
                Loading...
            </messagebox>
        </div>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.content_placement {
    @apply flex flex-col;
}

.box_title {
    @apply text-2xl font-bold place-self-center dark:text-white;
}

.schedule_table {
    @apply border w-full caption-bottom dark:border-white;
}

.schedule_table_name {
    @apply text-neutral-400 dark:text-neutral-500 font-extralight;
}

.schedule_table_col_title {
    @apply border w-1/3  font-semibold text-xl dark:border-white dark:text-white;
}

.schedule_table_col_element {
    @apply border font-normal text-xl text-center dark:text-white;
}
</style>