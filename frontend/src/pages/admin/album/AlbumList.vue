<template>
  <AdminLayout>
    <PageTitle title="Liste des albums" />
    <InputForm label="Filtrer" placeholder="vacances" v-model="filter" />
    <ul>
      <ListItem
        :key="album.slug"
        v-for="album of filteredAlbum"
        :to="{ name: 'admin_album_edit', params: { slug: album.slug }}"
        :title="album.title"
        :subtitle="`Créer par ${album.author}`"
      >
        <div v-if="album.medias">{{ album.medias.length }} médias </div>
        <div v-else>Aucun médias</div>
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
