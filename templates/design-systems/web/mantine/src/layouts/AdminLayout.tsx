import { AppShell } from '@mantine/core'
import { Outlet } from 'react-router-dom'
import { AdminSidebar } from './AdminSidebar'

export function AdminLayout() {
  return (
    <AppShell navbar={{ width: 224, breakpoint: 'sm' }} padding="md">
      <AdminSidebar />
      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  )
}
