<script setup lang="ts">
import { onMounted, ref } from "vue";
import GoodsItem from "@/components/GoodsItem.vue";
import type { Goods } from "@/types";
import type { GoodsUpdateRequest } from "@/types/api"
import { goodsAPI } from "@/api"
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores";
import error from "@/components/Error.vue"

const router = useRouter();
const route = useRoute();
const isLoaded = ref(false);
const isExist = ref(true)
const isUpdating = ref(false);

const getItem = async () => {
    item.value = await goodsAPI.getGood(Number(route.params.id));
	if (item.value.ID == 0) {
		isExist.value = false;
		return;
	} 
	isLoaded.value = true;
}

onMounted(async () => {
	getItem()
});

const item = ref<Goods>({
	ID: 0,
    Name: "",
    Image: "",
    Instruction: "",
    Description: "",
    IsPrescriptionNeeded: false,
	IsInStock: true,
    Price: 0,
    Producer: "",
    Tags: undefined,
	Type: "Goods"
})

const Update = async () => {
	if (isUpdating.value) return;
	isUpdating.value = true;
	try {
		if (!item.value || item.value.ID === 0) {
            alert("Error: Item does not exist");
            return;
        }

        const updateData: GoodsUpdateRequest = {
            ID: item.value.ID,
            Name: item.value.Name,
            Instruction: item.value.Instruction,
            Description: item.value.Description,
            Prescription: item.value.IsPrescriptionNeeded,
            IsInStock: item.value.IsInStock,
            Price: item.value.Price,
        };

        await goodsAPI.updateGood(updateData);
        alert("Done.");
		router.push({ path: "/announces"});
	}
	catch(err) {
		alert("Error: " + err)
	}
	finally {
        isUpdating.value = false;
    }
}
</script>

<template>
    <div class="main_box">
        <div class="content_placement">
            <div class="main_content" v-if="isExist">
				<div class="control_panel" v-if="useAuthStore().user?.permissions?.includes('Update_Goods')">
					<form class="update_form" v-if="isLoaded">
      					<h2 class="form_name">Data</h2>
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
					
      					<div>
      					  <label class="form_label" for="price">Price:</label>
      					  <input
						  	class="form_input"
      					    id="price"
      					    v-model="item.Price"
      					    type="number"
      					    required
      					  />
      					</div>

						<div>
      					  <label class="form_label" for="description">Description:</label>
      					  <input
							class="form_input"
      					    id="description"
      					    v-model="item.Description"
      					    type="text"
      					    required
      					  />
      					</div>

						<div>
      					  <label class="form_label" for="instruction">Instruction:</label>
      					  <input
							class="form_input"
      					    id="description"
      					    v-model="item.Instruction"
      					    type="text"
      					    required
      					  />
      					</div>

						<div class="form_radio_box_group">
      					  <label class="form_label" for="prescription">Prescription?</label>
      					  <input
							class="form_radio_button"
      					    id="prescription"
      					    v-model="item.IsPrescriptionNeeded"
      					    type="checkbox"
      					  />
      					</div>

						<div class="form_radio_box_group">
      					  <label class="form_label" for="instock">In stock?</label>
      					  <input
							class="form_radio_button"
      					    id="instock"
      					    v-model="item.IsInStock"
      					    type="checkbox"
      					  />
      					</div>
					
      					<button @click="Update" class="form_update_button" :disabled="!isLoaded || isUpdating" v-if=" useAuthStore().user?.permissions?.includes('Update_Goods')">
							{{ isUpdating ? "Updating..." : "Update" }}
      					</button>
    				</form>
            	</div>
                <div class="item_box">
                    <GoodsItem v-if="isLoaded" :product="item" />
					<p v-else>Loading...</p>
                </div>
            </div>
			<div class="main_content" v-else>
				<error/>
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
    @apply flex flex-row gap-5 h-fit w-full justify-center;
}

.control_panel {
    @apply flex flex-col w-1/2 bg-neutral-100 pt-2 p-5 shadow-md rounded-lg;
}

.item_box {
    @apply flex flex-col w-1/2 bg-neutral-100 pt-2 p-5 shadow-md rounded-lg;
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

.form_update_button {
  @apply w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-400 disabled:cursor-not-allowed;
}

</style>