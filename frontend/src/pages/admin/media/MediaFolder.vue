<template>
  <AdminLayout>
    <PageTitle :title="folderName" />
    <form @submit.prevent="onSubmit" class="mt-3">
      <InputSimple
        v-model="folderName"
        id="name"
        :placeholder="folderName"
        label="Nom du dossier"
        type="text"
      />
      <FormButton type="submit">
        Mettre à jour
      </FormButton>
    </form>
    <div v-if="mediasSelected.length" class="sticky top-0 left-0 shadow-md rounded z-20 bg-white">
      <div class="flex justify-between items-center px-4 py-3">
        <div class="title">{{ mediasSelected.length }} média(s) sélectionné(s)</div>
        <div>
          <InputSimple label="Dossier" placeholder="Nouveau dossier" id="folder" type="text" v-model="newFolderName"/>
          <button
            class="bg-transparent font-semibold py-2 px-4 border hover:border-transparent rounded"
            @click="onChangeFolder"
          >
            Changer de dossier
          </button>
        </div>
      </div>
    </div>
    <Grid :editable="true" :can-delete-media="true" :medias="medias" />
  </AdminLayout>
</template>

<script>
import { get, post } from '../../../utils/axiosHelper'
import errorHelper from '../../../utils/errorHelper'

import AdminLayout from '../../../components/layout/AdminLayout'
import Grid from '../../../components/grid/Grid'
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
    const folderName = this.$route.params.folder
    const res = await get(`medias/folder/${folderName}`)

    this.medias = res.data
  },
  methods: {
    onSubmit () {
      post('medias/folder/name', { folderToUpdate: this.$route.params.folder, newFolderName: this.folderName })
        .then(() => {
          this.$notify({ group: 'success', text: 'Enregistrement effectué' })
        })
        .catch(({ response }) => {
          let code = null
          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.$notify({ group: 'error', text: errorHelper(code) })
        })
    },
    async onChangeFolder () {
      if (!this.newFolderName || this.newFolderName === '') {
        return
      }

      await post(`medias/many/folder/name`, { folderName: this.newFolderName, medias: this.$store.state.mediaSelected })

      this.$store.commit('resetMediaSelection')

      get(`medias/folder/${this.folderName}`)
        .then(({ data }) => {
          this.medias = data
          this.$notify({ group: 'success', text: 'Les médias sélectionnés ont correctement été déplacés' })
        })
        .catch(({ response }) => {
          let code = null

          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.$notify({ group: 'error', text: errorHelper(code) })
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
