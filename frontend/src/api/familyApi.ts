import { api } from './client'

export type Family = {
  id: string
  name: string
  owner_id: string
  created_at: string
}

export async function listFamilies(): Promise<{ items: Family[] }> {
  const res = await api.get('/families')
  return res.data.data
}
