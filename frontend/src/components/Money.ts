export function formatMoneyIDR(value: string | number) {
  const n = typeof value === 'string' ? Number(value) : value
  const safe = Number.isFinite(n) ? n : 0
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(safe)
}
