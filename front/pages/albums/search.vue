<template>
    <div class="w-full p-4">
        <div class="text-sm breadcrumbs">
            <ul>
                <li><a>Albums</a></li> 
                <li><a>Recherche</a></li> 
            </ul>
        </div>
        <div>
            <div v-for="(album, index) in albums" :key="index">
                <NuxtLink :to="`/album/${album.slug}`">
                    <div class="card card-compact bg-base-100 shadow-xl mt-4">
                        <figure><img :src="album.favorites[0].urls.medium" :alt="album.title" /></figure>
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

        const variables = { limit: 200 }

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