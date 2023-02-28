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
      <el-input v-model="data.icon" clearable maxlength="50" placeholder="ElementPlus 图标" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="统一资源标识符" prop="uri">
      <el-input v-model="data.uri" clearable maxlength="100" show-word-limit></el-input>
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

<script setup>
import { ref } from "vue"
import { useEdit } from "@/crud"
import { useGet } from "@/http"

const props = defineProps({
  values: Object
})
const emits = defineEmits(["close"])
const dicts = ref([])
const { form, data, save } = useEdit(props.values, emits, async () => {
  dicts.value = await useGet("/mgt/sys/dicts?type=Resource")
})
const rules = {
  name: {
    required: true,
    message: "请输入名称",
    trigger: "blur",
  },
  uri: {
    required: true,
    message: "请输入统一资源标识符",
    trigger: "blur",
  },
}
</script>
