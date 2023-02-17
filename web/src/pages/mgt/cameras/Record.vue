<template>
  <el-form>
    <el-form-item label="日期" prop="date">
      <el-date-picker v-model="data.date" value-format="YYYY-MM-DDTHH:mm:ss.SSSZ" @change="play"></el-date-picker>
    </el-form-item>
    <el-form-item label="录像">
      <div id="player"></div>
    </el-form-item>
  </el-form>
</template>

<script>
import { inject, reactive, toRefs, onMounted, onUnmounted } from "vue"
import { useGet } from "@/http"

export default {
  props: {
    code: {
      type: String,
      required: true,
    },
  },
  setup(props, context) {
    const api = inject("api")

    const state = reactive({
      code: props.code,
      data: {
        date: new Date().toISOString(),
        url: "",
      },
      player: null,
    })

    const play = async () => {
      state.data.url = await useGet(`${api}/${state.code}/records/${state.data.date}`)
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

    return {
      ...toRefs(state),
      play,
    }
  },
}
</script>
