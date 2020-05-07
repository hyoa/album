<template>
  <LayoutGrid>
    <section data-e2e="notification-feature-section" v-if="hasAcceptedNotification === 'default'" class="border border-blue-500 p-2 rounded-sm">
      <h3 class="border-b border-blue-800 text-xl">{{ $t('home.notificationSection.title') }}</h3>
      <div class="my-2">
        <p>
          {{ $t('home.notificationSection.message') }}
        </p>
        <p class="mt-1" v-if="!isRunningAsPwa">
          {{ $t('home.notificationSection.advice') }}
          <a class="text-blue-700 " @click.stop="showPwaHelp = true">
            {{ $t('home.notificationSection.help') }}
            <i class="material-icons align-top">
              help_outline
            </i>
          </a>
        </p>
      </div>
      <div class="text-center mt-2 py-1">
        <SimpleButton
          :status="notificationValidationStatus"
          @click="acceptNotification"
          size="small"
        >
          {{ $t('home.notificationSection.form.accept') }}
        </SimpleButton>
        <button @click="declineNotification" class="text-red-500 mt-2">
          {{ $t('home.notificationSection.form.refuse') }}
        </button>
      </div>
    </section>
    <section data-e2e="search-form-section" class="mb-6">
      <h2 class="mb-3 text-xl">{{ $t('home.searchFormSection.title') }}</h2>
      <form @submit.prevent="onSearch">
        <input
          title="search"
          class="appearance-none block w-full bg-grey-lighter text-grey-darker border border-grey-lighter rounded py-3 px-4 leading-tight focus:outline-none focus:bg-white focus:border-grey mb-3"
          placeholder="ex: vacances"
          v-model="searchTerm"
        />
        <button class="bg-primary text-white font-bold py-1 px-2 rounded focus:outline-none focus:shadow-outline text-sm md:text-md" type="submit">
          {{ $t('home.searchFormSection.form.submit') }}
        </button>
      </form>
    </section>
    <section data-e2e="search-result-section" v-if="searchedTerm">
        <h2 class="text-xl">{{ $t('home.searchResultSection.title') }}</h2>
        <div class="searchResult">
          <AlbumCard v-for="album in albums" :key="album.slug" :album="album" />
        </div>
    </section>
    <section data-e2e="last-albums-section" v-else>
      <h2 class="text-xl">{{ $t('home.lastAlbumSection.title') }}</h2>
      <div class="showroom">
        <div class="grid1" v-if="albums[0]">
          <AlbumCard :isFirst="true" :key="albums[0].slug" :album="albums[0]" />
        </div>
        <div class="grid2" v-if="albums[1]">
          <AlbumCard :isFirst="false" :key="albums[1].slug" :album="albums[1]" />
        </div>
        <div class="grid3" v-if="albums[2]">
          <AlbumCard :isFirst="false" :key="albums[2].slug" :album="albums[2]" />
        </div>
      </div>
      <div>
        <div>
          <div class="md:flex" v-for="(albumsToLoad, index) in albumsMore" :key="index">
            <AlbumCard class="lg:w-1/3 md:w-1/2" v-for="album in albumsToLoad" :key="album.slug" :album="album" />
          </div>
        </div>
        <div class="flex justify-center" v-if="canLoadMore">
          <button
            @click="onLoadMore"
            v-if="!loadingMore"
            class="text-blue-500 text-underline text-sm mt-3 focus:outline-none active:outline-none border border-blue-500 p-1 rounded"
          >
            {{ $t('home.lastAlbumSection.loadMore') }}
          </button>
          <CubicLoader v-else />
        </div>
      </div>
    </section>
    <HelpPwa v-if="showPwaHelp" @onClose="showPwaHelp = false"/>
  </LayoutGrid>
</template>

<script>
import firebase from 'firebase/app'
import 'firebase/messaging'
import { get, post } from '../utils/axiosHelper'

import AlbumCard from '../components/album/AlbumCard'
import LayoutGrid from '../components/layout/LayoutGrid'
import CubicLoader from '../components/loader/CubicLoader'
import SimpleButton from '../components/form/button/SimpleAnimateButton'
import HelpPwa from '../components/help/PwaInstallation'

