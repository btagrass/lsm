<template>
  <div>
    <div class="row">
      <el-button type="danger" icon="Delete" @click="remove">删除</el-button>
      <el-button type="primary" icon="DocumentAdd" @click="open(Edit)">增加</el-button>
      <el-button type="primary" icon="Search" @click="list">查询</el-button>
    </div>
    <el-table ref="table" :data="data.records" border @selection-change="select">
      <el-table-column type="selection"></el-table-column>
      <el-table-column label="编码" prop="id" width="120"></el-table-column>
      <el-table-column label="部门" prop="dept.name"></el-table-column>
      <el-table-column label="用户名" prop="userName"></el-table-column>
      <el-table-column label="姓名" prop="fullName"></el-table-column>
      <el-table-column label="手机号码" prop="mobile"></el-table-column>
      <el-table-column label="密码" align="center">
        <template #default>******</template>
      </el-table-column>
      <el-table-column label="冻结" align="center">
        <template #default="scope">{{ scope.row.frozen ? "是" : "否" }}</template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="scope">
          <el-button-group>
            <el-button type="primary" icon="Edit" title="编辑" @click="open(Edit, { id: scope.row.id })"></el-button>
            <el-button type="danger" icon="Delete" title="删除" @click="remove(scope.row)"></el-button>
            <el-button type="warning" icon="Setting" title="角色" @click="open(UserRole, { id: scope.row.id })"></el-button>
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
import UserRole from "./UserRole.vue"

const { name, values, visible, open, close } = useComp()
const { table, params, data, list, remove, select } = useList("/mgt/sys/users")
</script>
