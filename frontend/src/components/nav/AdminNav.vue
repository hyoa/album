<template>
  <nav>
      <div class="text-center mt-4">
        <router-link class="no-underline hover:underline" :to="{ name: 'admin_home' }">{{ $t('admin.adminNav.title') }}</router-link>
      </div>
      <div class="px-4 mt-4 grid grid-cols-5 gap-2">
        <NavButton iconClass="bg-blue-500 text-white" link="admin_home">
          <template v-slot:icon>
            <v-icon name="home" scale="2" />
          </template>
        </NavButton>
        <NavButton iconClass="bg-teal-500 text-white" menuName="album">
          <template v-slot:icon>
            <v-icon name="book-open" scale="2" />
          </template>
          <template v-slot:submenu>
            <NavButton iconClass="bg-teal-500 text-white" :small="true" link="admin_album_add">
              <template v-slot:icon>
                <v-icon name="plus"></v-icon>
              </template>
            </NavButton>
            <NavButton iconClass="bg-teal-500 text-white ml-2" :small="true" link="admin_album_list">
              <template v-slot:icon>
                <v-icon name="eye"></v-icon>
              </template>
            </NavButton>
          </template>
        </NavButton>
        <NavButton iconClass="bg-green-500 text-white" menuName="mediaFolder">
          <template v-slot:icon>
            <v-icon name="folder" scale="2" />
          </template>
          <template v-slot:submenu>
            <NavButton iconClass="bg-green-500 text-white" :small="true" link="admin_medias_add">
              <template v-slot:icon>
                <v-icon name="plus"></v-icon>
              </template>
            </NavButton>
            <NavButton iconClass="bg-green-500 text-white ml-2" :small="true" link="admin_medias_folders_list">
              <template v-slot:icon>
                <v-icon name="eye"></v-icon>
              </template>
            </NavButton>
          </template>
        </NavButton>
        <NavButton iconClass="bg-pink-500 text-white" link="admin_users_list">
          <template v-slot:icon>
            <v-icon name="users" scale="2" />
          </template>
        </NavButton>
        <NavButton iconClass="bg-red-400 text-white" link="home">
          <template v-slot:icon>
            <v-icon name="sign-out-alt" scale="2" />
          </template>
        </NavButton>
      </div>
  </nav>
</template>

<script>
import NavButton from './AdminNavButton'
export default {
  name: 'AdminNav',
  components: { NavButton },
  data () {
    return {
      isMenuOpen: false
    }
  }
}
</script>
