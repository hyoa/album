<template>
  <nav class="nav flex justify-between items-center py-3 px-4 relative">
    <div>
      <router-link :to="{ name: 'home' }" class="text-3xl no-underline text-black font-manuscript text-primary">{{ title }}</router-link>
    </div>
    <i @click="isMenuOpen = !isMenuOpen" class="material-icons text-3xl text-darker-primary">account_circle</i>
    <div v-if="isMenuOpen" class="menu">
      <ul>
        <li v-if="isAdmin">
            <router-link class="text-black no-underline hover:underline" :to="{ name: 'admin_home' }">Administration</router-link>
        </li>
        <li v-if="isAdmin" class="border-t-2 border-darker-primary my-2"></li>
        <li @click="logout">
          Se déconnecter
        </li>
      </ul>
    </div>
  </nav>
</template>

<script>
export default {
  name: 'DefaultNav',
  props: {
    isAdmin: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      isMenuOpen: false,
      title: process.env.VUE_APP_SITE_TITLE
    }
  },
  methods: {
    logout () {
      localStorage.removeItem('album-token')
      this.$router.push({ name: 'auth' })
    }
  }
}
</script>

<style lang="scss" scoped>
  .menu {
    @apply absolute shadow-xl p-2 rounded-sm border border-primary z-10 bg-white;
    top: 3.5rem;
    right: 1rem;
  }
</style>
