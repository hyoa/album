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
    if (this.album.favorites.length > 1) {
      setInterval(() => {
        if (this.album.favorites.length - 1 > this.thumbnailIndex) {
          this.thumbnailIndex++
        } else {
          this.thumbnailIndex = 0
        }

        this.thumbnail = this.album.favorites[this.thumbnailIndex]
      }, 15000)
    }
  },
  methods: {
    getThumbnailAtLoad () {
      if (this.album.favorites.length > 0) {
        return this.album.favorites[0]
      }
    }
  },
  computed: {
    description () {
      if (this.album.description.length > 80) {
        return this.album.description.substring(0, 80) + '...'
      }

      return this.album.description
    }
  }
}
</script>
