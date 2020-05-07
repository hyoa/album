<template>
  <AdminLayout>
    <PageTitle :title="$t('admin.albumList.title')" />
    <InputForm :label="$t('admin.albumList.form.filter')" placeholder="vacances" v-model="filter" />
    <ul>
      <ListItem
        :key="album.slug"
        v-for="album of filteredAlbum"
        :to="{ name: 'admin_album_edit', params: { slug: album.slug }}"
        :title="album.title"
        :subtitle="$t('admin.albumList.item.createBy', { author: album.author })"
      >
        <div v-if="album.medias">{{ $t('admin.albumList.item.mediasCount', { number: album.medias.length }) }}</div>
        <div v-else>{{ $t('admin.albumList.item.noMedias') }}</div>
      </ListItem>
    </ul>
  </AdminLayout>
</template>

<script>
import { get, deleteMethod } from '../../../utils/axiosHelper'
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
    const response = await get('albums?limit=100&private=1&noMedias=1')
    this.albums = response.data
  },
  methods: {
    onDelete (slug) {
      deleteMethod(`album/${slug}`)
        .then(() => {
          this.albums = this.albums.filter(album => {
            return album.slug !== slug
          })
        })
    }
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
