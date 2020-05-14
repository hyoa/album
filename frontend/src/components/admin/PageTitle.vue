<template>
  <div class="text-white px-2 py-2 shadow-md" :class="color">
    <h1 class="text-3xl">
      <v-icon class="icon" v-if="icon" :name="icon" scale="1.2" :inverse="true"></v-icon>
      {{ title }}
    </h1>
  </div>
</template>

<script>
export default {
  name: 'PageTitle',
  props: ['title', 'icon', 'color']
}
</script>

<style scoped>
  .icon {
    position: relative;
    top: -1px;
  }
</style>
