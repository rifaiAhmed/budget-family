import { Box, LinearProgress } from '@mui/material'
import { useRef, useState } from 'react'
import type { ReactNode, TouchEvent } from 'react'

export default function PullToRefresh({ onRefresh, children }: { onRefresh: () => Promise<void>; children: ReactNode }) {
  const startY = useRef<number | null>(null)
  const [pull, setPull] = useState(0)
  const [refreshing, setRefreshing] = useState(false)

  const threshold = 70

  const handleTouchStart = (e: TouchEvent) => {
    if (window.scrollY > 0) return
    startY.current = e.touches[0].clientY
  }

  const handleTouchMove = (e: TouchEvent) => {
    if (startY.current == null) return
    const dy = e.touches[0].clientY - startY.current
    if (dy <= 0) return
    setPull(Math.min(100, dy))
  }

  const handleTouchEnd = async () => {
    if (startY.current == null) return
    startY.current = null
    if (pull >= threshold && !refreshing) {
      setRefreshing(true)
      try {
        await onRefresh()
      } finally {
        setRefreshing(false)
      }
    }
    setPull(0)
  }

  return (
    <Box onTouchStart={handleTouchStart} onTouchMove={handleTouchMove} onTouchEnd={handleTouchEnd}>
      {(pull > 0 || refreshing) && (
        <Box sx={{ mb: 1 }}>
          <LinearProgress variant={refreshing ? 'indeterminate' : 'determinate'} value={pull} />
        </Box>
      )}
      {children}
    </Box>
  )
}
