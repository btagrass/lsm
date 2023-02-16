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
    <el-form-item label="类型" prop="type">
      <el-select v-model="data.type" placeholder="请选择" filterable>
        <el-option v-for="dict in dicts" :key="dict.code" :label="dict.name" :value="dict.code"></el-option>
      </el-select>
    </el-form-item>
    <el-form-item label="图标" prop="icon">
      <el-input v-model="data.icon" clearable maxlength="50" placeholder="Element图标" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="链接" prop="url">
      <el-input v-model="data.url" clearable maxlength="100" placeholder="/user 或 /https://www.baidu.com"
        show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="动作" prop="act">
      <el-input v-model="data.act" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="次序" prop="sequence">
      <el-input-number v-model="data.sequence" :min="0"></el-input-number>
    </el-form-item>
    <el-form-item>
      <div class="row-center">
        <el-tooltip placement="top">
          <template #content>目前只支持三级资源</template>
          <el-button type="primary" @click="save">保存</el-button>
        </el-tooltip>
        <el-tooltip placement="top">
          <template #content>目前只支持三级资源</template>
          <el-button type="primary" @click="save(0)">保存增加</el-button>
        </el-tooltip>
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
    parentId: {
      type: Number,
      required: false,
    },
  },
  setup(props, context) {
    const { state, save } = useEdit(
      context,
      {
        id: props.id,
        parentId: props.parentId,
        rules: {
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
          type: {
            required: true,
            message: "请选择类型",
            trigger: "blur",
          },
          url: {
            pattern: "^/\\S*$",
            message: "请输入正确的链接",
            trigger: "blur",
          },
        },
        dicts: [],
      },
      async () => {
        state.dicts = await useGet("/mgt/sys/dicts?type=Resource")
      }
    )

    return {
      ...toRefs(state),
      save,
    }
  },
}
</script>
