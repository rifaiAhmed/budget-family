import { AppBar, Box, IconButton, Toolbar, Typography } from '@mui/material'
import DarkModeIcon from '@mui/icons-material/DarkMode'
import LightModeIcon from '@mui/icons-material/LightMode'
import { useColorMode } from '../context/ColorModeContext'
import { useAuth } from '../context/AuthContext'

export default function Header() {
  const { mode, toggle } = useColorMode()
  const { user } = useAuth()

  return (
    <AppBar position="sticky" elevation={0} color="transparent" sx={{ backdropFilter: 'blur(10px)' }}>
      <Toolbar sx={{ px: 2 }}>
        <Box sx={{ flex: 1 }}>
          <Typography variant="subtitle2" sx={{ opacity: 0.75 }}>
            Hello
          </Typography>
          <Typography variant="h6" sx={{ fontWeight: 800, lineHeight: 1.1 }}>
            {user?.name || 'Budget'}
          </Typography>
        </Box>
        <IconButton onClick={toggle} aria-label="toggle dark mode" size="large">
          {mode === 'dark' ? <LightModeIcon /> : <DarkModeIcon />}
        </IconButton>
      </Toolbar>
    </AppBar>
  )
}
