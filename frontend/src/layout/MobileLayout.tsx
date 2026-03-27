import { Box } from '@mui/material'
import { Outlet, useLocation } from 'react-router-dom'
import Header from './Header'
import BottomNavbar from './BottomNavbar'

const HIDE_CHROME_ROUTES = new Set(['/login', '/register'])

export default function MobileLayout() {
  const location = useLocation()
  const hide = HIDE_CHROME_ROUTES.has(location.pathname)

  return (
    <Box className="mobile-container" sx={{ bgcolor: 'background.paper', color: 'text.primary' }}>
      <Box sx={{ minHeight: '100%', display: 'flex', flexDirection: 'column' }}>
        {!hide && <Header />}

        <Box sx={{ flex: 1, p: 2, pb: hide ? 2 : 12 }}>
          <Outlet />
        </Box>

        {!hide && <BottomNavbar />}
      </Box>
    </Box>
  )
}
