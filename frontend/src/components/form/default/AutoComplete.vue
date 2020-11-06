<template>
  <div class="mb-3">
    <label class="block uppercase tracking-wide text-sm mb-2" :for="id">
      {{ label }}
    </label>
    <input
      class="appearance-none block w-full text-gray-darker border border-grey-lighter py-3 px-4 leading-tight rounded-sm focus:outline-none focus:shadow-md"
      :id="id"
      :type="type"
      :placeholder="placeholder"
      :value="value"
      @input="onInput($event.target.value)"
      v-model="folder"
    >
    <ul v-if="searchResults.length" class="list-reset border border-gray bg-white border-t-0 rounded-b autocomplete">
      <li
        class="my-1 text-grey-dark hover:bg-grey-lightest cursor-pointer p-2 flex"
        :key="searchResult.value"
        v-for="searchResult of searchResults"
        @click="onSelect(searchResult.value)"
      >
        {{ searchResult.label }}
        <div class="bg-blue-500 text-white rounded-md px-2 ml-3" v-if="searchResult.isNew">
          {{ $t('components.autocomplete.isNew') }}
        </div>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.autocomplete {
  max-height: 20vh;
  overflow-y: scroll;
}
</style>

<script>
import debounce from 'debounce'
import { get } from '../../../utils/axiosHelper'

export default {
  name: 'AutoComplete',
  props: ['label', 'placeholder', 'id', 'type', 'value', 'endpoint', 'allowNoCall'],
  data () {
    return {
      folder: '',
      searchResults: []
    }
  },
  methods: {
    onInput: debounce(async function (term) {
      if (term.trim() === '') {
        this.searchResults = []
        this.$emit('input', '')
        return
      }

      const res = await get(`${this.endpoint}?search=${term}`)

      if (this.allowNoCall) {
        if (res.data.length === 0) {
          this.$emit('input', term)
          this.searchResults = [{ label: term, value: term, isNew: true }]
        } else {
          this.searchResults = res.data

          if (!this.searchResults.some(({ label, value }) => {
            return label === term && value === term
          })) {
            this.searchResults.push({ label: term, value: term, isNew: true })
          }
        }
      } else {
        this.searchResults = res.data
      }
    }, 200),
    onSelect (value) {
      this.folder = value
      this.searchResults = []
      console.log(value)
      this.$emit('input', value)
    }
  }
}
</script>
