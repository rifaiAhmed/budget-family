import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import * as transactionApi from '../api/transactionApi'

export function useTransactions(params: transactionApi.TransactionListParams) {
  return useQuery({
    queryKey: ['transactions', params],
    queryFn: () => transactionApi.listTransactions(params),
    enabled: !!params.family_id
  })
}

export function useTransaction(id: string) {
  return useQuery({
    queryKey: ['transaction', id],
    queryFn: () => transactionApi.getTransaction(id),
    enabled: !!id
  })
}

export function useCreateTransaction() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: transactionApi.createTransaction,
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['transactions'] })
      await qc.invalidateQueries({ queryKey: ['transaction'] })
      await qc.invalidateQueries({ queryKey: ['summary'] })
      await qc.invalidateQueries({ queryKey: ['budgets-usage'] })
      await qc.invalidateQueries({ queryKey: ['wallets'] })
    }
  })
}

export function useTransactionSummary(params: { family_id: string; from?: string; to?: string }) {
  return useQuery({
    queryKey: ['summary', params],
    queryFn: () => transactionApi.transactionSummary(params),
    enabled: !!params.family_id
  })
}
