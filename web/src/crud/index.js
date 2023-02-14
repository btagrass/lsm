import { provide, inject, reactive, toRefs, watch, onMounted } from "vue"
import { ElMessage, ElMessageBox } from "element-plus"
import { useGet, usePost, useDelete } from "@/http"

export function useComponent() {
  const state = reactive({
    component: {
      id: 0,
      name: null,
      visible: false,
    },
  })

  const open = (id, name, props) => {
    Object.assign(state.component, props)
    state.component.id = id
    state.component.name = name
    state.component.visible = true
  }
  const close = () => {
    state.component.id = 0
    state.component.name = null
    state.component.visible = false
  }

  return {
    ...toRefs(state),
    open,
    close,
  }
}

export function useList(api, state, mounted) {
  provide("api", api)

  var o = {
    table: null,
    ids: [],
    params: {
      current: 1,
      size: 10,
    },
    data: {
      records: [],
      total: 0,
    },
  }
  Object.assign(o, state)
  const s = reactive(o)

  const select = (selection) => {
    if (selection instanceof Array) {
      s.ids = selection.map((i) => i.id)
    } else if (selection instanceof Object) {
      if (selection.children) {
        selection.children.forEach((s) => {
          select(s)
        })
      }
      s.ids.push(selection.id)
    }
  }
  const list = async () => {
    s.data = await useGet(api, {
      params: s.params,
    })
  }
  const remove = (selection) => {
    if (!(selection instanceof PointerEvent)) {
      select(selection)
    }
    if (s.ids.length > 0) {
      ElMessageBox.confirm("确定要删除吗？")
        .then(async () => {
          await useDelete(`${api}/${s.ids}`)
          s.params.current = 1
          await list()
        })
        .catch(() => {
          s.table.clearSelection()
          s.ids = []
        })
    } else {
      ElMessage.error("请选择记录")
    }
  }
  watch(s.params, list)
  onMounted(async () => {
    await list()
    if (mounted) {
      await mounted()
    }
  })

  return {
    state: s,
    select,
    list,
    remove,
  }
}

export function useEdit(context, state, mounted) {
  const api = inject("api")

  var o = {
    form: null,
    id: 0,
    data: {},
  }
  Object.assign(o, state)
  const s = reactive(o)

  const edit = async (id) => {
    id = id ?? s.id
    if (id > 0) {
      s.data = await useGet(`${api}/${id}`)
      if (!s.data) {
        ElMessage.error("该记录不存在")
      }
    } else {
      s.data = { ...state, id: 0 }
    }
  }
  const save = (id) => {
    s.form.validate(async (valid) => {
      if (valid) {
        await usePost(api, s.data)
        if (id instanceof PointerEvent) {
          context.emit("close")
        } else {
          edit(id)
        }
      }
    })
  }
  onMounted(async () => {
    await edit()
    if (mounted) {
      await mounted()
    }
  })

  return {
    api,
    state: s,
    save,
  }
}
