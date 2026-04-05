<script setup lang="ts" generic="T extends Goods | OrderedItem | Order | WorkTime | Announce | User | Role | Permission | Log">
import { imagesAPI, logsAPI } from "@/api"
import { defineEmits } from 'vue'

import type { Goods, Order, OrderedItem, WorkTime, Announce, User, Role, Permission, Log } from "@/types";

const props = defineProps<{
    item: T;
    is_in_window: boolean;
    is_announce: boolean;
}>();

const emit = defineEmits<{
    (e: 'item', value: T): void
}>();

const selectItem = () => {
    emit('item', props.item);
}
    
</script>

<template>
    <!-- Goods -->
    <div class="item_box" v-if="item.Type === 'Goods' || item.Type === 'OrderedItem'">
        <img :src="imagesAPI.getImageSRC() + item.Image" alt="No image" class="item_picture"/>
        <div class="item_text_box">
            <p class="item_name">
                <button v-if="props.is_in_window" @click="selectItem" class="item_name_button">{{ item.Name }}</button>
                <router-link :to="'/goodsedit/' + item.ID" class="item_name_link" v-else >
                    {{ item.Name }}
                </router-link>
            </p>
            <p class="item_discription">
                {{ item.Description && item.Description.length > 100 ? item.Description.slice(0, 100) + '...' :
                    item.Description || "No description available" }}
            </p>
        </div>
        <div class="item_price_box">
            <p class="item_price_numbers">{{ item.Price }}</p>
            <p class="item_price_text">Rub.</p>
        </div>
    </div>

    <!-- Order search -->
    <div class="item_box" v-if="item.Type === 'Order'">
        <div class="item_text_box">
            <p class="item_name">
                <router-link :to="'/ordersedit/' + item.ID + '?q='" class="item_name_link">
                    {{ item.Name }}
                </router-link>
            </p>
            <p class="item_discription">
                {{ item.Phone ? item.Phone : item.Email }}
            </p>
        </div>
    </div>

    <!-- WorkTime search -->
    <div class="item_box" v-if="item.Type === 'WorkTime'">
        <div class="item_text_box">
            <p class="item_name">
                <router-link :to="'/scheduleedit/' + item.ID" class="item_name_link">
                    {{ item.Date }}
                </router-link>
            </p>
            <p class="item_discription">
                {{ item.IsOpened ? (item.TimeStart + "-" + item.TimeEnd) : "Closed" }}
            </p>
        </div>
    </div>

    <!-- Announce search -->
    <div class="item_box" v-if="item.Type === 'Announce' && !is_announce">
        <div class="item_text_box">
            <p class="item_name">
                <router-link :to="'/announcesedit/' + item.ID" class="item_name_link">
                    {{ item.DateTime }}
                </router-link>
            </p>
            <p class="item_discription">
                From: {{ item.From }}
            </p>
        </div>
    </div>

    <!-- Announce -->
    <div class="item_box" v-if="item.Type === 'Announce' && is_announce">
        <div class="announce_box">
            <p class="announce_name ">
                {{ item.DateTime }}
            </p>
            <p class="item_discription">
                From: {{ item.From }}
            </p>
            <p class="announce_text">
                {{ item.Announce }}
            </p>
        </div>
    </div>

    <!-- User search -->
    <div class="item_box" v-if="item.Type === 'User'">
        <div class="item_text_box">
            <p class="item_name">
                <router-link :to="'/usersedit/' + item.ID" class="item_name_link">
                    {{ item.UserName }}
                </router-link>
            </p>
            <p class="item_discription" v-if="item.Role !== undefined">
                Role: {{ item.Role.Name }}
            </p>
        </div>
    </div>

    <!-- Role search -->
    <div class="item_box" v-if="item.Type === 'Role'">
        <div class="item_text_box">
            <p class="item_name">
                <router-link :to="'/rolesedit/' + item.ID + '?q='" class="item_name_link">
                    {{ item.Name }}
                </router-link>
            </p>
            <p class="item_discription">
                Permissions: {{ item.Permissions.map(item => "[" + item.Action + "]").join(" ") }}
            </p>
        </div> 
    </div>

    <!-- Permission search -->
    <div class="item_box" v-if="item.Type === 'Permission'">
        <div class="item_text_box">
            <p class="item_name">
                <button v-if="props.is_in_window" @click="selectItem" class="item_name_button">{{ item.Action }}</button>
                <p v-else>{{ item.Action }}</p>
            </p>
        </div>
    </div>

    <!-- Log search -->
    <div class="item_box" v-if="item.Type === 'Log'">
        <div class="item_text_box">
            <p class="item_name">
                <button @click="logsAPI.getLog(item.Name)" class="item_name_button">{{ item.Name }}</button>
            </p>
        </div>
    </div>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.item_box {
    @apply flex items-center bg-white shadow-md p-4 rounded-lg w-full gap-5;
}

.item_picture {
    @apply w-24 h-24 object-cover rounded-md;
}

.item_text_box {
    @apply ml-4 mr-4 w-11/12 flex flex-col justify-start self-start h-full;
}

.item_name_link {
    @apply text-blue-500 hover:underline;
}

.item_name_button {
    @apply underline hover:text-blue-500 hover:drop-shadow-lg hover:drop-shadow-blue-400;
}

.item_name {
    @apply text-2xl font-bold;
}

.item_discription {
    @apply text-gray-600 text-sm h-full;
}

.item_price_box {
    @apply flex flex-row gap-1 h-full w-1/12 items-center justify-end;
}

.item_price_numbers {
    @apply font-bold text-xl w-fit h-fit;
}

.item_price_text {
    @apply font-normal text-sm w-fit h-fit;
}

.announce_name {
    @apply text-2xl font-bold self-center;
}

.announce_text {
    @apply text-xl w-fit h-fit;
}

.announce_box {
    @apply ml-4 mr-4 w-full flex flex-col justify-start self-start h-full;
}
</style>