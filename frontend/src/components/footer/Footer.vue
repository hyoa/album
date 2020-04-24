<template>
  <footer class="footer text-center">
    <div>
      Développé avec <span class="text-red-500">♥</span> par {{ author }}
    </div>
    <div v-if="isAdmin">
      <router-link class="text-black no-underline hover:underline" :to="{ name: 'admin_home' }">Administration</router-link>
    </div>
  </footer>
</template>

<script>
export default {
  name: 'Footer',
  props: ['isAdmin'],
  data () {
    return {
      author: process.env.VUE_APP_AUTHOR
    }
  }
}
</script>
