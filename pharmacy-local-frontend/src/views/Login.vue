<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from "@/stores"
import type { LoginData } from "@/types/auth";

const router = useRouter()
const auth = useAuthStore();

const Login = async () => {
	try {
		await auth.login(loginData.value);
		router.push("/");
	}
	catch (error) {
    loginData.value.login = "";
	loginData.value.password = ""
  }
}

const loginData = ref<LoginData>({
	login: '',
	password: ''
});
</script>

<template>
    <div class="main_box">
        <div class="content_placement">
			<form @submit.prevent="Login" class="login_form">
      			<h2 class="form_name">Login</h2>
      			<div class="form_group">
      			  <label class="form_label" for="username">Username:</label>
      			  <input
					class="form_input"
      			    id="username"
      			    v-model="loginData.login"
      			    type="text"
      			    required
      			  />
      			</div>
			
      			<div class="form_group">
      			  <label class="form_label" for="password">Password:</label>
      			  <input
				    class="form_input"
      			    id="password"
      			    v-model="loginData.password"
      			    type="password"
      			    required
      			  />
      			</div>
			
      			<button type="submit" class="button" :disabled="auth.isLoading">
      			  {{ auth.isLoading ? 'Logging in...' : 'Login' }}
      			</button>
    		</form>
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

.login_form {
  @apply flex flex-col max-w-md mx-auto mt-8 p-6 bg-white rounded-lg shadow-md ;
}

.form_name {
	@apply block text-xl font-medium text-gray-700 mb-2 self-center;
}

.form_group {
  @apply mb-4;
}

.form_label {
  @apply block text-sm font-medium text-gray-700 mb-2;
}

.form_input {
  @apply w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500;
}

.button {
  @apply w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-400 disabled:cursor-not-allowed;
}

.error-message {
  @apply mt-4 p-2 bg-red-100 border border-red-400 text-red-700 rounded-md;
}
</style>