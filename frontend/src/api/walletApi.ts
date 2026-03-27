import { api } from './client'
import type { ApiResponse, Paginated } from '../types/api'
import type { Wallet } from '../types/wallet'

export type CreateWalletPayload = {
  family_id: string
  name: string
  type: 'cash' | 'bank' | 'ewallet' | 'card'
  balance?: string
}

export async function listWallets(params: { family_id: string; page?: number; limit?: number }) {
  const { data } = await api.get<ApiResponse<Paginated<Wallet>>>('/wallets', { params })
  return data.data
}

export async function createWallet(payload: CreateWalletPayload) {
  const { data } = await api.post<ApiResponse<{ wallet: Wallet }>>('/wallets', payload)
  return data.data.wallet
}

export async function updateWallet(id: string, payload: { name: string; type: 'cash' | 'bank' | 'ewallet' | 'card' }) {
  const { data } = await api.put<ApiResponse<{ wallet: Wallet }>>(`/wallets/${id}`, payload)
  return data.data.wallet
}

export async function deleteWallet(id: string) {
  const { data } = await api.delete<ApiResponse<Record<string, never>>>(`/wallets/${id}`)
  return data
}
