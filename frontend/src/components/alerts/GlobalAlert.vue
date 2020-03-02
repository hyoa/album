<template>
  <div @click="onClose" class="fixed bottom-0 px-2 py-4 w-full z-50" :class="getClass">
    <div @click="onClose" class="absolute text-white right-0 top-0 mr-1 mt-1">
      <i class="material-icons">close</i>
    </div>
    {{ message }}
  </div>
</template>

<script>
export default {
  name: 'GlobalAlert',
  created () {
    setTimeout(() => {
      this.$store.commit('hideAlert')
    }, 10000)
  },
  methods: {
    onClose () {
      this.$store.commit('hideAlert')
    }
  },
  computed: {
    message () {
      return this.$store.state.alert.message
    },
    getClass () {
      switch (this.$store.state.alert.type) {
        case 'success':
          return 'bg-green-400 text-white border-green-500'
        case 'warning':
          return 'bg-orange-400 text-white border-orange-500'
        case 'error':
          return 'bg-red-400 text-white border-red-500'
        default:
          return 'bg-blue-400 text-white border-blue-500'
      }
    }
  }
}
</script>

<style scoped>

</style>
