<template>
    <div v-if="album != null">
        <section class="text-center px-4">
            <h1 class="text-2xl">{{ album.title }}</h1>
            <hr class="border-t border-grey-400 w-1/2 md:w-1/4 mx-auto">
            <div class="text-gray-800">
                <div :class="{ 'h-12' : !readMore }" class="text-gray-700 text-base italic overflow-hidden description">
                    {{ album.description }}
                </div>
                <div>
                    <button 
                        @click:="!readMore" 
                        class="btn btn-xs"
                        v-if="album.description && album.description.length > 200"
                    >
                        {{ readMoreLabel }}
                    </button>
                </div>
                <div class="mt-2 text-grey-700 text-sm">
                    <span>Crée le {{ getCreationDate(album.creationDate) }} par {{ album.author }}</span>
                </div>
            </div>
        </section>
        <section class="mt-8">
            <div w-full class="bg-slate-900">
                <div v-for="(media, index) in mediasToShow" :key="index" class="mt-2">
                    <img v-if="media.kind === 'PHOTO'" loading="lazy" :src="media.urls.large">
                    <video v-else-if="media.kind === 'VIDEO'" controls preload="metadata">
                        <source :src="media.urls.large" type="video/mp4">
                        Your browser does not support the video tag.
                    </video>
                </div>
            </div>
        </section>
    </div>
    <div v-else class="text-center">
        <span class="loading loading-dots loading-lg"></span>
    </div>
</template>

<script>
export default {
    name: 'Album',
    components: {
    },
    data() {
        return {
            album: null,
            readMore: false,
            mediasToShow: [],
            indexDisplay: 0,
        }
    },
    async created() {
        const query = gql`
        query getAlbum($slug: String!) {
            album(input: {slug: $slug}) {
                title
                slug
                description
                author
                creationDate
                medias {
                    key
                    kind
                    urls {
                        large
                    }
                }
            }
        }
        `

        const variables = { slug: this.$route.params.slug }

        try {
            const { data: { _rawValue: { album } } } = await useAsyncQuery(query, variables)
            this.album = album
            this.mediasToShow = album.medias.slice(0, 10)
        } catch (e) {
            console.log(e)
            this.errorMessage = e;
        }
    },
    mounted () {
        this.scroll()
    },
    computed: {
        readMoreLabel () {
            if (this.album.description && this.album.description.length > 200) {
                return this.readMore ? 'Cacher' : 'Lire la suite'
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
        },
        load() {
            if (this.album === null) {
                return
            }

            if (this.indexDisplay >= this.album.medias.length) {
                return
            }

            this.indexDisplay += 10

            // merge new medias
            this.mediasToShow = this.mediasToShow.concat(this.album.medias.slice(this.indexDisplay, this.indexDisplay + 10))

            return 
        },
        scroll () {
            window.onscroll = () => {
                let bottomOfWindow = Math.max(window.pageYOffset, document.documentElement.scrollTop, document.body.scrollTop) + window.innerHeight === document.documentElement.offsetHeight

                if (bottomOfWindow) {
                    this.load()
                }
            }
        },
    }
}
</script>