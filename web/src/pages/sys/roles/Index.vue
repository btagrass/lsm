<template>
  <div>
    <div class="row">
      <el-button type="danger" icon="Delete" @click="remove">删除</el-button>
      <el-button type="primary" icon="DocumentAdd" @click="open(0, 'Edit')">增加</el-button>
      <el-button type="warning" icon="Refresh" @click="list">刷新</el-button>
    </div>
    <el-table ref="table" :data="data.records" border @selection-change="select">
      <el-table-column type="selection"></el-table-column>
      <el-table-column label="编码" prop="id" width="120"></el-table-column>
      <el-table-column label="名称" prop="name"></el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="scope">
          <el-button-group>
            <el-button type="primary" icon="Edit" title="编辑" @click="open(scope.row.id, 'Edit')"></el-button>
            <el-button type="danger" icon="Delete" title="删除" @click="remove(scope.row.id)"></el-button>
            <el-button type="warning" icon="Setting" title="资源" @click="open(scope.row.id, 'RoleResource')">
            </el-button>
          </el-button-group>
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
    RoleResource: defineAsyncComponent(() => import("./RoleResource.vue")),
  },
  setup() {
    const { component, open, close } = useComponent()
    const { state, select, list, remove } = useList("/mgt/sys/roles")

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
