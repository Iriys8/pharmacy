<script setup lang="ts">
import textbox from "@/components/TextBox.vue";
import scheduletable from "@/components/ScheduleTable.vue";
import advertbox from "@/components/AdvertBox.vue";
import announces from "@/components/SearchWindow.vue"

const now = new Date();

const pad = (num: number): string => {
    return num.toString().padStart(2, "0");
}

const formatDate = (date: Date): string => {
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
}

const addDays = (date: Date, days: number): Date => {
    const copy = new Date(date);
    copy.setDate(copy.getDate() + days);
    return copy;
}
</script>

<template>
    <div class="main_box">
        <div class="content_placement">
            <div class="main_content">
                <div class="message_box">
                    <announcenox class="message_text_box">
                        <textbox #text>
                            Welcome to Apteka-17: you chose it yourself or it was chosen for you - this is the best pharmacy left. 
                            I think so highly of Apteka-17 that I decided to make my home here, at the checkout counter so thoughtfully provided by my Patrons. 
                            I am proud to call Apteka-17 my home. So, whether you are going to buy something or are in for a stroke, welcome to Apteka-17! It is profitable here.
                        </textbox>
                    </announcenox>
                </div>
                <div class="announce_box">
                    <p class="box_title">Announces</p>
                    <announces page_addres="" :is_advert_box="true" :is_search_window="false" type="Announce"/>
                </div>  
                <div class="schedule_box">
                    <div class="schedule_box_content">
                        <scheduletable   :start-date="formatDate(now)" :end-date="formatDate(addDays(now, 10))" >
                            <template #title>
                                Schedule for the next month
                            </template>10 working days
                        </scheduletable>
                    </div>
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
    @apply flex w-full h-fit gap-5;
}

.main_content {
    @apply flex flex-col gap-5 w-2/3 h-full;
}

.message_box {
    @apply flex flex-col bg-neutral-100 dark:bg-cyan-950 pt-2 p-5 shadow-md rounded-lg;
}

.message_text_box {
    @apply w-5/6 place-self-center max-h-1/2;
}

.announce_box {
    @apply bg-neutral-100 dark:bg-cyan-950 flex flex-col pt-2 h-1/2 p-5 shadow-md rounded-lg items-center overflow-y-auto max-h-dvh;
}

.adverting {
    @apply w-5/6 max-w-[128] h-fit place-self-center;
}

.schedule_box {
    @apply min-h-32 bg-neutral-100 dark:bg-cyan-950 flex flex-col pt-2 h-1/2 p-5 shadow-md rounded-lg;
}

.schedule_box_content {
    @apply h-full flex justify-center w-5/6 place-self-center;
}

.additional_content {
    @apply w-1/3 bg-neutral-100 dark:bg-cyan-950 flex flex-col pt-2 p-5 shadow-md rounded-lg overflow-y-auto max-h-dvh;
}

.box_title {
    @apply text-2xl font-bold place-self-center dark:text-white;
}

</style>