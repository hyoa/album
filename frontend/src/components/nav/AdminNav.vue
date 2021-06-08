<template>
  <nav class="border-b pt-8 pb-4 shadow-md">
      <div class="px-4 grid grid-cols-5 gap-2">
        <NavButton iconClass="bg-gray-500" link="admin_home">
          <template v-slot:icon>
            <v-icon name="home" scale="2" :inverse="true"/>
          </template>
        </NavButton>
        <NavButton iconClass="bg-indigo-500" menuName="album">
          <template v-slot:icon>
            <v-icon name="book-open" scale="2" :inverse="true"/>
          </template>
          <template v-slot:submenu>
            <NavButton iconClass="bg-indigo-500" :small="true" link="admin_album_add">
              <template v-slot:icon>
                <v-icon name="plus" :inverse="true"></v-icon>
              </template>
            </NavButton>
            <NavButton iconClass="bg-indigo-500 ml-2" :small="true" link="admin_album_list">
              <template v-slot:icon>
                <v-icon name="eye" :inverse="true"></v-icon>
              </template>
            </NavButton>
          </template>
        </NavButton>
        <NavButton iconClass="bg-green-500" menuName="mediaFolder">
          <template v-slot:icon>
            <v-icon name="folder" scale="2" :inverse="true"/>
          </template>
          <template v-slot:submenu>
            <NavButton iconClass="bg-green-500" :small="true" link="admin_medias_add">
              <template v-slot:icon>
                <v-icon name="plus" :inverse="true"></v-icon>
              </template>
            </NavButton>
            <NavButton iconClass="bg-green-500 ml-2" :small="true" link="admin_medias_folders_list">
              <template v-slot:icon>
                <v-icon name="eye" :inverse="true"></v-icon>
              </template>
            </NavButton>
          </template>
        </NavButton>
        <NavButton iconClass="bg-pink-500" link="admin_users_list">
          <template v-slot:icon>
            <v-icon name="users" scale="2" :inverse="true"/>
          </template>
        </NavButton>
        <NavButton iconClass="bg-red-400" link="home">
          <template v-slot:icon>
            <v-icon name="sign-out-alt" scale="2" :inverse="true"/>
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
