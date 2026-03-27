import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import * as budgetApi from '../api/budgetApi'

export function useBudgets(params: { family_id: string; month?: number; year?: number; page?: number; limit?: number }) {
  return useQuery({
    queryKey: ['budgets', params],
    queryFn: () => budgetApi.listBudgets(params),
    enabled: !!params.family_id
  })
}

export function useBudgetUsage(params: { family_id: string; month: number; year: number }) {
  return useQuery({
    queryKey: ['budgets-usage', params],
    queryFn: () => budgetApi.budgetUsage(params),
    enabled: !!params.family_id && !!params.month && !!params.year
  })
}

export function useUpsertBudget() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: budgetApi.upsertBudget,
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['budgets'] })
      await qc.invalidateQueries({ queryKey: ['budgets-usage'] })
    }
  })
}
