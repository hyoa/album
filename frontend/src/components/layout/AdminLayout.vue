<template>
  <div class="min-h-screen">
    <AdminNav/>
    <slot name="sidebar"></slot>
    <main class="container mx-auto px-4 pb-4 mt-4 relative">
      <slot></slot>
    </main>
    <GlobalAlert v-if="isGlobalAlertVisible" />
  </div>
</template>

<script>
import AdminNav from '../nav/AdminNav'
import GlobalAlert from '../alerts/GlobalAlert'

export default {
  name: 'AdminLayout',
  components: { AdminNav, GlobalAlert },
  computed: {
    isGlobalAlertVisible () {
      return this.$store.state.alert.isVisible
    }
  }
}
</script>
