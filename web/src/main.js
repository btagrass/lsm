import { createApp } from "vue"
import ElementPlus from "element-plus"
import * as ElementIcons from "@element-plus/icons-vue"
import zhCn from "element-plus/es/locale/lang/zh-cn"
import "element-plus/dist/index.css"
import App from "@/App.vue"
import router from "@/router"
import store from "@/store"
import "@/assets/styles/index.less"

const app = createApp(App).
  use(router).
  use(store).
  use(ElementPlus, {
    size: "small",
    zIndex: 3000,
    locale: zhCn,
  })
for (const [key, icon] of Object.entries(ElementIcons)) {
  app.component(key, icon)
}
app.mount("#app")
