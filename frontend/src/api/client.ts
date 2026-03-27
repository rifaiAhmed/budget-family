import axios from 'axios'

const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export const api = axios.create({
  baseURL,
  timeout: 15000
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers = {
      ...(config.headers || {}),
      Authorization: `Bearer ${token}`
    }
  }
  return config
})

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
