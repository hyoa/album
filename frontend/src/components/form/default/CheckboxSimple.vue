<template>
  <label class="md:w-2/3 block text-white font-bold mb-3">
    <input
      :value="value"
      @input="$emit('input', !value)"
      class="mr-2 leading-tight"
      type="checkbox"
      :checked="value === true"
    >
    <span class="text-sm">
        {{ label }}
      </span>
  </label>
</template>

<script>
export default {
  name: 'CheckboxSimple',
  props: ['label', 'value']
}
</script>
