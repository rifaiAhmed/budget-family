import { Box, Button, Card, CardContent, MenuItem, Skeleton, Stack, TextField, Typography } from '@mui/material'
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts'
import { useMemo, useState } from 'react'
import PullToRefresh from '../../components/PullToRefresh'
import EmptyState from '../../components/EmptyState'
import MobileSectionTitle from '../../components/MobileSectionTitle'
import BudgetProgress from '../../ui/BudgetProgress'
import { useBudgetUsage, useUpsertBudget } from '../../hooks/useBudgets'
import { useCategories } from '../../hooks/useCategories'
import { useSnackbar } from '../../components/SnackbarProvider'
import { useNavigate } from 'react-router-dom'

function getFamilyId(): string {
  return localStorage.getItem('family_id') || ''
}

export default function BudgetPage() {
  const familyId = getFamilyId()
  const { notify } = useSnackbar()
  const nav = useNavigate()
  const now = new Date()
  const [month, setMonth] = useState<number>(now.getMonth() + 1)
  const [year, setYear] = useState<number>(now.getFullYear())

  const [categoryId, setCategoryId] = useState<string>('')
  const [amount, setAmount] = useState<string>('')

  const q = useBudgetUsage({ family_id: familyId, month, year })
  const upsertM = useUpsertBudget()
  const categoriesQ = useCategories({ family_id: familyId, type: 'expense', page: 1, limit: 200 })

  const categoryById = useMemo(() => {
    const map = new Map<string, string>()
    for (const c of categoriesQ.data?.items || []) map.set(c.id, c.name)
    return map
  }, [categoriesQ.data?.items])

  const pieData = useMemo(() => {
    const items = q.data?.items ?? []
    return items.map((r, idx) => ({ name: categoryById.get(r.category_id) || `Category ${idx + 1}`, value: Number(r.used_amount) || 0 }))
  }, [q.data?.items])

  const pieColors = useMemo(() => {
    const palette = ['#3b82f6', '#22c55e', '#f97316', '#a855f7', '#ef4444', '#14b8a6', '#eab308', '#0ea5e9', '#ec4899']
    const items = q.data?.items ?? []
    return items.map((r, idx) => {
      const s = r.category_id || String(idx)
      let h = 0
      for (let i = 0; i < s.length; i++) h = (h * 31 + s.charCodeAt(i)) >>> 0
      return palette[h % palette.length]
    })
  }, [q.data?.items])

  const onRefresh = async () => {
    await q.refetch()
  }

  if (!familyId) {
    return <EmptyState title="Family belum diset" subtitle="Masukkan Family ID di Settings." />
  }

  return (
    <PullToRefresh onRefresh={onRefresh}>
      <Stack spacing={2}>
        <MobileSectionTitle title="Ringkasan" />

        <Card variant="outlined">
          <CardContent>
            <Stack spacing={1.5}>
              <Stack direction="row" spacing={1.5}>
                <TextField select fullWidth label="Month" value={month} onChange={(e) => setMonth(Number(e.target.value))}>
                  {Array.from({ length: 12 }).map((_, i) => (
                    <MenuItem key={i + 1} value={i + 1}>
                      {new Intl.DateTimeFormat('en-US', { month: 'long' }).format(new Date(2020, i, 1))}
                    </MenuItem>
                  ))}
                </TextField>
                <TextField
                  fullWidth
                  label="Year"
                  value={year}
                  onChange={(e) => setYear(Number(e.target.value) || now.getFullYear())}
                  inputMode="numeric"
                />
              </Stack>

              <Box sx={{ height: 200 }}>
                {q.isLoading ? (
                  <Skeleton variant="rounded" height={200} />
                ) : (
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Tooltip />
                      <Pie data={pieData} dataKey="value" nameKey="name" innerRadius={52} outerRadius={80}>
                        {pieData.map((_, idx) => (
                          <Cell key={idx} fill={pieColors[idx] || '#3b82f6'} />
                        ))}
                      </Pie>
                    </PieChart>
                  </ResponsiveContainer>
                )}
              </Box>

              <Typography variant="body2" sx={{ opacity: 0.7 }}>
                Spending overview by category.
              </Typography>
            </Stack>
          </CardContent>
        </Card>

        <MobileSectionTitle title="Buat budget" />
        <Card variant="outlined">
          <CardContent>
            <Stack spacing={1.5}>
              {categoriesQ.isLoading ? (
                <Skeleton variant="rounded" height={56} />
              ) : (
                <TextField
                  select
                  fullWidth
                  label="Category"
                  value={categoryId}
                  onChange={(e) => setCategoryId(e.target.value)}
                >
                  {(categoriesQ.data?.items || []).map((c) => (
                    <MenuItem key={c.id} value={c.id}>
                      {c.name}
                    </MenuItem>
                  ))}
                </TextField>
              )}

              <TextField
                label="Budget amount"
                inputMode="numeric"
                value={amount}
                onChange={(e) => setAmount((e.target.value || '').replace(/\D/g, ''))}
              />

              <Button
                variant="contained"
                disabled={!categoryId || !amount || upsertM.isPending}
                onClick={async () => {
                  try {
                    await upsertM.mutateAsync({
                      family_id: familyId,
                      category_id: categoryId,
                      amount,
                      month,
                      year
                    })
                    setAmount('')
                    setCategoryId('')
                    notify('Budget saved', 'success')
                    await q.refetch()
                  } catch (e: any) {
                    notify(e?.response?.data?.message || 'Failed to save budget', 'error')
                  }
                }}
              >
                Save
              </Button>
            </Stack>
          </CardContent>
        </Card>

        <MobileSectionTitle title="Budgets" />

        {q.isLoading ? (
          <Stack spacing={1.5}>
            <Skeleton variant="rounded" height={92} />
            <Skeleton variant="rounded" height={92} />
          </Stack>
        ) : q.data?.items?.length ? (
          <Stack spacing={1.5}>
            {q.data.items.map((row) => (
              <BudgetProgress
                key={row.category_id}
                row={row}
                categoryName={categoryById.get(row.category_id)}
                onClick={() => nav(`/budget/${row.category_id}?month=${month}&year=${year}`)}
              />
            ))}
          </Stack>
        ) : (
          <EmptyState title="Belum ada budget" subtitle="Buat budget pertama kamu untuk mulai memantau pengeluaran." />
        )}
      </Stack>
    </PullToRefresh>
  )
}
