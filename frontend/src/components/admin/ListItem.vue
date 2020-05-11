<template>
  <li class="bg-gray-200 p-2 rounded-md mt-2">
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
