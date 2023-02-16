<template>
  <el-form>
    <el-form-item label="资源集合">
      <el-tree ref="tree" :data="resources" :props="{ label: 'name' }" default-expand-all node-key="id"
        show-checkbox></el-tree>
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
      tree: null,
      resources: [],
    })

    const list = async () => {
      state.resources = await useGet("/mgt/sys/resources")
    }
    const save = async () => {
      await usePost(`/mgt/sys/roles/${props.id}/resources`, state.tree.getCheckedNodes())
      context.emit("close")
    }
    onMounted(async () => {
      await list()
      state.tree.setCheckedKeys(await useGet(`/mgt/sys/roles/${props.id}/resources`))
    })

    return {
      ...toRefs(state),
      save,
    }
  },
}
</script>
