<template>
    <div class="main_box">
        <div class="content_placement">
			<div class="content">
				<p class="box_title">Access denied</p>
                <img src="@/assets/Error.png" alt="error" class="img">
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

.content {
    @apply flex flex-col gap-5 w-full h-fit bg-neutral-100 pt-2 p-5 shadow-md rounded-lg;
}
	
.box_title {
    @apply text-2xl font-bold place-self-center;
}

.img {
    @apply place-self-center h-24 w-24;
}
</style>