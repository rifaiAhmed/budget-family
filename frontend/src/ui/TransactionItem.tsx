import { Card, CardActionArea, CardContent, Stack, Typography } from '@mui/material'
import { useNavigate } from 'react-router-dom'
import type { Transaction } from '../types/transaction'
import { formatMoneyIDR } from '../components/Money'

export default function TransactionItem({ tx }: { tx: Transaction }) {
  const nav = useNavigate()
  const positive = tx.type === 'income'

  const dateLabel = (() => {
    const d = new Date(tx.transaction_date)
    if (Number.isNaN(d.getTime())) return tx.transaction_date
    return new Intl.DateTimeFormat('en-US', { day: '2-digit', month: 'short', year: 'numeric' }).format(d)
  })()

  return (
    <Card variant="outlined">
      <CardActionArea onClick={() => nav(`/transactions/${tx.id}`)}>
        <CardContent>
          <Stack direction="row" spacing={1.5} alignItems="center">
            <div
              style={{
                width: 38,
                height: 38,
                borderRadius: 12,
                background: positive ? 'rgba(34,197,94,0.15)' : 'rgba(239,68,68,0.15)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontWeight: 900
              }}
            >
              {positive ? '+' : '-'}
            </div>
            <div style={{ flex: 1, minWidth: 0 }}>
              <Typography sx={{ fontWeight: 900 }} noWrap>
                {tx.note || (positive ? 'Income' : 'Expense')}
              </Typography>
              <Typography variant="body2" sx={{ opacity: 0.7 }}>
                {dateLabel}
              </Typography>
            </div>
            <Typography sx={{ fontWeight: 900, color: positive ? 'success.main' : 'error.main' }}>
              {positive ? '' : '-'}
              {formatMoneyIDR(tx.amount)}
            </Typography>
          </Stack>
        </CardContent>
      </CardActionArea>
    </Card>
  )
}
