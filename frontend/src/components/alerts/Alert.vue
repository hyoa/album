<template>
  <div class="border-l-4 p-4" role="alert" :class="getClass">
    <p class="font-bold" v-if="title">{{ title }}</p>
    <p>{{ message }}</p>
  </div>
</template>

<script>
export default {
  name: 'Alert',
  props: {
    type: {
      type: String,
      default: 'success'
    },
    title: {
      type: String,
      default: null
    },
    message: {
      type: String,
      default: 'Something happened! But what ?!'
    }
  },
  computed: {
    getClass () {
      switch (this.type) {
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
