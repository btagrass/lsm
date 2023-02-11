<template>
  <div>
    <div class="row">
      <el-button type="primary" @click="open(0, 'Edit')">增加</el-button>
      <el-button type="primary" @click="list">刷新</el-button>
    </div>
    <div class="row">
      <el-table ref="table" :data="data.records" :tree-props="{ children: 'children' }" border default-expand-all
        row-key="id" @current-change="change">
        <el-table-column label="编码" prop="id" header-align="center"></el-table-column>
        <el-table-column label="名称" prop="name" header-align="center"></el-table-column>
        <el-table-column label="电话" prop="phone" header-align="center"></el-table-column>
        <el-table-column label="地址" prop="addr" header-align="center"></el-table-column>
        <el-table-column label="次序" prop="sequence" header-align="center"></el-table-column>
        <el-table-column label="操作" width="150px" align="center">
          <template #default="scope">
            <el-button type="warning" circle icon="Plus" title="增加下级" @click="open(0, 'Edit', scope.row.id)">
            </el-button>
            <el-button type="primary" circle icon="Plus" title="增加同级" @click="open(0, 'Edit', scope.row.parentId)">
            </el-button>
            <el-button type="primary" circle icon="edit" title="编辑" @click="open(scope.row.id, 'Edit')"></el-button>
            <el-button type="danger" circle icon="delete" title="删除" @click="remove()"></el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-drawer v-model="component.visible" destroy-on-close @close="list">
        <component :id="component.id" :parentId="component.parentId" :is="component.name" @close="close"></component>
      </el-drawer>
    </div>
  </div>
</template>

<script>
import { defineAsyncComponent, toRefs } from "vue"
import { useComponent, useList } from "@/crud"

export default {
  components: {
    Edit: defineAsyncComponent(() => import("./Edit.vue")),
  },
  setup() {
    const { component, open, close } = useComponent()
    const { state, change, select, list, remove } = useList("/mgt/sys/depts")

    return {
      component,
      open,
      close,
      ...toRefs(state),
      change,
      select,
      list,
      remove,
    }
  },
}
</script>
