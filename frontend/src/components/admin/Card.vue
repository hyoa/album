<template>
  <section>
    <header class="bg-light-primary p-2">
      <slot name="header"></slot>
    </header>
    <main class="bg-white p-3">
      <slot></slot>
    </main>
    <footer class="bg-lightest-primary p-2 flex justify-between">
      <slot name="footer"></slot>
    </footer>
  </section>
</template>

<script>
export default {
  name: 'AdminCard'
}
</script>

<style scoped>

</style>
