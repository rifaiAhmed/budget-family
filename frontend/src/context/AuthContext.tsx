import { createContext, useContext, useEffect, useMemo, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import * as authApi from '../api/authApi'
import type { User } from '../types/user'

type AuthContextValue = {
  user: User | null
  isAuthenticated: boolean
  loading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (name: string, email: string, phone: string | undefined, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const qc = useQueryClient()

  const isAuthenticated = !!localStorage.getItem('access_token')

  useEffect(() => {
    const boot = async () => {
      try {
        if (!localStorage.getItem('access_token')) {
          setLoading(false)
          return
        }
        const me = await authApi.me()
        setUser(me)
      } catch {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        setUser(null)
      } finally {
        setLoading(false)
      }
    }
    void boot()
  }, [])

  const login = async (email: string, password: string) => {
    const res = await authApi.login({ email, password })
    localStorage.setItem('access_token', res.tokens.access_token)
    localStorage.setItem('refresh_token', res.tokens.refresh_token)
    setUser(res.user)
    await qc.invalidateQueries()
  }

  const register = async (name: string, email: string, phone: string | undefined, password: string) => {
    const res = await authApi.register({ name, email, phone, password })
    localStorage.setItem('access_token', res.tokens.access_token)
    localStorage.setItem('refresh_token', res.tokens.refresh_token)
    setUser(res.user)
    await qc.invalidateQueries()
  }

  const logout = () => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    setUser(null)
    qc.clear()
  }

  const value = useMemo(
    () => ({ user, isAuthenticated: !!user && isAuthenticated, loading, login, register, logout }),
    [user, isAuthenticated, loading]
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
