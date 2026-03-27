import { Button, Card, CardContent, Dialog, DialogActions, DialogContent, DialogTitle, MenuItem, Skeleton, Stack, TextField, Typography } from '@mui/material'
import { useMemo, useState } from 'react'
import PullToRefresh from '../../components/PullToRefresh'
import EmptyState from '../../components/EmptyState'
import MobileSectionTitle from '../../components/MobileSectionTitle'
import { formatMoneyIDR } from '../../components/Money'
import { useCreateWallet, useDeleteWallet, useUpdateWallet, useWallets } from '../../hooks/useWallets'
import { useSnackbar } from '../../components/SnackbarProvider'

function getFamilyId(): string {
  return localStorage.getItem('family_id') || (import.meta.env.VITE_DEFAULT_FAMILY_ID as string) || ''
}

export default function WalletPage() {
  const familyId = getFamilyId()
  const { notify } = useSnackbar()

  const q = useWallets({ family_id: familyId, page: 1, limit: 100 })
  const createM = useCreateWallet()
  const updateM = useUpdateWallet()
  const deleteM = useDeleteWallet()

  const [open, setOpen] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)

  const [name, setName] = useState('')
  const [type, setType] = useState<'cash' | 'bank' | 'ewallet' | 'card'>('cash')
  const [balance, setBalance] = useState('0')

  const onRefresh = async () => {
    await q.refetch()
  }

  const openCreate = () => {
    setEditingId(null)
    setName('')
    setType('cash')
    setBalance('0')
    setOpen(true)
  }

  const openEdit = (w: any) => {
    setEditingId(w.id)
    setName(w.name)
    setType(w.type)
    setBalance(String(w.balance ?? '0'))
    setOpen(true)
  }

  const save = async () => {
    try {
      if (!familyId) return
      if (!name.trim()) {
        notify('Name is required', 'error')
        return
      }
      if (editingId) {
        await updateM.mutateAsync({ id: editingId, payload: { name: name.trim(), type } })
        notify('Wallet updated', 'success')
      } else {
        await createM.mutateAsync({ family_id: familyId, name: name.trim(), type, balance })
        notify('Wallet created', 'success')
      }
      setOpen(false)
    } catch (e: any) {
      notify(e?.response?.data?.message || 'Failed', 'error')
    }
  }

  const wallets = useMemo(() => q.data?.items ?? [], [q.data?.items])

  if (!familyId) {
    return <EmptyState title="No family selected" subtitle="Go to Settings and set your family_id." />
  }

  return (
    <PullToRefresh onRefresh={onRefresh}>
      <Stack spacing={2}>
        <MobileSectionTitle
          title="Wallets"
          right={
            <Button variant="contained" onClick={openCreate}>
              Add
            </Button>
          }
        />

        {q.isLoading ? (
          <Stack spacing={1.5}>
            <Skeleton variant="rounded" height={92} />
            <Skeleton variant="rounded" height={92} />
          </Stack>
        ) : wallets.length ? (
          <Stack spacing={1.5}>
            {wallets.map((w) => (
              <Card key={w.id} variant="outlined" onClick={() => openEdit(w)} sx={{ cursor: 'pointer' }}>
                <CardContent>
                  <Stack direction="row" justifyContent="space-between" alignItems="center">
                    <div>
                      <Typography sx={{ fontWeight: 900 }}>{w.name}</Typography>
                      <Typography variant="body2" sx={{ opacity: 0.7 }}>
                        {w.type}
                      </Typography>
                    </div>
                    <Typography sx={{ fontWeight: 900 }}>{formatMoneyIDR(w.balance)}</Typography>
                  </Stack>
                </CardContent>
              </Card>
            ))}
          </Stack>
        ) : (
          <EmptyState title="No wallets" subtitle="Create a wallet to start tracking." />
        )}

        <Dialog open={open} onClose={() => setOpen(false)} fullWidth>
          <DialogTitle sx={{ fontWeight: 900 }}>{editingId ? 'Edit Wallet' : 'Add Wallet'}</DialogTitle>
          <DialogContent>
            <Stack spacing={2} sx={{ pt: 1 }}>
              <TextField label="Name" value={name} onChange={(e) => setName(e.target.value)} />
              <TextField select label="Type" value={type} onChange={(e) => setType(e.target.value as any)}>
                <MenuItem value="cash">Cash</MenuItem>
                <MenuItem value="bank">Bank</MenuItem>
                <MenuItem value="ewallet">E-wallet</MenuItem>
                <MenuItem value="card">Card</MenuItem>
              </TextField>
              {!editingId && <TextField label="Initial Balance" value={balance} onChange={(e) => setBalance(e.target.value)} />}
            </Stack>
          </DialogContent>
          <DialogActions>
            {editingId && (
              <Button
                color="error"
                onClick={async () => {
                  await deleteM.mutateAsync(editingId)
                  notify('Wallet deleted', 'success')
                  setOpen(false)
                }}
              >
                Delete
              </Button>
            )}
            <Button onClick={() => setOpen(false)}>Cancel</Button>
            <Button variant="contained" onClick={save} disabled={createM.isPending || updateM.isPending}>
              Save
            </Button>
          </DialogActions>
        </Dialog>
      </Stack>
    </PullToRefresh>
  )
}
