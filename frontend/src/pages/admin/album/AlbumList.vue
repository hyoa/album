<template>
  <AdminLayout>
    <template v-slot:title>
      <PageTitle :title="$t('admin.albumList.title')" icon="list" color="bg-indigo-500"/>
    </template>
    <InputForm :label="$t('admin.albumList.form.filter')" placeholder="vacances" v-model="filter" />
    <ul>
      <ListItem
        :key="album.slug"
        v-for="album of filteredAlbum"
        :to="{ name: 'admin_album_edit', params: { slug: album.slug }}"
        :title="album.title"
        :subtitle="$t('admin.albumList.item.createBy', { author: album.author })"
        :iconTitle="album.private ? 'lock' : null"
      >
        <div v-if="album.medias">{{ $t('admin.albumList.item.mediasCount', { number: album.medias.length }) }}</div>
        <div v-else>{{ $t('admin.albumList.item.noMedias') }}</div>
      </ListItem>
    </ul>
  </AdminLayout>
</template>

<script>
import { graphql } from '../../../utils/axiosHelper'
import AdminLayout from '../../../components/layout/AdminLayout'
import ListItem from '../../../components/admin/ListItem'
import PageTitle from '../../../components/admin/PageTitle'
import InputForm from '../../../components/form/default/InputSimple'

export default {
  name: 'AdminAlbumList',
  components: { PageTitle, ListItem, AdminLayout, InputForm },
  data () {
    return {
      albums: [],
      filter: ''
    }
  },
  async created () {
    const query = `
      query {
        albums: albums(input: {includePrivate: true, includeNoMedias: true, limit: 1000}) {
          title
          slug
          author
          medias {
            kind
          }
        }
      }
    `

    const response = await graphql(query, 'v3')
    this.albums = response.albums
  },
  computed: {
    filteredAlbum () {
      if (this.filter === '') {
        return this.albums
      }

      return this.albums.filter(album => album.title.toUpperCase().includes(this.filter.toUpperCase()))
    }
  }
}
</script>
