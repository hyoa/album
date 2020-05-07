<template>
  <section>
    <h1 class="mb-4">{{ $t('auth.askResetPasswordPage.title') }}</h1>
    <form @submit.prevent="onAskResetPassword">
      <div class="mb-4">
        <label class="block text-grey-darker text-sm font-bold mb-2" for="email">
          {{ $t('auth.askResetPasswordPage.form.email') }}
        </label>
        <input v-model="email" class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline" id="email" type="email" placeholder="email@gmail.com">
      </div>
      <div class="flex items-center justify-between">
        <SimpleAnimateButton :status="formStatus">
          {{ $t('auth.askResetPasswordPage.form.submit') }}
        </SimpleAnimateButton>
      </div>
    </form>
  </section>
</template>

<script>
import SimpleAnimateButton from '../form/button/SimpleAnimateButton'
export default {
  name: 'ResetPassword',
  components: { SimpleAnimateButton },
  props: ['formStatus'],
  data () {
    return {
      email: ''
    }
  },
  methods: {
    onAskResetPassword () {
      const url = window.location.href
      const arr = url.split('/')
      const callbackUri = `${arr[0]}//${arr[2]}/?#/reset-password`

      const data = {
        email: this.email,
        callbackUri
      }

      this.$emit('onAskResetPassword', data)
    }
  }
}
</script>
