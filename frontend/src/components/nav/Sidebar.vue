<template>
  <div class="absolute h-screen-minus-header w-screen z-40">
    <div class="bg-white h-screen-minus-header absolute z-50 w-2/3" :class="[ side === 'right' ? 'right-0' : 'left-0' ]">
      <slot></slot>
    </div>
    <div class="bg-black opacity-50 h-screen-minus-header"></div>
  </div>
</template>

<script>
export default {
  name: 'Sidebar',
  props: ['side', 'size'],
  created () {
    window.scrollTo(0, 0)
    this.$store.commit('toggleSideBar', true)
  },
  destroyed () {
    this.$store.commit('toggleSideBar', false)
  }
}
</script>

<style scoped>
 .h-screen-minus-header {
   height: calc(100vh - 48px)
 }
</style>
