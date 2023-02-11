<template>
  <el-form ref="form" :model="data" :rules="rules">
    <el-form-item label="编码" prop="id">
      <el-input v-model="data.id" disabled></el-input>
    </el-form-item>
    <el-form-item label="名称" prop="name">
      <el-input v-model="data.name" maxlength="20" clearable></el-input>
    </el-form-item>
    <el-form-item label="摄像头集合" prop="cameras">
      <el-select v-model="data.cameras" placeholder="请选择" multiple>
        <el-option v-for="camera in cameras" :key="camera.Code" :label="camera.name" :value="camera.code"></el-option>
      </el-select>
    </el-form-item>
    <el-form-item>
      <div class="row-center">
        <el-button type="primary" @click="save()">保存</el-button>
        <el-button type="primary" @click="save(true)">保存增加</el-button>
      </div>
    </el-form-item>
  </el-form>
</template>

<script>
import { toRefs } from "vue"
import { useEdit } from "@/crud"
import { useGet } from "@/http"

export default {
  props: {
    id: {
      type: Number,
      required: true,
    },
  },
  setup(props, context) {
    const { state, save } = useEdit(
      context,
      {
        id: props.id,
        rules: {
          name: {
            required: true,
            message: "请输入名称",
            trigger: "blur",
          },
          cameras: {
            type: "array",
            min: 1,
            max: 4,
            message: "请选择摄像头, 数量在1-6之间",
            trigger: "change",
          },
        },
        cameras: [],
      },
      async () => {
        state.cameras = (
          await useGet("/mgt/cameras", {
            params: {
              current: 1,
              size: 1000,
            },
          })
        ).records
      },
    )

    return {
      ...toRefs(state),
      save,
    }
  },
}
</script>
