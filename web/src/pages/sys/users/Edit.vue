<template>
  <el-form ref="form" :model="data" :rules="rules">
    <el-form-item label="编码" prop="id">
      <el-input v-model="data.id" disabled></el-input>
    </el-form-item>
    <el-form-item label="部门" prop="deptId">
      <el-tree-select v-model="data.deptId" :data="depts" :props="{ label: 'name' }" check-strictly clearable
        node-key="id">
      </el-tree-select>
    </el-form-item>
    <el-form-item label="用户名" prop="userName">
      <el-input v-model="data.userName" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="姓名" prop="fullName">
      <el-input v-model="data.fullName" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="手机" prop="mobile">
      <el-input v-model="data.mobile" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="密码" prop="password">
      <el-input v-model="data.password" clearable maxlength="30" show-password></el-input>
    </el-form-item>
    <el-form-item label="冻结" prop="frozen">
      <el-switch v-model="data.frozen"></el-switch>
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
import { ref } from "vue"
import { useEdit } from "@/crud"
import { useGet } from "@/http"

const props = defineProps({
  values: Object
})
const emits = defineEmits(["close"])
const depts = ref([])
const { form, data, save } = useEdit(props.values, emits, async () => {
  depts.value = await useGet("/mgt/sys/depts")
})
const rules = {
  deptId: {
    required: true,
    message: "请选择部门",
    trigger: "blur",
  },
  userName: {
    required: true,
    message: "请输入用户名",
    trigger: "blur",
  },
  mobile: {
    required: true,
    pattern: "^[1][3-9][0-9]{9}$",
    message: "请输入正确的手机号码",
    trigger: "blur",
  },
  password: {
    required: true,
    message: "请输入密码",
    trigger: "blur",
  },
}
</script>
