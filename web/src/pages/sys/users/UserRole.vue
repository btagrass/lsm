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

<script>
import { reactive, toRefs, onMounted } from "vue"
import { useGet, usePost } from "@/http"

export default {
  props: {
    id: {
      type: Number,
      required: true,
    },
  },
  setup(props, context) {
    const state = reactive({
      data: {},
      roles: [],
    })
    const list = async () => {
      state.roles = await useGet("/mgt/sys/roles")
    }
    const save = async () => {
      await usePost(`/mgt/sys/users/${props.id}/roles`, state.data)
      context.emit("close")
    }
    onMounted(async () => {
      await list()
      state.data = await useGet(`/mgt/sys/users/${props.id}/roles`)
    })

    return {
      ...toRefs(state),
      save,
    }
  },
}
</script>
