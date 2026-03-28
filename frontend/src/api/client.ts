import axios, { AxiosHeaders, InternalAxiosRequestConfig } from 'axios'

const baseURL = import.meta.env.VITE_API_URL || 'https://budget-api.medgroup.my.id'

export const api = axios.create({
  baseURL,
  timeout: 15000,
})

// REQUEST INTERCEPTOR
api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = localStorage.getItem('access_token')

  if (token) {
    // pastikan headers adalah AxiosHeaders
    config.headers = new AxiosHeaders(config.headers)

    config.headers.set('Authorization', `Bearer ${token}`)
  }

  return config
})

// RESPONSE INTERCEPTOR
api.interceptors.response.use(
  (res) => res,
  (err) => {
    const status = err?.response?.status

    if (status === 401) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    }

    return Promise.reject(err)
  }
)
