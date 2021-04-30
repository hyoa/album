<template>
  <div class="parent min-h-screen">
    <DefaultNav :isAdmin="isAdmin"></DefaultNav>
    <main class="main px-4">
      <slot></slot>
    </main>
  </div>
</template>

<script>
import DefaultNav from '../nav/DefaultNav'
export default {
  name: 'LayoutGrid',
  components: { DefaultNav },
  computed: {
    isAdmin () {
      return this.$store.state.token.role === 9
    }
  }
}
</script>

<style scoped lang="scss">
  .parent {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr 1fr 1fr;
    grid-template-rows: 80px 1fr 1fr 1fr 50px;
    grid-column-gap: 0;
    grid-row-gap: 0;
    .nav { grid-area: 1 / 1 / 2 / 6; }
    .footer { grid-area: 5 / 1 / 6 / 6; }
    .main { grid-area: 2 / 1 / 5 / 6; }
  }
</style>
