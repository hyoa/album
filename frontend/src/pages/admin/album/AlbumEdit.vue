<template>
  <AdminLayout>
    <template v-slot:sidebar>
      <Sidebar v-if="addMediaMenuIsVisible" side="right" size="big">
        <section class="border-b px-5 py-3">
          <div class="flex justify-between items-center mb-3">
            <h3 class="text-xl">{{ $t('admin.albumEdit.sidebar.title') }}</h3>
            <span @click="addMediaMenuIsVisible = false" class="text-red-400">Fermer</span>
          </div>
          <InputSimple v-model="folderFilter" :placeholder="$t('admin.albumEdit.sidebar.searchPlaceholder')"/>
          <ul class="foldersList">
            <li class="p-1" :class="{ 'bg-light-primary text-white' : isFolderSelected(folder) }" @click="selectFolder(folder)" v-for="folder of foldersFiltered" :key="folder">{{ folder }}</li>
          </ul>
        </section>
        <section class="px-5 py-3" v-if="folderSelected">
          <div>
            <SimpleAnimateButton
              @click.prevent="onAddToAlbum"
              :status="formStatus.addToAlbum"
              v-if="mediasSelected.length"
            >
              {{ $t('admin.albumEdit.sidebar.addButton') }}
            </SimpleAnimateButton>
            <div class="text-sm" v-else>{{ $t('admin.albumEdit.sidebar.selectMedia') }}</div>
          </div>
          <div v-if="formStatus.loadingMediaAvailable === 'ready'" class="flex flex-wrap mediasToAdd mt-3">
            <div
              @click="selectMedia(media.key)"
              class="w-1/2 p-1"
              v-for="media of mediasAvailable"
              :key="media.id"
              :class="{ 'opacity-25': isMediaSelected(media.key) }"
            >
              <img v-if="media.type === 'image'" class="w-full h-full" :src="media.uris.small" alt="">
              <video v-else preload="metadata">
                <source :src="media.uris.original" type="video/mp4">
              </video>
            </div>
          </div>
          <div v-else>
            <CubicLoader />
          </div>
        </section>
      </Sidebar>
    </template>
    <PageTitle :title="title" />
    <form @submit.prevent="onSubmit" class="mt-3">
      <InputSimple v-model="title" id="title" placeholder="Mon super album" :label="$t('admin.albumEdit.form.title')" type="text" />
      <TextareaSimple v-model="description" id="description" placeholder="Raconte moi une histoire..." :label="$t('admin.albumEdit.form.description')" />
      <CheckboxSimple v-model="isPrivate" :label="$t('admin.albumEdit.form.private')" />
      <div class="flex justify-between">
        <div class="w-1/2 px-1">
          <SimpleAnimateButton :status="formStatus.editAlbum">
            {{ $t('admin.albumEdit.form.submit') }}
          </SimpleAnimateButton>
        </div>

        <div class="w-1/2 px-1">
          <FormButton @click.prevent.stop="toggleAddMediaMenu">
            {{ $t('admin.albumEdit.library') }}
          </FormButton>
        </div>
      </div>
    </form>
    <hr class="border-t border-grey-light">
    <div v-if="mediasSelected.length > 0 && !addMediaMenuIsVisible" class="sticky pin-t pin-l shadow-md rounded z-20 bg-white">
      <div class="flex justify-between items-center px-4 py-3">
        <div class="title">{{ $t('admin.albumEdit.mediaSelected.count', { count: mediasSelected.length}) }}</div>
        <div>
          <SimpleAnimateButton
            @click="onRemoveFromAlbum"
            :status="formStatus.removeFromAlbum"
          >
            {{ $t('admin.albumEdit.mediaSelected.remove') }}
          </SimpleAnimateButton>
        </div>
      </div>
    </div>
    <div class="relative">
    <Grid :medias="medias" :editable="true" :canStar="true" @toggleFavorite="onToggleFavorite" />
    </div>
  </AdminLayout>
</template>

<style scoped lang="scss">
  .mediasToAdd {
    max-height: 40vh;
    overflow-y: scroll;
  }

  .foldersList {
    max-height: 15vh;
    overflow-y: scroll;
  }
</style>

<script>
import { get, post, put } from '../../../utils/axiosHelper'
import errorHelper from '../../../utils/errorHelper'

import AdminLayout from '../../../components/layout/AdminLayout'
import InputSimple from '../../../components/form/default/InputSimple'
import TextareaSimple from '../../../components/form/default/TextareaSimple'
import CheckboxSimple from '../../../components/form/default/CheckboxSimple'
import Grid from '../../../components/grid/Grid'
import PageTitle from '../../../components/admin/PageTitle'
import FormButton from '../../../components/form/default/FormButton'
import SimpleAnimateButton from '../../../components/form/button/SimpleAnimateButton'
import Sidebar from '../../../components/nav/Sidebar'
import CubicLoader from '../../../components/loader/CubicLoader'

