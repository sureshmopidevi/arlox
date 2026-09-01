import { NavLink, useNavigate } from 'react-router-dom'
import { AppShell, Button, NavLink as MantineNavLink, Stack, Text } from '@mantine/core'
import { envConfig } from '@/config/envConfig'
import { useAuthStore } from '@/stores/authStore'

const navItems = [{ label: 'Home', to: '/' }]

export function AdminSidebar() {
  const token = useAuthStore((s) => s.token)
  const clearAuth = useAuthStore((s) => s.clearAuth)
  const navigate = useNavigate()

  return (
    <AppShell.Navbar p="md" w={224}>
      <Text fw={600} size="sm" mb="md">
        {envConfig.appName}
      </Text>
      <Stack gap={4} style={{ flex: 1 }}>
        {navItems.map(({ label, to }) => (
          <MantineNavLink key={to} component={NavLink} to={to} label={label} />
        ))}
      </Stack>
      {token ? (
        <Button variant="subtle" fullWidth onClick={() => { clearAuth(); navigate('/', { replace: true }) }}>
          Sign out
        </Button>
      ) : (
        <Button component={NavLink} to="/login" variant="light" fullWidth>
          Sign in
        </Button>
      )}
    </AppShell.Navbar>
  )
}
