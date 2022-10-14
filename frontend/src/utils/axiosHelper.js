import axios from 'axios'

const BASE_URI = process.env.VUE_APP_API_URI || 'http://localhost:8011'
const GRAPHQL_URI = process.env.VUE_APP_GRAPHQL_URI || 'http://localhost:3118'
const getConfig = (headers = {}) => {
  const defaultHeaders = {
    headers: { 'Authorization': 'Bearer ' + localStorage.getItem('album-token') }
  }

  return { ...defaultHeaders, ...headers }
}

export const post = (route, data, headers, version = 'v1') => {
  return axios.post(`${BASE_URI}/${version}/${route}`, data, getConfig(headers))
}

export const put = (route, data, version = 'v1') => {
  return axios.put(`${BASE_URI}/${version}/${route}`, data, getConfig())
}

export const get = (route, version = 'v1') => {
  return axios.get(`${BASE_URI}/${version}/${route}`, getConfig())
}

export const deleteMethod = (route, version = 'v1') => {
  return axios.delete(`${BASE_URI}/${version}/${route}`, getConfig())
}

export const graphql = (payload, version = 'v1', variables = {}) => {
  return new Promise(async (resolve, reject) => {
    try {
      let res = await axios.post(
        `${GRAPHQL_URI}/${version}/graphql`,
        {
          query: payload,
          variables
        },
        getConfig()
      )

      if (res.data.errors !== undefined) {
        let messages = []

        for (let error of res.data.errors) {
          messages.push(error.message)
        }

        reject(messages.join(', '))
      }
      resolve(res.data.data)
    } catch (e) {
      reject(e)
    }
  })
}
