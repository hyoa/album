<template>
    <div class="w-full p-4">
        <div class="text-sm breadcrumbs">
            <ul>
                <li><NuxtLink to="/">Albums</NuxtLink></li> 
                <li>Recherche</li> 
            </ul>
        </div>
        <div>
            <input v-model="search" type="text" placeholder="Rechercher" class="input input-bordered w-full max-w-xs" />
            <div v-for="(album, index) in filteredAlbums" :key="index">
                <NuxtLink :to="`/album/${album.slug}`">
                    <div class="card card-compact bg-base-100 shadow-xl mt-4">
                        <figure>
                            <img v-if="album.favorites[0].kind === 'PHOTO'" :src="album.favorites[0].urls.medium" :alt="album.title" />
                            <video autoplay muted loop v-else-if="album.favorites[0].kind === 'VIDEO'" :src="album.favorites[0].urls.medium" :alt="album.title"></video>
                        </figure>
                        <div class="card-body">
                            <p class="card-title">{{ album.title }}</p>
                        </div>
                    </div>
                </NuxtLink>
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
            albums: [],
            search: '',
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

        const variables = { limit: 200 }

        try {
            const { data: { _rawValue: { albums } } } = await useAsyncQuery(query, variables)
            this.albums = albums
        } catch (e) {
            console.log(e)
            this.errorMessage = e;
        }
    },
    computed: {
        filteredAlbums() {
            return this.albums.filter(album => album.title.toLowerCase().includes(this.search.toLowerCase()))
        }
    }
}
</script>