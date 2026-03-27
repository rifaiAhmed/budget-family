import { Paper, BottomNavigation, BottomNavigationAction } from '@mui/material'
import DashboardIcon from '@mui/icons-material/Dashboard'
import ListAltIcon from '@mui/icons-material/ListAlt'
import AddCircleIcon from '@mui/icons-material/AddCircle'
import PieChartIcon from '@mui/icons-material/PieChart'
import SettingsIcon from '@mui/icons-material/Settings'
import { useLocation, useNavigate } from 'react-router-dom'
import { useMemo } from 'react'
import type { SyntheticEvent } from 'react'

const tabs = [
  { label: 'Home', value: '/dashboard', icon: <DashboardIcon /> },
  { label: 'Transaksi', value: '/transactions', icon: <ListAltIcon /> },
  { label: 'Add', value: '/add', icon: <AddCircleIcon sx={{ fontSize: 34 }} /> },
  { label: 'Budget', value: '/budget', icon: <PieChartIcon /> },
  { label: 'Settings', value: '/settings', icon: <SettingsIcon /> }
]

export default function BottomNavbar() {
  const nav = useNavigate()
  const location = useLocation()

  const value = useMemo(() => {
    const found = tabs.find((t) => location.pathname.startsWith(t.value))
    return found?.value ?? '/dashboard'
  }, [location.pathname])

  return (
    <Paper
      elevation={10}
      sx={{
        position: 'fixed',
        bottom: 12,
        left: '50%',
        transform: 'translateX(-50%)',
        width: 'calc(100% - 24px)',
        maxWidth: 480,
        borderTopLeftRadius: 18,
        borderTopRightRadius: 18,
        overflow: 'hidden'
      }}
    >
      <BottomNavigation
        showLabels
        value={value}
        onChange={(_event: SyntheticEvent, next: string) => nav(next)}
        sx={{
          '& .MuiBottomNavigationAction-root': {
            minHeight: 60,
            py: 0.75,
            px: 0.25,
            flex: 1,
            minWidth: 0
          },
          '& .MuiBottomNavigationAction-iconOnly': { py: 0.75 },
          '& .MuiBottomNavigationAction-label': {
            fontSize: 10.5,
            fontWeight: 800,
            lineHeight: 1.05,
            textAlign: 'center',
            display: '-webkit-box',
            WebkitBoxOrient: 'vertical',
            WebkitLineClamp: 2,
            overflow: 'hidden',
            maxWidth: '100%'
          },
          '& .MuiSvgIcon-root': { fontSize: 22 },
          '& .Mui-selected .MuiSvgIcon-root': { fontSize: 23 }
        }}
      >
        {tabs.map((t) => (
          <BottomNavigationAction key={t.value} label={t.label} value={t.value} icon={t.icon} />
        ))}
      </BottomNavigation>
    </Paper>
  )
}
