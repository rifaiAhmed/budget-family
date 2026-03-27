import { useQuery } from '@tanstack/react-query'
import * as familyApi from '../api/familyApi'

export function useFamilies() {
  return useQuery({
    queryKey: ['families'],
    queryFn: familyApi.listFamilies
  })
}

export function useFamilyMembers(familyId: string) {
  return useQuery({
    queryKey: ['family-members', familyId],
    queryFn: () => familyApi.listFamilyMembers(familyId),
    enabled: !!familyId
  })
}
