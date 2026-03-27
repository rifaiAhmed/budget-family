import { ThemeProvider } from '@mui/material/styles'
import { SnackbarProvider } from './components/SnackbarProvider'
import { AuthProvider } from './context/AuthContext'
import { ColorModeProvider, useColorMode } from './context/ColorModeContext'
import AppRouter from './routes/AppRouter'

function InnerApp() {
  const { theme } = useColorMode()
  return (
    <ThemeProvider theme={theme}>
      <SnackbarProvider>
        <AuthProvider>
          <AppRouter />
        </AuthProvider>
      </SnackbarProvider>
    </ThemeProvider>
  )
}

export default function App() {
  return (
    <ColorModeProvider>
      <InnerApp />
    </ColorModeProvider>
  )
}
