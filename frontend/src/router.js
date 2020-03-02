import Vue from 'vue'
import Router from 'vue-router'

import Home from './pages/Home'
import Auth from './pages/Auth'
import UserResetPassword from './pages/UserResetPassword'
import AlbumView from './pages/AlbumView'
import AdminHome from './pages/admin/Home'
import AdminAlbumList from './pages/admin/album/AlbumList'
import AdminAlbumEdit from './pages/admin/album/AlbumEdit'
import AdminAlbumAdd from './pages/admin/album/AlbumAdd'
import AdminMediaFoldersList from './pages/admin/media/MediaFolderList'
import AdminMediaAdd from './pages/admin/media/MediaAdd'
import AdminMediaFolder from './pages/admin/media/MediaFolder'
import AdminUsersList from './pages/admin/user/UsersList'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
      meta: {
        requireAuth: true
      }
    },
    {
      path: '/admin',
      name: 'admin_home',
      component: AdminHome,
      meta: {
        requireAuth: true
      }
    },
    {
      path: '/auth',
      name: 'auth',
      component: Auth,
      meta: {
        rejectIfAuth: true
      }
    },
    {
      path: '/reset-password',
      name: 'resetPassword',
      component: UserResetPassword,
      meta: {
        rejectIfAuth: true
      }
    },
    {
      path: '/album/:slug',
      name: 'album_view',
      component: AlbumView
    },
    {
      path: '/admin/albums',
      name: 'admin_album_list',
      component: AdminAlbumList,
      meta: {
        requireAuth: true
      }
    },
    {
      path: '/admin/album/add',
      name: 'admin_album_add',
      component: AdminAlbumAdd,
      meta: {
        requireAuth: true
      }
    },
    {
      path: '/admin/album/:slug',
      name: 'admin_album_edit',
      component: AdminAlbumEdit,
      meta: {
        requireAuth: true
      }
    },
    {
      path: '/admin/medias/folders',
      name: 'admin_medias_folders_list',
      component: AdminMediaFoldersList,
      meta: {
        requireAuth: true
      }
    },
    {
      path: '/admin/medias/add',
      name: 'admin_medias_add',
      component: AdminMediaAdd,
      meta: {
        requireAuth: true
      }
    },
    {
      path: '/admin/medias/folder/:folder',
      name: 'admin_medias_folder',
      component: AdminMediaFolder,
      meta: {
        requireAuth: true
      }
    },
    {
      path: '/admin/users',
      name: 'admin_users_list',
      component: AdminUsersList,
      meta: {
        requireAuth: true
      }
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import(/* webpackChunkName: "about" */ './views/About.vue')
    }
  ]
})
