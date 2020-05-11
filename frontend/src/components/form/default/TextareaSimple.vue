<template>
  <div class="mb-3">
    <label class="block uppercase tracking-wide text-sm mb-2" :for="id">
      {{ label }}
    </label>
    <textarea
      class="appearance-none block w-full text-gray-darker border border-grey-lighter py-3 px-4 leading-tight rounded-sm focus:outline-none focus:shadow-md"
      :id="id"
      :type="type"
      :placeholder="placeholder"
      @input="$emit('input', $event.target.value)"
      :value="value"
      rows="8"
    ></textarea>
  </div>
</template>

<script>
export default {
  name: 'TextareaSimple',
  props: ['label', 'placeholder', 'id', 'type', 'value']
}
</script>
