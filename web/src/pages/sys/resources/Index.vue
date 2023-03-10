<template>
  <div>
    <div class="row">
      <el-button type="primary" icon="Search" @click="list">查询</el-button>
    </div>
    <el-table ref="table" :data="data.records" :tree-props="{ children: 'children' }" border default-expand-all
      row-key="id">
      <el-table-column label="编码" prop="id" width="170"></el-table-column>
      <el-table-column label="代码" prop="code"></el-table-column>
      <el-table-column label="名称" prop="name"></el-table-column>
      <el-table-column label="类型" prop="type"></el-table-column>
      <el-table-column label="图标" prop="icon"></el-table-column>
      <el-table-column label="链接" prop="url"></el-table-column>
      <el-table-column label="次序" prop="sequence" align="center"></el-table-column>
      <el-table-column label="操作" width="170">
        <template #default="scope">
          <el-button-group>
            <el-button type="primary" icon="Edit" title="编辑"
              @click="open(Edit, { id: scope.row.id, parentId: scope.row.parentId })"></el-button>
            <el-button type="danger" icon="Delete" title="删除" @click="remove(scope.row)"></el-button>
            <el-button type="primary" icon="DocumentAdd" title="增加同级"
              @click="open(Edit, { parentId: scope.row.parentId })"></el-button>
            <el-button type="warning" icon="DocumentAdd" title="增加下级"
              @click="open(Edit, { parentId: scope.row.id })"></el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination v-model:current-page="params.current" v-model:page-size="params.size" :total="data.total" background
      layout="total,prev,pager,next,sizes"></el-pagination>
    <el-drawer v-model="visible" destroy-on-close @close="list">
      <component :is="name" :values="values" @close="close"></component>
    </el-drawer>
  </div>
</template>

<script setup>
import { useComp, useList } from "@/crud"
import Edit from "./Edit.vue"

const { name, values, visible, open, close } = useComp()
const { table, params, data, list, remove } = useList("/mgt/sys/resources")
</script>
