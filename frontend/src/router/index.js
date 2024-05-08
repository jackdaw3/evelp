import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import Order from '../views/Order.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/order/:corporationId/:offerId/:itemId',
    name: 'Order',
    component: Order,
    props: (route) => ({
      ...route.params,
      ...route.query
    })
  },
]

const router = new VueRouter({
  mode: 'history',
  routes
})

export default router
