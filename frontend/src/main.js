import Vue from 'vue'
import VueLazyload from 'vue-lazyload'
import App from './App.vue'
import router from './router'
import store from './store'
import './registerServiceWorker'
import 'vue-awesome/icons'
// import 'vue-awesome/icons/regular'
import Icon from 'vue-awesome/components/Icon'
import Notifications from 'vue-notification'
import VueI18n from 'vue-i18n'
import messages from './i18n/fr'

Vue.component('v-icon', Icon)

Vue.use(VueLazyload, {
  lazyComponent: true
})

Vue.use(Notifications)

Vue.config.productionTip = false

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requireAuth) && !isTokenValid()) {
    next({
      name: 'auth'
    })
  } else if (
    to.matched.some(record => record.meta.rejectIfAuth) &&
    isTokenValid()
  ) {
    next({
      name: 'home'
    })
  } else {
    next()
  }
})

const isTokenValid = () => {
  return localStorage.getItem('album-token')
}

Vue.filter('truncate', (value, limit) => {
  if (value.length > limit) {
    value = `${value.substring(0, (limit - 3))}...`
  }

  return value
})

Vue.use(VueI18n)
const i18n = new VueI18n({
  locale: 'fr',
  messages
})

new Vue({
  router,
  store,
  i18n,
  render: h => h(App)
}).$mount('#app')
