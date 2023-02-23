<template>
  <el-form ref="form" :model="data" :rules="rules">
    <el-form-item label="编码" prop="id">
      <el-input v-model="data.id" disabled></el-input>
    </el-form-item>
    <el-form-item label="应用名称" prop="appName">
      <el-input v-model="data.appName" disabled></el-input>
    </el-form-item>
    <el-form-item label="名称" prop="name">
      <el-input v-model="data.name" disabled></el-input>
    </el-form-item>
    <el-form-item label="音频解码" prop="audioCodec">
      <el-input v-model="data.audioCodec" disabled></el-input>
    </el-form-item>
    <el-form-item label="视频解码" prop="videoCodec">
      <el-input v-model="data.videoCodec" disabled></el-input>
    </el-form-item>
    <el-form-item label="视频宽" prop="videoWidth">
      <el-input v-model="data.videoWidth" disabled></el-input>
    </el-form-item>
    <el-form-item label="视频高" prop="videoHeight">
      <el-input v-model="data.videoHeight" disabled></el-input>
    </el-form-item>
    <el-form-item label="会话" prop="session">
      <el-input v-model="data.session" disabled></el-input>
    </el-form-item>
    <el-form-item label="协议" prop="protocol">
      <el-input v-model="data.protocol" disabled></el-input>
    </el-form-item>
    <el-form-item label="类型" prop="type">
      <el-input v-model="data.type" disabled></el-input>
    </el-form-item>
    <el-form-item label="远程地址" prop="remoteAddr">
      <el-input v-model="data.remoteAddr" disabled></el-input>
    </el-form-item>
    <el-form-item label="码率" prop="codeRate">
      <el-input v-model="data.codeRate" disabled>
        <template #append>Kb</template>
      </el-input>
    </el-form-item>
    <el-form-item label="接收字节数" prop="receivedBytes">
      <el-input v-model="data.receivedBytes" disabled>
        <template #append>Kb</template>
      </el-input>
    </el-form-item>
    <el-form-item label="发送字节数" prop="sentBytes">
      <el-input v-model="data.sentBytes" disabled>
        <template #append>Kb</template>
      </el-input>
    </el-form-item>
  </el-form>
</template>

<script>
import { toRefs } from "vue"
import { useEdit } from "@/crud"

export default {
  props: {
    id: {
      type: Number,
      required: true,
    },
  },
  setup(props, context) {
    const { state } = useEdit(context, {
      id: props.id,
    })

    return {
      ...toRefs(state),
    }
  },
}
</script>
