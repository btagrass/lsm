<template>
  <el-form>
    <el-form-item label="实况">
      <div id="player"></div>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { inject, defineProps, reactive, onMounted, onUnmounted } from "vue"
import { usePost } from "@/http"

const api = inject("api")
const props = defineProps({
  values: Object
})
const state = reactive({
  name: props.values.name,
  data: {
    url: "",
  },
  player: null,
})

const start = async () => {
  state.data.url = await usePost(`${api}/${state.code}/streams/1/start`)
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
</script>
