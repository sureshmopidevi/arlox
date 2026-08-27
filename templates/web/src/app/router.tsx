import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AdminLayout } from '@/layouts/AdminLayout'
import { LoginPage } from '@/pages/LoginPage'
import { HomePage } from '@/pages/HomePage'

export function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route element={<AdminLayout />}>
          <Route path="/" element={<HomePage />} />
          {/* Wrap routes that require auth: <Route element={<ProtectedRoute />}>...</Route> */}
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
