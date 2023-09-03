<template>
  <AdminLayout>
    <template v-slot:title>
      <PageTitle :title="folderName" icon="regular/plus-square" color="bg-green-500"/>
    </template>
    <form @submit.prevent="onSubmit" class="mt-3">
      <InputSimple
        v-model="folderName"
        id="name"
        :placeholder="folderName"
        :label="$t('admin.mediaFolder.form.folderName')"
        type="text"
      />
      <FormButton type="submit">
        {{ $t('admin.mediaFolder.form.submit') }}
      </FormButton>
    </form>
    <div v-if="mediasSelected.length" class="sticky top-0 left-0 shadow-md rounded z-20 bg-white">
      <div class="flex justify-between items-center px-4 py-3">
        <div class="title">{{ $t('admin.mediaFolder.mediaSelected.count', { count: mediasSelected.length }) }}</div>
        <div>
          <InputSimple :label="$t('admin.mediaFolder.mediaSelected.form.folder')" :placeholder="$t('admin.mediaFolder.mediaSelected.form.folderPlaceholder')" id="folder" type="text" v-model="newFolderName"/>
          <button
            class="bg-transparent font-semibold py-2 px-4 border hover:border-transparent rounded"
            @click="onChangeFolder"
          >
            {{ $t('admin.mediaFolder.mediaSelected.form.submit') }}
          </button>
        </div>
      </div>
    </div>
    <div class="relative">
        <Grid v-if="medias.length > 0" :editable="true" :can-delete-media="true" :medias="medias" />
    </div>
  </AdminLayout>
</template>

<script>
import { graphql } from '../../../utils/axiosHelper'

import AdminLayout from '../../../components/layout/AdminLayout'
import Grid from '../../../components/grid/Grid2'
import InputSimple from '../../../components/form/default/InputSimple'
import FormButton from '../../../components/form/default/FormButton'
import PageTitle from '../../../components/admin/PageTitle'
export default {
  name: 'MediaFolder.vue',
  components: { PageTitle, FormButton, InputSimple, Grid, AdminLayout },
  data () {
    return {
      medias: [],
      albumName: '',
      toggleAddToAlbum: false,
      folderName: this.$route.params.folder,
      newFolderName: ''
    }
  },
  async created () {
    const query = `
    query {
      folder: folder(input: {name: "${this.folderName}"}){
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
    const { folder } = await graphql(query, 'v3')

    this.medias = folder.medias
  },
  methods: {
    onSubmit () {
      const query = `
        mutation {
          changeFolderName(input: {oldName: "${this.$route.params.folder}", newName: "${this.folderName}"}) {
            name
          }
        }
      `

      graphql(query, 'v3')
        .then(() => {
          this.$notify({ group: 'success', text: this.$t('admin.mediaFolder.notify.submitSuccess') })
        })
        .catch(message => {
          this.$notify({ group: 'error', text: message })
        })
    },
    async onChangeFolder () {
      if (!this.newFolderName || this.newFolderName === '') {
        return
      }

      const query = `
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
        keys: this.$store.state.mediaSelected
      }

      graphql(query, 'v3', variables)
        .then(({ folder: { medias } }) => {
          const mediasToKeep = []

          for (let mediaInFolder of this.medias) {
            let found = false
            for (let mediaRemove of medias) {
              if (mediaInFolder.key === mediaRemove.key) {
                found = true
                break
              }
            }

            if (!found) {
              mediasToKeep.push(mediaInFolder)
            }
          }

          this.medias = mediasToKeep
          this.$notify({ group: 'success', text: this.$t('admin.mediaFolder.notify.moveSuccess') })
        })
        .catch(message => {
          this.$notify({ group: 'error', text: message })
        })
        .finally(() => {
          this.$store.commit('resetMediaSelection')
        })
    }
  },
  computed: {
    mediasSelected () {
      return this.$store.state.mediaSelected
    }
  }
}
</script>

<style scoped>

</style>
