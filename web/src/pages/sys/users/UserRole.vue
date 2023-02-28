<template>
  <el-form>
    <el-form-item label="角色集合">
      <el-checkbox-group v-model="data">
        <el-checkbox v-for="role in roles" :key="role.id" :label="role.id">{{ role.name }}</el-checkbox>
      </el-checkbox-group>
    </el-form-item>
    <el-form-item>
      <div class="row-center">
        <el-button type="primary" @click="save">保存</el-button>
      </div>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { inject, reactive, toRefs, onMounted } from "vue"
import { useGet, usePost } from "@/http"

const api = inject("api")
const props = defineProps({
  values: Object
})
const emits = defineEmits(["close"])
const state = reactive({
  id: props.values.id,
  data: {},
  roles: [],
})

const save = async () => {
  await usePost(`${api}/${state.id}/roles`, state.data)
  context.emit("close")
}
onMounted(async () => {
  state.roles = await useGet("/mgt/sys/roles")
  state.data = await useGet(`${api}/${state.id}/roles`)
})

const { data, roles } = toRefs(state)
</script>