export default {
  name: 'Home',
  components: { LayoutGrid, AlbumCard, CubicLoader, SimpleButton, HelpPwa },
  data () {
    return {
      albums: [],
      searchTerm: '',
      searchedTerm: null,
      albumsMore: [],
      loadingMore: false,
      canLoadMore: true,
      currentPage: 0,
      hasAcceptedNotification: 'unknown',
      notificationValidationStatus: 'ready',
      showPwaHelp: false
    }
  },
  async created () {
    try {
      const res = await get('albums')
      this.albums = res.data
    } catch ({ response: { status } }) {
      if (status === 401) {
        this.$store.commit('setFlashMessage', 'Par mesure de sécurité, vous avez été déconnecté. Vous pouvez vous reconnecter avec le formulaire ci-dessous.')
        localStorage.removeItem('album-token')
        this.$router.push({ name: 'auth' })
      }
    }

    if (localStorage.getItem('declineNotification') !== null) {
      this.hasAcceptedNotification = 'denied'
    } else if (Notification.permission === 'granted') {
      this.hasAcceptedNotification = 'granted'
    } else {
      this.hasAcceptedNotification = 'default'
    }
  },
  methods: {
    async onSearch () {
      try {
        const res = await get(`albums?search=${this.searchTerm}&limit=100`)

        this.albums = res.data
        this.searchedTerm = this.searchTerm
      } catch ({ response: { satus } }) {
        if (status === 401) {
          this.$store.commit('setFlashMessage', 'auth.alert.disconnected')
          localStorage.removeItem('album-token')
          this.$router.push({ name: 'auth' })
        }
      }
    },
    async onLoadMore () {
      try {
        this.loadingMore = true
        const limit = 6
        const offset = this.currentPage * limit + 3

        const res = await get(`albums?offset=${offset}&limit=${limit}`)

        if (res.data.length < limit) {
          this.canLoadMore = false
        }

        this.albumsMore.push(res.data)
        this.loadingMore = false
        this.currentPage++
      } catch ({ response: { satus } }) {
        if (status === 401) {
          this.$store.commit('setFlashMessage', 'auth.alert.disconnect')
          localStorage.removeItem('album-token')
          this.$router.push({ name: 'auth' })
        }
      }
    },
    async acceptNotification () {
      try {
        this.notificationValidationStatus = 'pending'
        const messaging = firebase.messaging()

        try {
          await Notification.requestPermission()
          await messaging.requestPermission()
        } catch (e) {
          this.hasAcceptedNotification = 'denied'
          return
        }

        const token = await messaging.getToken()
        localStorage.setItem('albumNotificationToken', token)

        await post('notification/subscribe', { token, channel: 'album' })

        if (this.$store.state.token.role === 9) {
          await post('notification/subscribe', { token, channel: 'admin' })
        }

        this.hasAcceptedNotification = 'granted'
      } catch (e) {
        this.hasAcceptedNotification = 'unknown'
        this.notificationValidationStatus = 'ready'
      }
    },
    async declineNotification () {
      localStorage.setItem('declineNotification', 1)
      this.hasAcceptedNotification = 'denied'
    }
  },
  computed: {
    isRunningAsPwa () {
      return !!window.matchMedia('(display-mode: standalone)').matches || !!window.navigator.standalone
    }
  }
}
</script>

<style scoped lang="scss">
  .showroom {
    min-height: 50vh;
    display: grid;
    grid-template-columns: 1fr 0.5fr;
    grid-template-rows: 1fr 1fr;
    grid-column-gap: 0px;
    grid-row-gap: 0px;
    .grid1 { grid-area: 1 / 1 / 3 / 2; }
    .grid2 { grid-area: 1 / 2 / 2 / 3; }
    .grid3 { grid-area: 2 / 2 / 3 / 3; }
  }

  .searchResult {
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: 1fr 1fr 1fr 1fr 1fr;
    grid-column-gap: 0px;
    grid-row-gap: 0px;
  }

  @screen md {
    .showroom {
      display: grid;
      grid-template-columns: 1fr 1fr 1fr;
      grid-template-rows: 1fr;
      grid-column-gap: 0px;
      grid-row-gap: 0px;
      .grid1 { grid-area: 1 / 1 / 2 / 2; }
      .grid2 { grid-area: 1 / 2 / 2 / 3; }
      .grid3 { grid-area: 1 / 3 / 2 / 4; }
    }

    .searchResult {
      display: grid;
      grid-template-columns: repeat(6, 1fr);
      grid-template-rows: 1fr 1fr 1fr 1fr 1fr;
      grid-column-gap: 0px;
      grid-row-gap: 0px;
    }
  }
</style>
