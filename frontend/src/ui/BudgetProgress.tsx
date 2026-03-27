import { Card, CardContent, LinearProgress, Stack, Typography } from '@mui/material'
import type { BudgetUsage } from '../types/budget'
import { formatMoneyIDR } from '../components/Money'

export default function BudgetProgress({
  row,
  categoryName,
  onClick
}: {
  row: BudgetUsage
  categoryName?: string
  onClick?: () => void
}) {
  const pct = Number(row.percentage_used)

  return (
    <Card variant="outlined" onClick={onClick} sx={onClick ? { cursor: 'pointer' } : undefined}>
      <CardContent>
        <Stack spacing={1}>
          <Stack direction="row" justifyContent="space-between">
            <Typography sx={{ fontWeight: 900 }} noWrap>
              {categoryName || 'Category'}
            </Typography>
            <Typography sx={{ fontWeight: 900 }}>{pct.toFixed(0)}%</Typography>
          </Stack>
          <Typography variant="body2" sx={{ opacity: 0.7 }}>
            Budget {formatMoneyIDR(row.budget_amount)} · Used {formatMoneyIDR(row.used_amount)}
          </Typography>
          <LinearProgress variant="determinate" value={Math.min(100, Math.max(0, pct))} sx={{ height: 10, borderRadius: 99 }} />
        </Stack>
      </CardContent>
    </Card>
  )
}
