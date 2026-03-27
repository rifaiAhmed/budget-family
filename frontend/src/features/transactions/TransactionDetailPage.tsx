import { Card, CardContent, Skeleton, Stack, Typography } from '@mui/material'
import { useMemo } from 'react'
import { useParams } from 'react-router-dom'
import EmptyState from '../../components/EmptyState'
import { formatMoneyIDR } from '../../components/Money'
import { useTransaction } from '../../hooks/useTransactions'

export default function TransactionDetailPage() {
  const { id } = useParams()
  const txId = useMemo(() => id || '', [id])

  const q = useTransaction(txId)

  if (!txId) return <EmptyState title="Transaction not found" />

  if (q.isLoading) {
    return (
      <Stack spacing={2}>
        <Skeleton variant="rounded" height={96} />
        <Skeleton variant="rounded" height={96} />
        <Skeleton variant="rounded" height={96} />
      </Stack>
    )
  }

  if (!q.data) return <EmptyState title="Transaction not found" />

  const tx = q.data
  const positive = tx.type === 'income'

  return (
    <Stack spacing={2}>
      <Typography variant="h6" sx={{ fontWeight: 900 }}>
        Transaction
      </Typography>

      <Card variant="outlined">
        <CardContent>
          <Typography variant="body2" sx={{ opacity: 0.7 }}>
            Amount
          </Typography>
          <Typography variant="h5" sx={{ fontWeight: 900, color: positive ? 'success.main' : 'error.main' }}>
            {positive ? '' : '-'}
            {formatMoneyIDR(tx.amount)}
          </Typography>
        </CardContent>
      </Card>

      <Card variant="outlined">
        <CardContent>
          <Stack spacing={1}>
            <div>
              <Typography variant="body2" sx={{ opacity: 0.7 }}>
                Type
              </Typography>
              <Typography sx={{ fontWeight: 800 }}>{positive ? 'Income' : 'Expense'}</Typography>
            </div>
            <div>
              <Typography variant="body2" sx={{ opacity: 0.7 }}>
                Date
              </Typography>
              <Typography sx={{ fontWeight: 800 }}>{tx.transaction_date}</Typography>
            </div>
            {tx.note ? (
              <div>
                <Typography variant="body2" sx={{ opacity: 0.7 }}>
                  Note
                </Typography>
                <Typography sx={{ fontWeight: 800 }}>{tx.note}</Typography>
              </div>
            ) : null}
          </Stack>
        </CardContent>
      </Card>
    </Stack>
  )
}
