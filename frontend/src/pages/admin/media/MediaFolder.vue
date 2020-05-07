<template>
  <AdminLayout>
    <PageTitle :title="folderName" />
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
          this.$notify({ group: 'success', text: this.$t('admin.mediaFolder.notify.submitSuccess') })
        })
        .catch(({ response }) => {
          let code = null
          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.$notify({ group: 'error', text: this.$t(errorHelper(code)) })
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
          this.$notify({ group: 'success', text: this.$t('admin.mediaFolder.notify.moveSuccess') })
        })
        .catch(({ response }) => {
          let code = null

          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.$notify({ group: 'error', text: this.$t(errorHelper(code)) })
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
