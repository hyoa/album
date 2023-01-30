<template>
  <Layout>
    <div class="text-center mb-6">
      <h1 data-e2e="album-title" class="text-3xl">{{ album.title }}</h1>
      <hr class="border-t border-grey-400 w-1/2 md:w-1/4 mx-auto">
      <div data-e2e="album-description" class="mt-2">
        <p v-if="album.description" class="text-grey-700 text-base italic overflow-hidden description" :class="{ 'h-12' : !readMore }">
          {{ album.description }}
        </p>
        <button
          @click="readMore = !readMore"
          class="text-blue-500 text-underline text-sm mt-3 focus:outline-none active:outline-none border border-blue-500 p-1 rounded"
          v-if="album.description && album.description.length > 200"
        >
          {{ $t(readMoreLabel) }}
        </button>
      </div>
      <p data-e2e="album-informations" class="mt-2 text-grey-700 text-sm">
        <span>{{ $t('album.createAtBy', { date: getCreationDate(album.creationDate), author: album.author }) }}</span>
      </p>
    </div>
    <Grid :medias="album.medias" :editable="false" :can-delete-media="false" />
  </Layout>
</template>

<script>
import { graphql } from '../utils/axiosHelper'

import Layout from '../components/layout/LayoutGrid'
import Grid from '../components/grid/Grid'
export default {
  name: 'AlbumView',
  components: { Grid, Layout },
  data () {
    return {
      album: {},
      readMore: false
    }
  },
  async created () {
    try {
      const query = `
      query {
        album: album(input: {slug: "${this.$route.params.slug}"}) {
          title
          slug
          description
          author
          creationDate
          medias {
            key
            kind
            urls {
              small
              medium
              large
            }
          }
        }
      }
      `

      const res = await graphql(query, 'v3')

      this.album = res.album
    } catch ({ response: { status } }) {
      if (status === 401) {
        this.$store.commit('setFlashMessage', 'auth.alert.disconnect.message')
        localStorage.removeItem('album-token')
        this.$router.push({ name: 'auth' })
      }
    }
  },
  computed: {
    readMoreLabel () {
      if (this.album.description && this.album.description.length > 200) {
        return this.readMore ? 'album.readMore.hide' : 'album.readMore.show'
      }

      return true
    }
  },
  methods: {
    getCreationDate (date) {
      if (date) {
        const dtf = new Intl.DateTimeFormat()
        return dtf.format(date * 1000)
      }
    }
  }
}
</script>
