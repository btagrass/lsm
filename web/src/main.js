import { createApp } from "vue"
import { createPinia } from "pinia"
import piniaPluginPersistedstate from "pinia-plugin-persistedstate"
import ElementPlus from "element-plus"
import * as ElementIcons from "@element-plus/icons-vue"
import zhCn from "element-plus/es/locale/lang/zh-cn"
import "element-plus/dist/index.css"
import App from "@/App.vue"
import router from "@/router"
import "@/assets/styles/index.less"

const app = createApp(App).
  use(createPinia().use(piniaPluginPersistedstate)).
  use(router).
  use(ElementPlus, {
    size: "small",
    zIndex: 3000,
    locale: zhCn,
  })
for (const [key, icon] of Object.entries(ElementIcons)) {
  app.component(key, icon)
}
app.mount("#app")
