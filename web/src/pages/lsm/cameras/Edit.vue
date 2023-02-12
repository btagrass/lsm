<template>
  <el-form ref="form" :model="data" :rules="rules">
    <el-form-item label="编码" prop="id">
      <el-input v-model="data.id" disabled></el-input>
    </el-form-item>
    <el-form-item label="代码" prop="code">
      <el-input v-model="data.code" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="名称" prop="name">
      <el-input v-model="data.name" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="类型" prop="type">
      <el-input v-model="data.type" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="厂商" prop="mfr">
      <el-input v-model="data.mfr" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="型号" prop="model">
      <el-input v-model="data.model" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="固件" prop="firmware">
      <el-input v-model="data.firmwate" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="地址" prop="addr">
      <el-input v-model="data.addr" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="用户名" prop="userName">
      <el-input v-model="data.userName" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="密码" prop="password">
      <el-input v-model="data.password" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item>
      <div class="row-center">
        <el-tooltip placement="top">
          <template #content>Onvif协议：填写地址，用户名，密码等信息<br />28181协议：自动保存摄像头信息</template>
          <el-button type="primary" @click="save">保存</el-button>
        </el-tooltip>
        <el-tooltip placement="top">
          <template #content>Onvif协议：填写地址，用户名，密码等信息<br />28181协议：自动保存摄像头信息</template>
          <el-button type="primary" @click="save(0)">保存增加</el-button>
        </el-tooltip>
      </div>
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
    const { state, save } = useEdit(context, {
      id: props.id,
      rules: {
        code: {
          required: true,
          message: "请输入代码",
          trigger: "blur",
        },
      },
    })

    return {
      ...toRefs(state),
      save,
    }
  },
}
</script>
