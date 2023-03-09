<template>
    <el-row v-for="r in rows" :key="r">
        <el-col v-for="c in cols" :key="c" :span="24 / cols">
            <div :id="`player${(r - 1) * cols + (c - 1)}`"></div>
        </el-col>
    </el-row>
</template>

<script setup>
import { reactive, toRefs, onMounted, onUpdated, onUnmounted } from "vue"
import { useGet } from "@/http"

const state = reactive({
    urls: [],
    rows: 0,
    cols: 0,
    players: [],
})

onMounted(async () => {
    state.urls = await useGet("/mgt/videowalls/default")
    state.rows = Math.round(Math.sqrt(state.urls.length))
    state.cols = Math.ceil(Math.sqrt(state.urls.length))
})
onUpdated(async () => {
    for (var i = 0; i < state.urls.length; i++) {
        const player = new WasmPlayer("", `player${i}`)
        player.play(state.urls[i], 1)
        state.players.push(player)
    }
})
onUnmounted(async () => {
    for (player of state.players) {
        player.pause()
        player.destroy()
        player = null
    }
})

const { rows, cols } = toRefs(state)
</script>
