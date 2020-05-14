<template>
  <AdminLayout>
    <template v-slot:title>
      <PageTitle :title="$t('admin.mediaAdd.title')" icon="regular/plus-square" color="bg-green-400"/>
    </template>
    <div v-if="upload.state === 'running'">
      <div>
        <h3>{{ $t('admin.mediaAdd.uploadRunning.title') }}</h3>
        <div v-if="upload.timeRemaining">{{ $t('admin.mediaAdd.uploadRunning.timeRemaining', { time: upload.timeRemaining }) }}</div>
        <div>{{ $t('admin.mediaAdd.uploadRunning.totalUploaded', { count: upload.total }) }}</div>
        <div>{{ $t('admin.mediaAdd.uploadRunning.successUploaded', { time: upload.uploaded }) }}</div>
        <div>{{ $t('admin.mediaAdd.uploadRunning.failUploaded', { time: upload.failed }) }}</div>
      </div>
      <ul>
        <li v-for="media of upload.medias" :key="media.key">
          <div>{{ media.key }}</div>
          <div class="shadow w-full bg-white mt-2 rounded">
            <div class="bg-blue-500 text-xs leading-none py-1 text-center text-white rounded" :style="`width: ${media.progress}%`">
              {{ media.progress }}%
            </div>
          </div>
        </li>
      </ul>
    </div>
    <form @submit.prevent="uploadMedias" v-show="upload.state !== 'running'">
      <AutoComplete v-model="folder" id="folder" placeholder="Un super dossier !" :label="$t('admin.mediaAdd.form.folder')" type="text" endpoint="medias/folders/autocomplete" :allow-no-call="true"/>
      <CheckboxSimple v-model="linkToAlbum" :label="$t('admin.mediaAdd.form.linkToAlbum')" />
      <AutoComplete v-if="linkToAlbum" v-model="album" id="album" placeholder="Un album" :label="$t('admin.mediaAdd.form.album')" type="text" endpoint="albums/autocomplete" :allow-no-call="true"/>
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
import { post } from '../../../utils/axiosHelper'
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
  timeRemaining: null
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
      this.noSleep.enable()
      this.upload.state = null
      this.upload.progress = 0
      const promises = []
      this.upload.total = this.medias.length

      let totalSize = 0

      let mediasConfigToUpload = []

      for (let media of this.medias) {
        totalSize += media.size
        const key = media.name.replace(/[^a-zA-Z0-9.]/g, '').normalize('NFD').replace(/[\u0300-\u036f]/g, '')
        const metadata = {
          author: this.$store.state.token.name,
          folder: this.folder
        }

        if (this.album) {
          metadata.album = this.album
        }

        this.upload.medias.push({ key, progress: 0, uploadedSize: 0 })

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
          putConfig,
          metadata
        })
      }

      const payload = {
        medias: mediasConfigToUpload.map(({ key, media: { type }, metadata }) => { return { type, file: key, metadata } })
      }

      if (this.album) {
        payload.album = {
          title: this.album,
          author: this.$store.state.token.name
        }
      }

      const signedUrisResponse = await post('medias/ingest', payload, {}, 'v2')

      for (let signedUri of signedUrisResponse.data) {
        for (let config of mediasConfigToUpload) {
          if (config.key === signedUri.key) {
            const req = axios
              .put(signedUri.uri, config.media, config.putConfig)
              .then(() => this.upload.uploaded++)
              .catch(() => this.upload.failed++)

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

      Promise.all(promises).then(() => {
        this.medias = []
        this.upload.state = 'end'

        this.$notify({ group: 'success', text: this.$t('admin.mediaAdd.notify.uploadSuccess') })
        clearInterval(intervalTimeEstimation)
        clearInterval(intervalDisplayTime)
        this.noSleep.disable()
        this.upload = uploadData()
      })
    }
  }
}
</script>
