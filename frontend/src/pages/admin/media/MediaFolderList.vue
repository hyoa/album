<template>
  <AdminLayout>
    <template v-slot:title>
      <PageTitle :title="$t('admin.mediaFolderList.title')" icon="regular/plus-square" color="bg-green-500"/>
    </template>
    <InputSimple placeholder="vacances" label="Filtre" v-model="filter" type="text"/>
    <ul>
      <ListItem :key="folder.name" v-for="folder of foldersToDisplay" :to="{ name: 'admin_medias_folder', params: { folder: folder.name } }" :title="folder.name" />
    </ul>
  </AdminLayout>
</template>

<script>
import { graphql } from '../../../utils/axiosHelper'

import AdminLayout from '../../../components/layout/AdminLayout'
import ListItem from '../../../components/admin/ListItem'
import PageTitle from '../../../components/admin/PageTitle'
import InputSimple from '../../../components/form/default/InputSimple'

export default {
  name: 'MediaFolderList',
  components: { PageTitle, ListItem, AdminLayout, InputSimple },
  data () {
    return {
      folders: [],
      filter: ''
    }
  },
  async created () {
    const query = `
      query {
        folders: folders(input: {}){
          name
        }
      }
    `

    const { folders } = await graphql(query, 'v3')
    this.folders = folders
  },
  computed: {
    foldersToDisplay () {
      if (this.filter === '') {
        return this.folders
      }

      return this.folders.filter(folder => {
        return folder.name.includes(this.filter)
      })
    }
  }
}
</script>
