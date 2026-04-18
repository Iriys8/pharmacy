<script setup lang="ts">
import { onMounted, ref } from "vue";
import type { Goods, Order, OrderedItem } from "@/types";
import { orderAPI } from '@/api';
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores";
import SearchWindow from "@/components/SearchWindow.vue";
import MiniItem from "@/components/MiniItem.vue";

const route = useRoute();
const isLoaded = ref(false);

const getItem = async () => {
	if (route.params.id != '0') {
    	item.value = await orderAPI.getOrderByID(Number(route.params.id));
	}
	isLoaded.value = true;
}

onMounted(async () => {
    await getItem();
});

const item = ref<Order>({
    ID: 0,
    Name: "",
    Email: "",
    Phone: "",
    Items: [],
    Type: "Order"
})

const inProgress = ref(false);
const router = useRouter();

const AddQuantity = (ordereditem: OrderedItem) => {
    ordereditem.Quantity++;
}

const RemoveQuantity = (ordereditem: OrderedItem) => {
    ordereditem.Quantity--;
    if (ordereditem.Quantity <= 0) {
        if (item.value.Items !== undefined) {
            const index = item.value.Items.indexOf(ordereditem);
            if (index !== -1) {
                item.value.Items?.splice(index, 1);
            }
        }
    }
}

const selectedItem = (value: Goods) => {
    if (item.value.Items !== undefined) {
        const newItem = item.value.Items.find(item => item.ID === value.ID);
        if (!newItem) {
            item.value.Items.push(<OrderedItem>{
                ID: value.ID,
                Name: value.Name,
                Image: value.Image,
                Description: value.Description,
                Price: value.Price,
                Quantity: 1,
                Type: "OrderedItem"
            })
        }
        else {
            newItem.Quantity++;
        }
    }
}

const Create = async () => {
    if (inProgress.value) return;
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!item.value.Phone || item.value.Phone.trim().length === 0 || !item.value.Name || item.value.Name.trim().length === 0 || !item.value.Items || item.value.Items.length === 0) {
        alert("Error: Missing one of required fields");
        return;
    }
    if (!(item.value.Email !== undefined && (emailRegex.test(item.value.Email) || item.value.Email !== ""))) {
        alert("Error: Wrong email format");
        return;
    }
    inProgress.value = true;
    try {        
        await orderAPI.createOrder(item.value);
        alert("Done.");
		router.push({ path: "/orders"});
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
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!item.value.Phone || item.value.Phone.trim().length === 0 || !item.value.Name || item.value.Name.trim().length === 0 || !item.value.Items || item.value.Items.length === 0) {
        alert("Error: Missing one of required fields");
        return;
    }
    if (!(item.value.Email !== undefined && (emailRegex.test(item.value.Email) || item.value.Email !== ""))) {
        alert("Error: Wrong email format");
        return;
    }
    inProgress.value = true;
    try {
        if (!item || item.value.ID === 0) {
            alert("Error: Item does not exist");
            return;
        }
   
        await orderAPI.updateOrderByID(item.value);
        alert("Done.");
		router.push({ path: "/orders"});
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
			await orderAPI.deleteOrderByID(item.value);
        	alert("Done.");
			router.push({ path: "/orders"});
		}	
    }
    catch(err) {
        alert("Error: " + err)
    }
}
</script>

