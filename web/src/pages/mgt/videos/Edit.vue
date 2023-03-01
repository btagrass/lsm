<template>
  <el-form ref="form" :model="data" :rules="rules">
    <el-form-item label="编码" prop="id">
      <el-input v-model="data.id" disabled></el-input>
    </el-form-item>
    <el-form-item label="名称" prop="name">
      <el-input v-model="data.name" clearable maxlength="50" show-word-limit></el-input>
    </el-form-item>
    <el-form-item label="来源" prop="source">
      <el-input v-model="data.source" clearable maxlength="100" show-word-limit>
        <template #append>
          <el-upload :action="`${env.MGT_URL}/mgt/sys/files/videos`" :headers="{ Authorization: `${user.token}` }"
            accept="video/*" :on-success="(response) => (data.source = response)">
            <el-icon size="20">
              <UploadFilled />
            </el-icon>
          </el-upload>
        </template>
      </el-input>
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
import { useStore } from "vuex"
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
  source: {
    required: true,
    message: "请输入来源或上传来源文件",
    trigger: "blur",
  },
}

const { user } = useStore().state
const env = import.meta.env
</script>
