import axios from 'axios'

const GRAPHQL_URI = process.env.VUE_APP_GRAPHQL_URI || 'http://127.0.0.1:3118'

const getConfig = (headers = {}) => {
    const defaultHeaders = {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('albumToken') }
    }
  
    return { ...defaultHeaders, ...headers }
  }

export default (payload, version = 'v1', variables = {}) => {
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
  