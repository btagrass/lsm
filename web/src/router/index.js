import { createRouter, createWebHashHistory } from "vue-router"
import store from "@/store"

const dynamicRoutes = (router) => {
  store.state.resources.forEach((r) => {
    r.children.forEach((c) => {
      router.addRoute("index", {
        name: c.code,
        path: c.url,
        meta: {
          title: c.name,
        },
        component: pages[`/src/pages${c.url}/Index.vue`],
      })
    })
  })
}
const pages = import.meta.glob("/src/pages/**/**.vue")
const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      name: "login",
      path: "/login",
      meta: {
        title: "登录",
      },
      component: () => import("@/pages/home/Login.vue"),
    },
    {
      name: "index",
      path: "/",
      redirect: "/",
      component: () => import("@/pages/home/Index.vue"),
      children: [
        {
          name: "dashboard",
          path: "/",
          meta: {
            title: "首页",
          },
          component: () => import("@/pages/home/Dashboard.vue"),
        },
      ],
    },
  ],
})
router.beforeEach((to, from, next) => {
  if (store.state.user.token == null) {
    if (to.fullPath.startsWith("/login")) {
      next()
    } else {
      next("/login")
    }
  } else {
    if (from.fullPath.startsWith("/login")) {
      dynamicRoutes(router)
    }
    if (to.fullPath.startsWith("/http")) {
      window.open(to.fullPath.substr(1), "_blank")
    } else {
      next()
    }
  }
})
router.afterEach((to, from) => {
  if (!to.fullPath.startsWith("/login")) {
    store.commit("savePage", {
      path: to.fullPath,
      title: to.meta.title,
    })
  }
})
dynamicRoutes(router)

export default router
