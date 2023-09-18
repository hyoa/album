<template>
    <div class="text-sm breadcrumbs px-4">
        <ul>
            <li><NuxtLink to="/admin">Admin</NuxtLink></li> 
            <li class="font-bold">Albums</li>
        </ul>
    </div>
    <section class="px-4 mt-4">
        <div class="w-full">
            <input v-model="search" type="text" placeholder="Chercher" class="input input-bordered w-full max-w-xs" />
        </div>
    </section>
    <section class="mt-6 px-4">
        <div
            v-for="(album, index) in filteredAlbums"
            :key="index"
            class="shadow-md rounded-md p-2"
        >
            <NuxtLink :to="`album/${album.slug}`">
                <div class="flex justify-between">
                    <div class="text-xl">{{ album.title }}</div>
                    <div>{{ album.medias?.length ?? 0 }} medias</div>
                </div>
                <div class="text-sm text-gray-600">{{ getCreationDate(album.creationDate) }}</div>
             </NuxtLink>
        </div>
    </section>
</template>

<script>
export default {
    name: "Albums",
    data() {
        return {
            albums: [],
            search: ''
        }
    },
    async created () {
        const queryAlbums = gql`
            query {
                albums: albums(input: {limit: 1000, includePrivate: true, includeNoMedias: true}) {
                    title
                    slug
                    medias {
                        kind
                    }
                    private
                    creationDate
                }
            }
        `

        const { data: { _rawValue: { albums } } } = await useAsyncQuery(queryAlbums, {})
        this.albums = albums
    },
    methods: {
        getCreationDate (date) {
            if (date) {
                const dtf = new Intl.DateTimeFormat()
                return dtf.format(date * 1000)
            }
        }
    },
    computed: {
        filteredAlbums() {
            if (this.search === '') return this.albums

            return this.albums.filter(album => album.title.toLowerCase().includes(this.search.toLowerCase()))
        }
    }
}
</script>