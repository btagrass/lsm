<template>
  <div>
    <div class="row">
      <el-button type="warning" icon="Search" @click="list">查询</el-button>
    </div>
    <el-table ref="table" :data="data.records" :tree-props="{ children: 'children' }" border default-expand-all
      row-key="id">
      <el-table-column label="编码" prop="id" width="170"></el-table-column>
      <el-table-column label="名称" prop="name"></el-table-column>
      <el-table-column label="电话" prop="phone"></el-table-column>
      <el-table-column label="地址" prop="addr"></el-table-column>
      <el-table-column label="次序" prop="sequence" align="center"></el-table-column>
      <el-table-column label="操作" width="170">
        <template #default="scope">
          <el-button-group>
            <el-button type="primary" icon="Edit" title="编辑"
              @click="open(scope.row.id, 'Edit', { parentId: scope.row.parentId })"></el-button>
            <el-button type="danger" icon="Delete" title="删除" @click="remove(scope.row)"></el-button>
            <el-button type="primary" icon="DocumentAdd" title="增加同级"
              @click="open(0, 'Edit', { parentId: scope.row.parentId })">
            </el-button>
            <el-button type="warning" icon="DocumentAdd" title="增加下级"
              @click="open(0, 'Edit', { parentId: scope.row.id })">
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination v-model:current-page="params.current" v-model:page-size="params.size" :total="data.total" background
      layout="total,prev,pager,next,sizes"></el-pagination>
    <el-drawer v-model="component.visible" destroy-on-close @close="list">
      <component :id="component.id" :parentId="component.parentId" :is="component.name" @close="close"></component>
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
    const { state, list, remove } = useList("/mgt/sys/depts")
    state.params.size = 1

    return {
      component,
      open,
      close,
      ...toRefs(state),
      list,
      remove,
    }
  },
}
</script>
