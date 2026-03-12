import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue'),
    },
    {
      path: '/stocks',
      name: 'stocks',
      component: () => import('@/views/StocksView.vue'),
    },
    {
      path: '/recommendations',
      name: 'recommendations',
      component: () => import('@/views/RecommendationsView.vue'),
    },
  ],
})

export default router
