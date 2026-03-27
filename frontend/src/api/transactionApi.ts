import { api } from './client'
import type { ApiResponse, Paginated } from '../types/api'
import type { Transaction } from '../types/transaction'

export type TransactionListParams = {
  family_id: string
  page?: number
  limit?: number
  wallet_id?: string
  category_id?: string
  type?: 'income' | 'expense'
  from?: string
  to?: string
}

export type CreateTransactionPayload = {
  family_id: string
  wallet_id: string
  category_id: string
  amount: string
  type: 'income' | 'expense'
  note?: string
  transaction_date: string
}

export async function listTransactions(params: TransactionListParams) {
  const { data } = await api.get<ApiResponse<Paginated<Transaction>>>('/transactions', { params })
  return data.data
}

export async function createTransaction(payload: CreateTransactionPayload) {
  const { data } = await api.post<ApiResponse<{ transaction: Transaction }>>('/transactions', payload)
  return data.data.transaction
}

export async function getTransaction(id: string) {
  const { data } = await api.get<ApiResponse<{ transaction: Transaction }>>(`/transactions/${id}`)
  return data.data.transaction
}

export async function transactionSummary(params: { family_id: string; from?: string; to?: string }) {
  const { data } = await api.get<ApiResponse<{ income: string; expense: string; net: string }>>('/transactions/summary', { params })
  return data.data
}
