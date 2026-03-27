import { useQuery } from '@tanstack/react-query'
import * as familyApi from '../api/familyApi'

export function useFamilies() {
  return useQuery({
    queryKey: ['families'],
    queryFn: familyApi.listFamilies
  })
}
