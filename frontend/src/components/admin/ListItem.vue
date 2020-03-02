<template>
  <li class="bg-white p-2 rounded-sm mt-2">
    <router-link :to="to">
      <div>
        <h2 class="text-xl">{{ title }}</h2>
        <div class="text-xs text-gray-700">{{ subtitle }}</div>
      </div>
      <div class="mt-2">
        <slot></slot>
      </div>
    </router-link>
  </li>
</template>

<script>
export default {
  name: 'ListItem',
  props: ['to', 'title', 'subtitle']
}
</script>

<style scoped>

</style>