<template>
    <div class="main_box" v-if="useAuthStore().user?.permissions?.includes('Read_Orders')">
        <div class="content_placement">
            <div class="main_content">
				<div class="control_panel" v-if="(useAuthStore().user?.permissions?.includes('Read_Orders') && item.ID !== 0) || (useAuthStore().user?.permissions?.includes('Create_Orders') && item.ID === 0)">
                    <div class="update_form">
                        <h2 class="box_name">{{ item.ID !== 0 ? `Order from ${item.Name}` : 'Create new order' }}</h2>
                    
	            		<div v-if="useAuthStore().user?.permissions?.includes('Update_Orders') || (useAuthStore().user?.permissions?.includes('Create_Orders') && item.ID === 0)">
                            <label class="form_label" for="name">Name:</label>
                            <input
                                class="form_input"
                                id="name"
                                v-model="item.Name"
                                type="text"
                                required
                            />
                        </div>
                        <div v-else>
                            <label class="form_label">Name:</label>
                            <p class="form_input">{{ item.Name }}</p>
                        </div>
                    
                    
                        <div v-if="useAuthStore().user?.permissions?.includes('Update_Orders') || (useAuthStore().user?.permissions?.includes('Create_Orders') && item.ID === 0)">
                            <label class="form_label" for="email">Email:</label>
                            <input
                                class="form_input"
                                id="email"
                                v-model="item.Email"
                                type="email"
                                required
                            />
                        </div>
                        <div v-else>
                            <label class="form_label">Email:</label>
                            <p class="form_input">{{ item.Email }}</p>
                        </div>
                    
                        <div v-if="useAuthStore().user?.permissions?.includes('Update_Orders') || (useAuthStore().user?.permissions?.includes('Create_Orders') && item.ID === 0)">
                            <label class="form_label" for="phone">Phone:</label>
                            <input
                                class="form_input"
                                id="phone"
                                v-model="item.Phone"
                                type="tel"
                                required
                            />
                        </div>
                        <div v-else>
                            <label class="form_label">Phone:</label>    
                            <p class="form_input">{{ item.Phone }}</p>
                        </div>
                    </div>
                
                    <div class="order_items_box">
                        <div class="order_item_placement">
                            <div v-for="orderedItem in item.Items" class="items_box" v-if="item.Items !== undefined && item.Items.length !== 0">
                                <MiniItem :item="orderedItem" :is_in_window="false" :is_announce="false" class="order_item"/>
                                <div class="item_info_box">
                                    <div class="item_info_sub_box">
                                        <p>Quantity: </p>
                                        <p class="item_info_value">{{ orderedItem.Quantity }}</p>
                                    </div>
                                    <div class="item_info_sub_box">
                                        <p>Price: </p>
                                        <p class="item_info_value">{{ orderedItem.Quantity * orderedItem.Price }}</p>
                                        <p> Rub.</p>
                                    </div>
                                </div>
                                <div class="items_buttons_box" v-if="useAuthStore().user?.permissions?.includes('Update_Orders') || (useAuthStore().user?.permissions?.includes('Create_Orders') && item.ID === 0)">
                                    <button class="form_blue_button" @click="AddQuantity(orderedItem)">Add</button>
                                    <button class="form_red_button" @click="RemoveQuantity(orderedItem)">Remove</button>
                                </div>
                            </div>
                            <div v-else class="no_results">
                                <img src="@/assets/Error.png" alt="Error" class="no_results_picture"/>
                                <p class="no_results_text">Nothing found</p>
                            </div>
                        </div>
                        <div class="order_add_item_box" v-if="useAuthStore().user?.permissions?.includes('Update_Orders') || (useAuthStore().user?.permissions?.includes('Create_Orders') && item.ID === 0)">
                            <div class="order_add_catalog_box">
                                <SearchWindow :is_advert_box="false" :is_search_line="true" :is_enabled="true" :is_search_window="true" :page_addres="'/ordersedit/' + item.ID" @item="value => selectedItem(value as Goods)" type="Goods" needed_access_type="Goods"/>
                            </div>
                        </div>
                    </div>
                
                    <div class="control_buttons" v-if="useAuthStore().user?.permissions?.includes('Create_Orders') || useAuthStore().user?.permissions?.includes('Update_Orders') || useAuthStore().user?.permissions?.includes('Delete_Orders')">
                        <button @click.prevent="Create" class="form_blue_button" :disabled="inProgress" v-if="item.ID === 0 && useAuthStore().user?.permissions?.includes('Create_Orders')">
                            Create
                        </button>
                        <button @click.prevent="Update" class="form_blue_button" :disabled="inProgress" v-if="item.ID !== 0 && useAuthStore().user?.permissions?.includes('Update_Orders')">
                            Update
                        </button>
	            	    <button @click.prevent="Delete" class="form_red_button" :disabled="inProgress" v-if="item.ID !== 0 && useAuthStore().user?.permissions?.includes('Delete_Orders')">
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
  @apply w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500;
}

.order_items_box{
    @apply flex flex-row h-fit gap-5 bg-neutral-100 p-5 shadow-md rounded-lg mb-2 justify-center;
}

.order_item_placement {
    @apply flex flex-col gap-5 overflow-y-auto max-h-128 p-5 bg-neutral-200 rounded-lg w-1/2;
}

.items_box {
    @apply flex flex-row gap-2;
}

.order_item {
    @apply w-9/12;
}

.item_info_box {
    @apply flex flex-col bg-white rounded-lg shadow-md gap-5 justify-center items-center w-3/12;
}

.item_info_sub_box {
    @apply flex flex-row gap-2 items-center;
}

.items_buttons_box {
    @apply flex flex-col justify-around p-5 bg-white rounded-lg;
}

.item_info_value {
    @apply text-xl font-bold self-center;
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

.order_add_item_box {
    @apply w-1/2 h-128 gap-5 p-5 bg-neutral-200 rounded-lg;
}

.order_add_catalog_box {
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