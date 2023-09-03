<template>
  <div>
    <vue-masonry-wall :items="items" :options="{width: 150, padding: 10}" @append="append">
      <template v-slot:default="{item}">
        <div class="relative p-1">
          <div
            v-if="editable"
            class="actionBtn absolute left-0 top-0 w-5 h-5 z-10 text-white flex justify-center items-center"
            @click="() => selectMedia(item.key)"
          >
            <i v-if="isSelected(item.key)" class="material-icons">check_box</i>
            <i v-else class="material-icons">check_box_outline_blank</i>
          </div>
          <div
            v-if="canStar"
            class="actionBtn absolute right-0 top-0 w-5 h-5 z-10 text-white flex justify-center items-center"
            @click="$emit('toggleFavorite', item)"
          >
            <i v-if="item.favorite" class="material-icons">star</i>
            <i v-else class="material-icons">star_border</i>
          </div>
          <img
            v-if="item.kind === 'PHOTO'"
            :src="item.urls.small"
            :alt="item.key"
            :class="{ 'opacity-25': isSelected(item.key) }"
            @click="onImageClick(item.key)"
          />
          <video
            v-else-if="item.kind === 'VIDEO'"
            controls
            :class="{ 'opacity-25': isSelected(item.key) }"
            preload="metadata"
          >
            <source :src="item.urls.small" type="video/mp4">
            Your browser does not support the video tag.
          </video>
        </div>
      </template>
    </vue-masonry-wall>
    <VueGallery
      v-if="!editable"
      :images="galleryUris"
      :index="indexGallery"
      @close="indexGallery = null"
      :options="{ disableScroll: false, hidePageScrollbars: false }"
    ></VueGallery>
  </div>
</template>

<script>
import VueMasonryWall from 'vue-masonry-wall'
import VueGallery from 'vue-gallery'

export default {
  name: 'Grid2',
  components: { VueMasonryWall, VueGallery },
  props: ['medias', 'editable', 'canDeleteMedia', 'canStar'],
  data () {
    return {
      indexGallery: null,
      indexAppend: 0,
      items: []
    }
  },
  methods: {
    append () {
      this.items = this.items.concat(this.medias.slice(this.indexAppend, this.indexAppend + 5))
      this.indexAppend += 5
    },
    selectMedia (key) {
      if (this.editable) {
        this.$store.commit('toggleMediaSelection', key)
      } else {
        let index = this.medias.findIndex(media => media.key === key)
        this.indexGallery = index
      }
    },
    isSelected (key) {
      return this.$store.state.mediaSelected.includes(key)
    },
    onImageClick (key) {
      if (!this.editable) {
        let index = this.medias.findIndex(media => media.key === key)

        this.indexGallery = index
      }
    }
  },
  mounted () {
    console.log(this.medias.length)
    this.items = this.medias.slice(0, 5)
    this.indexAppend = 5
  },
  computed: {
    galleryUris () {
      if (this.medias !== undefined) {
        return this.medias.map(media => {
          return {
            href: media.kind === 'VIDEO' ? media.urls.small : media.urls.large,
            urlset: `${media.urls.medium} 600w, ${media.urls.large} 1000w`,
            description: this.$t('grid.media.description', { author: media.author }),
            type: media.kind === 'VIDEO' ? 'video/mp4' : 'image/jpg'
          }
        })
      }

      return []
    }
  }
}
</script>

<style lang="scss" scoped>
  .actionBtn {
    background-color: rgba(0,0,0,0.7);
  }
</style>
