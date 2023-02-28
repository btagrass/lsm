<template>
  <el-form ref="form" :model="data" :rules="rules">
    <el-form-item label="编码" prop="id">
      <el-input v-model="data.id" disabled></el-input>
    </el-form-item>
    <el-form-item label="名称" prop="name">
      <el-input v-model="data.name" clearable maxlength="50" show-word-limit></el-input>
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

const props = defineProps({
  values: Object
})
const emits = defineEmits(["close"])
const { form, data, save } = useEdit(props.values, emits)
const rules = {
  name: {
    required: true,
    message: "请输入名称",
    trigger: "blur",
  },
}
</script>
