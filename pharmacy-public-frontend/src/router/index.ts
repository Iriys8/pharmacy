import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/Home.vue'),
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/About.vue'),
    },
    {
      path: '/catalog',
      name: 'catalog',
      component: () => import('../views/Сatalog.vue'),
    },
    {
      path: '/schedule',
      name: 'schedule',
      component: () => import('../views/Schedule.vue'),
    },
    {
      path: '/shoppingcart',
      name: 'shoppingcart',
      component: () => import('../views/ShoppingCart.vue'),
    },
    {
      path: '/item/:id',
      name: 'item',
      component: () => import('../views/ItemPage.vue'),
    },
  ],
  scrollBehavior(to, from, savedPosition) {
    return { top: 0 };
  },
})

export default router
