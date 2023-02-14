import axios from "axios"
import { ElMessage } from "element-plus"
import store from "@/store"

const http = axios.create({
  baseURL: import.meta.env.VITE_MGT_URL,
  timeout: 15000,
})
http.interceptors.request.use((config) => {
  config.headers.Authorization = store.state.user.token

  return config
})
http.interceptors.response.use(
  (response) => {
    if (response.data.code == 200) {
      return Promise.resolve(response.data.data)
    } else {
      ElMessage.error(response.data.msg)

      return Promise.reject(response.data)
    }
  },
  (error) => {
    ElMessage.error("服务器网络异常")

    return Promise.reject(error)
  }
)

export function useGet(url, params) {
  return http.get(url, params)
}

export function usePost(url, data) {
  return http.post(url, data)
}

export function useDelete(url, data) {
  return http.delete(url, data)
}
