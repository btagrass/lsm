<template>
  <div>
    <div class="row">
      <el-button type="danger" icon="Delete" @click="remove()">删除</el-button>
      <el-button type="primary" icon="DocumentAdd" @click="open(0, 'Edit')">增加</el-button>
      <el-button type="warning" icon="Refresh" @click="list()">刷新</el-button>
    </div>
    <el-table ref="table" :data="data.records" border @selection-change="select">
      <el-table-column type="selection"></el-table-column>
      <el-table-column label="编码" prop="id"></el-table-column>
      <el-table-column label="名称" prop="name"></el-table-column>
      <el-table-column label="操作" width="150px" align="center">
        <template #default="scope">
          <el-button type="primary" circle icon="edit" title="编辑" @click="open(scope.row.id, 'Edit')"></el-button>
          <el-button type="danger" circle icon="delete" title="删除" @click="remove([scope.row.id])"></el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination v-model:current-page="params.current" v-model:page-size="params.size" :total="data.total" background
      layout="total,prev,pager,next,sizes"></el-pagination>
    <el-drawer v-model="component.visible" destroy-on-close @close="list">
      <component :id="component.id" :is="component.name" @close="close"></component>
    </el-drawer>
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
    const { state, select, list, remove } = useList("/mgt/videowalls")

    return {
      component,
      open,
      close,
      ...toRefs(state),
      select,
      list,
      remove,
    }
  },
}
</script>
