<template>
  <div class="mb-3">
    <label class="block uppercase tracking-wide text-white text-xs font-bold mb-2" :for="id">
      {{ label }}
    </label>
    <input
      class="appearance-none block w-full text-gray-darker border border-grey-lighter py-3 px-4 leading-tight rounded-sm focus:outline-none focus:shadow-md"
      :id="id"
      :type="type"
      :placeholder="placeholder"
      :value="value"
      @input="$emit('input', $event.target.value)"
    >
  </div>
</template>

<script>
export default {
  name: 'InputSimple',
  props: ['label', 'placeholder', 'id', 'type', 'value']
}
</script>
