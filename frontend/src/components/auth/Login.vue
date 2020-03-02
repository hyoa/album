<template>
  <section>
    <h1 class="mb-2">Se connecter</h1>
    <div class="mb-2">
      Vous n'avez pas de compte ? <a class="text-blue-500" href="#" @click="$emit('onRegister')">Créer en un maintenant !</a>
    </div>
    <form @submit.prevent="onLogin">
      <div class="mb-4">
        <label class="block text-grey-darker text-sm font-bold mb-2" for="email">
          Email
        </label>
        <input required v-model="email" class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline" id="email" type="text" placeholder="email@gmail.com">
      </div>
      <div class="mb-6">
        <label class="block text-grey-darker text-sm font-bold mb-2" for="password">
          Mot de passe
        </label>
        <input required v-model="password" class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline" id="password" type="password" placeholder="******************">
      </div>
      <div class="md:flex items-center md:justify-between">
        <SimpleAnimateButton :status="formStatus">
          Se connecter
        </SimpleAnimateButton>
        <div class="text-center md:text-left">
          <a @click.prevent="$emit('onPasswordForget')" class="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-700" href="#">
            Mot de passe oublié ?
          </a>
        </div>
      </div>
    </form>
  </section>
</template>

<script>
import SimpleAnimateButton from '../form/button/SimpleAnimateButton'

export default {
  name: 'Login',
  components: { SimpleAnimateButton },
  props: ['formStatus'],
  data () {
    return {
      email: '',
      password: ''
    }
  },
  methods: {
    onLogin () {
      const loginData = {
        email: this.email,
        password: this.password
      }
      this.$emit('onLogin', loginData)
    }
  }
}
</script>
