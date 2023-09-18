<template>
    <section class="px-4">
        <div class="text-sm breadcrumbs">
            <ul>
                <li><NuxtLink to="/admin">Admin</NuxtLink></li> 
                <li class="font-bold">Dossiers</li>
            </ul>
        </div>
    </section>
    <section class="px-4">
        <div class="w-full">
            <input v-model="searchFolder" type="text" placeholder="Chercher" class="input input-bordered w-full max-w-xs" />
        </div>
    </section>
    <section class="mt-6 px-4">
        <div
            v-for="folder in foldersFiltered"
            :key="folder.name"
        >
            <div class="shadow-md rounded-md p-2">
                <NuxtLink :to="`media/${folder.name}`" class="flex justify-between">
                    <div class="text-xl">{{ folder.name }}</div>
                    <div>{{ folder.medias.length  }} medias</div>
                </NuxtLink>
            </div>
        </div>
    </section>
</template>

<script>
export default {
    name: "Medias",
    data() {
        return {
            searchFolder: '',
            folders: [],
        }
    },
    async created () {
        const queryFolder = gql`
            query {
                folders: folders(input: {}){
                    name
                    medias {
                        key
                    }
                }
            }
        `

        try {
            const { data: { _rawValue: { folders } } } = await useAsyncQuery(queryFolder)
            this.folders = folders
        } catch (e) {
            
        }
    },
    computed: {
        foldersFiltered () {
            return this.folders.filter(folder => folder.name.includes(this.searchFolder))
        }
    }
}
</script>