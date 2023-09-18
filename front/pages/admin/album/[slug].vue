<template>
    <div v-if="album">
        <section class="px-4">
            <div class="text-sm breadcrumbs">
                <ul>
                    <li><NuxtLink to="/admin">Admin</NuxtLink></li> 
                    <li><NuxtLink to="/admin/albums">Albums</NuxtLink></li> 
                    <li class="font-bold">Editer {{ album.title }}</li>
                </ul>
            </div>
        </section>
        <div class="drawer drawer-end">
            <input id="my-drawer-4" type="checkbox" class="drawer-toggle" />
            <div class="px-4 drawer-content">
                <div @click="notify.type = null" v-if="notify.type" class="bg-red-400 p-4 rounded-md" :class="[ notify.type === 'error'  ? 'bg-red-400' : 'bg-green-400']">
                    {{ notify.message }}
                </div>
                <div class="form-control w-full max-w-xs">
                    <label class="label">
                        <span class="label-text">Titre</span>
                    </label>
                    <input v-model="album.title" type="text" placeholder="Type here" class="input input-bordered w-full max-w-xs" />
                </div>
                <div class="form-control mt-4">
                    <label class="label">
                        <span class="label-text">Description</span>
                    </label>
                    <textarea v-model="album.description" class="textarea textarea-bordered h-24" placeholder="Bio"></textarea>
                </div>
                <div class="form-control w-1/4 mt-4">
                    <label class="label cursor-pointer">
                        <input type="checkbox" class="checkbox" v-model="album.private" />
                        <span class="label-text">Privé</span> 
                    </label>
                </div>
                <div class="mt-4 flex justify-between">
                    <button @click="onEdit" class="btn btn-primary btn-sm">Enregistrer</button>
                    <label for="my-drawer-4" class="drawer-button btn btn-secondary btn-sm">Bibliothèque</label>
                </div>
                <div class="grid gap-2 grid-cols-2 grid-rows-3 mt-6">
                    <div 
                        v-for="(media, index) in album.medias"
                        :class="{'border-2 border-primary': isSelected(media.key, 'remove')}"
                    >
                        <div v-if="media.favorite">
                            <Icon @click="onToggleFavorite(media)" name="teenyicons:star-solid" class="text-2xl m-2 absolute z-10" />
                        </div>
                        <div v-else>
                            <Icon @click="onToggleFavorite(media)" name="teenyicons:star-outline" class="text-2xl m-2 absolute z-10" />
                        </div>
                        <div @click="selectMedia(media.key, 'remove')">
                            <img v-if="media.kind === 'PHOTO'" :src="media.urls.small">
                            <video v-else-if="media.kind === 'VIDEO'" :src="media.urls.small"></video>
                        </div>
                    </div>
                </div>
                <div
                    v-if="mediasSelectedToRemove.length > 0"
                    class="fixed bottom-10"
                >
                    <div v-if="actionInProgress">
                        <span class="loading loading-dots loading-sm"></span>
                    </div>
                    <div v-else>
                        <button @click="onChangeMedias('remove')" class="btn btn-primary btn-sm mt-4">Supprimer {{ mediasSelectedToRemove.length }} médias</button>
                    </div>
                </div>
            </div>
            <div class="drawer-side z-20" id="drawer-side-modal">
                <label for="my-drawer-4" class="drawer-overlay"></label>
                <div class="menu p-4 w-80 min-h-full bg-base-200 text-base-content">
                    <div>
                        <h3 class="text-xl">Dossiers</h3>
                        <input v-model="searchFolder" type="text" placeholder="Chercher un dossier" class="input input-bordered w-full max-w-xs mt-2" />
                    </div>
                    <div class="mt-4">
                        <h3 class="text-xl">Résultat</h3>
                        <ul>
                            <li @click="loadFolder(folder.name)" v-for="(folder, index) in filteredFolder">{{ folder.name }}</li>
                        </ul>
                    </div>
                    <div class="divider"></div>
                    <div class="grid gap-2 grid-cols-2 grid-rows-3 mt-6">
                        <div 
                            @click="selectMedia(media.key, 'add')" 
                            v-for="(media, index) in folderToDisplayMedias" 
                            :key="index"
                            :class="{'border-2 border-primary': isSelected(media.key, 'add')}"
                        >
                            <img v-if="media.kind === 'PHOTO'" :src="media.urls.small">
                            <video v-else-if="media.kind === 'VIDEO'" :src="media.urls.small"></video>
                        </div>
                    </div>
                </div>
                <button @click="onChangeMedias('add')" v-if="mediasSelectedToAdd.length > 0" class="btn btn-sm btn-primary mt-4 mr-4 fixed bottom-5">
                    <span class="loading loading-dots loading-sm" v-if="actionInProgress"></span>
                    <span v-else>Ajouter {{ mediasSelectedToAdd.length }} médias</span>  
                </button>
            </div>
        </div>
    </div>
    <div v-else class="text-center">
        <span class="loading loading-dots loading-lg"></span>
    </div>
</template>

