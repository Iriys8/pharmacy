import { useAuthStore } from '@/stores';
import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: () => import('@/views/Home.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'Login', 
      component: () => import('@/views/Login.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/accessdenied',
      name: 'AccessDenied',
      component: () => import('@/views/AccessDenied.vue'),
	  meta: { requiresAuth: true }
    },
	  {
      path: '/goods',
      name: 'Catalog',
      component: () => import('@/views/Goods/Goods.vue'),
      meta: { requiresAuth: true }
    },
	  {
      path: '/goodsedit/:id',
      name: 'Edititem',
      component: () => import('@/views/Goods/GoodEdit.vue'),
	  meta: { requiresAuth: true }
    },
		{
      path: '/schedule',
      name: 'Schedule',
      component: () => import('@/views/Schedule/Schedule.vue'),
      meta: { requiresAuth: true }
    },
		{
      path: '/scheduleedit/:id',
      name: 'EditSchedule',
      component: () => import('@/views/Schedule/ScheduleEdit.vue'),
	  meta: { requiresAuth: true }
    },
    {
      path: '/orders',
      name: 'Orders',
      component: () => import('@/views/Orders/Orders.vue'),
	  meta: { requiresAuth: true }
    },
    {
      path: '/ordersedit/:id',
      name: 'EditOrder',
      component: () => import('@/views/Orders/OrderEdit.vue'),
	  meta: { requiresAuth: true }
    },
    {
      path: '/announces',
      name: 'Announces',
      component: () => import('@/views/Announces/Announces.vue'),
	  meta: { requiresAuth: true }
    },
    {
      path: '/announcesedit/:id',
      name: 'EditAnnounce',
      component: () => import('@/views/Announces/AnnounceEdit.vue'),
	  meta: { requiresAuth: true }
    },
    {
      path: '/users',
      name: 'Users',
      component: () => import('@/views/Users/Users.vue'),
	  meta: { requiresAuth: true }
    },
    {
      path: '/usersedit/:id',
      name: 'EditUsers',
      component: () => import('@/views/Users/UserEdit.vue'),
	  meta: { requiresAuth: true }
    },
    {
      path: '/roles',
      name: 'Roles',
      component: () => import('@/views/Role/Roles.vue'),
	  meta: { requiresAuth: true }
    },
    {
      path: '/rolesedit/:id',
      name: 'EditRole',
      component: () => import('@/views/Role/RoleEdit.vue'),
	  meta: { requiresAuth: true }
    },
    {
      path: '/logs',
      name: 'Logs',
      component: () => import('@/views/Logs.vue'),
	  meta: { requiresAuth: true }
    },
  ],
  scrollBehavior(to, from, savedPosition) {
    return { top: 0 };
  },
});

router.beforeEach(async (to, _, next) => {
  const authStore = useAuthStore();

  if (!authStore.isAppInitialized) {
    await authStore.checkAuth();
  }

  if (to.meta.requiresAuth) {
    if (authStore.isAuthenticated) {
      next();
    } else {
      next('/login');
    }
  } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/');
  } else {
    next();
  }
});

export default router;