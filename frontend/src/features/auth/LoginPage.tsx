import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { Button, Card, CardContent, Link, Stack, TextField, Typography } from '@mui/material'
import { Link as RouterLink, useNavigate } from 'react-router-dom'
import useAuth from '../../hooks/useAuth'
import { useSnackbar } from '../../components/SnackbarProvider'
import AuthLayout from '../../layout/AuthLayout'

const schema = z.object({
  email: z.string().email(),
  password: z.string().min(8)
})

type FormValues = z.infer<typeof schema>

export default function LoginPage() {
  const nav = useNavigate()
  const { login } = useAuth()
  const { notify } = useSnackbar()

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting }
  } = useForm<FormValues>({ resolver: zodResolver(schema) })

  const onSubmit = async (values: FormValues) => {
    try {
      await login(values.email, values.password)
      notify('Welcome back', 'success')
      nav('/dashboard')
    } catch (e: any) {
      notify(e?.response?.data?.message || 'Login failed', 'error')
    }
  }

  return (
    <AuthLayout>
      <Card sx={{ width: '100%', bgcolor: 'background.paper' }}>
        <CardContent>
          <Stack spacing={2}>
            <Typography variant="h5" sx={{ fontWeight: 900 }}>
              Login
            </Typography>
            <Typography variant="body2" sx={{ opacity: 0.7 }}>
              Sign in to manage your family budget.
            </Typography>

            <TextField
              label="Email"
              type="email"
              fullWidth
              {...register('email')}
              error={!!errors.email}
              helperText={errors.email?.message}
            />
            <TextField
              label="Password"
              type="password"
              fullWidth
              {...register('password')}
              error={!!errors.password}
              helperText={errors.password?.message}
            />

            <Button variant="contained" size="large" onClick={handleSubmit(onSubmit)} disabled={isSubmitting}>
              Login
            </Button>

            <Stack direction="row" justifyContent="space-between">
              <Link component={RouterLink} to="/register" underline="hover">
                Register
              </Link>
              <Link href="#" underline="hover" onClick={(e) => e.preventDefault()}>
                Forgot password
              </Link>
            </Stack>
          </Stack>
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
