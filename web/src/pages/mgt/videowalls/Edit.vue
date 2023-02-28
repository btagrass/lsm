<template>
  <el-form ref="form" :model="data" :rules="rules">
    <el-form-item label="编码" prop="id">
      <el-input v-model="data.id" disabled></el-input>
    </el-form-item>
    <el-form-item label="名称" prop="name">
      <el-input v-model="data.name" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="摄像头集合" prop="cameras">
      <el-select v-model="data.cameras" :multiple-limit="4" multiple>
        <el-option v-for="camera in cameras" :key="camera.Code" :label="camera.name" :value="camera.code"></el-option>
      </el-select>
    </el-form-item>
    <el-form-item>
      <div class="row-center">
        <el-button type="primary" @click="save">保存</el-button>
        <el-button type="primary" @click="save(0)">保存增加</el-button>
      </div>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { useEdit } from "@/crud"
import { useGet } from "@/http"

const props = defineProps({
  values: Object
})
const emits = defineEmits(["close"])
const cameras = ref([])
const { form, data, save } = useEdit(props.values, emits, async () => {
  cameras.value = await useGet("/mgt/cameras")
})
const rules = {
  name: {
    required: true,
    message: "请输入名称",
    trigger: "blur",
  },
  cameras: {
    required: true,
    message: "请选择摄像头",
    trigger: "blur",
  },
}
</script>
