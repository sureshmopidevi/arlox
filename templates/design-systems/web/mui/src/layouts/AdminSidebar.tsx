import { NavLink, useNavigate } from 'react-router-dom'
import { Box, Button, List, ListItemButton, ListItemText, Typography } from '@mui/material'
import { envConfig } from '@/config/envConfig'
import { useAuthStore } from '@/stores/authStore'

const navItems = [{ label: 'Home', to: '/' }]

export function AdminSidebar() {
  const token = useAuthStore((s) => s.token)
  const clearAuth = useAuthStore((s) => s.clearAuth)
  const navigate = useNavigate()

  return (
    <Box
      component="aside"
      sx={{ width: 224, flexShrink: 0, borderRight: 1, borderColor: 'divider', display: 'flex', flexDirection: 'column', bgcolor: 'background.paper' }}
    >
      <Box sx={{ px: 2, py: 2, borderBottom: 1, borderColor: 'divider' }}>
        <Typography variant="subtitle2" fontWeight={600}>{envConfig.appName}</Typography>
      </Box>
      <List sx={{ flex: 1 }}>
        {navItems.map(({ label, to }) => (
          <ListItemButton key={to} component={NavLink} to={to}>
            <ListItemText primary={label} />
          </ListItemButton>
        ))}
      </List>
      <Box sx={{ p: 1.5, borderTop: 1, borderColor: 'divider' }}>
        {token ? (
          <Button fullWidth onClick={() => { clearAuth(); navigate('/', { replace: true }) }}>Sign out</Button>
        ) : (
          <Button component={NavLink} to="/login" fullWidth>Sign in</Button>
        )}
      </Box>
    </Box>
  )
}
