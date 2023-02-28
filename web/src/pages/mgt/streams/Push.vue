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

<script setup>
import { inject, reactive, toRefs, onMounted } from "vue"
import { useGet, usePost } from "@/http"

const api = inject("api")
const props = defineProps({
  values: Object
})
const emits = defineEmits(["close"])
const state = reactive({
  form: null,
  name: props.values.name,
  data: {},
})
const rules = {
  remoteAddr: {
    required: true,
    message: "请输入远程地址",
    trigger: "blur",
  },
}

const start = async () => {
  state.form.validate(async (valid) => {
    if (valid) {
      await usePost(`${api}/start`, data)
      emits("close")
    }
  })
}
const stop = async () => {
  state.form.validate(async (valid) => {
    if (valid) {
      await usePost(`${api}/stop`, data)
      emits("close")
    }
  })
}
onMounted(async () => {
  state.data = await useGet(`${api}/${state.name}/push`)
  if (!state.data) {
    state.data = { ...state }
  }
})

const { form, data } = toRefs(state)
</script>

<!-- <script>
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
</script> -->
