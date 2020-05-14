<template>
  <section class="rounded-md shadow-md">
    <header class="p-2" :class="headColor">
      <slot name="header"></slot>
    </header>
    <main class="bg-white p-3">
      <slot></slot>
    </main>
    <footer class="p-2 flex justify-between" :class="footerColor">
      <slot name="footer"></slot>
    </footer>
  </section>
</template>

<script>
export default {
  name: 'AdminCard',
  props: ['headColor', 'footerColor']
}
</script>

<style scoped>

</style>
