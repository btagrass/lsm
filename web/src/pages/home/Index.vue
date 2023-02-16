<template>
  <el-container>
    <el-header>
      <div class="row-side">
        <div class="row-center">
          <div class="hover" @click="toggleCollapse">
            <el-icon v-if="collapse">
              <expand></expand>
            </el-icon>
            <el-icon v-else>
              <fold></fold>
            </el-icon>
          </div>
          <span>Lsm</span>
        </div>
        <div class="row-center">
          <div class="hover" @click="clearPages">
            <el-tooltip content="清空">
              <el-icon>
                <FolderDelete />
              </el-icon>
            </el-tooltip>
          </div>
          <div class="hover" @click="toggleScreen">
            <el-tooltip content="全屏">
              <el-icon>
                <FullScreen />
              </el-icon>
            </el-tooltip>
          </div>
          <el-dropdown class="hover" @command="commandUser">
            <div>
              <el-icon>
                <User />
              </el-icon>
              <span>{{ user.userName }}</span>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">注销</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </el-header>
    <el-container>
      <el-aside :width="collapse ? '75px' : '302px'">
        <el-menu :collapse="collapse" :default-active="fullPath" unique-opened router>
          <template v-for="r in resources" :key="r.id">
            <el-sub-menu v-if="r.children.length" :index="r.url">
              <template #title>
                <el-icon>
                  <component :is="r.icon"></component>
                </el-icon>
                <span>{{ r.name }}</span>
              </template>
              <el-menu-item v-for="c in r.children" :key="c.id" :index="c.url">
                <el-icon>
                  <component :is="c.icon"></component>
                </el-icon>
                <span>{{ c.name }}</span>
              </el-menu-item>
            </el-sub-menu>
            <el-menu-item v-else :index="r.url">
              <el-icon>
                <component :is="r.icon"></component>
              </el-icon>
              <span>{{ r.name }}</span>
            </el-menu-item>
          </template>
        </el-menu>
      </el-aside>
      <el-main>
        <el-tabs v-model="fullPath" type="card" @tab-click="clickTab" @tab-remove="removeTab">
          <el-tab-pane v-for="page in pages" :key="page.path" :label="page.title" :name="page.path"
            :closable="page.path != '/'"></el-tab-pane>
        </el-tabs>
        <div id="content">
          <router-view v-slot="{ Component }">
            <!-- <keep-alive> -->
            <component :is="Component"></component>
            <!-- </keep-alive> -->
          </router-view>
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import { toRefs } from "vue"
import { useRoute, useRouter } from "vue-router"
import { useStore } from "vuex"
import screenfull from "screenfull"

export default {
  setup() {
    const router = useRouter()
    const route = useRoute()
    const store = useStore()

    const toggleCollapse = () => store.commit("toggleCollapse")
    const clearPages = () => {
      store.commit("clearPages")
      router.push("/")
    }
    const toggleScreen = () => screenfull.toggle()
    const commandUser = (command) => {
      if (command == "logout") {
        store.commit("clearUser")
        store.commit("clearPages")
        store.commit("clearResources")
        router.push("/login")
      }
    }
    const clickTab = (tab) => router.push(tab.paneName)
    const removeTab = (name) => {
      store.commit("removePage", name)
      router.push("/")
    }

    return {
      ...toRefs(route),
      ...toRefs(store.state),
      toggleCollapse,
      clearPages,
      toggleScreen,
      commandUser,
      clickTab,
      removeTab,
    }
  },
}
</script>
