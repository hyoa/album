<template>
  <div class="absolute top-0 left-0 w-full mt-10">
    <section class="w-4/5 mx-auto bg-white p-4">
      <article v-if="navigator === 'chrome' || navigator === 'firefox'">
        <div>
          <span class="text-xl">Etape 1:</span>
          <p>
            Cliquer sur l'icone comme sur l'image ci-dessous:
            <img v-if="navigator === 'chrome'" :src="images.chromesPwaHelp.page" alt="">
            <img v-else :src="images.chromesPwaHelp.page" alt="">
          </p>
        </div>
        <div class="border-t border-darker-primary">
          <span class="text-xl">Etape 2:</span>
          <p>
            Cliquer sur "Ajouter à l'écran d'acceuil":
            <img v-if="navigator === 'chrome'" :src="images.chromesPwaHelp.add" alt="">
            <img v-else :src="images.chromesPwaHelp.add" alt="">
          </p>
        </div>
        <div class="border-t border-darker-primary">
          <span class="text-xl">Etape 3:</span>
          <p>
            Valider les différentes fenêtres qui s'affichent, l'icône du site sera maintenant présent sur votre écran d'acceuil.
          </p>
        </div>
      </article>
      <article v-else>
        <p>L'application ne peut-être installer qu'à l'aide de Google Chrome ou Firefox.</p>
        <p class="mt-1">Cependant, vous pouvez quand même utiliser l'application sans l'installer ! Vous ne recevrez cependant pas les notifications de façon optimal.</p>
      </article>
      <button class="mt-2 text-red-500" @click="$emit('onClose')">Fermer</button>
    </section>
  </div>
</template>

<script>
export default {
  name: 'PwaInstallation',
  data () {
    return {
      navigator: null,
      images: {
        chromesPwaHelp: {
          page: require('../../assets/help/chrome_page.jpg'),
          add: require('../../assets/help/chrome_add.jpg')
        },
        firefoxPwaHelp: {
          page: require('../../assets/help/firefox_page.jpg'),
          add: require('../../assets/help/firefox_add.jpg')
        }
      }
    }
  },
  created () {
    if (!!window.chrome && (!!window.chrome.webstore || !!window.chrome.runtime)) {
      this.navigator = 'chrome'
    } else if (typeof InstallTrigger !== 'undefined') {
      this.navigator = 'firefox'
    }
  }
}
</script>
