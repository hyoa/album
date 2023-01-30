<template>
  <AdminLayout>
    <template v-slot:title>
      <PageTitle :title="$t('admin.mediaAdd.title')" icon="regular/plus-square" color="bg-green-500"/>
    </template>
    <div v-if="upload.state === 'running'">
      <div>
        <h3>{{ $t('admin.mediaAdd.uploadRunning.title') }}</h3>
        <div v-if="upload.timeRemaining">{{ $t('admin.mediaAdd.uploadRunning.timeRemaining', { time: upload.timeRemaining }) }}</div>
        <div>{{ $t('admin.mediaAdd.uploadRunning.totalUploaded', { count: upload.total }) }}</div>
        <div>{{ $t('admin.mediaAdd.uploadRunning.successUploaded', { count: upload.uploaded }) }}</div>
        <div>{{ $t('admin.mediaAdd.uploadRunning.failUploaded', { count: upload.failed }) }}</div>
      </div>
      <ul>
        <li v-for="media of upload.medias" :key="media.key">
          <div v-if="media.status  === 'uploading'">
            <div>{{ media.key }}</div>
            <div class="shadow w-full bg-white mt-2 rounded">
              <div class="bg-blue-500 text-xs leading-none py-1 text-center text-white rounded" :style="`width: ${media.progress}%`">
                {{ media.progress }}%
              </div>
            </div>
          </div>
          <div class="flex items-center" v-if="media.status === 'success'">
            <div class="text-green-500"><v-icon class="mr-2" name="check"></v-icon>{{ media.key }}</div>
          </div>
          <div class="flex items-center" v-if="media.status === 'failed'">
            <div class="text-red-500"><v-icon class="mr-2" name="times"></v-icon>{{ media.key }}</div>
          </div>
        </li>
      </ul>
    </div>
    <div v-if="upload.state === 'error'">
      <div class="bg-red-500 text-white">{{ $t('admin.mediaAdd.filesNotUploaded') }}</div>
      <ul>
        <li v-for="media of upload.failedMedias" :key="media.key">
          <div class="flex items-center">
            <div>{{ media.key }}</div>
          </div>
        </li>
      </ul>
      <hr class="border-4 mb-8">
    </div>
    <form @submit.prevent="uploadMedias" v-show="upload.state !== 'running'">
      <AutoComplete v-model="folder" id="folder" placeholder="Un super dossier !" :label="$t('admin.mediaAdd.form.folder')" type="text" entity="folder" :allow-no-call="true"/>
      <CheckboxSimple v-model="linkToAlbum" :label="$t('admin.mediaAdd.form.linkToAlbum')" />
      <AutoComplete v-if="linkToAlbum" v-model="album" id="album" placeholder="Un album" :label="$t('admin.mediaAdd.form.album')" type="text" entity="album" :allow-no-call="true"/>
      <div class="mb-3" v-if="folder.trim() !== ''">
        <label for="files">{{ $t('admin.mediaAdd.form.media') }}</label>
        <div
          class="shadow-inner bg-white hover:bg-white border-2 border-gray-lighter rounded-sm relative h-16"
        >
          <div class="shadow-inner absolute w-full h-full pin-t flex justify-center items-center">
            <span>{{ $t('admin.mediaAdd.form.dragAndDrop') }}</span>
          </div>
          <input
            id="files"
            type="file"
            class="opacity-0 w-full h-full"
            accept="image/png,image/jpg,video/mp4,image/jpeg"
            multiple
            @change="loadMedia"
          />
        </div>
      </div>
    </form>
  </AdminLayout>
</template>

<style lang="scss" scoped>
  ul {
    li {
      @apply border-b border-darker-primary py-4;

      &:last-child {
        border: none
      }
    }
  }
</style>

<script>
import axios from 'axios'
import prettyMs from 'pretty-ms'
import NoSleep from 'nosleep.js'
import { graphql } from '../../../utils/axiosHelper'
import AdminLayout from '../../../components/layout/AdminLayout'
import AutoComplete from '../../../components/form/default/AutoComplete'
import CheckboxSimple from '../../../components/form/default/CheckboxSimple'
import PageTitle from '../../../components/admin/PageTitle'

const uploadData = () => ({
  state: null,
  progress: 0,
  failed: 0,
  medias: [],
  total: 0,
  uploaded: 0,
  timeRemaining: null,
  failedMedias: []
})

