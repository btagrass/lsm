<template>
  <div>
    <div class="row">
      <el-button type="warning" icon="Search" @click="list">查询</el-button>
    </div>
    <el-table ref="table" :data="data.records" border>
      <el-table-column label="应用名称" prop="appName"></el-table-column>
      <el-table-column label="名称" prop="name"></el-table-column>
      <el-table-column label="音频" prop="audioCodec" align="center"></el-table-column>
      <el-table-column label="视频" prop="videoCodec" align="center" width="120">
        <template #default="scope">
          {{ scope.row.videoCodec }}({{ scope.row.videoWidth }} * {{ scope.row.videoHeight }})
        </template>
      </el-table-column>
      <el-table-column label="会话" prop="session"></el-table-column>
      <el-table-column label="协议" prop="protocol" align="center"></el-table-column>
      <el-table-column label="类型" prop="type" align="center"></el-table-column>
      <el-table-column label="远程地址" prop="remoteAddr"></el-table-column>
      <el-table-column label="码率 (Kb)" prop="codeRate" align="right"></el-table-column>
      <el-table-column label="接收字节数 (Kb)" prop="receivedBytes" align="right"></el-table-column>
      <el-table-column label="发送字节数 (Kb)" prop="sentBytes" align="right"></el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="scope">
          <el-button-group>
            <el-button type="primary" icon="VideoPlay" title="直播"
              @click="open(Live, { name: scope.row.name })"></el-button>
            <el-button type="warning" icon="RefreshRight" title="转推"
              @click="open(Push, { name: scope.row.name })"></el-button>
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
import Live from "./Live.vue"
import Push from "./Push.vue"

const { name, values, visible, open, close } = useComp()
const { table, params, data, list } = useList("/mgt/streams")
</script>