<script>
export default {
    name: "AlbumView",
    data() {
        return {
            album: null,
            folders: [],
            searchFolder: '',
            folderToDisplayMedias: [],
            mediasSelectedToAdd: [],
            mediasSelectedToRemove: [],
            actionInProgress: false,
            notify: {
                message: '',
                type: null
            }
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
                    author
                    key
                    kind
                    favorite
                    urls {
                        small
                    }
                }
                private
            }
        }
        `

        const variables = { slug: this.$route.params.slug }

        try {
            const { data: { _rawValue: { album } } } = await useAsyncQuery(query, variables)
            this.album = album

            if (this.album.medias === null) {
                this.album.medias = []
            }
        } catch (e) {
            console.log(e)
            this.errorMessage = e;
        }

        const queryFolder = gql`
            query {
                folders: folders(input: {}){
                    name
                }
            }
        `

        try {
            const { data: { _rawValue: { folders } } } = await useAsyncQuery(queryFolder)
            this.folders = folders
        } catch (e) {
            console.log(e)
            this.errorMessage = e;
        }
    },
    methods: {
        async loadFolder (folder) {
            this.mediasSelected = []

            const queryFolder = gql`
            query {
                folder: folder(input: {name: "${folder}"}){
                    medias {
                        key
                        kind
                        author
                        urls {
                            small
                        }
                    }
                }
            }
            `

            try {
                const { data: { _rawValue: { folder } } } = await useAsyncQuery(queryFolder)
                this.folderToDisplayMedias = folder.medias.filter(media => !this.album.medias.map(m => m.key).includes(media.key))
            } catch (e) {
            }
        },
        selectMedia (media, kind) {
            if (kind === 'add') {
                if (this.mediasSelectedToAdd.includes(media)) {
                    this.mediasSelectedToAdd = this.mediasSelectedToAdd.filter(m => m !== media)
                    return
                }

                this.mediasSelectedToAdd.push(media)
            } else {
                if (this.mediasSelectedToRemove.includes(media)) {
                    this.mediasSelectedToRemove = this.mediasSelectedToRemove.filter(m => m !== media)
                    return
                }

                this.mediasSelectedToRemove.push(media)
            }
        },
        isSelected (media, kind) {
            if (kind === 'add') {
                return this.mediasSelectedToAdd.includes(media)
            } else {
                return this.mediasSelectedToRemove.includes(media)
            }
        },
        async onChangeMedias (kind) {
            if (this.actionInProgress) {
                return
            }

            this.actionInProgress = true
            const query = gql`
                mutation ($medias: [MediaAlbumInput!]!) {
                    album: updateAlbumMedias(input: {slug: "${this.$route.params.slug}", medias: $medias, action: ${kind.toUpperCase()}}) {
                        medias {
                            author
                            key
                            urls {
                                small
                            }
                            kind
                            favorite
                        }
                    }
                }
            `

            let mediasSelected = []
            let mediasToCompare = []
            if (kind === 'add') {
                mediasSelected = this.mediasSelectedToAdd
                mediasToCompare = this.folderToDisplayMedias
            } else {
                mediasSelected = this.mediasSelectedToRemove
                mediasToCompare = this.album.medias
            }


            const mediasObject = mediasSelected.map(media => {
                
                let mediaObject = mediasToCompare.find(m => m.key === media)
                
                return {
                    key: mediaObject.key,
                    kind: mediaObject.kind,
                    author: mediaObject.author
                }
            })

            const variables = {
                medias: mediasObject,
            }

            try {
                const { data: { _rawValue: { album } } } = await useAsyncQuery(query, variables)
                this.album.medias = album.medias
                this.mediasSelectedToAdd = []
                this.mediasSelectedToRemove = []


                this.notify = {
                    message: 'Mise à jour effectuée.',
                    type: 'success'
                }
            } catch (e) {
                this.notify = {
                    message: 'Une erreur est survenue.',
                    type: 'error'
                }
            }

            this.actionInProgress = false
        },
        async onEdit () {
            if (this.actionInProgress) {
                return
            }

            this.actionInProgress = true

            const query = gql`
                mutation ($input: UpdateAlbumInput!) {
                    album: updateAlbum(input: $input) {
                        title
                        slug
                        description
                        author
                        creationDate
                        medias {
                            author
                            key
                            urls {
                                small
                            }
                        }
                        private
                    }
                }
            `

            const variables = {
                input: {
                    slug: this.$route.params.slug,
                    title: this.album.title,
                    description: this.album.description,
                    private: this.album.private,
                    author: this.album.author,
                }
            }

            try {
                const { data: { _rawValue: { album } } } = await useAsyncQuery(query, variables)
                this.album = album

                this.notify = {
                    message: 'Mise à jour effectuée.',
                    type: 'success'
                }
            } catch (e) {
                this.errorMessage = e;

                this.notify = {
                    message: 'Une erreur est survenue.',
                    type: 'error'
                }
            }

            this.actionInProgress = false
        },
        async onToggleFavorite (media) {
            const query = gql`
                mutation {
                    updateAlbumFavorite(input: {slug: "${this.$route.params.slug}", mediaKey: "${media.key}"}) {
                        title
                    }
                }
            `

            try {
                await useAsyncQuery(query)
                this.album.medias.forEach(({ key }, index) => {
                    if (key === media.key) {
                        this.album.medias[index].favorite = !media.favorite
                    }
                })

                this.notify = {
                    message: 'Mise à jour effectuée.',
                    type: 'success'
                }
            } catch (e) {
                this.notify = {
                    message: 'Une erreur est survenue.',
                    type: 'error'
                }
            }
        },
    },
    computed: {
        filteredFolder () {
            return this.folders.filter(folder => folder.name.includes(this.searchFolder))
        },
    }
}
</script>