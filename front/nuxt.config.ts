// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss', '@nuxtjs/apollo', 'nuxt-icon'],
  apollo: {
    clients: {
      default: {
        httpEndpoint: 'http://127.0.0.1:3118/v3/graphql',
        tokenStorage: "localStorage"
      },
    },
  },
})
