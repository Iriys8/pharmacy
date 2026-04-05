<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import type { WorkTime } from "@/types";
import { scheduleAPI } from '@/api';
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores";

const route = useRoute();
const isLoaded = ref(false);

const getItem = async () => {
	if (route.params.id != '0') {
    	item.value = await scheduleAPI.getScheduleByID(Number(route.params.id));
	}
	isLoaded.value = true;
}

onMounted(async () => {
    await getItem();
});

const item = ref<WorkTime>({
    ID: 0,
    Date: "",
    TimeStart: "",
    TimeEnd: "",
    IsOpened: false,
    Type: "WorkTime"
})

const inProgress = ref(false);
const router = useRouter();

watch(() => item.value.IsOpened, (newValue) => {
    if (!newValue) {
        item.value.TimeStart = "";
        item.value.TimeEnd = "";
    }
});

const Create = async () => {
    if (inProgress.value) return;
    inProgress.value = true;
    try {
        if (!item.value.IsOpened) {
            item.value.TimeStart = null;
            item.value.TimeEnd = null;
        }
        
        await scheduleAPI.createSchedule(item.value);
        alert("Done.");
		router.push({ path: "/schedule"});
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
    inProgress.value = true;
    try {
        if (!item.value || item.value.ID === 0) {
            alert("Error: Item does not exist");
            return;
        }
        
        if (!item.value.IsOpened) {
            item.value.TimeStart = null;
            item.value.TimeEnd = null;
        }
        
        await scheduleAPI.updateScheduleByID(item.value);
        alert("Done.");
		router.push({ path: "/schedule"});
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
		
		const result = confirm("Delete shedule for " + item.value.Date + "?" )
		if (result) {
			await scheduleAPI.deleteScheduleByID(item.value);
        	alert("Done.");
			router.push({ path: "/schedule"});
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
					<div class="control_panel" v-if="item.ID !== 0 || (useAuthStore().user?.permissions?.includes('Create_WorkTime') && item.ID === 0)">
                        <div class="update_form">
                            <h2 class="form_name">{{ item.ID !== 0 ? item.Date : 'Create new special shedule' }}</h2>
                        
	                		<div v-if="useAuthStore().user?.permissions?.includes('Update_WorkTime') || (useAuthStore().user?.permissions?.includes('Create_WorkTime') && item.ID === 0)">
                                <label class="form_label" for="dateat">Date:</label>
                                <input
                                    class="form_input"
                                    id="dateat"
                                    v-model="item.Date"
                                    type="date"
                                    required
                                />
                            </div>
                            <div v-else>
                                <label class="form_label">Date:</label>
                                <p class="form_input">{{ item.Date }}</p>
                            </div>
                        
                            <div v-if="useAuthStore().user?.permissions?.includes('Update_WorkTime') || (useAuthStore().user?.permissions?.includes('Create_WorkTime') && item.ID === 0)">
                                <label class="form_label" for="timestart">Open time:</label>
                                <input
                                    class="form_input"
                                    id="timestart"
                                    v-model="item.TimeStart"
                                    type="time"
                                    :disabled="!item.IsOpened"
                                    required
                                />
                            </div>
                            <div v-else>
                                <label class="form_label">Open time:</label>
                                <p class="form_input">{{ item.TimeStart }}</p>
                            </div>
                        
                            <div v-if="useAuthStore().user?.permissions?.includes('Update_WorkTime') || (useAuthStore().user?.permissions?.includes('Create_WorkTime') && item.ID === 0)">
                                <label class="form_label" for="timeend">Close time:</label>
                                <input
                                    class="form_input"
                                    id="timeend"
                                    v-model="item.TimeEnd"
                                    type="time"
                                    :disabled="!item.IsOpened"
                                    required
                                />
                            </div>
                            <div v-else>
                                <label class="form_label">Close time:</label>
                                <p class="form_input">{{ item.TimeEnd }}</p>
                            </div>
                        
                            <div class="form_radio_box_group" v-if="useAuthStore().user?.permissions?.includes('Update_WorkTime') || (useAuthStore().user?.permissions?.includes('Create_WorkTime') && item.ID === 0)">
                                <label class="form_label" for="isopened">Is opened?</label>
                                <input
                                    class="form_radio_button"
                                    id="isopened"
                                    v-model="item.IsOpened"
                                    type="checkbox"
                                />
                            </div>
                        
                        	<button @click.prevent="Create" class="form_blue_button" :disabled="inProgress" v-if="useAuthStore().user?.permissions?.includes('Create_WorkTime') && item.ID === 0">
                                Create
                            </button>
                            <button @click.prevent="Update" class="form_blue_button" :disabled="inProgress" v-if="useAuthStore().user?.permissions?.includes('Update_WorkTime') && item.ID !== 0">
                                Update
                            </button>
	                		<button @click.prevent="Delete" class="form_red_button" :disabled="inProgress" v-if="useAuthStore().user?.permissions?.includes('Delete_WorkTime') && item.ID !== 0">
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
    @apply flex flex-col w-full bg-neutral-100 pt-2 p-5 shadow-md rounded-lg;
}

.update_form {
	@apply flex flex-col gap-5;
}

.form_name {
	@apply block text-2xl font-medium self-center;
}

.form_label {
  @apply w-fit block text-sm font-medium text-gray-700;
}

.form_input {
  @apply w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500;
}

.form_radio_box_group {
  @apply flex flex-row w-fit border border-gray-300 rounded-md px-3 py-2; 
}

.form_radio_button {
	@apply w-fit h-fit self-center ml-5;
}

.form_blue_button {
  @apply w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-400 disabled:cursor-not-allowed;
}

.form_red_button {
  @apply w-full bg-red-500 text-white py-2 px-4 rounded-md hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500 disabled:bg-gray-400 disabled:cursor-not-allowed;
}

</style>