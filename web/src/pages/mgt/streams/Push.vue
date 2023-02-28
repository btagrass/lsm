<template>
  <el-form ref="form" :model="data" :rules="rules">
    <el-form-item label="编码" prop="id">
      <el-input v-model="data.id" disabled></el-input>
    </el-form-item>
    <el-form-item label="流名称" prop="name">
      <el-input v-model="data.name" disabled></el-input>
    </el-form-item>
    <el-form-item label="远程地址" prop="remoteAddr">
      <el-input v-model="data.remoteAddr" clearable maxlength="50" placeholder="rtsp://localhost:5544/live/test"
        show-word-limit></el-input>
    </el-form-item>
    <el-form-item>
      <div class="row-center">
        <el-button type="primary" @click="start">开始转推</el-button>
        <el-button type="primary" @click="stop">停止转推</el-button>
      </div>
    </el-form-item>
  </el-form>
</template>

<script>
import { inject, toRefs } from "vue"
import { useEdit } from "@/crud"
import { useGet, usePost } from "@/http"

export default {
  props: {
    v: Object,
  },
  setup(props, context) {
    const api = inject("api")

    const { state } = useEdit(context, {
      ...props.v,
      rules: {
        remoteAddr: {
          required: true,
          message: "请输入远程地址",
          trigger: "blur",
        },
      },
    }, async () => {
      state.data = await useGet(`${api}/${props.v.name}/push`)
      if (!state.data) {
        state.data = { ...state }
      }
    })

    const start = async () => {
      state.form.validate(async (valid) => {
        if (valid) {
          await usePost(`${api}/start`, state.data)
          context.emit("close")
        }
      })
    }
    const stop = async () => {
      state.form.validate(async (valid) => {
        if (valid) {
          await usePost(`${api}/stop`, state.data)
          context.emit("close")
        }
      })
    }

    return {
      ...toRefs(state),
      start,
      stop,
    }
  },
}
</script>
