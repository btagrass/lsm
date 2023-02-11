<template>
  <div class="background center">
    <el-form ref="form" :model="data" :rules="rules" :show-message="false" @keyup.enter="login">
      <el-form-item prop="userName">
        <el-input v-model="data.userName" maxlength="20" prefix-icon="user"></el-input>
      </el-form-item>
      <el-form-item prop="password">
        <el-input v-model="data.password" maxlength="50" show-password prefix-icon="lock" @keyup.enter="login"></el-input>
      </el-form-item>
      <el-form-item>
        <div class="row-center">
          <el-button type="primary" @click="login">登录</el-button>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { reactive, toRefs } from "vue"
import { useRouter } from "vue-router"
import { useStore } from "vuex"
import { ElMessage } from "element-plus"
import { useGet, usePost } from "@/http"

export default {
  setup() {
    const router = useRouter()
    const store = useStore()

    const state = reactive({
      form: null,
      data: {
        userName: "admin",
        password: "",
      },
      rules: {
        userName: {
          required: true,
          message: "请输入用户名",
          trigger: "blur",
        },
        password: {
          required: true,
          message: "请输入密码",
          trigger: "blur",
        },
      },
    })
    const login = () => {
      state.form.validate(async (valid) => {
        if (valid) {
          const user = await usePost("/mgt/login", state.data)
          if (user) {
            store.commit("saveUser", {
              userName: user.userName,
              token: user.token,
            })
            store.commit("saveResources", await useGet("/mgt/sys/resources/menu"))
            router.push("/")
          } else {
            ElMessage.error("用户名或密码错误")
          }
        } else {
          ElMessage.error("请输入用户名或密码")
        }
      })
    }

    return {
      ...toRefs(state),
      login,
    }
  },
}
</script>
