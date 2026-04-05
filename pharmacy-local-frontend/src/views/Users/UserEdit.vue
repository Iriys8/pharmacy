<script setup lang="ts">
import { onMounted, ref } from "vue";
import type { UserResponse } from "@/types/api";
import { usersAPI } from '@/api';
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores";

const route = useRoute();
const router = useRouter();
const isLoaded = ref(false);
const inProgress = ref(false);
const isPasswordChange = ref(false);

const getItem = async () => {
    item.value = await usersAPI.getUserByID(Number(route.params.id));
	isLoaded.value = true;
}

const item = ref<UserResponse>({
    User: {
        ID: 0,
        Login: "",
        UserName: "",
        RoleID: 0,
        Role: undefined,
        Password: "",
        Type: "User"
    },
    Roles: undefined,
})

onMounted(async () => {
    await getItem();
    if (Number(route.params.id) === 0) {
        isPasswordChange.value = true
    }
});

const Create = async () => {
    if (inProgress.value) return;
    inProgress.value = true;
    try {        
        item.value.User.Role = undefined
        if (item.value.User.Password.length < 4) {
            alert("Error: To short password");
            return;
        }
        await usersAPI.createUser(item.value.User);
        alert("Done.");
		router.push({ path: "/users"});
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
        if (!item.value || item.value.User.ID === 0) {
            alert("Error: Item does not exist");
            return;
        }
        if (item.value.User.Password.length < 4 && isPasswordChange) {
            alert("Error: To short password");
            return;
        }
        await usersAPI.updateUserByID(item.value.User);
        alert("Done.");
		router.push({ path: "/users"});
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
        if (!item.value || item.value.User.ID === 0) {
            alert("Error: Item does not exist");
            return;
        }
		
		const result = confirm("Delete user " + item.value.User.UserName + "?" )
		if (result) {
			await usersAPI.deleteUserByID(item.value.User);
        	alert("Done.");
			router.push({ path: "/users"});
		}	
    }
    catch(err) {
        alert("Error: " + err)
    }
}
</script>

<template>
    <div class="main_box" v-if="useAuthStore().user?.permissions?.includes('Change_Users'), isLoaded">
        <div class="content_placement">
            <div class="main_content">
				<div class="control_panel" v-if="isLoaded">
                    <div class="update_form">
                        <h2 class="form_name">{{ item.User.ID !== 0 ? item.User.UserName : 'Create new user' }}</h2>
                    
	            		<div>
                            <label class="form_label" for="login">Login:</label>
                            <input
                                class="form_input"
                                id="login"
                                v-model="item.User.Login"
                                type="text"
                                required
                            />
                        </div>
                    
                        <div>
                            <label class="form_label" for="username">UserName:</label>
                            <input
                                class="form_input"
                                id="username"
                                v-model="item.User.UserName"
                                type="text"
                                required
                            />
                        </div>

                        <div>
                            <label class="form_label" for="role">Role:</label>
                            <select id="role" v-model="item.User.RoleID" class="form_radio_box_group">
                                <option v-for="role in item.Roles" :value="role.ID">{{ role.Name }}</option>
                            </select>
                        </div>

                        <div v-if="isPasswordChange">
                            <label class="form_label" for="role">Password:</label>
                            <input
                                class="form_input"
                                id="username"
                                v-model="item.User.Password"
                                type="password"
                                required
                            />
                        </div>

                        <div class="form_radio_box_group" v-if="useAuthStore().user?.permissions?.includes('Update_Users') && item.User.ID !== 0">
                            <label class="form_label" for="isopened">Change Password?</label>
                            <input
                                class="form_radio_button"
                                id="isopened"
                                v-model="isPasswordChange"
                                type="checkbox"
                            />
                        </div>
                    
                    	<button @click.prevent="Create" class="form_blue_button" :disabled="inProgress" v-if="item.User.ID === 0">
                            Create
                        </button>
                        <button @click.prevent="Update" class="form_blue_button" :disabled="inProgress" v-if="item.User.ID !== 0">
                            Update
                        </button>
	            		<button @click.prevent="Delete" class="form_red_button" :disabled="inProgress" v-if="item.User.ID !== 0">
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