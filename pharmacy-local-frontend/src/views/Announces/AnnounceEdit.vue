<script setup lang="ts">
import { onMounted, ref } from "vue";
import type { Announce } from "@/types";
import { announcesAPI } from '@/api';
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores";

const route = useRoute();
const isLoaded = ref(false);

const getItem = async () => {
	if (route.params.id != '0') {
    	item.value = await announcesAPI.getAnnounceByID(Number(route.params.id));
	}
	isLoaded.value = true;
}

onMounted(async () => {
    await getItem();
});

const item = ref<Announce>({
    ID: 0,
    DateTime: "",
    From: "",
    Announce: "",
    Type: "Announce",
})

const inProgress = ref(false);
const router = useRouter();
const auth = useAuthStore()
const username = auth.user?.username

const Create = async () => {
    if (inProgress.value) return;
    if (!item.value.Announce || item.value.Announce.trim().length === 0) {
        alert("Error: Empty announce");
        return;
    }
    inProgress.value = true;

    if (username === undefined) {
        alert("Error: not a user");
        return;
    }

    item.value.From = username;

    try {        
        await announcesAPI.createAnnounce(item.value);
        alert("Done.");
		router.push({ path: "/announces"});
    }
    catch(err) {
        alert("Error: " + err)
    }
    finally {
        inProgress.value = false;
    }
}

const Update = async () => {
    if (inProgress.value) return;
    if (!item.value.Announce || item.value.Announce.trim().length === 0) {
        alert("Error: Empty announce");
        return;
    }
    inProgress.value = true;
    try {
        if (!item.value || item.value.ID === 0) {
            alert("Error: Item does not exist");
            return;
        }
   
        await announcesAPI.updateAnnounceByID(item.value);;
        alert("Done.");
		router.push({ path: "/announces"});
    }
    catch(err) {
        alert("Error: " + err)
    }
    finally {
        inProgress.value = false;
    }
}

const Delete = async () => {
	if (inProgress.value) return;
	try {
        if (!item.value || item.value.ID === 0) {
            alert("Error: Item does not exist");
            return;
        }
		
		const result = confirm("Delete announce from " + item.value.From+ "?" )
		if (result) {
			await announcesAPI.deleteAnnounceByID(item.value);
        	alert("Done.");
			router.push({ path: "/announces"});
		}	
    }
    catch(err) {
        alert("Error: " + err)
    }
}

</script>

<template>
    <div class="main_box">
        <div class="content_placement">
            <div class="main_content">
					<div class="control_panel" v-if="item.ID !== 0 || (useAuthStore().user?.permissions?.includes('Create_Announces') && item.ID === 0)">
                    <div class="update_form">
                        <h2 class="box_name">{{ item.ID !== 0 ? `Announce from ${item.From}` : 'Create new announce' }}</h2>
                    
	            		<div>
                            <label class="form_label" for="announce">Announce:</label>
                            <textarea
                                class="form_input"
                                id="announce"
                                v-model="item.Announce"
                                :disabled="!((useAuthStore().user?.permissions?.includes('Update_Announces') && item.ID !== 0) || (useAuthStore().user?.permissions?.includes('Create_Announces') && item.ID === 0))"
                            />
                        </div>
                    </div>
                
                    <div class="control_buttons" v-if="useAuthStore().user?.permissions?.includes('Create_Announces') || useAuthStore().user?.permissions?.includes('Update_Announces') || useAuthStore().user?.permissions?.includes('Delete_Announces')">
                        <button @click.prevent="Create" class="form_blue_button" :disabled="inProgress" v-if="item.ID === 0 && useAuthStore().user?.permissions?.includes('Create_Announces')">
                            Create
                        </button>
                        <button @click.prevent="Update" class="form_blue_button" :disabled="inProgress" v-if="item.ID !== 0 && useAuthStore().user?.permissions?.includes('Update_Announces')">
                            Update
                        </button>
	            	    <button @click.prevent="Delete" class="form_red_button" :disabled="inProgress" v-if="item.ID !== 0 && useAuthStore().user?.permissions?.includes('Delete_Announces')">
                            Delete
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.main_box {
    @apply flex flex-col mr-10 ml-10 mb-10;
}

.content_placement {
    @apply flex w-full gap-5;
}

.main_content {
    @apply flex flex-col gap-5 h-fit w-full;
}

.control_panel {
    @apply flex flex-col w-full;
}

.update_form {
	@apply flex flex-col gap-5 bg-neutral-100 pt-2 p-5 mb-2 shadow-md rounded-lg;
}

.box_name {
	@apply block text-2xl font-medium self-center;
}

.box_label {
  @apply w-fit block text-sm font-medium text-gray-700;
}

.form_input {
  @apply w-full px-3 py-2 h-96 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500;
}

.control_buttons {
    @apply flex flex-col bg-white rounded-lg shadow-md gap-5 p-5;
}

.form_blue_button {
  @apply w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-400 disabled:cursor-not-allowed;
}

.form_red_button {
  @apply w-full bg-red-500 text-white py-2 px-4 rounded-md hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500 disabled:bg-gray-400 disabled:cursor-not-allowed;
}

</style>