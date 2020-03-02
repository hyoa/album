<template>
  <button @click="e => $emit('click', e)" class="bg-darker-primary hover:bg-dark-primary text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full md:w-1/3" type="submit">
    <slot></slot>
  </button>
</template>

<script>
export default {
  name: 'FormButton'
}
</script>

<style scoped>

</style>
