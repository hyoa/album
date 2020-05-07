<template>
  <div class="auth flex w-screen h-screen md:justify-center md:items-center">
    <div v-if="disabled" class="rounded w-full md:w-5/6 lg:w-1/3 md:shadow-md bg-white">
      <div class="text-3xl px-4 py-3 text-center">
        {{ $t('auth.unavailableWebsite') }}
      </div>
    </div>
    <div v-else class="rounded w-full md:w-5/6 lg:w-1/3 md:shadow-md bg-white">
      <div class="flex">
        <div
          class="w-1/2 h-12 flex justify-center items-center rounded-tl cursor-pointer bg-primary text-white "
          :class="[ nav === 'login' ? '' : 'opacity-50' ]"
          @click="nav = 'login'"
        >
          <span>{{ $t('auth.connection') }}</span>
        </div>
        <div
          class="w-1/2 h-12 flex justify-center items-center rounded-tr cursor-pointer text-white bg-primary "
          :class="[ nav === 'register' ? '' : 'opacity-50' ]"
          @click="nav = 'register'"
        >
          <span>
            {{ $t('auth.register') }}
          </span>
        </div>
      </div>
      <div class="px-8 pt-6 pb-8 mb-4">
        <Alert v-if="alert.type" :type="alert.type" :message="$t(alert.message)" :title="$t(alert.title)" />
        <Login @onPasswordForget="nav = 'reset-password'" @onRegister="nav = 'register'" @onLogin="onLogin" v-if="nav === 'login'" :formStatus="formStatus.login"/>
        <Register @onRegister="onRegister" v-else-if="nav === 'register'" :formStatus="formStatus.register"/>
        <ResetPassword @onAskResetPassword="onAskResetPassword" v-else-if="nav === 'reset-password'" :formStatus="formStatus.reset"/>
      </div>
    </div>
  </div>
</template>

<script>
import { post } from '../utils/axiosHelper'
import errorHelper from '../utils/errorHelper'

import Login from '../components/auth/Login'
import Register from '../components/auth/Register'
import Alert from '../components/alerts/Alert'
import ResetPassword from '../components/auth/ResetPassword'
export default {
  name: 'Auth',
  components: { ResetPassword, Alert, Register, Login },
  data () {
    return {
      nav: 'login',
      alert: {
        type: null,
        title: null,
        message: null
      },
      formStatus: {
        login: 'ready',
        register: 'ready',
        reset: 'ready'
      },
      disabled: false
    }
  },
  created () {
    if (this.$store.state.flashMessage) {
      this.alert = {
        type: 'info',
        title: 'alert.info.title',
        message: this.$store.state.flashMessage
      }
    }
  },
  methods: {
    async onLogin (data) {
      this.formStatus.login = 'pending'
      this.alert.type = null
      post('user/login', data)
        .then(res => {
          localStorage.setItem('album-token', res.data.token)
          this.$router.push({ name: 'home' })
          this.$store.commit('setToken', res.data.token)
        })
        .catch(({ response }) => {
          let code = null
          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.alert = {
            type: 'error',
            title: 'alert.info.title',
            message: errorHelper(code)
          }
        })
        .finally(() => {
          this.formStatus.login = 'ready'
        })
    },
    onRegister (data) {
      this.alert.type = null
      this.formStatus.register = 'pending'
      post('user/register', data)
        .then(res => {
          this.alert = {
            type: 'success',
            title: 'auth.registerPage.alert.success.title',
            message: 'auth.registerPage.alert.success.message'
          }
        })
        .catch(({ response }) => {
          let code = null
          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.alert = {
            type: 'error',
            title: 'alert.error.title',
            message: errorHelper(code)
          }
        })
        .finally(() => {
          this.formStatus.register = 'ready'
        })
    },
    onAskResetPassword (data) {
      this.alert.type = null
      this.formStatus.reset = 'pending'
      post('user/reset-password/ask', data)
        .then(res => {
          this.alert = {
            type: 'success',
            title: 'auth.askResetPasswordPage.alert.success ',
            message: 'auth.askResetPasswordPage.alert.message'
          }
        })
        .catch(({ response }) => {
          let code = null
          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.alert = {
            type: 'error',
            title: 'auth.alert.error',
            message: errorHelper(code)
          }
        })
        .finally(() => {
          this.formStatus.reset = 'ready'
        })
    }
  }
}
</script>

<style scoped>
  .auth {
    background: url('../../public/img/background-auth.jpg') no-repeat;
    background-size: cover;
  }
</style>
