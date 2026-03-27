import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import * as walletApi from '../api/walletApi'

export function useWallets(params: { family_id: string; page?: number; limit?: number }) {
  return useQuery({
    queryKey: ['wallets', params],
    queryFn: () => walletApi.listWallets(params),
    enabled: !!params.family_id
  })
}

export function useCreateWallet() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: walletApi.createWallet,
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['wallets'] })
    }
  })
}

export function useUpdateWallet() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: { name: string; type: 'cash' | 'bank' | 'ewallet' | 'card' } }) =>
      walletApi.updateWallet(id, payload),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['wallets'] })
    }
  })
}

export function useDeleteWallet() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: walletApi.deleteWallet,
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['wallets'] })
    }
  })
}
