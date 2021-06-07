<template>
  <div data-e2e="medias-grid" class="relative">
    <div
      v-masonry="masonryId"
      transition-duration="0.3s"
      item-selector=".mediaTile"
      :class="{ 'opacity-100': isVisible, 'opacity-0': !isVisible }"
    >
      <lazy-component
        v-masonry-tile="masonryId"
        class="mediaTile w-1/2"
        v-for="(media, index) in medias"
        :key="index"
        @show="onShow"
      >
        <div class="relative p-1">
          <div
            v-if="editable"
            class="star absolute left-0 top-0 w-5 h-5 z-10 text-white flex justify-center items-center"
            @click="() => selectMedia(media.key, index)"
          >
            <i v-if="isSelected(media.key)" class="material-icons">check_box</i>
            <i v-else class="material-icons">check_box_outline_blank</i>
          </div>
          <div
            v-if="canStar"
            class="star absolute right-0 top-0 w-5 h-5 z-10 text-white flex justify-center items-center"
            @click="$emit('toggleFavorite', media)"
          >
            <i v-if="media.favorite" class="material-icons">star</i>
            <i v-else class="material-icons">star_border</i>
          </div>
          <img
            v-if="media.type === 'image'"
            :src="media.uris.small"
            :alt="media.key"
            :class="{ 'opacity-25': isSelected(media.key) }"
            @click="onImageClick(index)"
          />
          <video
            v-else-if="media.type === 'video'"
            controls
            :class="{ 'opacity-25': isSelected(media.key) }"
            preload="metadata"
          >
            <source :src="media.uris.original" type="video/mp4">
            Your browser does not support the video tag.
          </video>
        </div>
      </lazy-component>
    </div>
    <div
      class="text-center absolute pin-y w-full mt-10 z-0"
      :class="{ 'opacity-100': !isVisible, 'opacity-0': isVisible }"
      v-if="!isMobile"
    >
      {{ $t('grid.loading') }}
    </div>
    <VueGallery
      v-if="!editable"
      :images="galleryUris"
      :index="indexGallery"
      @close="indexGallery = null"
      :options="{ disableScroll: false, hidePageScrollbars: false }"
    ></VueGallery>
  </div>
</template>

<style scoped lang="scss">
.grid-item {
  max-width: 100%;

  @screen md {
    max-width: 33%;
  }

  video, img  {
    width: 100%;
  }
}
</style>

<script>
import VueGallery from 'vue-gallery'
import MobileDetect from 'mobile-detect'

export default {
  name: 'Grid',
  data () {
    return {
      isVisible: this.isMobile(),
      indexGallery: null,
      masonry: null,
      masonryId: 'mediasTiles'
    }
  },
  components: { VueGallery },
  props: ['medias', 'editable', 'canDeleteMedia', 'canStar'],
  mounted () {
    setTimeout(() => {
      this.updateGrid()
    }, 4000)
  },
  methods: {
    updateGrid () {
      console.log('updategrid')
      this.$redrawVueMasonry(this.masonryId)
    },
    selectMedia (key, index) {
      if (this.editable) {
        this.$store.commit('toggleMediaSelection', key)
      } else {
        this.indexGallery = index
      }
    },
    isSelected (key) {
      return this.$store.state.mediaSelected.includes(key)
    },
    onShow () {
      this.updateGrid()
    },
    isMobile () {
      const md = new MobileDetect(window.navigator.userAgent)
      return !!md.mobile()
    },
    onImageClick (index) {
      if (!this.editable) {
        this.indexGallery = index
      }
    }
  },
  computed: {
    galleryUris () {
      if (this.medias !== undefined) {
        return this.medias.map(media => {
          return {
            href: media.uris.original,
            urlset: `${media.uris.medium} 600w, ${media.uris.original} 1000w`,
            description: this.$t('grid.media.description', { author: media.author }),
            type: media.type === 'video' ? 'video/mp4' : 'image/jpg'
          }
        })
      }

      return []
    }
  }
}
</script>

<style lang="scss" scoped>
  .star {
    background-color: rgba(0,0,0,0.7);
  }
</style>
