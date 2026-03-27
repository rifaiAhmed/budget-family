import { Box } from '@mui/material'
import type { ReactNode } from 'react'

export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <Box
      className="mobile-container"
      sx={{
        minHeight: 'calc(100vh - 24px)',
        backgroundImage: "url(/family_budget_bg.png)",
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat',
        position: 'relative'
      }}
    >
      <Box
        sx={{
          position: 'absolute',
          inset: 0,
          background: 'linear-gradient(180deg, rgba(0,0,0,0.55), rgba(0,0,0,0.65))'
        }}
      />
      <Box sx={{ position: 'relative', minHeight: '100%', p: 2, display: 'flex', alignItems: 'center' }}>{children}</Box>
    </Box>
  )
}
