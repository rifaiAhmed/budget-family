import { api } from './client'
import type { ApiResponse, Paginated } from '../types/api'
import type { Category } from '../types/category'

export async function listCategories(params: { family_id: string; type?: 'income' | 'expense'; page?: number; limit?: number }) {
  const { data } = await api.get<ApiResponse<Paginated<Category>>>('/categories', { params })
  return data.data
}

export async function createCategory(payload: { family_id: string; name: string; type: 'income' | 'expense'; icon?: string }) {
  const { data } = await api.post<ApiResponse<{ category: Category }>>('/categories', payload)
  return data.data.category
}
