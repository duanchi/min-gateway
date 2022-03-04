import Vue from 'vue'
import VueRouter from 'vue-router'
import Services from '../components/services'
import Routes from '../components/routes'
import Authorize from '../components/authorize'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: '/routes'
  },
  {
    path: '/authorize',
    name: 'authorize',
    component: Authorize
  },
  {
    path: '/routes',
    name: 'routes',
    component: Routes
  },
  {
    path: '/services',
    name: 'services',
    component: Services
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    // component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
  linkActiveClass: 'active',
  linkExactActiveClass: 'active'
})

export default router
