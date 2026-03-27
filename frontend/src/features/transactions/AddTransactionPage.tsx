import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { Button, MenuItem, Skeleton, Stack, TextField } from '@mui/material'
import { useMemo, useState } from 'react'
import type { ChangeEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { useCreateTransaction } from '../../hooks/useTransactions'
import { useWallets } from '../../hooks/useWallets'
import { useCategories } from '../../hooks/useCategories'
import { useSnackbar } from '../../components/SnackbarProvider'
import EmptyState from '../../components/EmptyState'

function getFamilyId(): string {
  return localStorage.getItem('family_id') || (import.meta.env.VITE_DEFAULT_FAMILY_ID as string) || ''
}

const schema = z.object({
  type: z.enum(['income', 'expense']),
  wallet_id: z.string().uuid(),
  category_id: z.string().uuid(),
  amount: z.string().min(1),
  transaction_date: z.string().min(10),
  note: z.string().optional()
})

type FormValues = z.infer<typeof schema>

export default function AddTransactionPage() {
  const familyId = getFamilyId()
  const { notify } = useSnackbar()
  const nav = useNavigate()

  const createM = useCreateTransaction()

  const walletsQ = useWallets({ family_id: familyId, page: 1, limit: 100 })
  const [type, setType] = useState<'income' | 'expense'>(
    ((localStorage.getItem('tx_type') as 'income' | 'expense' | null) || 'income')
  )

  const categoriesQ = useCategories({ family_id: familyId, type, page: 1, limit: 200 })

  const today = useMemo(() => new Date().toISOString().slice(0, 10), [])

  const {
    register,
    handleSubmit,
    setValue,
    reset,
    formState: { errors }
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      type: 'income',
      transaction_date: today
    }
  })

  if (!familyId) {
    return <EmptyState title="No family selected" subtitle="Go to Settings and set your family_id." />
  }

  const onSubmit = async (v: FormValues) => {
    try {
      const created = await createM.mutateAsync({
        family_id: familyId,
        wallet_id: v.wallet_id,
        category_id: v.category_id,
        amount: v.amount,
        type: v.type,
        note: v.note,
        transaction_date: v.transaction_date
      })
      notify('Transaction saved', 'success')

	  // Reset form for next entry
	  localStorage.setItem('tx_type', 'income')
	  setType('income')
	  reset({
		type: 'income',
		transaction_date: today,
		wallet_id: '',
		category_id: '',
		amount: '',
		note: ''
	  } as any)
	  nav(`/transactions/${created.id}`)
    } catch (e: any) {
      notify(e?.response?.data?.message || 'Failed to save', 'error')
    }
  }

  const loading = walletsQ.isLoading || categoriesQ.isLoading

  return (
    <Stack spacing={2}>
      <TextField
        select
        label="Type"
        value={type}
        {...register('type')}
        onChange={(e: ChangeEvent<HTMLInputElement>) => {
          const t = e.target.value as 'income' | 'expense'
          setValue('type', t)
          setType(t)
          localStorage.setItem('tx_type', t)
        }}
      >
        <MenuItem value="income">Income</MenuItem>
        <MenuItem value="expense">Expense</MenuItem>
      </TextField>

      {loading ? (
        <Skeleton variant="rounded" height={56} />
      ) : (
        <TextField select label="Wallet" {...register('wallet_id')} error={!!errors.wallet_id} helperText={errors.wallet_id?.message}>
          {walletsQ.data?.items?.map((w) => (
            <MenuItem key={w.id} value={w.id}>
              {w.name}
            </MenuItem>
          ))}
        </TextField>
      )}

      {loading ? (
        <Skeleton variant="rounded" height={56} />
      ) : (
        <TextField
          select
          label="Category"
          {...register('category_id')}
          error={!!errors.category_id}
          helperText={errors.category_id?.message}
        >
          {categoriesQ.data?.items?.map((cat) => (
            <MenuItem key={cat.id} value={cat.id}>
              {cat.name}
            </MenuItem>
          ))}
        </TextField>
      )}

      <TextField label="Amount" inputMode="decimal" {...register('amount')} error={!!errors.amount} helperText={errors.amount?.message} />
      <TextField
        label="Date"
        type="date"
        InputLabelProps={{ shrink: true }}
        {...register('transaction_date')}
        error={!!errors.transaction_date}
        helperText={errors.transaction_date?.message}
      />
      <TextField label="Note" {...register('note')} />

      <Button variant="contained" size="large" onClick={handleSubmit(onSubmit)} disabled={createM.isPending}>
        Save Transaction
      </Button>
    </Stack>
  )
}
