<template>
  <div>
    <div class="row">
      <el-button type="warning" icon="Search" @click="list">查询</el-button>
    </div>
    <el-table ref="table" :data="data.records" border @selection-change="select">
      <el-table-column type="selection"></el-table-column>
      <el-table-column label="编码" prop="id" width="120"></el-table-column>
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
            <el-button type="primary" icon="View" title="查看" @click="open(scope.row.id, 'View')"></el-button>
            <el-button type="warning" icon="RefreshRight" title="转推"
              @click="open(scope.row.id, 'Push', { streamName: scope.row.name })"></el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination v-model:current-page="params.current" v-model:page-size="params.size" :total="data.total" background
      layout="total,prev,pager,next,sizes"></el-pagination>
    <el-drawer v-model="component.visible" destroy-on-close @close="list">
      <component :id="component.id" :streamName="component.streamName" :is="component.name" @close="close"></component>
    </el-drawer>
  </div>
</template>

<script>
import { defineAsyncComponent, toRefs } from "vue"
import { useComponent, useList } from "@/crud"

export default {
  components: {
    Push: defineAsyncComponent(() => import("./Push.vue")),
    View: defineAsyncComponent(() => import("./View.vue")),
  },
  setup() {
    const { component, open, close } = useComponent()
    const { state, select, list } = useList("/mgt/streams")

    return {
      component,
      open,
      close,
      ...toRefs(state),
      select,
      list,
    }
  },
}
</script>
