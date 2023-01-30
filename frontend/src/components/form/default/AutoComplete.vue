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
        @click="onSelect(searchResult.value, searchResult.isNew)"
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
import { graphql } from '../../../utils/axiosHelper'

export default {
  name: 'AutoComplete',
  props: ['label', 'placeholder', 'id', 'type', 'value', 'entity', 'allowNoCall'],
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

      let res
      if (this.entity === 'album') {
        res = await this.getAlbumAutocomplete(term)
      } else if (this.entity === 'folder') {
        res = await this.getFolderAutocomplete(term)
      } else {
        this.searchResults = []
        this.$emit('input', '')
        return
      }

      if (this.allowNoCall) {
        if (res.length === 0) {
          this.$emit('input', term)
          this.searchResults = [{ label: term, value: term, isNew: true }]
        } else {
          this.searchResults = res

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
    onSelect (value, isNew) {
      this.folder = value
      this.searchResults = []
      this.$emit('input', isNew ? 'new|' + value : value)
    },
    async getAlbumAutocomplete (value) {
      const query = `
      query {
        albums: albums(input: {includePrivate: true, includeNoMedias: true, limit: 100000, term: "${value}"}) {
          title
          slug
        }
      }
      `

      const { albums } = await graphql(query, 'v3')

      return albums.map(album => { return { label: album.title, value: album.slug } })
    },
    async getFolderAutocomplete (value) {
      const query = `
        query {
          folders: folders(input: {name: "${value}"}) {
            name
          }
        }
      `
      const { folders } = await graphql(query, 'v3')

      return folders.map(folder => { return { label: folder.name, value: folder.name } })
    }
  }
}
</script>
