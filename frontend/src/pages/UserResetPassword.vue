<template>
  <div class="flex w-screen h-screen justify-center items-center">
    <div class="rounded w-1/3 shadow-md">
      <div class="px-8 pt-6 pb-8 mb-4">
        <h1>{{ $t('userResetPassword.title') }}</h1>
        <Alert v-if="alert.type" :type="alert.type" :message="alert.message" :title="alert.title" />
        <form @submit.prevent="onUpdatePassword" class="mt-3">
          <div class="mb-6">
            <label class="block text-grey-darker text-sm font-bold mb-2" for="password">
              {{ $t('userResetPassword.form.password') }}
            </label>
            <input v-model="password" class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline" id="password" type="password" placeholder="******************">
          </div>
          <div class="mb-6">
            <label class="block text-grey-darker text-sm font-bold mb-2" for="passwordCheck">
              {{ $t('userResetPassword.form.checkPassword') }}
            </label>
            <input v-model="checkPassword" class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline" id="passwordCheck" type="password" placeholder="******************">
          </div>
          <div class="flex items-center justify-between">
            <button class="bg-blue hover:bg-blue-dark text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">
              {{ $t('userResetPassword.form.submit') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { post } from '../utils/axiosHelper'
import errorHelper from '../utils/errorHelper'
import Alert from '../components/alerts/Alert'

export default {
  name: 'UserResetPassword',
  components: { Alert },
  data () {
    return {
      token: null,
      password: '',
      checkPassword: '',
      alert: {
        type: null,
        title: null,
        message: null
      }
    }
  },
  created () {
    this.token = this.$route.query.token
  },
  methods: {
    onUpdatePassword () {
      this.alert.type = null

      const data = {
        'password': this.password,
        'passwordCheck': this.checkPassword,
        'token': this.token
      }

      post('user/reset-password', data)
        .then(() => {
          this.alert = {
            type: 'success',
            title: 'Succès ',
            message: 'Votre mot a été modifier avec succès. Vous allez être redirigé sur la page de connexion dans quelques secondes.'
          }

          setTimeout(() => {
            this.$router.push({ name: 'auth' })
          }, 10000)
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
            title: 'Oups',
            message: errorHelper(code)
          }
        })
    }
  }
}
</script>
