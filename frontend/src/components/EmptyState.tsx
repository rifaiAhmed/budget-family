import { Box, Typography } from '@mui/material'

export default function EmptyState({ title, subtitle }: { title: string; subtitle?: string }) {
  return (
    <Box sx={{ py: 6, textAlign: 'center', opacity: 0.8 }}>
      <Typography variant="h6" sx={{ fontWeight: 900, mb: 1 }}>
        {title}
      </Typography>
      {subtitle && <Typography variant="body2">{subtitle}</Typography>}
    </Box>
  )
}
