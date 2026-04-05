<script setup lang="ts">
import { RouterLink, useRouter} from 'vue-router';
import { useAuthStore } from "@/stores"
import logo from "@/components/Logo.vue";

const router = useRouter()
const auth = useAuthStore()
const username = auth.user?.username

const Logout = async () => {
	try {
		await auth.logout()
		router.push("/login")
	}
	finally {
		console.log("Logout")
	}
}
</script>

<template>
    <header class="header">
        <div class="header_box">
            <nav class="nav_links">
				<RouterLink class="nav_link" to="/">
            		<logo/>
        		</RouterLink>
				<RouterLink to="/goods" class="nav_link">Catalog</RouterLink>
                <RouterLink to="/schedule" class="nav_link">Schedule</RouterLink>
				<RouterLink to="/orders" class="nav_link" v-if="useAuthStore().user?.permissions?.includes('Read_Orders')">Orders</RouterLink>
				<RouterLink to="/announces" class="nav_link">Announces</RouterLink>
				<RouterLink to="/users" class="nav_link" v-if="useAuthStore().user?.permissions?.includes('Change_Users')">Users</RouterLink>
				<RouterLink to="/roles" class="nav_link" v-if="useAuthStore().user?.permissions?.includes('Change_Roles')">Roles</RouterLink>
				<RouterLink to="/logs" class="nav_link" v-if="useAuthStore().user?.permissions?.includes('Download_Logs')">Logs</RouterLink>
            </nav>
			<div class="user_info">
				<p class="user_name">{{ username }}</p>
				<button class="logout_button" @click="Logout">Logout</button>
			</div>
        </div>
    </header>
</template>

<style scoped lang="postcss">
@import "tailwindcss";

.header {
    @apply flex flex-col;
}

.header_box {
    @apply flex justify-between h-12 w-full bg-neutral-200 pl-5 pr-5 mb-5 shadow-md;
}

.nav_links {
    @apply flex justify-center gap-6 place-self-center;
}

.nav_link {
    @apply text-xl font-medium hover:text-blue-500 hover:drop-shadow-lg hover:drop-shadow-blue-400 place-self-center;
}

.logout_button {
	@apply w-full text-xl font-medium place-self-center bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500
}

.user_name {
	@apply text-xl font-medium place-self-center;
}

.user_info {
	@apply flex gap-6 justify-center p-2;
}

</style>