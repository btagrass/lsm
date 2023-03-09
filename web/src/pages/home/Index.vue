<template>
  <el-container>
    <el-header>
      <div class="row-side">
        <div class="row-center">
          <div class="hover" @click="toggleCollapsed">
            <el-icon v-if="collapsed">
              <expand></expand>
            </el-icon>
            <el-icon v-else>
              <fold></fold>
            </el-icon>
          </div>
          <span>Lsm</span>
        </div>
        <div class="row-center">
          <div class="hover" @click="clearTabs">
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
      <el-aside :width="collapsed ? '75px' : '302px'">
        <el-menu :collapse="collapsed" :default-active="route.fullPath" unique-opened router>
          <template v-for="r in resources" :key="r.id">
            <el-sub-menu v-if="r.children.length" :index="r.uri">
              <template #title>
                <el-icon>
                  <component :is="r.icon"></component>
                </el-icon>
                <span>{{ r.name }}</span>
              </template>
              <el-menu-item v-for="c in r.children" :key="c.id" :index="c.uri">
                <el-icon>
                  <component :is="c.icon"></component>
                </el-icon>
                <span>{{ c.name }}</span>
              </el-menu-item>
            </el-sub-menu>
            <el-menu-item v-else :index="r.uri">
              <el-icon>
                <component :is="r.icon"></component>
              </el-icon>
              <span>{{ r.name }}</span>
            </el-menu-item>
          </template>
        </el-menu>
      </el-aside>
      <el-main>
        <el-tabs v-model="route.fullPath" type="card" @tab-click="clickTab" @tab-remove="removeTab">
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

<script setup>
import { useRoute, useRouter } from "vue-router"
import screenfull from "screenfull"
import { useStore } from "@/store"

const router = useRouter()
const route = useRoute()
const { collapsed, pages, resources, user, toggleCollapsed, clearPages, removePage, clearUser, clearResources } = useStore()

const toggleScreen = () => screenfull.toggle()
const commandUser = (command) => {
  if (command == "logout") {
    clearUser()
    clearPages()
    clearResources()
    router.push("/login")
  }
}
const clearTabs = () => {
  clearPages()
  router.push("/")
}
const clickTab = (tab) => router.push(tab.paneName)
const removeTab = (name) => {
  removePage(name)
  router.push("/")
}
</script>
