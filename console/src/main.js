import Vue from 'vue'
import App from './App.vue'
import router from './router'
import BootstrapVue from 'bootstrap-vue'
import Component from 'vue-class-component'
import { setDefaults } from './api'

import './assets/css/bootstrap.min.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import './assets/css/icons.min.css'
import './assets/css/app.min.css'

Vue.use(BootstrapVue)

Vue.config.productionTip = false

Component.registerHooks([
  'beforeRouteEnter',
  'beforeRouteLeave',
  'beforeRouteUpdate'
])

const vue = new Vue({
  router,
  render: h => h(App)
})

setDefaults(defaults => {
  defaults.prefix = '/api'
  defaults.fallback = (error, resolve, reject) => {
    if (error?.status === 401 || error?.response?.status === 401) {
      if (vue.$route.name !== 'authorize') {
        vue.$router.push('/authorize')
        resolve(true)
      } else {
        reject(error)
      }
    } else {
      reject(error)
    }
  }
})

window.vue = vue
vue.$mount('#app')
export { vue }
