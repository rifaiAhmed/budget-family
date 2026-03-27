export type Budget = {
  id: string
  family_id: string
  category_id: string
  amount: string
  month: number
  year: number
}

export type BudgetUsage = {
  family_id: string
  category_id: string
  month: number
  year: number
  budget_amount: string
  used_amount: string
  remaining_amount: string
  percentage_used: string
}
