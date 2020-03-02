<template>
  <button
    @click="e => $emit('click', e)"
    class="font-bold  rounded focus:outline-none focus:shadow-outline w-full md:w-1/3 bg-darker-primary hover:bg-dark-primary text-white"
    type="submit"
    :disabled="status === 'pending'"
    :class="[size === 'small' ? 'py-1 px-2' : 'py-2 px-4']"
  >
    <BubbleLoader v-if="status === 'pending'" />
    <slot v-else></slot>
  </button>
</template>

<script>
import BubbleLoader from '../../loader/BubbleLoader'

export default {
  name: 'SimpleAnimateButton',
  components: { BubbleLoader },
  props: ['status', 'size']
}
</script>
