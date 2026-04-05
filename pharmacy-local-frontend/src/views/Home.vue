<script setup lang="ts">
import announcesbox from "@/components/SearchWindow.vue"
</script>

<template>
    <div class="main_box">
        <div class="content_placement">
			<div class="content">
                <p class="box_title">Announces</p>
				<announcesbox page_addres="" :is_advert_box="true" :is_search_line="false" :is_search_window="false" type="Announce"/>
    		</div>
        </div>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.main_box {
    @apply flex flex-col mr-40 ml-40 mb-10 bg-neutral-100 pt-2 p-5 shadow-md rounded-lg;
}

.content_placement {
    @apply flex w-full gap-5 ;
}

.box_title {
    @apply text-2xl font-bold place-self-center;
}

.content {
    @apply flex flex-col h-full w-full;
}
	
</style>