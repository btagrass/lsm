<template>
  <el-form>
    <el-form-item label="实况">
      <div id="player"></div>
    </el-form-item>
  </el-form>
</template>

<script>
import { inject, reactive, toRefs, onMounted, onUnmounted } from "vue"
import { usePost } from "@/http"

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

    return {
      ...toRefs(state),
    }
  },
}
</script>
