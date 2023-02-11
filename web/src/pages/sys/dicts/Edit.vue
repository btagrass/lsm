<template>
  <el-form ref="form" :model="data" :rules="rules">
    <el-form-item label="编码" prop="id">
      <el-input v-model="data.id" disabled></el-input>
    </el-form-item>
    <el-form-item label="类型" prop="type">
      <el-input v-model="data.type" maxlength="20" clearable></el-input>
    </el-form-item>
    <el-form-item label="代码" prop="code">
      <el-input-number v-model="data.code" :min="0"></el-input-number>
    </el-form-item>
    <el-form-item label="名称" prop="name">
      <el-input v-model="data.name" maxlength="50" clearable></el-input>
    </el-form-item>
    <el-form-item label="次序" prop="sequence">
      <el-input-number v-model="data.sequence" :min="0"></el-input-number>
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
        type: {
          required: true,
          message: "请输入类型",
          trigger: "blur",
        },
        code: {
          required: true,
          message: "请输入代码",
          trigger: "blur",
        },
        name: {
          required: true,
          message: "请输入名称",
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