export default {
  name: 'AdminAlbumEdit',
  components: { Sidebar, SimpleAnimateButton, FormButton, PageTitle, Grid, CheckboxSimple, TextareaSimple, InputSimple, AdminLayout, CubicLoader },
  data () {
    return {
      title: '',
      description: '',
      isPrivate: false,
      medias: [],
      addMediaMenuIsVisible: false,
      folders: [],
      folderSelected: null,
      mediasAvailable: [],
      formStatus: {
        editAlbum: 'ready',
        removeFromAlbum: 'ready',
        addToAlbum: 'ready',
        loadingMediaAvailable: 'ready'
      },
      folderFilter: ''
    }
  },
  async created () {
    try {
      const { data } = await get(`/album/${this.$route.params.slug}`)
      this.title = data.title
      this.description = data.description
      this.isPrivate = data.private
      this.medias = data.medias

      const res = await get('medias/folders')
      this.folders = res.data
    } catch (e) {
      this.$notify({ group: 'info', text: this.$t('admin.albumEdit.notify.albumDoesNotExist') })

      this.$router.push({ name: 'admin_album_add' })
    }
  },
  methods: {
    onSubmit () {
      this.formStatus.editAlbum = 'pending'
      const data = {
        title: this.title,
        description: this.description,
        private: !!this.isPrivate
      }

      post(`album/${this.$route.params.slug}`, data)
        .then(({ data }) => {
          this.$notify({ group: 'success', text: this.$t('admin.albumEdit.notify.editSuccess') })

          if (data.slug !== this.$route.params.slug) {
            this.$router.push({ name: 'admin_album_edit', params: { slug: data.slug } })
          }
        })
        .catch(({ response }) => {
          let code = null
          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.$notify({ group: 'error', text: this.$t(errorHelper(code)) })
        })
        .finally(() => {
          this.formStatus.editAlbum = 'ready'
        })
    },
    onRemoveFromAlbum () {
      this.formStatus.removeFromAlbum = 'pending'
      const mediasToSend = this.$store.state.mediaSelected

      const mediasObject = mediasToSend.map(mediaToSend => {
        for (let media of this.medias) {
          if (media.key === mediaToSend) {
            return {
              key: media.key,
              author: media.author,
              type: media.type
            }
          }
        }
      })

      post(`album/${this.$route.params.slug}/medias/remove`, mediasObject)
        .then(({ data }) => {
          this.$notify({ group: 'success', text: this.$t('admin.albumEdit.notify.mediaRemoveSuccess') })
          this.$store.commit('resetMediaSelection')
          this.medias = data.medias
        })
        .catch(({ response }) => {
          let code = null
          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.$notify({ group: 'error', text: this.$t(errorHelper(code)) })
        })
        .finally(() => {
          this.formStatus.removeFromAlbum = 'ready'
        })
    },
    async selectFolder (folder) {
      this.formStatus.loadingMediaAvailable = 'pending'
      this.folderSelected = folder
      this.mediasAvailable = []

      const { data } = await get(`medias/folder/${folder}`)

      this.mediasAvailable = data.filter(media => {
        for (let mediaAlbum of this.medias) {
          if (mediaAlbum.key === media.key) {
            return false
          }
        }

        return true
      })

      this.formStatus.loadingMediaAvailable = 'ready'
    },
    isFolderSelected (folder) {
      return folder === this.folderSelected
    },
    toggleAddMediaMenu () {
      this.addMediaMenuIsVisible = !this.addMediaMenuIsVisible

      if (this.addMediaMenuIsVisible) {
        this.$store.commit('resetMediaSelection')
      }
    },
    selectMedia (key) {
      this.$store.commit('toggleMediaSelection', key)
    },
    isMediaSelected (key) {
      return this.$store.state.mediaSelected.includes(key)
    },
    onAddToAlbum () {
      this.formStatus.addToAlbum = 'pending'
      const mediasToSend = this.$store.state.mediaSelected

      const mediasObject = mediasToSend.map(mediaToSend => {
        for (let media of this.mediasAvailable) {
          if (media.key === mediaToSend) {
            return {
              key: media.key,
              author: media.author,
              type: media.type
            }
          }
        }
      })

      post(`album/${this.$route.params.slug}/medias/add`, mediasObject)
        .then(({ data: { medias } }) => {
          this.medias = medias

          this.$notify({ group: 'success', text: this.$t('admin.albumEdit.notify.mediaAddSuccess') })
          this.$store.commit('resetMediaSelection')
        })
        .catch(({ response }) => {
          let code = null
          try {
            code = response.data.code
          } catch (e) {
            code = 999
          }

          this.$notify({ group: 'success', text: this.$t(errorHelper(code)) })
        })
        .finally(() => {
          this.formStatus.addToAlbum = 'ready'

          this.mediasAvailable = this.mediasAvailable.filter((media) => {
            for (let mediaAlbum of this.medias) {
              if (mediaAlbum.key === media.key) {
                return false
              }
            }

            return true
          })
        })
    },
    async onToggleFavorite (media) {
      const url = media.favorite ? `album/${this.$route.params.slug}/favorite/remove` : `album/${this.$route.params.slug}/favorite/add`

      await put(url, { favorite: media.key })

      this.medias.forEach(({ key }, index) => {
        if (key === media.key) {
          this.medias[index].favorite = !media.favorite
        }
      })
    }
  },
  computed: {
    mediasSelected () {
      return this.$store.state.mediaSelected
    },
    foldersFiltered () {
      if (this.folderFilter === '') {
        return this.folders
      }

      return this.folders.filter(folder => folder.toUpperCase().includes(this.folderFilter.toUpperCase()))
    }
  }
}
</script>
