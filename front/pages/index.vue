<template>
  <div>
      <div class="h-screen">
        <div class="carousel carousel-center max-w-md p-4 space-x-4 bg-slate-900 h-4/6 w-ful">
          <div v-for="(album) in albums" :key="album.slug" class="carousel-item w-full">
            <NuxtLink :to="`/album/${album.slug}`" class="block">
              <img :src="album.favorites[0].urls.medium" class="h-full w-full">
            </NuxtLink>
          </div>
        </div>
        <div class="flex align-middle justify-center mt-12">
          <NuxtLink to="/albums/search" class="btn btn-primary">Voir plus</NuxtLink>
        </div>
      </div>

  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})
</script>

<script>
export default {
name: 'Albums',
components: {
},
data() {
  return {
      albums: []
  }
},
async created() {
  const query = gql`
  query getAlbums($limit: Int!) {
      albums(input: {limit: $limit}) {
        title
        slug
        favorites {
          kind
          urls {
            small
            medium
          }
        }
      }
    }
  `

  const variables = { limit: 5 }

  try {
      const { data: { _rawValue: { albums } } } = await useAsyncQuery(query, variables)
      this.albums = albums
  } catch (e) {
      console.log(e)
      this.errorMessage = e;
  }
},
}
</script>