import { createStore } from "vuex"
import createPersistedState from "vuex-persistedstate"

const store = createStore({
  state: {
    collapse: false,
    pages: [
      {
        path: "/",
        title: "首页",
      },
    ],
    resources: [],
    user: {},
  },
  mutations: {
    // Collapse
    toggleCollapse: (state) => {
      state.collapse = !state.collapse
    },
    // Page
    clearPages: (state) => {
      state.pages = [
        {
          path: "/",
          title: "首页",
        },
      ]
    },
    removePage: (state, path) => {
      const index = state.pages.findIndex((i) => {
        return i.path == path
      })
      state.pages.splice(index, 1)
    },
    savePage: (state, page) => {
      const exist = state.pages.some((i) => {
        return i.path == page.path
      })
      if (!exist) {
        state.pages.push(page)
      }
    },
    // Resource
    clearResources: (state) => {
      state.resources = []
    },
    saveResources: (state, resources) => {
      state.resources = resources
    },
    // User
    clearUser: (state) => {
      state.user = {}
    },
    saveUser: (state, user) => {
      state.user = user
    },
  },
  plugins: [
    createPersistedState({
      storage: window.sessionStorage,
    }),
  ],
})

export default store
