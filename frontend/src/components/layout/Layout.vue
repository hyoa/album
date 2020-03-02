<template>
  <div>
    <DefaultNav/>
    <div class="container mx-auto px-4">
      <slot></slot>
    </div>
    <Footer :isAdmin="isAdmin" />
  </div>
</template>

<script>
import DefaultNav from '../nav/DefaultNav'
import Footer from '../footer/Footer'

export default {
  name: 'Layout',
  components: { Footer, DefaultNav },
  data () {
    return {
      tokenDecoded: {}
    }
  },
  computed: {
    isAdmin () {
      return this.$store.state.token.role === 9
    }
  }
}
</script>
