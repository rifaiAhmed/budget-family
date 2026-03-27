import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { Button, Card, CardContent, Link, Stack, TextField, Typography } from '@mui/material'
import { Link as RouterLink, useNavigate } from 'react-router-dom'
import useAuth from '../../hooks/useAuth'
import { useSnackbar } from '../../components/SnackbarProvider'
import AuthLayout from '../../layout/AuthLayout'

const schema = z
  .object({
    name: z.string().min(2),
    email: z.string().email(),
    phone: z.string().optional(),
    password: z.string().min(8),
    confirmPassword: z.string().min(8)
  })
  .refine((v) => v.password === v.confirmPassword, { path: ['confirmPassword'], message: 'Passwords do not match' })

type FormValues = z.infer<typeof schema>

export default function RegisterPage() {
  const nav = useNavigate()
  const { register: doRegister } = useAuth()
  const { notify } = useSnackbar()

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting }
  } = useForm<FormValues>({ resolver: zodResolver(schema) })

  const onSubmit = async (values: FormValues) => {
    try {
      await doRegister(values.name, values.email, values.phone, values.password)
      notify('Account created', 'success')
      nav('/dashboard')
    } catch (e: any) {
      notify(e?.response?.data?.message || 'Register failed', 'error')
    }
  }

  return (
    <AuthLayout>
      <Card sx={{ width: '100%', bgcolor: 'background.paper' }}>
        <CardContent>
          <Stack spacing={2}>
            <Typography variant="h5" sx={{ fontWeight: 900 }}>
              Register
            </Typography>

            <TextField label="Name" fullWidth {...register('name')} error={!!errors.name} helperText={errors.name?.message} />
            <TextField
              label="Email"
              type="email"
              fullWidth
              {...register('email')}
              error={!!errors.email}
              helperText={errors.email?.message}
            />
            <TextField label="Phone (optional)" fullWidth {...register('phone')} error={!!errors.phone} helperText={errors.phone?.message} />
            <TextField
              label="Password"
              type="password"
              fullWidth
              {...register('password')}
              error={!!errors.password}
              helperText={errors.password?.message}
            />
            <TextField
              label="Confirm Password"
              type="password"
              fullWidth
              {...register('confirmPassword')}
              error={!!errors.confirmPassword}
              helperText={errors.confirmPassword?.message}
            />

            <Button variant="contained" size="large" onClick={handleSubmit(onSubmit)} disabled={isSubmitting}>
              Create account
            </Button>

            <Stack direction="row" justifyContent="space-between">
              <Link component={RouterLink} to="/login" underline="hover">
                Back to login
              </Link>
            </Stack>
          </Stack>
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
