<template>
    <div class="px-4">
        <section class="shadow-md rounded-md p-4">
            <div>
                <h2 class="text-xl">Albums</h2>
            </div>
            <div class="divider"></div> 
            <div class="mt-2">
                <p>{{ albumsData.publicCount }} albums publics</p>
                <p>{{ albumsData.privateCount }} albums privé</p>
            </div>
            <div class="flex justify-between mt-2">
                <NuxtLink to="/admin/albums" class="btn btn-sm btn-secondary">Voir</NuxtLink>
                <NuxtLink to="/admin/album/add" class="btn btn-sm btn-primary">Ajouter</NuxtLink>
            </div>
        </section>
        <section class="shadow-md rounded-md p-4 mt-8">
            <div>
                <h2 class="text-xl">Medias</h2>
            </div>
            <div class="divider"></div> 
            <div class="mt-2">
                <p>{{ mediasData.imagesCount }} photos</p>
                <p>{{ mediasData.videosCount }} vidéos</p>
            </div>
            <div class="flex justify-between mt-2">
                <NuxtLink to="/admin/medias" class="btn btn-sm btn-secondary">Voir</NuxtLink>
            </div>
        </section>
        <section class="shadow-md rounded-md p-4 mt-8">
            <div>
                <h2 class="text-xl">Utilisateurs</h2>
            </div>
            <div class="divider"></div> 
            <div class="mt-2">
                <p>{{ usersData.total }} utilisateurs</p>
                <p>{{ usersData.unverifiedCount }} en attente de validation</p>
            </div>
            <div class="flex justify-between mt-2">
                <NuxtLink to="/admin/users" class="btn btn-sm btn-secondary">Voir</NuxtLink>
            </div>
        </section>
    </div>
</template>

<script setup>
definePageMeta({
middleware: ['auth', 'admin']
})
</script>


<script>
export default {
    layout: 'admin',
    head() {
        return {
        }
    },
    data() {
        return {
            albumsData: {
                publicCount: 0,
                privateCount: 0
            },
            mediasData: {
                imagesCount: 0,
                videosCount: 0
            },
            usersData: {
                total: 0,
                unverifiedCount: 0
            }
        }
    },
    async created() {
        const queryAlbum = gql`
            query {
                albums: albums(input: {limit: 1000, includePrivate: true, includeNoMedias: true}) {
                    private
                }
            }
        `

        const { data: { _rawValue: { albums } } } = await useAsyncQuery(queryAlbum, {})
        for (let album of albums) {
            if (album.private) {
                this.albumsData.privateCount++
            } else {
                this.albumsData.publicCount++
            }
        }

        const queryFolder = gql`
            query {
                folders: folders(input: {}){
                medias {
                    kind
                }
                }
            }
        `

        const { data: { _rawValue: { folders } } } = await useAsyncQuery(queryFolder, {})
        for (let folder of folders) {
            for (let media of folder.medias) {
                if (media.kind === 'PHOTO') {
                    this.mediasData.imagesCount++
                } else if (media.kind === 'VIDEO') {
                    this.mediasData.videosCount++
                }
            }
        }

        const queryUser = gql`
            query {
                users: users{
                    role
                }
            }
        `

        const { data: { _rawValue: { users } } } = await useAsyncQuery(queryUser, {})
        this.usersData.total = users.length
        for (let user of users) {
            if (user.role === 'UNIDENTIFIED') {
                this.usersData.unverifiedCount++
            }
        }
    }
}
</script>