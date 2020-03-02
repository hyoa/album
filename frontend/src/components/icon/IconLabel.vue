<template>
  <div class="flex">
    <i class="material-icons mr-1">{{ icon }}</i> {{ label }}
  </div>
</template>

<script>
export default {
  name: 'IconLabel',
  props: ['icon', 'label']
}
</script>

<style scoped>

</style>
