import { Button, Card, CardContent, Divider, Stack, TextField, Typography } from '@mui/material'
import useAuth from '../../hooks/useAuth'
import { useSnackbar } from '../../components/SnackbarProvider'

export default function SettingsPage() {
  const { user, logout } = useAuth()
  const { notify } = useSnackbar()

  const familyId = localStorage.getItem('family_id') || ''

  return (
    <Stack spacing={2}>
      <Card variant="outlined">
        <CardContent>
          <Typography sx={{ fontWeight: 900, mb: 1 }}>Profile</Typography>
          <Typography variant="body2" sx={{ opacity: 0.7 }}>
            {user?.email}
          </Typography>
          <Typography sx={{ fontWeight: 900 }}>{user?.name}</Typography>
        </CardContent>
      </Card>

      <Card variant="outlined">
        <CardContent>
          <Typography sx={{ fontWeight: 900, mb: 1 }}>Family</Typography>
          <Typography variant="body2" sx={{ opacity: 0.7, mb: 1 }}>
            Enter your Family ID once. We will remember it for next time.
          </Typography>
          <TextField
            label="Family ID"
            defaultValue={familyId}
            fullWidth
            disabled={!!familyId}
            onBlur={(e) => {
              if (familyId) return
              const v = e.target.value.trim()
              if (!v) return
              localStorage.setItem('family_id', v)
              notify('Saved', 'success')
            }}
          />
        </CardContent>
      </Card>

      <Card variant="outlined">
        <CardContent>
          <Typography sx={{ fontWeight: 900, mb: 1 }}>Family members</Typography>
          <Typography variant="body2" sx={{ opacity: 0.7 }}>
            Coming soon.
          </Typography>
        </CardContent>
      </Card>

      <Divider />

      <Button
        variant="contained"
        color="error"
        onClick={() => {
          logout()
          notify('Logged out', 'info')
        }}
      >
        Logout
      </Button>
    </Stack>
  )
}
