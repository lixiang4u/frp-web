import axios from "axios";

const request = axios.create({
  baseURL: import.meta.env.VITE_BASE_URL,
  timeout: 50000,
})


request.interceptors.request.use(value => {
  if (!value.headers['token']) {
    value.headers['token'] = 'frp-token'
  }

  // console.log('[request.data]', value)
  return value
}, error => {
  // console.log('[request.error]', error)
  return Promise.reject(error)
})


request.interceptors.response.use(resp => {
  // console.log('[response.data]', resp)
  if (resp.data.code === 200) {
    return resp
  }
  return Promise.reject(resp.data)
}, error => {
  // console.log('[response.error]', error)
  return Promise.reject(error)
})


export default request
