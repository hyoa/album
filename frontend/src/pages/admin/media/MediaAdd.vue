<template>
  <AdminLayout>
    <PageTitle title="Ajouter des médias" />
    <div v-if="upload.state === 'running'">
      <div class="text-white">
        <h3>Téléchargement en cours...</h3>
        <div v-if="upload.timeRemaining">Temps restant: {{ upload.timeRemaining }}</div>
        <div>Total: {{ upload.total }}</div>
        <div>Réussi: {{ upload.uploaded }} </div>
        <div>Echec: {{ upload.failed }} </div>
      </div>
      <ul>
        <li v-for="media of upload.medias" :key="media.key">
          <div class="text-white">{{ media.key }}</div>
          <div class="shadow w-full bg-white mt-2 rounded">
            <div class="bg-blue-500 text-xs leading-none py-1 text-center text-white rounded" :style="`width: ${media.progress}%`">
              {{ media.progress }}%
            </div>
          </div>
        </li>
      </ul>
    </div>
    <form @submit.prevent="uploadMedias" v-show="upload.state !== 'running'">
      <AutoComplete v-model="folder" id="folder" placeholder="Un super dossier !" label="Dossier" type="text" endpoint="medias/folders/autocomplete" :allow-no-call="true"/>
      <div class="mb-3" v-if="folder.trim() !== ''">
        <label class="text-white" for="files">Medias</label>
        <div
          class="shadow-inner bg-white hover:bg-white border border-white hover:border-gray-400 rounded-sm relative h-16"
        >
          <div class="shadow-inner absolute w-full h-full pin-t flex justify-center items-center">
            <span>Clique ou dépose tes fichiers ici </span>
          </div>
          <input
            id="files"
            type="file"
            class="opacity-0 w-full h-full"
            accept="image/png,image/jpg,video/mp4"
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
  components: { PageTitle, AutoComplete, AdminLayout },
  data () {
    return {
      folder: '',
      medias: [],
      upload: uploadData(),
      noSleep: new NoSleep()
    }
  },
  methods: {
    loadMedia ({ target }) {
      if (!target.validity.valid && target.files.length <= 0) {
        return null
      }

      if (target.files.length > 20) {
        this.$notify({ group: 'warning', text: 'Il n\'est pas possible d\'envoyer plus de 20 fichier à la fois' })
        return null
      }

      this.medias = [...this.medias, ...target.files]

      this.uploadMedias()
    },
    uploadMedias () {
      this.noSleep.enable()
      this.upload.state = null
      this.upload.progress = 0
      const promises = []
      this.upload.total = this.medias.length

      let totalSize = 0

      for (let media of this.medias) {
        totalSize += media.size
        const nameClean = media.name.replace(/[^a-zA-Z0-9.]/g, '').normalize('NFD').replace(/[\u0300-\u036f]/g, '')
        const folderClean = this.folder.replace(/\s+/g, '-').toLowerCase().normalize('NFD').replace(/[\u0300-\u036f]/g, '')
        const key = `${this.$store.state.token.name}_${folderClean}_${nameClean}`

        this.upload.medias.push({ key: nameClean, progress: 0, uploadedSize: 0 })

        const putConfig = {
          onUploadProgress: progressEvent => {
            for (let i in this.upload.medias) {
              if (this.upload.medias[i].key === nameClean) {
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

        const req = post('media/signed-uri', { file: key, type: media.type })
          .then(res => {
            return axios.put(res.data.uri, media, putConfig).then(() => this.upload.uploaded++)
          }).catch(() => {
            this.upload.failed++
          })

        promises.push(req)
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
          this.upload.timeRemaining = 'Calcul du temps restant en cours'
        } else {
          const sizeUploaded = this.upload.medias.reduce(reduceAdd, 0)

          const timeRemaining = (totalSize - sizeUploaded) / uploadSpeed

          if (timeRemaining < 1000) {
            this.upload.timeRemaining = 'C\'est bientôt terminé !'
          }

          this.upload.timeRemaining = `~ ${prettyMs(timeRemaining)}`
        }
      }, 1000)

      Promise.all(promises).then(() => {
        this.medias = []
        this.upload.state = 'end'

        this.$notify({ group: 'success', text: 'Les fichiers ont été transférés sur le serveur' })
        clearInterval(intervalTimeEstimation)
        clearInterval(intervalDisplayTime)
        this.noSleep.disable()
        this.upload = uploadData()
      })
    }
  }
}
</script>
