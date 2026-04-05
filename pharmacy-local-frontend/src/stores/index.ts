import { defineStore } from 'pinia';
import { ref, computed, onMounted } from 'vue';
import { type AxiosResponse, type InternalAxiosRequestConfig } from 'axios';
import type { User, LoginData, AuthResponse } from '@/types/auth';
import { api, authAPI } from "@/api"
import { useRouter } from 'vue-router';

export const useAuthStore = defineStore('auth', () => {
  const error = ref<string | null>(null);
  const accessToken = ref<string | null>(null);
  const router = useRouter()
  
  const user = ref<User | null>(null);
  const isLoading = ref(false);
  const isAppInitialized = ref(false);
  const isAuthenticated = computed(() => !!(user.value && accessToken.value));

  const checkAuth = async () => {      
      isLoading.value = true;
      try {
        await refreshToken();
        return true;
      } catch (err) {
        accessToken.value = null;
        user.value = null;
        return false;
      } finally {
        isLoading.value = false;
        isAppInitialized.value = true;
      }
  };
  
  const refreshToken = async (): Promise<AuthResponse> => {
    const response = await authAPI.refreshToken()
    user.value = response.user;
	accessToken.value =	response.access_token;
    return response;
  };
  
  const login = async (loginData: LoginData): Promise<AuthResponse> => {
    isLoading.value = true;
    error.value = null;
    
    try {
      const response = await authAPI.login(loginData);
      accessToken.value = response.access_token;
      user.value = response.user;
      return response;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Login failed';
      throw err;
    } finally {
      isLoading.value = false;
    }
  };
  
  const logout = async (): Promise<void> => {
    try {
      await authAPI.logout()
    } catch (err) {
      console.error('Logout error:', err);
    } finally {
      accessToken.value = null;
      user.value = null;
      isAppInitialized.value = false;
    }
  };
  
  onMounted(() => {

    api.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        if (accessToken.value) {
          config.headers.Authorization = `Bearer ${accessToken.value}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    api.interceptors.response.use(
      (response: AxiosResponse) => response,
      async (error) => {
        const originalRequest = error.config;
        
        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;
          
          try {
            await refreshToken();
            return api(originalRequest);
          } catch (refreshError) {
            await logout();
            router.push('/login');
            return Promise.reject(refreshError);
          }
        }

        if (error.response?.status === 403) {
          router.push('/accessdenied')
        }
        
        return Promise.reject(error);
      }
    );
	
  });
  return{
	user,
	isLoading,
	isAppInitialized,
	isAuthenticated,
	checkAuth,
	refreshToken,
	login,
	logout,
  }
});
