<template>
  <div class="fixed h-full w-screen z-40 top-0">
    <div class="bg-white h-full absolute z-50 w-2/3" :class="[ side === 'right' ? 'right-0' : 'left-0' ]">
      <slot></slot>
    </div>
    <div class="bg-black h-full opacity-50"></div>
  </div>
</template>

<script>
export default {
  name: 'Sidebar',
  props: ['side', 'size'],
  created () {
    this.$store.commit('toggleSideBar', true)
  },
  destroyed () {
    this.$store.commit('toggleSideBar', false)
  }
}
</script>
