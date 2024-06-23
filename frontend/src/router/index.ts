import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import InvoiceView from '@/views/InvoiceView.vue'
import PettyCashView from '@/views/PettyCashView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/invoice',
      name: 'invoice',
      component: InvoiceView
    },{
     path: '/petty_cash',
      name: 'petty_cash',
      component: PettyCashView
    }
  ]
})

export default router
