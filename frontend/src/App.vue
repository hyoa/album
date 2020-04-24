<template>
  <div id="app" class="min-h-screen relative" :class="{ 'max-h-screen overflow-hidden' : isSideBarOpen }">
    <router-view/>
    <div
      v-if="reloadApplication"
      @click="reload"
      class="sticky bottom-0 right-0 bg-primary p-3 text-white md:rounded-tl"
    >
      L'application a été mis à jour ! Cliquez pour raffraichir
    </div>
    <notifications group="success" position="bottom left" classes="success" />
    <notifications group="error" position="bottom left" classes="error" />
    <notifications group="info" position="bottom left" classes="info" />
    <notifications group="warning" position="bottom left" classes="warning" />
  </div>
</template>

<style lang="scss">
@import 'assets/styles/tailwind.postcss';
body {
  font-family: 'Source Sans Pro', sans-serif;
}
.font-manuscript {
  font-family: "Fredericka the Great", cursive;
}
.notifications {
  width: 100% !important;

  @screen md {
    width: 300px !important;
  }

  .vue-notification-template {
    @apply p-3 bg-blue-400 text-white border-blue-500 border-l-4;
  }

  .success {
    @apply bg-green-400 text-white border-green-500;
  }

  .warning {
    @apply bg-orange-400 text-white border-orange-500;
  }

  .error {
    @apply bg-red-400 text-white border-red-500;
  }
}
</style>

<script>
import firebase from 'firebase/app'
import 'firebase/messaging'
import { post } from './utils/axiosHelper'

export default {
  data () {
    return {
      tokenDecoded: null,
      intervalId: null,
      reloadApplication: false
    }
  },
  async created () {
    const token = localStorage.getItem('album-token')
    document.title = process.env.VUE_APP_SITE_TITLE

    if (token) {
      this.$store.commit('setToken', token)
    }

    document.addEventListener('swUpdated', function () {
      this.reloadApplication = true
    }.bind(this))

    if (localStorage.getItem('declineNotification') === null) {
      firebase.initializeApp({
        messagingSenderId: process.env.VUE_APP_FIREBASE_SENDER_ID
      })

      if (Notification.permission === 'granted') {
        const messaging = firebase.messaging()
        await messaging.requestPermission()
        const tokenMessaging = await messaging.getToken()

        await post('notification/subscribe', { token: tokenMessaging, channel: 'album' })

        if (this.$store.state.token.role === 9) {
          await post('notification/subscribe', { token: tokenMessaging, channel: 'admin' })
        }
      }
    }
  },
  watch: {
    token (newToken, oldToken) {
      if (oldToken === null && newToken !== null) {
        this.watchTokenValidity()
      }
    }
  },
  methods: {
    watchTokenValidity () {
      this.intervalId = setInterval(() => {
        const currentToken = this.$store.state.token

        if (currentToken.exp * 1000 < (new Date()).getTime()) {
          // this.$store.commit('setToken', null)
          this.$store.commit('setFlashMessage', 'Par mesure de sécurité, vous avez été déconnecté. Vous pouvez vous reconnecter avec le formulaire ci-dessous.')
          localStorage.removeItem('album-token')

          if (localStorage.getItem('album-token') === null) {
            this.$router.push({ name: 'auth' })
            clearInterval(this.intervalId)
          }
        }
      }, 1000)
    },
    reload () {
      location.reload(true)
    }
  },
  computed: {
    token () {
      return this.$store.state.token
    },
    isSideBarOpen () {
      return this.$store.state.sideBarIsOpen
    }
  }
}
</script>
