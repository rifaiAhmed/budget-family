import { api } from './client'
import type { ApiResponse, Paginated } from '../types/api'
import type { Budget, BudgetUsage } from '../types/budget'

// LIST BUDGETS
export async function listBudgets(params: {
  family_id: string
  month?: number
  year?: number
  page?: number
  limit?: number
}) {
  const { data } = await api.get<ApiResponse<Paginated<Budget>>>(
    '/budgets',
    { params }
  )

  return data.data
}

// BUDGET USAGE
export async function budgetUsage(params: {
  family_id: string
  month: number
  year: number
}) {
  const { data } = await api.get<
    ApiResponse<{
      items: BudgetUsage[]
      month: number
      year: number
    }>
  >('/budgets/usage', { params })

  return data.data
}

// CREATE / UPDATE (UPSERT) BUDGET
export async function upsertBudget(payload: {
  family_id: string
  category_id: string
  amount: string
  month: number
  year: number
}) {
  const { data } = await api.post<ApiResponse<{ budget: Budget }>>(
    '/budgets',
    payload
  )

  return data.data.budget
}
