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
            <li class="p-1" :class="{ 'bg-light-primary text-white' : isFolderSelected(folder.name) }" @click="selectFolder(folder.name)" v-for="folder of foldersFiltered" :key="folder.name">{{ folder.name }}</li>
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
              :key="media.key"
              :class="{ 'opacity-25': isMediaSelected(media.key) }"
            >
              <img v-if="media.kind === 'PHOTO'" class="w-full h-full" :src="media.urls.small" alt="">
              <video v-else preload="metadata">
                <source :src="media.urls.small" type="video/mp4">
              </video>
            </div>
          </div>
          <div v-else>
            <CubicLoader />
          </div>
        </section>
      </Sidebar>
    </template>
    <template v-slot:title>
      <PageTitle :title="title" icon="regular/edit" color="bg-indigo-500" />
    </template>
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
import { graphql } from '../../../utils/axiosHelper'

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
      const queryAlbum = `
        query {
          album: album(input: {slug: "${this.$route.params.slug}"}) {
            title
            slug
            description
            author
            creationDate
            medias {
              key
              kind
              favorite
              author
              urls {
                small
                medium
                large
              }
            }
          }
        }
      `

      const { album } = await graphql(queryAlbum, 'v3')
      this.title = album.title
      this.description = album.description
      this.isPrivate = album.private
      this.medias = album.medias ?? []

      const queryFolder = `
        query {
          folders: folders(input: {}){
            name
            medias {
              key
              kind
              author
              urls {
                small
              }
            }
          }
        }
      `
      const { folders } = await graphql(queryFolder, 'v3')
      this.folders = folders
    } catch (e) {
      this.$notify({ group: 'info', text: this.$t('admin.albumEdit.notify.albumDoesNotExist') })

      this.$router.push({ name: 'admin_album_add' })
    }
  },
  methods: {
    onSubmit () {
      this.formStatus.editAlbum = 'pending'

      const query = `
        mutation {
          album: updateAlbum(input: {title: "${this.title}", author: "${this.$store.state.token.name}", description: "${this.description}", private: ${!!this.isPrivate}, slug: "${this.$route.params.slug}"}) {
            slug
          }
        }
      `

      graphql(query, 'v3')
        .then(({ album }) => {
          this.$notify({ group: 'success', text: this.$t('admin.albumEdit.notify.editSuccess') })

          if (album.slug !== this.$route.params.slug) {
            this.$router.push({ name: 'admin_album_edit', params: { slug: album.slug } })
          }
        })
        .catch(message => {
          this.$notify({ group: 'error', text: message })
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
              kind: media.kind
            }
          }
        }
      })

      const query = `
        mutation ($medias: [MediaAlbumInput!]!) {
          album: updateAlbumMedias(input: {slug: "${this.$route.params.slug}", medias: $medias, action: REMOVE}) {
            medias {
              key
              author
              urls {
                small
              }
              kind
              favorite
            }
          }
        }
      `

      const variables = {
        medias: mediasObject
      }

      graphql(query, 'v3', variables)
        .then(({ album: { medias } }) => {
          this.$notify({ group: 'success', text: this.$t('admin.albumEdit.notify.mediaRemoveSuccess') })
          this.$store.commit('resetMediaSelection')
          this.medias = medias
        })
        .catch(message => {
          this.$notify({ group: 'error', text: message })
        })
        .finally(() => {
          this.formStatus.removeFromAlbum = 'ready'
        })
    },
    async selectFolder (folderName) {
      this.formStatus.loadingMediaAvailable = 'pending'
      this.folderSelected = folderName
      this.mediasAvailable = []

      const folderSelected = this.folders.find(folder => folder.name === folderName)

      this.mediasAvailable = folderSelected.medias.filter(media => {
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
              kind: media.kind
            }
          }
        }
      })

      const query = `
        mutation ($medias: [MediaAlbumInput!]!) {
          album: updateAlbumMedias(input: {slug: "${this.$route.params.slug}", medias: $medias, action: ADD}) {
            medias {
              key
              urls {
                small
              }
              kind
              favorite
            }
          }
        }
      `

      const variables = {
        medias: mediasObject
      }

      graphql(query, 'v3', variables)
        .then(({ album: { medias } }) => {
          this.medias = medias

          this.$notify({ group: 'success', text: this.$t('admin.albumEdit.notify.mediaAddSuccess') })
          this.$store.commit('resetMediaSelection')
        })
        .catch(message => {
          this.$notify({ group: 'error', text: message })
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
      const query = `
        mutation {
          updateAlbumFavorite(input: {slug: "${this.$route.params.slug}", mediaKey: "${media.key}"}) {
            title
          }
        }
      `

      await graphql(query, 'v3')

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

      return this.folders.filter(folder => folder.name.toUpperCase().includes(this.folderFilter.toUpperCase()))
    }
  }
}
</script>
