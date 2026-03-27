export type Transaction = {
  id: string
  family_id: string
  wallet_id: string
  category_id: string
  amount: string
  type: 'income' | 'expense'
  note?: string
  transaction_date: string
  created_by: string
  created_at: string
}
