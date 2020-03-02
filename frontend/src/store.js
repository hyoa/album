import Vue from 'vue'
import Vuex from 'vuex'
import jwtDecode from 'jwt-decode'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    mediaSelected: [],
    token: null,
    flashMessage: null,
    sideBarIsOpen: false,
    alert: {
      message: null,
      type: null,
      isVisible: false
    }
  },
  mutations: {
    toggleMediaSelection (state, key) {
      if (state.mediaSelected.includes(key)) {
        state.mediaSelected = state.mediaSelected.filter(
          mediaToFilter => key !== mediaToFilter
        )
      } else {
        state.mediaSelected.push(key)
      }
    },
    resetMediaSelection (state) {
      state.mediaSelected = []
    },
    setToken (state, token) {
      state.token = jwtDecode(token)
    },
    setFlashMessage (state, message) {
      state.flashMessage = message
    },
    toggleSideBar (state, toggle) {
      state.sideBarIsOpen = toggle
    },
    displayAlert (state, { message, type }) {
      state.alert = {
        message,
        type,
        isVisible: true
      }
    },
    hideAlert (state) {
      state.alert = {
        message: null,
        type: null,
        isVisible: false
      }
    }
  },
  actions: {

  }
})
