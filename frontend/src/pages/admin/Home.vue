<template>
  <AdminLayout>
    <div class="px-3">
      <AdminCard>
        <template v-slot:header>
          Albums
        </template>
        <template v-slot:default>
          <div>{{ albumsData.publicCount }} albums publics</div>
          <div>{{ albumsData.privateCount }} albums privés</div>
        </template>
        <template v-slot:footer>
          <router-link :to="{ name: 'admin_album_list' }" class="">
            <IconLabel icon="remove_red_eye" label="Voir" />
          </router-link>
          <router-link :to="{name: 'admin_album_add' }">
            <IconLabel icon="add" label="Créer" />
          </router-link>
        </template>
      </AdminCard>
      <AdminCard class="mt-8">
        <template v-slot:header>
          Médias
        </template>
        <template v-slot:default>
          <div>{{ mediasData.imagesCount }} photos</div>
          <div>{{ mediasData.videosCount }} vidéos</div>
        </template>
        <template v-slot:footer>
          <router-link :to="{ name: 'admin_medias_folders_list' }">
            <IconLabel icon="remove_red_eye" label="Voir" />
          </router-link>
          <router-link :to="{ name: 'admin_medias_add' }">
            <IconLabel icon="add" label="Ajouter" />
          </router-link>
        </template>
      </AdminCard>
      <AdminCard class="mt-8">
        <template v-slot:header>
          Utilisateurs
        </template>
        <template v-slot:default>
          <div>{{ usersData.total }} utilisateurs</div>
          <div>{{ usersData.unverifiedCount }} en attente de validation</div>
        </template>
        <template v-slot:footer>
          <router-link :to="{ name: 'admin_users_list' }">
            <IconLabel icon="remove_red_eye" label="Voir" />
          </router-link>
        </template>
      </AdminCard>
    </div>
  </AdminLayout>
</template>

<script>
import { get } from '../../utils/axiosHelper'
import AdminLayout from '../../components/layout/AdminLayout'
import AdminCard from '../../components/admin/Card'
import IconLabel from '../../components/icon/IconLabel'

export default {
  name: 'AdminHome',
  components: { IconLabel, AdminCard, AdminLayout },
  data () {
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
  async created () {
    const albumsData = await get('albums/resume')
    const mediasData = await get('medias/resume')
    const usersData = await get('users/resume')

    this.albumsData = albumsData.data
    this.mediasData = mediasData.data
    this.usersData = usersData.data
  }
}
</script>
