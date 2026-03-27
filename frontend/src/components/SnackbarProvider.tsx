import { createContext, useContext, useMemo, useState } from 'react'
import { Alert, Snackbar } from '@mui/material'

type SnackbarState = {
  open: boolean
  message: string
  severity: 'success' | 'error' | 'info' | 'warning'
}

type SnackbarContextValue = {
  notify: (message: string, severity?: SnackbarState['severity']) => void
}

const SnackbarContext = createContext<SnackbarContextValue | undefined>(undefined)

export function SnackbarProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<SnackbarState>({ open: false, message: '', severity: 'info' })

  const notify = (message: string, severity: SnackbarState['severity'] = 'info') => {
    setState({ open: true, message, severity })
  }

  const value = useMemo(() => ({ notify }), [])

  return (
    <SnackbarContext.Provider value={value}>
      {children}
      <Snackbar
        open={state.open}
        autoHideDuration={2500}
        onClose={() => setState((s) => ({ ...s, open: false }))}
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      >
        <Alert severity={state.severity} variant="filled" sx={{ width: '100%' }}>
          {state.message}
        </Alert>
      </Snackbar>
    </SnackbarContext.Provider>
  )
}

export function useSnackbar() {
  const ctx = useContext(SnackbarContext)
  if (!ctx) throw new Error('useSnackbar must be used within SnackbarProvider')
  return ctx
}
