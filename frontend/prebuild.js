const fs = require('fs')

const firebaseSender = process.env.CONTEXT === 'production' ? process.env.FIREBASE_SENDER_ID : process.env.FIREBASE_SENDER_ID_DEV
const apiUri = process.env.CONTEXT === 'production' ? process.env.API_URI : process.env.API_URI_DEV

// Create correct firebase-messaging-sw.js
fs.copyFileSync(`${__dirname}/public/firebase-messaging-sw.js.dist`, `${__dirname}/public/firebase-messaging-sw.js`)
const firebaseWorkerFile = fs.readFileSync(`${__dirname}/public/firebase-messaging-sw.js`, 'utf-8')
const firebaseWorkerFileReplaced = firebaseWorkerFile.replace(/SENDER_ID/, firebaseSender)
fs.writeFileSync(`${__dirname}/public/firebase-messaging-sw.js`, firebaseWorkerFileReplaced)

// Create .env file for production
fs.writeFileSync(`${__dirname}/.env`, `VUE_APP_API_URI=${apiUri}\nVUE_APP_FIREBASE_SENDER_ID=${firebaseSender}`)

if (process.env.CONTEXT !== 'production') {
  // Create .env file for pr review
  fs.writeFileSync(`${__dirname}/.env.e2e`, `VUE_APP_API_URI=${process.env.API_URI_DEV}\nVUE_APP_FIREBASE_SENDER_ID=${process.env.FIREBASE_SENDER_ID_DEV}`)

  const env = {
    'VUE_APP_API_URI': apiUri,
    'VUE_APP_FIREBASE_SENDER_ID': firebaseSender
  }
  fs.writeFileSync(`${__dirname}/cypress.env.json`, JSON.stringify(env), 'utf-8')
}
