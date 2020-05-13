<template>
  <AdminLayout>
    <template v-slot:title>
      <PageTitle :title="$t('admin.mediaFolderList.title')" icon="regular/plus-square" color="green"/>
    </template>
    <ul>
      <ListItem :key="folder" v-for="folder of folders" :to="{ name: 'admin_medias_folder', params: { folder } }" :title="folder" />
    </ul>
  </AdminLayout>
</template>

<script>
import { get, deleteMethod } from '../../../utils/axiosHelper'

import AdminLayout from '../../../components/layout/AdminLayout'
import ListItem from '../../../components/admin/ListItem'
import PageTitle from '../../../components/admin/PageTitle'
export default {
  name: 'MediaFolderList',
  components: { PageTitle, ListItem, AdminLayout },
  data () {
    return {
      folders: []
    }
  },
  async created () {
    const res = await get('medias/folders')
    this.folders = res.data
  },
  methods: {
    onDelete (folderName) {
      deleteMethod(`medias/folder/${folderName}`)
        .then(() => {
          this.folders = this.folders.filter(folder => {
            return folder !== folderName
          })
        })
    }
  }
}
</script>
