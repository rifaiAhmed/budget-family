import { Box, Button, Skeleton, Stack, Typography } from '@mui/material'
import { useMemo } from 'react'
import { useNavigate, useParams, useSearchParams } from 'react-router-dom'
import EmptyState from '../../components/EmptyState'
import MobileSectionTitle from '../../components/MobileSectionTitle'
import PullToRefresh from '../../components/PullToRefresh'
import { formatMoneyIDR } from '../../components/Money'
import { useBudgetUsage } from '../../hooks/useBudgets'
import { useCategories } from '../../hooks/useCategories'
import { useTransactions } from '../../hooks/useTransactions'
import TransactionItem from '../../ui/TransactionItem'

function getFamilyId(): string {
  return localStorage.getItem('family_id') || ''
}

function monthRange(month: number, year: number) {
  const from = new Date(year, month - 1, 1)
  const to = new Date(year, month, 0)
  const fmt = (d: Date) => d.toISOString().slice(0, 10)
  return { from: fmt(from), to: fmt(to) }
}

export default function BudgetDetailPage() {
  const familyId = getFamilyId()
  const nav = useNavigate()
  const params = useParams()
  const categoryId = params.categoryId || ''
  const [sp] = useSearchParams()

  const now = new Date()
  const month = Number(sp.get('month') || now.getMonth() + 1)
  const year = Number(sp.get('year') || now.getFullYear())
  const range = useMemo(() => monthRange(month, year), [month, year])

  const usageQ = useBudgetUsage({ family_id: familyId, month, year })
  const categoriesQ = useCategories({ family_id: familyId, type: 'expense', page: 1, limit: 200 })
  const txQ = useTransactions({ family_id: familyId, category_id: categoryId, type: 'expense', from: range.from, to: range.to, page: 1, limit: 50 })

  const categoryName = useMemo(() => {
    const c = (categoriesQ.data?.items || []).find((x) => x.id === categoryId)
    return c?.name || 'Budget'
  }, [categoriesQ.data?.items, categoryId])

  const row = useMemo(() => {
    return (usageQ.data?.items || []).find((r) => r.category_id === categoryId)
  }, [usageQ.data?.items, categoryId])

  const onRefresh = async () => {
    await Promise.all([usageQ.refetch(), categoriesQ.refetch(), txQ.refetch()])
  }

  if (!familyId) return <EmptyState title="Family belum diset" subtitle="Masukkan Family ID di Settings." />
  if (!categoryId) return <EmptyState title="Budget tidak ditemukan" subtitle="Kembali ke halaman Budget dan pilih salah satu budget." />

  const loading = usageQ.isLoading || categoriesQ.isLoading || txQ.isLoading

  return (
    <PullToRefresh onRefresh={onRefresh}>
      <Stack spacing={2}>
        <MobileSectionTitle
          title={categoryName}
          right={
            <Button size="small" onClick={() => nav(-1)} sx={{ fontWeight: 800 }}>
              Kembali
            </Button>
          }
        />

        {loading ? (
          <Skeleton variant="rounded" height={92} />
        ) : row ? (
          <Box>
            <Typography sx={{ fontWeight: 900 }}>{formatMoneyIDR(row.remaining_amount)} sisa</Typography>
            <Typography variant="body2" sx={{ opacity: 0.7 }}>
              Budget {formatMoneyIDR(row.budget_amount)} · Terpakai {formatMoneyIDR(row.used_amount)}
            </Typography>
          </Box>
        ) : (
          <EmptyState title="Budget tidak aktif" subtitle="Budget ini mungkin sudah expired atau belum dibuat untuk bulan ini." />
        )}

        <MobileSectionTitle title="Dipakai untuk" />

        {txQ.isLoading ? (
          <Stack spacing={1.5}>
            <Skeleton variant="rounded" height={78} />
            <Skeleton variant="rounded" height={78} />
          </Stack>
        ) : txQ.data?.items?.length ? (
          <Stack spacing={1.5}>
            {txQ.data.items.map((t) => (
              <TransactionItem key={t.id} tx={t} />
            ))}
          </Stack>
        ) : (
          <EmptyState title="Belum ada pengeluaran" subtitle="Belum ada transaksi yang memakai budget ini pada periode ini." />
        )}
      </Stack>
    </PullToRefresh>
  )
}
