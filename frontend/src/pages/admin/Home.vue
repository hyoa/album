<template>
  <AdminLayout>
    <div class="px-3">
      <AdminCard headColor="bg-indigo-400" footerColor="bg-indigo-200">
        <template v-slot:header>
          {{ $t('admin.home.albumCard.title') }}
        </template>
        <template v-slot:default>
          <div>{{ albumsData.publicCount }} {{ $t('admin.home.albumCard.publicCount') }}</div>
          <div>{{ albumsData.privateCount }} {{ $t('admin.home.albumCard.privateCount') }}</div>
        </template>
        <template v-slot:footer>
          <router-link :to="{ name: 'admin_album_list' }" class="">
            <IconLabel icon="remove_red_eye" :label="$t('admin.home.albumCard.see')" />
          </router-link>
          <router-link :to="{name: 'admin_album_add' }">
            <IconLabel icon="add" :label="$t('admin.home.albumCard.add')" />
          </router-link>
        </template>
      </AdminCard>
      <AdminCard class="mt-8" headColor="bg-green-400" footerColor="bg-green-200">
        <template v-slot:header>
          {{ $t('admin.home.mediaCard.title') }}
        </template>
        <template v-slot:default>
          <div>{{ mediasData.imagesCount }} {{ $t('admin.home.mediaCard.photos') }}</div>
          <div>{{ mediasData.videosCount }} {{ $t('admin.home.mediaCard.videos') }}</div>
        </template>
        <template v-slot:footer>
          <router-link :to="{ name: 'admin_medias_folders_list' }">
            <IconLabel icon="remove_red_eye" :label="$t('admin.home.mediaCard.see')" />
          </router-link>
          <router-link :to="{ name: 'admin_medias_add' }">
            <IconLabel icon="add" :label="$t('admin.home.mediaCard.add')" />
          </router-link>
        </template>
      </AdminCard>
      <AdminCard class="mt-8" headColor="bg-pink-400" footerColor="bg-pink-200">
        <template v-slot:header>
          {{ $t('admin.home.userCard.title') }}
        </template>
        <template v-slot:default>
          <div>{{ usersData.total }} {{ $t('admin.home.userCard.count') }}</div>
          <div>{{ usersData.unverifiedCount }} {{ $t('admin.home.userCard.waitingValidation') }}</div>
        </template>
        <template v-slot:footer>
          <router-link :to="{ name: 'admin_users_list' }">
            <IconLabel icon="remove_red_eye" :label="$t('admin.home.userCard.see')" />
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
