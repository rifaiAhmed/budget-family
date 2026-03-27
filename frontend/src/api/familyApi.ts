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

export type FamilyMember = {
  id: string
  family_id: string
  user_id: string
  role: string
  name: string
  email: string
}

export async function listFamilyMembers(familyId: string): Promise<{ items: FamilyMember[] }> {
  const res = await api.get(`/families/${familyId}/members`)
  return res.data.data
}

export async function joinFamily(familyId: string): Promise<{ family: Family }> {
  const res = await api.post('/families/join', { family_id: familyId })
  return res.data.data
}

export async function createFamily(name: string): Promise<{ family: Family }> {
  const res = await api.post('/families', { name })
  return res.data.data
}
