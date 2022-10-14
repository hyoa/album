<template>
  <div data-e2e="album-card" class="w-full h-full p-1 md:p-3" :key="album.slug">
    <router-link
      :to="{ name: 'album_view', params: { slug: album.slug } }"
    >
      <div class="overflow-hidden h-full rounded-sm hover:shadow-md">
        <div class="thumbnail relative h-full">

          <img
            :src="thumbnail"
            class="h-full rounded-sm"
          />
        </div>
      </div>
    </router-link>
  </div>
</template>

<style scoped lang="scss">
.thumbnail {
  img {
    width: 100%;
    object-fit: cover;
  }
}

@screen md {
  .thumbnail {
    img {
      width: 100%;
      object-fit: cover;
    }
  }
}

@screen lg {
  .thumbnail {
    img {
      width: 100%;
      object-fit: cover;
      max-height: 400px;
    }
  }
}
</style>

<script>
export default {
  name: 'AlbumCard',
  props: ['album', 'isFirst'],
  data () {
    return {
      thumbnail: this.getThumbnailAtLoad(),
      thumbnailIndex: 0
    }
  },
  mounted () {
    const favorites = this.album.medias.filter(media => media.favorite)

    if (favorites > 1) {
      setInterval(() => {
        if (favorites.length - 1 > this.thumbnailIndex) {
          this.thumbnailIndex++
        } else {
          this.thumbnailIndex = 0
        }

        this.thumbnail = favorites[this.thumbnailIndex]
      }, 15000)
    }
  },
  methods: {
    getThumbnailAtLoad () {
      if (this.album.medias.length > 0) {
        const favorites = this.album.medias.filter(media => media.favorite)

        if (favorites.length > 0) {
          return favorites[0].urls.small
        }

        const photos = this.album.medias.filter(media => media.kind === 'PHOTO')

        if (photos.length > 0) {
          return photos[0].urls.small
        }
      }
    }
  }
}
</script>
