<template>
    <div v-if="folder">
        <section class="px-4">
            <div class="text-sm breadcrumbs">
                <ul>
                    <li><NuxtLink to="/admin">Admin</NuxtLink></li> 
                    <li><NuxtLink to="/admin/medias">Dossiers</NuxtLink></li> 
                    <li class="font-bold">Modifier {{ folder.name }}</li>
                </ul>
            </div>
        </section>
        <section class="px-4">
            <div @click="notify.type = null" v-if="notify.type" class="bg-red-400 p-4 rounded-md" :class="[ notify.type === 'error'  ? 'bg-red-400' : 'bg-green-400']">
                {{ notify.message }}
            </div>
        </section>
        <section class="px-4 mt-4" v-if="mediasSelected.length == 0">
            <h2>Mettre à jour le nom du dosser</h2>
            <div>
                <input v-model="folder.name" type="text" placeholder="Changer nom" class="input input-bordered w-full max-w-xs" />
            </div>
            <div class="mt-2">
                <button @click="onChangeName" class="btn btn-sm btn-primary">
                    <span class="loading loading-dots loading-sm" v-if="actionInProgress"></span>
                    <span v-else>Enregistrer</span>
                </button>
            </div>
        </section>
        <section class="px-4 mt-4" v-if="mediasSelected.length > 0">
            <h2>Changer les médias de dossier</h2>
            <div>
                <input v-model="newFolderName" type="text" placeholder="Nouveau dossier" class="input input-bordered w-full max-w-xs" />
            </div>
            <div class="mt-2">
                <button @click="onChangeFolder" class="btn btn-sm btn-primary">
                    <span class="loading loading-dots loading-sm" v-if="actionInProgress"></span>
                    <span v-else>Valider</span>
                </button>
            </div>
        </section>
        <div class="divider px-4"></div>
        <section class="px-4">
            <div class="flex flex-wrap">
                <div
                    v-for="(media, index) in folder.medias"
                    :key="index"
                    class="w-1/3 p-1"
                    @click="selectMedia(media)"
                    :class="{'border-4 border-blue-500': isSelected(media)}"
                >
                    <img :src="media.urls.small" alt="">
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
    name: "MediaUpdate",
    components: {
    },
    data() {
        return {
            folder: null,
            mediasSelected: [],
            actionInProgress: false,
            newFolderName: '',
            notify: {
                message: '',
                type: null
            }
        }
    },
    async created() {
        if (this.actionInProgress) {
            return
        }

        this.actionInProgress = true

        const queryFolder = gql`
            query {
                folder: folder(input: {name: "${this.$route.params.id}"}){
                    name
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
            this.folder = folder
        } catch (e) {
            
        }

        this.actionInProgress = false
    },
    methods: {
        selectMedia (media) {
            if (this.mediasSelected.includes(media)) {
                this.mediasSelected = this.mediasSelected.filter(m => m !== media)
            } else {
                this.mediasSelected.push(media)
            }
        },
        isSelected (media) {
            return this.mediasSelected.includes(media)
        },
        async onChangeFolder() {
            if (!this.newFolderName || this.newFolderName === '') {
                return
            }

            const query = gql`
                mutation ($keys: [String!]!) {
                    folder: changeMediasFolder(input: {folderName: "${this.newFolderName}" , keys: $keys}) {
                        medias {
                            key
                            author
                            kind
                            urls {
                                small
                            }
                        }
                    }
                }
            `

            const variables = {
                keys: this.mediasSelected.map(media => media.key)
            }

            try {
                const { data: { _rawValue: { folder } } } = await useAsyncQuery(query, variables)

                this.folder.medias = this.folder.medias.filter(media => !folder.medias.map(m => m.key).includes(media.key))

                this.mediasSelected = []

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
        async onChangeName() {
            if (this.actionInProgress) {
                return
            }

            if (!this.folder.name || this.folder.name === '') {
                return
            }

            if (this.folder.name === this.$route.params.id) {
                return
            }

            this.actionInProgress = true

            const query = gql`
                mutation {
                    folder: changeFolderName(input: {oldName: "${this.$route.params.id}", newName: "${this.folder.name}"}) {
                        name
                    }
                }
            `

            try {
                const { data: { _rawValue: { folder } } } = await useAsyncQuery(query)
                this.actionInProgress = false
                return navigateTo(`/admin/media/${folder.name}`)
            } catch (e) {
                this.notify = {
                    message: 'Une erreur est survenue.',
                    type: 'error'
                }
            }

            this.actionInProgress = false
        }
    },
}

</script>