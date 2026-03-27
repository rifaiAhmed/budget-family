import { api } from './client'
import type { ApiResponse, TokenPair } from '../types/api'
import type { User } from '../types/user'

export type LoginPayload = {
  email: string
  password: string
}

export type RegisterPayload = {
  name: string
  email: string
  phone?: string
  password: string
}

export async function login(payload: LoginPayload) {
  const { data } = await api.post<ApiResponse<{ user: User; tokens: TokenPair }>>('/auth/login', payload)
  return data.data
}

export async function register(payload: RegisterPayload) {
  const { data } = await api.post<ApiResponse<{ user: User; tokens: TokenPair }>>('/auth/register', payload)
  return data.data
}

export async function me() {
  const { data } = await api.get<ApiResponse<{ user: User }>>('/auth/me')
  return data.data.user
}
