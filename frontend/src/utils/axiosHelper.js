import axios from 'axios'

const BASE_URI = process.env.VUE_APP_API_URI || 'http://localhost:8011'
const getConfig = (headers = {}) => {
  const defaultHeaders = {
    headers: { 'Authorization': 'Bearer ' + localStorage.getItem('album-token') }
  }

  return { ...defaultHeaders, ...headers }
}

export const post = (route, data, headers) => {
  return axios.post(`${BASE_URI}/${route}`, data, getConfig(headers))
}

export const put = (route, data) => {
  return axios.put(`${BASE_URI}/${route}`, data, getConfig())
}

export const get = (route) => {
  return axios.get(`${BASE_URI}/${route}`, getConfig())
}

export const deleteMethod = route => {
  return axios.delete(`${BASE_URI}/${route}`, getConfig())
}
