import { createContext, useContext, useMemo, useState } from 'react'
import type { ReactNode } from 'react'
import { createTheme } from '@mui/material/styles'

type ColorModeContextValue = {
  mode: 'light' | 'dark'
  toggle: () => void
  theme: ReturnType<typeof createTheme>
}

const ColorModeContext = createContext<ColorModeContextValue | undefined>(undefined)

export function ColorModeProvider({ children }: { children: ReactNode }) {
  const stored = localStorage.getItem('color_mode')
  const [mode, setMode] = useState<'light' | 'dark'>((stored === 'dark' || stored === 'light') ? stored : 'light')

  const toggle = () => {
    setMode((m) => {
      const next = m === 'light' ? 'dark' : 'light'
      localStorage.setItem('color_mode', next)
      return next
    })
  }

  const theme = useMemo(() => {
    return createTheme({
      palette: {
        mode,
        primary: { main: '#3b82f6' },
        secondary: { main: '#22c55e' }
      },
      shape: { borderRadius: 16 },
      typography: {
        fontFamily: 'Inter, system-ui, -apple-system, Segoe UI, Roboto, Arial, sans-serif'
      },
      components: {
        MuiButton: {
          styleOverrides: {
            root: {
              minHeight: 48,
              borderRadius: 14,
              textTransform: 'none',
              fontWeight: 700
            }
          }
        },
        MuiCard: {
          styleOverrides: {
            root: {
              borderRadius: 18
            }
          }
        }
      }
    })
  }, [mode])

  const value = useMemo(() => ({ mode, toggle, theme }), [mode, theme])

  return <ColorModeContext.Provider value={value}>{children}</ColorModeContext.Provider>
}

export function useColorMode() {
  const ctx = useContext(ColorModeContext)
  if (!ctx) throw new Error('useColorMode must be used within ColorModeProvider')
  return ctx
}
