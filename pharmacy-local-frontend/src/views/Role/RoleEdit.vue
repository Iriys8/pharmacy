<script setup lang="ts">
import { onMounted, ref } from "vue";
import type { Goods, Role, Permission } from "@/types";
import { rolesAPI } from '@/api';
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores";
import SearchWindow from "@/components/SearchWindow.vue";
import MiniItem from "@/components/MiniItem.vue";

const route = useRoute();
const isLoaded = ref(false);

const getItem = async () => {
	if (route.params.id != '0') {
    	item.value = await rolesAPI.getRoleByID(Number(route.params.id));
	}
	isLoaded.value = true;
}

onMounted(async () => {
    await getItem();
});

const item = ref<Role>({
    ID: 0,
    Name: "",
    Permissions: [],
    Type: "Role"
})

const inProgress = ref(false);
const router = useRouter();

const removePermission = (permission: Permission) => {
    if (item.value.Permissions !== undefined) {
        const index = item.value.Permissions.indexOf(permission);
        if (index !== -1) {
            item.value.Permissions.splice(index, 1);
        }
    }
}

const selectedItem = (value: Permission) => {
    if (item.value.Permissions !== undefined) {
        const newItem = item.value.Permissions.find(item => item.ID === value.ID);
        if (!newItem) {
            item.value.Permissions.push(<Permission>{
                ID: value.ID,
                Action: value.Action,
                Type: "Permission"
            })
        }
        else {
            removePermission(newItem)
        }
    }
}

const Create = async () => {
    if (inProgress.value) return;
    if (!item.value.Name || item.value.Name.trim().length === 0) {
        alert("Error: Missing one of required fields");
        return;
    }
    inProgress.value = true;
    try {        
        await rolesAPI.createRole(item.value);
        alert("Done.");
		router.push({ path: "/roles"});
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
    if (!item.value.Name || item.value.Name.trim().length === 0) {
        alert("Error: Missing one of required fields");
        return;
    }
    inProgress.value = true;
    try {
        if (!item || item.value.ID === 0) {
            alert("Error: Item does not exist");
            return;
        }
   
        await rolesAPI.updateRoleByID(item.value);
        alert("Done.");
		router.push({ path: "/roles"});
    }
    catch(err) {
        alert("Error: " + err);
    }
    finally {
        inProgress.value = false;
    }
}

const Delete = async () => {
	if (inProgress.value) return;
	try {
        if (!item || item.value.ID === 0) {
            alert("Error: Item does not exist");
            return;
        }
		
		const result = confirm("Delete order from " + item.value.Name + "?" );
		if (result) {
			await rolesAPI.deleteRoleByID(item.value);
        	alert("Done.");
			router.push({ path: "/roles"});
		}	
    }
    catch(err) {
        alert("Error: " + err)
    }
}
</script>

<template>
    <div class="main_box" v-if="useAuthStore().user?.permissions?.includes('Change_Roles')">
        <div class="content_placement">
            <div class="main_content">
				<div class="control_panel" v-if="isLoaded">
                    <div class="update_form">
                        <h2 class="box_name">{{ item.ID !== 0 ? item.Name : 'Create new role' }}</h2>
                    
	            		<div>
                            <label class="form_label" for="name">Name:</label>
                            <input
                                class="form_input"
                                id="name"
                                v-model="item.Name"
                                type="text"
                                required
                            />
                        </div>
                
                        <div class="permission_items_box">
                            <div class="permission_item_placement">
                                <div v-for="permission in item.Permissions" class="items_box" v-if="item.Permissions !== undefined && item.Permissions.length !== 0">
                                    <MiniItem :item="permission" :is_in_window="false" :is_announce="false" class="permission_item"/>
                                    <div class="items_buttons_box">
                                        <button class="form_red_button" @click="removePermission(permission )">Remove</button>
                                    </div>
                                </div>
                                <div v-else class="no_results">
                                    <img src="@/assets/Error.png" alt="Error" class="no_results_picture"/>
                                    <p class="no_results_text">Nothing found</p>
                                </div>
                            </div>
                            <div class="permission_add_item_box">
                                <div class="permission_add_catalog_box">
                                    <SearchWindow :is_advert_box="false" :is_search_line="true" :is_enabled="true" :is_search_window="true" :page_addres="'/rolesedit/' + item.ID" @item="value => selectedItem(value as Permission)" type="Permission"/>
                                </div>
                            </div>
                        </div>
                
                        <div class="control_buttons">
                            <button @click.prevent="Create" class="form_blue_button" :disabled="inProgress" v-if="item.ID === 0">
                                Create
                            </button>
                            <button @click.prevent="Update" class="form_blue_button" :disabled="inProgress" v-if="item.ID !== 0">
                                Update
                            </button>
	            	        <button @click.prevent="Delete" class="form_red_button" :disabled="inProgress" v-if="item.ID !== 0">
                                Delete
                            </button>
                        </div>
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
  @apply w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500;
}

.permission_items_box{
    @apply flex flex-row h-fit gap-5 bg-neutral-100 p-5 shadow-md rounded-lg mb-2 justify-center;
}

.permission_item_placement {
    @apply flex flex-col gap-5 overflow-y-auto max-h-128 p-5 bg-neutral-200 rounded-lg w-1/2;
}

.items_box {
    @apply flex flex-row gap-2;
}

.permission_item {
    @apply w-9/12;
}

.items_buttons_box {
    @apply flex flex-col justify-around p-5 bg-white rounded-lg;
}

.no_results {
    @apply flex flex-col border rounded-lg pt-2 p-5;
}

.no_results_picture {
    @apply place-self-center h-24 w-24;
}

.no_results_text {
    @apply place-self-center text-2xl;
}

.permission_add_item_box {
    @apply w-1/2 h-128 gap-5 p-5 bg-neutral-200 rounded-lg;
}

.permission_add_catalog_box {
    @apply h-full overflow-y-auto;
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