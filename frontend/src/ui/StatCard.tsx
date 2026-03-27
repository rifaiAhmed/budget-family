import { Card, CardContent, Stack, Typography } from '@mui/material'
import type { ReactNode } from 'react'

export default function StatCard({ title, amount, icon }: { title: string; amount: string; icon: ReactNode }) {
  return (
    <Card variant="outlined">
      <CardContent>
        <Stack direction="row" spacing={1.5} alignItems="center">
          <div style={{ width: 34, display: 'flex', justifyContent: 'center' }}>{icon}</div>
          <div style={{ flex: 1, minWidth: 0 }}>
            <Typography variant="body2" sx={{ opacity: 0.7, fontWeight: 700 }}>
              {title}
            </Typography>
            <Typography
              variant="h6"
              sx={{ fontWeight: 900, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', fontSize: 18, lineHeight: 1.15 }}
            >
              {amount}
            </Typography>
          </div>
        </Stack>
      </CardContent>
    </Card>
  )
}
