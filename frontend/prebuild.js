const fs = require('fs')

const apiUri = process.env.CONTEXT === 'production' ? process.env.API_URI : process.env.API_URI_DEV
const siteTitle = process.env.SITE_TITLE
const siteShortTitle = process.env.SITE_SHORT_TITLE

// Create correct manifest.json
fs.copyFileSync(`${__dirname}/public/manifest.json.dist`, `${__dirname}/public/manifest.json`)
const manifestFile = fs.readFileSync(`${__dirname}/public/manifest.json`, 'utf-8')
let manifestFileReplaced = manifestFile.replace(/SITE_TITLE/, siteTitle)
manifestFileReplaced = manifestFileReplaced.replace(/SITE_SHORT_TITLE/, siteShortTitle)
fs.writeFileSync(`${__dirname}/public/manifest.json`, manifestFileReplaced)

// Create .env file for production
fs.writeFileSync(`${__dirname}/.env`, `VUE_APP_API_URI=${apiUri}\nVUE_APP_SITE_TITLE=${siteTitle}\nVUE_APP_SITE_SHORT_TITLE=${siteShortTitle}`)

if (process.env.CONTEXT !== 'production') {
  // Create .env file for pr review
  fs.writeFileSync(`${__dirname}/.env.e2e`, `VUE_APP_API_URI=${process.env.API_URI_DEV}\nVUE_APP_SITE_TITLE=${siteTitle}\nVUE_APP_SITE_SHORT_TITLE=${siteShortTitle}`)

  const env = {
    'VUE_APP_API_URI': apiUri,
    'VUE_APP_SITE_TITLE': siteTitle,
    'VUE_APP_SITE_SHORT_TITLE': siteShortTitle
  }
  fs.writeFileSync(`${__dirname}/cypress.env.json`, JSON.stringify(env), 'utf-8')
}
