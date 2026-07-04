import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory('/wol/'),
  routes: [
    {
      path: '/',
      name: 'home-zh',
      component: HomeView,
    },
    {
      path: '/en',
      name: 'home-en',
      component: HomeView,
    },
  ],
})

export default router
