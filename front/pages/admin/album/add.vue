<template>
    <section class="px-4">
        <div class="text-sm breadcrumbs">
            <ul>
                <li><NuxtLink to="/admin">Admin</NuxtLink></li> 
                <li><NuxtLink to="/admin/albums">Albums</NuxtLink></li> 
                <li class="font-bold">Créer un album</li>
            </ul>
        </div>
    </section>
    <div class="px-4">
        <div v-if="notify.type" class="bg-red-400 p-4 rounded-md">
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
                <input v-model="album.private" type="checkbox" class="checkbox" />
                <span class="label-text">Privé</span> 
            </label>
        </div>
        <div class="mt-4">
            <button @click="onAdd" class="btn btn-primary btn-sm">
                <span class="loading loading-dots loading-sm" v-if="actionInProgress"></span>
                <span v-else>Enregistrer</span>  
            </button>
        </div>
    </div>
</template>

<script>
const { getToken } = useApollo()

export default {
    name: "AlbumAdd",
    data() {
        return {
            album: {
                title: '',
                description: '',
                private: false
            },
            actionInProgress: false,
            notify: {
                message: '',
                type: null
            }
        }
    },
    methods: {
        async onAdd() {
            this.notify = {
                message: '',
                type: null
            }

            if (this.actionInProgress) {
                return
            }

            if (!this.album.title) {
                this.notify = {
                    message: 'Le titre est obligatoire.',
                    type: 'error'
                }

                return
            }

            this.actionInProgress = true

            const query = gql`
                mutation createAlbum($input: CreateAlbumInput!) {
                    createAlbum(input: $input) {
                        title
                    }
                }
            `

            const token = await getToken()

            let decodedToken = atob(token.split('.')[1])
            let tokenParsed = JSON.parse(decodedToken)

            const variables = {
                input: {
                    title: this.album.title,
                    description: this.album.description,
                    private: this.album.private,
                    author: tokenParsed.name
                }
            }

            try {
                const { data : { _rawValue: { createAlbum } } } =  await useAsyncQuery(query, variables)                

                return navigateTo('/admin')
            } catch (e) {
                this.notify = {
                    message: 'Une erreur est survenue.',
                    type: 'error'
                }
            }

            this.actionInProgress = false
        }
    }
}
</script>