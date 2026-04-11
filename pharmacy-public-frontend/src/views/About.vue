<script setup lang="ts">
import textbox from "@/components/TextBox.vue";
import advertbox from "@/components/AdvertBox.vue";
</script>

<template>
    <div class="main_box">
        <div class="content_placement">
            <div class="main_content">
                <div class="about_box">
                    <p class="box_title">About us</p>
                    <div class="about_box_content">
                        <div class="about_box_text">
                            <textbox>
                                <template #title>
                                    Director's speech
                                </template>
                                <template #text>
                                    "Did the lungfish ever refuse air? No. It crawled bravely across the land while its kind remained in the black abyss of the ocean, 
                                    peering into the darkness, ignorant and doomed despite their eternal vigilance. Will we follow in the footsteps of the trilobites? 
                                    Are all the achievements of mankind destined to become mere shards of plastic scattered on the petrified bed between the Burgess Shale 
                                    and the thousand-year-old mud? To be true to ourselves and our destiny, we must strive for more. We have outgrown our cradle. It is useless 
                                    to demand mother's milk when true help awaits us among the opiates. And only a universal union, which the narrow-minded call "hallucinations", 
                                    will help us achieve them." - Dr. Wallace Breen, head of Citadel corp.
                                </template>
                            </textbox>
                        </div>
                        <div class="about_box_picture">
                            <img class="about_picture" src="@/assets/Breen.png" alt="Director"/>
                        </div>
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

.box_title {
    @apply text-2xl font-bold place-self-center pb-2 dark:text-white;
}

.content_placement {
    @apply flex w-full gap-5;
}

.main_content {
    @apply flex flex-col gap-5 w-2/3;
}

.about_box {
    @apply flex flex-col bg-neutral-100 dark:bg-cyan-950 pt-2 p-5 shadow-md rounded-lg;
}

.additional_content {
    @apply flex flex-col w-1/3 bg-neutral-100 dark:bg-cyan-950 pt-2 p-5 shadow-md rounded-lg overflow-y-auto max-h-dvh;
}

.about_box_content {
    @apply flex gap-5;
}

.about_box_text {
    @apply w-2/3 h-full;
}

.about_box_picture {
    @apply flex flex-col w-1/3 h-full justify-center;
}

</style>