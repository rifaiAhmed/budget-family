import { useQuery } from '@tanstack/react-query'
import * as categoryApi from '../api/categoryApi'

export function useCategories(params: { family_id: string; type?: 'income' | 'expense'; page?: number; limit?: number }) {
  return useQuery({
    queryKey: ['categories', params],
    queryFn: () => categoryApi.listCategories(params),
    enabled: !!params.family_id
  })
}
