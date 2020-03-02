<template>
  <div data-e2e="album-card" class="w-full h-full p-1 md:p-3" :key="album.slug">
    <div class="overflow-hidden h-full rounded-sm" :class="{ 'shadow-md lg:shadow-none' : isFirst }">
      <div class="thumbnail relative h-full">
        <router-link
          :to="{ name: 'album_view', params: { slug: album.slug } }"
        >
          <img
            :src="thumbnail"
            class="h-full rounded-sm"
          />
        </router-link>
        <div
          class="absolute bottom-0 left-0 w-full text-white card-title p-1 md:p-2 rounded-b-sm"
          :class="[ isFirst ? 'text-2xl pl-4' : 'text-center' ]"
        >
          {{ album.title }}
        </div>
      </div>
    </div>
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

.card-title {
  background-color: rgba(77, 134, 168, 0.9);
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
