import { Navigate, Route, Routes } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import MobileLayout from '../layout/MobileLayout'

import LoginPage from '../features/auth/LoginPage'
import RegisterPage from '../features/auth/RegisterPage'
import DashboardPage from '../features/dashboard/DashboardPage'
import TransactionListPage from '../features/transactions/TransactionListPage'
import AddTransactionPage from '../features/transactions/AddTransactionPage'
import TransactionDetailPage from '../features/transactions/TransactionDetailPage'
import BudgetPage from '../features/budget/BudgetPage'
import WalletPage from '../features/wallet/WalletPage'
import SettingsPage from '../features/settings/SettingsPage'

function Protected({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, loading } = useAuth()
  if (loading) return null
  if (!isAuthenticated) return <Navigate to="/login" replace />
  return <>{children}</>
}

export default function AppRouter() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />

      <Route
        path="/"
        element={
          <Protected>
            <MobileLayout />
          </Protected>
        }
      >
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<DashboardPage />} />
        <Route path="transactions" element={<TransactionListPage />} />
        <Route path="transactions/:id" element={<TransactionDetailPage />} />
        <Route path="add" element={<AddTransactionPage />} />
        <Route path="budget" element={<BudgetPage />} />
        <Route path="wallets" element={<WalletPage />} />
        <Route path="settings" element={<SettingsPage />} />
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