export default {
  name: 'MediaAdd',
  components: { PageTitle, AutoComplete, AdminLayout, CheckboxSimple },
  data () {
    return {
      folder: '',
      album: '',
      medias: [],
      upload: uploadData(),
      noSleep: new NoSleep(),
      linkToAlbum: false
    }
  },
  methods: {
    async loadMedia ({ target }) {
      if (!target.validity.valid && target.files.length <= 0) {
        return null
      }

      if (target.files.length > 100) {
        this.$notify({ group: 'warning', text: this.$t('admin.mediaAdd.notify.tooManyMedia') })
        return null
      }

      this.medias = [...this.medias, ...target.files]

      await this.uploadMedias()
    },
    async uploadMedias () {
      this.upload = uploadData()
      this.noSleep.enable()
      this.upload.state = null
      this.upload.progress = 0
      const promises = []
      this.upload.total = this.medias.length

      let totalSize = 0

      let mediasConfigToUpload = []
      const author = this.$store.state.token.name
      const folder = this.folder.replace('new|', '').replace(/[^a-zA-Z0-9. ]/g, '').normalize('NFD').replace(/[\u0300-\u036f]/g, '')

      for (let media of this.medias) {
        totalSize += media.size
        const key = media.name.replace(/[^a-zA-Z0-9.]/g, '').normalize('NFD').replace(/[\u0300-\u036f]/g, '')

        this.upload.medias.push({ key, progress: 0, uploadedSize: 0, status: 'uploading', kind: media.type.includes('video') ? 'VIDEO' : 'PHOTO' })

        const putConfig = {
          onUploadProgress: progressEvent => {
            for (let i in this.upload.medias) {
              if (this.upload.medias[i].key === key) {
                this.upload.medias[i].progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
                this.upload.medias[i].uploadedSize = progressEvent.loaded
                break
              }
            }
          },
          headers: {
            'Content-Type': media.type,
            'Cache-Control': 'max-age=43200'
          }
        }

        mediasConfigToUpload.push({
          key,
          media,
          putConfig
        })
      }

      const querySignedUri = `
        query($medias: [GetIngestMediaInput!]!) {
          medias: ingest (input: {medias: $medias}) {
            key
            signedUri
          }
        }
      `

      const variablesSignedUri = {
        medias: mediasConfigToUpload.map(({ key, media: { type } }) => {
          return { kind: type.includes('video') ? 'VIDEO' : 'PHOTO', key: key }
        })
      }

      const signedUrisResponse = await graphql(querySignedUri, 'v3', variablesSignedUri)

      for (let media of signedUrisResponse.medias) {
        for (let config of mediasConfigToUpload) {
          if (config.key === media.key) {
            const req = axios
              .put(media.signedUri, config.media, config.putConfig)
              .then(() => {
                this.upload.uploaded++
                const index = this.upload.medias.findIndex(m => {
                  return m.key === config.key
                })

                this.upload.medias[index].status = 'success'
              })
              .catch(() => {
                const index = this.upload.medias.findIndex(m => {
                  return m.key === config.key
                })

                this.upload.medias[index].status = 'failed'
                this.upload.failed++
                this.upload.failedMedias.push(this.upload.medias[index])
              })

            promises.push(req)

            break
          }
        }
      }

      this.upload.state = 'running'
      let uploadSpeed = null

      const reduceAdd = (accumulator, currentValue) => accumulator + currentValue.uploadedSize

      const intervalTimeEstimation = setInterval(() => {
        const uploadedSizeRef = this.upload.medias.reduce(reduceAdd, 0)
        setTimeout(() => {
          const uploadedSizeAfterTime = this.upload.medias.reduce(reduceAdd, 0) - uploadedSizeRef
          uploadSpeed = uploadedSizeAfterTime / 1000
        }, 1000)
      }, 5000)

      const intervalDisplayTime = setInterval(() => {
        if (uploadSpeed === null) {
          this.upload.timeRemaining = this.$t('admin.mediaAdd.uploadRunning.timeRemaingingCalculation')
        } else {
          const sizeUploaded = this.upload.medias.reduce(reduceAdd, 0)

          const timeRemaining = (totalSize - sizeUploaded) / uploadSpeed

          if (timeRemaining < 1000) {
            this.upload.timeRemaining = this.$t('admin.mediaAdd.uploadRunning.almostDone')
          }

          this.upload.timeRemaining = `~ ${prettyMs(timeRemaining)}`
        }
      }, 1000)

      Promise.all(promises).then(async () => {
        this.medias = []
        try {
          if (this.upload.uploaded > 0) {
            const queryIngest = `
              mutation($medias: [PutIngestMediaInput!]!) {
                ingest(input: {medias: $medias}) {
                  key
                  status
                }
              }
            `

            const mediasUploaded = this.upload.medias.filter(media => media.status === 'success').map(media => { return { author, folder, key: media.key, kind: media.kind } })
            const variablesIngest = {
              medias: mediasUploaded
            }

            await graphql(queryIngest, 'v3', variablesIngest)

            if (this.album) {
              let slug = this.album
              if (slug.includes('new|')) {
                const albumName = slug.replace('new|', '')
                const queryCreate = `
                  mutation {
                    album: createAlbum(input: {title: "${albumName}", author: "${this.$store.state.token.name}", description: "", private: true}) {
                      title
                      slug
                    }
                  }
                `

                const { album } = await graphql(queryCreate, 'v3')
                slug = album.slug
              }

              const queryAdd = `
                mutation ($medias: [MediaAlbumInput!]!) {
                  album: updateAlbumMedias(input: {slug: "${slug}", medias: $medias, action: ADD}) {
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

              const variablesAdd = {
                medias: mediasUploaded.map(({ key, author, kind }) => { return { key, author, kind } })
              }

              await graphql(queryAdd, 'v3', variablesAdd)
            }
          }

          if (this.upload.failed === 0 && this.upload.uploaded > 0) {
            this.$notify({ group: 'success', text: this.$t('admin.mediaAdd.notify.uploadSuccess') })

            this.upload.state = 'end'
            this.upload = uploadData()
          } else if (this.upload.failed > 0 && this.upload.uploaded > 0) {
            this.upload.state = 'error'
            this.$notify({ group: 'warning', text: this.$t('admin.mediaAdd.notify.uploadPartial') })
          } else {
            this.upload.state = 'error'
            this.$notify({ group: 'error', text: this.$t('admin.mediaAdd.notify.uploadFailed') })
          }
        } catch (e) {
          console.log(e)
          this.upload.state = 'error'
          this.$notify({ group: 'error', text: this.$t('admin.mediaAdd.notify.uploadFailed') })
        }

        clearInterval(intervalTimeEstimation)
        clearInterval(intervalDisplayTime)
        this.noSleep.disable()
      })
    }
  }
}
</script>
