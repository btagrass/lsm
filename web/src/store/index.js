import { reactive, toRefs } from "vue"
import { defineStore } from "pinia"

export const useStore = defineStore("store", () => {
    const state = reactive({
        collapsed: false,
        pages: [{
            path: "/",
            title: "首页",
        }],
        resources: [],
        user: {
            token: "",
        },
    })

    const toggleCollapsed = () => {
        state.collapsed = !state.collapsed
    }
    const clearPages = () => {
        state.pages.splice(1, state.pages.length)
    }
    const removePage = (path) => {
        const index = state.pages.findIndex((i) => {
            return i.path == path
        })
        state.pages.splice(index, 1)
    }
    const savePage = (page) => {
        const exist = state.pages.some((i) => {
            return i.path == page.path
        })
        if (!exist) {
            state.pages.push(page)
        }
    }
    const clearResources = () => {
        state.resources = []
    }
    const saveResources = (resources) => {
        state.resources = resources
    }
    const clearUser = () => {
        state.user = {
            token: "",
        }
    }
    const saveUser = (user) => {
        state.user = user
    }

    return {
        ...toRefs(state),
        toggleCollapse: toggleCollapsed,
        clearPages,
        removePage,
        savePage,
        clearResources,
        saveResources,
        clearUser,
        saveUser,
    }
}, {
    persist: {
        enabled: true,
        strategies: [{
            storage: localStorage,
        }],
    }
})
