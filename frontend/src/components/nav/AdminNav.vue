<template>
  <nav>
    <div class="py-3 px-4 border-grey-lighter bg-primary text-white">
      <div class="float-left" @click="isMenuOpen = !isMenuOpen">
        <i class="material-icons">menu</i>
      </div>
      <div class="text-center">
        <router-link class="no-underline hover:underline" :to="{ name: 'admin_home' }">{{ $t('admin.adminNav.title') }}</router-link>
      </div>
    </div>
    <Sidebar v-if="isMenuOpen">
      <ul>
        <li>
          <div class="bg-light-primary p-2">{{ $t('admin.adminNav.album.title') }}</div>
          <ul class="px-3 py-1">
            <li class="menu-link">
              <router-link :to="{name: 'admin_album_add' }">
                {{ $t('admin.adminNav.album.add') }}
              </router-link>
            </li>
            <li class="menu-link">
              <router-link :to="{ name: 'admin_album_list' }">
                {{ $t('admin.adminNav.album.see') }}
              </router-link>
            </li>
          </ul>
        </li>
        <li>
          <div class="bg-light-primary p-2">{{ $t('admin.adminNav.medias.title') }}</div>
          <ul class="px-3 py-1">
            <li class="menu-link">
              <router-link :to="{ name: 'admin_medias_add' }">
                {{ $t('admin.adminNav.medias.add') }}
              </router-link>
            </li>
            <li class="menu-link">
              <router-link :to="{ name: 'admin_medias_folders_list' }">
                {{ $t('admin.adminNav.medias.see') }}
              </router-link>
            </li>
          </ul>
        </li>
        <li>
          <div class="bg-light-primary p-2">{{ $t('admin.adminNav.users.title') }}</div>
          <ul class="px-3 py-1">
            <li class="menu-link">
              <router-link :to="{ name: 'admin_users_list' }">
                {{ $t('admin.adminNav.users.see') }}
              </router-link>
            </li>
          </ul>
        </li>
        <li>
          <hr class="border-b border-light-primary">
          <router-link :to="{ name: 'home' }" class="flex items-center pl-2">
            {{ $t('admin.adminNav.leave') }} <i class="ml-1 material-icons">exit_to_app</i>
          </router-link>
        </li>
      </ul>
    </Sidebar>
  </nav>
</template>

<style lang="scss" scoped>
  .menu-link {
    @apply py-2;

    a {
      @apply block w-full;
    }
  }
</style>

<script>
import Sidebar from './Sidebar'
export default {
  name: 'AdminNav',
  components: { Sidebar },
  data () {
    return {
      isMenuOpen: false
    }
  }
}
</script>
