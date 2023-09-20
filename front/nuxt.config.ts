// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss', '@nuxtjs/apollo', 'nuxt-icon', '@vite-pwa/nuxt'],
  apollo: {
    clients: {
      default: {
        httpEndpoint: `${process.env.API_URI}/v3/graphql`,
        tokenStorage: "localStorage"
      },
    },
  },
  pwa: {
    registerType: 'autoUpdate',
    manifest: {
      name: 'Pauline & Jules',
      short_name: 'Pauline&Jules',
      theme_color: '#3490dc',
      background_color: '#ffffff',
      icons: [
        {
          src: 'logo-192x192.png',
          sizes: '192x192',
          type: 'image/png',
        },
        {
          src: 'logo-512x512.png',
          sizes: '512x512',
          type: 'image/png',
        },
      ],
    },
    workbox: {
      navigateFallback: '/',
      globPatterns: ['**/*.{js,css,html,png,svg,ico}'],
    },
    client: {
      installPrompt: true,
    },
    devOptions: {
      enabled: true,
      suppressWarnings: true,
      navigateFallbackAllowlist: [/^\/$/],
      type: 'module',
    },
  }
})
