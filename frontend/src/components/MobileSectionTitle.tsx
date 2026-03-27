import { Stack, Typography } from '@mui/material'
import type { ReactNode } from 'react'

export default function MobileSectionTitle({ title, right }: { title: string; right?: ReactNode }) {
  return (
    <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 1.5, mt: 1 }}>
      <Typography variant="subtitle1" sx={{ fontWeight: 900 }}>
        {title}
      </Typography>
      {right}
    </Stack>
  )
}
