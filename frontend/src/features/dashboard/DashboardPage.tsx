import { Box, Button, Skeleton, Stack, Typography } from '@mui/material'
import TrendingUpIcon from '@mui/icons-material/TrendingUp'
import TrendingDownIcon from '@mui/icons-material/TrendingDown'
import AccountBalanceWalletIcon from '@mui/icons-material/AccountBalanceWallet'
import { useMemo } from 'react'
import StatCard from '../../ui/StatCard'
import TransactionItem from '../../ui/TransactionItem'
import BudgetProgress from '../../ui/BudgetProgress'
import MobileSectionTitle from '../../components/MobileSectionTitle'
import PullToRefresh from '../../components/PullToRefresh'
import EmptyState from '../../components/EmptyState'
import { formatMoneyIDR } from '../../components/Money'
import { useTransactionSummary, useTransactions } from '../../hooks/useTransactions'
import { useBudgetUsage } from '../../hooks/useBudgets'
import { useNavigate } from 'react-router-dom'
import { useCategories } from '../../hooks/useCategories'

function getFamilyId(): string {
  return localStorage.getItem('family_id') || (import.meta.env.VITE_DEFAULT_FAMILY_ID as string) || ''
}

function monthRange() {
  const now = new Date()
  const from = new Date(now.getFullYear(), now.getMonth(), 1)
  const to = new Date(now.getFullYear(), now.getMonth() + 1, 0)
  const fmt = (d: Date) => d.toISOString().slice(0, 10)
  return { from: fmt(from), to: fmt(to), month: now.getMonth() + 1, year: now.getFullYear() }
}

export default function DashboardPage() {
  const familyId = getFamilyId()
  const range = useMemo(() => monthRange(), [])
  const nav = useNavigate()

  const summaryQ = useTransactionSummary({ family_id: familyId, from: range.from, to: range.to })
  const txQ = useTransactions({ family_id: familyId, page: 1, limit: 5 })
  const usageQ = useBudgetUsage({ family_id: familyId, month: range.month, year: range.year })
  const categoriesQ = useCategories({ family_id: familyId, type: 'expense', page: 1, limit: 200 })

  const categoryById = useMemo(() => {
    const map = new Map<string, string>()
    for (const c of categoriesQ.data?.items || []) map.set(c.id, c.name)
    return map
  }, [categoriesQ.data?.items])

  const onRefresh = async () => {
    await Promise.all([summaryQ.refetch(), txQ.refetch(), usageQ.refetch()])
  }

  if (!familyId) return <EmptyState title="Family belum diset" subtitle="Masukkan Family ID di Settings." />

  const income = summaryQ.data?.income ?? '0'
  const expense = summaryQ.data?.expense ?? '0'
  const net = summaryQ.data?.net ?? '0'

  return (
    <PullToRefresh onRefresh={onRefresh}>
      <Stack spacing={2.25}>
        <Stack spacing={1.5}>
          {summaryQ.isLoading ? (
            <Skeleton variant="rounded" height={92} />
          ) : (
            <StatCard title="Income" amount={formatMoneyIDR(income)} icon={<TrendingUpIcon color="success" />} />
          )}

          {summaryQ.isLoading ? (
            <Skeleton variant="rounded" height={92} />
          ) : (
            <StatCard title="Expense" amount={formatMoneyIDR(expense)} icon={<TrendingDownIcon color="error" />} />
          )}

          {summaryQ.isLoading ? (
            <Skeleton variant="rounded" height={92} />
          ) : (
            <StatCard title="Balance" amount={formatMoneyIDR(net)} icon={<AccountBalanceWalletIcon color="primary" />} />
          )}
        </Stack>

        <Box>
          <MobileSectionTitle title="Budget" />
          {usageQ.isLoading ? (
            <Stack spacing={1.5}>
              <Skeleton variant="rounded" height={92} />
              <Skeleton variant="rounded" height={92} />
            </Stack>
          ) : usageQ.data?.items?.length ? (
            <Stack spacing={1.5}>
              {usageQ.data.items.slice(0, 3).map((row) => (
                <BudgetProgress key={row.category_id} row={row} categoryName={categoryById.get(row.category_id)} />
              ))}
            </Stack>
          ) : (
            <EmptyState title="Belum ada budget" subtitle="Buat budget untuk mulai memantau pengeluaran." />
          )}
        </Box>

        <Box>
          <MobileSectionTitle
            title="Transaksi terbaru"
            right={
              <Button size="small" onClick={() => nav('/transactions')} sx={{ fontWeight: 800 }}>
                Lihat semua
              </Button>
            }
          />
          {txQ.isLoading ? (
            <Stack spacing={1.5}>
              <Skeleton variant="rounded" height={78} />
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
            <EmptyState title="Belum ada transaksi" subtitle="Tambahkan pemasukan atau pengeluaran pertama kamu." />
          )}
        </Box>

        <Box sx={{ height: 8 }} />
      </Stack>
    </PullToRefresh>
  )
}
