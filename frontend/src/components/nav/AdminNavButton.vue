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
      <div class="submenu mt-2 flex relative -left-20% z-50" :class="[visible ? 'block-animate' : 'hidden']">
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

<style lang="css" scoped>
  @keyframes show {
    from { display: 'none'; opacity: 0;}
    to { display: 'block'; opacity: 100%; }
  }

  @keyframes hide {
    from { display: 'block'; opacity: 100%; }
    to { display: 'none'; opacity: 0;}
  }

  .block-animate {
    animation-name: show;
    animation-duration: 0.2s;
  }

  .hidden-animate {
    animation-name: hide;
    animation-duration: 0.5s;
  }
</style>
