<template>
  <el-form ref="form" :model="data" :rules="rules">
    <el-form-item label="编码" prop="id">
      <el-input v-model="data.id" disabled></el-input>
    </el-form-item>
    <el-form-item label="父编码" prop="parentId">
      <el-input-number v-model="data.parentId" :min="0"></el-input-number>
    </el-form-item>
    <el-form-item label="名称" prop="name">
      <el-input v-model="data.name" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="电话" prop="phone">
      <el-input v-model="data.phone" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="地址" prop="addr">
      <el-input v-model="data.addr" clearable maxlength="100" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="次序" prop="sequence">
      <el-input-number v-model="data.sequence" :min="0"></el-input-number>
    </el-form-item>
    <el-form-item>
      <div class="row-center">
        <el-tooltip placement="top">
          <template #content>目前只支持三级部门</template>
          <el-button type="primary" @click="save">保存</el-button>
        </el-tooltip>
        <el-tooltip placement="top">
          <template #content>目前只支持三级部门</template>
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
    parentId: {
      type: Number,
      required: true,
    },
  },
  setup(props, context) {
    const { state, save } = useEdit(context, {
      id: props.id,
      parentId: props.parentId,
      rules: {
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
