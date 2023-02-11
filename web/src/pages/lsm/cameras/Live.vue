<template>
  <div>
    <div id="player"></div>
    <div>
      <el-button-group>
        <el-button type="primary" icon="TopLeft" @click="control('LeftUp')" />
        <el-button type="primary" icon="ArrowUpBold" @click="control('Up')" />
        <el-button type="primary" icon="TopRight" @click="control('RightUp')" />
        <el-button type="primary" icon="ZoomIn" @click="control('ZoomIn')" />
      </el-button-group>
    </div>
    <div>
      <el-button-group>
        <el-button type="primary" icon="ArrowLeftBold" @click="control('Left')" />
        <el-button type="primary" icon="House" />
        <el-button type="primary" icon="ArrowRightBold" @click="control('Right')" />
        <el-button type="primary" icon="ZoomOut" @click="control('ZoomOUt')" />
      </el-button-group>
    </div>
    <div>
      <el-button-group>
        <el-button type="primary" icon="BottomLeft" @click="control('LeftDown')" />
        <el-button type="primary" icon="ArrowDownBold" @click="control('Down')" />
        <el-button type="primary" icon="BottomRight" @click="control('RightDown')" />
        <el-button type="primary" icon="Camera" @click="snapshot()" />
      </el-button-group>
    </div>
  </div>
</template>

<script>
import { inject, reactive, toRefs, onMounted, onUnmounted } from "vue"
import { usePost } from "@/http"

export default {
  props: {
    id: {
      type: String,
      required: true,
    },
  },
  setup(props, context) {
    const api = inject("api")

    const state = reactive({
      id: props.id,
      data: {
        player: "",
        url: "",
      },
    })

    const start = async () => {
      state.data.url = await usePost(`${api}/${state.id}/streams/1/start`)
    }
    const control = async (command) => {
      await usePost(`${api}/${state.id}/ptz/${command}/2`)
    }
    const snapshot = async () => {
      window.open(await usePost(`${api}/${state.id}/streams/1/snapshot`), "_blank")
    }
    onMounted(async () => {
      await start()
      state.player = new WasmPlayer(null, "player")
      state.player.play(state.data.url, 1)
    })
    onUnmounted(async () => {
      state.player.pause()
      state.player.destroy()
      state.player = null
    })

    return {
      ...toRefs(state),
      control,
      snapshot,
    }
  },
}
</script>
