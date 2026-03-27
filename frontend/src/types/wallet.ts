export type Wallet = {
  id: string
  family_id: string
  name: string
  type: 'cash' | 'bank' | 'ewallet' | 'card'
  balance: string
  created_at: string
}
