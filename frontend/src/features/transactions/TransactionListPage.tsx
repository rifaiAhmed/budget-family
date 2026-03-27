import { Box, Chip, MenuItem, Skeleton, Stack, TextField } from '@mui/material'
import { useMemo, useState } from 'react'
import TransactionItem from '../../ui/TransactionItem'
import EmptyState from '../../components/EmptyState'
import PullToRefresh from '../../components/PullToRefresh'
import MobileSectionTitle from '../../components/MobileSectionTitle'
import { useTransactions } from '../../hooks/useTransactions'

function getFamilyId(): string {
  return localStorage.getItem('family_id') || ''
}

function getMonthOptions() {
  const fmt = new Intl.DateTimeFormat('en-US', { month: 'long' })
  return Array.from({ length: 12 }).map((_, i) => ({ label: fmt.format(new Date(2020, i, 1)), value: i + 1 }))
}

function monthToRange(month: number, year: number) {
  const from = new Date(year, month - 1, 1)
  const to = new Date(year, month, 0)
  const fmt = (d: Date) => d.toISOString().slice(0, 10)
  return { from: fmt(from), to: fmt(to) }
}

export default function TransactionListPage() {
  const familyId = getFamilyId()
  const now = new Date()

  const [month, setMonth] = useState<number>(now.getMonth() + 1)
  const [year, setYear] = useState<number>(now.getFullYear())
  const [type, setType] = useState<'income' | 'expense' | ''>('')

  const range = useMemo(() => monthToRange(month, year), [month, year])

  const q = useTransactions({
    family_id: familyId,
    page: 1,
    limit: 50,
    type: type || undefined,
    from: range.from,
    to: range.to
  })

  const onRefresh = async () => {
    await q.refetch()
  }

  if (!familyId) {
    return <EmptyState title="No family selected" subtitle="Go to Settings and set your family_id." />
  }

  return (
    <PullToRefresh onRefresh={onRefresh}>
      <Stack spacing={2}>
        <MobileSectionTitle title="Filters" />

        <Stack direction="row" spacing={1.5}>
          <TextField select fullWidth label="Month" value={month} onChange={(e) => setMonth(Number(e.target.value))}>
            {getMonthOptions().map((m) => (
              <MenuItem key={m.value} value={m.value}>
                {m.label}
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

        <Stack direction="row" spacing={1}>
          <Chip label="All" color={type === '' ? 'primary' : 'default'} onClick={() => setType('')} />
          <Chip label="Income" color={type === 'income' ? 'primary' : 'default'} onClick={() => setType('income')} />
          <Chip label="Expense" color={type === 'expense' ? 'primary' : 'default'} onClick={() => setType('expense')} />
          <Box sx={{ flex: 1 }} />
        </Stack>

        <MobileSectionTitle title="Transactions" />

        {q.isLoading ? (
          <Stack spacing={1.5}>
            <Skeleton variant="rounded" height={78} />
            <Skeleton variant="rounded" height={78} />
            <Skeleton variant="rounded" height={78} />
          </Stack>
        ) : q.data?.items?.length ? (
          <Stack spacing={1.5}>
            {q.data.items.map((t) => (
              <TransactionItem key={t.id} tx={t} />
            ))}
          </Stack>
        ) : (
          <EmptyState title="No transactions" subtitle="Try changing filters or add a new transaction." />
        )}
      </Stack>
    </PullToRefresh>
  )
}
