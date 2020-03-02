<template>
  <section>
    <h1 class="mb-4">S'inscrire</h1>
    <form @submit.prevent="onRegister">
      <div class="mb-4">
        <label class="block text-grey-darker text-sm font-bold mb-2" for="email">
          Email
        </label>
        <input
          v-model="email"
          class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline"
          id="email"
          type="email"
          placeholder="email@gmail.com"
          :class="{ 'border border-red-600': errors.email }"
        >
        <span class="text-red-600" v-if="errors.email">L'email doit être renseigné</span>
      </div>
      <div class="mb-4">
        <label class="block text-grey-darker text-sm font-bold mb-2" for="name">
          Nom/Prénom
        </label>
        <input required v-model="name" class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline" id="name" type="text" placeholder="Jules" :class="{ 'border border-red-600': errors.name }">
        <span class="text-red-600" v-if="errors.name">Un nom/prénom doit être renseigné</span>
      </div>
      <div class="mb-6">
        <label class="block text-grey-darker text-sm font-bold mb-2" for="password">
          Mot de passe
        </label>
        <input required v-model="password" class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline" id="password" type="password" placeholder="******************" :class="{ 'border border-red-600': errors.password }">
        <span class="text-red-600" v-if="errors.password">Les mot de passe ne correspondent pas</span>
      </div>
      <div class="mb-6">
        <label class="block text-grey-darker text-sm font-bold mb-2" for="passwordCheck">
          Confirmation du mot de passe
        </label>
        <input required v-model="checkPassword" class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline" id="passwordCheck" type="password" placeholder="******************" :class="{ 'border border-red-600': errors.checkPassword }">
        <span class="text-red-600" v-if="errors.checkPassword">Les mot de passe ne correspondent pas</span>
      </div>
      <div class="flex items-center justify-between">
        <SimpleAnimateButton :status="formStatus">
          S'inscrire
        </SimpleAnimateButton>
      </div>
    </form>
  </section>
</template>

<script>
import SimpleAnimateButton from '../form/button/SimpleAnimateButton'
export default {
  name: 'Register',
  components: { SimpleAnimateButton },
  props: ['formStatus'],
  data () {
    return {
      email: '',
      password: '',
      checkPassword: '',
      name: '',
      errors: {
        email: false,
        password: false,
        checkPassword: false,
        name: false
      }
    }
  },
  methods: {
    onRegister () {
      const data = {
        email: this.email,
        password: this.password,
        checkPassword: this.checkPassword,
        name: this.name
      }

      this.errors = {
        email: false,
        password: false,
        checkPassword: false,
        name: false
      }

      if (!data.email.includes('@') || data.email.trim() === '') {
        this.errors.email = true

        return
      }

      if (data.password !== data.checkPassword) {
        this.errors.password = true
        this.errors.checkPassword = true

        return
      }

      this.$emit('onRegister', data)
    }
  }
}
</script>
