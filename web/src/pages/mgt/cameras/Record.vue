<template>
  <el-form v-model="data">
    <el-form-item label="日期" prop="date">
      <el-date-picker v-model="data.date" value-format="YYYY-MM-DDTHH:mm:ss.SSSZ" @change="play"></el-date-picker>
    </el-form-item>
    <el-form-item label="录像">
      <div id="player"></div>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { inject, reactive, toRefs, onMounted, onUnmounted } from "vue"
import { useGet } from "@/http"

const api = inject("api")
const props = defineProps({
  values: Object
})
const state = reactive({
  code: props.values.code,
  date: new Date().toISOString(),
  data: {
    url: "",
  },
  player: null,
})

const play = async () => {
  if (state.player) {
    state.player.pause()
    state.player.destroy()
    state.player = null
  }
  state.data.url = await useGet(`${api}/${state.code}/records/${state.date}`)
  state.player = new WasmPlayer(null, "player")
  state.player.play(state.data.url, 1)
}
onMounted(async () => {
  await play()
})
onUnmounted(async () => {
  state.player.pause()
  state.player.destroy()
  state.player = null
})

const { data } = toRefs(state)
</script>
