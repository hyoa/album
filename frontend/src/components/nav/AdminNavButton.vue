<template>
  <div v-if="link">
    <router-link :to="{ name: link }">
      <div :class="[iconClass, small ? 'w-8 h-8' : 'w-12 h-12']" class="rounded-full p-2 flex justify-center items-center">
        <slot name="icon"></slot>
      </div>
    </router-link>
  </div>
  <div v-else>
      <div @click=onDisplaySubMenu :class="[iconClass, small ? 'w-8 h-8' : 'w-12 h-12']" class="rounded-full p-2 flex justify-center items-center">
        <slot name="icon"></slot>
      </div>
      <div class="mt-2 flex" :class="[visible ? 'block' : 'hidden']">
        <slot name="submenu"></slot>
      </div>
  </div>
</template>

<script>
export default {
  props: ['link', 'iconClass', 'small', 'menuName'],
  data () {
    return {
      visible: false
    }
  },
  mounted () {
    this.$root.$on('onSubMenuDisplay', (menuName) => {
      if (this.menuName !== menuName && this.visible) {
        this.visible = false
      }
    })
  },
  methods: {
    onDisplaySubMenu () {
      if (!this.visible) {
        this.$root.$emit('onSubMenuDisplay', this.menuName)
      }

      this.visible = !this.visible
    }
  }
}
</script>
