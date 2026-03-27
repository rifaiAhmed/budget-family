import { Button, Card, CardContent, Divider, MenuItem, Stack, TextField, Typography } from '@mui/material'
import { useEffect, useMemo, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import useAuth from '../../hooks/useAuth'
import { useSnackbar } from '../../components/SnackbarProvider'
import { useFamilies, useFamilyMembers } from '../../hooks/useFamilies'
import * as familyApi from '../../api/familyApi'

function isUUID(v: string) {
  return /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(v)
}

export default function SettingsPage() {
  const { user, logout } = useAuth()
  const { notify } = useSnackbar()
  const qc = useQueryClient()

  const [familyId, setFamilyId] = useState(localStorage.getItem('family_id') || '')
  const [familyInput, setFamilyInput] = useState('')
  const [familyNameInput, setFamilyNameInput] = useState('')
  const familiesQ = useFamilies()

  const families = familiesQ.data?.items || []
  const isMemberOfAnyFamily = families.length > 0

  const activeFamily = useMemo(() => {
    if (!isMemberOfAnyFamily) return null
    if (!familyId) return null
    return families.find((f) => f.id === familyId) || null
  }, [families, familyId, isMemberOfAnyFamily])

  useEffect(() => {
    if (!isMemberOfAnyFamily) return
    if (!familyId) return
    const exists = families.some((f) => f.id === familyId)
    if (exists) return
    localStorage.removeItem('family_id')
    setFamilyId('')
  }, [activeFamily?.id, familyId, isMemberOfAnyFamily])

  const [selectedFamilyId, setSelectedFamilyId] = useState('')

  useEffect(() => {
    if (!isMemberOfAnyFamily) return
    if (familyId) return
    if (selectedFamilyId) return
    if (!families[0]?.id) return
    setSelectedFamilyId(families[0].id)
  }, [families, familyId, isMemberOfAnyFamily, selectedFamilyId])

  const membersQ = useFamilyMembers(activeFamily?.id || '')

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
          {isMemberOfAnyFamily ? (
            <>
              <Typography variant="body2" sx={{ opacity: 0.7, mb: 1 }}>
                Kamu sudah terdaftar di keluarga.
              </Typography>

              {familyId ? (
                <TextField label="Keluarga" fullWidth value={activeFamily?.name || ''} InputProps={{ readOnly: true }} />
              ) : (
                <Stack spacing={1.25}>
                  <Typography variant="body2" sx={{ opacity: 0.7 }}>
                    Pilih keluarga untuk mulai melihat data.
                  </Typography>
                  <TextField
                    select
                    label="Keluarga"
                    fullWidth
                    value={selectedFamilyId}
                    onChange={(e) => setSelectedFamilyId(e.target.value)}
                  >
                    {families.map((f) => (
                      <MenuItem key={f.id} value={f.id}>
                        {f.name}
                      </MenuItem>
                    ))}
                  </TextField>
                  <Button
                    variant="contained"
                    disabled={!selectedFamilyId}
                    onClick={() => {
                      localStorage.setItem('family_id', selectedFamilyId)
                      setFamilyId(selectedFamilyId)
                      notify('Saved', 'success')
                    }}
                  >
                    Gunakan keluarga ini
                  </Button>
                </Stack>
              )}

              {!!familyId && (
                <>
                  <Divider sx={{ my: 2 }} />

                  <Typography sx={{ fontWeight: 900, mb: 1 }}>Anggota keluarga</Typography>
                  {membersQ.isLoading ? (
                    <Typography variant="body2" sx={{ opacity: 0.7 }}>
                      Loading...
                    </Typography>
                  ) : membersQ.data?.items?.length ? (
                    <Stack spacing={1}>
                      {membersQ.data.items.map((m) => (
                        <Card key={m.id} variant="outlined">
                          <CardContent sx={{ py: 1.25, '&:last-child': { pb: 1.25 } }}>
                            <Stack direction="row" alignItems="center" justifyContent="space-between" spacing={1.5}>
                              <div style={{ flex: 1, minWidth: 0 }}>
                                <Typography sx={{ fontWeight: 900 }} noWrap>
                                  {m.name}
                                </Typography>
                                <Typography variant="body2" sx={{ opacity: 0.7 }} noWrap>
                                  {m.email}
                                </Typography>
                              </div>
                              <Typography sx={{ fontWeight: 900, opacity: 0.9 }}>
                                {m.role === 'owner' ? 'Kepala Keluarga' : 'Anggota'}
                              </Typography>
                            </Stack>
                          </CardContent>
                        </Card>
                      ))}
                    </Stack>
                  ) : (
                    <Typography variant="body2" sx={{ opacity: 0.7 }}>
                      Tidak ada data anggota.
                    </Typography>
                  )}
                </>
              )}
            </>
          ) : (
            <>
              <Typography variant="body2" sx={{ opacity: 0.7, mb: 1.5 }}>
                Kamu belum terdaftar di keluarga manapun.
              </Typography>

              <Stack spacing={2}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography sx={{ fontWeight: 900, mb: 1 }}>Buat keluarga baru</Typography>
                    <Stack spacing={1.25}>
                      <TextField
                        label="Nama keluarga"
                        value={familyNameInput}
                        fullWidth
                        onChange={(e) => setFamilyNameInput(e.target.value)}
                      />
                      <Button
                        variant="contained"
                        disabled={!familyNameInput.trim()}
                        onClick={async () => {
                          try {
                            const name = familyNameInput.trim()
                            if (!name) return
                            const res = await familyApi.createFamily(name)
                            const createdId = res.family.id
                            localStorage.setItem('family_id', createdId)
                            setFamilyId(createdId)
                            setFamilyNameInput('')
                            await qc.invalidateQueries({ queryKey: ['families'] })
                            await qc.invalidateQueries({ queryKey: ['family-members'] })
                            await familiesQ.refetch()
                            notify('Saved', 'success')
                          } catch (e: any) {
                            notify(e?.response?.data?.message || 'Gagal membuat keluarga', 'error')
                          }
                        }}
                      >
                        Buat keluarga
                      </Button>
                    </Stack>
                  </CardContent>
                </Card>

                <Card variant="outlined">
                  <CardContent>
                    <Typography sx={{ fontWeight: 900, mb: 1 }}>Gabung dengan Family UUID</Typography>
                    <Stack spacing={1.25}>
                      <TextField
                        label="Family UUID"
                        value={familyInput}
                        fullWidth
                        onChange={(e) => setFamilyInput(e.target.value)}
                        error={!!familyInput && !isUUID(familyInput.trim())}
                        helperText={familyInput && !isUUID(familyInput.trim()) ? 'Format UUID tidak valid' : 'Wajib diisi'}
                      />
                      <Button
                        variant="contained"
                        disabled={!isUUID(familyInput.trim())}
                        onClick={async () => {
                          try {
                            const v = familyInput.trim()
                            if (!isUUID(v)) return
                            const res = await familyApi.joinFamily(v)
                            const joinedId = res.family.id
                            localStorage.setItem('family_id', joinedId)
                            setFamilyId(joinedId)
                            setFamilyInput('')
                            await qc.invalidateQueries({ queryKey: ['families'] })
                            await qc.invalidateQueries({ queryKey: ['family-members'] })
                            await familiesQ.refetch()
                            notify('Saved', 'success')
                          } catch (e: any) {
                            notify(e?.response?.data?.message || 'Gagal bergabung ke keluarga', 'error')
                          }
                        }}
                      >
                        Simpan
                      </Button>
                    </Stack>
                  </CardContent>
                </Card>
              </Stack>
            </>
          )}
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
